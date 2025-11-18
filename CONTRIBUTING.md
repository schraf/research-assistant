# Contributing to Research Assistant

Thank you for your interest in contributing to Research Assistant! This document provides guidelines and instructions for contributing.

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers and help them learn
- Focus on constructive feedback
- Respect different viewpoints and experiences

## How to Contribute

### Reporting Bugs

1. Check if the bug has already been reported in the Issues section
2. If not, create a new issue with:
   - A clear, descriptive title
   - Steps to reproduce the issue
   - Expected vs. actual behavior
   - Environment details (OS, Go version, etc.)
   - Any relevant error messages or logs

### Suggesting Features

1. Check if the feature has already been suggested
2. Create a new issue with:
   - A clear description of the feature
   - Use cases and benefits
   - Any implementation ideas (if you have them)

### Submitting Code Changes

1. **Fork the repository** and clone your fork
2. **Create a branch** for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. **Make your changes** following the coding standards below
4. **Test your changes** thoroughly
5. **Commit your changes** with clear, descriptive commit messages
6. **Push to your fork** and create a Pull Request

## Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/schraf/research-assistant.git
   cd research-assistant
   ```

2. Install dependencies:
   ```bash
   make deps
   ```

3. Build the project:
   ```bash
   make build
   ```

4. Run tests:
   ```bash
   make test
   ```

## Coding Standards

- **Go formatting**: Run `make fmt` before committing
- **Code quality**: Run `make vet` to check for issues
- **Testing**: Add tests for new functionality
- **Documentation**: Update README.md or relevant docs for user-facing changes
- **Commit messages**: Use clear, descriptive messages in the present tense (e.g., "Add feature X" not "Added feature X")

## Pull Request Process

1. Ensure your code follows the coding standards
2. Update documentation if needed
3. Add or update tests as appropriate
4. Ensure all tests pass: `make test`
5. Create a clear PR description explaining:
   - What changes were made
   - Why the changes were made
   - How to test the changes

## Questions?

If you have questions about contributing, feel free to open an issue with the "question" label.

Thank you for contributing to Research Assistant!
