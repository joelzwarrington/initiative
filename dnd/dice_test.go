package dnd

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

func TestRoll_BasicDice(t *testing.T) {
	// Set deterministic seed for testing
	rng = rand.New(rand.NewSource(1))

	tests := []struct {
		notation string
		wantErr  bool
	}{
		{"d6", false},
		{"1d6", false},
		{"2d6", false},
		{"1d20", false},
		{"3d8", false},
		{"", true},
		{"invalid", true},
		{"d0", true},
		{"0d6", true},
		{"-1d6", true},
	}

	for _, tt := range tests {
		t.Run(tt.notation, func(t *testing.T) {
			result, err := Roll(tt.notation)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Roll(%q) error = nil, want error", tt.notation)
				}
				return
			}

			if err != nil {
				t.Errorf("Roll(%q) error = %v, want nil", tt.notation, err)
				return
			}

			if result == nil {
				t.Errorf("Roll(%q) returned nil result", tt.notation)
				return
			}

			if result.Total <= 0 {
				t.Errorf("Roll(%q) Total = %d, want > 0", tt.notation, result.Total)
			}

			if result.Formula != tt.notation {
				t.Errorf("Roll(%q) Formula = %q, want %q", tt.notation, result.Formula, tt.notation)
			}

			if len(result.Details) == 0 {
				t.Errorf("Roll(%q) Details is empty", tt.notation)
			}
		})
	}
}

func TestRoll_WithModifiers(t *testing.T) {
	tests := []struct {
		notation string
		wantErr  bool
	}{
		{"1d6+1", false},
		{"1d6-1", false},
		{"1d20+5", false},
		{"2d6+6", false},
		{"3d8-2", false},
		{"1d6+0", true},     // zero modifiers not allowed
		{"1d6-0", true},     // zero modifiers not allowed
		{"1d4+1d6+2", true}, // multiple dice expressions not allowed
		{"d6+", true},
		{"1d6+-", true},
	}

	for _, tt := range tests {
		t.Run(tt.notation, func(t *testing.T) {
			result, err := Roll(tt.notation)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Roll(%q) error = nil, want error", tt.notation)
				}
				return
			}

			if err != nil {
				t.Errorf("Roll(%q) error = %v, want nil", tt.notation, err)
				return
			}

			if result.Formula != tt.notation {
				t.Errorf("Roll(%q) Formula = %q, want %q", tt.notation, result.Formula, tt.notation)
			}
		})
	}
}

func TestRoll_DeterministicResults(t *testing.T) {
	// Test with fixed seed for predictable results
	rng = rand.New(rand.NewSource(42))

	result, err := Roll("2d6+3")
	if err != nil {
		t.Fatalf("Roll() error = %v", err)
	}

	// With seed 42, this should give us specific results
	if result.Total == 0 {
		t.Errorf("Roll() Total should not be 0")
	}

	if len(result.Details) != 2 {
		t.Errorf("Roll() Details length = %d, want 2", len(result.Details))
	}

	// Verify the dice part shows multiple rolls
	diceDetail := result.Details[0]
	if !strings.Contains(diceDetail, "[") || !strings.Contains(diceDetail, "]") {
		t.Errorf("Roll() dice detail %q should show multiple rolls in brackets", diceDetail)
	}

	// Verify modifier is present
	modifierDetail := result.Details[1]
	if modifierDetail != "3" {
		t.Errorf("Roll() modifier detail = %q, want '3'", modifierDetail)
	}
}

func TestRoll_SingleDieFormat(t *testing.T) {
	rng = rand.New(rand.NewSource(1))

	result, err := Roll("1d6")
	if err != nil {
		t.Fatalf("Roll() error = %v", err)
	}

	// Single die should not use brackets
	detail := result.Details[0]
	if strings.Contains(detail, "[") || strings.Contains(detail, "]") {
		t.Errorf("Single die detail %q should not use brackets", detail)
	}

	// Should be just a number
	if _, err := strconv.Atoi(detail); err != nil {
		t.Errorf("Single die detail %q should be a valid number", detail)
	}
}

func TestRoll_MultipleDiceFormat(t *testing.T) {
	rng = rand.New(rand.NewSource(1))

	result, err := Roll("3d6")
	if err != nil {
		t.Fatalf("Roll() error = %v", err)
	}

	// Multiple dice should use brackets and commas
	detail := result.Details[0]
	if !strings.Contains(detail, "[") || !strings.Contains(detail, "]") {
		t.Errorf("Multiple dice detail %q should use brackets", detail)
	}

	if !strings.Contains(detail, ",") {
		t.Errorf("Multiple dice detail %q should contain commas", detail)
	}

	// Should have format like [1,2,3]
	if !strings.HasPrefix(detail, "[") || !strings.HasSuffix(detail, "]") {
		t.Errorf("Multiple dice detail %q should be wrapped in brackets", detail)
	}
}

func TestRoll_NegativeModifiers(t *testing.T) {
	result, err := Roll("1d6-2")
	if err != nil {
		t.Fatalf("Roll() error = %v", err)
	}

	// Should have two details: dice roll and negative modifier
	if len(result.Details) != 2 {
		t.Errorf("Roll() Details length = %d, want 2", len(result.Details))
	}

	modifierDetail := result.Details[1]
	if modifierDetail != "-2" {
		t.Errorf("Roll() negative modifier detail = %q, want '-2'", modifierDetail)
	}
}

func TestRoll_InvalidComplexNotation(t *testing.T) {
	// Complex notation should now be rejected
	_, err := Roll("2d6+1d4+3")
	if err == nil {
		t.Errorf("Roll(\"2d6+1d4+3\") should fail with strict grammar")
	}

	// Test another complex case
	_, err = Roll("1d4+1d6")
	if err == nil {
		t.Errorf("Roll(\"1d4+1d6\") should fail with strict grammar")
	}
}

func TestDiceRoll_String(t *testing.T) {
	roll := DiceRoll{
		Total:   10,
		Details: []string{"[3,2]", "5"},
		Formula: "2d6+5",
	}

	expected := "2d6+5 = [3,2] + 5 = 10"
	if got := roll.String(); got != expected {
		t.Errorf("DiceRoll.String() = %q, want %q", got, expected)
	}
}

func TestParseDiceNotation(t *testing.T) {
	tests := []struct {
		notation string
		wantErr  bool
		wantLen  int
	}{
		{"1d6", false, 1},
		{"2d6+3", false, 2},
		{"1d20-1", false, 2},
		{"d6", false, 1},
		{"1d6+0", true, 0},     // zero modifier not allowed
		{"2d6+1d4+5", true, 0}, // multiple dice not allowed
		{"invalid", true, 0},
		{"", true, 0},
		{"d6+", true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.notation, func(t *testing.T) {
			components, err := parseDiceNotation(tt.notation)

			if tt.wantErr {
				if err == nil {
					t.Errorf("parseDiceNotation(%q) error = nil, want error", tt.notation)
				}
				return
			}

			if err != nil {
				t.Errorf("parseDiceNotation(%q) error = %v, want nil", tt.notation, err)
				return
			}

			if len(components) != tt.wantLen {
				t.Errorf("parseDiceNotation(%q) components length = %d, want %d", tt.notation, len(components), tt.wantLen)
			}
		})
	}
}

func TestRoll_GrammarCompliance(t *testing.T) {
	tests := []struct {
		notation string
		valid    bool
		reason   string
	}{
		// Valid grammar
		{"d6", true, "basic die"},
		{"1d6", true, "explicit count"},
		{"12d20", true, "multi-digit count and sides"},
		{"1d6+1", true, "positive modifier"},
		{"2d8-3", true, "negative modifier"},

		// Invalid grammar
		{"0d6", false, "zero count"},
		{"1d0", false, "zero sides"},
		{"1d6+0", false, "zero modifier"},
		{"1d6-0", false, "zero modifier (negative)"},
		{"2d6+1d4", false, "multiple dice"},
		{"1d6+2+3", false, "multiple modifiers"},
		{"d6d8", false, "malformed"},
		{"", false, "empty"},
		{"invalid", false, "not dice notation"},
	}

	for _, tt := range tests {
		t.Run(tt.notation, func(t *testing.T) {
			_, err := Roll(tt.notation)
			hasError := err != nil

			if tt.valid && hasError {
				t.Errorf("Roll(%q) should be valid (%s) but got error: %v", tt.notation, tt.reason, err)
			} else if !tt.valid && !hasError {
				t.Errorf("Roll(%q) should be invalid (%s) but was accepted", tt.notation, tt.reason)
			}
		})
	}
}

func TestSetSeed(t *testing.T) {
	// Test that seeding produces deterministic results
	SetSeed(42)
	result1, _ := Roll("2d6")

	SetSeed(42)
	result2, _ := Roll("2d6")

	if result1.Total != result2.Total {
		t.Errorf("Same seed should produce same results: %d != %d", result1.Total, result2.Total)
	}

	if result1.String() != result2.String() {
		t.Errorf("Same seed should produce identical output: %q != %q", result1.String(), result2.String())
	}

	// Test different seeds produce different results (most of the time)
	SetSeed(1)
	result3, _ := Roll("2d6")

	SetSeed(2)
	result4, _ := Roll("2d6")

	// This isn't guaranteed but very likely
	if result3.String() == result4.String() {
		t.Logf("Different seeds produced same result (rare but possible): %q", result3.String())
	}
}

func TestRollDice(t *testing.T) {
	// Test with deterministic seed
	rng = rand.New(rand.NewSource(123))

	result := rollDice(2, 6)

	if result.sum <= 0 {
		t.Errorf("rollDice(2, 6) sum = %d, want > 0", result.sum)
	}

	if result.sum < 2 || result.sum > 12 {
		t.Errorf("rollDice(2, 6) sum = %d, want between 2 and 12", result.sum)
	}

	if !strings.Contains(result.detail, "[") || !strings.Contains(result.detail, "]") {
		t.Errorf("rollDice(2, 6) detail = %q, should contain brackets", result.detail)
	}

	// Test single die
	singleResult := rollDice(1, 6)
	if strings.Contains(singleResult.detail, "[") {
		t.Errorf("rollDice(1, 6) detail = %q, should not contain brackets", singleResult.detail)
	}
}
