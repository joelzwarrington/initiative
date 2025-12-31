# VHS tape to demonstrate tracking health
#
# 1. Applying damage with 'd' key
# 2. Healing creatures with 'h' key
# 3. Health status visualization

## Output
Output health.gif

## Settings
Set WindowBar Colorful
Set Height 800

Hide
Type "go run ../. --game game.yaml"
Enter
Wait+Screen /Encounter/

Type "n"
Type "Ambush at Emerald Grove"
Enter

Sleep 0.25
Type "/Karlach"
Sleep 1
Enter
Sleep 0.25
Space
Sleep 0.25
Enter

Hide

Sleep 0.25
Type "/Goblin"
Sleep 0.25
Enter
Sleep 0.25
Space
Sleep 0.25
Type "/"
Backspace@0ms 6

Sleep 0.25
Type "Bugbear"
Sleep 0.25
Enter
Sleep 0.25
Space
Sleep 0.25
Type "/"
Backspace@0ms 7 

Sleep 0.25
Type "Worg"
Sleep 0.25
Enter
Sleep 0.25
Space
Sleep 0.25
Enter

Sleep 0.25
Type "4"
Sleep 0.25
Enter

Sleep 0.25
Type "10"
Sleep 0.25
Enter

Sleep 0.25
Type "17"
Sleep 0.25
Enter

Sleep 0.25
Type "11"
Sleep 0.25
Enter

Sleep 0.25
Enter
Sleep 0.25
Type "3"
Sleep 0.25
Enter
Sleep 0.25
Type "19"
Sleep 0.25
Enter

Sleep 0.25
Enter
Sleep 0.25
Type "7"
Sleep 0.25
Enter
Sleep 0.25
Type "14"
Sleep 0.25
Enter

Sleep 0.25
Enter
Sleep 0.25
Type "2"
Sleep 0.25
Enter
Sleep 0.25
Type "3"
Sleep 0.25
Enter
Sleep 0.25
Right
Sleep 0.25
Up
Sleep 0.25

Show

Sleep 2
Type "d"
Sleep 1
Type "9"
Sleep 1
Enter 2
Sleep 2

Down@250ms 2
Type "d"
Sleep 1
Type "7"
Sleep 1
Enter 2
Sleep 2

Right
Up@250ms 2
Type "h"
Sleep 1
Type "4"
Sleep 1
Enter 2
Sleep 2