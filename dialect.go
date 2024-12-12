package gormx

import (
	"fmt"

	"gorm.io/gorm"
)

var (
	drivers     = map[string]func(string) gorm.Dialector{}
	driverAlias = map[string]string{}
)

type DialectOpen = func(string) gorm.Dialector

// RegisterDriver 注册一个新的数据库驱动及其方言。
// 该函数接受驱动名称、方言实现以及可选的别名列表。
// 主要用于在数据库驱动和其对应方言之间建立映射，同时支持通过别名来引用相同的方言。
//
// 参数:
//
//	name - 驱动的唯一名称，用作在drivers map中的键。
//	dialect - DialectOpen类型，表示特定数据库的方言实现。
//	alias - 可变参数，包含驱动的别名，每个别名也将与主驱动名称建立映射。
func RegisterDriver(name string, dialect DialectOpen, alias ...string) {
	// 将驱动名称与方言实现建立映射，以便后续可以通过驱动名称获取方言。
	drivers[name] = dialect
	// 遍历别名列表，将每个别名与主驱动名称建立映射，支持通过别名引用相同的方言。
	for _, a := range alias {
		driverAlias[a] = name
	}
}

// Open 是一个用于初始化数据库连接的函数。
// 它接受数据库驱动名称、数据源名称（DSN）以及可选的 GORM 配置选项作为参数。
// 函数返回一个 *gorm.DB 实例，用于与数据库进行交互，或者返回一个错误，如果连接失败。
func Open(driver, dsn string, opts ...gorm.Option) (*gorm.DB, error) {
	// 使用 driver 参数值初始化 name 变量，用于后续查找对应的数据库方言。
	name := driver

	// 尝试根据数据库名称获取对应的数据库方言构造函数。
	dialect, ok := drivers[name]
	// 如果没有找到对应的方言，检查是否存在该数据库驱动的别名。
	if !ok {
		// 如果存在别名，使用别名再次尝试获取数据库方言构造函数。
		if name, ok = driverAlias[driver]; ok {
			dialect, ok = drivers[name]
		}
	}

	// 如果仍然没有找到对应的方言，返回一个未知驱动的错误。
	if !ok {
		return nil, fmt.Errorf("unknown driver: %s", driver)
	}

	// 使用找到的数据库方言构造函数和提供的 DSN 初始化数据库连接。
	// 同时传入所有的 GORM 配置选项。
	return gorm.Open(dialect(dsn), opts...)
}
