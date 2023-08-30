package devslog

type (
	foregroundColor   []byte
	backgroundColor   []byte
	commonValuesColor []byte
)

var (
	// Foreground colors
	fgBlack   foregroundColor = []byte("\x1b[30m")
	fgRed     foregroundColor = []byte("\x1b[31m")
	fgGreen   foregroundColor = []byte("\x1b[32m")
	fgYellow  foregroundColor = []byte("\x1b[33m")
	fgBlue    foregroundColor = []byte("\x1b[34m")
	fgMagenta foregroundColor = []byte("\x1b[35m")
	fgCyan    foregroundColor = []byte("\x1b[36m")
	fgWhite   foregroundColor = []byte("\x1b[37m")

	// Background colors
	bgBlack   backgroundColor = []byte("\x1b[40m")
	bgRed     backgroundColor = []byte("\x1b[41m")
	bgGreen   backgroundColor = []byte("\x1b[42m")
	bgYellow  backgroundColor = []byte("\x1b[43m")
	bgBlue    backgroundColor = []byte("\x1b[44m")
	bgMagenta backgroundColor = []byte("\x1b[45m")
	bgCyan    backgroundColor = []byte("\x1b[46m")
	bgWhite   backgroundColor = []byte("\x1b[47m")

	// Common consts
	resetColor     commonValuesColor = []byte("\x1b[0m")
	faintColor     commonValuesColor = []byte("\x1b[2m")
	underlineColor commonValuesColor = []byte("\x1b[4m")
)

// Color string foreground
func cs(b []byte, fgColor foregroundColor) []byte {
	b = append(fgColor, b...)
	b = append(b, resetColor...)
	return b
}

// Color string fainted
func csf(b []byte, fgColor foregroundColor) []byte {
	b = append(fgColor, b...)
	b = append(faintColor, b...)
	b = append(b, resetColor...)
	return b
}

// Color string background
func csb(b []byte, fgColor foregroundColor, bgColor backgroundColor) []byte {
	b = append(fgColor, b...)
	b = append(bgColor, b...)
	b = append(b, resetColor...)
	return b
}

// Underline text
func ul(b []byte) []byte {
	b = append(underlineColor, b...)
	b = append(b, resetColor...)
	return b
}
