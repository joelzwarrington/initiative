# CLAUDE.md

## Overview

**initiative** is a CLI and TUI application for running Dungeons & Dragons encounters. It tracks initiative order, health, effects, and turns.

## Purpose

Managing Dungeon & Dragon encounters as a Dungeon Master is difficult as there is lots of information to keep track of. This CLI and TUI aims to help solve this problem.

## Technologies

- **Language**: Go
- **CLI Framework**: [Cobra](https://github.com/spf13/cobra)
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **UI Components**: [Bubbles](https://github.com/charmbracelet/bubbles), [Skeleton](https://github.com/joelzwarrington/skeleton)
- **Styling**: [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **Forms/Prompts**: [Huh?](https://github.com/charmbracelet/huh)
- **Release**: GoReleaser

## Architecture

### CLI

CLI commands are defined in `cmd/` using Cobra patterns:

- Each command is typically its own file
- Commands can launch TUI mode or perform direct CLI operations
- Use persistent flags for global options (e.g., `--game`, `--source`)

### TUI

The TUI follows the Elm architecture pattern used by Bubble Tea:

1. **Model** - Application state struct
2. **Init** - Returns initial command
3. **Update** - Handles messages, returns updated model and commands
4. **View** - Renders the UI as a string

When adding new TUI components:

- Implement the `tea.Model` interface
- Use `tea.Cmd` for async operations (which themselves return messages which be passed to Update())
- Compose models for complex views
- Use Lipgloss for consistent styling

### Component Architecture

The application has 1 parent model (program.go) which renders a skeleton which multiple pages through tabs. Each page can implement multiple models to render different views, but shouldn't be nested further than this. Each of the views may use components to ultimately render the content.

**Application Layout:**

Here's a brief example layout showcasing the Encounter page and component structure:

```
Skeleton (program.go) (Tab Container)
└── Encounter Tab
    ├── Encounter form (which itself uses a `huh.Form` from huh?)
    │   └── uses multiple huh.Form for each step in the multi-step creation form
    ├── Encounter delegate (shows an individual encounter) (which itself uses a `list.List` from bubbles)
    │   ├── initiative group list (list.List)
    │   └── health component
    ├── Hit point form
    └── among other components
```

**Component Guidelines:**

- **Root Level**: Use Skeleton for multi-tab layout and global widgets (turn timer, round counter)
  - Most widgets should only be shown when viewing a specific tab, so when navigating to/from a tab it should be removed
  - When performing an action which requires focused, tabs should be locked to prevent accidental navigation
- **Tab Level**: Each tab should be a self-contained `tea.Model` managing its domain
- **Widget Level**: Use Bubbles components (list, table, input) for specific interactions
- **Custom Components**: Create domain-specific models for D&D concepts (creature cards, initiative tracker)

**State Management:**

- Each tab manages its own state and communicates via message passing, and the cmd process
- Shared state (encounter data, settings) should be managed at the Skeleton level
- Use dependency injection to pass shared services (dice roller, persistence) to components

### Styling Standards

Use a consistent visual language using Lipgloss adaptive theming:

**Color Palette:**

```go
// Define adaptive colors for cross-terminal compatibility
var Theme = struct {
    Primary   lipgloss.AdaptiveColor
    Secondary lipgloss.AdaptiveColor
    Success   lipgloss.AdaptiveColor
    Warning   lipgloss.AdaptiveColor
    Error     lipgloss.AdaptiveColor
    Text      lipgloss.AdaptiveColor
    Muted     lipgloss.AdaptiveColor
}{
    Primary:   lipgloss.AdaptiveColor{Light: "#0969DA", Dark: "#58A6FF"},
    Secondary: lipgloss.AdaptiveColor{Light: "#656D76", Dark: "#8B949E"},
    Success:   lipgloss.AdaptiveColor{Light: "#1A7F37", Dark: "#3FB950"},
    Warning:   lipgloss.AdaptiveColor{Light: "#9A6700", Dark: "#D29922"},
    Error:     lipgloss.AdaptiveColor{Light: "#CF222E", Dark: "#F85149"},
    Text:      lipgloss.AdaptiveColor{Light: "#24292F", Dark: "#F0F6FC"},
    Muted:     lipgloss.AdaptiveColor{Light: "#656D76", Dark: "#8B949E"},
}
```

**Base Styles:**

- Use foundational styles for common patterns (cards, headers, buttons)
- Use consistent spacing (padding: 1, margin: 1)
- Establish border styles for different component types
- Define typography scale (bold for headers, regular for body, muted for meta)

**Component-Specific Guidelines:**

- **Creature Cards**: Use borders to distinguish HP status (green=healthy, yellow=bloodied, red=critical)
- **Initiative Order**: Highlight current turn with primary color, use muted for waiting turns
- **Status Effects**: Color-code by type (red=harmful, green=beneficial, gray=neutral)
- **Action Feedback**: Use success/warning/error colors for feedback and action results

**Responsive Design:**

- Components should adapt to terminal width using `lipgloss.Width()`
- Use flexible layouts that gracefully handle resizing
- Prioritize information density for smaller terminals

### Message Patterns

**Message Flow Guidelines:**

- **Parent Models**: Should respond to messages when applicable and forward messages to child components for handling
- **Child Components**: Use `tea.Cmd` to send custom messages back to their parent model
- **Message Propagation**: Models should pass messages down to children and bubble results back up
- **Command Handling**: Always check if child models return commands in Update() and use `tea.Batch` or `tea.Sequence` to combine and return them
- **Early Returns**: May early return when handling specific messages that shouldn't be passed to child components
- **Async Operations**: Use `tea.Cmd` for dice rolling, file I/O, and timers

**Key Binding System:**

- Always use the key binding system for key press handling
- Key bindings should remain consistent but be enabled/disabled based on context
- Always show help text at the bottom of the screen
- Disable help from sub-models and aggregate their key bindings into custom help text for better control and customization

## Key Domain Concepts

- **Source**: A source represents game materials which provide rules, lore and content (such as monster)
- **Encounter**: A combat session with multiple participants, at least 1 character and any number of monsters
- **Player characters**: Player characters are played by human players and will generally have more complex backstory and a full character sheet instead of a simple statblock
- **Non-playable character**: Peristent characters in a campaign that fill the supporting cast role in the story—such as a friendly shopkeeper, a quest-giving king, or a rival adventurer
- **Monster**: Technically, a monster is any creature that can be interacted with and potentially fought or killed. While "monster" often implies a terrifying beast, mechanically it includes human enemies like bandits or cultists
- **Creature**: It refers to any entity that has a stat block, can take actions, and has life (or unlife). A creature is any one of the following: player character, non-playable character, or monster.
- **Initiative**: Turn order tracking during combat
- **Health**: HP tracking for creatures during encounters
- **Effects**: Status conditions and temporary modifiers

## Code Style Guidelines

Follow standard Go conventions. This guide is based on [Google's Go Style Guide](https://google.github.io/styleguide/go/).

### Core Principles (in priority order)

1. **Clarity** - Code's purpose and rationale are evident to readers
2. **Simplicity** - Goals accomplished in the most straightforward way
3. **Concision** - High signal-to-noise ratio
4. **Maintainability** - Easy for future programmers to modify correctly
5. **Consistency** - Aligned with broader codebase patterns

### Formatting

- All code must conform to `gofmt` output
- Run `go fmt ./...` before task completion or after significant code changes
- No fixed line length; prefer refactoring over arbitrary line splits

### Development Commands

- `go build ./...` - Build all packages
- `go test ./...` - Run all tests
- `go vet ./...` - Run Go vet for potential issues
- `go mod tidy` - Clean up module dependencies
- `go fmt ./...` - Format all code

**Note:** Avoid `go run .` in non-TTY environments as the TUI requires a terminal.

### Naming Conventions

- Use `MixedCaps` or `mixedCaps`, never underscores (except test functions, generated code)
- Acronyms maintain consistent case: `URL`, `HTTP`, `ID`, `url`, `http`, `id` (not `Url`, `Http`, `Id`)
- Package names: lowercase, single-word, no underscores; avoid generic names like `util` or `common`
- Keep names short but descriptive; prefer `i` to `index` in short scopes
- Avoid redundancy: `Parse()` in package `yamlconfig`, not `ParseYAMLConfig()`
- Getters omit `Get` prefix: use `Name()` not `GetName()`
- Error variables start with `Err`: `ErrNotFound` not `NotFoundError`
- Interfaces with one method use `-er` suffix: `Reader`, `Writer`
- Receiver names: short (1-2 letters), consistent across all methods of a type

### Error Handling

- Always check errors; never ignore with blank identifier (unless intentional)
- Handle errors early with early returns; keep normal code flow unindented
- Wrap errors with context: `fmt.Errorf("operation failed: %w", err)` (place `%w` at end)
- Use `errors.Is()` and `errors.As()` for inspection (not direct comparison)
- Error strings: lowercase, no punctuation (they often appear in larger context)
- Return errors to callers; only log at the top level
- Don't panic except for truly unrecoverable situations
- Avoid in-band errors (returning -1 or empty string); return additional validity value instead

### Functions & Methods

- Keep functions focused on a single responsibility
- Prefer early returns to reduce nesting
- Context should be first parameter: `func DoThing(ctx context.Context, ...)`
- For many parameters, use option structs or variadic functional options

### Code Organization

- Group imports: standard library, third-party, local packages
- Rename imports only when necessary (collisions, uninformative names)
- Define interfaces in consuming packages, not producing packages
- Return concrete types from constructors
- Group related code by concept, not "one type per file"

### Comments & Documentation

- Every exported function/type needs a doc comment starting with its name
- Doc comments: complete sentences, capitalized, punctuated
- Comments explain _why_, not _what_ (code shows _what_)
- Package comments go in `doc.go` or at top of main file (no blank line before `package`)
- Wrap comments at ~80 characters

## Testing

- Tests live alongside source files (`*_test.go`)
- Use the standard `testing` package; avoid assertion libraries
- Use table-driven tests where appropriate
- Test both success and error paths
- Use "got before want" format: `YourFunc(%v) = %v, want %v`
- Prefer `t.Error` over `t.Fatal` to allow multiple failures per run
- Use `cmp.Diff` for complex comparisons
- Never call `t.Fatal()` from goroutines; use `t.Error()` with return
- Keep test logic in the `Test` function, not assertion helpers
- Prefer fewer high-quality tests over many lower-quality tests
- Focus on testing important functionality; avoid tests for trivial code
- Prefer updating existing tests to avoid test bloat
