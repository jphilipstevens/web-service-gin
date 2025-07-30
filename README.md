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

The application is configured using a YAML file located next to `main.go`. Edit it to set your specific configuration:

   server:
     port: 8080
     timeout: 10s

   uptrace:
     dsn: http://project2_secret_token@localhost:14317/2
   enable_open_telemetry: true

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

   Adjust the values according to your environment and requirements. The `maxOpenConns`, `maxIdleConns`, and `connMaxLifetime` settings control database connection pooling.

You can override any configuration value with an environment variable by
replacing dots in the key with underscores and uppercasing it. For example,
`REDIS_HOST=cache.example.com` overrides `redis.host`. This is useful when
deploying to containerized or cloud environments where settings are provided via
environment variables.

## Features

1. **Album Management**: CRUD operations for managing albums.
2. **Error Handling**: Centralized error handling with custom error types.
3. **Caching**: Request caching to improve response times.
4. **Database Integration**: Configured database layer with connection pooling.
5. **Middleware**: Custom middleware for various purposes like error handling.
6. **Logging**: Structured logging for better traceability.
7. **Graceful Shutdown**: Proper shutdown procedure to ensure all resources are released.
8. **Docker Support**: Dockerized application for easy deployment and scaling.
9. **Tracing and Metrics**: Integrated tracing and metrics for monitoring and performance analysis.

## Getting Started

1. Clone the repository
2. Create a configuration file as shown in [docs/basic_server_setup.md](docs/basic_server_setup.md)
3. Run `go mod tidy` to install dependencies
4. Run `go run main.go` to start the server

## Starting the Server

To start the server, follow these steps:

1. Ensure you have Go installed on your system.
2. Open a terminal and navigate to the project root directory.
3. Run the following command:

   go run main.go

4. You should see output similar to this:

   2023/06/10 15:30:45 Starting server on :8080

5. The server is now running and listening on port 8080 (or the port specified in your `config.yaml`).

You can now send requests to `http://localhost:8080` to interact with the API.

To stop the server, press `Ctrl+C` in the terminal. The application will perform a graceful shutdown, ensuring all resources are properly released.

For more detailed information on each component, please refer to the respective files in the project structure.

## Swagger API Documentation

Swagger documentation is generated with [swag](https://github.com/swaggo/swag) and served at `/docs/index.html` when the server is running.

To update the documentation:

1. Install the `swag` CLI: `go install github.com/swaggo/swag/cmd/swag@latest`.
2. Run `swag init -g main.go` from the repository root.
3. Start the server and browse to `http://localhost:8080/docs/index.html`.

Each handler and middleware includes Swagger comments so new routes should follow the existing pattern. Middleware attached to a route is documented using the `@Middleware` annotation.

### Adding Middleware

Middleware can be registered globally, for route groups, or for individual routes. Global middleware is added in `app/server.go`, group middleware in each module's `init.go`, and route middleware when declaring the handler.

When implementing new middleware, document it with a comment starting with `// @Middleware` so Swagger includes the description.

#### Registering Custom Middleware

Client projects can extend the request pipeline by registering their own middleware before starting the server. Middleware functions must follow the Gin handler signature:

```go
func(c *gin.Context)
```

Use `UseBefore`, `UseAfter`, and `UseFinal` to extend the middleware stack. All
middleware must be registered before calling `RegisterRoutes` so it can be
applied in the correct order. The execution sequence is:

1. `UseBefore` middleware
2. built-in middleware (recovery, tracing, logging, etc.)
3. `UseAfter` middleware
4. route handlers
5. `UseFinal` middleware

`UseFinal` handlers should always call `c.Next()` so the request reaches the
actual route handler before running their cleanup logic:

```go
srv := server.New(cfg)

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

srv.RegisterRoutes(registerRoutes)
srv.Run()
```

Core middleware for tracing, context propagation, error handling and logging always runs first and cannot be replaced. Custom middleware executes next, followed by any route‑specific middleware configured within modules.

### Advanced server control

For specialized scenarios you may want to bootstrap the HTTP server yourself. Use
`ApplyMiddleware()` to attach all middleware, retrieve the router with `Router()`,
and manage shutdown using `GracefulShutdown`:

```go
srv.ApplyMiddleware()
s := &http.Server{Handler: srv.Router()}
go s.ListenAndServe()

// trigger shutdown somehow
server.GracefulShutdown(s, time.Second*5)
```

### Extending beyond HTTP

The HTTP server exposes composable building blocks so you can run additional
protocols alongside Gin. For example, start a gRPC or GraphQL server in the same
process and manage their lifecycles together:

```go
go srv.Run() // existing HTTP service

grpcSrv := grpc.NewServer()
go grpcSrv.Serve(lis)

// shutdown logic
server.GracefulShutdown(httpSrv, time.Second*5)
grpcSrv.GracefulStop()
```

The template keeps middleware, configuration and shutdown utilities decoupled,
making it straightforward to integrate other protocols without altering the core
modules.

## Structure

```
├── app                         // Core application modules
│   ├── apiErrors               // API error helpers
│   ├── cache                   // request caching layer
│   ├── db                      // database connectors
│   ├── middleware              // shared middleware
│   └── server.go               // server bootstrap
├── example
│   ├── config                  // application configuration
│   │   └── config.yaml
│   ├── features
│   │   └── albums              // Albums domain implementation
│   ├── seed                    // seed data
│   └── main.go
├── go.mod
└── go.sum                      // Go module checksum file

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
