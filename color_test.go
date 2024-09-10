package devslog

import (
	"bytes"
	"testing"
)

func Test_Color(t *testing.T) {
	h := NewHandler(nil, nil)
	test_GetColor(t, h)

	b := []byte("Hello")
	test_ColorCs(t, b, h)
	test_ColorCsf(t, b, h)
	test_ColorCsb(t, b, h)
	test_ColorUl(t, b, h)
}

func test_GetColor(t *testing.T, h *developHandler) {
	result := h.getColor(Black)
	expected := colors[1].fg

	if !bytes.Equal(expected, result.fg) {
		t.Errorf("\nExpected: %s\nResult:   %s\nExpected: %[1]q\nResult:   %[2]q", expected, result)
	}

	result = h.getColor(Color(20))
	expected = colors[8].fg

	if !bytes.Equal(expected, result.fg) {
		t.Errorf("\nExpected: %s\nResult:   %s\nExpected: %[1]q\nResult:   %[2]q", expected, result)
	}
}

func test_ColorCs(t *testing.T, b []byte, h *developHandler) {
	result := h.cs(b, fgGreen)

	expected := []byte("\x1b[32mHello\x1b[0m")
	if !bytes.Equal(expected, result) {
		t.Errorf("\nExpected: %s\nResult:   %s\nExpected: %[1]q\nResult:   %[2]q", expected, result)
	}
}

func test_ColorCsf(t *testing.T, b []byte, h *developHandler) {
	result := h.csf(b, fgBlue)

	expected := []byte("\x1b[2m\x1b[34mHello\x1b[0m")
	if !bytes.Equal(expected, result) {
		t.Errorf("\nExpected: %s\nResult:   %s\nExpected: %[1]q\nResult:   %[2]q", expected, result)
	}
}

func test_ColorCsb(t *testing.T, b []byte, h *developHandler) {
	result := h.csb(b, fgYellow, bgRed)

	expected := []byte("\x1b[41m\x1b[33mHello\x1b[0m")
	if !bytes.Equal(expected, result) {
		t.Errorf("\nExpected: %s\nResult:   %s\nExpected: %[1]q\nResult:   %[2]q", expected, result)
	}
}

func test_ColorUl(t *testing.T, b []byte, h *developHandler) {
	result := h.ul(b)

	expected := []byte("\x1b[4mHello\x1b[0m")
	if !bytes.Equal(expected, result) {
		t.Errorf("\nExpected: %s\nResult:   %s\nExpected: %[1]q\nResult:   %[2]q", expected, result)
	}
}
