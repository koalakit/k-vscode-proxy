package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

var debugMode = false

func LogMode(debug bool) {
	debugMode = debug
}

func LogDebug(a ...any) {
	if !debugMode {
		return
	}

	var now = time.Now().Format("2006/01/02 - 15:04:05")
	var header = fmt.Sprintf("[%s] %s", color.HiYellowString("DEBUG"), now)
	fmt.Println(header, fmt.Sprint(a...))
}

func LogDebugf(format string, a ...any) {
	if !debugMode {
		return
	}

	var now = time.Now().Format("2006/01/02 - 15:04:05")
	var header = fmt.Sprintf("[%s] %s", color.HiYellowString("DEBUG"), now)
	fmt.Println(header, fmt.Sprintf(format, a...))
}

func Log(a ...any) {
	var now = time.Now().Format("2006/01/02 - 15:04:05")
	var header = fmt.Sprintf("[%s] %s", color.HiBlueString("HTTP"), now)
	fmt.Println(header, fmt.Sprint(a...))
}

func Logf(format string, a ...any) {
	var now = time.Now().Format("2006/01/02 - 15:04:05")
	var header = fmt.Sprintf("[%s] %s", color.HiBlueString("HTTP"), now)
	fmt.Println(header, fmt.Sprintf(format, a...))
}
