package devslog

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"testing"
	"time"
)

func Test_NewHandlerDefaults(t *testing.T) {
	sliceElementDivider = string(rune(29))
	opts := &Options{
		HandlerOptions: &slog.HandlerOptions{},
	}
	h := NewHandler(os.Stdout, opts)

	if h.opts.Level.Level() != slog.LevelInfo.Level() {
		t.Errorf("Expected default log level to be LevelInfo")
	}

	if h.opts.MaxSlicePrintSize != 50 {
		t.Errorf("Expected default MaxSlicePrintSize to be 50")
	}

	if h.opts.TimeFormat != "[15:06:05]" {
		t.Errorf("Expected default TimeFormat to be \"[15:06:05]\" ")
	}

	if h.out == nil {
		t.Errorf("Expected writer to be initialized")
	}

	if h.mu == nil {
		t.Errorf("Expected mutex to be initialized")
	}

	if !initialized {
		t.Error("Expected initialized to be true")
	}
}

func Test_NewHandlerWithOptions(t *testing.T) {
	handlerOpts := &Options{
		HandlerOptions:      &slog.HandlerOptions{Level: slog.LevelWarn},
		MaxSlicePrintSize:   10,
		SliceElementDivider: "||",
		TimeFormat:          "[06:05]",
	}
	h := NewHandler(nil, handlerOpts)

	if h.opts.Level.Level() != slog.LevelWarn.Level() {
		t.Errorf("Expected custom log level to be LevelWarn")
	}

	if h.opts.MaxSlicePrintSize != 10 {
		t.Errorf("Expected custom MaxSlicePrintSize to be 10")
	}

	if h.opts.TimeFormat != "[06:05]" {
		t.Errorf("Expected custom TimeFormat to be \"[06:05]\" ")
	}

	if h.opts.SliceElementDivider != "||" {
		t.Errorf("Expected custom SliceElementDivider to be '||'")
	}
}

func Test_NewHandlerWithNilOptions(t *testing.T) {
	h := NewHandler(nil, nil)

	if h.opts.HandlerOptions == nil || h.opts.HandlerOptions.Level != slog.LevelInfo {
		t.Errorf("Expected HandlerOptions to be initialized with default level")
	}

	if h.opts.MaxSlicePrintSize != 50 {
		t.Errorf("Expected MaxSlicePrintSize to be initialized with default value")
	}

	if h.out != nil {
		t.Errorf("Expected writer to be nil")
	}
}

func Test_Enabled(t *testing.T) {
	h := NewHandler(nil, nil)
	ctx := context.Background()

	if !h.Enabled(ctx, slog.LevelInfo) {
		t.Error("Expected handler to be enabled for LevelInfo")
	}

	if h.Enabled(ctx, slog.LevelDebug) {
		t.Error("Expected handler to be disabled for LevelDebug")
	}
}

func Test_WithGroup(t *testing.T) {
	h := NewHandler(nil, nil)
	h2 := h.WithGroup("myGroup")

	if h2 == h {
		t.Error("Expected a new handler instance")
	}
}

func Test_WithGroupEmpty(t *testing.T) {
	h := NewHandler(nil, nil)
	h2 := h.WithGroup("")

	if h2 != h {
		t.Error("Expected a original handler instance")
	}
}

func Test_WithAttrs(t *testing.T) {
	h := NewHandler(nil, nil)
	h2 := h.WithAttrs([]slog.Attr{slog.String("key", "value")})

	if h2 == h {
		t.Error("Expected a new handler instance")
	}
}

func Test_WithAttrsEmpty(t *testing.T) {
	h := NewHandler(nil, nil)
	h2 := h.WithAttrs([]slog.Attr{})

	if h2 != h {
		t.Error("Expected a original handler instance")
	}
}

func Test_IsURL(t *testing.T) {
	urlString := "https://www.example.com"
	if !isURL(urlString) {
		t.Errorf("Expected URL to be recognized as URL: %s", urlString)
	}

	nonURLString := "not-a-valid-url"
	if isURL(nonURLString) {
		t.Errorf("Expected non-URL string to not be recognized as URL: %s", nonURLString)
	}
}

func Test_IsMap(t *testing.T) {
	mapString := "map[key:value]"
	if !isMap(mapString) {
		t.Errorf("Expected string to be recognized as map: %s", mapString)
	}

	nonMapString := "not-a-valid-map"
	if isMap(nonMapString) {
		t.Errorf("Expected non-map string to not be recognized as map: %s", nonMapString)
	}
}

func Test_IsSlice(t *testing.T) {
	sliceString := "slice[1 2 3]"
	if !isSlice(sliceString) {
		t.Errorf("Expected string to be recognized as slice: %s", sliceString)
	}

	nonSliceString := "not-a-valid-slice"
	if isSlice(nonSliceString) {
		t.Errorf("Expected non-slice string to not be recognized as slice: %s", nonSliceString)
	}
}

func Test_ArrayString(t *testing.T) {
	h := NewHandler(nil, nil)
	data := []string{"apple", "ba na na"}
	expected := "\x1b[33m2\x1b[0m \x1b[32mslice[\x1b[0m\n    \x1b[32m0\x1b[0m: \x1b[34mapple\x1b[0m\n    \x1b[32m1\x1b[0m: \x1b[34mba na na\x1b[0m \x1b[32m]\x1b[0m"

	result := h.arrayString(fmt.Sprintf("slice[%v%v%v]", data[0], sliceElementDivider, data[1]), 0)

	if result != expected {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, result)
	}
}

func Test_ArrayStringEmpty(t *testing.T) {
	h := NewHandler(nil, nil)
	expected := "\x1b[33m0\x1b[0m \x1b[32mslice[]\x1b[0m"

	result := h.arrayString("slice[]", 0)

	if result != expected {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, result)
	}
}

func Test_ArrayStringBig(t *testing.T) {
	opts := &Options{
		MaxSlicePrintSize: 4,
	}

	h := NewHandler(nil, opts)
	slice := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		slice[i] = i + 1
	}

	expected := "\x1b[33m1000\x1b[0m \x1b[32mslice[\x1b[0m\n    \x1b[32m0\x1b[0m: \x1b[34m1\x1b[0m\n    \x1b[32m1\x1b[0m: \x1b[34m2\x1b[0m\n    \x1b[32m2\x1b[0m: \x1b[34m3\x1b[0m\n    \x1b[32m3\x1b[0m: \x1b[34m4\x1b[0m\n         \x1b[34m...\x1b[0m\x1b[32m]\x1b[0m"

	attr := Slice("key", slice)
	result := h.arrayString(attr.Value.String(), 0)

	if result != expected {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, result)
	}
}

func Test_MapString(t *testing.T) {
	h := NewHandler(nil, nil)
	data := map[string]string{"a": "1", "b": "2"}
	expected := "\x1b[33m2\x1b[0m \x1b[32mmap[\x1b[0m\n    \x1b[32ma\x1b[0m : \x1b[34m1\x1b[0m\n    \x1b[32mb\x1b[0m : \x1b[34m2\x1b[0m \x1b[32m]\x1b[0m"

	result := h.mapString(fmt.Sprintf("map[%s:%s%s%s:%s]", "a", data["a"], sliceElementDivider, "b", data["b"]), 0)

	if result != expected {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, result)
	}
}

func Test_MapStringEmpty(t *testing.T) {
	h := NewHandler(nil, nil)
	expected := "\x1b[33m0\x1b[0m \x1b[32mmap[]\x1b[0m"

	result := h.mapString("map[]", 0)

	if result != expected {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, result)
	}
}

func Test_LevelMessageDebug(t *testing.T) {
	h := NewHandler(nil, nil)
	buf := make([]byte, 0)
	record := &slog.Record{
		Level:   slog.LevelDebug,
		Message: "Debug message",
	}

	buf = h.levelMessage(buf, record)

	expected := "\x1b[30m\x1b[44m DEBUG \x1b[0m \x1b[34mDebug message\x1b[0m\n"
	result := string(buf)

	if result != expected {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, result)
	}
}

func Test_LevelMessageInfo(t *testing.T) {
	h := NewHandler(nil, nil)
	buf := make([]byte, 0)
	record := &slog.Record{
		Level:   slog.LevelInfo,
		Message: "Info message",
	}

	buf = h.levelMessage(buf, record)

	expected := "\x1b[30m\x1b[42m INFO \x1b[0m \x1b[32mInfo message\x1b[0m\n"
	result := string(buf)

	if result != expected {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, result)
	}
}

func Test_LevelMessageWarn(t *testing.T) {
	h := NewHandler(nil, nil)
	buf := make([]byte, 0)
	record := &slog.Record{
		Level:   slog.LevelWarn,
		Message: "Warning message",
	}

	buf = h.levelMessage(buf, record)

	expected := "\x1b[30m\x1b[43m WARN \x1b[0m \x1b[33mWarning message\x1b[0m\n"
	result := string(buf)

	if result != expected {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, result)
	}
}

func Test_LevelMessageError(t *testing.T) {
	h := NewHandler(nil, nil)
	buf := make([]byte, 0)
	record := &slog.Record{
		Level:   slog.LevelError,
		Message: "Error message",
	}

	buf = h.levelMessage(buf, record)

	expected := "\x1b[30m\x1b[41m ERROR \x1b[0m \x1b[31mError message\x1b[0m\n"
	result := string(buf)

	if result != expected {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, result)
	}
}

func (w *MockWriter) Write(p []byte) (int, error) {
	w.WrittenData = append(w.WrittenData, p...)
	return len(p), nil
}

type MockWriter struct {
	WrittenData []byte
}

func Test_WholeOutput(t *testing.T) {
	w := &MockWriter{}

	slogOpts := &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr { return a },
	}

	opts := &Options{
		HandlerOptions:    slogOpts,
		MaxSlicePrintSize: 4,
		SortKeys:          true,
		TimeFormat:        "[15:06]",
	}

	logger := slog.New(NewHandler(w, opts).WithAttrs([]slog.Attr{slog.String("attr", "string")}).WithGroup("with_group"))

	mapString := map[string]string{
		"apple":    "pear",
		"ba na na": "man go",
	}

	sliceSmall := []string{"dsa", "ba na na"}
	sliceBig := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		sliceBig[i] = i + 1
	}

	emptySlice := make([]int, 0)
	emptyMap := make(map[int]int, 0)

	timeString := csf(time.Now().Format("[15:06]"), fgWhite)
	logger.Info(
		"My INFO message",
		slog.String("test_string", "some string"),
		slog.String("empty", ""),
		slog.String("url", "https://go.dev/"),
		slog.Any("boolean", true),
		slog.Any("time", time.Date(2012, time.March, 28, 0, 0, 0, 0, time.UTC)),
		slog.Any("duration", time.Second),
		Map("map", mapString),
		slog.Any("map", mapString),
		Map("empty_map", emptyMap),
		Slice("slice", sliceSmall),
		Slice("slice_big", sliceBig),
		Slice("empty_slice", emptySlice),
		slog.Group("my_group",
			slog.Any("int", 1),
			slog.Any("float", 1.21),
		),
	)

	expected := fmt.Sprintf("%[1]v \x1b[30m\x1b[42m INFO \x1b[0m \x1b[32mMy INFO message\x1b[0m\n  \x1b[35mattr\x1b[0m       : string\n\x1b[32mG\x1b[0m \x1b[35mwith_group\x1b[0m : \x1b[32mgroup\x1b[0m\n  \x1b[31m#\x1b[0m \x1b[35mboolean\x1b[0m     : \x1b[31mtrue\x1b[0m\n  \x1b[36m@\x1b[0m \x1b[35mduration\x1b[0m    : \x1b[36m1s\x1b[0m\n    \x1b[35mempty\x1b[0m       : \x1b[37m\x1b[2mempty\x1b[0m\n  \x1b[32mM\x1b[0m \x1b[35mempty_map\x1b[0m   : \x1b[33m0\x1b[0m \x1b[32mmap[]\x1b[0m\n  \x1b[32mS\x1b[0m \x1b[35mempty_slice\x1b[0m : \x1b[33m0\x1b[0m \x1b[32mslice[]\x1b[0m\n  \x1b[32mM\x1b[0m \x1b[35mmap\x1b[0m         : \x1b[33m2\x1b[0m \x1b[32mmap[\x1b[0m\n      \x1b[32mapple\x1b[0m    : \x1b[34mpear\x1b[0m\n      \x1b[32mba na na\x1b[0m : \x1b[34mman go\x1b[0m \x1b[32m]\x1b[0m\n    \x1b[35mmap\x1b[0m         : map[apple:pear ba na na:man go]\n  \x1b[32mS\x1b[0m \x1b[35mslice\x1b[0m       : \x1b[33m2\x1b[0m \x1b[32mslice[\x1b[0m\n      \x1b[32m0\x1b[0m: \x1b[34mdsa\x1b[0m\n      \x1b[32m1\x1b[0m: \x1b[34mba na na\x1b[0m \x1b[32m]\x1b[0m\n  \x1b[32mS\x1b[0m \x1b[35mslice_big\x1b[0m   : \x1b[33m1000\x1b[0m \x1b[32mslice[\x1b[0m\n      \x1b[32m0\x1b[0m: \x1b[34m1\x1b[0m\n      \x1b[32m1\x1b[0m: \x1b[34m2\x1b[0m\n      \x1b[32m2\x1b[0m: \x1b[34m3\x1b[0m\n      \x1b[32m3\x1b[0m: \x1b[34m4\x1b[0m\n           \x1b[34m...\x1b[0m\x1b[32m]\x1b[0m\n    \x1b[35mtest_string\x1b[0m : some string\n  \x1b[36m@\x1b[0m \x1b[35mtime\x1b[0m        : \x1b[36m2012-03-28 00:00:00 +0000 UTC\x1b[0m\n  \x1b[34m*\x1b[0m \x1b[35murl\x1b[0m         : \x1b[34mhttps://go.dev/\x1b[0m\n  \x1b[32mG\x1b[0m \x1b[35mmy_group\x1b[0m    : \x1b[32mgroup\x1b[0m\n    \x1b[33m#\x1b[0m \x1b[35mfloat\x1b[0m : \x1b[33m1.21\x1b[0m\n    \x1b[33m#\x1b[0m \x1b[35mint\x1b[0m   : \x1b[33m1\x1b[0m\n\n", timeString)

	if !bytes.Equal(w.WrittenData, []byte(expected)) {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, w.WrittenData)
	}
}

func Test_EmptyLogs(t *testing.T) {
	w := &MockWriter{}

	slogOpts := &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr { return a },
	}

	opts := &Options{
		HandlerOptions:    slogOpts,
		MaxSlicePrintSize: 4,
		SortKeys:          true,
		TimeFormat:        "[15:06]",
	}

	logger := slog.New(NewHandler(w, opts))

	timeString := csf(time.Now().Format("[15:06]"), fgWhite)
	_, filename, l1, _ := runtime.Caller(0)
	logger.Debug("My DEBUG message")
	_, _, l2, _ := runtime.Caller(0)
	logger.Info("My INFO message")
	_, _, l3, _ := runtime.Caller(0)
	logger.Warn("My WARN message")
	_, _, l4, _ := runtime.Caller(0)
	logger.Error("My ERROR message")

	expected := fmt.Sprintf("%[1]v \x1b[34m@@@\x1b[0m \x1b[4m\x1b[33m%[2]s\x1b[0m\x1b[0m:\x1b[31m%[3]v\x1b[0m\n\x1b[30m\x1b[44m DEBUG \x1b[0m \x1b[34mMy DEBUG message\x1b[0m\n\n%[1]v \x1b[34m@@@\x1b[0m \x1b[4m\x1b[33m%[2]v\x1b[0m\x1b[0m:\x1b[31m%[4]v\x1b[0m\n\x1b[30m\x1b[42m INFO \x1b[0m \x1b[32mMy INFO message\x1b[0m\n\n%[1]v \x1b[34m@@@\x1b[0m \x1b[4m\x1b[33m%[2]v\x1b[0m\x1b[0m:\x1b[31m%[5]v\x1b[0m\n\x1b[30m\x1b[43m WARN \x1b[0m \x1b[33mMy WARN message\x1b[0m\n\n%[1]v \x1b[34m@@@\x1b[0m \x1b[4m\x1b[33m%[2]v\x1b[0m\x1b[0m:\x1b[31m%[6]v\x1b[0m\n\x1b[30m\x1b[41m ERROR \x1b[0m \x1b[31mMy ERROR message\x1b[0m\n\n", timeString, filename, l1+1, l2+1, l3+1, l4+1)

	if !bytes.Equal(w.WrittenData, []byte(expected)) {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, w.WrittenData)
	}
}

func Test_WithGroups(t *testing.T) {
	w := &MockWriter{}

	slogOpts := &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}

	opts := &Options{
		HandlerOptions:    slogOpts,
		MaxSlicePrintSize: 4,
		SortKeys:          true,
		TimeFormat:        "[15:06]",
	}

	logger := slog.New(NewHandler(w, opts).WithGroup("test_group"))

	timeString := csf(time.Now().Format("[15:06]"), fgWhite)
	logger.Info("My INFO message")

	expected := fmt.Sprintf("%[1]v \x1b[30m\x1b[42m INFO \x1b[0m \x1b[32mMy INFO message\x1b[0m\n\n", timeString)

	if !bytes.Equal(w.WrittenData, []byte(expected)) {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, w.WrittenData)
	}
}

const (
	LevelTrace     = slog.Level(-8)
	LevelDebug     = slog.LevelDebug
	LevelInfo      = slog.LevelInfo
	LevelNotice    = slog.Level(2)
	LevelWarning   = slog.LevelWarn
	LevelError     = slog.LevelError
	LevelEmergency = slog.Level(12)
)

func Test_ReplaceLevelAttributes(t *testing.T) {
	w := &MockWriter{}

	slogOpts := &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelDebug,
		ReplaceAttr: replaceAttributes,
	}

	opts := &Options{
		HandlerOptions:    slogOpts,
		MaxSlicePrintSize: 4,
		SortKeys:          true,
		TimeFormat:        "[15:06]",
	}

	logger := slog.New(NewHandler(w, opts))

	timeString := csf(time.Now().Format("[15:06]"), fgWhite)
	ctx := context.Background()
	logger.Log(ctx, LevelEmergency, "missing pilots")
	logger.Error("failed to start engines", "err", "missing fuel")
	logger.Warn("falling back to default value")
	logger.Log(ctx, LevelNotice, "all systems are running")
	logger.Info("initiating launch")
	logger.Debug("starting background job")
	logger.Log(ctx, LevelTrace, "button clicked")

	expected := fmt.Sprintf("%[1]v \x1b[30m\x1b[41m EMERGENCY \x1b[0m \x1b[31mmissing pilots\x1b[0m\n  \x1b[35msev\x1b[0m : EMERGENCY\n\n%[1]v \x1b[30m\x1b[41m ERROR \x1b[0m \x1b[31mfailed to start engines\x1b[0m\n  \x1b[35merr\x1b[0m : missing fuel\n  \x1b[35msev\x1b[0m : ERROR\n\n%[1]v \x1b[30m\x1b[43m WARNING \x1b[0m \x1b[33mfalling back to default value\x1b[0m\n  \x1b[35msev\x1b[0m : WARNING\n\n%[1]v \x1b[30m\x1b[42m NOTICE \x1b[0m \x1b[32mall systems are running\x1b[0m\n  \x1b[35msev\x1b[0m : NOTICE\n\n%[1]v \x1b[30m\x1b[42m INFO \x1b[0m \x1b[32minitiating launch\x1b[0m\n  \x1b[35msev\x1b[0m : INFO\n\n%[1]v \x1b[30m\x1b[44m DEBUG \x1b[0m \x1b[34mstarting background job\x1b[0m\n  \x1b[35msev\x1b[0m : DEBUG\n\n", timeString)

	if !bytes.Equal(w.WrittenData, []byte(expected)) {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, w.WrittenData)
	}
}

func replaceAttributes(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		// Rename the level key from "level" to "sev".
		a.Key = "sev"

		// Handle custom level values.
		level := a.Value.Any().(slog.Level)

		// This could also look up the name from a map or other structure, but
		// this demonstrates using a switch statement to rename levels. For
		// maximum performance, the string values should be constants, but this
		// example uses the raw strings for readability.
		switch {
		case level < LevelDebug:
			a.Value = slog.StringValue("TRACE")
		case level < LevelInfo:
			a.Value = slog.StringValue("DEBUG")
		case level < LevelNotice:
			a.Value = slog.StringValue("INFO")
		case level < LevelWarning:
			a.Value = slog.StringValue("NOTICE")
		case level < LevelError:
			a.Value = slog.StringValue("WARNING")
		case level < LevelEmergency:
			a.Value = slog.StringValue("ERROR")
		default:
			a.Value = slog.StringValue("EMERGENCY")
		}
	}

	return a
}
