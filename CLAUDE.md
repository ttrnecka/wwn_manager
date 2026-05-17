# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Toolchain

Versions are pinned in `mise.toml` (run `mise install` once). The repo expects `mise`, Go 1.26.x, golangci-lint 2.x, Node 26, `just`, and `tmux`. Most tasks are wrapped by `just` — read `Justfile` before inventing a new command.

## Common commands

- `just dev` — starts MongoDB in docker, runs `npm install`, then opens a tmux split with `npm run dev` (Vite) on one pane and `go run .` in `webapi/` on the other. The Go server listens on `:8888` by default (`ADDRESS` env), and Vite proxies/serves the frontend.
- `just docker` — full stack via docker-compose (nginx + webapi + mongodb + a one-shot `frontend-artifacts` container that copies the built SPA into `./dist/`). Open http://localhost and log in as `admin` / `password`.
- `just dev-mongo` / `just docker-stop` — start/stop just the Mongo container.
- `just install` — `npm install` inside `frontend/`.
- `just golint` — runs `golangci-lint run webapi/`.
- `just test` — autodetects podman vs docker, starts the user podman socket if needed, exports `DOCKER_HOST` to the socket path, and runs `go test ./webapi/...`. Tests use **testcontainers-go** to spin up real MongoDB, so a working container runtime socket is required — there is no mock mode. To run a single test: `DOCKER_HOST=… go test ./webapi/internal/service/ -run TestName -v`.
- `just build` — produces a deployable `wwn_manager_<date>.zip` in the repo root via `scripts/standalone_build/build.sh` (uses `docker-compose-build.yaml` to build both backend and frontend, then zips `build/`).
- `just scan` — vulnerability scan via `scripts/scan/scan.sh` (Trivy-based, see `Dockerfile.trivy`).
- `just update-backend` — `cd webapi && go get -u`. After upgrading Go, also update `mise.toml` and the `golang:` base image in `webapi/Dockerfile*`.

## Environment

`.env` is loaded by the webapi at startup (see `utils.LoadEnv`). Template in `.env.template`:

- `ADDRESS` — bind address for the Go HTTP server (default `:8888`).
- `MONGO_URI` — Mongo connection string. Compose overrides this to `mongodb://mongodb_wwnmgr:27017/`.
- `ITA_API_URI` + `ITA_TOKEN` + `ITA_FEED_ID` — credentials for the upstream ITA report API. If unset, `ita.NewITAClient` returns an error and the `/api/v1/import_api` endpoint will fail.

The Mongo database name is hardcoded to `wwn_identity` (`webapi/db/db.go`).

## Architecture

### Big picture

WWN Manager is a SPA + Go API for managing Fibre Channel **WWN** (World-Wide Name) entries across multiple customers. It ingests WWN data (CSV upload or pulled from the external **ITA** reporting API), applies user-defined regex rules to classify each WWN (host vs. array vs. backup, hostname extraction, customer reassignment, reconciliation), and supports taking immutable **snapshots** of the entry collection for point-in-time export.

### Backend layout (`webapi/`)

A single Go module (`github.com/ttrnecka/wwn_identity/webapi`) registered into the root `go.work`. It builds as a service binary using `kardianos/service` — the same binary runs as a Windows service, a Linux service, or in console mode (`service.Interactive()` branch in `main.go`).

Layered architecture, wired top-down in `server/router.go`:

```
entity   → BSON models + Mongo collection factories (FCWWNEntry, Rule, Snapshot, User)
repository → typed CRUD wrappers built on github.com/ttrnecka/agent_poc/common/db (CRUD[T])
service  → business logic; GenericService[T] in service/generic_service.go provides
           All/Find/Get/Update/Delete/SoftDelete/Restore/InsertAll on top of the repo,
           plus a DependencyDeleteFunc registration mechanism for cascade deletes
handler  → Echo HTTP handlers; validation via go-playground/validator
mapper   → entity ↔ shared/dto translation
server   → Echo router setup + middleware (auth, session, request/response capture, logging)
shared/dto → DTOs exchanged with the frontend (gob-registered in main.go for sessions)
ita      → client for the upstream ITA report API and its JSON schema
```

Routes are all under `/api/v1` (see `server/router.go`). Auth uses gorilla cookie sessions with a static secret — the `AuthMiddleware` protects everything except `/login`. Static SPA assets are served from `<binary_dir>/static/`, and any non-API route falls through to `index.html` (SPA history mode).

### Domain concepts to know before editing

- **Customer scoping**: every `FCWWNEntry` belongs to a `Customer`. Two sentinel customer names are reserved: `__GLOBAL__` (cross-customer view) and `<NO CUSTOMER>` (entries without an assigned customer) — both defined in `entity/common.go`. Treat them as constants, never hardcode.
- **Rules** (`entity/rules.go`) are typed regex rules. Categories: **host rules** (zone/alias/wwn_host_map) infer hostnames; **range rules** (wwn_range_array/backup/host/other) classify WWN ranges into device types; **reconcile rules** (wwn_customer_map, ignore_loaded, default_reconcile_rule_*) drive cross-customer reconciliation. Rules have a `Customer` field — a rule belonging to `__GLOBAL__` applies everywhere.
- **WWNSet** on an entry encodes provenance: `1 = SAN-discovered`, `2 = manual`, `3 = auto/derived`. Reconciliation logic checks this.
- **Snapshots** (`entity/snapshot.go`): `MakeSnapshot` uses a Mongo `$out` aggregation to copy `fc_wwn_entries` into a sibling collection `fc_wwn_entries_<unix_ts>` and also copies the indexes. Snapshot entry collections are read-only by convention; use `Snapshot.GetEntries(db)` to access them through the same `CRUD[FCWWNEntry]` wrapper.
- **Indexes** (`db/schema.go`): `fc_wwn_entries` has a unique compound `(customer, wwn)` index. Bulk inserts that may collide must use upsert / unordered writes — duplicate-key errors are not handled silently.
- **Default admin** is upserted on every startup (`db/schema.go::EnsureUserCollection`) with username `admin` / password `password`.

### Frontend (`frontend/`)

Vue 3 + Vite + Pinia + vue-router + Bootstrap 5. Source under `src/`:

- `views/` — top-level pages: `GlobalFCManager` (cross-customer), `FCManager` (single customer), `Summary`, `LoginView`, `AboutView`.
- `stores/` — Pinia stores; `apiStore.js` is the central data fetcher (initialized by the router's `beforeEach` when a route has `meta.requiresData`), `userStore.js` tracks auth, `flash.js` for notifications.
- `services/fcService.js` — axios calls to `/api/v1/...`; auth and alert helpers live in `services/auth.js` / `services/alert.js`.
- `components/` — table and modal widgets used by the views.

Dev server runs on the default Vite port and is expected to reach the Go API at `:8888` (proxy / `VITE_*` config in `vite.config.js`). In production the SPA is built into `dist/` and served by nginx (`nginx/conf/`); the standalone build instead serves it via the Go binary from `static/`.

## Linting and CI conventions

`.golangci.yml` enables `wrapcheck` — wrap errors returned across package boundaries with `fmt.Errorf("…: %w", err)`. Repository/service layers already follow this; the `// nolint:wrapcheck // delegating to repository, which already wraps errors` comments in `generic_service.go` are the canonical exception pattern (don't add new `nolint` without an inline reason). `gocyclo` is set to min-complexity 15. Echo `JSON`/`NoContent`/`Attachment` return values are exempt from `wrapcheck` via `extra-ignore-sigs`. Test files are exempted from `errcheck`, `gosec`, `wrapcheck`, and `revive`.

The only GitHub Actions workflow is `security.yml` — there is no CI-driven test/lint job, so run `just golint` and `just test` locally before pushing.
