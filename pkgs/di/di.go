package di

import "github.com/samber/do"

type Injector struct {
	injector *do.Injector
}

func NewInjector() *Injector {
	return &Injector{
		injector: do.New(),
	}
}

func Provide[T any](d *Injector, fn func(d *Injector) (T, error)) {
	do.Provide(d.injector, func(_ *do.Injector) (T, error) {
		return fn(d)
	})
}

func Invoke[T any](d *Injector) (T, error) {
	return do.Invoke[T](d.injector)
}
