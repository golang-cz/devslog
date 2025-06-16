package devslog

import (
	"bytes"
	"testing"
)

func Test_Color(t *testing.T) {
	h := NewHandler(nil, nil)
	testGetColor(t, h)

	b := []byte("Hello")
	testColorColorString(t, b, h)
	testColorColorStringFainted(t, b, h)
	testColorColorStringBackground(t, b, h)
	testColorUnderlineText(t, b, h)
	testColorFaintedText(t, b, h)
}

func testGetColor(t *testing.T, h *developHandler) {
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

func testColorColorString(t *testing.T, b []byte, h *developHandler) {
	result := h.colorString(b, fgGreen)

	expected := []byte("\x1b[32mHello\x1b[0m")
	if !bytes.Equal(expected, result) {
		t.Errorf("\nExpected: %s\nResult:   %s\nExpected: %[1]q\nResult:   %[2]q", expected, result)
	}
}

func testColorColorStringFainted(t *testing.T, b []byte, h *developHandler) {
	result := h.colorStringFainted(b, fgBlue)

	expected := []byte("\x1b[2m\x1b[34mHello\x1b[0m")
	if !bytes.Equal(expected, result) {
		t.Errorf("\nExpected: %s\nResult:   %s\nExpected: %[1]q\nResult:   %[2]q", expected, result)
	}
}

func testColorColorStringBackground(t *testing.T, b []byte, h *developHandler) {
	result := h.colorStringBackgorund(b, fgYellow, bgRed)

	expected := []byte("\x1b[41m\x1b[33mHello\x1b[0m")
	if !bytes.Equal(expected, result) {
		t.Errorf("\nExpected: %s\nResult:   %s\nExpected: %[1]q\nResult:   %[2]q", expected, result)
	}
}

func testColorUnderlineText(t *testing.T, b []byte, h *developHandler) {
	result := h.underlineText(b)

	expected := []byte("\x1b[4mHello\x1b[0m")
	if !bytes.Equal(expected, result) {
		t.Errorf("\nExpected: %s\nResult:   %s\nExpected: %[1]q\nResult:   %[2]q", expected, result)
	}
}

func testColorFaintedText(t *testing.T, b []byte, h *developHandler) {
	result := h.faintedText(b)

	expected := []byte("\x1b[2mHello\x1b[0m")
	if !bytes.Equal(expected, result) {
		t.Errorf("\nExpected: %s\nResult:   %s\nExpected: %[1]q\nResult:   %[2]q", expected, result)
	}
}
