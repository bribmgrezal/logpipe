# enrich

The `enrich` package adds static or derived fields to structured JSON log lines.

## Rules

Each rule specifies:
- `field` — the destination field name (required)
- `value` — a static string value to set
- `copy_of` — copy the value from an existing field

If `copy_of` references a field that does not exist, the destination field is left unset.

## Config file

```json
{
  "rules": [
    { "field": "env",    "value": "production" },
    { "field": "svc",    "copy_of": "service" }
  ]
}
```

## Usage

```go
cfg, _ := enrich.LoadConfig("enrich.json")
e := enrich.NewFromConfig(cfg)

out, err := e.Apply(line)

// or as middleware
next := e.Wrap(outputWriter.Write)
```
