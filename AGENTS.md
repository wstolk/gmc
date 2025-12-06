# GMC Agent Guidelines

## Build/Test Commands
- Build: `go build`
- Run: `go run .`
- Test all: `go test ./...`
- Test single package: `go test -v ./internal/package_name`
- Test single function: `go test -v -run TestFunctionName ./internal/package_name`
- Lint: `golangci-lint run`
- Format: `gofmt -w .`
- Vet: `go vet ./...`

## Code Style Guidelines
- **Go Version**: 1.25.4 minimum
- **Formatting**: Use `gofmt` for standard Go formatting
- **Linting**: Use `golangci-lint run` for code quality checks
- **Naming**: PascalCase for exported functions/types, camelCase for unexported/variables
- **Packages**: lowercase, single word when possible (e.g., `git`, `ui`)
- **Error handling**: return errors with `fmt.Errorf` and `%w` verb, don't panic except for exceptional cases
- **Imports**: standard library first, then third-party (alphabetized), then local packages
- **Dependencies**: Use Cobra for CLI commands, go-git for Git operations
- **Functions**: Keep small and focused on single responsibility
- **Testing**: Use standard Go testing with table-driven tests where appropriate
- **Types**: Use struct types for complex data, interfaces for abstraction