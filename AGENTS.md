# AGENTS Instructions for git-auto

## Commit Messages
- Use Conventional Commits style: `<type>: <description>`
  - Example: `feat: add custom flag for greet`
  - Types include: feat, fix, docs, test, chore, etc.
- Keep the summary imperative, under 50 characters, and start with a lowercase letter.

## Coding Style
- Format Go code using `go fmt` and ensure code quality by running `golangci-lint run ./...`.

## Testing
- Run unit tests using `go test ./... -v`.
- Integration tests are gated by setting the environment variable `INTEGRATION_TEST=true`.
  To run integration tests, export this variable in your shell before executing tests, for example:
  ```
  export INTEGRATION_TEST=true
  go test ./... -v
  ```
- Ensure `mise run go-test-all` passes if you use mise.

## Pull Request Notes
### Summary
- Clearly describe what was changed and why.

### Testing
- Include test results if applicable.
- If tests could not be run, explain why.

## Directory Overview
- `cmd/`    – subcommands (`greet`, `version`, etc.)
- `main.go` – entry point calling `cmd.Execute()`
- `mise.toml` – defines reusable tasks for building, testing, and linting (`mise run <task>`).
