package shared

import (
	"database/sql"

	"github.com/gorilla/mux"
)

type Functions struct {
	Database func() *sql.DB
}

type Family interface {
	Route(*mux.Router)
}
