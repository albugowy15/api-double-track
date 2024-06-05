# Double Track Recommendation REST API

This project is a Go-based REST API that offers a variety of endpoints designed to supply data for integration with the [Double Track Recommendation Web](https://github.com/albugowy15/double-track-recommendations-web).

## Tech Stack

- Go 1.22
- Chi (HTTP Router)
- SQLx (database library)
- Viper (configuration management)
- Air (hot reload for development)
- PostgreSQL (database)

## Docker Compose

To run the project with Docker Compose, follow these steps:

1. Start the PostgreSQL Docker container by running:

```bash
docker-compose up -d db
```

2. Run database migrations and seed data by executing:

```bash
go run cmd/database/migrations/up/main.go && go run cmd/database/seeder/main.go
```

3. Finally, launch the application using Docker Compose:

```bash
docker-compose up api
```

This will start the application in a Docker container, allowing you to interact with the REST API. You can interact with the API by navigating to [http://localhost:8080](http://localhost:8080).

## Development Mode

To run the project in development mode, follow these steps:

1. Make sure you have installed air for hot reload. Refer to [air docs](https://github.com/cosmtrek/air) for installation instructions.
2. Run the following command to start the application with air:

```bash
air
```

This will launch the application with hot reload enabled, allowing you to make changes to the code and see them reflected in real-time without restarting the server.
