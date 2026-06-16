# AGENTS.md — 喵玉殿论坛 (neonya)

Mono-repo with pnpm workspace. 喵玉殿新版论坛。

## Directory ownership

| Path           | Package                                        | Stack                                    |
| -------------- | ---------------------------------------------- | ---------------------------------------- |
| `apps/server/` | Go module `github.com/NSDN/neonya/apps/server` | Gin + GORM + PostgreSQL                  |
| `apps/web/`    | `@neonya/web`                                  | Vue 3 + Vite + Naive UI + Pinia + ofetch |
| `docs/`        | `@neonya/docs`                                 | VitePress                                |

## Developer commands

All pnpm commands run from repo root.

```sh
# Frontend
pnpm dev                     # start vite dev server (port 10123)
pnpm --filter @neonya/web dev
pnpm build                   # typecheck + build

# Backend (from apps/server/)
go run ./cmd/server          # start server (port 10127)
go build ./...               # compile check

# Database
./apps/server/podman-db.ps1 start     # start postgres (pod + container)
./apps/server/podman-db.ps1 stop      # stop and remove pod
./apps/server/podman-db.ps1 status    # check pod/container status
./apps/server/podman-db.ps1 remove    # remove pod and all resources
migrate -source file://apps/server/database/migrations \
  -database "postgres://forum_user:password@localhost:5432/forum?sslmode=disable" up

# Docs
pnpm docs:dev                # vitepress (port 10126)
```

## Setup checklist

1. Copy `apps/server/.env.example` → `apps/server/.env` (set APPLICATION_PORT, TOKEN_KEY)
2. Copy `apps/server/.env.postgres.example` → `apps/server/.env.postgres`
3. Start postgres with podman, then run migrations
4. `pnpm install` from root (workspace installs all 3 packages)
5. `pnpm --filter @neonya/web dev` to start frontend in a separate terminal

## Architecture conventions

### Go backend — functional style, no OOP layers

There is **no** `controllers → services → repositories` layering. Each domain package (`internal/auth/`, `internal/board/`, `internal/thread/`) contains handler + service + model in one place, using **pure functions with explicit dependency injection**:

```go
// Handlers return gin.HandlerFunc closures capturing deps
func HandleRegister(db *gorm.DB) gin.HandlerFunc { ... }

// Services are pure functions — deps passed as args, no struct receivers
func Register(db *gorm.DB, info *RegisterInfo) (*UserPublicInfo, error) { ... }
```

- `internal/` is Go-compiler-enforced isolation (cannot be imported by `apps/web`)
- `main.go` wires everything: loads env → opens DB → creates router → calls each domain's `RegisterRoutes()`
- `internal/shared/` holds cross-cutting code: database, middleware (CORS + JWT auth), response helpers, message constants
- `internal/config/` holds all constants (env var names, HTTP headers, page types, API paths)
- Package `appcontext` (not `context`) to avoid collision with Go stdlib

### Frontend — feature-based modules

- `src/features/authorization/`, `board/`, `thread/` — each contains `apis/`, `composables/`, `stores/`, `types.ts`
- `src/shared/` — cross-cutting: `components/`, `composables/`, `services/`, `utils/`, `constants/`
- `src/stores/` — global Pinia stores (not `store/`)
- `@` alias resolves to `/src`
- HTTP layer: `ofetch` (not axios), entry point is `shared/services/request.ts` returning `Result<T, ApiError>` (Rust-style)

### Both sides

- **No abbreviations** in names — use full words (`authorization` not `auth`, `appcontext` not `appctx`)

## Style

| Aspect                 | Rule                                                                                               |
| ---------------------- | -------------------------------------------------------------------------------------------------- |
| Go indent              | tabs, width 4                                                                                      |
| TS/Vue/JSON indent     | spaces, width 2                                                                                    |
| TS semicolons          | no                                                                                                 |
| TS quotes              | single                                                                                             |
| TS trailing commas     | none                                                                                               |
| Comments               | only add when logic is non-obvious; no AI-generated comments                                       |
| Go formatting          | `gofmt` standard                                                                                   |
| Prettier               | `.prettierrc` at root, `arrowParens: "avoid"`                                                      |
| Naming                 | use full words, no abbreviations (`address` not `addr`, `sourceName` not `dsn`)                    |
| One-liner calls        | prefer single call expressing intent (`strings.CutPrefix` over `HasPrefix` + `TrimPrefix`)         |
| Intermediate variables | omit when the inline expression is self-explanatory (`f(mustGetenv("K"))` not `v := ...; f(v)`)    |
| Markdown prose         | Chinese sentences end with `。`; break lines at sentence boundaries when a line would exceed width; use trailing `  ` (double space) for hard line breaks when a sentence continues across multiple lines in rendered output |

## Theme

Dark mode is default. `useNaiveUIGlobalConfig` returns dark theme first; CSS `:root` sets dark variables as base, `@media (prefers-color-scheme: light)` overrides for light.

## Container

**Podman only** (not Docker). DB is managed with `podman play kube pod.yaml` via `podman-db.ps1`. Build images with `podman build -f Containerfile`.

## Testing

No test runner configured yet. Backend: `go build ./...` is the only automated check. Frontend: `vue-tsc` typecheck.

## Hard constraints

- **Never install packages or tools without explicit user approval.** Always ask for consent before running any install command (npm install, go get, pip install, apt, winget, etc.).

## Decision records

Architecture and design decisions made during development are documented in
`docs/source/decisions/`. When context is ambiguous, check there first for
established decisions before proposing new approaches.
