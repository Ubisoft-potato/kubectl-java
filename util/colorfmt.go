package util

import "github.com/fatih/color"

// console color string format func
var (
	Cyan   = color.New(color.FgCyan).SprintfFunc()
	Yellow = color.New(color.FgYellow).SprintFunc()
	Red    = color.New(color.FgRed).SprintFunc()
)
