package ui

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

// Colors for different types of output
var (
	SuccessColor = color.New(color.FgGreen, color.Bold)
	InfoColor    = color.New(color.FgBlue)
	WarningColor = color.New(color.FgYellow)
	ErrorColor   = color.New(color.FgRed, color.Bold)
)

// PrintSuccess prints a success message
func PrintSuccess(format string, args ...interface{}) {
	SuccessColor.Printf("✓ "+format+"\n", args...)
}

// PrintInfo prints an info message
func PrintInfo(format string, args ...interface{}) {
	InfoColor.Printf("ℹ "+format+"\n", args...)
}

// PrintWarning prints a warning message
func PrintWarning(format string, args ...interface{}) {
	WarningColor.Printf("⚠ "+format+"\n", args...)
}

// PrintError prints an error message
func PrintError(format string, args ...interface{}) {
	ErrorColor.Printf("✗ "+format+"\n", args...)
}

// CreateProgressBar creates a progress bar for operations
func CreateProgressBar(max int, description string) *progressbar.ProgressBar {
	return progressbar.NewOptions(max,
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(50),
		progressbar.OptionThrottle(65),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
	)
}
