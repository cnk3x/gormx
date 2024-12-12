package gormx

import (
	"bytes"
	"cmp"
	"slices"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SortOptions 定义了排序选项的结构体。
// 主要用于指定对数据库表或模型进行操作所需的排序相关信息。
type SortOptions struct {
	Table      any    // 表名或模型。可以是表名或模型结构体。
	KeyColumn  string // 键列名，用于标识需要排序的记录，默认为主键。
	SortColumn string // 排序列名，用于指定排序依据的列，默认为 `sort`。
}

// SortPrep 为排序前的准备操作生成 SQL 表达式。
//
// 该函数根据提供的映射，创建一个 CASE 表达式来指定排序的值，并生成 WHERE 表达式来过滤结果。
//
// 参数:
//
//	values - 一个映射，其键和值将用于生成排序表达式。
//	kc - 一个表示排序键列的 clause.Column。
//	sc - 一个表示默认排序列的 clause.Column。
//
// 返回值:
//
//	where - 一个 clause.Expr，用于在查询中使用，以过滤出需要排序的记录。
//	value - 一个 clause.Expr，表示 CASE 表达式，用于指定排序的值。
func SortPrep[K cmp.Ordered, S cmp.Ordered](values map[K]S, kc, sc clause.Column) (where clause.Expr, value clause.Expr) {
	// 获取映射中键的数量，用于初始化键切片的容量。
	l := len(values)
	// 创建一个切片来存储映射的所有键。
	keys := make([]K, 0, l)
	// 遍历映射，将键添加到切片中。
	for key := range values {
		keys = append(keys, key)
	}
	// 对键切片进行排序，以确保 CASE 表达式的顺序。
	slices.Sort(keys)

	// 创建一个缓冲区来构建 CASE 表达式的 SQL 字符串。
	caseSql := bytes.NewBufferString(`(CASE ?`)
	// 创建一个切片来存储 CASE 表达式的所有参数。
	caseArg := make([]any, 0, l*2+2)
	// 将排序键列添加到参数中。
	caseArg = append(caseArg, kc)

	// 遍历排序后的键，构建 CASE 表达式的 WHEN THEN 部分。
	for _, key := range keys {
		caseSql.WriteString(` WHEN ? THEN ?`)
		caseArg = append(caseArg, key, values[key])
	}
	// 添加 CASE 表达式的 ELSE 部分和结束括号。
	caseSql.WriteString(` ELSE ? END)`)
	// 将默认排序列添加到参数中。
	caseArg = append(caseArg, sc)

	// 构建 WHERE 表达式，用于过滤出需要排序的记录。
	where = gorm.Expr(`? in (?)`, kc, keys)
	// 构建最终的 CASE 表达式 clause.Expr。
	value = gorm.Expr(caseSql.String(), caseArg...)
	// 返回 WHERE 和 VALUE 表达式。
	return
}

// SortExec 根据给定的键值对对数据库中的记录进行排序更新。
//
// 该函数使用 GORM 进行数据库操作，利用类型参数 K 和 S 来确保键和值的有序性。
//
// 参数:
//
//	tx - GORM 的数据库连接对象，如果为 nil，则使用默认连接。
//	values - 是一个映射，其键和值分别对应数据库记录中的键和排序值。
//	keyColumn 和 sortColumn - 是数据库表中的列名，分别用于标识键列和排序列。
//
// 函数返回更新操作后的 GORM DB 对象。
func SortExec[K cmp.Ordered, S cmp.Ordered](tx *gorm.DB, values map[K]S, keyColumn, sortColumn string) *gorm.DB {
	// 初始化键列和排序列的 Clause 对象。
	kc := column(keyColumn)
	sc := column(sortColumn)

	// 如果传入的 tx 为 nil，则使用默认的数据库连接。
	if tx == nil {
		tx = Default()
	}

	// 如果键列名为空，尝试从 Model 中获取主键名，如果 Model 为空，则默认为 "id"。
	if kc.Name == "" {
		if tx.Statement.Model != nil {
			kc.Name = clause.PrimaryKey
		} else {
			kc.Name = "id"
		}
	}

	// 如果排序列名为空，则默认为 "sort"。
	if sc.Name == "" {
		sc.Name = "sort"
	}

	// 调用 SortPrep 函数准备 WHERE 子句和更新值。
	where, value := SortPrep(values, kc, sc)

	// 执行更新操作，返回更新后的 DB 对象。
	return tx.Where(where).UpdateColumn(sc.Name, value)
}

// Sort 函数用于更新数据库中的排序信息。
//
// 该函数接收一个 *gorm.DB 类型的参数 tx，代表数据库事务，
// 以及一个映射类型 values，其键和值都实现了 Ordered 接口，
// 用于指定需要更新排序信息的实体及其新的排序值。
// 函数返回更新的行数和遇到的错误（如果有）。
func Sort[K cmp.Ordered, S cmp.Ordered](tx *gorm.DB, values map[K]S) (rowsUpdated int64, err error) {
	// 调用 SortExec 函数执行排序更新操作，传入空字符串作为排序字段和表名，
	// 这意味着 SortExec 需要根据上下文自行决定如何执行排序更新。
	tx = SortExec(tx, values, "", "")
	// 返回更新的行数和遇到的错误。
	return tx.RowsAffected, tx.Error
}
