# batch

The `batch` package buffers structured log lines and flushes them as a group either when a size threshold is reached or a time interval elapses.

## Usage

```go
b := batch.New(50, 2*time.Second, func(entries []map[string]any) {
    // handle flushed batch
})
defer b.Stop()

b.Write([]byte(`{"level":"info","msg":"hello"}`))
```

## Config file (JSON)

```json
{
  "size": 50,
  "interval": "2s"
}
```

| Field      | Type   | Default | Description                        |
|------------|--------|---------|------------------------------------|
| `size`     | int    | 100     | Max entries before forced flush    |
| `interval` | string | `"5s"`  | Max time between flushes           |
