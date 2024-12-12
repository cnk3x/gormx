//go:build mssql

package gormx

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func PostgresOpen(dsn string) gorm.Dialector {
	return sqlserver.Open(dsn)
}

func init() {
	RegisterDriver("mssql", sqlserver.Open)
	RegisterDriver("sqlserver", sqlserver.Open)
}
