# GMC Agent Guidelines

## Build/Test Commands
- Build: `go build`
- Run: `go run .`
- Test all: `go test ./...`
- Test single package: `go test -v ./pkg/package_name`
- Test single function: `go test -v -run TestFunctionName ./pkg/package_name`
- Lint: `golangci-lint run`
- Format: `gofmt -w .`
- Vet: `go vet ./...`

## Code Style Guidelines
- Follow standard Go formatting with `gofmt`
- Use `golangci-lint` for code quality checks
- Package naming: lowercase, single word when possible
- Function naming: PascalCase for exported, camelCase for unexported
- Variable naming: camelCase, descriptive names
- Error handling: return errors, don't panic except for truly exceptional cases
- Imports: standard library first, then third-party, then local packages
- Use `context.Context` for cancellation and timeouts
- Prefer `interface{}` over `any` for Go 1.18+
- Use struct embedding judiciously
- Keep functions small and focused on single responsibility