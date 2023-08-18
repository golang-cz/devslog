# 🧻 `devslog` - [`slog.Handler`](https://pkg.go.dev/log/slog#Handler) for developing code
 [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/golang-cz/devslog/blob/master/LICENSE)
 [![Go Report Card](https://goreportcard.com/badge/github.com/golang-cz/devslog)](https://goreportcard.com/report/github.com/golang-cz/devslog)
 [![Go Reference](https://pkg.go.dev/badge/github.com/golang-cz/devslog.svg)](https://pkg.go.dev/github.com/golang-cz/devslog)

`devslog` is zero dependency custom logging handler for Go's standard [`log/slog`](https://pkg.go.dev/log/slog) package that provides structured logging with colorful and indented structure for developing.

### Develop with this output
![image](https://github.com/golang-cz/devslog/assets/17728576/30a0d98a-a2de-4aa8-a4c4-e60d0c325049)

### Instead of these outputs
`TextHandler`
![image](https://github.com/golang-cz/devslog/assets/17728576/856f7e34-dc72-4f22-bd47-9fd5cbf7dd2f)
`JSONHandler`
![image](https://github.com/golang-cz/devslog/assets/17728576/3d4b091d-813a-461d-88e1-4cc95b9d6939)

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
| Parameter         | Description                                                                                                                                                        | Default          | Value  |
|-------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------|--------|
| MaxSlicePrintSize | Specifies the maximum number of elements to print for a slice.                                                                                                     | 50               | uint   |
| SortKeys          | Determines if attributes should be sorted by keys.                                                                                                                 | false            | bool   |
| TimeFormat        | Time format for timestamp.                                                                                                                                         | "[15:06:05]"     | string |
| ElementDivider    | Global variable. Character used to separate elements in slices and maps.  While change is possible, the default 'Group separator' (ASCII code 29)  is recommended. | string(rune(29)) | string |

## Custom functions
If `devslog` is not initialized, these functions simply pass the arguments to `slog.Any()`, ensuring no modification to the output when they are used in production environment.

### `Slice()`
Instead of relying on `slog.Any()`, you have the option to directly pass a slice of basic types, resulting in cleaner and formatted output. This approach enhances the readability of your logs. It uses `SliceElementDivider` for identifying single elements.
```go
sampleSlice := []string{"dsa", "ba na na"}

slog.Info(
	"some message",
	devslog.Slice("slice", sampleSlice),
)
```

### `Map()`
Similar to `Slice()`, you can pass a map of basic types.
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
