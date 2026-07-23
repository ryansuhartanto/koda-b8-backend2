# GoREST

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go)
![Gin](https://img.shields.io/badge/Gin-Web_Framework-008ECF?logo=gin)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-4169E1?logo=postgresql)
![JWT](https://img.shields.io/badge/JWT-Authentication-black)
![OpenAPI](https://img.shields.io/badge/OpenAPI-Docs-85EA2D?logo=swagger)
![Vite](https://img.shields.io/badge/Vite-Frontend-9135FF?logo=vite)
![React](https://img.shields.io/badge/React-19-61DAFB?logo=react)
![MIT](https://img.shields.io/badge/License-MIT-yellow)

Go and React exercise implementing a full stack REST server.

![Screenshot 1](docs/1.png)
![Screenshot 2](docs/2.png)

---

## Tech Stack

| Technology         | Description                  |
| ------------------ | ---------------------------- |
| Go                 | Backend language             |
| Gin                | HTTP web framework           |
| PostgreSQL         | Database                     |
| pgx/v5             | PostgreSQL driver            |
| golang-migrate     | Schema migrations            |
| JWT                | Authentication               |
| swag / gin-openapi | API documentation            |
| React + Vite       | Frontend                     |
| bun                | JS package manager / runtime |
| mise               | Tool versions, env, tasks    |

---

## Prerequisites

- [mise](https://mise.jdx.dev) installs the pinned Go toolchain and the `migrate` CLI, no manual `go install` needed
- PostgreSQL running locally
- bun (also managed via `devEngines` in `package.json`, or install separately)

---

## Running

Install tool versions (Go, `migrate`, `vp`) pinned in `.config/mise.toml`:

```bash
mise install
```

Copy the env file and fill in your Postgres credentials:

```bash
cp .env.example .env
```

Install JS dependencies:

```bash
bun install
```

Run backend + frontend together:

```bash
mise run dev
```

This runs `dev:be` (`air`, live-reloading Go server) and `dev:fe` (`bun run dev`, Vite) in parallel. Migrations run automatically on backend startup (see `main.go`).

| Stack    | URL                     |
| -------- | ----------------------- |
| Frontend | <http://localhost:5143> |
| Backend  | <http://localhost:8080> |

Other useful tasks:

| Task                        | Description                           |
| --------------------------- | ------------------------------------- |
| `mise run dev`              | Run backend + frontend together       |
| `mise run dev:be`           | Run backend only (`air`)              |
| `mise run dev:fe`           | Run frontend only (`bun run dev`)     |
| `mise run docs`             | Regenerate OpenAPI docs (`swag init`) |
| `mise run db:create <name>` | Create a new migration                |
| `mise run db:up`            | Apply migrations                      |
| `mise run db:down`          | Roll back migrations                  |

---

## API Documentation

Interactive OpenAPI UI is served by the app itself:

```text
http://localhost:8080/docs
```

---

## API Endpoints

### Auth

| Method | Endpoint         | Description       |
| ------ | ---------------- | ----------------- |
| POST   | `/auth/register` | Register new user |
| POST   | `/auth/login`    | Login user        |

### Users (requires auth)

| Method | Endpoint             | Description            |
| ------ | -------------------- | ---------------------- |
| GET    | `/users/`            | List users             |
| PATCH  | `/users/:id`         | Update user            |
| PUT    | `/users/:id/picture` | Upload profile picture |
| DELETE | `/users/:id`         | Delete user            |

---

## ERD

```mermaid
---
title: GoREST
---
erDiagram

users ||--o| profiles : "has"

users {
    bigint id PK

    timestamptz created_at
    timestamptz updated_at

    string email    UK
    string password
}

profiles {
    bigint id PK,FK

    timestamptz created_at
    timestamptz updated_at

    string  name
    string? picture_url
}
```

## License

[MIT](LICENSE)
