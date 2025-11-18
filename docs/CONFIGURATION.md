# Configuration Guide

This guide explains how configuration works in this application, including the structure, hot-reload capabilities, and best practices.

## Table of Contents

- [Configuration Structure](#configuration-structure)
- [Application Configs](#application-configs)
- [Pprof Configuration (Realtime Hot-Reload)](#pprof-configuration-realtime-hot-reload)
- [How Configuration Works Internally](#how-configuration-works-internally)

## Configuration Structure

The application uses JSON-based configuration files (`env.json`) that can be hot-reloaded during development. The configuration is separated by application type for better modularity and independence.

**Config file location:** `env.json` (create from `env.example.json`)

**Important:** After cloning the repository, run `make generate` to generate all required code before starting the application.

## Application Configs

Each application (REST API, gRPC API, Scheduler) has its own independent configuration:

### REST API Configuration

**`app_rest_api`:**

```json
{
    "app_rest_api": {
        "name": "directory-service-rest-api",   // Service name for logging/tracing
        "env": "development",                    // Environment: development/staging/production
        "debug_mode": true,                      // Enable debug logging and SQL query logging
        "port": 8080,                            // HTTP server port
        "pprof": {                               // Pprof configuration (nested)
            "enable": true,                      // Enable/disable pprof server
            "port": 8080,                        // Pprof HTTP server port
            "static_token": "your-secret-token"  // Static token for authentication
        }
    }
}
```

### gRPC API Configuration

**`app_grpc_api`:**

```json
{
    "app_grpc_api": {
        "name": "directory-service-grpc-api",   // Service name for logging/tracing
        "env": "development",                    // Environment: development/staging/production
        "debug_mode": true,                      // Enable debug logging and SQL query logging
        "port": 9090,                            // gRPC server port
        "pprof": {                               // Pprof configuration (nested)
            "enable": false,                     // Enable/disable pprof server
            "port": 7070,                        // Pprof HTTP server port
            "static_token": "your-secret-token"  // Static token for authentication
        }
    }
}
```

### Scheduler Configuration

**`app_scheduler` - Background Jobs:**

```json
{
    "app_scheduler": {
        "name": "directory-service-scheduler",     // Service name for logging/tracing
        "env": "development",                       // Environment: development/staging/production
        "debug_mode": true,                         // Enable debug logging
        "healthcheck_interval": "0 */5 * * * *",   // Cron expression (every 5 minutes)
        "pprof": {                                  // Pprof configuration (nested)
            "enable": true,                         // Enable/disable pprof server
            "port": 7070,                           // Pprof HTTP server port
            "static_token": "your-secret-token"     // Static token for authentication
        }
    }
}
```

### Why Separate Configs?

- Each service can have different names for logging/monitoring
- Independent debug modes per service
- Different environment settings (e.g., HTTP in production, gRPC in staging)
- Each service has its own pprof configuration (can enable/disable independently)
- Better multi-service architecture support

### Database Configuration

Each application has its own database configuration nested within its config:

```json
{
    "app_rest_api": {
        "database": {
            "dsn": "user:password@tcp(localhost:3306)/dbname?parseTime=true",
            "max_open_conns": 25,
            "max_idle_conns": 25,
            "conn_max_lifetime": "300s",
            "conn_max_idle_time": "60s"
        }
    },
    "app_grpc_api": {
        "database": {
            "dsn": "user:password@tcp(localhost:3306)/dbname?parseTime=true",
            "max_open_conns": 25,
            "max_idle_conns": 25,
            "conn_max_lifetime": "300s",
            "conn_max_idle_time": "60s"
        }
    },
    "app_scheduler": {
        "database": {
            "dsn": "user:password@tcp(localhost:3306)/dbname?parseTime=true",
            "max_open_conns": 25,
            "max_idle_conns": 25,
            "conn_max_lifetime": "300s",
            "conn_max_idle_time": "60s"
        }
    }
}
```

**Why Nested Database Config?**

- Each service can connect to different databases if needed
- Independent connection pool settings per service
- Better isolation and flexibility for microservices architecture

## Pprof Configuration (Realtime Hot-Reload)

Each application (REST API, gRPC API, Scheduler) has its own **independent pprof configuration** nested within its config. This allows you to enable/disable profiling per service.

### Pprof Nested Structure

```json
{
    "app_rest_api": {
        "pprof": {
            "enable": true,
            "port": 8080,
            "static_token": "your-secret-token"
        }
    },
    "app_grpc_api": {
        "pprof": {
            "enable": false,
            "port": 7070,
            "static_token": "your-secret-token"
        }
    },
    "app_scheduler": {
        "pprof": {
            "enable": true,
            "port": 6060,
            "static_token": "your-secret-token"
        }
    }
}
```

### Hot-Reload Behavior

The pprof server monitors configuration changes and automatically starts/stops based on the `enable` flag in each app's config:

1. **Enable pprof:**
   - Set `"enable": true` in the respective app's pprof config
   - Server starts automatically on the configured port
   - No application restart required

2. **Disable pprof:**
   - Set `"enable": false` in the respective app's pprof config
   - Server stops gracefully
   - No application restart required

3. **Change pprof port:**
   - **Step 1:** Set `"enable": false` (stops current server)
   - **Step 2:** Change `"port": 8080` to desired port
   - **Step 3:** Set `"enable": true` (starts server on new port)

   **Important:** The restart trigger is based on changes to the `enable` flag. Direct port changes without toggling `enable` will not take effect.

**Accessing Pprof Endpoints:**

All pprof endpoints require authentication via `Authotization` header (note: typo in header name is intentional - matches code):

```bash
# REST API pprof (port 8080)
curl -H "Authotization: your-secret-token" \
  http://localhost:8080/debug/pprof/heap

# gRPC API pprof (port 7070)
curl -H "Authotization: your-secret-token" \
  http://localhost:7070/debug/pprof/goroutine

# Scheduler pprof (port 6060)
curl -H "Authotization: your-secret-token" \
  http://localhost:6060/debug/pprof/

# Available endpoints:
# /debug/pprof/          - Index page
# /debug/pprof/heap     - Memory heap profile
# /debug/pprof/goroutine - Goroutine profile
# /debug/pprof/profile  - CPU profile
# /debug/pprof/trace    - Execution trace
# /debug/pprof/allocs   - Memory allocations
# /debug/pprof/block    - Blocking profile
# /debug/pprof/mutex    - Mutex contention
```

**Security Note:** In production, always use a strong `static_token` and restrict access to pprof endpoints via network policies or authentication middleware.

### Why Nested Pprof Config?

- Each service can enable/disable profiling independently
- Different services can use different pprof ports
- No port conflicts when running multiple services simultaneously
- Better isolation and control per service

## How Configuration Works Internally

**For beginners:** Understanding how configuration flows through the application:

### 1. Configuration Files

**`env.json`:**

- Contains all application settings
- Watched for changes during development
- Automatically reloaded when modified

### 2. Config Loading

**`internal/config/load_config.go`:**

```go
// In cmd layer - Application startup (cmd.go PersistentPreRun)
config.LoadConfig("restapi")  // Initialize with command name
appCfg := config.GetAppRestApi() // Get REST API config
```

**Command Name:** The parameter `"restapi"` (or `"grpcapi"`, `"scheduler"`) tells the config loader which app config to use and watches for changes to its `debug_mode`. When `debug_mode` changes, the config automatically reloads.

### 3. Config Types

**`internal/config/type.go`:**

- Defines struct types for all configurations
- Uses struct tags for JSON mapping: `env:"field_name"`

### 4. Usage in Application

```go
// In cmd layer - Application startup
config.LoadConfig("restapi")  // Set command context
infrastructure.NewLogging(slogHookOption, zerologHookOption)  // Auto-reads config

// In app layer - Feature initialization
db, err := infrastructure.NewDB()  // Auto-reads config based on cmd context
```

### Configuration Flow

```text
env.json → LoadConfig() → root struct → Getter functions → Application
```

### Getter Functions

**App Config Getters:**

- `config.GetAppRestApi()` - Get complete REST API config (includes nested pprof & database)
- `config.GetAppGrpcApi()` - Get complete gRPC API config (includes nested pprof & database)
- `config.GetAppScheduler()` - Get complete Scheduler config (includes nested pprof & database)

**Dynamic Context-Aware Getters:**

These functions automatically return the config for the currently running application based on the `cmdName` set during `LoadConfig()`:

- `config.GetPprof()` - Get pprof config for current app
- `config.GetDatabase()` - Get database config for current app
- `config.GetDebugMode()` - Get debug mode for current app
- `config.GetAppName()` - Get app name for current app
- `config.GetEnv()` - Get environment for current app

**Utility:**

- `config.UnwatchLoader()` - Stop watching config file (called on shutdown)

### Why This Design?

- **Separation of Concerns:** Config loading separated from business logic
- **Type Safety:** Strongly typed configuration access
- **Hot-Reload:** Changes detected automatically in development
- **Testability:** Easy to mock config in tests
- **Independence:** Each app can use different configs
