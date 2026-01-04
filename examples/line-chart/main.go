// Example program demonstrating line chart features.
package main

import (
	"fmt"

	"github.com/neilpeterson/termcharts/pkg/termcharts"
)

func main() {
	fmt.Println("=== Line Chart Examples ===")
	fmt.Println()

	// Basic line chart
	fmt.Println("1. Basic Line Chart (ASCII mode)")
	fmt.Println("--------------------------------")
	data := []float64{1, 5, 2, 8, 3, 7, 4, 6}
	line := termcharts.NewLineChart(
		termcharts.WithData(data),
		termcharts.WithStyle(termcharts.StyleASCII),
		termcharts.WithWidth(50),
		termcharts.WithHeight(8),
	)
	fmt.Println(line.Render())

	// Unicode line chart with title and color
	fmt.Println("2. Unicode Line Chart with Title and Color")
	fmt.Println("------------------------------------------")
	line2 := termcharts.NewLineChart(
		termcharts.WithData(data),
		termcharts.WithStyle(termcharts.StyleUnicode),
		termcharts.WithTitle("Sales Trend"),
		termcharts.WithColor(true),
		termcharts.WithWidth(50),
		termcharts.WithHeight(8),
	)
	fmt.Println(line2.Render())

	// Braille high-resolution line chart
	fmt.Println("3. Braille High-Resolution Line Chart")
	fmt.Println("-------------------------------------")
	line3 := termcharts.NewLineChart(
		termcharts.WithData(data),
		termcharts.WithStyle(termcharts.StyleBraille),
		termcharts.WithColor(true),
		termcharts.WithWidth(50),
		termcharts.WithHeight(8),
	)
	fmt.Println(line3.Render())

	// Line chart with axes and labels
	fmt.Println("4. Line Chart with Axes and Labels")
	fmt.Println("----------------------------------")
	quarterData := []float64{150, 230, 180, 290}
	labels := []string{"Q1", "Q2", "Q3", "Q4"}
	line4 := termcharts.NewLineChart(
		termcharts.WithData(quarterData),
		termcharts.WithLabels(labels),
		termcharts.WithTitle("Quarterly Revenue ($K)"),
		termcharts.WithShowAxes(true),
		termcharts.WithColor(true),
		termcharts.WithWidth(50),
		termcharts.WithHeight(10),
	)
	fmt.Println(line4.Render())

	// Multi-series line chart
	fmt.Println("5. Multi-Series Line Chart")
	fmt.Println("--------------------------")
	series := []termcharts.Series{
		{Label: "Revenue", Data: []float64{100, 150, 130, 180, 200}},
		{Label: "Costs", Data: []float64{80, 90, 100, 110, 120}},
		{Label: "Profit", Data: []float64{20, 60, 30, 70, 80}},
	}
	line5 := termcharts.NewLineChart(
		termcharts.WithSeries(series),
		termcharts.WithTitle("Financial Overview"),
		termcharts.WithColor(true),
		termcharts.WithWidth(60),
		termcharts.WithHeight(12),
	)
	fmt.Println(line5.Render())

	// Convenience function examples
	fmt.Println("6. Convenience Functions")
	fmt.Println("------------------------")
	fmt.Println("Line():")
	fmt.Println(termcharts.Line([]float64{1, 4, 2, 5, 3, 6}))

	fmt.Println("LineBraille():")
	fmt.Println(termcharts.LineBraille([]float64{1, 4, 2, 5, 3, 6}))
}
