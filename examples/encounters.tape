# VHS tape to demonstrate managing an encounter
#
# 1. Creating encounters with pre-existing characters
# 2. Adding monsters and setting initiative

## Output
Output encounters.gif

## Settings
Set WindowBar Colorful
Set Height 800

Hide
Type "go run ../. --game game.yaml"
Enter
Wait+Screen /Encounter/
Show

Sleep 2
Type "n"

# Add encounter

Sleep 0.5
Type "Ambush at Emerald Grove"
Sleep 0.5
Enter

# Add characters

Sleep 0.5
Type "/Karlach"
Sleep 0.5
Enter
Sleep 0.75
Space
Sleep 0.75
Enter

# Add monsters

Sleep 0.5
Type "/Goblin"
Sleep 0.5
Enter
Sleep 0.75
Space
Sleep 0.75
Type "/"
Backspace@0ms 6

Sleep 0.5
Type "Bugbear"
Sleep 0.5
Enter
Sleep 0.75
Space
Sleep 0.75
Type "/"
Backspace@0ms 7 

Sleep 0.5
Type "Worg"
Sleep 0.5
Enter
Sleep 0.75
Space
Sleep 0.75
Enter

# Setup initiative

Sleep 1
Type "12"
Sleep 0.5
Enter
Sleep 0.5
Type "16"
Sleep 0.5
Enter
Sleep 0.5
Type "7"
Sleep 0.5
Enter
Sleep 0.5
Type "20"
Sleep 0.5
Enter
Sleep 0.5

# Add initiative / quantity for monsters

Sleep 1
Enter
Sleep 0.5
Type "2"
Sleep 0.5
Enter
Sleep 0.5
Type "13"
Sleep 0.5
Enter
Sleep 0.5

Enter
Sleep 0.5
Type "6"
Sleep 0.5
Enter
Sleep 0.5
Type "4"
Sleep 0.5
Enter

Enter
Sleep 0.5
Type "1"
Sleep 0.5
Enter
Sleep 0.5
Type "17"
Sleep 0.5
Enter

Sleep 3
