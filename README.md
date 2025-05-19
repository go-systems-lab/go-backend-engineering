# GO Backend Engineering

## Getting Started

### 1. Initialize Go Module

```bash
go mod init github.com/kuluruvineeth/social-go
```

### 2. Project Structure

```
social-go/
├── bin/                    # Compiled binaries and executables
├── cmd/                    # Main applications of the project
│   ├── api/               # HTTP API server implementation
│   │   └── main.go       # Entry point for the API server
│   └── migrate/          # Database migration tooling
│       └── migrations/   # Migration files for database schema changes
├── internal/              # Private application and library code
├── scripts/               # Setting up server and starting scripts
└── docs/                  # OpenAPI/Swagger specs, JSON schema files, protocol definition files
```

The structure follows standard Go project layout conventions:
- `bin/`: Stores compiled applications
- `cmd/`: Contains the main applications of the project
- `cmd/api/`: Houses the HTTP server and API endpoints
- `cmd/migrate/`: Contains database migration tools and scripts
- `internal/`: Private code that's specific to this project
- `pkg/`: Reusable packages that can be imported by other projects
- `docs/`: API documentation and specifications

### 3. Install Dependencies

```bash
# Install external packages
go get -u github.com/go-chi/chi/v5
```

This adds the dependency to our `go.mod` file and creates a `go.sum` file.

### 4. Development with Hot Reload

This project uses [Air](github.com/air-verse/air) for hot reloading during development. Air watches your files, rebuilds, and restarts the application when changes are detected.

#### Install Air

```bash
# Install Air globally
go install github.com/air-verse/air@latest
echo 'export PATH=$PATH:~/go/bin' >> ~/.zshrc
source ~/.zshrc
```

#### Initialize Air

```bash
# Create Air configuration file
air init
```

#### Run the application with Air

```bash
# Start the API server with hot reload
air
```

#### Configuration

The Air configuration is in `.air.toml` with the following settings:
- Watches `.go`, `.tpl`, `.tmpl`, and `.html` files
- Excludes certain directories from watching (assets, bin, vendor, etc.)
- Builds the API server from `./cmd/api`
- Outputs the binary to `./bin/main`

### 5. Environment Management with direnv

This project uses [direnv](https://direnv.net/) to manage environment variables. direnv loads and unloads environment variables based on the current directory.

#### Install direnv

```bash
# macOS
brew install direnv

# Ubuntu/Debian
sudo apt-get install direnv

# Other platforms
# See: https://direnv.net/docs/installation.html
```

#### Setup

1. Add direnv hook to your shell:

```bash
# For bash
echo 'eval "$(direnv hook bash)"' >> ~/.bashrc

# For zsh
echo 'eval "$(direnv hook zsh)"' >> ~/.zshrc
source ~/.zshrc
```

2. Create a `.envrc` file in your project root:

```bash
# Example .envrc
export ADDR=":8081"
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=admin
export DB_PASSWORD=adminpassword
export DB_NAME=social_go
```

3. Allow the `.envrc` file:

```bash
direnv allow .
```

Now, whenever you navigate to your project directory, direnv will automatically load these environment variables, and unload them when you leave.

### 6. Database Setup

#### Start PostgreSQL Container

```bash
docker compose up --build
```

#### Database Migrations

This project uses [golang-migrate](https://github.com/golang-migrate/migrate) for database migrations.

##### Install golang-migrate

```bash
brew install golang-migrate
```

##### Create Migrations

```bash
# Create a new migration
migrate create -seq -ext sql -dir ./cmd/migrate/migrations create_users
```

The flags used:
- `-seq`: Creates sequential migration files (001, 002, etc.)
- `-ext sql`: Sets the file extension to .sql
- `-dir`: Specifies the migrations directory

##### Run Migrations

```bash
# Apply migrations
migrate -path=./cmd/migrate/migrations -database="postgres://admin:adminpassword@localhost/social_go?sslmode=disable" up

# Rollback migrations
migrate -path=./cmd/migrate/migrations -database="postgres://admin:adminpassword@localhost/social_go?sslmode=disable" down
```

### 7. Using Makefile

The project includes a Makefile to simplify common commands. Here's what each target does:

```makefile
# Create a new migration
make migration create_posts

# Apply all pending migrations
make migrate-up

# Rollback migrations (specify number to rollback)
make migrate-down 1

# Seed the database
make seed
```

The Makefile:
- Imports environment variables from `.envrc`
- Sets the migrations path
- Provides shortcuts for migration commands
- Handles arguments for migration commands
- Includes a seed command for database seeding