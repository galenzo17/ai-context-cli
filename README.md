# AI Context CLI

A Terminal User Interface (TUI) wrapper for AI models with advanced context engineering capabilities, built with Go and Bubble Tea.

## Features

- Interactive terminal interface using Bubble Tea
- Support for multiple AI model providers
- Context engineering with customizable templates
- Session management for chat conversations
- Configuration management
- Extensible architecture

## Installation

```bash
git clone <repository-url>
cd ai-context-cli
go build -o ai-context-cli cmd/ai-context-cli/main.go
```

## Usage

### Interactive Mode
```bash
./ai-context-cli
```

### Command Line Options
```bash
./ai-context-cli help     # Show help
./ai-context-cli version  # Show version
```

## Project Structure

```
├── cmd/ai-context-cli/     # Main application entry point
├── internal/
│   ├── app/               # Bubble Tea application logic
│   ├── config/            # Configuration management
│   ├── models/            # AI model integrations
│   ├── ui/                # UI components
│   ├── context/           # Context engineering
│   └── ai/                # AI provider implementations
├── pkg/
│   ├── types/             # Shared types and structs
│   └── utils/             # Utility functions
└── test/
    ├── integration/       # Integration tests
    └── unit/              # Unit tests
```

## Configuration

The application stores configuration in `~/.ai-context-cli/config.json`. This includes:

- AI model configurations
- API keys and endpoints
- Context templates
- User preferences

## Development

### Running Tests
```bash
go test ./...
```

### Running Integration Tests
```bash
go test ./test/integration/...
```

### Building
```bash
go build -o ai-context-cli cmd/ai-context-cli/main.go
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.