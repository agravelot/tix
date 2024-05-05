package color

import (
	"log"
	"runtime"
)

// Coloe represents a color in the terminal
type Color string

var (
	EnableColor       = true
	Reset       Color = "\033[0m"
	Red         Color = "\033[31m"
	Green       Color = "\033[32m"
	Yellow      Color = "\033[33m"
	Blue        Color = "\033[34m"
	Purple      Color = "\033[35m"
	Cyan        Color = "\033[36m"
	Gray        Color = "\033[37m"
	White       Color = "\033[97m"
	allColors         = []Color{Red, Green, Yellow, Blue, Purple, Cyan, Gray, White}
)

func init() {
	if runtime.GOOS == "windows" {
		log.Println("Disabling colors on Windows")
		EnableColor = false
	}
}

func Colorize(color Color, s string) string {
	if !EnableColor {
		return s
	}
	return string(color) + s + string(Reset)
}

func ColorByIndex(index int) Color {
	if !EnableColor {
		return ""
	}
	return allColors[index%len(allColors)]
}
