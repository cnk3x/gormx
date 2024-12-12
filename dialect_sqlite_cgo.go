//go:build (!mssql && !mysql && !postgres && !sqlite) || (sqlite && cgo)

package gormx

import (
	"gorm.io/driver/sqlite"
)

func init() { RegisterDriver("sqlite", sqlite.Open) }
