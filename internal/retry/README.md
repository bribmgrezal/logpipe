# retry

The `retry` package provides a writer wrapper that retries failed writes up to a configurable number of attempts with an optional delay between each attempt.

## Usage

```go
w := retry.New(3, 50*time.Millisecond, func(line []byte) error {
    return forwardToRemote(line)
})

err := w.Write([]byte(`{"level":"error","msg":"disk full"}`))
```

## Config file (JSON)

```json
{
  "max_attempts": 3,
  "delay_ms": 100
}
```

## Fields

| Field | Type | Description |
|---|---|---|
| `max_attempts` | int | Maximum number of write attempts (default: 1) |
| `delay_ms` | int | Milliseconds to wait between retries (default: 0) |

## Notes

- Invalid JSON input is rejected immediately without retrying.
- If `max_attempts` is less than 1, it defaults to 1.
