package gormx

import (
	"strings"
	"unicode"

	"gorm.io/gorm/clause"
)

func column(column string) (col clause.Column) {
	col.Name = column
	if i := strings.Index(strings.ToLower(column), " as "); i > 0 {
		col.Alias = column[i+4:]
		col.Name = column[:i]
	}

	if t, n, ok := strings.Cut(column, "."); ok {
		col.Table = t
		col.Name = n
	}

	col.Name = strings.TrimFunc(col.Name, nameClean)
	col.Table = strings.TrimFunc(col.Table, nameClean)
	col.Alias = strings.TrimFunc(col.Alias, nameClean)

	if col.Table == "" {
		col.Table = clause.CurrentTable
	}

	return
}

func nameClean(r rune) bool {
	return r == '"' || r == '`' || r == '\'' || r == '[' || r == ']' || unicode.IsSpace(r)
}

type (
	Int interface {
		~int | ~int8 | ~int16 | ~int32 | ~int64
	}
	Uint interface {
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
	}
	Integer interface{ Uint | Int }
	// Float   interface{ ~float32 | ~float64 }
	// Ordered interface{ Integer | Float | ~string }
)
