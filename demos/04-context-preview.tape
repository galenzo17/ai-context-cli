# Demo 4: Context Preview and Editing
# This demo shows the context preview functionality after scanning

Output demos/04-context-preview.gif
Set Theme "Catppuccin Mocha"
Set FontSize 16
Set Width 1200
Set Height 800

# Start the application
Type "go run cmd/ai-context-cli/main.go"
Sleep 1s
Enter
Sleep 2s

# First scan all files to have context
Enter
Sleep 3s

# Wait for scan to complete and go back to menu
Sleep 5s
Type "r"
Sleep 2s

# Navigate to Context Preview
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms
Enter
Sleep 1s

# Navigate through sections
Type "l"
Sleep 1s
Type "l"
Sleep 1s
Type "h"
Sleep 1s

# Toggle full content view
Enter
Sleep 2s
Enter
Sleep 1s

# Show template selection
Type "t"
Sleep 2s

# Navigate templates
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms

# Cancel template selection
Type "esc"
Sleep 1s

# Try edit mode
Type "e"
Sleep 2s

# Cancel edit
Type "esc"
Sleep 1s

# Exit preview
Type "esc"
Sleep 1s

# Exit application
Type "q"