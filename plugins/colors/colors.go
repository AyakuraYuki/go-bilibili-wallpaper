package colors

import (
	"github.com/fatih/color"
)

func White(format string, v ...any) string {
	f := color.New(color.FgWhite).SprintfFunc()
	return f(format, v...)
}

func Red(format string, v ...any) string {
	f := color.New(color.FgRed).SprintfFunc()
	return f(format, v...)
}

func Yellow(format string, v ...any) string {
	f := color.New(color.FgYellow).SprintfFunc()
	return f(format, v...)
}

func Green(format string, v ...any) string {
	f := color.New(color.FgGreen).SprintfFunc()
	return f(format, v...)
}
