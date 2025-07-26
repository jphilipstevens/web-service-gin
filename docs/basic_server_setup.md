# Basic Server Setup

This guide shows how to create a minimal HTTP server using the packages in this repository. The example demonstrates registering a couple of routes and starting the server.

## 1. Create a configuration file

Save a YAML file with the following fields. Adjust values as needed for your environment.

```yaml
app_name: demo
server:
  host: 0.0.0.0
  port: 8080
redis:
  host: localhost
  port: 6379
  password: ""
  db: 0
database:
  host: localhost
  port: 5432
  user: demo
  password: demo
  driver: postgres
  dbname: demo
  maxOpenConns: 5
  maxIdleConns: 2
  connMaxLifetime: 30m
uptrace:
  dsn: ""
  endpoint: ""
```

Save the file as `config.yaml` in a directory accessible to your application.

Any configuration value can be overridden using environment variables. Replace
dots in the key with underscores and convert it to uppercase. For example,
`REDIS_HOST=cache.example.com` overrides `redis.host`.

## 2. Implement `main.go`

```go
package main

import (
    "log"

    "github.com/gin-gonic/gin"
    "github.com/jphilipstevens/web-service-gin/v2/pkg/config"
    "github.com/jphilipstevens/web-service-gin/v2/pkg/dependencies"
    "github.com/jphilipstevens/web-service-gin/v2/pkg/server"
)

func registerRoutes(deps *dependencies.Dependencies) {
    r := deps.Router
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })

    r.GET("/hello", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "hello"})
    })
}

func main() {
    srv, err := server.NewServer(config.ConfigOptions{Path: ".", Name: "config", Type: "yaml"})
    if err != nil {
        log.Fatalf("failed to create server: %v", err)
    }

    srv.RegisterRoutes(registerRoutes)

    if err := srv.Run(); err != nil {
        log.Fatalf("server exited: %v", err)
    }
}
```

## 3. Run the example

```bash
go run main.go
```

The server listens on the configured port and exposes `/ping` and `/hello` routes.
