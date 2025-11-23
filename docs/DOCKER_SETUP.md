# Docker Compose Setup

This Docker Compose configuration provides a PostgreSQL database for local development.

## Quick Start

```bash
# Start PostgreSQL
docker-compose up -d

# View logs
docker-compose logs -f postgres

# Check status
docker-compose ps

# Stop services
docker-compose down

# Stop and remove volumes (delete all data)
docker-compose down -v
```

## Database Configuration

**Connection Details:**

- Host: `localhost`
- Port: `5432`
- Database: `go_bootstrap`
- Username: `postgres`
- Password: `postgres`

**Connection String (DSN):**

```text
postgres://postgres:postgres@localhost:5432/go_bootstrap?sslmode=disable
```

## Migrations

Migrations are automatically executed when the container starts for the first time.

**How it works:**

- All `.sql` files in the `migrations/` directory are mounted to `/docker-entrypoint-initdb.d/` in the container
- PostgreSQL automatically runs all scripts in this directory in alphabetical order on first startup
- Scripts are only executed on initial database creation (when volume is empty)

**To re-run migrations:**

```bash
# Remove the volume and restart
docker-compose down -v
docker-compose up -d
```

## Accessing PostgreSQL

### Using psql (inside container)

```bash
# Connect to database
docker-compose exec postgres psql -U postgres -d go_bootstrap

# List tables
\dt

# Describe table
\d users

# Exit
\q
```

### Using psql (from host)

```bash
psql -h localhost -p 5432 -U postgres -d go_bootstrap
# Password: postgres
```

### Using pgAdmin or DataGrip

**Connection settings:**

- Host: `localhost`
- Port: `5432`
- Database: `go_bootstrap`
- Username: `postgres`
- Password: `postgres`

## Volume Management

Database data is persisted in a Docker volume named `postgres_data`.

```bash
# List volumes
docker volume ls | grep postgres

# Inspect volume
docker volume inspect go-bootstrap_postgres_data

# Remove volume (WARNING: deletes all data)
docker volume rm go-bootstrap_postgres_data
```

## Troubleshooting

### Port already in use

If port 5432 is already in use:

```bash
# Check what's using the port
lsof -i :5432

# Option 1: Stop local PostgreSQL
brew services stop postgresql@16

# Option 2: Change port in docker-compose.yml
ports:
  - "5433:5432"  # Use 5433 on host instead
```

### Container won't start

```bash
# View detailed logs
docker-compose logs postgres

# Remove everything and start fresh
docker-compose down -v
docker-compose up -d
```

### Migrations not running

Migrations only run on first container startup. If you need to re-run:

```bash
# Remove volume and restart
docker-compose down -v
docker-compose up -d

# Or connect and run manually
docker-compose exec postgres psql -U postgres -d go_bootstrap -f /docker-entrypoint-initdb.d/001_create_auth_and_user_tables.sql
```

## Production Considerations

**DO NOT use this configuration for production!**

This setup is for local development only. For production:

- Use strong, unique passwords
- Enable SSL/TLS connections
- Configure proper backup strategies
- Use managed database services (RDS, Cloud SQL, etc.)
- Implement proper security groups/firewall rules
- Use secrets management (not plain text passwords)
