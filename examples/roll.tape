# VHS tape to demonstrate rolling dice
#
# 1. Rolling multiple dice (8d6)
# 2. Rolling with modifiers (1d20+4)

## Output
Output roll.gif

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

# Roll 8d6
Type "initiative roll 8d6"
Sleep 1
Enter
Sleep 2

# Roll 1d20+4
Type "initiative roll 1d20+4"
Sleep 1
Enter
Sleep 2
