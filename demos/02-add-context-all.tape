# Demo 2: Add Context to All Files
# This demo shows the full project scanning functionality

Output demos/02-add-context-all.gif
Set Theme "Catppuccin Mocha"
Set FontSize 16
Set Width 1200
Set Height 800

# Start the application
Type "go run cmd/ai-context-cli/main.go"
Sleep 1s
Enter
Sleep 2s

# Select "Add Context to All Files" (first option)
Enter
Sleep 1s

# Watch the scanning process
Sleep 3s

# Let the scanning complete
Sleep 5s

# Navigate through the results
Type "r"
Sleep 2s

# Exit
Type "q"
Sleep 1s