# Development Rules

## Code Generation
- When generating structs, always include JSON tags: `json:"fieldName"`.
- If using SQL, use `sqlx` or `pgx` styles (define which one here).
- Avoid `panic` in production code; always return `error`.

## Documentation
- All exported functions and types must have GoDoc comments.
- Comments should be full sentences starting with the function name.

## Dependencies
- Prefer standard library (`net/http`, `encoding/json`) over external dependencies where possible to keep binary size small.