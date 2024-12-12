//go:build sqlite && !cgo

package gormx

import (
	_ "github.com/ncruces/go-sqlite3/embed"
	sqlite "github.com/ncruces/go-sqlite3/gormlite"
)

func init() { RegisterDriver("sqlite", sqlite.Open) }
