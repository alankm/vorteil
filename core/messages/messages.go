package messages

import (
	"database/sql"
	"time"

	"github.com/alankm/vorteil/core/shared"
)

type Severity int

const (
	Debug Severity = iota
	Info
	Warning
	Error
	Critical
	Alert
	All
)

var (
	cleanupInterval = time.Second * 60
	cleanupAge      = time.Minute * 2
	severities      = []string{
		"debug",
		"info",
		"warning",
		"error",
		"critical",
		"alert",
		"all",
	}
)

type Messages struct {
	data      *sql.DB
	timer     *time.Timer
	inbox     chan *message
	outbox    map[Severity](chan *message)
	customers map[Severity](map[chan *message]bool)
}

type message struct {
	Sev  Severity          `json:"severity"`
	Time int64             `json:timestamp`
	Msg  string            `json:"message"`
	Code string            `json:"code"`
	Args map[string]string `json:"info"`
}

func New(functions *shared.Functions) (*Messages, error) {
	msgs := &Messages{
		data: functions.Database(),
	}
	err := msgs.init()
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func (m *Messages) init() error {
	_, err := m.data.Exec(
		"CREATE TABLE IF NOT EXISTS messages (" +
			"id INTEGER PRIMARY KEY, " +
			"time INTEGER NULL, " +
			"severity VARCHAR(32) NULL, " +
			"messagetext VARCHAR(128) NULL, " +
			"messagecode VARCHAR(128) NULL, " +
			"rulesowner VARCHAR(128) NULL, " +
			"rulesgroup VARCHAR(128) NULL, " +
			"rulesmode VARCHAR(4) NULL" +
			")",
	)
	if err != nil {
		return err
	}
	_, err = m.data.Exec(
		"CREATE TABLE IF NOT EXISTS args (" +
			"id INTEGER, " +
			"key VARCHAR(128) NULL, " +
			"value VARCHAR(128) NULL, " +
			"PRIMARY KEY (id, key), " +
			"FOREIGN KEY(id) REFERENCES messages(id) ON DELETE CASCADE" +
			")",
	)
	if err != nil {
		return err
	}
	m.inbox = make(chan *message)
	m.outbox = make(map[Severity](chan *message))
	m.customers = make(map[Severity](map[chan *message]bool))
	m.timer = time.NewTimer(cleanupInterval)
	go m.dispatch()
	go m.janitor()
	for i := Debug; i <= All; i++ {
		m.outbox[i] = make(chan *message)
		m.customers[i] = make(map[chan *message]bool)
		go m.postie(i)
	}
	return nil
}

func (m *Messages) PostMessage(severity Severity, owner, group, mode, plaintext, code string, args map[string]string) error {
	x := &message{
		Sev:  severity,
		Msg:  plaintext,
		Code: code,
		Args: args,
		Time: time.Now().Unix(),
	}
	return m.Post(owner, group, mode, x)
}

func (m *Messages) Post(owner, group, mode string, message *message) error {
	m.inbox <- message

	// insert into database
	tx, err := m.data.Begin()
	if err != nil {
		panic(nil)
	}
	defer tx.Rollback()

	result, err := tx.Exec("INSERT INTO alerts(time, severity, messagetext, messagecode, rulesowner, rulesgroup, rulesmode) VALUES(?,?,?,?,?,?,?)", message.Time, message.Sev, message.Msg, message.Code, owner, group, mode)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	if message.Args != nil {
		for key, value := range message.Args {
			_, err = tx.Exec("INSERT INTO args(id, key, value) VALUES(?,?,?)", id, key, value)
			if err != nil {
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		panic(nil)
	}

	return nil

}

func (m *Messages) janitor() {
	for {
		select {
		case <-m.timer.C:
			m.data.Exec("DELETE FROM alerts WHERE time < ?", time.Now().Unix()-int64(cleanupAge.Seconds()))
			m.timer.Reset(cleanupInterval)
		}
	}
}

func (m *Messages) dispatch() {
	for {
		select {
		case x := <-m.inbox:
			m.outbox[x.Sev] <- x
			m.outbox[All] <- x
			if x.Sev < Warning {
				m.outbox[Alert] <- x
			}
		}
	}
}

func (m *Messages) postie(department Severity) {
	customers := m.customers[department]
	for {
		select {
		case x := <-m.outbox[department]:
			for c := range customers {
				c <- x
			}
		}
	}
}
