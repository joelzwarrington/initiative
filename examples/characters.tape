# Showcase management of characters in the party
#
# 1. Adding characters
# 2. Editing a character
# 3. Removing a character

## Output
Output characters.gif

## Settings
Set WindowBar Colorful

Hide
Type "rm game.yaml"
Enter
Type "go run ../. --game game.yaml"
Enter
Wait+Screen /Encounter/
Tab
Show

Sleep 2
Type "n"

Sleep 0.5
Type "Astarion"
Sleep 0.5
Enter

Sleep 2
Type "n"

Sleep 0.5
Type "Gale"
Sleep 0.5
Enter

Sleep 2
Type "n"

Sleep 0.5
Type "Halsin"
Sleep 0.5
Enter

Sleep 1
Up
Sleep 1
Type "e"

Sleep 0.5
Type " Dekarios"
Sleep 0.5
Enter

Sleep 1
Down
Sleep 1
Type "d"

Sleep 2

Hide
git restore game.yaml
