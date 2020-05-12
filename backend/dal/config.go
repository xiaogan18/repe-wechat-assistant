package dal

import "errors"

type DbConfig struct {
	DriverName string
	DataSource string
}

var ErrNotFound error = errors.New("data not found")
