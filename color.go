package devslog

import "fmt"

type (
	foregroundColor   string
	backgroundColor   string
	commonValuesColor string
)

const (
	// Foreground colors
	fgBlack   foregroundColor = "\x1b[30m"
	fgRed     foregroundColor = "\x1b[31m"
	fgGreen   foregroundColor = "\x1b[32m"
	fgYellow  foregroundColor = "\x1b[33m"
	fgBlue    foregroundColor = "\x1b[34m"
	fgMagenta foregroundColor = "\x1b[35m"
	fgCyan    foregroundColor = "\x1b[36m"
	fgWhite   foregroundColor = "\x1b[37m"

	// Background colors
	bgBlack   backgroundColor = "\x1b[40m"
	bgRed     backgroundColor = "\x1b[41m"
	bgGreen   backgroundColor = "\x1b[42m"
	bgYellow  backgroundColor = "\x1b[43m"
	bgBlue    backgroundColor = "\x1b[44m"
	bgMagenta backgroundColor = "\x1b[45m"
	bgCyan    backgroundColor = "\x1b[46m"
	bgWhite   backgroundColor = "\x1b[47m"

	// Common consts
	resetColor     commonValuesColor = "\x1b[0m"
	faintColor     commonValuesColor = "\x1b[2m"
	underlineColor commonValuesColor = "\x1b[4m"
)

// Color string foreground
func cs(text string, fgColor foregroundColor) string {
	return fmt.Sprintf("%v%v%v", fgColor, text, resetColor)
}

// Color string fainted
func csf(text string, fgColor foregroundColor) string {
	return fmt.Sprintf("%v%v%v%v", fgColor, faintColor, text, resetColor)
}

// Color string background
func csb(text string, fgColor foregroundColor, bgColor backgroundColor) string {
	return fmt.Sprintf("%v%v%v%v", fgColor, bgColor, text, resetColor)
}

// Underline text
func ul(text string) string {
	return fmt.Sprintf("%v%v%v", underlineColor, text, resetColor)
}
