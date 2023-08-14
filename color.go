package devslog

import "fmt"

type color string

const (
	// Foreground colors
	fgBlack   color = "\x1b[30m"
	fgRed     color = "\x1b[31m"
	fgGreen   color = "\x1b[32m"
	fgYellow  color = "\x1b[33m"
	fgBlue    color = "\x1b[34m"
	fgMagenta color = "\x1b[35m"
	fgCyan    color = "\x1b[36m"
	fgWhite   color = "\x1b[37m"

	// Background colors
	bgBlack   color = "\x1b[40m"
	bgRed     color = "\x1b[41m"
	bgGreen   color = "\x1b[42m"
	bgYellow  color = "\x1b[43m"
	bgBlue    color = "\x1b[44m"
	bgMagenta color = "\x1b[45m"
	bgCyan    color = "\x1b[46m"
	bgWhite   color = "\x1b[47m"

	// Common consts
	resetColor     color = "\x1b[0m"
	faintColor     color = "\x1b[2m"
	underlineColor color = "\x1b[4m"
)

// Color string foreground
func cs(text string, fgColor color) string {
	return fmt.Sprintf("%v%v%v", fgColor, text, resetColor)
}

// Color string fainted
func csf(text string, fgColor color) string {
	return fmt.Sprintf("%v%v%v%v", fgColor, faintColor, text, resetColor)
}

// Color string background
func csb(text string, fgColor color, bgColor color) string {
	return fmt.Sprintf("%v%v%v%v", fgColor, bgColor, text, resetColor)
}

// Underline text
func ul(text string) string {
	return fmt.Sprintf("%v%v%v", underlineColor, text, resetColor)
}
