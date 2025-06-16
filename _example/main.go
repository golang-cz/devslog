package main

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-cz/devslog"
)

func main() {
	production := false
	json := false

	w := os.Stdout

	slogOpts := &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelDebug,
		ReplaceAttr: replaceAttr,
	}

	var logger *slog.Logger
	if production {
		if json {
			logger = slog.New(slog.NewJSONHandler(w, slogOpts))
		} else {
			logger = slog.New(slog.NewTextHandler(w, slogOpts))
		}
	} else {
		opts := &devslog.Options{
			HandlerOptions:    slogOpts,
			MaxSlicePrintSize: 5,
			SortKeys:          true,
			DebugColor:        devslog.Magenta,
			StringIndentation: true,
		}

		logger = slog.New(devslog.NewHandler(w, opts))
	}

	slog.SetDefault(logger)

	log1(false)
	levels(true)
	longStrings(false)
	logSmall(false)
	printInfiniteLoop(false)
	printNoColor(false)
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

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.LevelKey:
		// Rename the level key from "level" to "sev".
		// a.Key = "sev"

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

	case slog.SourceKey:
		source := a.Value.Any().(*slog.Source)
		source.File = filepath.Base(source.File)
		// return slog.Attr{}
	}

	return a
}

func levels(log bool) {
	if !log {
		return
	}

	ctx := context.Background()
	slog.Log(ctx, LevelTrace, "sunrise")
	slog.Debug("adding fuel")
	slog.Info("initiating launch")
	slog.Log(ctx, LevelNotice, "button clicked")
	slog.Warn("overheating")
	slog.Error("it is broken")
	slog.Log(ctx, LevelEmergency, "very much broken")
}

func log1(log bool) {
	if !log {
		return
	}

	structSlice := []*A{
		{Num: 2, Str: "2x str"},
		{Num: 4, Str: "4x str"},
	}

	mapString := map[string]string{
		"ba na na": "man go",
		"apple":    "pear",
	}

	sliceMap := []map[string]string{
		mapString,
		mapString,
	}

	str1 := "Hello"
	str2 := "World"
	str3 := "Goodbye"
	str4 := "Curled"

	str2p := &str2
	str4p := &str4

	mapStringPtrs := map[*string]**string{
		&str1: &str2p,
		&str3: &str4p,
	}

	type customString string
	sliceCustom := []customString{"orange", "ba na na"}

	sliceSmall := []string{"dsa", "ba na na"}
	var sliceBig []float64
	for i := 0; i < 10; i++ {
		sliceBig = append(sliceBig, math.Pow(2, float64(i)))
	}

	intMap := map[int]int{1: 2, 3: 4, 5: 6}
	intMapPtr := &intMap

	sliceSlice := []*[]float64{
		&sliceBig,
		&sliceBig,
	}

	person := Person{
		FirstName: "John",
		Age:       30,
		Contact: Contact{
			Email: "john@example.com",
			Address: &Address{
				Street: "123 Main St",
				Slice:  &sliceBig,
				Map:    mapStringPtrs,
			},
		},
	}

	mapStruct := make(map[int]*Person, 2)
	mapStruct[21] = &person

	type slicer []int
	mapSlice := make(map[int][]float64, 2)
	mapSlice[0] = sliceBig

	err := fmt.Errorf("so much broken")
	err = fmt.Errorf("load config: %w", err)
	err = fmt.Errorf("init app: %w", err)

	a := A{Str: "dsadsa", Num: 45}
	b := &a
	c := &b
	d := &c

	var in uint8 = 21
	fl := 1.21
	flp := &fl
	bool := true

	date := time.Date(2012, time.March, 28, 0, 0, 0, 0, time.UTC)
	duration := time.Since(date)

	slog.Info(
		"My INFO message",
		slog.String("string", "some string"),
		slog.String("url", "https://go.dev/"),
		slog.Any("bool", true),
		slog.Any("time", date),
		slog.Any("time2", &date),
		slog.Any("duration", duration),
		slog.Any("map", mapString),
		slog.Any("intMap", intMap),
		slog.Any("sliceSmall", sliceSmall),
		slog.Any("sliceBig", &sliceBig),
		slog.Any("sliceSlice", &sliceSlice),
		slog.Any("sliceCustom", sliceCustom),
		slog.Any("emptySlice", make([]*int, 0)),
		slog.Any("mapStringPtrs", mapStringPtrs),
		slog.Any("sliceMap", sliceMap),
		slog.Group("innerGroup",
			slog.Any("flp", &flp),
			slog.Any("mapStruct", &mapStruct),
			slog.Any("mapSlice", mapSlice),
			slog.Any("intPtr", &in),
			slog.Any("boolPtr", &bool),
		),
		slog.Any("err", err),
		slog.Any("errPtr", &err),
		slog.Any("person", &person),
		slog.Any("structs1", &structSlice),
		slog.Any("optrmap", intMapPtr),
		slog.Any("structSlice", structSlice),
		slog.Any("aPtr", &a),
		slog.Any("a", a),
		slog.Any("b", B{}),
		slog.Any("pointersToStruct", d),
	)
}

func longStrings(log bool) {
	if !log {
		return
	}

	type LongStrings struct {
		S1 string
		S2 string
	}

	slog.Info("string_group",
		slog.String("string1", "this is a long string\non multiple lines.\nit would read better if we kept indentation"),
		slog.String("string11", "this is a long string\non multiple lines.\nit would read better if we kept indentation"),
		slog.Group("group",
			slog.String("string1", "this is a long string\non multiple lines.\nit would read better if we kept indentation"),
			slog.String("string11", "this is a long string\non multiple lines.\nit would read better if we kept indentation"),
			slog.Any("longStrings", LongStrings{
				S1: "this is a long string\non multiple lines.\nit would read better if we kept indentation",
				S2: "this is a long string\non multiple lines.\nit would read better if we kept indentation",
			}),
			slog.Group("group",
				slog.Any("longStrings", LongStrings{
					S1: "this is a long string\non multiple lines.\nit would read better if we kept indentation",
					S2: "this is a long string\non multiple lines.\nit would read better if we kept indentation",
				}),
				slog.String("string1", "this is a long string\non multiple lines.\nit would read better if we kept indentation"),
				slog.String("string11", "this is a long string\non multiple lines.\nit would read better if we kept indentation"),
			),
		),
	)

	strings := []string{
		"abc\ndef",
		"dsa\nqwe",
		"abc\ndef",
		"abc\ndef",
	}

	mapString := map[int]string{
		1:    "abc\ndef",
		1234: "abc\ndef",
	}

	slog.Info("longStrings",
		slog.Any("slice", strings),
		slog.Any("map", mapString),
	)
}

func logSmall(log bool) {
	if !log {
		return
	}

	s1 := &A{
		Num: 2,
		Str: "haha",
	}

	slog.Info("in",
		slog.Any("s1", s1),
		slog.Any("s2", &B{}),
	)

	slog.Info("stats", "mem", ByteCount(9503159), "num", ByteString("hafananananana"))

	a := 32
	b := &a
	slog.Info("bl", "b", &b)
}

type ByteCount int // Note this could be a struct as well.

func (b ByteCount) String() string {
	return fmt.Sprintf("from stringer: %[1]d %[1]d %[1]d", b)
}

type A struct {
	Num int
	Str string
}

func (a A) String() string {
	return fmt.Sprintf("%d - %s", a.Num, a.Str)
}

type B struct {
	pff int
}

type Person struct {
	FirstName string
	Age       int
	Contact   Contact
}

type Contact struct {
	Email   string
	Address *Address
}

type Address struct {
	Street string
	Slice  *[]float64
	Map    map[*string]**string
}

type ByteString string // Note this could be a struct as well.

type testStruct struct {
	A string
	B int
	C *int
}

type lazyTestStruct testStruct

func (l lazyTestStruct) LogValue() slog.Value {
	return slog.StringValue("dsadas")
}

type tt struct {
	t time.Time
}

func printInfiniteLoop(log bool) {
	if !log {
		return
	}

	type Infinite struct {
		I *Infinite
	}

	v1 := Infinite{}
	v2 := Infinite{}
	v3 := Infinite{}

	v1.I = &v2
	v2.I = &v3
	v3.I = &v1

	slog.Info("infinite", slog.Any("a", v3))
}

func printNoColor(log bool) {
	if !log {
		return
	}

	opts := &devslog.Options{
		HandlerOptions: &slog.HandlerOptions{
			AddSource: true,
		},
		NoColor: true,
	}

	l := slog.New(devslog.NewHandler(os.Stdout, opts))
	l.Info("msg",
		slog.String("string", "str"),
		slog.Any("map", map[int]string{3: "three", 4: "four"}),
	)
}
