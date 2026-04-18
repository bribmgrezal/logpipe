# format

The `format` package provides log line formatting for logpipe.

## Overview

Convert structured JSON log lines into human-readable or custom-formatted strings using a simple template syntax.

## Template Syntax

Use `{field}` placeholders to reference top-level JSON fields.

```
{level} [{time}] {msg}
```

If a field is missing from the log line, it is replaced with an empty string.

If no template is provided, the output is pretty-printed JSON.

## Config File

```json
{
  "template": "{level} {msg}"
}
```

## Usage

```go
cfg, _ := format.LoadConfig("format.json")
f := format.NewFromConfig(cfg)
out, err := f.Apply(line)
```
