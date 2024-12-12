package gormx

import (
	"os"
	"strconv"
	"strings"
)

var (
	envPrefix  = ""
	getOptions func(name string) (opts Options)
)

// SetOptionsFunc 是一个用于设置选项的函数。
// 它接受一个函数作为参数，该函数应接受一个字符串参数 name，并返回一个 Options 对象。
// 当需要动态设置或获取配置选项时，可以通过调用 SetOptionsFunc 来实现。
// 这种设计允许在运行时灵活地调整或定制选项，而无需在编译时固定这些选项。
//
// 参数:
//
//	fn - 一个函数，根据提供的名称返回相应的配置选项。
//
// 设置:
//
//	getOptions - 通过调用 fn，可以动态地获取或设置配置选项。
func SetOptionsFunc(fn func(name string) Options) { getOptions = fn }

// SetEnvPrefix 设置环境变量的前缀。
// 此函数允许用户在全局范围内更改环境变量的前缀，以便在大型项目或复杂环境中更好地管理配置。
//
// 参数:
//
//	prefix: 要设置的新环境变量前缀。这应该是一个简洁且具有描述性的字符串，用于标识项目或应用程序的环境变量。
func SetEnvPrefix(prefix string) { envPrefix = prefix }

func getOpts(name string) Options {
	get := getOptions
	if get == nil {
		get = defaultOptions
	}
	return get(name)
}

func defaultOptions(name string) (opts Options) {
	opts.Driver = fromEnv("DRIVER", name)
	opts.DSN = fromEnv("DSN", name)
	opts.Debug, _ = strconv.ParseBool(fromEnv("DEBUG", name))
	return
}

func fromEnv(field, name string) string {
	if name == DEFAULT || name == "" {
		name = ""
	} else {
		name = strings.ToUpper(name)
	}

	var p string
	if envPrefix != "" {
		p = envPrefix + "_DB_" + field
	} else {
		p = "DB_" + field
	}

	if name != "" {
		p += "_"
	}

	return os.Getenv(p + name)
}
