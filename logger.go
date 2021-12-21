package hotomata

import (
	"fmt"
	"io"
)

const logFillerRune = '-'
const loglineLength = 100

type Logger struct {
	Verbose    bool
	writer     io.Writer
	colored    bool
	fillerRune rune
}

func NewLogger(writer io.Writer, colored bool, verbose bool) *Logger {
	return &Logger{writer: writer, colored: colored, fillerRune: logFillerRune, Verbose: verbose}
}

func (l *Logger) Write(c Color, value string) {
	if l.colored {
		l.writer.Write([]byte(Colorize(value, c)))
	} else {
		l.writer.Write([]byte(value))
	}
}

// Write no color
func (l *Logger) Writenc(value string) {
	l.writer.Write([]byte(value))
}

func (l *Logger) WriteLine(c Color, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	for len(msg) < loglineLength {
		msg = msg + string(l.fillerRune)
	}
	l.Write(c, msg+"\n")
}
