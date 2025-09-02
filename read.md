# EV Station Finder and Booking API Documentation

## Project Overview
This backend RESTful API service helps electric vehicle users locate nearby charging stations, view details, and book charging slots efficiently. It is built using Golang with Gin framework, GORM ORM, PostgreSQL database, Redis caching, and JWT-based authentication.

## Running the Project
1. Ensure you have Go 1.24+, PostgreSQL, and Redis installed and running.
2. Set environment variables in a `.env` file, including `DATABASE_URL`, `REDIS_ADDR`, `REDIS_PASSWORD`, and optionally `PORT`.
3. Run the project with:
   ```
   go run cmd/main.go
   ```
4. The server will start on the port specified in `PORT` or default to `8080`.
5. The project uses GORM auto migrations to create/update database tables automatically on startup.

## API Endpoints

### Authentication
- `POST /api/auth/register`  
  Register a new user.  
  Request body: JSON with user details (name, email, password).

- `POST /api/auth/login`  
  Login user and receive access and refresh tokens.  
  Request body: JSON with email and password.

### Stations
- `GET /api/stations`  
  Retrieve all EV charging stations.

- `GET /api/stations/:id`  
  Retrieve details of a specific station by ID.

- `POST /api/stations`  
  Create a new charging station.  
  Request body: JSON with station details.

### Bookings
- `POST /api/bookings`  
  Create a new booking for a charging slot.  
  Request body: JSON with booking details (station ID, user ID, start and end times).

- `GET /api/bookings/user?user_id={user_id}`  
  Retrieve bookings for a specific user.

### Health Checks
- `GET /health`  
  Returns the health status of the service and its dependencies.

- `GET /readiness`  
  Returns readiness status indicating if the service is ready to accept traffic.

### Metrics
- `GET /metrics/cache`  
  Retrieve cache hit/miss metrics.

- `POST /metrics/cache/reset`  
  Reset cache metrics.

## Notes
- The service uses JWT for secure authentication.
- Passwords are hashed using Bcrypt.
- Rate limiting middleware is applied globally.
- Database migrations are handled automatically on startup using GORM's AutoMigrate.
