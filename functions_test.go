package devslog

import (
	"fmt"
	"testing"
)

func Test_FunctionsDevslogSlice(t *testing.T) {
	initialized = true
	elementDivider = string(rune(29))
	attr := Slice("testSlice", []int{1, 2})
	expect := fmt.Sprintf("testSlice=slice[1%c2]", rune(29))

	if attr.String() != expect {
		t.Errorf("Expected: %v, but got: %v", expect, attr.String())
	}
}

func Test_FunctionsDevslogMap(t *testing.T) {
	initialized = true
	elementDivider = string(rune(29))
	attr := Map("testMap", map[string]int{"a": 1, "b": 2})
	expect1 := fmt.Sprintf("testMap=map[a:1%cb:2]", rune(29))
	expect2 := fmt.Sprintf("testMap=map[b:2%ca:1]", rune(29))

	if attr.String() != expect1 && attr.String() != expect2 {
		t.Errorf("Expected: %v or %v, but got: %v", expect1, expect2, attr.String())
	}
}

func Test_FunctionsProductionSlice(t *testing.T) {
	initialized = false
	elementDivider = string(rune(29))
	attr := Slice("testSlice", []int{1, 2})
	expect := "testSlice=[1 2]"

	if attr.String() != expect {
		t.Errorf("Expected: %v, but got: %v", expect, attr.String())
	}
}

func Test_FunctionsProductionMap(t *testing.T) {
	initialized = false
	elementDivider = string(rune(29))
	attr := Map("testMap", map[string]int{"a": 1, "b": 2})
	expect1 := "testMap=map[a:1 b:2]"
	expect2 := "testMap=map[b:2 a:1]"

	if attr.String() != expect1 && attr.String() != expect2 {
		t.Errorf("Expected: %v or %v, but got: %v", expect1, expect2, attr.String())
	}
}
