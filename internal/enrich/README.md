# enrich

The `enrich` package adds static or derived fields to structured JSON log lines.

## Rules

Each rule specifies:
- `field` — the destination field name (required)
- `value` — a static string value to set
- `copy_of` — copy the value from an existing field

Exactly one of `value` or `copy_of` must be set per rule. If both are provided,
`value` takes precedence. If `copy_of` references a field that does not exist,
the destination field is left unset.

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
cfg, err := enrich.LoadConfig("enrich.json")
if err != nil {
    log.Fatal(err)
}
e := enrich.NewFromConfig(cfg)

out, err := e.Apply(line)
if err != nil {
    // line was not valid JSON; original line is returned unchanged
}

// or as middleware
next := e.Wrap(outputWriter.Write)
```

## Error handling

`Apply` returns an error if the input is not a valid JSON object. In that case
the original line is returned unmodified so callers can decide whether to drop,
pass through, or log the malformed line.
