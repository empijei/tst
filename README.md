# tst

`tst` is a collection of small, focused helpers designed to make Go tests leaner, more readable, and easier to maintain. It provides intuitive functions for common testing patterns, reducing boilerplate and making the intent of your tests clearer.

## Key Features

- **Lean Error Handling:** Reduce `if err != nil { t.Fatalf(...) }` blocks to a single line.
- **Value Unwrapping:** Extract values from functions that return `(value V, err error)` or `(value V, ok bool)` directly in your assertions.
- **Deep Equality:** Built-in support for `google/go-cmp` for powerful, expressive and readable diffs.
- **Cascading Failure Prevention:** Easily stop tests early if a critical assertion fails.
- **Concurrency Helpers:** Shorthand for parallel tests and context management.

## Installation

```bash
go get github.com/empijei/tst
```

## Usage Examples

### Error Handling and Value Unwrapping

Instead of:

```go
f, err := os.Open("config.json")
if err != nil {
    t.Fatalf("failed to open config: %v", err)
}
defer f.Close()
```

Use:

```go
f := tst.Do(os.Open("config.json"))(t)
defer f.Close()
```

### Deep Equality with `tst.Is`

`tst.Is` uses `go-cmp` to provide detailed diffs when values don't match.

```go
want := &User{Name: "Alice", Age: 30}
got := FetchUser(1)
tst.Is(want, got, t)
```

### Asserting Errors with `tst.Err`

Verify that an error is not nil and optionally contains a specific substring.

```go
_, err := ProcessData(invalidInput)
tst.Err("invalid input", err, t)
```

### Stopping Tests Early with `tst.Ko`

Prevent a flood of error messages by stopping the test if a previous assertion failed.

```go
tst.Is(expectedHeader, actualHeader, t)
tst.Ko(t) // Stop here if header check failed, as subsequent tests might be meaningless.

tst.Is(expectedBody, actualBody, t)
```

### Parallel Tests and Context

`tst.Go` is a convenient shorthand for `t.Parallel()` that also returns the test context.

```go
func TestSomething(t *testing.T) {
    ctx := tst.Go(t)
    // Run your test using ctx...
}
```

## API Reference

### `Do[V any](v V, err error) func(t Test) V`

Unwraps a result and stops the test immediately (`t.Fatalf`) if an error occurred.

### `Do2[V1, V2 any](v1 V1, v2 V2, err error) func(t Test) (V1, V2)`

Like `Do`, but for functions that return two values and an error.

### `DoB[V any](v V, ok bool) func(t Test) V`

Unwraps a result and stops the test immediately (`t.Fatalf`) if `ok` is false.

### `No(err error, t Test)`

Stops the test immediately (`t.Fatalf`) if the provided error is not nil.

### `Be(ok bool, t Test)`

Stops the test immediately (`t.Fatalf`) if `ok` is false.

### `Is[T any](want, got T, t Test, opts ...cmp.Option)`

Checks that `want` matches `got` via `cmp.Diff`. Calls `t.Errorf` if there's a mismatch.

### `Err(errorSubMessage string, err error, t Test)`

Checks if the provided error is not nil and contains an optional message.

### `Ko(t Test)`

Stops the test immediately (`t.Fatalf`) if the test has already failed.

### `Go(t PTest) context.Context`

Shorthand for `t.Parallel()` and returns the test context.
