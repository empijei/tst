# tst

`tst` is a collection of small helpers designed to make Go tests more readable and easier to maintain.
It provides functions for common testing patterns, making the intent of your tests clearer.

## Key Features

- **Lean Error Handling:** Reduce `if err != nil { t.Fatalf(...) }` blocks to a single line.
- **Value Unwrapping:** Extract values from functions that return `(value V, err error)` or `(value V, ok bool)` directly in your assertions.
- **Deep Equality:** Built-in support for `google/go-cmp` for expressive and readable diffs.
- **Concurrency Helpers:** Shorthands for parallel tests.

## Installation

```bash
go get github.com/empijei/tst
```

## Usage Examples

### Error Handling and Value Unwrapping

Instead of:

```go
t.Parallel()
f, err := os.Open("config.json")
if err != nil {
    t.Fatalf("failed to open config: %v", err)
}
defer f.Close()
```

Use:

```go
a := tst.Go(t)
f := a.Do(os.Open("config.json"))
defer f.Close()
```

### Deep Equality with `Is`

`(*Assertions).Is` uses `go-cmp` to provide detailed diffs when values don't match.

```go
want := &User{Name: "Alice", Age: 30}
got := FetchUser(1)
a.Is(want, got)
```

### Asserting Errors with `Err`

Verify that an error is not nil and optionally contains a specific substring.

```go
_, err := ProcessData(invalidInput)
a.Err("invalid input", err, t)
```

### Stopping Tests Early with `Ko`

Prevent a flood of error messages by stopping the test if a previous assertion failed.

```go
a := tst.Go(t)
a.Is(expectedHeader, actualHeader, t)
a.Ko(t) // Stop here if header check failed, as subsequent tests might be meaningless.

a.Is(expectedBody, actualBody, t)
```

### Parallel Tests

`Go` is a shorthand for `t.Parallel()` that also returns the assertions set.

```go
func TestSomething(t *testing.T) {
    // Defaults to calling t.Parallel()
    a := tst.Go(t)
    // Run your assertions with a.[...]
}
```

But tst can be used synchronously by calling `Sync`.

```go
func TestSomething(t *testing.T) {
    // This doesn't call t.Parallel()
    a := tst.Sync(t)
    // Run your assertions with a.[...]
}
```
