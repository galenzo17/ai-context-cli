# Demo 5: Navigation System and Breadcrumbs
# This demo shows the navigation breadcrumb system

Output demos/05-navigation-system.gif
Set Theme "Catppuccin Mocha"
Set FontSize 16
Set Width 1200
Set Height 800

# Start the application
Type "go run cmd/ai-context-cli/main.go"
Sleep 1s
Enter
Sleep 2s

# Navigate to folder browser to show breadcrumbs
Type "j"
Sleep 500ms
Enter
Sleep 2s

# Show navigation breadcrumbs in action
Sleep 2s

# Go back with ESC
Type "esc"
Sleep 1s

# Navigate to context preview (after scanning)
Enter
Sleep 3s
Sleep 5s
Type "r"
Sleep 2s

# Go to context preview
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms
Enter
Sleep 2s

# Show breadcrumbs in context preview
Sleep 2s

# Go back
Type "esc"
Sleep 1s

# Exit
Type "q"