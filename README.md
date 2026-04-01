# Go Auth Backend API

Explore and test the API using the interactive documentation:
- 🔗 **Live API Docs (Try it out)**  
[![API Docs](https://img.shields.io/badge/API-Docs-blue)](https://lunar-shuttle-398631.docs.buildwithfern.com/go-auth-backend-api)

Production-ready auth service built with Go · Gin · GORM · PostgreSQL · Redis.
Covers the full identity lifecycle — registration, email verification, login, session management, password reset, and Google OAuth 2.0.
Supports server and Vercel serverless from the same codebase.

---

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Runtime Modes](#runtime-modes)
- [Architecture](#architecture)
- [API Reference](#api-reference)
- [Auth Flows](#auth-flows)
- [Google OAuth Flow](#google-oauth-flow)
- [Password Management](#password-management)
- [Configuration](#configuration)
- [Getting Started](#getting-started)
- [Deployment](#deployment)
- [Security Design](#security-design)
- [Production Checklist](#production-checklist)

---

## Features

| Category | Details |
|---|---|
| **Auth Methods** | Email/password · Google OAuth 2.0 |
| **Sessions** | Rotating refresh tokens · HttpOnly cookies · DB-backed with revocation |
| **Password Security** | Bcrypt hashing · old ≠ new enforcement · OTP-based reset |
| **OTP** | 6-digit · 5-min expiry · 30s cooldown · 3-attempt lockout |
| **Rate Limiting** | Per-IP Redis limiter · 6 req/min on all auth routes |
| **Email Verification** | Token sent on register · login blocked until verified |
| **Account Linking** | Google ↔ password linking with explicit user confirmation |
| **Dual Runtime** | Server (`cmd/api/main.go`) + Vercel serverless (`api/index.go`) |

- **Auth Methods** — email/password and Google OAuth 2.0; Google login auto-creates or links accounts on first use.
- **Sessions** — rotating refresh tokens stored as SHA-256 hashes; each session tracks IP and device info for auditing.
- **Password Security** — bcrypt hashed, change requires current password as proof, new must differ enforced before any write.
- **OTP** — `crypto/rand` generated, bcrypt-hashed in Redis; 30s cooldown and 3-attempt lockout before invalidation.
- **Rate Limiting** — per-IP Redis counter on the entire `/api/v1/auth` group; rejected before reaching any handler.
- **Email Verification** — one-time SHA-256 hashed token with 10-min expiry; login blocked until verified.
- **Account Linking** — Google login on an existing password account returns `ErrLinkRequired`; no silent merge, user must confirm.
- **Dual Runtime** — same codebase, same handler; only the entry point differs between server and Vercel serverless.

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.21+ |
| HTTP Framework | Gin |
| ORM | GORM |
| Database | PostgreSQL |
| Cache / OTP / Rate Limiting | Redis |
| Auth Tokens | JWT (custom `tokenJWT` package) |
| OAuth | Google OAuth 2.0 (`golang.org/x/oauth2`) |
| Email | SMTP (custom mailer package) |
| Serverless | Vercel (`@vercel/go`) |
| Hot Reload | Air |

---

## Project Structure

```
Go-Auth-Backend-Api/
├── api/index.go                      # Vercel serverless entry point
├── app/app.go                        # Shared app setup (sync.Once) — used by both runtimes
├── cmd/
│   ├── api/main.go                   # Server entry point
│   └── migrate/main.go               # Migration runner
├── internal/
│   ├── config/
│   │   ├── authconfig/               # Google OAuth2 config
│   │   ├── env/env.go                # Env var loader (os.Getenv + godotenv fallback)
│   │   ├── mailConfig/               # SMTP config
│   │   └── redisConfig/              # Redis client init
│   ├── handler/
│   │   ├── authHandler/
│   │   │   ├── authRequests.go       # Gin binding structs + validation tags
│   │   │   ├── loginHandler.go       # Login + Google-only account guard
│   │   │   ├── passwordHandler.go    # Change/forgot/verify/update password
│   │   │   ├── registerHandler.go    # Registration
│   │   │   └── sessionHandler.go     # Refresh · logout · email verification
│   │   └── oAuthHandler/
│   │       ├── googleOAuthHandler.go    # Google redirect + callback + cookie
│   │       └── googleRefreshHandler.go  # Refresh Google access token
│   ├── middleware/
│   │   ├── cors/cors.go              # CORS policy
│   │   ├── logger/access_log.go      # HTTP access logging
│   │   ├── rateLimiter/              # Per-IP Redis rate limiter
│   │   └── setup.go                  # Middleware registration
│   ├── model/AuthModel/authModel.go  # User · AuthMethod · Session · UserToken
│   ├── repository/authRepository.go  # All DB queries · inserts · transactions
│   ├── router/
│   │   ├── authRoutes.go             # /auth route declarations + rate limiter
│   │   └── router.go                 # Mounts groups under /api/v1
│   └── service/
│       ├── authService/
│       │   ├── authTypes.go          # Service input/output structs
│       │   ├── forgotPassword.go     # OTP generation · storage · cooldown
│       │   ├── helpers.go            # Token gen + hashing helpers
│       │   ├── login.go              # Login · refresh rotation · logout · verification
│       │   ├── password.go           # Change + forgot-password update
│       │   └── register.go           # User creation · token gen · email trigger
│       ├── oAuthService/
│       │   ├── oAuthlogin.go         # Find/create/link user · create session
│       │   └── oAuthTypes.go         # GoogleUser struct
│       └── rateLimiterService/
│           └── rateLimiterService.go # Redis counter with window config
├── migrations/
│   ├── authModelMigration.go         # GORM AutoMigrate definitions
│   └── migrate.go                    # Migration entry point
├── pkg/
│   ├── database/connection.go        # PostgreSQL connection (sync.Once guarded)
│   ├── logger/logger.go              # Logger setup
│   ├── mailer/mailer.go              # Verification + OTP emails
│   ├── redis/redis.go                # OTP · attempts · rate limit primitives
│   ├── tokenJWT/jwt.go               # JWT generation + validation
│   └── utils/
│       ├── crypto.go                 # Secure token gen + SHA-256 hash
│       ├── fomatValidationErrors.go  # Gin error formatter
│       └── passwordHashing.go        # Bcrypt hash + compare
└── vercel.json                       # Routes all traffic to api/index.go
```

---

## Runtime Modes

Both runtimes share the same routes, handlers, services, and repositories. `app/app.go` wires everything up once via `sync.Once` — env load, Redis, PostgreSQL, Gin engine, middleware, and routes. Only the entry point differs.

| | Server | Serverless |
|---|---|---|
| Entry point | `cmd/api/main.go` | `api/index.go` |
| HTTP lifecycle | `r.Run(":8000")` | Managed by Vercel |
| Connections | Persistent | Reused across warm instances |
| Email sending | Goroutine (async) | Synchronous (no goroutines) |
| Cold starts | None | ~300–500ms after inactivity |

---

## Architecture

```
HTTP Request
      │
      ▼
[ Entry Point ]     cmd/api/main.go (server) · api/index.go (serverless)
      │
      ▼
[ app.Handler() ]   Initialized once via sync.Once
      │
      ▼
[ Middleware ]       Rate Limiter · CORS · Logger · Recovery
      │
      ▼
[ Handler ]          Parse · Validate · Map errors to HTTP status
      │
      ▼
[ Service ]          Business logic · Token gen · Redis ops
      │
      ▼
[ Repository ]       GORM queries · Transactions · No logic
      │
      ▼
[ PostgreSQL ]    [ Redis ]
```

- **Handler** — parses request, validates input, maps errors to HTTP status; never touches the DB directly
- **Service** — all business logic, token generation, Redis operations; no HTTP knowledge
- **Repository** — GORM queries and transactions only; returns `nil, nil` for not-found
- Multi-step DB writes are always wrapped in transactions
- All init is `sync.Once` guarded — safe across server restarts and serverless warm reuse

---

## API Reference

**Base URL:** `/api/v1/auth` · **Rate limit:** 6 req/min per IP

### Registration & Verification

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/register` | Create account · sends verification email |
| `GET` | `/verification-email?token=` | Verify email · expires in 10 min · single-use |

### Login & Session

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/login` | Auth with email + password · sets HttpOnly refresh cookie |
| `POST` | `/refresh` | Rotate tokens · revokes old session · issues new pair |
| `POST` | `/logout` | Revokes session · clears cookie |

**Login errors:** `403` email not verified · `403` account `pending_otp` · `403` not active · `401` everything else

**Refresh cookie:** `HttpOnly · Path=/api/v1/auth/refresh · MaxAge=7d`

### Password

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/change-password` | Requires current password · enforces old ≠ new |
| `POST` | `/forgot-password` | Sends OTP · sets status `pending_otp` · 30s cooldown |
| `POST` | `/forgot-password-verify` | Verifies OTP · max 3 attempts |
| `POST` | `/forgot-password-update` | Sets new password · resets status to `active` |

### Google OAuth

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/google` | Redirect to Google consent screen |
| `GET` | `/google/callback` | Exchange code · create/link account · set cookie |
| `POST` | `/google/refresh` | Get new Google access token via stored refresh token |

---

## Auth Flows

### Login
```
POST /login
  → Block if Google-only account (no password method)
  → Bcrypt compare password
  → Assert email_verified · account_status = active
  → Generate access + refresh JWT
  → Hash refresh token → insert session row
  ← 200 access_token + HttpOnly cookie
```

### Token Refresh
```
POST /refresh
  → Validate JWT · hash token → lookup session
  → Revoke old session
  → Issue new access + refresh token pair → insert new session
  ← 200 new access_token + new cookie
```

### Logout
```
POST /logout
  → Hash token → revoke session
  → Clear cookie (MaxAge = -1)
  ← 200
```

**Session rules:**
- Raw token never stored — only SHA-256 hash in `sessions` table
- Each refresh rotates the token; a stolen token is invalidated on next legitimate use
- Sessions expire after 7 days (`expires_at > NOW()` in all queries)
- Cookie is scoped to `/api/v1/auth/refresh` — not sent with every request

---

## Google OAuth Flow

```
Callback received
  → Provider match (google + provider_user_id)?
      YES → update stored tokens + profile → login
      NO  →
        Email exists?
          NO  → create new user + google auth method → login
          YES →
            Has password method?
              YES → return ErrLinkRequired (frontend must confirm)
              NO  → link google method → login
  → Issue app JWT pair + session (same as password login)
```

- Google access token is used immediately to fetch user info then discarded
- Google refresh token is stored in `authentication_methods.oauth_refresh_token`

---

## Password Management

### Change Password
```
→ Verify old_password against stored hash
→ Reject if new_password = current (enforced via bcrypt compare)
→ Hash new_password → update authentication_methods
```

### Forgot Password
```
Step 1 — POST /forgot-password
  → Check 30s Redis cooldown
  → Generate 6-digit OTP · bcrypt-hash → store in Redis (TTL 5 min)
  → Set account_status = pending_otp
  → Send OTP email (synchronous)

Step 2 — POST /forgot-password-verify
  → Increment attempt counter · reject if > 3
  → Compare OTP against stored hash
  → On match → clear OTP + counter from Redis

Step 3 — POST /forgot-password-update
  → Reject if new = current password
  → Hash + update password · set account_status = active
```

---

## Configuration

Env vars are read via `os.Getenv`. Locally, a `.env` file is loaded automatically as a fallback.

```env
APP_ENV=development
GIN_MODE=debug

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=ats_db

REDIS_ADDR=localhost:6379
REDIS_USERNAME=
REDIS_PASSWORD=
REDIS_DB=0

JWT_SECRET=your_access_secret
JWT_REFRESH_SECRET=your_refresh_secret

CLIENT_ID=your_google_client_id
CLIENT_SECRET=your_google_client_secret
REDIRECT_URL=http://localhost:8000/api/v1/auth/google/callback

SMTP_FROM=noreply@example.com
SMTP_PASSWORD=smtp_password
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_URL=http://localhost:8000

PORT=8000
```

---

## Getting Started

**Prerequisites:** Go 1.21+ · PostgreSQL · Redis

```bash
git clone https://github.com/your-org/go-auth-backend-api.git
cd go-auth-backend-api

cp .env.example .env
go mod download
go run ./cmd/migrate/main.go

air                          # hot reload
# or
go run ./cmd/api/main.go     # direct run
```

---

## Deployment

### Docker
```bash
docker build -t go-auth-api .
docker run --env-file .env -p 8000:8000 go-auth-api
```

### Railway / Render / Fly.io
1. Connect GitHub repo — auto-detects Go + Dockerfile
2. Add PostgreSQL and Redis as managed add-ons
3. Set env vars in platform dashboard
4. Update `REDIRECT_URL` and `SMTP_URL` to your deployed domain

### Vercel
`api/index.go` and `vercel.json` are already in the repo — no extra setup needed.

1. Set all env vars in **Vercel → Project Settings → Environment Variables**
   - `REDIRECT_URL` → `https://your-app.vercel.app/api/v1/auth/google/callback`
   - `SMTP_URL` → `https://your-app.vercel.app`
   - Use [Neon](https://neon.tech) for PostgreSQL and [Upstash](https://upstash.com) for Redis
2. Add the callback URL to **Google Cloud Console → Authorized redirect URIs**
3. Deploy
```bash
vercel --prod
```

---

## Security Design

| Concern | Approach |
|---|---|
| **Token storage** | Raw refresh token never stored — SHA-256 hash only |
| **XSS** | Refresh token in `HttpOnly` cookie scoped to `/refresh` path |
| **Token theft** | Rotation on every refresh — stolen token invalidated on next use |
| **OTP abuse** | 3-attempt lockout + 30s cooldown tracked in Redis |
| **Brute force** | Per-IP rate limiter (6 req/min) on all endpoints |
| **OAuth takeover** | Google linking requires password confirmation — no silent merge |
| **Enumeration** | All login failures return `"invalid credentials"` |
| **Partial writes** | Multi-step DB operations wrapped in GORM transactions |
| **Password reuse** | Old = new enforced via bcrypt compare before any update |

---

## Production Checklist

- Strong unique values for `JWT_SECRET` and `JWT_REFRESH_SECRET`
- `APP_ENV=production` · `GIN_MODE=release`
- CORS restricted to your frontend domain
- `Secure=true` on cookies (HTTPS required)
- `REDIRECT_URL` updated in Google Cloud Console
- PostgreSQL with SSL (`sslmode=require`)
- Managed Redis reachable over the internet (Upstash for Vercel)
- Migrations run against production DB before first deploy
