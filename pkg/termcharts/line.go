package termcharts

import (
	"fmt"
	"strings"

	"github.com/neilpeterson/termcharts/internal"
)

// LineChart represents a line chart visualization.
// Line charts can be rendered using ASCII box-drawing characters,
// Unicode characters, or high-resolution Braille patterns.
type LineChart struct {
	opts *Options
}

// Box-drawing characters for ASCII line rendering.
const (
	lineHorizontal = '─'
	lineVertical   = '│'
	lineCornerTL   = '╭'
	lineCornerTR   = '╮'
	lineCornerBL   = '╰'
	lineCornerBR   = '╯'
	lineUp         = '╱'
	lineDown       = '╲'
	lineDot        = '•'
)

// ASCII fallback characters.
const (
	asciiHorizontal = '-'
	asciiVertical   = '|'
	asciiUp         = '/'
	asciiDown       = '\\'
	asciiDot        = '*'
)

// Braille patterns for high-resolution rendering.
// Braille characters use a 2x4 dot matrix per character cell.
// Pattern: dots are numbered 1-8:
// 1 4
// 2 5
// 3 6
// 7 8
const brailleBase = 0x2800 // Unicode Braille pattern blank

// brailleDots maps (row, col) within a cell to the bit position.
// Row 0-3 (top to bottom), Col 0-1 (left to right).
var brailleDots = [4][2]int{
	{0x01, 0x08}, // Row 0: dots 1, 4
	{0x02, 0x10}, // Row 1: dots 2, 5
	{0x04, 0x20}, // Row 2: dots 3, 6
	{0x40, 0x80}, // Row 3: dots 7, 8
}

// NewLineChart creates a new line chart with the given options.
// At minimum, data must be provided via WithData option or WithSeries for multi-series.
//
// Example:
//
//	line := termcharts.NewLineChart(
//	    termcharts.WithData([]float64{1, 5, 2, 8, 3, 7}),
//	    termcharts.WithWidth(60),
//	    termcharts.WithHeight(10),
//	)
//	fmt.Println(line.Render())
func NewLineChart(opts ...Option) *LineChart {
	options := NewOptions(opts...)
	return &LineChart{
		opts: options,
	}
}

// Render generates the line chart as a multi-line string.
func (l *LineChart) Render() string {
	// Get all data series
	allSeries := l.getAllSeries()
	if len(allSeries) == 0 {
		return ""
	}

	// Check for invalid values
	for _, series := range allSeries {
		if !internal.AllValid(series.Data) {
			return ""
		}
	}

	// Render based on style
	if l.opts.Style == StyleBraille {
		return l.renderBraille(allSeries)
	}
	return l.renderASCII(allSeries)
}

// getAllSeries returns all data series to render.
func (l *LineChart) getAllSeries() []Series {
	// If explicit series are provided, use them
	if len(l.opts.Series) > 0 {
		return l.opts.Series
	}

	// Otherwise, create a single series from the data
	if len(l.opts.Data) > 0 {
		return []Series{{
			Label: "",
			Data:  l.opts.Data,
			Color: "",
		}}
	}

	return nil
}

// renderASCII renders the line chart using ASCII/Unicode box-drawing characters.
//
//nolint:gocyclo // Complex rendering logic
func (l *LineChart) renderASCII(allSeries []Series) string {
	// Determine dimensions
	width := l.opts.Width
	height := l.opts.Height

	// Reserve space for title and axes
	chartHeight := height
	if l.opts.Title != "" {
		chartHeight--
	}
	if l.opts.ShowAxes {
		chartHeight -= 2 // Bottom axis and labels
	}
	if chartHeight < 3 {
		chartHeight = 10
	}

	// Calculate chart width (leave room for Y axis if showing)
	chartWidth := width
	yAxisWidth := 0
	if l.opts.ShowAxes {
		yAxisWidth = 8 // Space for Y axis labels
		chartWidth -= yAxisWidth
	}
	if chartWidth < 10 {
		chartWidth = 60
	}

	// Find global min/max across all series
	globalMin, globalMax := l.findGlobalMinMax(allSeries)
	if globalMin == globalMax {
		globalMax = globalMin + 1
	}

	// Get styling
	useUnicode := l.shouldUseUnicode()
	colorEnabled := l.isColorEnabled()
	theme := l.opts.Theme
	if theme == nil {
		theme = DefaultTheme
	}

	// Create the chart grid
	grid := make([][]rune, chartHeight)
	colors := make([][]string, chartHeight)
	for i := range grid {
		grid[i] = make([]rune, chartWidth)
		colors[i] = make([]string, chartWidth)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	// Render each series
	for seriesIdx, series := range allSeries {
		color := series.Color
		if color == "" {
			color = theme.GetSeriesColor(seriesIdx)
		}

		l.renderSeriesASCII(grid, colors, series.Data, chartWidth, chartHeight, globalMin, globalMax, useUnicode, color)
	}

	// Build result
	var result strings.Builder

	// Render title if provided
	if l.opts.Title != "" {
		titleText := l.opts.Title
		if colorEnabled {
			titleText = Colorize(titleText, theme.Text, true)
		}
		result.WriteString(titleText)
		result.WriteString("\n")
	}

	// Render chart rows
	for row := 0; row < chartHeight; row++ {
		// Y axis label
		if l.opts.ShowAxes {
			// Calculate value at this row
			rowValue := globalMax - (float64(row)/float64(chartHeight-1))*(globalMax-globalMin)
			label := fmt.Sprintf("%7.1f ", rowValue)
			if colorEnabled {
				label = Colorize(label, theme.Muted, true)
			}
			result.WriteString(label)
		}

		// Chart content
		for col := 0; col < chartWidth; col++ {
			char := string(grid[row][col])
			if colorEnabled && colors[row][col] != "" {
				char = Colorize(char, colors[row][col], true)
			}
			result.WriteString(char)
		}
		result.WriteString("\n")
	}

	// Render X axis if showing axes
	if l.opts.ShowAxes {
		// Axis line
		if yAxisWidth > 0 {
			result.WriteString(strings.Repeat(" ", yAxisWidth))
		}
		axisLine := strings.Repeat("─", chartWidth)
		if !useUnicode {
			axisLine = strings.Repeat("-", chartWidth)
		}
		if colorEnabled {
			axisLine = Colorize(axisLine, theme.Muted, true)
		}
		result.WriteString(axisLine)
		result.WriteString("\n")

		// X axis labels
		if len(l.opts.Labels) > 0 {
			if yAxisWidth > 0 {
				result.WriteString(strings.Repeat(" ", yAxisWidth))
			}
			l.renderXAxisLabels(&result, chartWidth, colorEnabled, theme)
			result.WriteString("\n")
		}
	}

	// Render legend for multi-series
	if len(allSeries) > 1 {
		result.WriteString("\n")
		for i, series := range allSeries {
			color := series.Color
			if color == "" {
				color = theme.GetSeriesColor(i)
			}
			marker := "●"
			if !useUnicode {
				marker = "*"
			}
			if colorEnabled {
				marker = Colorize(marker, color, true)
			}
			label := series.Label
			if label == "" {
				label = fmt.Sprintf("Series %d", i+1)
			}
			result.WriteString(fmt.Sprintf("%s %s  ", marker, label))
		}
		result.WriteString("\n")
	}

	return result.String()
}

// renderSeriesASCII renders a single data series onto the grid.
func (l *LineChart) renderSeriesASCII(grid [][]rune, colors [][]string, data []float64, width, height int, minVal, maxVal float64, useUnicode bool, color string) {
	if len(data) == 0 {
		return
	}

	// Map data points to grid coordinates
	points := make([][2]int, len(data))
	for i, val := range data {
		// X position: spread across width
		x := int(float64(i) / float64(len(data)-1) * float64(width-1))
		if len(data) == 1 {
			x = width / 2
		}

		// Y position: scale to height (0 = top, height-1 = bottom)
		y := int((maxVal - val) / (maxVal - minVal) * float64(height-1))
		y = internal.ClampInt(y, 0, height-1)
		x = internal.ClampInt(x, 0, width-1)

		points[i] = [2]int{x, y}
	}

	// Draw lines between consecutive points
	for i := 0; i < len(points)-1; i++ {
		x1, y1 := points[i][0], points[i][1]
		x2, y2 := points[i+1][0], points[i+1][1]

		l.drawLine(grid, colors, x1, y1, x2, y2, useUnicode, color)
	}

	// Draw data points
	for _, p := range points {
		x, y := p[0], p[1]
		if useUnicode {
			grid[y][x] = lineDot
		} else {
			grid[y][x] = asciiDot
		}
		colors[y][x] = color
	}
}

// drawLine draws a line between two points using Bresenham-style algorithm.
func (l *LineChart) drawLine(grid [][]rune, colors [][]string, x1, y1, x2, y2 int, useUnicode bool, color string) {
	dx := internal.Abs(x2 - x1)
	dy := internal.Abs(y2 - y1)

	sx := 1
	if x1 > x2 {
		sx = -1
	}
	sy := 1
	if y1 > y2 {
		sy = -1
	}

	err := dx - dy

	x, y := x1, y1
	for {
		// Choose character based on direction
		char := l.getLineChar(x, y, x1, y1, x2, y2, useUnicode)
		if grid[y][x] == ' ' || grid[y][x] == lineHorizontal || grid[y][x] == asciiHorizontal {
			grid[y][x] = char
			colors[y][x] = color
		}

		if x == x2 && y == y2 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

// getLineChar returns the appropriate character for a line segment.
func (l *LineChart) getLineChar(x, y, x1, y1, x2, y2 int, useUnicode bool) rune {
	dy := y2 - y1
	dx := x2 - x1

	if dx == 0 {
		// Vertical line
		if useUnicode {
			return lineVertical
		}
		return asciiVertical
	}

	if dy == 0 {
		// Horizontal line
		if useUnicode {
			return lineHorizontal
		}
		return asciiHorizontal
	}

	// Diagonal line
	if (dy > 0 && dx > 0) || (dy < 0 && dx < 0) {
		// Going down-right or up-left
		if useUnicode {
			return lineDown
		}
		return asciiDown
	}
	// Going up-right or down-left
	if useUnicode {
		return lineUp
	}
	return asciiUp
}

// renderXAxisLabels renders X axis labels.
func (l *LineChart) renderXAxisLabels(result *strings.Builder, width int, colorEnabled bool, theme *Theme) {
	labels := l.opts.Labels
	if len(labels) == 0 {
		return
	}

	// Distribute labels across width
	labelPositions := make([]int, len(labels))
	for i := range labels {
		labelPositions[i] = int(float64(i) / float64(len(labels)-1) * float64(width-1))
		if len(labels) == 1 {
			labelPositions[i] = width / 2
		}
	}

	// Build label line
	line := make([]byte, width)
	for i := range line {
		line[i] = ' '
	}

	for i, label := range labels {
		pos := labelPositions[i]
		// Center the label around the position
		start := pos - len(label)/2
		if start < 0 {
			start = 0
		}
		if start+len(label) > width {
			start = width - len(label)
		}
		for j, c := range label {
			if start+j < width {
				line[start+j] = byte(c)
			}
		}
	}

	text := string(line)
	if colorEnabled {
		text = Colorize(text, theme.Muted, true)
	}
	result.WriteString(text)
}

// renderBraille renders the line chart using high-resolution Braille patterns.
//
//nolint:gocyclo // Complex rendering logic
func (l *LineChart) renderBraille(allSeries []Series) string {
	// Determine dimensions
	width := l.opts.Width
	height := l.opts.Height

	// Reserve space for title
	chartHeight := height
	if l.opts.Title != "" {
		chartHeight--
	}
	if l.opts.ShowAxes {
		chartHeight -= 2
	}
	if chartHeight < 3 {
		chartHeight = 10
	}

	// Calculate chart width
	chartWidth := width
	yAxisWidth := 0
	if l.opts.ShowAxes {
		yAxisWidth = 8
		chartWidth -= yAxisWidth
	}
	if chartWidth < 10 {
		chartWidth = 60
	}

	// Braille resolution: each character is 2x4 dots
	brailleWidth := chartWidth
	brailleHeight := chartHeight * 4 // 4 vertical dots per character

	// Find global min/max
	globalMin, globalMax := l.findGlobalMinMax(allSeries)
	if globalMin == globalMax {
		globalMax = globalMin + 1
	}

	// Get styling
	colorEnabled := l.isColorEnabled()
	theme := l.opts.Theme
	if theme == nil {
		theme = DefaultTheme
	}

	// Create Braille dot grid
	dotGrid := make([][]bool, brailleHeight)
	for i := range dotGrid {
		dotGrid[i] = make([]bool, brailleWidth*2) // 2 horizontal dots per char
	}

	// Create color grid for each character cell
	colorGrid := make([][]string, chartHeight)
	for i := range colorGrid {
		colorGrid[i] = make([]string, chartWidth)
	}

	// Render each series
	for seriesIdx, series := range allSeries {
		color := series.Color
		if color == "" {
			color = theme.GetSeriesColor(seriesIdx)
		}

		l.renderSeriesBraille(dotGrid, colorGrid, series.Data, brailleWidth*2, brailleHeight, chartWidth, chartHeight, globalMin, globalMax, color)
	}

	// Build result
	var result strings.Builder

	// Render title if provided
	if l.opts.Title != "" {
		titleText := l.opts.Title
		if colorEnabled {
			titleText = Colorize(titleText, theme.Text, true)
		}
		result.WriteString(titleText)
		result.WriteString("\n")
	}

	// Convert dot grid to Braille characters
	for row := 0; row < chartHeight; row++ {
		// Y axis label
		if l.opts.ShowAxes {
			rowValue := globalMax - (float64(row)/float64(chartHeight-1))*(globalMax-globalMin)
			label := fmt.Sprintf("%7.1f ", rowValue)
			if colorEnabled {
				label = Colorize(label, theme.Muted, true)
			}
			result.WriteString(label)
		}

		// Chart content
		for col := 0; col < chartWidth; col++ {
			// Calculate Braille pattern for this cell
			pattern := 0
			for dotRow := 0; dotRow < 4; dotRow++ {
				for dotCol := 0; dotCol < 2; dotCol++ {
					gridRow := row*4 + dotRow
					gridCol := col*2 + dotCol
					if gridRow < brailleHeight && gridCol < brailleWidth*2 {
						if dotGrid[gridRow][gridCol] {
							pattern |= brailleDots[dotRow][dotCol]
						}
					}
				}
			}

			char := string(rune(brailleBase + pattern))
			if colorEnabled && colorGrid[row][col] != "" {
				char = Colorize(char, colorGrid[row][col], true)
			}
			result.WriteString(char)
		}
		result.WriteString("\n")
	}

	// Render X axis if showing axes
	if l.opts.ShowAxes {
		if yAxisWidth > 0 {
			result.WriteString(strings.Repeat(" ", yAxisWidth))
		}
		axisLine := strings.Repeat("─", chartWidth)
		if colorEnabled {
			axisLine = Colorize(axisLine, theme.Muted, true)
		}
		result.WriteString(axisLine)
		result.WriteString("\n")

		if len(l.opts.Labels) > 0 {
			if yAxisWidth > 0 {
				result.WriteString(strings.Repeat(" ", yAxisWidth))
			}
			l.renderXAxisLabels(&result, chartWidth, colorEnabled, theme)
			result.WriteString("\n")
		}
	}

	// Render legend for multi-series
	if len(allSeries) > 1 {
		result.WriteString("\n")
		for i, series := range allSeries {
			color := series.Color
			if color == "" {
				color = theme.GetSeriesColor(i)
			}
			marker := "●"
			if colorEnabled {
				marker = Colorize(marker, color, true)
			}
			label := series.Label
			if label == "" {
				label = fmt.Sprintf("Series %d", i+1)
			}
			result.WriteString(fmt.Sprintf("%s %s  ", marker, label))
		}
		result.WriteString("\n")
	}

	return result.String()
}

// renderSeriesBraille renders a single data series onto the Braille dot grid.
func (l *LineChart) renderSeriesBraille(dotGrid [][]bool, colorGrid [][]string, data []float64, dotWidth, dotHeight, charWidth, charHeight int, minVal, maxVal float64, color string) {
	if len(data) == 0 {
		return
	}

	// Map data points to dot coordinates
	for i := 0; i < len(data)-1; i++ {
		// Start point
		x1 := int(float64(i) / float64(len(data)-1) * float64(dotWidth-1))
		y1 := int((maxVal - data[i]) / (maxVal - minVal) * float64(dotHeight-1))
		y1 = internal.ClampInt(y1, 0, dotHeight-1)
		x1 = internal.ClampInt(x1, 0, dotWidth-1)

		// End point
		x2 := int(float64(i+1) / float64(len(data)-1) * float64(dotWidth-1))
		y2 := int((maxVal - data[i+1]) / (maxVal - minVal) * float64(dotHeight-1))
		y2 = internal.ClampInt(y2, 0, dotHeight-1)
		x2 = internal.ClampInt(x2, 0, dotWidth-1)

		// Draw line between points using Bresenham
		l.drawBrailleLine(dotGrid, colorGrid, x1, y1, x2, y2, charWidth, charHeight, color)
	}

	// Ensure single point is drawn
	if len(data) == 1 {
		x := dotWidth / 2
		y := int((maxVal - data[0]) / (maxVal - minVal) * float64(dotHeight-1))
		y = internal.ClampInt(y, 0, dotHeight-1)
		dotGrid[y][x] = true
		colorGrid[y/4][x/2] = color
	}
}

// drawBrailleLine draws a line on the Braille dot grid.
func (l *LineChart) drawBrailleLine(dotGrid [][]bool, colorGrid [][]string, x1, y1, x2, y2, charWidth, charHeight int, color string) {
	dx := internal.Abs(x2 - x1)
	dy := internal.Abs(y2 - y1)

	sx := 1
	if x1 > x2 {
		sx = -1
	}
	sy := 1
	if y1 > y2 {
		sy = -1
	}

	err := dx - dy

	x, y := x1, y1
	for {
		if y >= 0 && y < len(dotGrid) && x >= 0 && x < len(dotGrid[0]) {
			dotGrid[y][x] = true
			// Set color for the character cell
			charRow := y / 4
			charCol := x / 2
			if charRow < charHeight && charCol < charWidth {
				colorGrid[charRow][charCol] = color
			}
		}

		if x == x2 && y == y2 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

// findGlobalMinMax finds the min and max values across all series.
func (l *LineChart) findGlobalMinMax(allSeries []Series) (float64, float64) {
	var allData []float64
	for _, series := range allSeries {
		allData = append(allData, series.Data...)
	}
	return internal.MinMax(allData)
}

// shouldUseUnicode determines whether to use Unicode characters.
func (l *LineChart) shouldUseUnicode() bool {
	if l.opts.Style == StyleASCII {
		return false
	}
	if l.opts.Style == StyleUnicode || l.opts.Style == StyleBraille {
		return true
	}
	return internal.SupportsUnicode()
}

// isColorEnabled determines whether colors should be used.
func (l *LineChart) isColorEnabled() bool {
	if l.opts.ColorEnabled != nil {
		return *l.opts.ColorEnabled
	}
	return internal.SupportsColor()
}

// Line is a convenience function that creates and renders a line chart.
//
// Example:
//
//	fmt.Println(termcharts.Line([]float64{1, 5, 2, 8, 3, 7}))
func Line(data []float64) string {
	line := NewLineChart(
		WithData(data),
		WithHeight(10),
		WithWidth(60),
	)
	return line.Render()
}

// LineBraille creates a high-resolution line chart using Braille patterns.
//
// Example:
//
//	fmt.Println(termcharts.LineBraille([]float64{1, 5, 2, 8, 3, 7}))
func LineBraille(data []float64) string {
	line := NewLineChart(
		WithData(data),
		WithHeight(10),
		WithWidth(60),
		WithStyle(StyleBraille),
	)
	return line.Render()
}

// LineMultiSeries creates a line chart with multiple data series.
//
// Example:
//
//	series := []termcharts.Series{
//	    {Label: "Sales", Data: []float64{10, 20, 15, 25}},
//	    {Label: "Costs", Data: []float64{8, 12, 10, 18}},
//	}
//	fmt.Println(termcharts.LineMultiSeries(series))
func LineMultiSeries(series []Series) string {
	line := NewLineChart(
		WithSeries(series),
		WithHeight(10),
		WithWidth(60),
	)
	return line.Render()
}
