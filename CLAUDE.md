# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

krapp is a CLI tool written in Go for creating and managing daily notes and inbox notes. It uses YAML configuration files and supports custom editors for quick note-taking workflows.

## Development Commands

### Build and Run
```bash
# Build the application
go build -o krapp ./cmd/krapp

# Install the application globally
go install github.com/ishida722/krapp-go/cmd/krapp@HEAD

# Run directly without building
go run ./cmd/krapp [command]
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test ./config
go test ./models
go test ./usecase
```

### Code Quality
```bash
# Format code
go fmt ./...

# Vet code for potential issues
go vet ./...

# Tidy dependencies
go mod tidy
```

## Architecture

### Core Components

- **cmd/krapp/main.go**: Main entry point with Cobra CLI commands
- **config/**: Configuration management with YAML support and merging of global/local configs
- **models/**: Core data structures (Note, FrontMatter) with YAML frontmatter parsing
- **usecase/**: Business logic for note creation, file operations, and Git sync

### Configuration System

The application uses a hierarchical configuration system:
1. Default configuration (hardcoded)
2. Global configuration (`~/.config/krapp/config.yaml` or `$XDG_CONFIG_HOME/krapp/config.yaml`)
3. Local configuration (`./.krapp_config.yaml`)

Local settings override global, which override defaults using the `mergo` library.

The global configuration follows XDG Base Directory Specification:
- Uses `$XDG_CONFIG_HOME/krapp/config.yaml` if `XDG_CONFIG_HOME` is set
- Falls back to `~/.config/krapp/config.yaml` if `XDG_CONFIG_HOME` is not set
- Automatically migrates legacy config from `~/.krapp_config.yaml` on first run

### Note Structure

Notes are Markdown files with optional YAML frontmatter:
- Daily notes: `notes/daily/YYYY/MM/YYYY-MM-DD.md`
- Inbox notes: `notes/inbox/YYYY-MM-DD-title.md`

### Key Interfaces

- `Config`: Provides base directory and note directory paths
- `InboxConfig`: Extends Config with inbox-specific methods

## Application Commands

- `krapp create-daily` (alias: `cd`): Create today's daily note
- `krapp create-inbox "title"` (alias: `ci`): Create inbox note with title
- `krapp config`: Print current configuration as YAML (also available as `krapp print-config`)
- `krapp sync`: Sync notes using Git
- `krapp import-notes [dir]` (alias: `in`): Import notes from directory
- `krapp import-issues`: Import GitHub issues as notes

All create commands support `-e/--edit` flag to open in editor after creation.

## Configuration Options

- `base_dir`: Root directory for notes (default: "./notes")
- `daily_note_dir`: Subdirectory for daily notes (default: "daily") 
- `inbox_dir`: Subdirectory for inbox notes (default: "inbox")
- `editor`: Command to open files (default: "vim")
- `editor_option`: Additional options for editor
- `with_always_open_editor`: Always open editor after note creation
- `daily_template`: YAML template for daily notes frontmatter
- `inbox_template`: YAML template for inbox notes frontmatter

## Development Guidelines

### Architecture
- **Follow Clean Architecture principles**
- Dependencies flow inward: cmd → usecase → models
- Business logic resides in the usecase layer
- Models contain core entities and business rules
- External dependencies (CLI, file system) are kept in outer layers

### Testing
- **Always write tests** for new functions and modifications
- Use table-driven tests where appropriate
- Test both success and error cases
- Follow the existing test patterns in `*_test.go` files

### Function Parameters
- **Use value receivers by default** when functions don't modify the receiver
- Use pointer receivers only when:
  - The function needs to modify the receiver
  - The struct is large and copying would impact performance significantly
  - Consistency is needed (if some methods use pointer receivers, use them for all methods on that type)

### Development Workflow
- **Always create a new branch** when implementing from an issue
- **Create a Pull Request** at the end of the implementation
- Branch naming convention: Use descriptive names like `feature/issue-123-add-new-command` or `fix/issue-456-config-bug`