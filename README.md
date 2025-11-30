# Research Assistant

A content generator plugin for the [assistant](https://github.com/schraf/assistant) project that performs research on topics using Google Cloud APIs.

## Overview

This project provides a researcher content generator that integrates with the assistant framework. It ties together the researcher content generator with Google Cloud API endpoints to perform in-depth research on specified topics.

## Features

- **Research Generator**: Implements the `ContentGenerator` interface from the assistant project
- **Configurable Depth**: Supports three research depth levels (basic, medium, long)
- **Google Cloud Integration**: Uses Google Cloud APIs for research capabilities
- **Standalone CLI**: Can be run as a standalone command-line tool for testing

## Usage

### As a Plugin

This project is imported by the assistant project and automatically registers the "researcher" generator. The generator can be invoked through the assistant's API endpoints.

### As a Standalone Tool

Build and run the CLI tool:

```bash
make build
./researcher -topic "your research topic" -depth basic
```

Available depth options:
- `basic` - Short research depth
- `medium` - Medium research depth  
- `long` - Long research depth

## Development

```bash
# Build the project
make build

# Run tests
make test

# Format code
make fmt

# Run vetting
make vet
```

## Project Structure

- `cmd/main.go` - Standalone CLI application
- `pkg/generator/` - Generator plugin implementation
- `internal/researcher/` - Core research functionality

## Dependencies

- `github.com/schraf/assistant` - Assistant framework
- Google Cloud APIs (via `google.golang.org/genai`)
