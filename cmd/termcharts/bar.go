// Last modified: 2026-01-03

package main

import (
	"bufio"
	"fmt"
	"os"
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
  termcharts bar 10 20 30 --color`,
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
}

func runBar(cmd *cobra.Command, args []string) error {
	// Parse data from various sources
	data, err := parseBarData(args)
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

// readBarDataWithLabels reads data and labels from a file.
// File format: value,label (one per line)
// Or just values (one per line)
func readBarDataWithLabels(filename string) ([]float64, []string, error) {
	file, err := os.Open(filename) // #nosec G304 - filename is provided by user via CLI
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to close file: %v\n", closeErr)
		}
	}()

	var data []float64
	var labels []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}

		// Check if line contains comma (value,label format)
		if strings.Contains(line, ",") {
			parts := strings.SplitN(line, ",", 2)
			if len(parts) == 2 {
				nums, err := parseNumbers([]string{parts[0]})
				if err != nil {
					return nil, nil, fmt.Errorf("invalid data in file %s: %s", filename, line)
				}
				if len(nums) > 0 {
					data = append(data, nums[0])
					labels = append(labels, strings.TrimSpace(parts[1]))
				}
				continue
			}
		}

		// Otherwise just parse as numbers
		nums, err := parseNumberLine(line)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid data in file %s: %s", filename, line)
		}
		data = append(data, nums...)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return data, labels, nil
}
