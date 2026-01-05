package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/neilpeterson/termcharts/pkg/termcharts"
	"github.com/spf13/cobra"
)

var (
	barWidth      int
	barHeight     int
	barColor      bool
	barASCII      bool
	barNoColor    bool
	barVertical   bool
	barShowValues bool
	barTitle      string
	barLabels     string
	barGrouped    bool
	barStacked    bool
	barShowLegend bool
	barSeries     string
)

var barCmd = &cobra.Command{
	Use:   "bar [values...]",
	Short: "Create a bar chart",
	Long: `Create a bar chart to visualize data.

Bar charts can be rendered horizontally or vertically and support
labels, titles, and customizable styling options.

Data can be provided as:
  - Command-line arguments: termcharts bar 10 20 30 25
  - File path: termcharts bar data.txt
  - Stdin: cat data.txt | termcharts bar

Data format:
  - One number per line, or
  - Space-separated numbers on one line, or
  - Comma-separated numbers

For multi-series charts (grouped/stacked), use --series flag with JSON:
  --series '[{"label":"2023","data":[10,20,30]},{"label":"2024","data":[15,25,35]}]'

Labels can be provided via --labels flag as comma-separated values.

Examples:
  # Simple horizontal bar chart
  termcharts bar 10 20 30 40

  # With labels
  termcharts bar 10 25 15 30 --labels "Q1,Q2,Q3,Q4"

  # Vertical bar chart
  termcharts bar 10 25 15 30 --vertical

  # With title and values
  termcharts bar 10 25 15 30 --title "Sales Report" --show-values

  # From file with custom width
  termcharts bar data.txt --width 60

  # ASCII mode for compatibility
  termcharts bar 10 20 30 --ascii

  # With color
  termcharts bar 10 20 30 --color

  # Grouped bar chart (multiple series side-by-side)
  termcharts bar --series '[{"label":"2023","data":[10,20,30]},{"label":"2024","data":[15,25,35]}]' --grouped --labels "Q1,Q2,Q3"

  # Stacked bar chart (multiple series stacked)
  termcharts bar --series '[{"label":"Product A","data":[10,20,30]},{"label":"Product B","data":[5,10,15]}]' --stacked --labels "Q1,Q2,Q3"

  # Vertical grouped bar chart with legend
  termcharts bar --series '[{"label":"2023","data":[10,20,30]},{"label":"2024","data":[15,25,35]}]' --grouped --vertical --legend`,
	RunE: runBar,
}

func init() {
	rootCmd.AddCommand(barCmd)

	barCmd.Flags().IntVarP(&barWidth, "width", "w", 80, "chart width in characters")
	barCmd.Flags().IntVar(&barHeight, "height", 15, "chart height in rows (vertical mode)")
	barCmd.Flags().BoolVarP(&barColor, "color", "c", false, "enable colored output")
	barCmd.Flags().BoolVar(&barASCII, "ascii", false, "use ASCII characters only")
	barCmd.Flags().BoolVar(&barNoColor, "no-color", false, "disable colored output")
	barCmd.Flags().BoolVarP(&barVertical, "vertical", "v", false, "render vertical bar chart")
	barCmd.Flags().BoolVar(&barShowValues, "show-values", false, "display numeric values on bars")
	barCmd.Flags().StringVarP(&barTitle, "title", "t", "", "chart title")
	barCmd.Flags().StringVarP(&barLabels, "labels", "l", "", "comma-separated labels for each bar")
	barCmd.Flags().BoolVarP(&barGrouped, "grouped", "g", false, "display multiple series as grouped bars")
	barCmd.Flags().BoolVarP(&barStacked, "stacked", "s", false, "display multiple series as stacked bars")
	barCmd.Flags().BoolVar(&barShowLegend, "legend", false, "show legend for multi-series charts")
	barCmd.Flags().StringVar(&barSeries, "series", "", "JSON array of series: [{\"label\":\"name\",\"data\":[1,2,3]}]")
}

func runBar(cmd *cobra.Command, args []string) error {
	// Build options
	var opts []termcharts.Option

	// Check if multi-series data is provided
	if barSeries != "" {
		series, err := parseSeriesJSON(barSeries)
		if err != nil {
			return fmt.Errorf("failed to parse series JSON: %w", err)
		}
		if len(series) == 0 {
			return fmt.Errorf("no series data provided")
		}
		opts = append(opts, termcharts.WithSeries(series))

		// Set bar mode
		if barStacked {
			opts = append(opts, termcharts.WithBarMode(termcharts.BarModeStacked))
		} else {
			opts = append(opts, termcharts.WithBarMode(termcharts.BarModeGrouped))
		}

		// Show legend by default for multi-series, or if explicitly requested
		if barShowLegend {
			opts = append(opts, termcharts.WithShowLegend(true))
		}
	} else {
		// Parse single-series data from various sources
		data, err := parseBarData(args)
		if err != nil {
			return fmt.Errorf("failed to parse data: %w", err)
		}

		if len(data) == 0 {
			return fmt.Errorf("no data provided")
		}

		opts = append(opts, termcharts.WithData(data))
	}

	// Apply width
	if barWidth > 0 {
		opts = append(opts, termcharts.WithWidth(barWidth))
	}

	// Apply height if vertical mode
	if barVertical {
		opts = append(opts, termcharts.WithDirection(termcharts.Vertical))
		if barHeight > 0 {
			opts = append(opts, termcharts.WithHeight(barHeight))
		}
	}

	// Apply title if specified
	if barTitle != "" {
		opts = append(opts, termcharts.WithTitle(barTitle))
	}

	// Apply labels if specified
	if barLabels != "" {
		labels := parseLabels(barLabels)
		opts = append(opts, termcharts.WithLabels(labels))
	}

	// Apply show values
	if barShowValues {
		opts = append(opts, termcharts.WithShowValues(true))
	}

	// Apply style
	if barASCII {
		opts = append(opts, termcharts.WithStyle(termcharts.StyleASCII))
	}

	// Apply color settings
	if barNoColor {
		colorEnabled := false
		opts = append(opts, termcharts.WithColor(colorEnabled))
	} else if barColor {
		colorEnabled := true
		opts = append(opts, termcharts.WithColor(colorEnabled))
	}

	// Create and render bar chart
	bar := termcharts.NewBarChart(opts...)
	fmt.Print(bar.Render())

	return nil
}

// parseBarData parses data from command-line args, files, or stdin.
func parseBarData(args []string) ([]float64, error) {
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

// parseLabels parses comma-separated labels.
func parseLabels(labelsStr string) []string {
	parts := strings.Split(labelsStr, ",")
	labels := make([]string, 0, len(parts))
	for _, p := range parts {
		label := strings.TrimSpace(p)
		if label != "" {
			labels = append(labels, label)
		}
	}
	return labels
}

// seriesJSON is used for JSON parsing of series data.
type seriesJSON struct {
	Label string    `json:"label"`
	Data  []float64 `json:"data"`
	Color string    `json:"color,omitempty"`
}

// parseSeriesJSON parses JSON array of series data.
func parseSeriesJSON(jsonStr string) ([]termcharts.Series, error) {
	var seriesData []seriesJSON
	if err := json.Unmarshal([]byte(jsonStr), &seriesData); err != nil {
		return nil, err
	}

	result := make([]termcharts.Series, len(seriesData))
	for i, s := range seriesData {
		result[i] = termcharts.Series{
			Label: s.Label,
			Data:  s.Data,
			Color: s.Color,
		}
	}
	return result, nil
}
