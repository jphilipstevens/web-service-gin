# Basic Server Setup

This guide shows how to create a minimal HTTP server using the packages in this repository. The example demonstrates registering a couple of routes and starting the server.

## 1. Create a configuration file

Save a YAML file with the following fields. Adjust values as needed for your environment.

```yaml
app_name: demo
server:
  host: 0.0.0.0
  port: 8080
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
    "github.com/jphilipstevens/web-service-gin/v2/pkg/server"
)

func registerRoutes(r *gin.Engine) {
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })

    r.GET("/hello", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "hello"})
    })
}


func main() {
    if err := config.Init(config.ConfigOptions{Path: ".", Name: "config", Type: "yaml"}); err != nil {
        log.Fatalf("failed to load config: %v", err)
    }

    srv := server.New(config.Get())

    srv.UseBefore(func(c *gin.Context) {
        log.Println("before built-in")
        c.Next()
    })

    srv.UseAfter(func(c *gin.Context) {
        log.Println("before routes")
        c.Next()
    })

    srv.UseFinal(func(c *gin.Context) {
        c.Next()
        log.Println("after route")
    })

    // Middleware runs in the following order:
    //   1. UseBefore middleware
    //   2. built-in middleware
    //   3. UseAfter middleware
    //   4. route handlers
    //   5. UseFinal middleware

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
