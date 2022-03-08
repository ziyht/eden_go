package elog

import "github.com/fatih/color"

type colorSwitch int

const (
  ColorOff colorSwitch = 0 + iota
	ColorOn
	ColorAuto

)

const (
	COLOR_BLACK = 30 + iota
	COLOR_RED
	COLOR_GREEN
	COLOR_YELLOW
	COLOR_BLUE
	COLOR_MAGENTA
	COLOR_CYAN
	COLOR_WHITE
)

// Foreground Hi-Intensity text colors
const (
	COLOR_HI_BLACK = 90 + iota
	COLOR_HI_RED
	COLOR_HI_GREEN
	COLOR_HI_YELLOW
	COLOR_HI_BLUE
	COLOR_HI_MAGENTA
	COLOR_HI_CYAN
	COLOR_HI_WHITE
)

// defaultLevelColor defines the default level and its mapping prefix string.
var defaultLevelColor = map[Level]int{
	LEVEL_DEBUG : COLOR_YELLOW,
	LEVEL_INFO  : COLOR_GREEN,
	LEVEL_WARN  : COLOR_MAGENTA,
	LEVEL_ERROR : COLOR_RED,
	LEVEL_DPANIC: COLOR_HI_RED,
	LEVEL_PANIC : COLOR_HI_RED,
	LEVEL_FATAL : COLOR_HI_RED,
}

// getColoredStr returns a string that is colored by given color.
func getColoredStr(c int, s string) string {
	c_ := color.New(color.Attribute(c))
	c_.EnableColor()
	return c_.Sprint(s)
}
func getAutoColoredStr(c int, s string) string {
	return color.New(color.Attribute(c)).Sprint(s)
}

func getColorByLevel(level Level) int {
	return defaultLevelColor[level]
}
