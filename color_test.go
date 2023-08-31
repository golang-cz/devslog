package devslog

import (
	"bytes"
	"testing"
)

func Test_Color(t *testing.T) {
	b := []byte("Hello")
	test_ColorCs(t, b)
	test_ColorCsf(t, b)
	test_ColorCsb(t, b)
	test_ColorUl(t, b)
}

func test_ColorCs(t *testing.T, b []byte) {
	result := cs(b, fgGreen)

	expected := []byte("\x1b[32mHello\x1b[0m")
	if !bytes.Equal(expected, result) {
		t.Errorf("\nExpected: %s\nResult:   %s\nExpected: %[1]q\nResult:   %[2]q", expected, result)
	}
}

func test_ColorCsf(t *testing.T, b []byte) {
	result := csf(b, fgBlue)

	expected := []byte("\x1b[2m\x1b[34mHello\x1b[0m")
	if !bytes.Equal(expected, result) {
		t.Errorf("\nExpected: %s\nResult:   %s\nExpected: %[1]q\nResult:   %[2]q", expected, result)
	}
}

func test_ColorCsb(t *testing.T, b []byte) {
	result := csb(b, fgYellow, bgRed)

	expected := []byte("\x1b[41m\x1b[33mHello\x1b[0m")
	if !bytes.Equal(expected, result) {
		t.Errorf("\nExpected: %s\nResult:   %s\nExpected: %[1]q\nResult:   %[2]q", expected, result)
	}
}

func test_ColorUl(t *testing.T, b []byte) {
	result := ul(b)

	expected := []byte("\x1b[4mHello\x1b[0m")
	if !bytes.Equal(expected, result) {
		t.Errorf("\nExpected: %s\nResult:   %s\nExpected: %[1]q\nResult:   %[2]q", expected, result)
	}
}
