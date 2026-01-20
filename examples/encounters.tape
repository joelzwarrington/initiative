# VHS tape to demonstrate managing an encounter
#
# 1. Creating encounters with pre-existing characters
# 2. Adding monsters and setting initiative

## Output
Output encounters.gif

## Settings
Set LoopOffset 95%
Set Shell fish
Set WindowBar Colorful
Set Padding 10
Set Height 800
Set Width 1200
Set TypingSpeed 100ms

Hide
Type "go run ../. --game game.yaml"
Enter
Wait+Screen /Encounter/
Show

Sleep 2s
Type "n"

# Add encounter

Sleep 1s
Type "Ambush at Emerald Grove"
Sleep 2s
Enter

# Set character initiatives

Type "19"
Sleep 1s
Down
Sleep 1s

Type "6"
Enter
Sleep 1s
Down
Sleep 1s

Type "11"
Sleep 1s
Down
Sleep 1s

Type "d"
Sleep 1s

Type "15"
Sleep 1s

# Add monsters

## Add goblins
Sleep 1s
Type "a"
Sleep 1.5s
Type "/Goblin"
Sleep 1s
Enter
Type "6"
Enter 2
Sleep 1s
Type "8"
Enter

## Add hobgoblins
Sleep 1s
Type "a"
Sleep 1.5s
Type "/Hobgoblin"
Sleep 1s
Enter
Type "2"
Enter 2
Sleep 1s
Type "13"

## Add worg
Sleep 1s
Type "a"
Sleep 1.5s
Type "/Worg"
Sleep 1s
Enter 3
Sleep 1s
Type "4"
Enter 2

Sleep 4s
