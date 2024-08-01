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
- [ ] Add tests
- [x] logger
- [x] graceful shutdown
- [x] create a docker image
- [ ] CI/CD