package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Health represents a creature's health status
type Health struct {
	maxHealth     int
	currentHealth int
	progress      progress.Model
	width         int
}

// NewHealth creates a new health component
func NewHealth(maxHealth int, width int) Health {
	prog := progress.New(progress.WithDefaultGradient())
	prog.Width = width

	return Health{
		maxHealth:     maxHealth,
		currentHealth: maxHealth,
		progress:      prog,
		width:         width,
	}
}

// SetHealth updates the current and maximum health values
func (h *Health) SetHealth(current int, maximum int) {
	h.maxHealth = maximum
	if current < 0 {
		current = 0
	}
	if current > h.maxHealth {
		current = h.maxHealth
	}
	h.currentHealth = current
}

// getPercent returns the health as a percentage (0.0 to 1.0)
func (h *Health) getPercent() float64 {
	if h.maxHealth == 0 {
		return 0.0
	}
	return float64(h.currentHealth) / float64(h.maxHealth)
}

// GetStatus returns a string describing the health status
func (h *Health) status() string {
	percent := h.getPercent()
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
func (h *Health) SetWidth(width int) {
	h.width = width
	h.progress.Width = width
}

// Init implements tea.Model interface (not needed for component usage)
func (h Health) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model interface (not needed for component usage)
func (h Health) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return h, nil
}

// View renders the health component as a string
func (h Health) View() string {
	percent := h.getPercent()

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
	healthText := fmt.Sprintf("%d/%d (%s)", h.currentHealth, h.maxHealth, h.status())

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
