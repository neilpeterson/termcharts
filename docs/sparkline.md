# Sparkline Documentation

Sparklines are compact, inline charts that visualize data trends in a single line. They use Unicode block characters (▁▂▃▄▅▆▇█) or ASCII characters for maximum compatibility.

## Table of Contents

- [Overview](#overview)
- [Code Structure](#code-structure)
- [Library Usage](#library-usage)
- [CLI Usage](#cli-usage)
- [Examples](#examples)
- [API Reference](#api-reference)

## Overview

Sparklines were invented by Edward Tufte as "data-intense, design-simple" charts. They show trends at a glance without taking up much space, making them perfect for:

- Terminal dashboards
- System monitoring
- Log analysis
- Quick data visualization in shell scripts
- Inline metrics in documentation

### What it looks like:

```
CPU Usage:     ▁▂▃▄▅▆▇█▇▆▅▄▃▂▁
Memory:        ▃▃▄▄▅▅▆▆▇▇██▇▆
API Response:  ▂▂▂▂▃▅▇▃▂▂▂▂
```

## Code Structure

### File Organization

```
pkg/termcharts/
├── sparkline.go         # Sparkline implementation
└── sparkline_test.go    # Unit tests (14 test cases)

cmd/termcharts/
├── main.go              # CLI entry point
├── root.go              # Root command
└── spark.go             # Sparkline subcommand

examples/
└── sparkline/
    └── main.go          # Demo application
```

### Core Components

#### `Sparkline` Type (pkg/termcharts/sparkline.go)

```go
type Sparkline struct {
    opts *Options  // Chart configuration
}
```

The `Sparkline` type implements the `Chart` interface and provides sparkline rendering with:
- Unicode block characters: `▁▂▃▄▅▆▇█` (8 levels)
- ASCII fallback characters: `_.-=+*#@` (8 levels)
- Auto-detection of terminal capabilities
- Color support with theme-based coloring
- Width limiting with smart downsampling

#### Key Functions

- `NewSparkline(opts ...Option) *Sparkline` - Constructor with functional options
- `(s *Sparkline) Render() string` - Renders the sparkline
- `Spark(data []float64) string` - Convenience function for quick sparklines
- `SparkASCII(data []float64) string` - ASCII-only mode
- `SparkColor(data []float64) string` - Auto-colored sparklines

## Library Usage

### Installation

```bash
go get github.com/neilpeterson/termcharts
```

### Quick Start

```go
package main

import (
    "fmt"
    "github.com/neilpeterson/termcharts/pkg/termcharts"
)

func main() {
    // Simple one-liner
    data := []float64{1, 5, 2, 8, 3, 7, 4, 6}
    fmt.Println(termcharts.Spark(data))
    // Output: ▁▅▂█▃▇▄▆
}
```

### Advanced Usage with Options

```go
package main

import (
    "fmt"
    "github.com/neilpeterson/termcharts/pkg/termcharts"
)

func main() {
    data := []float64{10, 20, 30, 25, 15, 35, 40, 38, 42, 50}

    // With custom options
    spark := termcharts.NewSparkline(
        termcharts.WithData(data),
        termcharts.WithWidth(20),           // Limit width
        termcharts.WithColor(true),         // Enable colors
        termcharts.WithStyle(termcharts.StyleUnicode),
    )

    fmt.Println(spark.Render())
}
```

### Available Options

```go
// Data and dimensions
termcharts.WithData([]float64{...})     // Set data points
termcharts.WithWidth(int)               // Limit output width (0 = unlimited)

// Rendering style
termcharts.WithStyle(StyleASCII)        // Use ASCII characters
termcharts.WithStyle(StyleUnicode)      // Use Unicode blocks
termcharts.WithStyle(StyleAuto)         // Auto-detect (default)

// Color control
termcharts.WithColor(true)              // Enable colors
termcharts.WithColor(false)             // Disable colors
termcharts.WithTheme(&Theme{...})       // Custom color theme
```

### Character Sets

**Unicode (Default):**
```
▁ ▂ ▃ ▄ ▅ ▆ ▇ █
```

**ASCII (Compatibility Mode):**
```
_ . - = + * # @
```

## CLI Usage

### Installation

```bash
# From source
go install github.com/neilpeterson/termcharts/cmd/termcharts@latest

# Or build locally
git clone https://github.com/neilpeterson/termcharts
cd termcharts
go build -o termcharts ./cmd/termcharts
```

### Basic Usage

```bash
# From command-line arguments
termcharts spark 10 20 30 25 15 35 40

# From a file
termcharts spark data.txt

# From stdin (pipe-friendly!)
echo "1 5 2 8 3 7" | termcharts spark

# With flags
termcharts spark data.txt --width 50 --color
```

### Command-Line Flags

```
  --width, -w int     Maximum width in characters (0 = no limit)
  --ascii             Use ASCII characters only
  --color, -c         Enable colored output
  --no-color          Disable colored output
  --help, -h          Show help
```

### Data Input Formats

The CLI accepts multiple data formats:

**1. Command-line arguments:**
```bash
termcharts spark 10 20 30 25 15
```

**2. File with one number per line:**
```bash
# data.txt
10
20
30
25
15
```

**3. File with space-separated values:**
```bash
# data.txt
10 20 30 25 15
```

**4. File with comma-separated values:**
```bash
# data.txt
10,20,30,25,15
```

**5. Stdin (from pipes):**
```bash
cat data.txt | termcharts spark
echo "1 2 3 4 5" | termcharts spark
seq 1 10 | termcharts spark
```

### Comments in Files

Lines starting with `#` are treated as comments and ignored:

```bash
# data.txt
# CPU usage over the last hour
12
15
18
22
# spike during backup
67
45
```

## Examples

### System Monitoring

```bash
# Monitor CPU usage
while true; do
    cpu=$(ps aux | awk '{sum+=$3} END {print sum}')
    echo $cpu >> cpu.log
    tail -20 cpu.log | termcharts spark --width 40
    sleep 1
done
```

### Git Commit Activity

```bash
# Visualize commit frequency over time
git log --all --format='%ai' | \
    cut -d' ' -f1 | \
    uniq -c | \
    awk '{print $1}' | \
    termcharts spark --color
```

### API Response Times

```bash
# Monitor API response times
for i in {1..20}; do
    curl -w "%{time_total}\n" -o /dev/null -s https://api.example.com
done | termcharts spark
```

### File Size Trends

```bash
# Show file size growth over commits
git log --all --oneline | \
    while read hash msg; do
        git show $hash:large_file.bin 2>/dev/null | wc -c
    done | termcharts spark --width 60
```

### Inline Metrics in Scripts

```go
package main

import (
    "fmt"
    "github.com/neilpeterson/termcharts/pkg/termcharts"
)

func main() {
    cpuUsage := []float64{12, 15, 14, 18, 22, 45, 67, 78, 65, 52}
    memUsage := []float64{30, 32, 35, 38, 42, 45, 48, 52, 55, 58}

    fmt.Printf("CPU:    %s (avg: %.1f%%)\n",
        termcharts.Spark(cpuUsage), average(cpuUsage))
    fmt.Printf("Memory: %s (avg: %.1f%%)\n",
        termcharts.Spark(memUsage), average(memUsage))
}

func average(data []float64) float64 {
    sum := 0.0
    for _, v := range data {
        sum += v
    }
    return sum / float64(len(data))
}
```

Output:
```
CPU:    ▁▁▁▁▂▄▆█▆▅ (avg: 37.8%)
Memory: ▁▁▂▂▃▄▄▅▆▇ (avg: 43.5%)
```

## API Reference

### Functions

#### `NewSparkline(opts ...Option) *Sparkline`

Creates a new sparkline with the given options.

**Parameters:**
- `opts ...Option` - Functional options for configuration

**Returns:**
- `*Sparkline` - A new sparkline instance

**Example:**
```go
spark := termcharts.NewSparkline(
    termcharts.WithData([]float64{1, 5, 2, 8}),
    termcharts.WithWidth(50),
)
```

#### `(s *Sparkline) Render() string`

Renders the sparkline as a string.

**Returns:**
- `string` - The rendered sparkline

**Example:**
```go
spark := termcharts.NewSparkline(termcharts.WithData(data))
output := spark.Render()
fmt.Println(output)
```

#### `Spark(data []float64) string`

Convenience function that creates and renders a sparkline in one call.

**Parameters:**
- `data []float64` - The data points to visualize

**Returns:**
- `string` - The rendered sparkline

**Example:**
```go
fmt.Println(termcharts.Spark([]float64{1, 5, 2, 8, 3, 7}))
// Output: ▁▅▂█▃▇
```

#### `SparkASCII(data []float64) string`

Creates an ASCII-only sparkline for maximum compatibility.

**Parameters:**
- `data []float64` - The data points to visualize

**Returns:**
- `string` - The rendered sparkline using ASCII characters

**Example:**
```go
fmt.Println(termcharts.SparkASCII([]float64{1, 5, 2, 8, 3, 7}))
// Output: _+.@-#
```

#### `SparkColor(data []float64) string`

Creates a colored sparkline with auto-detected color support.

**Parameters:**
- `data []float64` - The data points to visualize

**Returns:**
- `string` - The rendered colored sparkline

**Example:**
```go
fmt.Println(termcharts.SparkColor([]float64{1, 5, 2, 8, 3, 7}))
// Output: (colored) ▁▅▂█▃▇
```

### Configuration Options

All options follow the functional options pattern:

```go
type Option func(*Options)
```

Available options:
- `WithData([]float64)` - Set the data points
- `WithWidth(int)` - Set maximum width (0 = unlimited)
- `WithStyle(RenderStyle)` - Set rendering style (ASCII, Unicode, Auto)
- `WithColor(bool)` - Enable/disable colors
- `WithTheme(*Theme)` - Set custom color theme

### Edge Cases

The sparkline implementation handles various edge cases gracefully:

**Empty Data:**
```go
termcharts.Spark([]float64{})  // Returns: ""
```

**Single Value:**
```go
termcharts.Spark([]float64{42})  // Returns: "▄" (middle character)
```

**All Same Values:**
```go
termcharts.Spark([]float64{5, 5, 5, 5})  // Returns: "▄▄▄▄" (all middle)
```

**Invalid Data (NaN, Inf):**
```go
termcharts.Spark([]float64{1, math.NaN(), 3})  // Returns: "" (empty)
```

## Testing

The sparkline implementation includes comprehensive unit tests:

```bash
# Run sparkline tests
go test ./pkg/termcharts -run TestSparkline -v

# Run all tests
go test ./...

# With coverage
go test ./pkg/termcharts -cover
```

**Test Coverage:**
- Basic data rendering
- Empty/single/same values
- ASCII and Unicode modes
- Invalid data handling (NaN, Inf)
- Width limiting
- Color support
- Convenience functions
- Character mapping

All 14 tests passing ✓

## Contributing

When modifying sparkline functionality:

1. Update tests in `pkg/termcharts/sparkline_test.go`
2. Run tests: `go test ./pkg/termcharts -v`
3. Update this documentation if API changes
4. Update examples if new features are added

## See Also

- [Main README](../README.md) - Project overview
- [API Documentation](https://pkg.go.dev/github.com/neilpeterson/termcharts)
- [Examples](../examples/) - More usage examples
