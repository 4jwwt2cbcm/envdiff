# envdiff

Compare `.env` files across environments and report missing or mismatched keys.

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git
cd envdiff && go build -o envdiff .
```

## Usage

```bash
envdiff [flags] <base> <target> [target...]
```

### Example

```bash
# Compare .env.example against your local and production env files
envdiff .env.example .env .env.production
```

**Sample output:**

```
Missing in .env:
  - DATABASE_URL
  - REDIS_HOST

Missing in .env.production:
  - DEBUG
  - LOG_LEVEL

Mismatched keys:
  APP_ENV: "development" (.env) vs "production" (.env.production)
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--strict` | `false` | Exit with non-zero status if any differences are found |
| `--values` | `false` | Include value comparison in output |
| `--json` | `false` | Output results as JSON |
| `--ignore` | | Comma-separated list of keys to ignore during comparison |

## Why envdiff?

Keeping `.env` files in sync across environments is error-prone. `envdiff` makes it easy to catch missing variables before they cause issues in staging or production.

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

## License

MIT © [yourusername](https://github.com/yourusername)
