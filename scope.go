package gormx

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Scope 是一个定义了如何对数据库查询进行条件化修改的函数类型。
// 它接收一个 *gorm.DB 类型的参数，代表 GORM 的数据库连接实例，
// 并返回一个经过修改的 *gorm.DB 类型的实例。
// 使用 Scope 可以动态地修改数据库查询，例如添加额外的条件、排序规则等。
type Scope func(*gorm.DB) *gorm.DB

// Like 创建一个查询范围，用于在数据库查询中添加LIKE条件。
// 该函数主要用于实现模糊查询，通过在指定列中搜索包含查询字符串q的项。
//
// 参数:
//
//	column: 数据库列名，表示要在哪一列中进行模糊查询。
//	q: 查询字符串，表示要搜索的关键字。
//
// 返回值:
//
//	Scope: 返回一个函数，该函数接收一个*gorm.DB实例，并返回添加了LIKE条件的*gorm.DB实例。
func Like(column, q string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("? LIKE ?", clause.Column{Name: column}, "%"+q+"%")
	}
}

// Prefix 生成一个查询范围，用于在指定列上应用前缀匹配
// 它允许通过在查询字符串后添加百分号来匹配列值的前缀
//
// 参数:
//
//	col: 要应用前缀匹配的数据库列名
//	q: 要查询的前缀字符串
//
// 返回值:
//
//	Scope: 一个函数，接受一个*gorm.DB实例，返回应用了前缀匹配条件的*gorm.DB实例
func Prefix(col, q string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("? LIKE ?", column(col), q+"%")
	}
}

// Suffix 生成一个查询范围，用于在指定列上应用后缀匹配
// 它允许通过在查询字符串后添加百分号来匹配列值的后缀
//
// 参数:
//
//	column: 要应用后缀匹配的数据库列名
//	q: 要查询的后缀字符串
//
// 返回值:
//
//	Scope: 一个函数，接受一个*gorm.DB实例，返回应用了后缀匹配条件的*gorm.DB实例
func Suffix(col, q string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("? LIKE ?", column(col), "%"+q)
	}
}

// Paging 是一个泛型函数，用于创建一个分页查询的范围。
// 它接受页码（page）、每页大小（size）和一个可选的默认每页大小（defSize）作为参数。
// 该函数返回一个 Scope 函数，该函数对传入的 *gorm.DB 实例应用分页逻辑。
func Paging[T1 Integer, T2 Integer, T3 Integer](page T1, size T2, defSize ...T3) Scope {
	// 将页码和每页大小转换为 int 类型。
	p, s, d := int(page), int(size), 1000

	// 遍历 defSize 参数，如果传入了大于 0 的值，则将其作为默认每页大小。
	for _, v := range defSize {
		if v > 0 {
			d = int(v)
			break
		}
	}

	// 如果每页大小为 0，则使用默认的每页大小。
	if s == 0 {
		s = d
	}

	// 返回一个 Scope 函数，该函数对传入的 *gorm.DB 实例应用分页逻辑。
	return func(db *gorm.DB) *gorm.DB {
		// 如果页码大于 1，则应用 OFFSET 分页逻辑。
		if p > 1 {
			db = db.Offset((p - 1) * s)
		}
		// 应用 LIMIT 限制查询结果的数量。
		return db.Limit(s)
	}
}

// OrderBy 根据传入的排序参数构建排序查询。
// 该函数接收两个参数：orderBy 是用户指定的排序参数，def 是默认的排序参数。
// 它返回一个 Scope 函数，该函数可以应用于 gorm.DB 对象以添加排序条件。
func OrderBy(orderBy string, def string) Scope {
	// calc 是一个内部函数，用于处理排序字符串。
	// 它接收一个字符串 in，并返回一个格式化后的排序字符串。
	calc := func(in string) string {
		if in != "" {
			var x int
			// 将输入字符串按逗号分割成多个排序项。
			orders := strings.Split(in, ",")
			for _, it := range orders {
				// 去除排序项的前后空格。
				if it = strings.TrimSpace(it); it != "" && it != "-" {
					// 检查排序项是否以 '-' 开头，以确定是升序还是降序。
					if it[0] == '-' {
						// 如果是降序，移除 '-' 并添加 ' DESC'。
						orders[x] = it[1:] + " DESC"
					} else {
						// 如果是升序，直接添加 ' ASC'。
						orders[x] = it + " ASC"
					}
					x++
				}
			}
			// 如果有有效的排序项，将它们连接成一个字符串并返回。
			if len(orders) > 0 {
				return strings.Join(orders[:x], ", ")
			}
		}
		// 如果输入字符串为空，返回空字符串。
		return ""
	}

	// 使用 calc 函数处理传入的 orderBy 参数。
	orders := calc(orderBy)
	// 如果处理结果为空，使用默认的排序参数 def 进行处理。
	if orders == "" {
		orders = calc(def)
	}

	// 返回一个 Scope 函数，用于应用排序条件到 gorm.DB 对象。
	return func(d *gorm.DB) *gorm.DB {
		// 如果有有效的排序条件，使用 Order 方法添加排序条件到 d。
		if orders != "" {
			d = d.Order(orders)
		}
		// 返回 d。
		return d
	}
}
