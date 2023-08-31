# 🧻 `devslog` - [`slog.Handler`](https://pkg.go.dev/log/slog#Handler) for developing code
 [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/golang-cz/devslog/blob/master/LICENSE)
 [![Go Report Card](https://goreportcard.com/badge/github.com/golang-cz/devslog)](https://goreportcard.com/report/github.com/golang-cz/devslog)
 [![Go Reference](https://pkg.go.dev/badge/github.com/golang-cz/devslog.svg)](https://pkg.go.dev/github.com/golang-cz/devslog)

`devslog` is zero dependency custom logging handler for Go's standard [`log/slog`](https://pkg.go.dev/log/slog) package that provides structured logging with colorful and indented structure for developing.

### Develop with this output
![image](https://github.com/golang-cz/devslog/assets/17728576/cfdc1634-16fe-4dd0-a643-21bf519cd4fe)

### Instead of these outputs
`TextHandler`
![image](https://github.com/golang-cz/devslog/assets/17728576/49aab1c0-93ba-409d-8637-a96eeeaaf0e1)

`JSONHandler`
![image](https://github.com/golang-cz/devslog/assets/17728576/775af693-2f96-47e8-9190-5ead77b41a27)

## Install
```
go get github.com/golang-cz/devslog@latest
```

## Examples
### Logger without options
```go
w := os.Stdout

logger := slog.New(devslog.NewHandler(w, nil))

// optional: set global logger
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
| Parameter         | Description                                                    | Default      | Value  |
|-------------------|----------------------------------------------------------------|--------------|--------|
| MaxSlicePrintSize | Specifies the maximum number of elements to print for a slice. | 50           | uint   |
| SortKeys          | Determines if attributes should be sorted by keys.             | false        | bool   |
| TimeFormat        | Time format for timestamp.                                     | "[15:06:05]" | string |
