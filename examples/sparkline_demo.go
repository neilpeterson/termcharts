// Last modified: 2026-01-02

package main

import (
	"fmt"

	"github.com/neilpeterson/termcharts/pkg/termcharts"
)

func main() {
	fmt.Println("Sparkline Demo")
	fmt.Println("==============")

	// Simple sparkline
	data1 := []float64{1, 5, 2, 8, 3, 7, 4, 6}
	fmt.Printf("Simple:        %s\n", termcharts.Spark(data1))

	// ASCII sparkline
	fmt.Printf("ASCII:         %s\n", termcharts.SparkASCII(data1))

	// Colored sparkline
	fmt.Printf("Color:         %s\n", termcharts.SparkColor(data1))

	// More data points
	data2 := []float64{10, 20, 30, 25, 15, 35, 40, 38, 42, 50, 45, 40, 35, 30, 25, 20}
	fmt.Printf("\nTrend:         %s\n", termcharts.Spark(data2))

	// CPU usage simulation
	cpu := []float64{12, 15, 14, 18, 22, 45, 67, 78, 65, 52, 38, 25, 18, 15, 12}
	fmt.Printf("CPU Usage:     %s\n", termcharts.Spark(cpu))

	// Stock price simulation
	stock := []float64{100, 102, 98, 105, 110, 108, 115, 112, 118, 122, 120, 125}
	fmt.Printf("Stock Price:   %s\n", termcharts.Spark(stock))

	// Memory usage
	memory := []float64{30, 32, 35, 38, 42, 45, 48, 52, 55, 58, 60, 62, 65, 68, 70}
	fmt.Printf("Memory Usage:  %s\n", termcharts.Spark(memory))

	// Using NewSparkline with options
	fmt.Println("\nWith Options:")
	spark := termcharts.NewSparkline(
		termcharts.WithData(data1),
		termcharts.WithWidth(20),
	)
	fmt.Printf("Limited Width: %s\n", spark.Render())

	// Edge cases
	fmt.Println("\nEdge Cases:")
	single := []float64{42}
	fmt.Printf("Single value:  %s\n", termcharts.Spark(single))

	same := []float64{5, 5, 5, 5, 5}
	fmt.Printf("All same:      %s\n", termcharts.Spark(same))

	empty := []float64{}
	fmt.Printf("Empty:         '%s'\n", termcharts.Spark(empty))
}
