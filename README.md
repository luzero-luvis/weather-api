# Weather API

A production-ready weather API service built with Go that fetches weather data from a third-party weather service and provides it through a clean REST API interface.

## Features

### Core Functionality
- **Weather Data Retrieval**: Fetch current weather information for any city
- **RESTful API**: Clean, versioned API endpoints following REST principles
- **Health Check**: Built-in health monitoring endpoint for container orchestration

### Logging & Monitoring
- **Structured Logging**: Uses Go's `slog` package for structured, searchable logs
- **Environment-Based Logging**:
  - **Development**: Text format with DEBUG level for detailed debugging
  - **Production**: JSON format with INFO level for efficient log aggregation
- **Request Logging Middleware**: Automatically logs all incoming requests with:
  - HTTP method, path, and status code
  - Request duration in milliseconds
  - Client IP address and user agent
  - Source file location for debugging

### Production Ready
- **Docker Support**: Multi-stage Dockerfile for optimized image size
- **Security**: Runs as non-root user in container
- **CI/CD Pipeline**: Automated Docker image building and deployment via GitHub Actions
- **Environment Configuration**: Environment-based configuration with `.env` support

## API Endpoints

### Health Check
```
GET /healthz
```
Returns `200 OK` with "ok" message when service is healthy.

### Get Weather Data
```
GET /api/v1/weather?city={city_name}
```

**Query Parameters:**
- `city` (required): Name of the city to fetch weather for

**Success Response (200 OK):**
```json
{
  "resolvedAddress": "London, England, United Kingdom",
  "description": "Clear conditions throughout the day.",
  "timezone": "Europe/London",
  "days": [
    {
      "temp": 15.5,
      "windspeed": 12.3,
      "humidity": 65.0,
      "Conditions": "Clear"
    }
  ]
}
```

**Error Responses:**
- `400 Bad Request`: Missing city parameter
- `500 Internal Server Error`: Failed to fetch weather data or API error

## Architecture

### Project Structure
```
weather-api/
├── main.go                          # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go               # Environment configuration
│   ├── client/
│   │   └── weather-api.go          # Weather API client
│   ├── logger/
│   │   └── logger.go               # Structured logging setup
│   └── middleware/
│       └── logger.go               # HTTP request logging middleware
├── Dockerfile                       # Multi-stage Docker build
└── .github/workflows/
    └── docker-build-push.yml       # CI/CD pipeline
```

### Tech Stack
- **Language**: Go 1.25
- **Router**: [Chi](https://github.com/go-chi/chi) - Lightweight, idiomatic HTTP router
- **Logging**: `log/slog` - Go's native structured logging
- **Configuration**: [godotenv](https://github.com/joho/godotenv) - Environment variable management
- **Container**: Docker with Alpine Linux base

## Setup & Installation

### Prerequisites
- Go 1.25 or higher
- Docker (optional, for containerized deployment)
- Weather API key from your provider

### Local Development

1. **Clone the repository**
```bash
git clone <repository-url>
cd weather-api
```

2. **Install dependencies**
```bash
go mod download
```

3. **Configure environment variables**

Create a `.env` file in the project root:
```env
API_KEY=your_weather_api_key_here
PORT=8000
BASE_URL=https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline
ENV=development
```

Environment variables:
- `API_KEY`: Your weather service API key (required)
- `PORT`: Port number for the server (default: 8000)
- `BASE_URL`: Weather API base URL (required)
- `ENV`: Environment mode - `development` or `production` (affects logging)

4. **Run the application**
```bash
go run main.go
```

The server will start on the configured port (default: 8000).

### Docker Deployment

1. **Build the Docker image**
```bash
docker build -t weather-api:latest .
```

2. **Run the container**
```bash
docker run -d \
  -p 8000:8000 \
  -e API_KEY=your_api_key \
  -e PORT=8000 \
  -e BASE_URL=https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline \
  -e ENV=production \
  --name weather-api \
  weather-api:latest
```

## CI/CD

The project includes a GitHub Actions workflow that automatically:
- Builds the Docker image on push to `main` branch
- Pushes the image to Docker Hub with the `dev` tag
- Generates SBOM (Software Bill of Materials) and provenance attestations

### Required GitHub Secrets
- `DOCKER_USERNAME`: Docker Hub username
- `DOCKER_PASSWORD`: Docker Hub access token

## Logging Examples

### Development Mode (ENV=development)
```
time=2026-02-15T10:30:45.123Z level=INFO msg="starting weather api server" port=8000 env=development
time=2026-02-15T10:30:50.456Z level=INFO msg="incoming request" method=GET path=/api/v1/weather remote_addr=127.0.0.1:54321 user_agent=curl/7.81.0
time=2026-02-15T10:30:50.457Z level=DEBUG msg="calling weather api" city=London
time=2026-02-15T10:30:50.789Z level=DEBUG msg="weather api responded successfully" city=London status_code=200
time=2026-02-15T10:30:50.790Z level=INFO msg="weather data fetched successfully" city=London resolved_address="London, England, United Kingdom" days_count=7
time=2026-02-15T10:30:50.791Z level=INFO msg="request completed" method=GET path=/api/v1/weather status=200 duration_ms=335
```

### Production Mode (ENV=production)
```json
{"time":"2026-02-15T10:30:45.123Z","level":"INFO","msg":"starting weather api server","port":"8000","env":"production"}
{"time":"2026-02-15T10:30:50.456Z","level":"INFO","msg":"incoming request","method":"GET","path":"/api/v1/weather","remote_addr":"127.0.0.1:54321","user_agent":"curl/7.81.0"}
{"time":"2026-02-15T10:30:50.790Z","level":"INFO","msg":"weather data fetched successfully","city":"London","resolved_address":"London, England, United Kingdom","days_count":7}
{"time":"2026-02-15T10:30:50.791Z","level":"INFO","msg":"request completed","method":"GET","path":"/api/v1/weather","status":200,"duration_ms":335}
```

## Usage Examples

### Using cURL
```bash
# Health check
curl http://localhost:8000/healthz

# Get weather for a city
curl http://localhost:8000/api/v1/weather?city=London

# Get weather for a city with spaces
curl "http://localhost:8000/api/v1/weather?city=New%20York"
```

### Using HTTPie
```bash
# Get weather data
http GET localhost:8000/api/v1/weather city==London
```

### Using a browser
```
http://localhost:8000/api/v1/weather?city=London
```

## DON'T BE LIKE ME: DO NOT LOG YOUR API KEY 😂

Seriously, I accidentally logged my API key to the console and had to scramble to remove it. Learn from my mistake -- **NEVER** log sensitive environment variables like API keys. If you see your secrets in your logs, rotate them IMMEDIATELY. Don't be a clown like me. 😂😂😂

## Development Timeline

Here's how the project evolved over time:

1. **Initial Setup** - Basic Go project scaffolding and module setup
2. **HTTP Server & Dockerfile** - Set up a basic HTTP server and created a multi-stage Dockerfile
3. **Chi Router** - Integrated the Chi router for cleaner, idiomatic routing
4. **Health Check & Routes** - Added `/healthz` endpoint and organized API route structure
5. **Weather API Integration** - Built the weather API client to fetch JSON data from VisualCrossing
6. **Docker Security** - Added non-root user (`appuser`) to the Dockerfile
7. **CI/CD Pipeline** - Set up GitHub Actions workflow to auto-build and push Docker images
8. **Structured Logging** - Added `slog`-based logger with environment-aware log levels (debug for dev, info for prod)
9. **Request Logging Middleware** - Built middleware to log every incoming request (method, path, status, duration, IP, user agent)
10. **Environment Configuration** - Restructured env variable handling with `.env` support and proper config loading
11. **Redis Caching** - Added Redis-based caching for weather API responses with a 6-hour TTL to reduce API calls
12. **Code Documentation** - Added comments throughout the codebase for better code understanding
13. **`.dockerignore`** - Added `.dockerignore` to keep Docker builds clean

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Security Features

- **Non-root Container**: Application runs as unprivileged user `appuser` in Docker
- **Input Sanitization**: URL query parameters are properly escaped
- **Error Handling**: Comprehensive error handling with proper logging
- **Secure Defaults**: Production environment uses structured JSON logging for better security monitoring

## License

This project is open source and available under the MIT License.
