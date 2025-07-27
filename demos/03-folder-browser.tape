# Demo 3: Folder Browser Navigation
# This demo shows the folder selection and browsing functionality

Output demos/03-folder-browser.gif
Set Theme "Catppuccin Mocha" 
Set FontSize 16
Set Width 1200
Set Height 800

# Start the application
Type "go run cmd/ai-context-cli/main.go"
Sleep 1s
Enter
Sleep 2s

# Navigate to "Add Context to Specific Folder"
Type "j"
Sleep 500ms
Enter
Sleep 1s

# Browse the folder tree
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms

# Expand a folder
Type "l"
Sleep 1s

# Navigate into subfolders
Type "j"
Sleep 500ms
Type "l" 
Sleep 1s

# Collapse folder
Type "h"
Sleep 1s

# Toggle stats
Type "s"
Sleep 1s
Type "s"
Sleep 1s

# Select a folder for confirmation
Type "j"
Sleep 500ms
Type "c"
Sleep 2s

# Confirm selection
Type "y"
Sleep 1s

# Watch scanning process
Sleep 3s

# Exit after completion
Type "r"
Sleep 2s
Type "q"