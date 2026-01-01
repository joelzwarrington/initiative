package dnd

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// SetSeed sets the random seed for dice rolls
// Useful for testing or deterministic results
func SetSeed(seed int64) {
	rng = rand.New(rand.NewSource(seed))
}

// DiceRoll represents the result of rolling dice
type DiceRoll struct {
	Total   int      // Total value after all modifiers
	Details []string // Detailed breakdown of each roll
	Formula string   // Original dice notation
}

// String returns a formatted string representation of the dice roll
func (dr DiceRoll) String() string {
	return fmt.Sprintf("%s = %s = %d", dr.Formula, strings.Join(dr.Details, " + "), dr.Total)
}

// Roll parses dice notation and returns the result
// Supports formats like: "d6", "2d6", "1d20+5", "3d8-2"
// Follows grammar: [count]d<sides>[+/-modifier]
func Roll(notation string) (*DiceRoll, error) {
	// Remove all spaces for easier parsing
	notation = strings.ReplaceAll(notation, " ", "")

	if notation == "" {
		return nil, fmt.Errorf("empty dice notation")
	}

	// Parse the dice notation into components
	components, err := parseDiceNotation(notation)
	if err != nil {
		return nil, err
	}

	var total int
	var details []string

	for _, component := range components {
		switch component.Type {
		case "dice":
			diceResult := rollDice(component.Count, component.Sides)
			total += diceResult.sum
			details = append(details, diceResult.detail)
		case "modifier":
			total += component.Value
			if component.Value >= 0 {
				details = append(details, fmt.Sprintf("%d", component.Value))
			} else {
				details = append(details, fmt.Sprintf("%d", component.Value))
			}
		}
	}

	return &DiceRoll{
		Total:   total,
		Details: details,
		Formula: notation,
	}, nil
}

// Component represents a part of the dice notation
type Component struct {
	Type  string // "dice" or "modifier"
	Count int    // number of dice (for dice type)
	Sides int    // sides of dice (for dice type)
	Value int    // modifier value (for modifier type)
}

// parseDiceNotation breaks down dice notation into components
// Follows strict grammar: [count]d<sides>[+/-modifier]
func parseDiceNotation(notation string) ([]Component, error) {
	var components []Component

	// Use regex to match the exact grammar: optional count, 'd', sides, optional modifier
	pattern := regexp.MustCompile(`^(\d*)d(\d+)([+-]\d+)?$`)
	matches := pattern.FindStringSubmatch(notation)

	if matches == nil {
		return nil, fmt.Errorf("invalid dice notation: %s", notation)
	}

	countStr := matches[1]
	sidesStr := matches[2]
	modifierStr := matches[3]

	// Parse dice count (default to 1 if not specified)
	count := 1
	if countStr != "" {
		var err error
		count, err = strconv.Atoi(countStr)
		if err != nil || count <= 0 {
			return nil, fmt.Errorf("invalid dice count: %s", countStr)
		}
	}

	// Parse dice sides
	sides, err := strconv.Atoi(sidesStr)
	if err != nil || sides <= 0 {
		return nil, fmt.Errorf("invalid dice sides: %s", sidesStr)
	}

	// Add dice component
	components = append(components, Component{
		Type:  "dice",
		Count: count,
		Sides: sides,
	})

	// Parse modifier if present
	if modifierStr != "" {
		value, err := strconv.Atoi(modifierStr)
		if err != nil {
			return nil, fmt.Errorf("invalid modifier: %s", modifierStr)
		}

		// Reject zero modifiers (not allowed in grammar)
		if value == 0 {
			return nil, fmt.Errorf("zero modifier not allowed: %s", modifierStr)
		}

		components = append(components, Component{
			Type:  "modifier",
			Value: value,
		})
	}

	return components, nil
}

// diceResult holds the result of rolling a set of dice
type diceResult struct {
	sum    int
	detail string
}

// rollDice rolls the specified number of dice with given sides
func rollDice(count, sides int) diceResult {
	var rolls []int
	var sum int

	for range count {
		roll := rng.Intn(sides) + 1
		rolls = append(rolls, roll)
		sum += roll
	}

	rollsStr := make([]string, len(rolls))
	for i, roll := range rolls {
		rollsStr[i] = strconv.Itoa(roll)
	}

	var detail string
	if count == 1 {
		detail = strconv.Itoa(sum)
	} else {
		detail = fmt.Sprintf("[%s]", strings.Join(rollsStr, ","))
	}

	return diceResult{
		sum:    sum,
		detail: detail,
	}
}
