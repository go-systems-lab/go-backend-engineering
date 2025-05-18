# Social Go

## 1. Initialize go module

```bash
go mod init github.com/kuluruvineeth/social-go
```

## 2. Project Structure

```
social-go/
├── bin/            # Compiled binaries
├── cmd/            # Application entry points
│   ├── api/        # API server
│   └── migrate/    # Database migration tool
├── docs/           # Documentation
├── internal/       # Private application code
├── scripts/        # Build/automation scripts
└── web/            # Web assets and frontend code
```

## 3. Development with Hot Reload

This project uses [Air](github.com/air-verse/air) for hot reloading during development. Air watches your files, rebuilds, and restarts the application when changes are detected.

### Install Air

```bash
# Install Air globally
go install github.com/air-verse/air@latest
```

### Run the application with Air

```bash
# Start the API server with hot reload
air init
air
```

### Configuration

The Air configuration is in `.air.toml` with the following settings:
- Watches `.go`, `.tpl`, `.tmpl`, and `.html` files
- Excludes certain directories from watching (assets, bin, vendor, etc.)
- Builds the API server from `./cmd/api`
- Outputs the binary to `./bin/main`