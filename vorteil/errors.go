package vorteil

import "errors"

var (
	ErrMode = errors.New("the mode setting from the config file is missing or invalid")
)
