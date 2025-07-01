# Makefile

# Load variables from .env file
include .env
export DATABASE_URL

# Apply all up migrations
migrate-up:
	migrate -path db/migrations -database "$(DATABASE_URL)" up

# Roll back the last migration
migrate-down:
	migrate -path db/migrations -database "$(DATABASE_URL)" down

# Drop and recreate schema version (used for troubleshooting)
migrate-force-zero:
	migrate -path db/migrations -database "$(DATABASE_URL)" force 0

# Create a new migration (usage: make migrate-create name=create_tasks_table)
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=create_table_name"; \
		exit 1; \
	fi
	migrate create -ext sql -dir db/migrations -seq $(name)