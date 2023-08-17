# 🧻 `devslog` - [`slog.Handler`](https://pkg.go.dev/log/slog#Handler) for developing code
 [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/golang-cz/devslog/blob/master/LICENSE)
 [![Go Report Card](https://goreportcard.com/badge/github.com/golang-cz/devslog)](https://goreportcard.com/report/github.com/golang-cz/devslog)
 [![Go Reference](https://pkg.go.dev/badge/github.com/golang-cz/devslog.svg)](https://pkg.go.dev/github.com/golang-cz/devslog)

`devslog` is zero dependency custom logging handler for Go's standard [`log/slog`](https://pkg.go.dev/log/slog) package that provides structured logging with colorful and indented structure for developing.

![image](https://github.com/golang-cz/devslog/assets/17728576/7c870a4e-1f77-4bb1-a042-db2704ba8547)

## Install
```
go get github.com/golang-cz/devslog@latest
```

## Examples
### Logger without options
```go
w := os.Stdout

logger := slog.New(devslog.NewHandler(w, nil))

// set global logger
slog.SetDefault(logger)
```

### Logger with custom options
```go
w := os.Stdout

// new logger with options
opts := &devslog.Options{
	MaxSlicePrintSize: 4,
	SortKeys:          true,
	TimeFormat:        "[06:05]"
}

logger := slog.New(devslog.NewHandler(w, opts))

// optional: set global logger
slog.SetDefault(logger)
```

### Logger with default slog options
Handler accept default [slog.HandlerOptions](https://pkg.go.dev/golang.org/x/exp/slog#HandlerOptions)
```go
w := os.Stdout

// slog.HandlerOptions
slogOpts := &slog.HandlerOptions{
	AddSource:   true,
	Level:       slog.LevelDebug,
}

// new logger with options
opts := &devslog.Options{
	HandlerOptions:    slogOpts,
	MaxSlicePrintSize: 4,
	SortKeys:          true,
}

logger := slog.New(devslog.NewHandler(w, opts))

// optional: set global logger
slog.SetDefault(logger)
```

### Example of initialization with production
```go
production := false

w := os.Stdout

slogOpts := &slog.HandlerOptions{
	AddSource: true,
	Level:     slog.LevelDebug,
}

var logger *slog.Logger
if production {
	logger = slog.New(slog.NewJSONHandler(w, slogOpts))
} else {
	opts := &devslog.Options{
		HandlerOptions:    slogOpts,
		MaxSlicePrintSize: 10,
		SortKeys:          true,
	}

	logger = slog.New(devslog.NewHandler(w, opts))
}

// optional: set global logger
slog.SetDefault(logger)
```

## Options
| Parameter           | Description                                                                                                                                                                         | Default          | Value  |
|---------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------|--------|
| MaxSlicePrintSize   | Defines the maximum number of elements for slice, that would be printed.                                                                                                            | 50               | uint   |
| SortKeys            | If the attributes should be sorted by keys.                                                                                                                                         | false            | bool   |
| TimeFormat          | Time format for timestamp.                                                                                                                                                          | "[15:06:05]"     | string |
| SliceElementDivider | Global variable. Character for splitting elements in slices and maps. I don't recommend to change this value, but i made it optional. Default is 'Group separator' - ASCII code 29. | string(rune(29)) | string |

## Custom functions
If the `devslog` is not initialized, then these functions just pass the arguments to `slog.Any()`, so in the production the output wouldn't be modified.

### Slice()
You can pass slice with basic types instead of use `slog.Any()`
```go
sampleSlice := []string{"dsa", "ba na na"}

slog.Info(
	"some message",
	devslog.Slice("slice", sampleSlice),
)
```

### Map()
You can pass map with basic types instead of use `slog.Any()`
```go
sampleMap := map[string]string{
	"apple":    "pear",
	"ba na na": "man go",
}

slog.Info(
	"some message",
	devslog.Map("map", sampleMap),
)
```
