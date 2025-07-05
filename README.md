# VASST COMMUNICATION AGENT APP

This is the backend API for VASST Communication agents application, built with Go and Gin framework.

## API Documentation

The API documentation is available through Swagger UI at `/swagger/index.html` when running the server.

### Base URL
- Development: `http://localhost:8080/v1`
- Production: `https://api.vasst.id/v1`

### Authentication

The API uses JWT (JSON Web Token) for authentication. To access protected endpoints:

1. Include the JWT token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

2. The token should be obtained through the login endpoint.

### Health Check

The API provides a health check endpoint at `/health-check` that monitors:
- Database connection
- Third-party service availability

## Development

### Prerequisites

- Go 1.21 or higher
- PostgreSQL
- Make (optional)

### Setup

1. Clone the repository
2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env
```

4. Run the application:
```bash
go run cmd/api/main.go
```

### Testing

Run tests with:
```bash
go test ./...
```

### Building

Build the application with:
```bash
go build -o bin/api cmd/api/main.go
```

## License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details.