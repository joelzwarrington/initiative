package initiative

import "fmt"

type icon struct {
	icon     string
	fallback string
}

type iconMap struct {
	encounterTab icon
	partyTab     icon
	sourcesTab   icon
}

var icons = iconMap{
	encounterTab: icon{icon: "󰓥", fallback: ""},
	partyTab:     icon{icon: "", fallback: ""},
	sourcesTab:   icon{icon: "󰧮", fallback: ""},
}

func (i icon) Render(fallback *string) string {
	return i.renderWithSupport(hasNerdFontSupport(), fallback)
}

func (i icon) renderWithSupport(supported bool, fallback *string) string {
	if supported {
		return i.icon
	} else if fallback != nil {
		return *fallback
	}

	return i.fallback
}

func (i icon) Join(str string, fallback *string) string {
	return i.joinWithSupport(str, hasNerdFontSupport(), fallback)
}

func (i icon) joinWithSupport(str string, supported bool, fallback *string) string {
	iconStr := i.renderWithSupport(supported, fallback)

	if iconStr == "" {
		return str
	}

	return fmt.Sprintf("%s %s", iconStr, str)
}

// hasNerdFontSupport attempts to detect NerdFont support
func hasNerdFontSupport() bool {
	// This is hardcoded for now, it will be implemented in the future.
	return true
}
