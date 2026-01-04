package main

import (
	"fmt"

	"github.com/neilpeterson/termcharts/pkg/termcharts"
	"github.com/spf13/cobra"
)

var (
	lineWidth     int
	lineHeight    int
	lineColor     bool
	lineASCII     bool
	lineBraille   bool
	lineNoColor   bool
	lineShowAxes  bool
	lineTitle     string
	lineLabels    string
	lineThemeName string
)

var lineCmd = &cobra.Command{
	Use:   "line [values...]",
	Short: "Create a line chart",
	Long: `Create a line chart to visualize data trends.

Line charts can be rendered using ASCII box-drawing characters,
Unicode characters, or high-resolution Braille patterns for
maximum detail.

Data can be provided as:
  - Command-line arguments: termcharts line 10 20 30 25
  - File path: termcharts line data.txt
  - Stdin: cat data.txt | termcharts line

Data format:
  - One number per line, or
  - Space-separated numbers on one line, or
  - Comma-separated numbers

Examples:
  # Simple line chart
  termcharts line 1 5 2 8 3 7 4 6

  # High-resolution Braille rendering
  termcharts line 1 5 2 8 3 7 4 6 --braille

  # With title and axes
  termcharts line 10 25 15 30 20 --title "Sales Trend" --axes

  # With X-axis labels
  termcharts line 10 25 15 30 --labels "Jan,Feb,Mar,Apr"

  # From file with custom dimensions
  termcharts line data.txt --width 80 --height 15

  # ASCII mode for compatibility
  termcharts line 10 20 30 --ascii

  # With color
  termcharts line 10 20 30 --color`,
	RunE: runLine,
}

func init() {
	rootCmd.AddCommand(lineCmd)

	lineCmd.Flags().IntVarP(&lineWidth, "width", "w", 60, "chart width in characters")
	lineCmd.Flags().IntVar(&lineHeight, "height", 12, "chart height in rows")
	lineCmd.Flags().BoolVarP(&lineColor, "color", "c", false, "enable colored output")
	lineCmd.Flags().BoolVar(&lineASCII, "ascii", false, "use ASCII characters only")
	lineCmd.Flags().BoolVarP(&lineBraille, "braille", "b", false, "use high-resolution Braille patterns")
	lineCmd.Flags().BoolVar(&lineNoColor, "no-color", false, "disable colored output")
	lineCmd.Flags().BoolVar(&lineShowAxes, "axes", true, "show axes and labels")
	lineCmd.Flags().StringVarP(&lineTitle, "title", "t", "", "chart title")
	lineCmd.Flags().StringVarP(&lineLabels, "labels", "l", "", "comma-separated X-axis labels")
	lineCmd.Flags().StringVar(&lineThemeName, "theme", "default", "color theme (default, dark, light, mono)")
}

func runLine(cmd *cobra.Command, args []string) error {
	// Parse data from various sources
	data, err := parseLineData(args)
	if err != nil {
		return fmt.Errorf("failed to parse data: %w", err)
	}

	if len(data) == 0 {
		return fmt.Errorf("no data provided")
	}

	// Build options
	opts := []termcharts.Option{
		termcharts.WithData(data),
	}

	// Apply dimensions
	if lineWidth > 0 {
		opts = append(opts, termcharts.WithWidth(lineWidth))
	}
	if lineHeight > 0 {
		opts = append(opts, termcharts.WithHeight(lineHeight))
	}

	// Apply title if specified
	if lineTitle != "" {
		opts = append(opts, termcharts.WithTitle(lineTitle))
	}

	// Apply labels if specified
	if lineLabels != "" {
		labels := parseLabels(lineLabels)
		opts = append(opts, termcharts.WithLabels(labels))
	}

	// Apply axes setting
	opts = append(opts, termcharts.WithShowAxes(lineShowAxes))

	// Apply style
	if lineBraille {
		opts = append(opts, termcharts.WithStyle(termcharts.StyleBraille))
	} else if lineASCII {
		opts = append(opts, termcharts.WithStyle(termcharts.StyleASCII))
	}

	// Apply color settings
	if lineNoColor {
		colorEnabled := false
		opts = append(opts, termcharts.WithColor(colorEnabled))
	} else if lineColor {
		colorEnabled := true
		opts = append(opts, termcharts.WithColor(colorEnabled))
	}

	// Apply theme
	theme := getTheme(lineThemeName)
	if theme != nil {
		opts = append(opts, termcharts.WithTheme(theme))
	}

	// Create and render line chart
	line := termcharts.NewLineChart(opts...)
	fmt.Print(line.Render())

	return nil
}

// parseLineData parses data from command-line args, files, or stdin.
func parseLineData(args []string) ([]float64, error) {
	// If no args, read from stdin
	if len(args) == 0 {
		return readDataFromStdin()
	}

	// If single arg and it's a file, read from file
	if len(args) == 1 {
		if fileExists(args[0]) {
			return readDataFromFile(args[0])
		}
	}

	// Otherwise, parse args as numbers
	return parseNumbers(args)
}

// getTheme returns a theme by name.
func getTheme(name string) *termcharts.Theme {
	switch name {
	case "dark":
		return termcharts.DarkTheme
	case "light":
		return termcharts.LightTheme
	case "mono", "monochrome":
		return termcharts.MonochromeTheme
	default:
		return termcharts.DefaultTheme
	}
}
