# aggregate

The `aggregate` package counts occurrences of a log field value over a time window and emits a summary JSON line on each flush.

## Usage

```go
import "logpipe/internal/aggregate"

a := aggregate.New("level", 10*time.Second, func(b []byte) error {
    fmt.Println(string(b))
    return nil
})
a.Start()
defer a.Stop()

a.Record([]byte(`{"level":"info","msg":"started"}`))
```

## Config file (JSON)

```json
{
  "field": "level",
  "interval": "10s"
}
```

| Key        | Type   | Default | Description                          |
|------------|--------|---------|--------------------------------------|
| `field`    | string | —       | JSON field to aggregate by (required)|
| `interval` | string | `10s`   | Flush interval (Go duration string)  |

## Output format

Each flush emits one JSON line:

```json
{"_aggregate":"level","counts":{"info":42,"warn":5,"error":1}}
```
