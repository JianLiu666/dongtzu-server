// Package color adds coloring functionality for TTY output.
package color

import (
	"fmt"
	"strconv"
	"strings"
)

const template = "\x1b[%sm%s\x1b[0m"

// Foreground colors.
const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

const (
	Blue1 Color = 94
	Blue2 Color = 104
	Blue3 Color = 44
)

var (
	Blue4 Style = New(40, 34)
)

// Color represents a text color.
type Color uint8
type Style []Color

func New(colors ...Color) Style {
	return colors
}

// Add adds the coloring to the given string.
func (c Color) Add(s string) string {
	return fmt.Sprintf(template, strconv.Itoa(int(c)), s)
}

func (s Style) Add(msg string) string {
	return fmt.Sprintf(template, s.colors2code(), msg)
}

// convert colors to code. return like "32;45;3"
func (s Style) colors2code() string {
	if len(s) == 0 {
		return ""
	}

	var codes []string
	for _, c := range s {
		codes = append(codes, strconv.Itoa(int(c)))
	}

	return strings.Join(codes, ";")
}
