package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestCLI_Bar tests the bar command.
func TestCLI_Bar(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		contains []string
	}{
		{
			name:    "simple bar chart",
			args:    []string{"bar", "10", "20", "30", "40"},
			wantErr: false,
		},
		{
			name:     "bar chart with labels",
			args:     []string{"bar", "10", "20", "30", "--labels", "A,B,C"},
			wantErr:  false,
			contains: []string{"A", "B", "C"},
		},
		{
			name:     "bar chart with title",
			args:     []string{"bar", "10", "20", "30", "--title", "Test Title"},
			wantErr:  false,
			contains: []string{"Test Title"},
		},
		{
			name:    "vertical bar chart",
			args:    []string{"bar", "10", "20", "30", "--vertical"},
			wantErr: false,
		},
		{
			name:    "ascii mode",
			args:    []string{"bar", "10", "20", "30", "--ascii"},
			wantErr: false,
			contains: []string{"#"},
		},
		{
			name:    "no data error",
			args:    []string{"bar"},
			wantErr: true,
		},
	}

	binary := buildBinary(t)
	defer os.Remove(binary)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v\nstderr: %s", err, stderr.String())
				return
			}

			output := stdout.String()
			for _, want := range tt.contains {
				if !strings.Contains(output, want) {
					t.Errorf("output should contain %q, got: %s", want, output)
				}
			}
		})
	}
}

// TestCLI_BarGrouped tests grouped bar charts.
func TestCLI_BarGrouped(t *testing.T) {
	binary := buildBinary(t)
	defer os.Remove(binary)

	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		contains []string
	}{
		{
			name: "grouped bar chart",
			args: []string{
				"bar",
				"--series", `[{"label":"2023","data":[10,20,30]},{"label":"2024","data":[15,25,35]}]`,
				"--grouped",
				"--labels", "Q1,Q2,Q3",
			},
			wantErr:  false,
			contains: []string{"Q1", "Q2", "Q3"},
		},
		{
			name: "grouped bar chart with legend",
			args: []string{
				"bar",
				"--series", `[{"label":"2023","data":[10,20,30]},{"label":"2024","data":[15,25,35]}]`,
				"--grouped",
				"--legend",
			},
			wantErr:  false,
			contains: []string{"2023", "2024"},
		},
		{
			name: "vertical grouped bar chart",
			args: []string{
				"bar",
				"--series", `[{"label":"A","data":[10,20]},{"label":"B","data":[15,25]}]`,
				"--grouped",
				"--vertical",
			},
			wantErr: false,
		},
		{
			name: "invalid JSON series",
			args: []string{
				"bar",
				"--series", `invalid json`,
				"--grouped",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v\nstderr: %s", err, stderr.String())
				return
			}

			output := stdout.String()
			for _, want := range tt.contains {
				if !strings.Contains(output, want) {
					t.Errorf("output should contain %q, got: %s", want, output)
				}
			}
		})
	}
}

// TestCLI_BarStacked tests stacked bar charts.
func TestCLI_BarStacked(t *testing.T) {
	binary := buildBinary(t)
	defer os.Remove(binary)

	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		contains []string
	}{
		{
			name: "stacked bar chart",
			args: []string{
				"bar",
				"--series", `[{"label":"A","data":[10,20,30]},{"label":"B","data":[5,10,15]}]`,
				"--stacked",
				"--labels", "Q1,Q2,Q3",
			},
			wantErr:  false,
			contains: []string{"Q1", "Q2", "Q3"},
		},
		{
			name: "stacked bar chart with title",
			args: []string{
				"bar",
				"--series", `[{"label":"A","data":[10,20]},{"label":"B","data":[5,10]}]`,
				"--stacked",
				"--title", "Stacked Chart",
			},
			wantErr:  false,
			contains: []string{"Stacked Chart"},
		},
		{
			name: "vertical stacked bar chart",
			args: []string{
				"bar",
				"--series", `[{"label":"A","data":[10,20]},{"label":"B","data":[5,10]}]`,
				"--stacked",
				"--vertical",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v\nstderr: %s", err, stderr.String())
				return
			}

			output := stdout.String()
			for _, want := range tt.contains {
				if !strings.Contains(output, want) {
					t.Errorf("output should contain %q, got: %s", want, output)
				}
			}
		})
	}
}

// TestCLI_Spark tests the spark command.
func TestCLI_Spark(t *testing.T) {
	binary := buildBinary(t)
	defer os.Remove(binary)

	tests := []struct {
		name     string
		args     []string
		wantErr  bool
	}{
		{
			name:    "simple sparkline",
			args:    []string{"spark", "1", "5", "2", "8", "3", "7"},
			wantErr: false,
		},
		{
			name:    "ascii sparkline",
			args:    []string{"spark", "1", "2", "3", "4", "--ascii"},
			wantErr: false,
		},
		{
			name:    "sparkline with width",
			args:    []string{"spark", "1", "2", "3", "4", "5", "6", "7", "8", "--width", "4"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v\nstderr: %s", err, stderr.String())
				return
			}

			output := stdout.String()
			if output == "" {
				t.Error("expected non-empty output")
			}
		})
	}
}

// TestCLI_Pie tests the pie command.
func TestCLI_Pie(t *testing.T) {
	binary := buildBinary(t)
	defer os.Remove(binary)

	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		contains []string
	}{
		{
			name:    "simple pie chart",
			args:    []string{"pie", "30", "25", "20", "15", "10"},
			wantErr: false,
		},
		{
			name:     "pie chart with labels",
			args:     []string{"pie", "30", "20", "50", "--labels", "A,B,C"},
			wantErr:  false,
			contains: []string{"A", "B", "C"},
		},
		{
			name:     "pie chart with title",
			args:     []string{"pie", "30", "70", "--title", "Pie Chart"},
			wantErr:  false,
			contains: []string{"Pie Chart"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v\nstderr: %s", err, stderr.String())
				return
			}

			output := stdout.String()
			for _, want := range tt.contains {
				if !strings.Contains(output, want) {
					t.Errorf("output should contain %q, got: %s", want, output)
				}
			}
		})
	}
}

// TestCLI_Line tests the line command.
func TestCLI_Line(t *testing.T) {
	binary := buildBinary(t)
	defer os.Remove(binary)

	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		contains []string
	}{
		{
			name:    "simple line chart",
			args:    []string{"line", "10", "20", "15", "30", "25"},
			wantErr: false,
		},
		{
			name:     "line chart with title",
			args:     []string{"line", "10", "20", "30", "--title", "Line Chart"},
			wantErr:  false,
			contains: []string{"Line Chart"},
		},
		{
			name:    "braille line chart",
			args:    []string{"line", "10", "20", "30", "--braille"},
			wantErr: false,
		},
		{
			name:    "ascii line chart",
			args:    []string{"line", "10", "20", "30", "--ascii"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v\nstderr: %s", err, stderr.String())
				return
			}

			output := stdout.String()
			for _, want := range tt.contains {
				if !strings.Contains(output, want) {
					t.Errorf("output should contain %q, got: %s", want, output)
				}
			}
		})
	}
}

// TestCLI_Help tests help commands.
func TestCLI_Help(t *testing.T) {
	binary := buildBinary(t)
	defer os.Remove(binary)

	tests := []struct {
		name     string
		args     []string
		contains []string
	}{
		{
			name:     "root help",
			args:     []string{"--help"},
			contains: []string{"termcharts", "Available Commands"},
		},
		{
			name:     "bar help",
			args:     []string{"bar", "--help"},
			contains: []string{"bar chart", "--grouped", "--stacked"},
		},
		{
			name:     "spark help",
			args:     []string{"spark", "--help"},
			contains: []string{"sparkline"},
		},
		{
			name:     "pie help",
			args:     []string{"pie", "--help"},
			contains: []string{"pie chart"},
		},
		{
			name:     "line help",
			args:     []string{"line", "--help"},
			contains: []string{"line chart"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			var stdout bytes.Buffer
			cmd.Stdout = &stdout

			err := cmd.Run()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			output := stdout.String()
			for _, want := range tt.contains {
				if !strings.Contains(output, want) {
					t.Errorf("output should contain %q, got: %s", want, output)
				}
			}
		})
	}
}

// TestCLI_Stdin tests reading data from stdin.
func TestCLI_Stdin(t *testing.T) {
	binary := buildBinary(t)
	defer os.Remove(binary)

	tests := []struct {
		name    string
		args    []string
		stdin   string
		wantErr bool
	}{
		{
			name:    "bar from stdin",
			args:    []string{"bar"},
			stdin:   "10\n20\n30\n40",
			wantErr: false,
		},
		{
			name:    "spark from stdin",
			args:    []string{"spark"},
			stdin:   "1 2 3 4 5 6 7 8",
			wantErr: false,
		},
		{
			name:    "pie from stdin",
			args:    []string{"pie"},
			stdin:   "30,25,20,15,10",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			cmd.Stdin = strings.NewReader(tt.stdin)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v\nstderr: %s", err, stderr.String())
				return
			}

			output := stdout.String()
			if output == "" {
				t.Error("expected non-empty output")
			}
		})
	}
}

// buildBinary builds the CLI binary for testing.
func buildBinary(t *testing.T) string {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "termcharts-test-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpFile.Close()
	binaryPath := tmpFile.Name()

	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	cmd.Dir = "."
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build binary: %v\nstderr: %s", err, stderr.String())
	}

	return binaryPath
}
