#!/bin/bash

# Script to record all VHS demos
# Make sure you have vhs installed: go install github.com/charmbracelet/vhs@latest

set -e

echo "🎬 Recording CLI demos with VHS..."

# Create demos directory if it doesn't exist
mkdir -p demos

# Record each demo
echo "📹 Recording Demo 1: Main Menu Navigation"
vhs demos/01-main-menu.tape

echo "📹 Recording Demo 2: Add Context to All Files" 
vhs demos/02-add-context-all.tape

echo "📹 Recording Demo 3: Folder Browser Navigation"
vhs demos/03-folder-browser.tape

echo "📹 Recording Demo 4: Context Preview and Editing"
vhs demos/04-context-preview.tape

echo "📹 Recording Demo 5: Navigation System"
vhs demos/05-navigation-system.tape

echo "✅ All demos recorded successfully!"
echo "📁 GIFs saved in demos/ directory"

# List the generated files
echo ""
echo "Generated files:"
ls -la demos/*.gif