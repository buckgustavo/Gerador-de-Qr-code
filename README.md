# QR Code Generator

> REST API service for QR code generation with participant data persistence, HTTPS automation via Let's Encrypt, and LGPD-compliant data management.

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=flat&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-336791?style=flat&logo=postgresql&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-green?style=flat)

---

## Overview

This service exposes an HTTP API that generates QR codes on demand, stores participant data (name, e-mail, social networks) in PostgreSQL, and automates TLS certificate provisioning through Let's Encrypt. The codebase follows a layered architecture (`handler → service → repository → domain`) to keep concerns separated and the business logic independently testable.

Key design decisions:

- **Standard library HTTP server** — no external router dependency, keeps the binary lean.
- **`skip2/go-qrcode`** for QR encoding — battle-tested, zero-CGO dependency.
- **`lib/pq`** as the PostgreSQL driver — stable and widely adopted in production Go services.
- **LGPD art. 18 compliance** — a dedicated `DELETE /api/dados/:email` endpoint allows participants to exercise their right to erasure, protected by an `ADMIN_KEY` header.
- **Automatic HTTPS** — when `DOMAIN` is set, the server binds to `:443` with an auto-renewed certificate and redirects `:80` to HTTPS; without it, the server runs on plain HTTP for local development.

---

## Architecture

```
.
├── cmd/                  # Application entrypoints (main packages)
├── internal/
│   ├── domain/           # Entities and business rules (no external deps)
│   ├── service/          # Use-case orchestration
│   ├── repository/       # PostgreSQL data access layer
│   └── handler/          # HTTP handlers and routing
├── web/                  # Static assets / frontend (if applicable)
├── .env.example          # Reference for required environment variables
├── go.mod
└── go.sum
```

---

## Requirements

- Go 1.25+
- PostgreSQL 14+ (with SSL enabled for production)
- A publicly accessible domain for automatic HTTPS (optional for local dev)

---

## Getting Started

**1. Clone and install dependencies**

```bash
git clone https://github.com/buckgustavo/Gerador-de-Qr-code.git
cd Gerador-de-Qr-code
go mod download
```

**2. Configure environment**

```bash
cp .env.example .env
# Edit .env with your credentials
```

| Variable       | Required | Description |
|----------------|----------|-------------|
| `DATABASE_URL` | ✅       | PostgreSQL connection string (`postgresql://user:pass@host/db?sslmode=require`) |
| `PORT`         | ❌       | HTTP port (default: `3000`) |
| `DOMAIN`       | ❌       | Public domain for automatic HTTPS via Let's Encrypt |
| `ADMIN_KEY`    | ❌       | Secret key for the LGPD data-deletion endpoint |

**3. Run**

```bash
go run ./cmd/...
```

For production, build a static binary:

```bash
CGO_ENABLED=0 GOOS=linux go build -o qrcode-server ./cmd/...
./qrcode-server
```

---

## API Reference

### Generate QR Code

```
POST /api/qrcode
Content-Type: application/json

{
  "content": "https://example.com",
  "name":    "Gustavo Buck",
  "email":   "gustavo@example.com"
}
```

Returns a PNG image of the generated QR code.

---

### Delete Participant Data (LGPD art. 18)

```
DELETE /api/dados/:email
X-Admin-Key: <ADMIN_KEY>
```

Permanently removes all data associated with the given e-mail address.  
This endpoint is **disabled** if `ADMIN_KEY` is not set.

---

## Privacy & Compliance

Participant data (name, e-mail, social networks) is stored under the obligations of **LGPD (Lei nº 13.709/2018)**. Operators must:

- Never commit the `.env` file — it contains database credentials and secret keys.
- Generate `ADMIN_KEY` with a cryptographically secure source: `openssl rand -hex 32`.
- Use `sslmode=require&channel_binding=require` in `DATABASE_URL` for production.

---

## Contributing

1. Fork the repository and create a feature branch (`git checkout -b feat/my-feature`).
2. Commit following [Conventional Commits](https://www.conventionalcommits.org/).
3. Open a Pull Request with a clear description of the change and its motivation.

---

## License

Distributed under the MIT License. See [`LICENSE`](LICENSE) for details.
