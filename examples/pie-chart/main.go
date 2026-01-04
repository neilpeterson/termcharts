package main

import (
	"fmt"

	"github.com/neilpeterson/termcharts/pkg/termcharts"
)

func main() {
	fmt.Println("=== Basic Pie Chart ===")
	fmt.Println(termcharts.Pie([]float64{30, 25, 20, 15, 10}))

	fmt.Println("\n=== Pie Chart with Labels ===")
	fmt.Println(termcharts.PieWithLabels(
		[]float64{30, 25, 20, 15, 10},
		[]string{"Chrome", "Firefox", "Safari", "Edge", "Other"},
	))

	fmt.Println("\n=== Pie Chart with Title and Labels ===")
	pie := termcharts.NewPieChart(
		termcharts.WithData([]float64{45, 30, 15, 10}),
		termcharts.WithLabels([]string{"North", "South", "East", "West"}),
		termcharts.WithTitle("Sales by Region"),
	)
	fmt.Println(pie.Render())

	fmt.Println("\n=== Pie Chart with Values ===")
	pieWithValues := termcharts.NewPieChart(
		termcharts.WithData([]float64{1250.50, 980.25, 750.00, 520.75}),
		termcharts.WithLabels([]string{"Product A", "Product B", "Product C", "Product D"}),
		termcharts.WithTitle("Revenue by Product ($)"),
		termcharts.WithShowValues(true),
	)
	fmt.Println(pieWithValues.Render())

	fmt.Println("\n=== Colored Pie Chart ===")
	coloredPie := termcharts.NewPieChart(
		termcharts.WithData([]float64{40, 35, 25}),
		termcharts.WithLabels([]string{"Desktop", "Mobile", "Tablet"}),
		termcharts.WithTitle("Traffic by Device"),
		termcharts.WithColor(true),
	)
	fmt.Println(coloredPie.Render())

	fmt.Println("\n=== ASCII Mode Pie Chart ===")
	asciiPie := termcharts.NewPieChart(
		termcharts.WithData([]float64{50, 30, 20}),
		termcharts.WithLabels([]string{"Yes", "No", "Maybe"}),
		termcharts.WithTitle("Survey Results"),
		termcharts.WithStyle(termcharts.StyleASCII),
	)
	fmt.Println(asciiPie.Render())

	fmt.Println("\n=== Custom Width Pie Chart ===")
	narrowPie := termcharts.NewPieChart(
		termcharts.WithData([]float64{60, 25, 15}),
		termcharts.WithLabels([]string{"Majority", "Minority", "Undecided"}),
		termcharts.WithWidth(50),
	)
	fmt.Println(narrowPie.Render())

	fmt.Println("\n=== Dark Theme Pie Chart ===")
	darkPie := termcharts.NewPieChart(
		termcharts.WithData([]float64{35, 30, 20, 15}),
		termcharts.WithLabels([]string{"Q1", "Q2", "Q3", "Q4"}),
		termcharts.WithTitle("Quarterly Distribution"),
		termcharts.WithTheme(termcharts.DarkTheme),
		termcharts.WithColor(true),
	)
	fmt.Println(darkPie.Render())

	fmt.Println("\n=== Convenience Functions ===")
	fmt.Println("Pie():")
	fmt.Println(termcharts.Pie([]float64{50, 30, 20}))

	fmt.Println("\nPieWithLabels():")
	fmt.Println(termcharts.PieWithLabels(
		[]float64{50, 30, 20},
		[]string{"Alpha", "Beta", "Gamma"},
	))

	fmt.Println("\nPieWithValues():")
	fmt.Println(termcharts.PieWithValues(
		[]float64{100, 75, 50},
		[]string{"Large", "Medium", "Small"},
	))
}
