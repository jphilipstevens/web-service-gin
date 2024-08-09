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

The application is configured using a YAML file located at `config/config.yaml`. Here's how to set it up:

1. Copy the `config.yaml.example` to `config.yaml`.
2. Edit the `config.yaml` file to set your specific configuration:

   server:
     port: 8080
     timeout: 10s

   database:
     host: localhost
     port: 5432
     user: youruser
     password: yourpassword
     dbname: yourdbname

   cache:
     address: localhost:6379
     password: ""
     db: 0

   log:
     level: info
     format: json

   Adjust the values according to your environment and requirements.

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
2. Configure the `config.yaml` file
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

## Structure

```
├── app                         // Our application and all dependent code
│   ├── albums                  // Our Albums domain, including all APIs, services, and models
│   │   ├── controller.go       // API controller for the Album domain
│   │   ├── service.go          // service layer for all business logic
│   │   ├── repository.go       // repository layer for all data access to an album
│   │   ├── models.go           // Models for presenting an Album
│   │   └── init.go             // the bootstrapping of the entire api, including routes, and versioning
│   ├── apiErrors                  
│   │   └── error.go            // API Error creation and model definition
│   ├── cache                  
│   │   └── cache.go            // request caching layer, model and initialization
│   ├── db                      // database layer module
│   │   ├── db.go               // database connection and initialization
│   │   ├── error.go            // error mapping from db specific to application error
│   │   └── models.go           // shared models from the database E.G. pagination models
│   ├── middleware              // middleware used for the application
│   │   └── errorHandler.go     // error handling code to return standardized error models
├── config
│   └── config.yaml             // yaml file for all configuration
├── seed                        // seed data for the application locally
│   ├── albums.go               // Albums seed data
│   └── seed.go                 // main seed script for all models
└── main.go
└── go.sum                      // Go module checksum file

```

## TODOs

- [ ] Add swagger
- [ ] add versioning
- [ ] Add tests
- [x] logger
- [x] graceful shutdown
- [x] create a docker image
- [x] add tracing and metrics
- [ ] CI/CD
