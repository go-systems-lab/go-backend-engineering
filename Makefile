include .envrc
MIGRATIONS_PATH = ./cmd/migrate/migrations

.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up

.PHONY: migrate-force
migrate-force:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) force $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-version
migrate-version:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) version

.PHONY: migrate-goto
migrate-goto:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) goto $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-drop
migrate-drop:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) drop

.PHONY: seed
seed:
	@DB_ADDR="postgres://admin:adminpassword@localhost/social_go?sslmode=disable" go run cmd/migrate/seed/main.go

.PHONY: gen-docs
gen-docs:
	@echo "Generating Swagger documentation..."
	@swag init --generalInfo cmd/api/main.go --output ./docs --parseDependency --parseInternal --quiet && swag fmt --dir ./docs
	@echo "Swagger documentation generated in ./docs"

.PHONY: test
test:
	@go test -v ./...

%:
	@: