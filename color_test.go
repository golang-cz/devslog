package devslog

import (
	"testing"
)

func Test_ColorCs(t *testing.T) {
	expected := "\x1b[32mHello\x1b[0m"
	result := cs("Hello", fgGreen)

	if result != expected {
		t.Errorf("Expected: %q, but got: %q", expected, result)
	}
}

func Test_ColorCsf(t *testing.T) {
	expected := "\x1b[34m\x1b[2mHello\x1b[0m"
	result := csf("Hello", fgBlue)

	if result != expected {
		t.Errorf("Expected: %q, but got: %q", expected, result)
	}
}

func Test_ColorCsb(t *testing.T) {
	expected := "\x1b[35m\x1b[43mHello\x1b[0m"
	result := csb("Hello", fgMagenta, bgYellow)

	if result != expected {
		t.Errorf("Expected: %q, but got: %q", expected, result)
	}
}

func Test_ColorUl(t *testing.T) {
	expected := "\x1b[4mHello\x1b[0m"
	result := ul("Hello")

	if result != expected {
		t.Errorf("Expected: %q, but got: %q", expected, result)
	}
}
