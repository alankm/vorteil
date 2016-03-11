package privileges

import (
	"database/sql"

	"github.com/alankm/vorteil/core/shared"
)

type Privileges struct {
	data *sql.DB
}

func New(functions *shared.Functions) (*Privileges, error) {
	privs := &Privileges{
		data: functions.Database(),
	}
	return privs, nil
}

type Session struct {
}

func (s *Session) User() string {
	return ""
}

func (s *Session) Access(path string) error {
	return nil
}
