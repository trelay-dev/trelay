# Contributing to Trelay

Thanks for looking at the code. Trelay is under active development; patches that fit the project (self-hosted, privacy-minded, small surface area) are welcome.

## How to work on a change

1. Fork the repo and branch from `main`.
2. Make your change in small, reviewable steps.
3. Match existing style: `go fmt`, `go vet`, and the patterns already in `internal/core` vs API/CLI.
4. Run tests before you open a PR.
5. Open a pull request with a short title and enough context in the body that a reviewer can follow the intent.

## Code style (Go and layout)

- Idiomatic Go; run `go fmt` on touched files.
- Keep domain logic in `internal/core`; wire it through API or CLI without duplicating rules.
- Logging stays on `zerolog` like the rest of the tree.
- Respect the repo layout you see today (no new top-level packages without a good reason).

## Pull requests

- One logical change per PR when possible.
- Link related issues if there are any.
- Say what you changed and why; screenshots help for UI work.

## Issues

Use GitHub Issues for bugs and feature ideas. For bugs, include version or commit, what you did, what you expected, and what happened instead.

## Security

Do **not** open a public issue for security problems. Follow **[SECURITY.md](SECURITY.md)** and use private vulnerability reporting on the repo’s **Security** tab.
