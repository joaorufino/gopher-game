# CLAUDE.md - Development Guidelines

## Build Commands
- Build all: `make all` (clean, format, deps, build, wasm, assets)
- Build native: `make build`
- Build WebAssembly: `make wasm`
- Run locally: `go run cmd/main.go` (serves at localhost:8000)
- Run tests: `make test` or `go test ./...`
- Run single test: `go test ./pkg/animation -run TestAnimation/Creating_a_new_animation`
- Format code: `make format` or `go fmt ./...`

## Code Style Guidelines
- Use camelCase for variable/function names, PascalCase for exported symbols
- Group imports: standard lib, then third party, then local packages
- Use error handling with returns, not panics
- Add comprehensive comments for exported functions and types
- Use dependency injection with `go.uber.org/fx`
- Follow GoDoc conventions for documentation
- Use `sync.Mutex` for concurrent access protection
- Use interfaces for abstractions and loose coupling
- Employ GoLang BDD testing with Ginkgo/Gomega for feature tests

## Project Structure
- `/cmd`: Main application entry points
- `/internal`: Private application code
- `/pkg`: Public, reusable packages
- `/assets`: Game resources (images, audio, configs)