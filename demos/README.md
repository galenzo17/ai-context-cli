# 🎬 CLI Demos

This directory contains VHS recordings that demonstrate the key features of the AI Context CLI application.

## 📋 Available Demos

### 1. Main Menu Navigation (`01-main-menu.gif`)
- **Duration**: ~15 seconds
- **Features Shown**:
  - Professional CLI banner and interface
  - Menu navigation with j/k keys
  - Help system with ? key
  - Clean exit with q key
- **Key Highlights**: Shows the polished main interface and intuitive navigation

### 2. Add Context to All Files (`02-add-context-all.gif`)
- **Duration**: ~20 seconds  
- **Features Shown**:
  - Complete project scanning workflow
  - Real-time progress indicators
  - Loading animations and feedback
  - Context generation results
- **Key Highlights**: Demonstrates the core functionality of scanning entire projects

### 3. Folder Browser Navigation (`03-folder-browser.gif`)
- **Duration**: ~25 seconds
- **Features Shown**:
  - Interactive folder tree navigation
  - Expand/collapse folders with h/l keys
  - Folder statistics toggle with s key
  - Folder selection and confirmation dialog
  - Smart folder scanning workflow
- **Key Highlights**: Shows the powerful folder browser with real-time stats

### 4. Context Preview and Editing (`04-context-preview.gif`)
- **Duration**: ~30 seconds
- **Features Shown**:
  - Context preview interface with token estimation
  - Section navigation with arrow keys
  - Template selection system
  - Edit mode demonstration
  - Professional context management
- **Key Highlights**: Demonstrates advanced context editing and template features

### 5. Navigation System (`05-navigation-system.gif`)
- **Duration**: ~20 seconds
- **Features Shown**:
  - Breadcrumb navigation system
  - Screen transitions and back navigation
  - Consistent UI patterns across screens
  - ESC key functionality
- **Key Highlights**: Shows the robust navigation architecture

## 🎥 Recording New Demos

### Prerequisites
Install VHS (Video High-Speed):
```bash
go install github.com/charmbracelet/vhs@latest
```

### Recording Individual Demos
```bash
# Record a specific demo
vhs demos/01-main-menu.tape

# Record all demos at once
./demos/record-all.sh
```

### Creating New Demos
1. Create a new `.tape` file in the `demos/` directory
2. Use VHS syntax to define the recording:
   ```
   Output demos/my-demo.gif
   Set Theme "Catppuccin Mocha"
   Set FontSize 16
   Set Width 1200
   Set Height 800
   
   Type "command here"
   Sleep 1s
   Enter
   ```
3. Add the demo to `record-all.sh`

## 📐 VHS Configuration

All demos use consistent settings:
- **Theme**: Catppuccin Mocha (dark theme for better visibility)
- **Font Size**: 16 (readable in documentation)
- **Dimensions**: 1200x800 (good balance for web viewing)
- **Format**: GIF (universal compatibility)

## 🎯 Usage in Documentation

These GIFs can be embedded in:
- README.md files
- GitHub Issues and PRs
- Documentation websites
- Blog posts and tutorials

Example embedding:
```markdown
![Main Menu Demo](demos/01-main-menu.gif)
```

## 🔄 Updating Demos

When adding new features:
1. Update existing relevant demos
2. Create new demos for major features
3. Re-record all demos to maintain consistency
4. Update this README with new demo descriptions

## 📊 Demo Metrics

| Demo | Duration | File Size | Features |
|------|----------|-----------|----------|
| Main Menu | ~15s | ~500KB | Navigation, Help |
| Add Context All | ~20s | ~800KB | Scanning, Progress |
| Folder Browser | ~25s | ~1MB | Tree Navigation, Stats |
| Context Preview | ~30s | ~1.2MB | Preview, Editing, Templates |
| Navigation System | ~20s | ~700KB | Breadcrumbs, Transitions |

**Total**: ~110 seconds of demo content showcasing all major features