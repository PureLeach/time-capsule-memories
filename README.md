# Time Capsule Memories

Time Capsule Memories is a personal pet-project designed to store and send your memories into the future. Users can create “time capsules” containing messages and photos that will be delivered to the recipients of their choice at a specified date and time.

---

## Quick Start

1. Copy the example environment file:

   ```bash
   cp example.env .env
   ```
2. Build the Docker containers:

   ```bash
   docker-compose build
   ```
3. Start the app:

   ```bash
   docker-compose up
   ```

## Usage

- frontend - http://frontend.localhost
- backend - http://backend.localhost/swagger/
- minio - http://minio.localhost
- pgadmin - http://pgadmin.localhost

## Working with migrations (local)

This project uses the [Goose](https://github.com/pressly/goose) library to manage database migrations.

### **Installing Goose**

1. Install the CLI:

   ```bash
   go install github.com/pressly/goose/v3/cmd/goose@latest
   ```
2. (Optional) Install the library if you plan to use Goose in your Go application:

   ```bash
   go get -u github.com/pressly/goose/v3
   ```

### **Creating a migration**

1. Initialize the migrations directory:

   ```bash
   mkdir -p migrations
   ```
2. Create a new migration:

   ```bash
   goose create create_capsules_table sql -dir ./migrations
   ```

   This will generate a file like `20241124120000_create_capsules_table.sql` in the `migrations` folder.
3. Add SQL to the migration file

### **Applying migrations**

1. Apply all pending migrations:

   ```bash
   make migrate_up
   ```
2. Rollback the last migration:

   ```bash
   make migrate_down
   ```
