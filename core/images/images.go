package images

import (
	"database/sql"

	"github.com/alankm/vorteil/core/shared"
)

type Images struct {
	data *sql.DB
}

func New(functions *shared.Functions) (*Images, error) {
	imgs := &Images{
		data: functions.Database(),
	}
	return imgs, nil
}
