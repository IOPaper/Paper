package values

import (
	"context"
	"reflect"

	"github.com/pkg/errors"
)

func Zero[T any]() T {
	var zero T
	return zero
}

func IsZero[T any](value T) bool {
	return reflect.DeepEqual(value, Zero[T]())
}

func Equal[T any](a, b T) bool {
	return reflect.DeepEqual(a, b)
}

func ContextAssertion[T any](ctx context.Context, key string) (T, error) {
	v := ctx.Value(key)
	x, ok := v.(T)
	if !ok {
		return Zero[T](), errors.Errorf("assertion error, values mismatch: %T", v)
	}
	return x, nil
}

func SliceFind[T any](elem T, slice []T) bool {
	for _, value := range slice {
		if Equal[T](elem, value) {
			return true
		}
	}
	return false
}
