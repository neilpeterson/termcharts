// Last modified: 2026-01-02

// Package internal provides internal utilities for terminal interaction and data manipulation.
package internal

import (
	"os"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/term"
)

// TerminalSize represents the dimensions of the terminal.
type TerminalSize struct {
	Width  int
	Height int
}

// DefaultSize provides fallback terminal dimensions.
var DefaultSize = TerminalSize{
	Width:  80,
	Height: 24,
}

// GetTerminalSize returns the current terminal dimensions.
// If detection fails, returns DefaultSize.
func GetTerminalSize() TerminalSize {
	// Try to get terminal size from file descriptor
	fd := int(os.Stdout.Fd())
	width, height, err := term.GetSize(fd)
	if err != nil || width == 0 || height == 0 {
		// Fall back to environment variables
		if w, h := getSizeFromEnv(); w > 0 && h > 0 {
			return TerminalSize{Width: w, Height: h}
		}
		// Use default size
		return DefaultSize
	}

	return TerminalSize{Width: width, Height: height}
}

// getSizeFromEnv attempts to read terminal size from environment variables.
func getSizeFromEnv() (width, height int) {
	if cols := os.Getenv("COLUMNS"); cols != "" {
		if w, err := strconv.Atoi(cols); err == nil {
			width = w
		}
	}
	if lines := os.Getenv("LINES"); lines != "" {
		if h, err := strconv.Atoi(lines); err == nil {
			height = h
		}
	}
	return width, height
}

// SupportsColor detects whether the terminal supports ANSI colors.
// Checks environment variables and terminal capabilities.
func SupportsColor() bool {
	// Check if colors are explicitly disabled
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	// Check if colors are explicitly enabled
	if os.Getenv("FORCE_COLOR") != "" {
		return true
	}

	// Check TERM environment variable
	termType := os.Getenv("TERM")
	if termType == "" || termType == "dumb" {
		return false
	}

	// Common color-supporting terminal types
	colorTerms := []string{"color", "ansi", "xterm", "screen", "tmux", "rxvt"}
	for _, ct := range colorTerms {
		if strings.Contains(termType, ct) {
			return true
		}
	}

	// Windows Terminal and ConEmu support colors
	if runtime.GOOS == "windows" {
		if os.Getenv("WT_SESSION") != "" || os.Getenv("ConEmuANSI") == "ON" {
			return true
		}
	}

	// Check if stdout is a terminal
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// SupportsUnicode detects whether the terminal supports Unicode characters.
// Checks locale and environment variables.
func SupportsUnicode() bool {
	// Check if ASCII is forced
	if os.Getenv("LANG") == "C" || os.Getenv("LC_ALL") == "C" {
		return false
	}

	// Check for UTF-8 locale
	locale := os.Getenv("LANG")
	if locale == "" {
		locale = os.Getenv("LC_ALL")
	}
	if locale == "" {
		locale = os.Getenv("LC_CTYPE")
	}

	// UTF-8 support is common in modern terminals
	if strings.Contains(strings.ToUpper(locale), "UTF-8") || strings.Contains(strings.ToUpper(locale), "UTF8") {
		return true
	}

	// On macOS and Linux, assume UTF-8 support by default
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		return true
	}

	// Windows 10+ supports Unicode
	if runtime.GOOS == "windows" {
		return true
	}

	return false
}

// IsTTY returns true if stdout is connected to a terminal.
func IsTTY() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}
