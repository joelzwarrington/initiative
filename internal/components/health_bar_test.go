package components

import (
	"strings"
	"testing"
)

func TestHealthBar_View(t *testing.T) {
	width := 50

	tests := []struct {
		name     string
		current  int
		maximum  int
		width    int
		expected string
	}{
		{
			name:     "healthy state",
			current:  100,
			maximum:  100,
			expected: "100/100 (Healthy)",
		},
		{
			name:     "hurt state",
			current:  80,
			maximum:  100,
			expected: "80/100 (Hurt)",
		},
		{
			name:     "injured state",
			current:  60,
			maximum:  100,
			expected: "60/100 (Injured)",
		},
		{
			name:     "wounded state",
			current:  30,
			maximum:  100,
			expected: "30/100 (Wounded)",
		},
		{
			name:     "critical state",
			current:  10,
			maximum:  100,
			expected: "10/100 (Critical)",
		},
		{
			name:     "dead state",
			current:  0,
			maximum:  100,
			expected: "0/100 (Dead)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hb := NewHealthBar(width)
			result := hb.View(tt.current, tt.maximum)

			// Check that the result contains expected health text
			if !strings.Contains(result, tt.expected) {
				t.Errorf("View() result should contain %q, got %q", tt.expected, result)
			}

			// Basic format check - should have progress bar and health text
			if len(result) == 0 {
				t.Error("View() should not return empty string")
			}
		})
	}
}

func TestHealthBar_View_EdgeCases(t *testing.T) {
	width := 50
	tests := []struct {
		name     string
		current  int
		maximum  int
		expected string
	}{
		{
			name:     "negative current health",
			current:  -10,
			maximum:  100,
			expected: "0/100 (Dead)", // should clamp to 0
		},
		{
			name:     "current health exceeds maximum",
			current:  150,
			maximum:  100,
			expected: "100/100 (Healthy)", // should clamp to maximum
		},
		{
			name:     "zero maximum health",
			current:  50,
			maximum:  0,
			expected: "50/0 (Dead)", // should handle division by zero
		},
		{
			name:     "both zero",
			current:  0,
			maximum:  0,
			expected: "0/0 (Dead)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hb := NewHealthBar(width)
			result := hb.View(tt.current, tt.maximum)

			if !strings.Contains(result, tt.expected) {
				t.Errorf("View() result should contain %q, got %q", tt.expected, result)
			}
		})
	}
}
