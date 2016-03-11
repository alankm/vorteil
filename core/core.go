package core

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"os"

	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/yaml.v2"

	"github.com/alankm/vorteil/core/images"
	"github.com/alankm/vorteil/core/messages"
	"github.com/alankm/vorteil/core/privileges"
	"github.com/alankm/vorteil/core/shared"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/mattn/go-sqlite3"
	"github.com/sisatech/multiserver"
	"github.com/sisatech/raft"
)

func init() {
	hook := func(conn *sqlite3.SQLiteConn) error {
		_, err := conn.Exec("PRAGMA foreign_keys = ON", nil)
		return err
	}
	driver := &sqlite3.SQLiteDriver{ConnectHook: hook}
	sql.Register("sql_fk", driver)
}

type Core struct {
	logs log15.Logger
	data *sql.DB
	conf *CoreConfig
	imgs *images.Images
	msgs *messages.Messages
	priv *privileges.Privileges
	raft *Cluster
	web  *WebServer
}

type CoreConfig struct {
	Database string                 `yaml:"database"`
	Data     string                 `yaml:"data"`
	Bind     string                 `yaml:"bind"`
	Modules  map[string]interface{} `yaml:"modules"`
}

type WebServer struct {
	server   *multiserver.HTTPServer
	cookie   *securecookie.SecureCookie
	mux      *mux.Router
	services []string
}

type Cluster struct {
	config *raft.Config
	server *raft.Raft
	client *raft.Client
}

func loadConfig(file string) (*CoreConfig, error) {
	src, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	config := new(CoreConfig)
	err = yaml.Unmarshal(src, config)
	if err != nil {
		return nil, err
	}
	err = config.setDefaults()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *CoreConfig) setDefaults() error {
	if c.Data == "" {
		return errors.New("config file must identify a Data directory")
	}
	if c.Database == "" {
		c.Database = c.Data + "/vorteil.db"
	}
	return nil
}

func New(config string) (*Core, error) {
	cc, err := loadConfig(config)
	if err != nil {
		return nil, err
	}
	c := new(Core)
	c.logs = log15.New("server", cc.Bind)
	c.logs.Info("Launching Vorteil...")
	c.conf = cc
	err = os.MkdirAll(c.conf.Data, 0772)
	if err != nil {
		c.logs.Crit("creating data directory", "error", err.Error())
		return nil, err
	}
	c.data, err = sql.Open("sql_fk", c.conf.Database)
	if err != nil {
		c.logs.Crit("opening database", "error", err.Error())
		return nil, err
	}
	functions := &shared.Functions{
		Database: c.Database,
	}
	c.priv, err = privileges.New(functions)
	if err != nil {
		c.logs.Crit("starting privileges", "error", err.Error())
		return nil, err
	}
	c.msgs, err = messages.New(functions)
	if err != nil {
		c.logs.Crit("starting messages", "error", err.Error())
		return nil, err
	}
	c.imgs, err = images.New(functions)
	if err != nil {
		c.logs.Crit("starting images", "error", err.Error())
		return nil, err
	}
	c.web = new(WebServer)
	c.web.mux = mux.NewRouter()
	c.web.server = multiserver.NewHTTPServer(c.conf.Bind, c.web.mux, nil)
	c.web.cookie = securecookie.New(securecookie.GenerateRandomKey(64),
		securecookie.GenerateRandomKey(32))
	c.web.mux.Handle("/services/login", &GeneralHandler{c.loginHandler})
	return c, nil
}

func (c *Core) Start() error {
	err := c.web.server.Start()
	if err != nil {
		return err
	}
	return nil
}

func (c *Core) Stop() {
	c.data.Close()
}

func (c *Core) Database() *sql.DB {
	return c.data
}
