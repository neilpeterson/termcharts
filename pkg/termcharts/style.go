package termcharts

import "fmt"

// RenderStyle specifies the character set used for rendering charts.
type RenderStyle int

const (
	// StyleAuto automatically selects the best rendering style based on terminal capabilities.
	StyleAuto RenderStyle = iota
	// StyleASCII uses only basic ASCII characters for maximum compatibility.
	StyleASCII
	// StyleUnicode uses Unicode block characters for higher fidelity.
	StyleUnicode
	// StyleBraille uses Unicode Braille patterns for highest resolution (line charts).
	StyleBraille
)

// String returns the string representation of the RenderStyle.
func (s RenderStyle) String() string {
	switch s {
	case StyleAuto:
		return "auto"
	case StyleASCII:
		return "ascii"
	case StyleUnicode:
		return "unicode"
	case StyleBraille:
		return "braille"
	default:
		return "unknown"
	}
}

// Theme defines colors for chart elements.
// Colors are specified as ANSI color codes or hex values.
type Theme struct {
	// Primary is the primary chart color.
	Primary string
	// Secondary is the secondary chart color (for multi-series).
	Secondary string
	// Accent is used for highlights and emphasis.
	Accent string
	// Muted is used for axes, labels, and grid lines.
	Muted string
	// Background is the background color (rarely used in terminal).
	Background string
	// Text is the default text color.
	Text string
	// Series contains colors for multiple data series.
	Series []string
}

// Predefined themes
var (
	// DefaultTheme uses standard terminal colors.
	DefaultTheme = &Theme{
		Primary:    "blue",
		Secondary:  "green",
		Accent:     "yellow",
		Muted:      "gray",
		Background: "",
		Text:       "",
		Series:     []string{"red", "blue", "yellow", "magenta", "green", "cyan"},
	}

	// DarkTheme is optimized for dark terminal backgrounds.
	DarkTheme = &Theme{
		Primary:    "cyan",
		Secondary:  "magenta",
		Accent:     "yellow",
		Muted:      "gray",
		Background: "",
		Text:       "white",
		Series:     []string{"cyan", "magenta", "yellow", "green", "blue", "red"},
	}

	// LightTheme is optimized for light terminal backgrounds.
	LightTheme = &Theme{
		Primary:    "blue",
		Secondary:  "red",
		Accent:     "orange",
		Muted:      "gray",
		Background: "",
		Text:       "black",
		Series:     []string{"blue", "red", "green", "purple", "orange", "brown"},
	}

	// MonochromeTheme uses only grayscale colors.
	MonochromeTheme = &Theme{
		Primary:    "white",
		Secondary:  "gray",
		Accent:     "white",
		Muted:      "gray",
		Background: "",
		Text:       "white",
		Series:     []string{"white", "gray"},
	}
)

// ANSI color codes for terminal output.
const (
	colorReset   = "\033[0m"
	colorBlack   = "\033[30m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorWhite   = "\033[37m"
	colorGray    = "\033[90m"
)

// colorMap maps color names to ANSI codes.
var colorMap = map[string]string{
	"black":   colorBlack,
	"red":     colorRed,
	"green":   colorGreen,
	"yellow":  colorYellow,
	"orange":  colorYellow, // Alias for yellow
	"blue":    colorBlue,
	"magenta": colorMagenta,
	"purple":  colorMagenta, // Alias for magenta
	"cyan":    colorCyan,
	"white":   colorWhite,
	"gray":    colorGray,
	"grey":    colorGray, // Alternative spelling
	"brown":   colorRed,  // Alias for red
}

// Colorize wraps text with ANSI color codes.
// If colorEnabled is false, returns the text unchanged.
func Colorize(text, color string, colorEnabled bool) string {
	if !colorEnabled || color == "" {
		return text
	}

	code, ok := colorMap[color]
	if !ok {
		return text
	}

	return fmt.Sprintf("%s%s%s", code, text, colorReset)
}

// GetSeriesColor returns the color for a data series at the given index.
// It cycles through the theme's series colors if the index exceeds the array length.
func (t *Theme) GetSeriesColor(index int) string {
	if len(t.Series) == 0 {
		return t.Primary
	}
	return t.Series[index%len(t.Series)]
}
