package gormx

import (
	"log/slog"
	"testing"
	"unsafe"

	"gorm.io/gorm"
)

func TestX(t *testing.T) {
	t.Logf("%v", 111)

	var s string = "123"
	var s1 string = "45611"
	var i int = 123
	var i1 int = 456
	t.Log(getType(&s), getType(&s1))
	t.Log(getType(&i), getType(&i1))
	t.Log(unsafe.Pointer(&s))
}

func getType[T any](value *T) uintptr {
	u := (*[2]uintptr)(unsafe.Pointer(value))
	return u[1]
}

func TestSort(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	sorts := map[string]int{
		"5": 5,
		"1": 2, "2": 1,
		"3": 3, "4": 4,
	}

	sql := Default().ToSQL(func(tx *gorm.DB) *gorm.DB {
		return SortExec(tx.Model(&ZZ{}), sorts, "", "")
	})

	t.Logf("log: %s", sql)
}

type ZZ struct {
	ID        int
	Sort      int
	UpdatedAt int64
}
