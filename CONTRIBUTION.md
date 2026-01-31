# Contributing to Trelay

Thank you for your interest in contributing to Trelay. This project is under active development, and we welcome contributions that align with our goals of being developer-first and privacy-respecting.

## Development Process

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Implement your changes following the [Coding Standards](#coding-standards).
4. Ensure all tests pass and the code is formatted.
5. Submit a pull request with a clear description of the changes.

## Coding Standards

- Follow standard Go idioms and conventions.
- Use `go fmt` for formatting.
- Ensure `go vet` passes without issues.
- Keep the core logic in the `internal/core` package to maintain separation from delivery layers (API/CLI).
- Maintain structured logging using `zerolog`.
- Follow the existing monorepo structure.

## Pull Request Guidelines

- Provide a concise but descriptive title.
- Detail the "why" behind the change in the PR body.
- Reference any related issues.
- Keep PRs focused on a single change.

## Reporting Issues

Use GitHub Issues to report bugs or suggest features. Please provide as much detail as possible, including steps to reproduce for bugs.
