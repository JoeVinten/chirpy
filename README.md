# Chirpy

A learning project: building a Twitter-like RESTful API server in Go by following the [boot.dev](https://boot.dev) Web Servers course.

## About This Project

This is an **educational project** built to learn backend web development, HTTP servers, and Go fundamentals. I'm working through boot.dev's "Learn Web Servers" course, which teaches backend concepts by building a real API from scratch rather than relying on frameworks.

Chirpy is a social media API that allows users to create accounts, post short messages ("chirps"), and manage their content - similar to a simplified Twitter clone.

## Why I'm Building This

### Learning Go

I'm learning Go to strengthen my backend development skills:

- **Expanding beyond frontend/Node.js**: Coming from a frontend and Node.js background, I wanted to learn a statically-typed, compiled language that's built for backend systems
- **Understanding the fundamentals**: Go's minimalist approach forces you to understand HTTP, databases, and authentication from first principles rather than hiding complexity behind abstractions

## What I'm Learning

Through building Chirpy, I'm gaining hands-on experience with:

- **HTTP fundamentals**: Request/response cycle, status codes, headers, and REST principles
- **Go language features**: Error handling, interfaces, structs, slices, and Go idioms
- **Database operations**: Schema design, migrations, queries, and working with PostgreSQL
- **Authentication & security**: JWT tokens, refresh tokens, password hashing with bcrypt, and authorization patterns
- **API design**: RESTful endpoints, middleware patterns, and proper error responses
- **Testing in Go**: Table-driven tests and the testing package
- **Tools & ecosystem**: sqlc for type-safe queries, Goose for migrations, and the Go toolchain

## Features Implemented

- ✅ User registration and login with JWT authentication
- ✅ Create, read, and delete chirps (posts)
- ✅ Token refresh and revocation
- ✅ Authorization (users can only delete their own chirps)
- ✅ Filtering chirps by author
- ✅ Sorting chirps by date
- ✅ Middleware for authentication
- ✅ Password hashing and validation
- ✅ PostgreSQL database with migrations

## Tech Stack

- **Go 1.21+** - Backend language
- **PostgreSQL** - Database
- **sqlc** - Type-safe SQL query generation
- **Goose** - Database migrations
- **JWT** - Authentication tokens
- **argon2id** - Password hashing

## Project Structure
```
.
├── main.go                 # Server setup and routing
├── handler_*.go           # HTTP handlers for each endpoint
├── internal/
│   ├── auth/              # Authentication utilities (JWT, argon2id, API keys)
│   └── database/          # sqlc generated code
├── sql/
│   ├── schema/            # Database migrations
│   └── queries/           # SQL queries for sqlc
└── README.md
```

## API Endpoints

### Users
- `POST /api/users` - Create a new user
- `POST /api/login` - Login and receive JWT + refresh token
- `PUT /api/users` - Update user email/password (authenticated)

### Chirps
- `POST /api/chirps` - Create a chirp (authenticated)
- `GET /api/chirps` - Get all chirps (optional `?author_id=<uuid>` and `?sort=desc`)
- `GET /api/chirps/{chirpID}` - Get a specific chirp
- `DELETE /api/chirps/{chirpID}` - Delete your chirp (authenticated)

### Auth
- `POST /api/refresh` - Refresh access token using refresh token
- `POST /api/revoke` - Revoke refresh token (logout)

## Running Locally
```bash
# Install dependencies
go mod download

# Set up database
createdb chirpy

# Run migrations
goose -dir sql/schema postgres "postgresql://localhost/chirpy?sslmode=disable" up

# Generate sqlc code
sqlc generate

# Build and run
go build -o chirpy && ./chirpy
```

## Key Learnings So Far

- **Go's explicit error handling** makes you think about failure cases upfront
- **Middleware patterns** for authentication are cleaner than repeating auth logic
- **Context values** are the idiomatic way to pass request-scoped data in Go
- **Table-driven tests** make it easy to cover multiple scenarios
- **sqlc** generates type-safe database code, catching errors at compile time
- **JWT tokens** need both access tokens (short-lived) and refresh tokens (long-lived)

## Disclaimer
**This is a learning project** - built by following a structured course to understand backend development in Go. Not production-ready, but a real working API built from first principles!
