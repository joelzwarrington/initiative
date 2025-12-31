# VHS tape to demonstrate using custom sources
#
# 1. Using default SRD monsters
# 2. Loading custom monster sources
# 3. Mixing multiple sources in encounters

## Output
Output sources.gif

## Settings
Set WindowBar Colorful
Set Height 800

Hide
Type "go build ."
Enter
Type "go install"
Enter
Require initiative
Type "clear"
Enter
Show

Type "initiative -g game.yaml"
Sleep 1
Enter
Sleep 1

# Add new encounter to showcase monsters
Type "n"
Sleep 0.5
Type "Using source reference document"
Sleep 0.5
Enter@300ms 2

# View some monsters

Down@200ms 30

# Exit to demonstrate custom source
Ctrl+c
Sleep 2

# Start with custom source file
Type "initiative -g game.yaml --source source.yaml"
Sleep 3
Enter
Sleep 1

# Add new encounter to showcase monsters
Type "n"
Sleep 0.5
Type "Using custom source"
Sleep 0.5
Enter 2

Sleep 3

# Exit to demonstrate custom source and srd
Ctrl+c
Sleep 2

Type "initiative -g game.yaml --source srd --source source.yaml"
Sleep 3
Enter
Sleep 1

# Add new encounter to showcase monsters
Type "n"
Sleep 0.5
Type "Using source reference document & custom source"
Sleep 0.5
Enter 2

# Show both SRD and custom monsters
Down@100ms 30
Sleep 0.5
Type "/Void"
Sleep 2
