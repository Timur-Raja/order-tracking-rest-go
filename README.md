# Order Tracking REST API

A RESTful service for managing users, products, and orders, built with Go using GIN framework and Postgres as DB and Elastic search for text search and filetering, containerized with Docker.

## Improvments needed, not added due to time constraints 
- more exhaustive testing suite, failing tests
- independent tests
- avoid n+1 query problem in 2 handlers
- some extra error checking in business logic
- signout handler
- product related endpoints
- general more throught out code architecture improvements

## Project Structure

```
├── .env                   # Environment variables
├── compose.yaml           # Docker Compose for Postgres and test DB setup
├── cmd/
│   └── main.go            # Application entry point
├── api/
│   ├── endpoints.go       # endpoints list
│   └── middleware.go      # middleware functions
├── config/
│   ├── config.go          # env config loading
│   └── errors.go          # definition of eventual errors in config loading
├── db/
│   ├── init.go            # DB connection pool initialization
│   ├── migrate.go         # func to run migrations (golang-migrate) used to run them in test db
│   ├── schema.sql         # Full SQL schema
│   └── migrations/        # Up/down migration files
├── es/
│   └── utils.go          # util function to populate elastic search indexes
├── app/
│   ├── user/              # User business logic and SQL queries
│   ├── order/             # Order business logic and SQL queries
│   └── product/           # Product business logic and SQL queries
├── testing/
│   ├── main_test.go       # Test setup, migrations, HTTP server
│   └── endpoints_test.go  # Sequential endpoint tests
├── go.mod
└── go.sum
```

## Features

* User registration and authentication (session-based via cookies)
* Product listing and management
* Order CRUD operations (create, read, update, delete)
* Order full text search and date filtering through Elasticsearch
* Robust error handling and logging middleware
* Database migrations powered by `golang-migrate`
* Docker Compose for spinning up development and test databases
* integration test suite using `httptest`

## Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/timur-raja/order-tracking-rest-go.git
   cd order-tracking-rest-go
   ```
2. Start PostgreSQL and initialize the test database:

   ```bash
   docker-compose up -d
   ```

   (optionally setup your own .env variables)

## Running the Application

Start the server:

```bash
go run cmd/main.go
```

Access the API at `http://localhost:8080`.

## Database Migrations

All migration files live under `db/migrations`:

* `0001_init_db_tables.up.sql` / `.down.sql`
* `0002_create_views.up.sql` / `.down.sql`

Migrations are automatically applied in tests.
For prod can just run:

* `migrate -path db/migrations -database "${DATABASE_URL}" up`

## Testing

Integration and endpoint tests are located in the `testing/` directory. They:

1. Run migrations against a test postgres db
2. connect to a test Elasticsearch service
3. Spin up an `httptest` server
4. Execute endpoints sequentially and do a simple response check
5. Return the DB to the original state (only migrations applied)

To run tests:

```bash
go test ./testing -v
```

## API Endpoints

### Public Endpoints

* `POST /signup` — Register a new user
* `POST /signin` — Authenticate and receive a session cookie

### Protected Endpoints (require `session_token` cookie)

* `POST /orders` — Create a new order
* `GET /orders` — List your orders # includes Elasticsearch powered text search and date filtering
* `GET /orders/:order_id` — Retrieve a specific order
* `PATCH /orders/:order_id` — Update an existing order # includes editing and cancellation of orders

