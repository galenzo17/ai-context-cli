# Demo 1: Main Menu Navigation
# This demo shows the main menu interface and navigation

Output demos/01-main-menu.gif
Set Theme "Catppuccin Mocha"
Set FontSize 16
Set Width 1200
Set Height 800

# Type the command to run the CLI
Type "go run cmd/ai-context-cli/main.go"
Sleep 1s
Enter

# Wait for the app to load
Sleep 2s

# Show the main menu
Sleep 2s

# Add annotation for the banner
Type@100ms ""
Sleep 1s

# Navigate through menu items
Type "j"
Sleep 500ms
Type "j" 
Sleep 500ms
Type "k"
Sleep 500ms

# Show help for current item
Type "?"
Sleep 2s

# Close help
Type "esc"
Sleep 1s

# Navigate to different options
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms

# Show help for context preview
Type "?"
Sleep 2s
Type "esc"
Sleep 1s

# Exit application
Type "q"
Sleep 1s