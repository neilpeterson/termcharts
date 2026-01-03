package main

import (
	"fmt"

	"github.com/neilpeterson/termcharts/pkg/termcharts"
)

func main() {
	fmt.Println("=== Basic Bar Chart ===")
	data := []float64{10, 25, 15, 30, 20}
	chart := termcharts.NewBarChart(
		termcharts.WithData(data),
	)
	fmt.Println(chart.Render())

	fmt.Println("\n=== Bar Chart with Labels ===")
	quarterlyData := []float64{45, 67, 52, 78}
	labels := []string{"Q1", "Q2", "Q3", "Q4"}
	chartWithLabels := termcharts.NewBarChart(
		termcharts.WithData(quarterlyData),
		termcharts.WithLabels(labels),
		termcharts.WithTitle("Quarterly Sales (in thousands)"),
	)
	fmt.Println(chartWithLabels.Render())

	fmt.Println("\n=== Bar Chart with Values ===")
	salesData := []float64{120.5, 98.3, 145.7, 132.1}
	salesLabels := []string{"North", "South", "East", "West"}
	chartWithValues := termcharts.NewBarChart(
		termcharts.WithData(salesData),
		termcharts.WithLabels(salesLabels),
		termcharts.WithTitle("Regional Sales"),
		termcharts.WithShowValues(true),
	)
	fmt.Println(chartWithValues.Render())

	fmt.Println("\n=== Vertical Bar Chart ===")
	monthlyData := []float64{12, 19, 25, 18, 22, 30}
	monthLabels := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
	verticalChart := termcharts.NewBarChart(
		termcharts.WithData(monthlyData),
		termcharts.WithLabels(monthLabels),
		termcharts.WithTitle("Monthly Activity"),
		termcharts.WithDirection(termcharts.Vertical),
		termcharts.WithHeight(12),
	)
	fmt.Println(verticalChart.Render())

	fmt.Println("\n=== ASCII Mode Bar Chart ===")
	asciiData := []float64{8, 15, 12, 20}
	asciiChart := termcharts.NewBarChart(
		termcharts.WithData(asciiData),
		termcharts.WithLabels([]string{"A", "B", "C", "D"}),
		termcharts.WithStyle(termcharts.StyleASCII),
		termcharts.WithTitle("ASCII Bar Chart"),
	)
	fmt.Println(asciiChart.Render())

	fmt.Println("\n=== Convenience Functions ===")
	fmt.Println("Bar():")
	fmt.Println(termcharts.Bar([]float64{5, 10, 15, 20}))

	fmt.Println("\nBarWithLabels():")
	fmt.Println(termcharts.BarWithLabels(
		[]float64{25, 40, 35, 50},
		[]string{"Alpha", "Beta", "Gamma", "Delta"},
	))

	fmt.Println("\nBarVertical():")
	fmt.Println(termcharts.BarVertical([]float64{10, 20, 15, 25, 18}))
}
