# Golang Technology Stack & Standards

## Core Principles
- Use **Go 1.22+** features (iterators, min/max).
- Write **Idiomatic Go**: Prefer simplicity over abstraction. Avoid "Java-style" getters/setters or unnecessary interfaces.
- **Error Handling**:
  - Always handle errors immediately.
  - Use `errors.Is` and `errors.As` for checks.
  - Wrap errors with context using `fmt.Errorf("doing action: %w", err)`.
- **Concurrency**:
  - Always pass `context.Context` as the first argument to functions involving I/O.
  - Use `errgroup` over raw `sync.WaitGroup` for error propagation.

## Project Structure
- `/cmd`: Main applications.
- `/internal`: Private application and business logic (cannot be imported by others).
- `/pkg`: Library code ok to use by external applications.

## Testing
- Use the standard `testing` package.
- **Strictly use Table-Driven Tests** for logic.
- Use `testify/assert` for assertions if needed, but standard library is preferred.