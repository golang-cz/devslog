package devslog

import (
	"log/slog"
)

type attributes []slog.Attr

func (a attributes) Len() int      { return len(a) }
func (a attributes) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a attributes) Less(i, j int) bool {
	if a[i].Value.Kind() == slog.KindGroup && a[j].Value.Kind() != slog.KindGroup {
		return false
	} else if a[i].Value.Kind() != slog.KindGroup && a[j].Value.Kind() == slog.KindGroup {
		return true
	}

	return a[i].Key < a[j].Key
}

func (a attributes) padding(c foregroundColor, colorFunction func(b []byte, fgColor foregroundColor) []byte) int {
	var padding int
	for _, e := range a {
		color := len(e.Key)
		if c != nil {
			color = len(colorFunction([]byte(e.Key), c))
		}

		if color > padding {
			padding = color
		}
	}

	return padding
}
