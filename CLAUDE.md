# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a sample plugin for the UDB (Universal Display Board) project. It demonstrates how to create custom boards and datasources that integrate with the UDB core system. The plugin is written in Go and implements the `UdbPlugin` interface from the `udb-plugin-library`.

## Architecture

### Plugin Structure

The plugin is loaded by udb-core and must export a `Plugin` variable of type `udb_plugin_library.UdbPlugin`:

- **Main entry point**: `udb_sample_plugin.go` - Implements the `UdbPlugin` interface with methods to register boards and datasources
- **Boards**: Implement the `types.Board[T]` interface - visual displays that render to images
  - `boards/single_colour.go` - Static board that fills the display with a configurable colour; pre-builds its image in `Init()`
  - `boards/digital_clock.go` - Dynamic board that displays the current time; pre-computes font layout in `Init()`
- **Datasources**: Implement the `types.Datasource[T]` interface - provide data to boards
  - `datasources/current_time.go` - Returns the current time via `time.Now()`; no background goroutine needed

### Key Concepts

- **Boards** must implement: `GetId()`, `GetName()`, `GetSupportedDimensions()`, `GetType()`, `GetDatasourceType()`, `Init()`, and `Render()`
- **Board types**: `BoardTypeStatic` renders once and holds the image; `BoardTypeAnimated` returns a pre-baked frame sequence; `BoardTypeDynamic` is called repeatedly with each call returning the next frame
- **Dimensions in `Init()`**: Boards receive `types.BoardDimensions` in `Init()`, not in `Render()`. Pre-compute any layout values (font sizes, image buffers, drawing coordinates) that depend on the display size here. `Render()` takes no parameters.
- **Configuration**: Boards accept optional JSON configuration via `Init()` (e.g., `SingleColourBoard` accepts a hex colour string; `DigitalClockBoard` accepts `format`, `colour`, and `blinkColon`)
- **Datasource lifecycle**: Datasources must implement `Start(ctx context.Context) error` (called once at startup — launch background goroutines here) and `DataChanged() <-chan struct{}` (return a channel to trigger immediate re-renders on data change, or `nil` for no push notifications)
- **Rendering**: `Render()` calls `datasource.GetData()` to get the latest cached value, constructs an `image.Image`, and returns it (or an `AnimationFrame` for animated/dynamic boards)

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
1. Create a new file in the `datasources/` directory implementing the `types.Datasource[T]` interface (including `Start()` and `DataChanged()`)
2. Register it in `udb_sample_plugin.go` within `GetDatasourceMap()`
3. Return it from `GetAllDatasources()`

## Dependencies

- `github.com/benwiebe/udb-plugin-library` - Core plugin interfaces and types
