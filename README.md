# Task Management Service

This Go application provides an API endpoint for bulk creation of tasks. It connects to a PostgreSQL database and performs bulk inserts of task records.

## Getting Started

### Prerequisites

- Go 1.17 or higher
- PostgreSQL
- `github.com/lib/pq` and `github.com/joho/godotenv` dependencies

### Installation

1. **Clone the repository:**

```sh
  git clone https://github.com/mmubeenalikhan/go-task-service.git
  cd go-task-service
```

2. **Install Go dependencies:**

```sh
    go mod tidy
```

3. **Create and configure the .env file:**

Copy the sample environment file to .env:

```sh
    cp .env.sample .env
```

Update .env with your PostgreSQL connection details.

4. **Commands to run the project:**

```sh
     go run main.go
```
