# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a sample plugin for the UDB (Universal Display Board) project. It demonstrates how to create custom boards and datasources that integrate with the UDB core system. The plugin is written in Go and implements the `UdbPlugin` interface from the `udb-plugin-library`.

## Architecture

### Plugin Structure

The plugin is loaded by udb-core and must export a `Plugin` variable of type `udb_plugin_library.UdbPlugin`:

- **Main entry point**: `udb_sample_plugin.go` - Implements the `UdbPlugin` interface with methods to register boards and datasources
- **Boards**: Implement the `types.Board[T]` interface - visual displays that render to images
  - `boards/single_colour.go` - Example static board that fills the display with a configurable color
- **Datasources**: Implement the `types.Datasource[T]` interface - provide data to boards
  - `datasources/current_time.go` - Example datasource that returns the current time

### Key Concepts

- **Boards** must implement: `GetId()`, `GetName()`, `GetSupportedDimensions()`, `GetType()`, `GetDatasourceType()`, `Init()`, and `Render()`
- **Static vs Dynamic boards**: `BoardTypeStatic` boards don't need real-time updates; `BoardTypeDynamic` requires periodic data refresh
- **Configuration**: Boards accept optional JSON configuration via `Init()` (e.g., `SingleColourBoard` accepts a hex color string)
- **Rendering**: Boards return `image.Image` objects generated from their state

## Development Commands

```bash
# Build the plugin as a shared library (required for udb-core loading)
go build -o udb-plugin-samples.so -buildmode=plugin

# Run tests
go test ./...

# Run a specific test
go test -run TestName ./path/to/package

# Format code
go fmt ./...

# Run linter (if golangci-lint is installed)
golangci-lint run

# Check for vet issues
go vet ./...

# View dependencies
go mod tidy
go mod graph
```

## Adding New Content

To add a new board:
1. Create a new file in the `boards/` directory implementing the `types.Board[T]` interface
2. Instantiate it in `udb_sample_plugin.go` within `GetBoardMap()`
3. Return it from `GetAllBoards()`

To add a new datasource:
1. Create a new file in the `datasources/` directory implementing the `types.Datasource[T]` interface
2. Register it in the plugin if needed

## Dependencies

- `github.com/benwiebe/udb-plugin-library` - Core plugin interfaces and types
- `github.com/benwiebe/udb-core` - Board types and core functionality
