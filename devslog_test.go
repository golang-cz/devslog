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

func TestNewHandler(t *testing.T) {
	testNewHandlerDefaults(t)
	testNewHandlerWithOptions(t)
	testNewHandlerWithNilOptions(t)
	testNewHandlerWithNilSlogHandlerOptions(t)
}

func TestMethods(t *testing.T) {
	testEnabled(t)
	testWithGroup(t)
	testWithGroupEmpty(t)
	testWithAttrs(t)
	testWithAttrsEmpty(t)
}

func TestLevels(t *testing.T) {
	testLevelMessageDebug(t)
	testLevelMessageInfo(t)
	testLevelMessageWarn(t)
	testLevelMessageError(t)
}

func TestGroupsAndAttributes(t *testing.T) {
	testWithGroups(t)
	testWithGroupsEmpty(t)
	testWithAttributes(t)
	testWithAttributesRaceCondition()
}

func TestSourceAndReplace(t *testing.T) {
	testSource(t)
	testReplaceLevelAttributes(t)
}

func TestTypes(t *testing.T) {
	slogOpts := &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr { return a },
	}

	opts := &Options{
		HandlerOptions:    slogOpts,
		MaxSlicePrintSize: 4,
		SortKeys:          true,
		TimeFormat:        "[]",
		NewLineAfterLog:   true,
		StringerFormatter: true,
	}

	testString(t, opts)
	testIntFloat(t, opts)
	testBool(t, opts)
	testTime(t, opts)
	testError(t, opts)
	testSlice(t, opts)
	testSliceBig(t, opts)
	testMap(t, opts)
	testMapOfPointers(t, opts)
	testMapOfInterface(t, opts)
	testStruct(t, opts)
	testNilInterface(t, opts)
	testGroup(t, opts)
	testLogValuer(t, opts)
	testLogValuerPanic(t, opts)
	testStringer(t, opts)
	testStringerInner(t, opts)
	testNoColor(t, opts)
}

func testNewHandlerDefaults(t *testing.T) {
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

	if h.opts.TimeFormat != "[15:04:05]" {
		t.Errorf("Expected default TimeFormat to be \"[15:04:05]\" ")
	}

	if h.out == nil {
		t.Errorf("Expected writer to be initialized")
	}
}

func testNewHandlerWithOptions(t *testing.T) {
	handlerOpts := &Options{
		HandlerOptions:    &slog.HandlerOptions{Level: slog.LevelWarn},
		MaxSlicePrintSize: 10,
		TimeFormat:        "[04:05]",
	}
	h := NewHandler(nil, handlerOpts)

	if h.opts.Level.Level() != slog.LevelWarn.Level() {
		t.Errorf("Expected custom log level to be LevelWarn")
	}

	if h.opts.MaxSlicePrintSize != 10 {
		t.Errorf("Expected custom MaxSlicePrintSize to be 10")
	}

	if h.opts.TimeFormat != "[04:05]" {
		t.Errorf("Expected custom TimeFormat to be \"[04:05]\" ")
	}
}

func testNewHandlerWithNilOptions(t *testing.T) {
	h := NewHandler(nil, nil)

	if h.opts.HandlerOptions == nil || h.opts.Level != slog.LevelInfo {
		t.Errorf("Expected HandlerOptions to be initialized with default level")
	}

	if h.opts.MaxSlicePrintSize != 50 {
		t.Errorf("Expected MaxSlicePrintSize to be initialized with default value")
	}

	if h.out != nil {
		t.Errorf("Expected writer to be nil")
	}
}

func testNewHandlerWithNilSlogHandlerOptions(t *testing.T) {
	opts := &Options{}
	h := NewHandler(nil, opts)

	if h.opts.HandlerOptions == nil || h.opts.Level != slog.LevelInfo {
		t.Errorf("Expected HandlerOptions to be initialized with default level")
	}

	if h.opts.MaxSlicePrintSize != 50 {
		t.Errorf("Expected MaxSlicePrintSize to be initialized with default value")
	}

	if h.out != nil {
		t.Errorf("Expected writer to be nil")
	}
}

func testEnabled(t *testing.T) {
	h := NewHandler(nil, nil)
	ctx := context.Background()

	if !h.Enabled(ctx, slog.LevelInfo) {
		t.Error("Expected handler to be enabled for LevelInfo")
	}

	if h.Enabled(ctx, slog.LevelDebug) {
		t.Error("Expected handler to be disabled for LevelDebug")
	}
}

func testWithGroup(t *testing.T) {
	h := NewHandler(nil, nil)
	h2 := h.WithGroup("myGroup")

	if h2 == h {
		t.Error("Expected a new handler instance")
	}
}

func testWithGroupEmpty(t *testing.T) {
	h := NewHandler(nil, nil)
	h2 := h.WithGroup("")

	if h2 != h {
		t.Error("Expected a original handler instance")
	}
}

func testWithAttrs(t *testing.T) {
	h := NewHandler(nil, nil)
	h2 := h.WithAttrs([]slog.Attr{slog.String("key", "value")})

	if h2 == h {
		t.Error("Expected a new handler instance")
	}
}

func testWithAttrsEmpty(t *testing.T) {
	h := NewHandler(nil, nil)
	h2 := h.WithAttrs([]slog.Attr{})

	if h2 != h {
		t.Error("Expected a original handler instance")
	}
}

func testLevelMessageDebug(t *testing.T) {
	h := NewHandler(nil, nil)
	buf := make([]byte, 0)
	record := &slog.Record{
		Level:   slog.LevelDebug,
		Message: "Debug message",
	}

	buf = h.levelMessage(buf, record)

	expected := "\x1b[44m\x1b[30m DEBUG \x1b[0m \x1b[34mDebug message\x1b[0m\n"
	result := string(buf)

	if result != expected {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, result)
	}
}

func testLevelMessageInfo(t *testing.T) {
	h := NewHandler(nil, nil)
	buf := make([]byte, 0)
	record := &slog.Record{
		Level:   slog.LevelInfo,
		Message: "Info message",
	}

	buf = h.levelMessage(buf, record)

	expected := "\x1b[42m\x1b[30m INFO \x1b[0m \x1b[32mInfo message\x1b[0m\n"
	result := string(buf)

	if result != expected {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, result)
	}
}

func testLevelMessageWarn(t *testing.T) {
	h := NewHandler(nil, nil)
	buf := make([]byte, 0)
	record := &slog.Record{
		Level:   slog.LevelWarn,
		Message: "Warning message",
	}

	buf = h.levelMessage(buf, record)

	expected := "\x1b[43m\x1b[30m WARN \x1b[0m \x1b[33mWarning message\x1b[0m\n"
	result := string(buf)

	if result != expected {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, result)
	}
}

func testLevelMessageError(t *testing.T) {
	h := NewHandler(nil, nil)
	buf := make([]byte, 0)
	record := &slog.Record{
		Level:   slog.LevelError,
		Message: "Error message",
	}

	buf = h.levelMessage(buf, record)

	expected := "\x1b[41m\x1b[30m ERROR \x1b[0m \x1b[31mError message\x1b[0m\n"
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

func testSource(t *testing.T) {
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
		TimeFormat:        "[15:04]",
		NewLineAfterLog:   true,
	}

	h := NewHandler(w, opts)
	logger := slog.New(h)

	timeString := h.csf([]byte(time.Now().Format("[15:04]")), fgWhite)
	_, filename, l, _ := runtime.Caller(0)
	logger.Info("message")

	expected := fmt.Sprintf("%1s \x1b[34m@@@\x1b[0m \x1b[4m\x1b[33m%2s\x1b[0m\x1b[0m:\x1b[31m%v\x1b[0m\n\x1b[42m\x1b[30m INFO \x1b[0m \x1b[32mmessage\x1b[0m\n\n", timeString, filename, l+1)

	if !bytes.Equal(w.WrittenData, []byte(expected)) {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, w.WrittenData)
	}
}

func testWithGroups(t *testing.T) {
	w := &MockWriter{}

	slogOpts := &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}

	opts := &Options{
		HandlerOptions:    slogOpts,
		MaxSlicePrintSize: 4,
		SortKeys:          true,
		TimeFormat:        "[]",
		NewLineAfterLog:   true,
	}

	logger := slog.New(NewHandler(w, opts).WithGroup("test_group"))

	logger.Info("My INFO message",
		slog.Any("a", "1"),
	)

	expected := "\x1b[2m\x1b[37m[]\x1b[0m \x1b[42m\x1b[30m INFO \x1b[0m \x1b[32mMy INFO message\x1b[0m\n\x1b[32mG\x1b[0m \x1b[35mtest_group\x1b[0m: \n    \x1b[35ma\x1b[0m: 1\n\n"

	if !bytes.Equal(w.WrittenData, []byte(expected)) {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, w.WrittenData)
	}
}

func testWithGroupsEmpty(t *testing.T) {
	w := &MockWriter{}

	slogOpts := &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}

	opts := &Options{
		HandlerOptions:    slogOpts,
		MaxSlicePrintSize: 4,
		SortKeys:          true,
		TimeFormat:        "[]",
		NewLineAfterLog:   true,
	}

	logger := slog.New(NewHandler(w, opts).WithGroup("test_group"))

	logger.Info("My INFO message")

	expected := "\x1b[2m\x1b[37m[]\x1b[0m \x1b[42m\x1b[30m INFO \x1b[0m \x1b[32mMy INFO message\x1b[0m\n\n"

	if !bytes.Equal(w.WrittenData, []byte(expected)) {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, w.WrittenData)
	}
}

func testWithAttributes(t *testing.T) {
	w := &MockWriter{}

	slogOpts := &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}

	opts := &Options{
		HandlerOptions:    slogOpts,
		MaxSlicePrintSize: 4,
		SortKeys:          true,
		TimeFormat:        "[]",
		NewLineAfterLog:   true,
	}

	as := []slog.Attr{slog.Any("a", "1")}
	logger := slog.New(NewHandler(w, opts).WithAttrs(as))

	logger.Info("My INFO message")

	expected := "\x1b[2m\x1b[37m[]\x1b[0m \x1b[42m\x1b[30m INFO \x1b[0m \x1b[32mMy INFO message\x1b[0m\n  \x1b[35ma\x1b[0m: 1\n\n"

	if !bytes.Equal(w.WrittenData, []byte(expected)) {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, w.WrittenData)
	}
}

func testWithAttributesRaceCondition() {
	w := &MockWriter{}

	slogOpts := &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}

	opts := &Options{
		HandlerOptions:    slogOpts,
		MaxSlicePrintSize: 4,
		SortKeys:          true,
		TimeFormat:        "[]",
		NewLineAfterLog:   true,
	}

	logger := slog.New(NewHandler(w, opts))

	go func() {
		as := []slog.Attr{slog.Any("a", "1")}
		logger.Handler().WithAttrs(as)
	}()

	go func() {
		logger.Info("INFO message")
	}()
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

func testReplaceLevelAttributes(t *testing.T) {
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
		TimeFormat:        "[15:04]",
		NewLineAfterLog:   true,
	}

	h := NewHandler(w, opts)
	logger := slog.New(h)

	timeString := h.csf([]byte(time.Now().Format("[15:04]")), fgWhite)
	ctx := context.Background()
	logger.Log(ctx, LevelEmergency, "missing pilots")
	logger.Error("failed to start engines", "err", "missing fuel")
	logger.Warn("falling back to default value")
	logger.Log(ctx, LevelNotice, "all systems are running")
	logger.Info("initiating launch")
	logger.Debug("starting background job")
	logger.Log(ctx, LevelTrace, "button clicked")

	expected := fmt.Sprintf(
		"%[1]s \x1b[41m\x1b[30m EMERGENCY \x1b[0m \x1b[31mmissing pilots\x1b[0m\n  \x1b[35msev\x1b[0m: EMERGENCY\n\n%[1]s \x1b[41m\x1b[30m ERROR \x1b[0m \x1b[31mfailed to start engines\x1b[0m\n  \x1b[35merr\x1b[0m: missing fuel\n  \x1b[35msev\x1b[0m: ERROR\n\n%[1]s \x1b[43m\x1b[30m WARNING \x1b[0m \x1b[33mfalling back to default value\x1b[0m\n  \x1b[35msev\x1b[0m: WARNING\n\n%[1]s \x1b[42m\x1b[30m NOTICE \x1b[0m \x1b[32mall systems are running\x1b[0m\n  \x1b[35msev\x1b[0m: NOTICE\n\n%[1]s \x1b[42m\x1b[30m INFO \x1b[0m \x1b[32minitiating launch\x1b[0m\n  \x1b[35msev\x1b[0m: INFO\n\n%[1]s \x1b[44m\x1b[30m DEBUG \x1b[0m \x1b[34mstarting background job\x1b[0m\n  \x1b[35msev\x1b[0m: DEBUG\n\n",
		timeString,
	)

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

func testString(t *testing.T, o *Options) {
	w := &MockWriter{}
	logger := slog.New(NewHandler(w, o))

	s := "string"

	logger.Info("msg",
		slog.Any("s", s),
		slog.Any("sp", &s),
		slog.Any("empty", ""),
		slog.Any("url", "https://go.dev/"),
	)

	expected := []byte(
		"\x1b[2m\x1b[37m[]\x1b[0m \x1b[42m\x1b[30m INFO \x1b[0m \x1b[32mmsg\x1b[0m\n  \x1b[35mempty\x1b[0m: \x1b[2m\x1b[37mempty\x1b[0m\n  \x1b[35ms\x1b[0m    : string\n  \x1b[35msp\x1b[0m   : string\n\x1b[34m*\x1b[0m \x1b[35murl\x1b[0m  : \x1b[4m\x1b[34mhttps://go.dev/\x1b[0m\x1b[0m\n\n",
	)

	if !bytes.Equal(w.WrittenData, expected) {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, w.WrittenData)
	}
}

func testIntFloat(t *testing.T, o *Options) {
	w := &MockWriter{}
	logger := slog.New(NewHandler(w, o))

	f := 1.21
	fp := &f
	i := 1
	ip := &i
	logger.Info("msg",
		slog.Any("f", f),
		slog.Any("fp", fp),
		slog.Any("i", i),
		slog.Any("ip", ip),
	)

	expected := fmt.Sprintf(
		"\x1b[2m\x1b[37m[]\x1b[0m \x1b[42m\x1b[30m INFO \x1b[0m \x1b[32mmsg\x1b[0m\n\x1b[33m#\x1b[0m \x1b[35mf\x1b[0m : \x1b[33m1.21\x1b[0m\n\x1b[33m#\x1b[0m \x1b[35mfp\x1b[0m: \x1b[31m*\x1b[0m\x1b[33m1.21\x1b[0m\x1b[2m\x1b[37m \"%v\"\x1b[0m\n\x1b[33m#\x1b[0m \x1b[35mi\x1b[0m : \x1b[33m1\x1b[0m\n\x1b[33m#\x1b[0m \x1b[35mip\x1b[0m: \x1b[31m*\x1b[0m\x1b[33m1\x1b[0m\x1b[2m\x1b[37m \"%v\"\x1b[0m\n\n",
		fp,
		ip,
	)

	if !bytes.Equal(w.WrittenData, []byte(expected)) {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, w.WrittenData)
	}
}

func testBool(t *testing.T, o *Options) {
	w := &MockWriter{}
	logger := slog.New(NewHandler(w, o))

	b := true
	bp := &b
	logger.Info("msg",
		slog.Any("b", b),
		slog.Any("bp", bp),
	)

	expected := fmt.Sprintf("\x1b[2m\x1b[37m[]\x1b[0m \x1b[42m\x1b[30m INFO \x1b[0m \x1b[32mmsg\x1b[0m\n\x1b[31m#\x1b[0m \x1b[35mb\x1b[0m : \x1b[31mtrue\x1b[0m\n\x1b[31m#\x1b[0m \x1b[35mbp\x1b[0m: \x1b[31m*\x1b[0m\x1b[31mtrue\x1b[0m\x1b[2m\x1b[37m \"%v\"\x1b[0m\n\n", bp)

	if !bytes.Equal(w.WrittenData, []byte(expected)) {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, w.WrittenData)
	}
}

func testTime(t *testing.T, o *Options) {
	w := &MockWriter{}
	logger := slog.New(NewHandler(w, o))

	timeT := time.Date(2012, time.March, 28, 0, 0, 0, 0, time.UTC)
	timeE := time.Date(2023, time.August, 15, 12, 0, 0, 0, time.UTC)
	timeD := timeE.Sub(timeT)

	logger.Info("msg",
		slog.Any("t", timeT),
		slog.Any("tp", &timeT),
		slog.Any("d", timeD),
		slog.Any("tp", &timeD),
	)

	expected := []byte(
		"\x1b[2m\x1b[37m[]\x1b[0m \x1b[42m\x1b[30m INFO \x1b[0m \x1b[32mmsg\x1b[0m\n\x1b[36m@\x1b[0m \x1b[35md\x1b[0m : \x1b[36m99780h0m0s\x1b[0m\n\x1b[36m@\x1b[0m \x1b[35mt\x1b[0m : \x1b[36m2012-03-28 00:00:00 +0000 UTC\x1b[0m\n\x1b[36m@\x1b[0m \x1b[35mtp\x1b[0m: \x1b[36m2012-03-28 00:00:00 +0000 UTC\x1b[0m\n\x1b[36m@\x1b[0m \x1b[35mtp\x1b[0m: \x1b[36m99780h0m0s\x1b[0m\n\n",
	)

	if !bytes.Equal(w.WrittenData, expected) {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, w.WrittenData)
	}
}

func testError(t *testing.T, o *Options) {
	w := &MockWriter{}
	logger := slog.New(NewHandler(w, o))

	e := fmt.Errorf("broken")
	e = fmt.Errorf("err 1: %w", e)
	e = fmt.Errorf("err 2: %w", e)

	logger.Info("msg",
		slog.Any("e", e),
	)

	expected := []byte(
		"\x1b[2m\x1b[37m[]\x1b[0m \x1b[42m\x1b[30m INFO \x1b[0m \x1b[32mmsg\x1b[0m\n\x1b[31mE\x1b[0m \x1b[35me\x1b[0m: \n    \x1b[31m0\x1b[0m\x1b[37m: \x1b[0m\x1b[31merr 2\x1b[0m\n    \x1b[31m1\x1b[0m\x1b[37m: \x1b[0m\x1b[31merr 1\x1b[0m\n    \x1b[31m2\x1b[0m\x1b[37m: \x1b[0m\x1b[31mbroken\x1b[0m\n\n",
	)

	if !bytes.Equal(w.WrittenData, expected) {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, w.WrittenData)
	}
}

func testSlice(t *testing.T, o *Options) {
	w := &MockWriter{}
	logger := slog.New(NewHandler(w, o))

	s := []string{"apple", "ba na na"}

	logger.Info("msg",
		slog.Any("s", s),
	)

	expected := []byte(
		"\x1b[2m\x1b[37m[]\x1b[0m \x1b[42m\x1b[30m INFO \x1b[0m \x1b[32mmsg\x1b[0m\n\x1b[32mS\x1b[0m \x1b[35ms\x1b[0m: \x1b[34m2\x1b[0m \x1b[32m[\x1b[0m\x1b[32m]\x1b[0m\x1b[33ms\x1b[0m\x1b[33mt\x1b[0m\x1b[33mr\x1b[0m\x1b[33mi\x1b[0m\x1b[33mn\x1b[0m\x1b[33mg\x1b[0m\n    \x1b[32m0\x1b[0m: apple\n    \x1b[32m1\x1b[0m: ba na na\n\n",
	)

	if !bytes.Equal(w.WrittenData, expected) {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, w.WrittenData)
	}
}

func testSliceBig(t *testing.T, o *Options) {
	w := &MockWriter{}
	logger := slog.New(NewHandler(w, o))

	s := make([]int, 0)
	for i := 0; i < 11; i++ {
		s = append(s, i*2)
	}

	logger.Info("msg",
		slog.Any("s", s),
	)

	expected := []byte(
		"\x1b[2m\x1b[37m[]\x1b[0m \x1b[42m\x1b[30m INFO \x1b[0m \x1b[32mmsg\x1b[0m\n\x1b[32mS\x1b[0m \x1b[35ms\x1b[0m: \x1b[34m11\x1b[0m \x1b[32m[\x1b[0m\x1b[32m]\x1b[0m\x1b[33mi\x1b[0m\x1b[33mn\x1b[0m\x1b[33mt\x1b[0m\n    \x1b[32m0\x1b[0m: \x1b[33m0\x1b[0m\n    \x1b[32m1\x1b[0m: \x1b[33m2\x1b[0m\n    \x1b[32m2\x1b[0m: \x1b[33m4\x1b[0m\n    \x1b[32m3\x1b[0m: \x1b[33m6\x1b[0m\n       \x1b[34m...\x1b[0m\x1b[32m]\x1b[0m\n\n",
	)

	if !bytes.Equal(w.WrittenData, expected) {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, w.WrittenData)
	}
}

func testMap(t *testing.T, o *Options) {
	w := &MockWriter{}
	logger := slog.New(NewHandler(w, o))

	m := map[int]string{0: "a", 1: "b"}
	mp := &m

	logger.Info("msg",
		slog.Any("m", m),
		slog.Any("mp", mp),
		slog.Any("mpp", &mp),
	)

	expected := []byte(
		"\x1b[2m\x1b[37m[]\x1b[0m \x1b[42m\x1b[30m INFO \x1b[0m \x1b[32mmsg\x1b[0m\n\x1b[32mM\x1b[0m \x1b[35mm\x1b[0m  : \x1b[34m2\x1b[0m \x1b[33mm\x1b[0m\x1b[33ma\x1b[0m\x1b[33mp\x1b[0m\x1b[32m[\x1b[0m\x1b[33mi\x1b[0m\x1b[33mn\x1b[0m\x1b[33mt\x1b[0m\x1b[32m]\x1b[0m\x1b[33ms\x1b[0m\x1b[33mt\x1b[0m\x1b[33mr\x1b[0m\x1b[33mi\x1b[0m\x1b[33mn\x1b[0m\x1b[33mg\x1b[0m\n    \x1b[32m0\x1b[0m: a\n    \x1b[32m1\x1b[0m: b\n\x1b[32mM\x1b[0m \x1b[35mmp\x1b[0m : \x1b[34m2\x1b[0m \x1b[31m*\x1b[0m\x1b[33mm\x1b[0m\x1b[33ma\x1b[0m\x1b[33mp\x1b[0m\x1b[32m[\x1b[0m\x1b[33mi\x1b[0m\x1b[33mn\x1b[0m\x1b[33mt\x1b[0m\x1b[32m]\x1b[0m\x1b[33ms\x1b[0m\x1b[33mt\x1b[0m\x1b[33mr\x1b[0m\x1b[33mi\x1b[0m\x1b[33mn\x1b[0m\x1b[33mg\x1b[0m\n    \x1b[32m0\x1b[0m: a\n    \x1b[32m1\x1b[0m: b\n\x1b[32mM\x1b[0m \x1b[35mmpp\x1b[0m: \x1b[34m2\x1b[0m \x1b[31m*\x1b[0m\x1b[31m*\x1b[0m\x1b[33mm\x1b[0m\x1b[33ma\x1b[0m\x1b[33mp\x1b[0m\x1b[32m[\x1b[0m\x1b[33mi\x1b[0m\x1b[33mn\x1b[0m\x1b[33mt\x1b[0m\x1b[32m]\x1b[0m\x1b[33ms\x1b[0m\x1b[33mt\x1b[0m\x1b[33mr\x1b[0m\x1b[33mi\x1b[0m\x1b[33mn\x1b[0m\x1b[33mg\x1b[0m\n    \x1b[32m0\x1b[0m: a\n    \x1b[32m1\x1b[0m: b\n\n",
	)

	if !bytes.Equal(w.WrittenData, expected) {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, w.WrittenData)
	}
}

func testMapOfPointers(t *testing.T, o *Options) {
	w := &MockWriter{}
	logger := slog.New(NewHandler(w, o))

	s := "a"
	m := map[int]*string{0: &s, 1: &s}

	logger.Info("msg",
		slog.Any("m", m),
	)

	expected := []byte(
		"\x1b[2m\x1b[37m[]\x1b[0m \x1b[42m\x1b[30m INFO \x1b[0m \x1b[32mmsg\x1b[0m\n\x1b[32mM\x1b[0m \x1b[35mm\x1b[0m: \x1b[34m2\x1b[0m \x1b[33mm\x1b[0m\x1b[33ma\x1b[0m\x1b[33mp\x1b[0m\x1b[32m[\x1b[0m\x1b[33mi\x1b[0m\x1b[33mn\x1b[0m\x1b[33mt\x1b[0m\x1b[32m]\x1b[0m\x1b[31m*\x1b[0m\x1b[33ms\x1b[0m\x1b[33mt\x1b[0m\x1b[33mr\x1b[0m\x1b[33mi\x1b[0m\x1b[33mn\x1b[0m\x1b[33mg\x1b[0m\n    \x1b[32m0\x1b[0m: a\n    \x1b[32m1\x1b[0m: a\n\n",
	)

	if !bytes.Equal(w.WrittenData, expected) {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, w.WrittenData)
	}
}

func testMapOfInterface(t *testing.T, o *Options) {
	w := &MockWriter{}
	logger := slog.New(NewHandler(w, o))

	m := map[int]any{0: "a", 1: "b"}
	mp := &m

	logger.Info("msg",
		slog.Any("m", m),
		slog.Any("mp", mp),
		slog.Any("mpp", &mp),
	)

	expected := []byte(
		"\x1b[2m\x1b[37m[]\x1b[0m \x1b[42m\x1b[30m INFO \x1b[0m \x1b[32mmsg\x1b[0m\n\x1b[32mM\x1b[0m \x1b[35mm\x1b[0m  : \x1b[34m2\x1b[0m \x1b[33mm\x1b[0m\x1b[33ma\x1b[0m\x1b[33mp\x1b[0m\x1b[32m[\x1b[0m\x1b[33mi\x1b[0m\x1b[33mn\x1b[0m\x1b[33mt\x1b[0m\x1b[32m]\x1b[0m\x1b[33mi\x1b[0m\x1b[33mn\x1b[0m\x1b[33mt\x1b[0m\x1b[33me\x1b[0m\x1b[33mr\x1b[0m\x1b[33mf\x1b[0m\x1b[33ma\x1b[0m\x1b[33mc\x1b[0m\x1b[33me\x1b[0m\x1b[33m \x1b[0m\x1b[33m{\x1b[0m\x1b[33m}\x1b[0m\n    \x1b[32m0\x1b[0m: a\n    \x1b[32m1\x1b[0m: b\n\x1b[32mM\x1b[0m \x1b[35mmp\x1b[0m : \x1b[34m2\x1b[0m \x1b[31m*\x1b[0m\x1b[33mm\x1b[0m\x1b[33ma\x1b[0m\x1b[33mp\x1b[0m\x1b[32m[\x1b[0m\x1b[33mi\x1b[0m\x1b[33mn\x1b[0m\x1b[33mt\x1b[0m\x1b[32m]\x1b[0m\x1b[33mi\x1b[0m\x1b[33mn\x1b[0m\x1b[33mt\x1b[0m\x1b[33me\x1b[0m\x1b[33mr\x1b[0m\x1b[33mf\x1b[0m\x1b[33ma\x1b[0m\x1b[33mc\x1b[0m\x1b[33me\x1b[0m\x1b[33m \x1b[0m\x1b[33m{\x1b[0m\x1b[33m}\x1b[0m\n    \x1b[32m0\x1b[0m: a\n    \x1b[32m1\x1b[0m: b\n\x1b[32mM\x1b[0m \x1b[35mmpp\x1b[0m: \x1b[34m2\x1b[0m \x1b[31m*\x1b[0m\x1b[31m*\x1b[0m\x1b[33mm\x1b[0m\x1b[33ma\x1b[0m\x1b[33mp\x1b[0m\x1b[32m[\x1b[0m\x1b[33mi\x1b[0m\x1b[33mn\x1b[0m\x1b[33mt\x1b[0m\x1b[32m]\x1b[0m\x1b[33mi\x1b[0m\x1b[33mn\x1b[0m\x1b[33mt\x1b[0m\x1b[33me\x1b[0m\x1b[33mr\x1b[0m\x1b[33mf\x1b[0m\x1b[33ma\x1b[0m\x1b[33mc\x1b[0m\x1b[33me\x1b[0m\x1b[33m \x1b[0m\x1b[33m{\x1b[0m\x1b[33m}\x1b[0m\n    \x1b[32m0\x1b[0m: a\n    \x1b[32m1\x1b[0m: b\n\n",
	)

	if !bytes.Equal(w.WrittenData, expected) {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, w.WrittenData)
	}
}

func testStruct(t *testing.T, o *Options) {
	w := &MockWriter{}
	logger := slog.New(NewHandler(w, o))

	type StructTest struct {
		Slice      []int
		Map        map[int]int
		Struct     struct{ B bool }
		SliceP     *[]int
		MapP       *map[int]int
		StructP    *struct{ B bool }
		unexported int
	}

	s := &StructTest{
		Slice:      []int{},
		Map:        map[int]int{},
		Struct:     struct{ B bool }{},
		SliceP:     &[]int{},
		MapP:       &map[int]int{},
		StructP:    &struct{ B bool }{},
		unexported: 5,
	}

	logger.Info("msg",
		slog.Any("s", s),
	)

	expected := []byte(
		"\x1b[2m\x1b[37m[]\x1b[0m \x1b[42m\x1b[30m INFO \x1b[0m \x1b[32mmsg\x1b[0m\n\x1b[33mS\x1b[0m \x1b[35ms\x1b[0m: \x1b[31m*\x1b[0m\x1b[33md\x1b[0m\x1b[33me\x1b[0m\x1b[33mv\x1b[0m\x1b[33ms\x1b[0m\x1b[33ml\x1b[0m\x1b[33mo\x1b[0m\x1b[33mg\x1b[0m\x1b[33m.\x1b[0m\x1b[33mS\x1b[0m\x1b[33mt\x1b[0m\x1b[33mr\x1b[0m\x1b[33mu\x1b[0m\x1b[33mc\x1b[0m\x1b[33mt\x1b[0m\x1b[33mT\x1b[0m\x1b[33me\x1b[0m\x1b[33ms\x1b[0m\x1b[33mt\x1b[0m\n    \x1b[32mSlice\x1b[0m  : \x1b[34m0\x1b[0m \x1b[32m[\x1b[0m\x1b[32m]\x1b[0m\x1b[33mi\x1b[0m\x1b[33mn\x1b[0m\x1b[33mt\x1b[0m\n    \x1b[32mMap\x1b[0m    : \x1b[34m0\x1b[0m \x1b[33mm\x1b[0m\x1b[33ma\x1b[0m\x1b[33mp\x1b[0m\x1b[32m[\x1b[0m\x1b[33mi\x1b[0m\x1b[33mn\x1b[0m\x1b[33mt\x1b[0m\x1b[32m]\x1b[0m\x1b[33mi\x1b[0m\x1b[33mn\x1b[0m\x1b[33mt\x1b[0m\n    \x1b[32mStruct\x1b[0m : \x1b[33ms\x1b[0m\x1b[33mt\x1b[0m\x1b[33mr\x1b[0m\x1b[33mu\x1b[0m\x1b[33mc\x1b[0m\x1b[33mt\x1b[0m\x1b[33m \x1b[0m\x1b[33m{\x1b[0m\x1b[33m \x1b[0m\x1b[33mB\x1b[0m\x1b[33m \x1b[0m\x1b[33mb\x1b[0m\x1b[33mo\x1b[0m\x1b[33mo\x1b[0m\x1b[33ml\x1b[0m\x1b[33m \x1b[0m\x1b[33m}\x1b[0m\n      \x1b[32mB\x1b[0m: \x1b[31mfalse\x1b[0m\n    \x1b[32mSliceP\x1b[0m : \x1b[34m0\x1b[0m \x1b[31m*\x1b[0m\x1b[32m[\x1b[0m\x1b[32m]\x1b[0m\x1b[33mi\x1b[0m\x1b[33mn\x1b[0m\x1b[33mt\x1b[0m\n    \x1b[32mMapP\x1b[0m   : \x1b[34m0\x1b[0m \x1b[31m*\x1b[0m\x1b[33mm\x1b[0m\x1b[33ma\x1b[0m\x1b[33mp\x1b[0m\x1b[32m[\x1b[0m\x1b[33mi\x1b[0m\x1b[33mn\x1b[0m\x1b[33mt\x1b[0m\x1b[32m]\x1b[0m\x1b[33mi\x1b[0m\x1b[33mn\x1b[0m\x1b[33mt\x1b[0m\n    \x1b[32mStructP\x1b[0m: \x1b[31m*\x1b[0m\x1b[33ms\x1b[0m\x1b[33mt\x1b[0m\x1b[33mr\x1b[0m\x1b[33mu\x1b[0m\x1b[33mc\x1b[0m\x1b[33mt\x1b[0m\x1b[33m \x1b[0m\x1b[33m{\x1b[0m\x1b[33m \x1b[0m\x1b[33mB\x1b[0m\x1b[33m \x1b[0m\x1b[33mb\x1b[0m\x1b[33mo\x1b[0m\x1b[33mo\x1b[0m\x1b[33ml\x1b[0m\x1b[33m \x1b[0m\x1b[33m}\x1b[0m\n      \x1b[32mB\x1b[0m: \x1b[31mfalse\x1b[0m\n\n",
	)

	if !bytes.Equal(w.WrittenData, expected) {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, w.WrittenData)
	}
}

func testNilInterface(t *testing.T, o *Options) {
	w := &MockWriter{}
	logger := slog.New(NewHandler(w, o))

	type StructWithInterface struct {
		Data any
	}

	s := StructWithInterface{}

	logger.Info("msg",
		slog.Any("s", s),
	)

	expected := []byte(
		"\x1b[2m\x1b[37m[]\x1b[0m \x1b[42m\x1b[30m INFO \x1b[0m \x1b[32mmsg\x1b[0m\n\x1b[33mS\x1b[0m \x1b[35ms\x1b[0m: \x1b[33md\x1b[0m\x1b[33me\x1b[0m\x1b[33mv\x1b[0m\x1b[33ms\x1b[0m\x1b[33ml\x1b[0m\x1b[33mo\x1b[0m\x1b[33mg\x1b[0m\x1b[33m.\x1b[0m\x1b[33mS\x1b[0m\x1b[33mt\x1b[0m\x1b[33mr\x1b[0m\x1b[33mu\x1b[0m\x1b[33mc\x1b[0m\x1b[33mt\x1b[0m\x1b[33mW\x1b[0m\x1b[33mi\x1b[0m\x1b[33mt\x1b[0m\x1b[33mh\x1b[0m\x1b[33mI\x1b[0m\x1b[33mn\x1b[0m\x1b[33mt\x1b[0m\x1b[33me\x1b[0m\x1b[33mr\x1b[0m\x1b[33mf\x1b[0m\x1b[33ma\x1b[0m\x1b[33mc\x1b[0m\x1b[33me\x1b[0m\n    \x1b[32mData\x1b[0m: \x1b[31m<\x1b[0m\x1b[33mnil\x1b[0m\x1b[31m>\x1b[0m\n\n",
	)

	if !bytes.Equal(w.WrittenData, expected) {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, w.WrittenData)
	}
}

func testGroup(t *testing.T, o *Options) {
	w := &MockWriter{}
	logger := slog.New(NewHandler(w, o))

	logger.Info("msg",
		slog.Any("1", "a"),
		slog.Group("g",
			slog.Any("2", "b"),
		),
	)

	expected := []byte("\x1b[2m\x1b[37m[]\x1b[0m \x1b[42m\x1b[30m INFO \x1b[0m \x1b[32mmsg\x1b[0m\n  \x1b[35m1\x1b[0m: a\n\x1b[32mG\x1b[0m \x1b[35mg\x1b[0m: \n    \x1b[35m2\x1b[0m: b\n\n")

	if !bytes.Equal(w.WrittenData, expected) {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, w.WrittenData)
	}
}

type logValuerExample1 struct {
	A int
	B string
}

func (item logValuerExample1) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("A", item.A),
		slog.String("B", item.B),
	)
}

func testLogValuer(t *testing.T, o *Options) {
	w := &MockWriter{}
	logger := slog.New(NewHandler(w, o))
	item1 := logValuerExample1{
		A: 5,
		B: "test",
	}
	logger.Info("test_log_valuer",
		slog.Any("item1", item1),
	)

	expected := []byte("\x1b[2m\x1b[37m[]\x1b[0m \x1b[42m\x1b[30m INFO \x1b[0m \x1b[32mtest_log_valuer\x1b[0m\n\x1b[32mG\x1b[0m \x1b[35mitem1\x1b[0m: \n  \x1b[33m#\x1b[0m \x1b[35mA\x1b[0m: \x1b[33m5\x1b[0m\n    \x1b[35mB\x1b[0m: test\n\n")

	if !bytes.Equal(w.WrittenData, expected) {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, w.WrittenData)
	}
}

type logValuerExample2 struct {
	A int
	B string
}

func (item logValuerExample2) LogValue() slog.Value {
	panic("log valuer paniced")
}

func testLogValuerPanic(t *testing.T, o *Options) {
	w := &MockWriter{}
	logger := slog.New(NewHandler(w, o))
	item1 := logValuerExample2{
		A: 5,
		B: "test",
	}
	logger.Info("test_log_valuer_panic",
		slog.Any("item1", item1),
	)

	expectedPrefix := []byte("\x1b[2m\x1b[37m[]\x1b[0m \x1b[42m\x1b[30m INFO \x1b[0m \x1b[32mtest_log_valuer_panic\x1b[0m\n\x1b[31mE\x1b[0m \x1b[35mitem1\x1b[0m: \n    \x1b[31m0\x1b[0m\x1b[37m: \x1b[0m\x1b[31mLogValue panicked\n")
	if !bytes.HasPrefix(w.WrittenData, expectedPrefix) {
		t.Errorf("\nGot:\n%s\n , %[1]q expected it to contain panic stack trace", w.WrittenData)
	}
}

type logStringerExample1 struct {
	A []byte
}

func (item logStringerExample1) String() string {
	return fmt.Sprintf("A: %s", item.A)
}

func testStringer(t *testing.T, o *Options) {
	w := &MockWriter{}
	logger := slog.New(NewHandler(w, o))
	item1 := logStringerExample1{
		A: []byte("test"),
	}
	logger.Info("test_stringer",
		slog.Any("item1", item1),
	)

	expected := []byte("\x1b[2m\x1b[37m[]\x1b[0m \x1b[42m\x1b[30m INFO \x1b[0m \x1b[32mtest_stringer\x1b[0m\n  \x1b[35mitem1\x1b[0m: A: test\n\n")

	if !bytes.Equal(w.WrittenData, expected) {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, w.WrittenData)
	}
}

type logStringerExample2 struct {
	Inner logStringerExample1
	Other int
}

func testStringerInner(t *testing.T, o *Options) {
	w := &MockWriter{}
	logger := slog.New(NewHandler(w, o))
	item1 := logStringerExample2{
		Inner: logStringerExample1{
			A: []byte("test"),
		},
		Other: 42,
	}
	logger.Info("test_stringer_inner",
		slog.Any("item1", item1),
	)

	expected := []byte(
		"\x1b[2m\x1b[37m[]\x1b[0m \x1b[42m\x1b[30m INFO \x1b[0m \x1b[32mtest_stringer_inner\x1b[0m\n\x1b[33mS\x1b[0m \x1b[35mitem1\x1b[0m: \x1b[33md\x1b[0m\x1b[33me\x1b[0m\x1b[33mv\x1b[0m\x1b[33ms\x1b[0m\x1b[33ml\x1b[0m\x1b[33mo\x1b[0m\x1b[33mg\x1b[0m\x1b[33m.\x1b[0m\x1b[33ml\x1b[0m\x1b[33mo\x1b[0m\x1b[33mg\x1b[0m\x1b[33mS\x1b[0m\x1b[33mt\x1b[0m\x1b[33mr\x1b[0m\x1b[33mi\x1b[0m\x1b[33mn\x1b[0m\x1b[33mg\x1b[0m\x1b[33me\x1b[0m\x1b[33mr\x1b[0m\x1b[33mE\x1b[0m\x1b[33mx\x1b[0m\x1b[33ma\x1b[0m\x1b[33mm\x1b[0m\x1b[33mp\x1b[0m\x1b[33ml\x1b[0m\x1b[33me\x1b[0m\x1b[33m2\x1b[0m\n    \x1b[32mInner\x1b[0m: A: test\n    \x1b[32mOther\x1b[0m: \x1b[33m42\x1b[0m\n\n",
	)

	if !bytes.Equal(w.WrittenData, expected) {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, w.WrittenData)
	}
}

func testNoColor(t *testing.T, o *Options) {
	w := &MockWriter{}
	o.NoColor = true
	logger := slog.New(NewHandler(w, o))

	logger.Info("msg",
		slog.Any("i", 1),
		slog.Any("f", 2.2),
		slog.Any("s", "someString"),
		slog.Any("m", map[int]string{3: "three", 4: "four"}),
	)

	expected := []byte("[]  INFO  msg\n# f: 2.2\n# i: 1\nM m: 2 map[int]string\n    3: three\n    4: four\n  s: someString\n\n")

	if !bytes.Equal(w.WrittenData, expected) {
		t.Errorf("\nExpected:\n%s\nGot:\n%s\nExpected:\n%[1]q\nGot:\n%[2]q", expected, w.WrittenData)
	}
}
