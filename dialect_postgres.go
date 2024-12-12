//go:build postgres

package gormx

import (
	"gorm.io/driver/postgres"
)

func init() {
	RegisterDriver("postgres", postgres.Open)
	RegisterDriver("pg", postgres.Open)
	RegisterDriver("postgresql", postgres.Open)
}
