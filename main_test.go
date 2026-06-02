package devslog

import (
	"os"
	"testing"
)

// TestMain neutralizes color-related environment variables so the suite's
// ANSI-escape assertions are hermetic regardless of the developer's shell.
// Per-test cases that exercise env-driven color disabling use t.Setenv.
func TestMain(m *testing.M) {
	os.Unsetenv("NO_COLOR")
	os.Unsetenv("TERM")
	os.Exit(m.Run())
}
