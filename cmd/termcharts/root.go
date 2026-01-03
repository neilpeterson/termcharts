package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "termcharts",
	Short: "ASCII/Unicode terminal charting tool",
	Long: `termcharts is a CLI tool for creating beautiful charts in your terminal.

It supports multiple chart types including sparklines, bar charts, line charts,
and more. Charts can be rendered using pure ASCII for maximum compatibility or
Unicode block/Braille characters for higher fidelity.

Examples:
  # Create a sparkline from values
  termcharts spark 10 20 30 25 15

  # Read from a file
  termcharts spark data.txt

  # Read from stdin
  cat data.txt | termcharts spark

  # Create a bar chart
  termcharts bar 10 20 30 25 --labels "Q1,Q2,Q3,Q4"`,
	Version: "0.1.0",
}

func init() {
	// Global flags can be added here
	// rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
}
