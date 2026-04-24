# cast

The `cast` module converts JSON field values to a target type in log lines.

## Supported Target Types

| Target   | Description                        |
|----------|------------------------------------|
| `string` | Converts the value to a string     |
| `int`    | Parses the value as a 64-bit int   |
| `float`  | Parses the value as a float64      |
| `bool`   | Parses the value as a boolean      |

## Configuration

```json
{
  "rules": [
    { "field": "status_code", "target": "int" },
    { "field": "latency",     "target": "float" },
    { "field": "active",      "target": "bool" }
  ]
}
```

## Behaviour

- Fields not present in the log line are silently skipped.
- If a value cannot be cast to the target type, the field is left unchanged.
- Invalid JSON input returns an error.
