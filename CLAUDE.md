# Claude Context

## Project Overview
This is a Go CLI application built with Bubble Tea for managing tabletop RPG initiative tracking.

## Architecture
```
app.go (router) -> view models (business logic + keybindings) -> data (pure data)
```

### Key Principles
1. **Slim app.go**: Only handles view routing and delegates everything to view models
2. **Self-contained views**: Each view model handles its own keybindings and Update logic
3. **Navigation via commands**: Views return navigation commands (NavigateToGameList(), etc.) instead of direct state changes
4. **Private by default**: Only expose methods that need to be called from outside the package

## Current Structure
- `internal/data/` - Pure data structures (Game struct)
- `internal/ui/app.go` - Main app router, handles navigation messages
- `internal/ui/views/game.go` - View models: GameListModel, GameNewFormModel, GamePageModel

## Testing & Building
- Run tests: `go test ./...` 
- Build: `go build`
- All tests should pass before making changes

## Recent Refactoring
- Removed `internal/models` (moved to `internal/data`)
- Removed centralized keybindings (moved to individual view models)
- Views use standard `Update()` methods with navigation commands
- Made internal methods private (lowercase names)