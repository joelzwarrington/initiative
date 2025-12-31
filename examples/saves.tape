# VHS tape to demonstrate saved game files
#
# 1. Default game file creation and persistence
# 2. Using custom game file paths with --game
# 3. Loading saved party data across sessions

## Output
Output saves.gif

## Settings
Set WindowBar Colorful

Hide
# Remove the existing default game save
Type "rm ~/.config/initiative/game.yaml"
Enter
Type "go build ."
Enter
Type "go install"
Enter
Require initiative
Type "clear"
Enter
Show

# Start with default game file (fake it to not delete save)
Type "initiative"
Sleep 1
Enter
Wait+Screen /Encounter/

# Create some characters to save
Tab
Sleep 2
Type "n" 
Type "Shadowheart"
Sleep 1
Enter
Sleep 2

Ctrl+c
Sleep 0.1
Type "cat ~/.config/initiative/game.yaml"
Enter
Sleep 2

Type "initiative -g game.yaml"
Sleep 1
Enter
Wait+Screen /Encounter/

Tab
Sleep 2
Type "n" 
Type "Another character!"
Sleep 1
Enter
Sleep 2

Ctrl+c
Sleep 0.1
Type "cat game.yaml"
Enter
Sleep 2

Hide
Type "git restore game.yaml"
Enter
Sleep 0.1