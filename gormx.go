package gormx

import (
	"log/slog"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var fetch = SingleWrap(Create)

// Options 定义了数据库连接的配置选项。
// 它是一个结构体，包含了连接数据库所需的信息以及调试模式的配置。
type Options struct {
	// Driver 是数据库驱动的名称，例如 "mysql" 或 "postgres"。
	// 这个字段是可选的，如果不需要指定驱动或者使用默认驱动，可以留空。
	Driver string `json:"driver,omitempty"`

	// DSN 是数据源名称（Data Source Name），包含了连接数据库所需的详细信息。
	// 这个字段的格式依赖于具体的数据库驱动，通常包括用户名、密码、主机和数据库名等信息。
	// 和 Driver 一样，这个字段也是可选的，取决于是否需要在连接时指定这些详细信息。
	DSN string `json:"dsn,omitempty"`

	// Debug 指示是否启用调试模式。
	// 当设置为 true 时，数据库操作的相关信息会被记录下来，通常用于开发或者调试阶段。
	// 在生产环境中，通常将这个值设置为 false，以避免不必要的性能开销。
	Debug bool `json:"debug,omitempty"`
}

// Default 返回一个默认的 *gorm.DB 实例，主要用于数据库操作。
// 该函数尝试通过调用 fetch 函数来获取数据库实例。如果 fetch 函数返回错误，
// 则构建一个带有错误信息的 *gorm.DB 实例，并确保其内部状态正确初始化。
func Default() *gorm.DB {
	// 尝试调用 fetch 函数来获取数据库实例。
	d, err := fetch("")
	// 如果 fetch 函数返回错误，对数据库实例进行错误处理。
	if err != nil {
		// 创建一个带有错误信息的新 *gorm.DB 实例。
		d = &gorm.DB{Error: err, Config: &gorm.Config{}}
		// 初始化该实例的 Statement 属性，确保其内部状态正确。
		d.Statement = &gorm.Statement{DB: d}
	}
	// 返回数据库实例。
	return d
}

// Get 是一个用于获取数据库连接的函数。
// 它接受一个数据库名称作为参数，并返回一个指向 gorm.DB 的指针和一个错误值。
// 该函数主要负责调用 fetch 函数来实际进行数据库连接的获取。
//
// 参数:
//
//	name - 指定要连接的数据库名称。
//
// 返回值:
//
//	*gorm.DB - 返回一个指向 gorm.DB 的指针，代表数据库连接。
//	error - 如果获取数据库连接时发生错误，返回错误信息。
func Get(name string) (*gorm.DB, error) {
	return fetch(name)
}

// Create 是一个用于创建数据库连接的方法。
// 它接受一个数据库名称作为参数，并根据该名称获取数据库配置。
// 如果没有指定数据库驱动和DSN，则使用默认的SQLite数据库和内存存储。
// 该方法返回一个*gorm.DB实例和一个错误（如果有的话）。
func Create(name string) (*gorm.DB, error) {
	// 获取数据库配置
	opts := getOpts(name)
	// 如果未指定数据库驱动和DSN，则使用默认值
	if opts.Driver == "" && opts.DSN == "" {
		opts.Driver = "sqlite"
		opts.DSN = ":memory:"
	}

	// 输出调试信息
	slog.Debug("[sql] open", "driver", opts.Driver, "dsn", opts.DSN, "debug", opts.Debug)
	// 使用获取的配置打开数据库连接
	d, err := Open(opts.Driver, opts.DSN)
	if err != nil {
		// 如果发生错误，返回nil和错误信息
		return nil, err
	}
	// 如果启用了调试模式，配置数据库日志记录
	if opts.Debug {
		d.Config.Logger = logger.Default.LogMode(logger.Info)
	}
	// 返回数据库连接和nil，表示成功
	return d, nil
}
