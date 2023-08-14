package devslog

import (
	"fmt"
	"log/slog"
	"strings"
)

type basicTypes interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float32 | float64 | string | bool
}

// same as slog.Any() just pass a Slice with basic types.
// bool |
// string |
// float32 | float64 |
// int | int8 | int16 | int32 | int64 |
// uint | uint8 | uint16 | uint32 | uint64 |
// uintptr
func Slice[T basicTypes](key string, s []T) slog.Attr {
	if initialized {
		var elements []string
		for _, e := range s {
			elements = append(elements, fmt.Sprintf("%v", e))
		}

		merged := strings.Join(elements, string(sliceElementDivider))
		val := fmt.Sprintf("slice[%s]", merged)
		return slog.String(key, val)
	}

	return slog.Any(key, s)
}

// same as slog.Any() just pass a Map with basic types
// bool |
// string |
// float32 | float64 |
// int | int8 | int16 | int32 | int64 |
// uint | uint8 | uint16 | uint32 | uint64 |
// uintptr
func Map[K basicTypes, V basicTypes](key string, m map[K]V) slog.Attr {
	if initialized {
		var elements []string
		for k, v := range m {
			elements = append(elements, fmt.Sprintf("%v:%v", k, v))
		}

		merged := strings.Join(elements, string(sliceElementDivider))
		val := fmt.Sprintf("map[%s]", merged)
		return slog.String(key, val)
	}

	return slog.Any(key, m)
}
