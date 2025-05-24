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

### Debugging with VSCode and Delve

To debug your Go application in VSCode while using Air for live reloading, you can integrate Delve, the Go debugger.

**1. Install Delve:**

If you haven't already, install Delve:
```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

**2. Configure Air for Delve:**

Your `.air.toml` file needs to be configured to run your application with Delve. The key changes involve:
- Modifying the `cmd` under `[build]` to compile with debug flags:
  `cmd = "go build -gcflags="all=-N -l" -o ./bin/debug ./cmd/api"`
- Setting `bin` to the debug binary:
  `bin = "./bin/debug"`
- Updating `full_bin` to execute Delve, make it listen on a port (e.g., 2345), and run your debug binary:
  `full_bin = "dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient --continue --log exec ./bin/debug"`
  (These changes should already be in your `.air.toml` from our previous steps).

**3. Configure VSCode Launch:**

Create a `launch.json` file inside a `.vscode` directory in your project root (if it doesn't exist) with the following configuration:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Attach to Delve",
            "type": "go",
            "request": "attach",
            "mode": "remote",
            "remotePath": "${workspaceFolder}",
            "port": 2345,
            "host": "127.0.0.1",
            "showLog": true,
            "apiVersion": 2,
            "trace": "verbose"
        }
    ]
}
```
This configuration tells VSCode to attach to the Delve debugger listening on `127.0.0.1:2345`.

**4. Start Debugging:**

1.  Run Air in your terminal from the project root:
    ```bash
    air
    ```
    Air will compile your application and start Delve. Look for a message like `API server listening at: 127.0.0.1:2345` in Air's output.
2.  In VSCode, go to the "Run and Debug" view (usually the play button with a bug icon in the sidebar).
3.  Select the "Attach to Delve" configuration from the dropdown menu.
4.  Click the green play button (Start Debugging).

VSCode should now attach to Delve, allowing you to set breakpoints, inspect variables, etc.

**Important Note on Live Reloading:**

When you save a file:
- Air will detect the change, stop the current application (and Delve instance), rebuild, and restart the application with a *new* Delve instance.
- Your VSCode debugger, which was attached to the old Delve instance, will disconnect.
- **You will need to manually restart the "Attach to Delve" debugging session in VSCode** after Air has finished restarting your application. The `--continue` flag in the Delve command within `.air.toml` may help make this process smoother but does not guarantee automatic re-connection.

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

# For Generating Swagger Docs
```bash
go install github.com/swaggo/swag/cmd/swag@latest

swag init -d cmd/api # whereever the main.go file is there

go get -u github.com/swaggo/http-swagger

http://localhost:8080/v1/swagger/index.html
```

### 8. API Documentation with Swagger

This project uses [Swagger/OpenAPI](https://swagger.io/) for API documentation. Swagger provides interactive documentation that makes it easy to explore and test your API endpoints.

#### Setup Swagger

```bash
# Install Swag CLI tool
go install github.com/swaggo/swag/cmd/swag@latest

# Install Swagger UI handler for your Go HTTP server
go get -u github.com/swaggo/http-swagger
```

#### Generate Swagger Documentation

```bash
# Generate Swagger documentation
swag init -d cmd/api
```

The `-d` flag specifies the directory containing your `main.go` file, docs folder is generated

#### Configure Your API Server

Import and configure the Swagger UI in your main.go:

```go
import (
    httpSwagger "github.com/swaggo/http-swagger"
)

docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.apiURL)
r.Get("/swagger/*", httpSwagger.Handler(
	httpSwagger.URL(docsURL),
))
```

#### Access Swagger UI

Once your server is running, access the Swagger UI at:

```
http://localhost:8080/v1/swagger/index.html
```

### Performance Testing with Autocannon

[Autocannon](https://github.com/mcollina/autocannon) is a fast HTTP/1.1 benchmarking tool written in Node.js. You can use it to load test your API endpoints. `npx` is part of npm (Node Package Manager) and allows you to run Node.js packages without having to install them globally or in your project. Ensure you have Node.js and npm installed to use it.

Here's an example command to test the `/v1/users/106` endpoint:

```bash
npx autocannon http://localhost:8080/v1/users/106 --connections 10 --duration 5 -h "Authorization: Bearer asas"
```

Let's break down this command:
- `npx autocannon`: Executes autocannon.
- `http://localhost:8080/v1/users/106`: The URL to test.
- `--connections 10` (`-c 10`): The number of concurrent connections to use.
- `--duration 5` (`-d 5`): The duration of the test in seconds.
- `-h "Authorization: Bearer asas"`: Sets an HTTP header. In this case, it's an `Authorization` header for JWT authentication. Replace `asas` with a valid token.

You can customize the URL, number of connections, duration, headers, and other parameters as needed. Refer to the [Autocannon documentation](https://github.com/mcollina/autocannon) for more options.

### Testing rate limiter

```bash
npx autocannon -r 4000 -d 2 -c 10 --renderStatusCodes http://localhost:8080/v1/health
```