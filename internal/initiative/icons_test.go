package initiative

import "testing"

func TestIconRender(t *testing.T) {
	testIcon := icon{
		icon:     "🎲",
		fallback: "[D]",
	}

	tests := []struct {
		name      string
		supported bool
		fallback  *string
		expected  string
	}{
		{
			name:      "with nerd font support",
			supported: true,
			fallback:  nil,
			expected:  "🎲",
		},
		{
			name:      "without nerd font support, no fallback",
			supported: false,
			fallback:  nil,
			expected:  "[D]",
		},
		{
			name:      "without nerd font support, with custom fallback",
			supported: false,
			fallback:  stringPtr("DICE"),
			expected:  "DICE",
		},
		{
			name:      "with nerd font support, fallback ignored",
			supported: true,
			fallback:  stringPtr("DICE"),
			expected:  "🎲",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testIcon.renderWithSupport(tt.supported, tt.fallback)

			if result != tt.expected {
				t.Errorf("renderWithSupport() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestIconJoin(t *testing.T) {
	testIcon := icon{
		icon:     "🎲",
		fallback: "[D]",
	}

	emptyIcon := icon{
		icon:     "",
		fallback: "",
	}

	tests := []struct {
		name      string
		icon      icon
		str       string
		fallback  *string
		supported bool
		expected  string
	}{
		{
			name:      "join with nerd font support",
			icon:      testIcon,
			str:       "Dice",
			fallback:  nil,
			supported: true,
			expected:  "🎲 Dice",
		},
		{
			name:      "join without nerd font support",
			icon:      testIcon,
			str:       "Dice",
			fallback:  nil,
			supported: false,
			expected:  "[D] Dice",
		},
		{
			name:      "join with empty icon returns string only",
			icon:      emptyIcon,
			str:       "Dice",
			fallback:  nil,
			supported: true,
			expected:  "Dice",
		},
		{
			name:      "join with empty icon and no support returns string only",
			icon:      emptyIcon,
			str:       "Dice",
			fallback:  nil,
			supported: false,
			expected:  "Dice",
		},
		{
			name:      "join with custom fallback",
			icon:      testIcon,
			str:       "Dice",
			fallback:  stringPtr("DICE"),
			supported: false,
			expected:  "DICE Dice",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.icon.joinWithSupport(tt.str, tt.supported, tt.fallback)

			if result != tt.expected {
				t.Errorf("joinWithSupport(%q, %v, %v) = %q, want %q", tt.str, tt.supported, tt.fallback, result, tt.expected)
			}
		})
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
