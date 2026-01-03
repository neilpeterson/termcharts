package main

import (
	"fmt"
	"strings"

	"github.com/neilpeterson/termcharts/pkg/termcharts"
	"github.com/spf13/cobra"
)

var (
	pieWidth      int
	pieColor      bool
	pieASCII      bool
	pieNoColor    bool
	pieShowValues bool
	pieTitle      string
	pieLabels     string
)

var pieCmd = &cobra.Command{
	Use:   "pie [values...]",
	Short: "Create a pie chart",
	Long: `Create a pie chart to visualize proportional data.

Pie charts display data as proportional slices, showing both a
visual representation and a legend with percentages.

Data can be provided as:
  - Command-line arguments: termcharts pie 30 25 20 15 10
  - File path: termcharts pie data.txt
  - Stdin: cat data.txt | termcharts pie

Data format:
  - One number per line, or
  - Space-separated numbers on one line, or
  - Comma-separated numbers

Labels can be provided via --labels flag as comma-separated values.

Examples:
  # Simple pie chart
  termcharts pie 30 25 20 15 10

  # With labels
  termcharts pie 30 25 20 15 10 --labels "Chrome,Firefox,Safari,Edge,Other"

  # With title and values
  termcharts pie 30 25 20 --title "Market Share" --show-values

  # From file with color
  termcharts pie data.txt --color

  # ASCII mode for compatibility
  termcharts pie 50 30 20 --ascii`,
	RunE: runPie,
}

func init() {
	rootCmd.AddCommand(pieCmd)

	pieCmd.Flags().IntVarP(&pieWidth, "width", "w", 80, "chart width in characters")
	pieCmd.Flags().BoolVarP(&pieColor, "color", "c", false, "enable colored output")
	pieCmd.Flags().BoolVar(&pieASCII, "ascii", false, "use ASCII characters only")
	pieCmd.Flags().BoolVar(&pieNoColor, "no-color", false, "disable colored output")
	pieCmd.Flags().BoolVar(&pieShowValues, "show-values", false, "display numeric values")
	pieCmd.Flags().StringVarP(&pieTitle, "title", "t", "", "chart title")
	pieCmd.Flags().StringVarP(&pieLabels, "labels", "l", "", "comma-separated labels for each slice")
}

func runPie(cmd *cobra.Command, args []string) error {
	// Parse data from various sources
	data, err := parsePieData(args)
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

	// Apply width
	if pieWidth > 0 {
		opts = append(opts, termcharts.WithWidth(pieWidth))
	}

	// Apply title if specified
	if pieTitle != "" {
		opts = append(opts, termcharts.WithTitle(pieTitle))
	}

	// Apply labels if specified
	if pieLabels != "" {
		labels := parsePieLabels(pieLabels)
		opts = append(opts, termcharts.WithLabels(labels))
	}

	// Apply show values
	if pieShowValues {
		opts = append(opts, termcharts.WithShowValues(true))
	}

	// Apply style
	if pieASCII {
		opts = append(opts, termcharts.WithStyle(termcharts.StyleASCII))
	}

	// Apply color settings
	if pieNoColor {
		colorEnabled := false
		opts = append(opts, termcharts.WithColor(colorEnabled))
	} else if pieColor {
		colorEnabled := true
		opts = append(opts, termcharts.WithColor(colorEnabled))
	}

	// Create and render pie chart
	pie := termcharts.NewPieChart(opts...)
	fmt.Print(pie.Render())

	return nil
}

// parsePieData parses data from command-line args, files, or stdin.
func parsePieData(args []string) ([]float64, error) {
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

// parsePieLabels parses comma-separated labels.
func parsePieLabels(labelsStr string) []string {
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
