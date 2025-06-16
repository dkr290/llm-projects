# Desktop Automation CLI

Command-line tool for desktop automation tasks using robotgo.

## Installation

```bash
go install github.com/dmahlow/desktop-automation/cmd/desktop-automation@latest
```

## Build from Source

```bash
git clone https://github.com/dmahlow/desktop-automation.git
cd desktop-automation
task build
```

## Usage

### Basic Commands

Move mouse cursor:
```bash
desktop-automation move 100 200
```

Click at coordinates:
```bash
desktop-automation click 100 200
```

Type text:
```bash
desktop-automation type "Hello World"
```

Interactive terminal UI:
```bash
desktop-automation tui
```

### Available Commands

- `move [x] [y]` - Move mouse to coordinates
- `click [x] [y]` - Click at coordinates
- `type [text]` - Type text at cursor position
- `tui` - Launch interactive terminal interface

## Requirements

- Go 1.21+
- Platform-specific dependencies for robotgo (see robotgo documentation)

## Development

Use Task for development tasks:

```bash
task build    # Build binary
task run      # Build and run
task clean    # Clean artifacts
task deps     # Download dependencies
task test     # Run tests
```
