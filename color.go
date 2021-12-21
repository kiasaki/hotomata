package hotomata

import (
	"fmt"
)

// from https://github.com/mitchellh/cli/blob/master/ui_colored.go
// Color is a posix shell color code to use.
type Color struct {
	Code int
	Bold bool
}

// A list of colors that are useful. These are all non-bolded by default.
var (
	ColorNone    Color = Color{-1, false}
	ColorRed           = Color{31, false}
	ColorGreen         = Color{32, false}
	ColorYellow        = Color{33, false}
	ColorBlue          = Color{34, false}
	ColorMagenta       = Color{35, false}
	ColorCyan          = Color{36, false}
)

func Colorize(message string, color Color) string {
	if color.Code == -1 {
		return message
	}

	attr := 0
	if color.Bold {
		attr = 1
	}

	return fmt.Sprintf("\033[%d;%dm%s\033[0m", attr, color.Code, message)
}
