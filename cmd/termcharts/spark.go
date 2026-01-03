// Last modified: 2026-01-02

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/neilpeterson/termcharts/pkg/termcharts"
	"github.com/spf13/cobra"
)

var (
	sparkWidth   int
	sparkColor   bool
	sparkASCII   bool
	sparkNoColor bool
)

var sparkCmd = &cobra.Command{
	Use:   "spark [values...]",
	Short: "Create a sparkline chart",
	Long: `Create a compact sparkline chart to visualize data trends.

Sparklines use Unicode block characters (▁▂▃▄▅▆▇█) to visualize
data in a single line, perfect for dashboards and monitoring.

Data can be provided as:
  - Command-line arguments: termcharts spark 10 20 30 25 15
  - File path: termcharts spark data.txt
  - Stdin: cat data.txt | termcharts spark

Data format:
  - One number per line, or
  - Space-separated numbers on one line, or
  - Comma-separated numbers

Examples:
  # From arguments
  termcharts spark 10 20 30 25 15 35 40

  # From file
  termcharts spark metrics.txt

  # From stdin
  echo "1 5 2 8 3 7" | termcharts spark

  # With width limit
  termcharts spark data.txt --width 50

  # ASCII mode for compatibility
  termcharts spark 10 20 30 --ascii

  # With color
  termcharts spark 10 20 30 --color`,
	RunE: runSparkline,
}

func init() {
	rootCmd.AddCommand(sparkCmd)

	sparkCmd.Flags().IntVarP(&sparkWidth, "width", "w", 0, "maximum width in characters (0 = no limit)")
	sparkCmd.Flags().BoolVarP(&sparkColor, "color", "c", false, "enable colored output")
	sparkCmd.Flags().BoolVar(&sparkASCII, "ascii", false, "use ASCII characters only")
	sparkCmd.Flags().BoolVar(&sparkNoColor, "no-color", false, "disable colored output")
}

func runSparkline(cmd *cobra.Command, args []string) error {
	// Parse data from various sources
	data, err := parseSparklineData(args)
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

	// Apply width if specified
	if sparkWidth > 0 {
		opts = append(opts, termcharts.WithWidth(sparkWidth))
	}

	// Apply style
	if sparkASCII {
		opts = append(opts, termcharts.WithStyle(termcharts.StyleASCII))
	}

	// Apply color settings
	if sparkNoColor {
		colorEnabled := false
		opts = append(opts, termcharts.WithColor(colorEnabled))
	} else if sparkColor {
		colorEnabled := true
		opts = append(opts, termcharts.WithColor(colorEnabled))
	}

	// Create and render sparkline
	spark := termcharts.NewSparkline(opts...)
	fmt.Println(spark.Render())

	return nil
}

// parseSparklineData parses data from command-line args, files, or stdin.
func parseSparklineData(args []string) ([]float64, error) {
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

// readDataFromStdin reads numeric data from stdin.
func readDataFromStdin() ([]float64, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	// Check if stdin has data
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return nil, fmt.Errorf("no data provided via stdin")
	}

	var data []float64
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Try to parse as space-separated or comma-separated numbers
		nums, err := parseNumberLine(line)
		if err != nil {
			return nil, fmt.Errorf("invalid data on line: %s", line)
		}
		data = append(data, nums...)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

// readDataFromFile reads numeric data from a file.
func readDataFromFile(filename string) ([]float64, error) {
	file, err := os.Open(filename) // #nosec G304 - filename is provided by user via CLI
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// Log error but don't override return value since we're already returning
			fmt.Fprintf(os.Stderr, "Warning: failed to close file: %v\n", closeErr)
		}
	}()

	var data []float64
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}

		// Try to parse as space-separated or comma-separated numbers
		nums, err := parseNumberLine(line)
		if err != nil {
			return nil, fmt.Errorf("invalid data in file %s: %s", filename, line)
		}
		data = append(data, nums...)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

// parseNumberLine parses a line containing space-separated or comma-separated numbers.
func parseNumberLine(line string) ([]float64, error) {
	// Try comma-separated first
	if strings.Contains(line, ",") {
		parts := strings.Split(line, ",")
		return parseNumbers(parts)
	}

	// Otherwise space-separated
	parts := strings.Fields(line)
	return parseNumbers(parts)
}

// parseNumbers converts string values to float64.
func parseNumbers(values []string) ([]float64, error) {
	var nums []float64
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		num, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number: %s", v)
		}
		nums = append(nums, num)
	}
	return nums, nil
}

// fileExists checks if a file exists.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
