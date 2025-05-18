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

## 4. Environment Management with direnv

This project uses [direnv](https://direnv.net/) to manage environment variables. direnv loads and unloads environment variables based on the current directory.

### Install direnv

```bash
# macOS
brew install direnv

# Ubuntu/Debian
sudo apt-get install direnv

# Other platforms
# See: https://direnv.net/docs/installation.html
```

### Setup

1. Add direnv hook to your shell:

```bash
# For bash
echo 'eval "$(direnv hook bash)"' >> ~/.bashrc

# For zsh
echo 'eval "$(direnv hook zsh)"' >> ~/.zshrc
```

2. Create a `.envrc` file in your project root:

```bash
touch .envrc
```

3. Add your environment variables to the `.envrc` file:

```bash
# Example .envrc
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=social_go
export API_PORT=8080
```

4. Allow the `.envrc` file:

```bash
direnv allow
```

Now, whenever you navigate to your project directory, direnv will automatically load these environment variables, and unload them when you leave.