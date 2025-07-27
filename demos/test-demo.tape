# Test Demo: Quick CLI Test
Output demos/test-demo.gif
Set Theme "Catppuccin Mocha"
Set FontSize 16
Set Width 1200
Set Height 600

# Build the CLI first
Type "go build -o ai-context-cli cmd/ai-context-cli/main.go"
Enter
Sleep 2s

# Run the CLI
Type "./ai-context-cli"
Enter
Sleep 2s

# Navigate menu
Type "j"
Sleep 500ms
Type "k"
Sleep 500ms

# Show help
Type "?"
Sleep 2s
Type "esc"
Sleep 1s

# Exit
Type "q"
Sleep 1s