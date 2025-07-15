# Gin Web Service Template

## Overview

This project is a template for building web services using the Gin framework in Go. It provides a structured and scalable foundation for developing robust APIs.

## Benefits

- **Modular Structure**: Organized codebase for easy maintenance and scalability.
- **RESTful API**: Built-in support for creating RESTful endpoints.
- **Database Integration**: Pre-configured database layer for efficient data management.
- **Caching**: Integrated caching mechanism for improved performance.
- **Error Handling**: Standardized error handling and reporting.
- **Middleware Support**: Easy integration of custom middleware.
- **Logging**: Built-in logging for better debugging and monitoring.
- **Graceful Shutdown**: Ensures proper closure of resources and connections.
- **Docker Support**: Containerization ready for easy deployment.
- **Tracing and Metrics**: Integrated for better observability.

## Configuration

Configuration is loaded via Viper using the `config.Init` helper. Provide the
path, file name and type of your configuration file when creating a server. A
typical `config.yaml` looks like this:

   server:
     port: 8080
     timeout: 10s

   database:
     host: localhost
     port: 5432
     user: youruser
     password: yourpassword
     dbname: yourdbname
     driver: postgres
     maxOpenConns: 5
     maxIdleConns: 2
     connMaxLifetime: 30m

   redis:
     host: localhost
     port: 6379
     password: ""
     db: 0

   log:
     level: info
     format: json

Adjust the values according to your environment and requirements. The
`maxOpenConns`, `maxIdleConns` and `connMaxLifetime` fields control connection
pooling.

## Features

1. **Error Handling**: Centralized error handling with custom error types.
2. **Caching**: Request caching to improve response times.
3. **Database Integration**: Configured database layer with connection pooling.
4. **Middleware**: Custom middleware for various purposes like error handling.
5. **Logging**: Structured logging for better traceability.
6. **Graceful Shutdown**: Proper shutdown procedure to ensure all resources are released.
7. **Docker Support**: Dockerized application for easy deployment and scaling.
8. **Tracing and Metrics**: Integrated tracing and metrics for monitoring and performance analysis.

## Getting Started

1. Clone the repository.
2. Create a configuration file and initialize it using `config.Init`.
3. Build a `main.go` that constructs a server and registers your routes.

```go
cfgOpts := config.ConfigOptions{Path: "./config", Name: "config", Type: "yaml"}
srv, _ := app.NewServer(cfgOpts)
srv.RegisterRoutes(func(d *dependencies.Dependencies) {
    d.Router.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })
})
srv.Run()
```

## Starting the Server

With your `main.go` in place, simply run:

```bash
go run ./...
```

The server listens on the host and port specified in your configuration. Stop
the process with `Ctrl+C` to trigger a graceful shutdown.

For more information on each component see the packages under `pkg/`.

## Swagger API Documentation

Swagger documentation is generated with [swag](https://github.com/swaggo/swag) and served at `/docs/index.html` when the server is running.

To update the documentation:

1. Install the `swag` CLI: `go install github.com/swaggo/swag/cmd/swag@latest`.
2. Run `swag init -g path/to/your/main.go` from the repository root.
3. Start the server and browse to `http://localhost:8080/docs/index.html`.

Each handler and middleware includes Swagger comments so new routes should follow the existing pattern. Middleware attached to a route is documented using the `@Middleware` annotation.

### Adding Middleware

Middleware can be registered globally, for route groups, or for individual
routes. Global middleware is added in `pkg/server.go`, group middleware in each
module's `init.go`, and route middleware when declaring the handler.

When implementing new middleware, document it with a comment starting with `// @Middleware` so Swagger includes the description.

#### Registering Custom Middleware

Client projects can extend the request pipeline by registering their own middleware before starting the server. Middleware functions must follow the Gin handler signature:

```go
func(c *gin.Context)
```

Use the server's `Use` method in `main.go` to add middleware that will run for every route after the built‑in middleware:

```go
srv, _ := app.NewServer(cfg)

srvMw := func(c *gin.Context) {
    log.Printf("path: %s", c.Request.URL.Path)
    c.Next()
}

srv.Use(srvMw)
srv.RegisterRoutes(registerRoutes)
srv.Run()
```

Core middleware for tracing, context propagation, error handling and logging always runs first and cannot be replaced. Custom middleware executes next, followed by any route‑specific middleware configured within modules.

## Structure

```
├── pkg                     // Library packages
│   ├── apiErrors           // API error helpers
│   ├── appTracer           // OpenTelemetry integration
│   ├── cache               // Redis client
│   ├── clientContext       // request/response context
│   ├── config              // configuration loader
│   ├── db                  // SQL database connector
│   ├── dependencies        // lightweight DI container
│   ├── datastore           // data store abstraction
│   ├── middleware          // shared middleware
│   └── server.go           // server bootstrap
├── scripts                 // helper scripts
├── testUtils               // testing helpers
├── go.mod
├── go.sum
└── version.txt             // application version

```

## TODOs

- [x] Add swagger
- [x] add versioning
- [ ] Add tests
- [x] logger
- [x] graceful shutdown
- [x] create a docker image
- [x] add tracing and metrics
- [ ] CI/CD
