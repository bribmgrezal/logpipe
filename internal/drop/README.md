# drop

The `drop` module discards structured log lines that match one or more rules. Lines that do not match any rule are passed through unchanged.

## Rules

Each rule targets a single JSON field and supports the following operators:

| Operator   | Behaviour                                      |
|------------|------------------------------------------------|
| `eq`       | Drop if the field value equals `value`         |
| `contains` | Drop if the field value contains `value`       |
| `exists`   | Drop if the field is present (any value)       |
| `missing`  | Drop if the field is absent from the record    |

## Configuration

```json
{
  "rules": [
    { "field": "level",   "operator": "eq",       "value": "debug" },
    { "field": "msg",     "operator": "contains",  "value": "healthcheck" },
    { "field": "internal","operator": "exists" },
    { "field": "req_id",  "operator": "missing" }
  ]
}
```

## Usage

```go
cfg, err := drop.LoadConfig("drop.json")
d := drop.NewFromConfig(cfg)

out, err := d.Apply(line)
if err != nil {
    // invalid JSON
}
if out == "" {
    // line was dropped
}
```
