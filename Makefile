# Directory containing migration files
MIGRATIONS_DIR=./migrations

# Database connection URL for PostgreSQL
DB_URL=postgres://cl_tn:pa55word@localhost:5432/clean_town?sslmode=disable

# Create a new migration file with up and down SQL files
# Usage: make migrate-new name=create_users_table
migrate-new:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

# Apply all available "up" migrations to bring the database schema to the latest version
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

# Roll back the most recent migration (one step down)
migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down

# Roll back ALL applied migrations (reset database schema to empty state)
migrate-down-all:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down -all

# Force the schema version to a specific version (useful if migrations got stuck or failed)
# Example usage: make migrate-force version=3
migrate-force:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $(version)

# Show the current migration version and dirty state
migrate-status:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version
