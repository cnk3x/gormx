//go:build mysql

package gormx

import (
	"gorm.io/driver/mysql"
)

func init() { RegisterDriver("mysql", mysql.Open) }
