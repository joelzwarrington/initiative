package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
)

// HealthBar represents a creature's health status
type HealthBar struct {
	progress progress.Model
	width    int
}

// NewHealthBar creates a new health bar component
func NewHealthBar(width int) HealthBar {
	prog := progress.New(progress.WithDefaultGradient())
	prog.Width = width

	return HealthBar{
		progress: prog,
		width:    width,
	}
}

// getPercent returns the health as a percentage (0.0 to 1.0)
func getPercent(current, maximum int) float64 {
	if maximum == 0 {
		return 0.0
	}
	if current < 0 {
		current = 0
	}
	if current > maximum {
		current = maximum
	}
	return float64(current) / float64(maximum)
}

// getStatus returns a string describing the health status
func getStatus(current, maximum int) string {
	percent := getPercent(current, maximum)
	switch {
	case percent == 0:
		return "Dead"
	case percent < 0.25:
		return "Critical"
	case percent < 0.5:
		return "Wounded"
	case percent < 0.75:
		return "Injured"
	case percent < 1.0:
		return "Hurt"
	default:
		return "Healthy"
	}
}

// SetWidth updates the width of the progress bar
func (h *HealthBar) SetWidth(width int) {
	h.width = width
	h.progress.Width = width
}

// View renders the health component as a string with given current and maximum health
func (h *HealthBar) View(current, maximum int) string {
	// Clamp current health for display
	if current < 0 {
		current = 0
	}
	if current > maximum && maximum > 0 {
		current = maximum
	}

	percent := getPercent(current, maximum)

	// Create single color based on health level
	var color string
	switch {
	case percent >= 0.75:
		color = "#0f0" // Green for healthy
	case percent >= 0.5:
		color = "#ff0" // Yellow for injured
	case percent >= 0.25:
		color = "#f80" // Orange for wounded
	default:
		color = "#f00" // Red for critical/dead
	}

	// Create health text
	healthText := fmt.Sprintf("%d/%d (%s)", current, maximum, getStatus(current, maximum))

	// Style the health text based on status
	var textStyle lipgloss.Style
	switch {
	case percent <= 0:
		textStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#666")).Italic(true)
	case percent < 0.25:
		textStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f00")).Bold(true)
	case percent < 0.5:
		textStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f80"))
	case percent < 0.75:
		textStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0"))
	default:
		textStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#0f0"))
	}

	styledHealthText := textStyle.Render(healthText)

	// Use a constant progress bar width (e.g., 20 characters)
	progressBarWidth := 20

	// Create progress bar with single color (no percentage display)
	prog := progress.New(progress.WithSolidFill(color), progress.WithoutPercentage())
	prog.Width = progressBarWidth

	progressBar := prog.ViewAs(percent)

	// Calculate remaining space and right-align the text
	totalContentWidth := progressBarWidth + 1 + len(healthText) // bar + space + text
	if totalContentWidth < h.width {
		padding := h.width - totalContentWidth
		return fmt.Sprintf("%s %s%s", progressBar, lipgloss.NewStyle().Width(padding).Render(""), styledHealthText)
	}

	// If no room for padding, just use space
	return fmt.Sprintf("%s %s", progressBar, styledHealthText)
}
