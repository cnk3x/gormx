package gormx

import (
	"sync"

	"golang.org/x/sync/singleflight"
)

const DEFAULT = "DEFAULT"

// SingleWrap 是一个函数装饰器，用于缓存和去重处理。
// 它接受一个函数 get，该函数通过名称获取一个类型为 T 的实例。
// 返回一个新的函数，该函数会缓存 get 的调用结果，以避免重复获取相同的实例。
func SingleWrap[T any](get func(string) (T, error)) func(string) (T, error) {
	// ins 是一个缓存，用于存储通过名称创建的实例。
	var (
		ins = map[string]T{}
		// sfg 用于确保相同的 name 只会有一个 goroutine 在执行 get 操作。
		sfg singleflight.Group
		// mu 保护 ins 的读写操作，以确保并发安全。
		mu sync.RWMutex
	)

	// 返回一个新的函数，用于获取缓存或创建实例。
	return func(name string) (out T, err error) {
		// 如果 name 为空，则使用默认名称。
		if name == "" {
			name = DEFAULT
		}

		// 尝试从缓存中读取实例。
		mu.RLock()
		if instance, ok := ins[name]; ok {
			mu.RUnlock()
			// 如果找到实例，直接返回。
			return instance, nil
		}
		mu.RUnlock()

		// 使用 singleflight 机制，避免相同的 name 被同时多次调用。
		instance, err, _ := sfg.Do(name, func() (any, error) {
			// 调用原始的 get 函数获取实例。
			v, err := get(name)
			if err != nil {
				return nil, err
			}
			// 将获取的实例存储到缓存中。
			mu.Lock()
			ins[name] = v
			mu.Unlock()
			return v, nil
		})

		// 如果有错误发生，返回错误。
		if err != nil {
			return out, err
		}

		// 将结果转换为类型 T 并返回。
		return instance.(T), nil
	}
}
