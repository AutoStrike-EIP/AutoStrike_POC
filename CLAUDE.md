# AutoStrike - Developer Context

## Project Overview

AutoStrike is a **Breach and Attack Simulation (BAS)** platform for security testing based on the **MITRE ATT&CK** framework. It allows SOC teams and security professionals to validate their detection capabilities through automated attack simulations.

## Architecture

### Components

```
┌─────────────────────────────────────────────────────────────────┐
│                    Server (Go) - Port 8443                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │  Dashboard  │  │  REST API   │  │  WebSocket  │              │
│  │  (Static)   │  │  /api/v1/*  │  │  /ws/*      │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
└─────────────────────────────────────────────────────────────────┘
                              ▲
                              │ WebSocket
                    ┌─────────┴─────────┐
                    │      Agent        │
                    │     (Rust)        │
                    └───────────────────┘
```

**Single server on port 8443** serves:
- Dashboard (static files from `dashboard/dist`)
- REST API (`/api/v1/*`)
- WebSocket (`/ws/agent`, `/ws/dashboard`)
- Health check (`/health`)

### Directory Structure

```
autostrike/
├── server/          # Go backend (Hexagonal Architecture)
│   ├── cmd/         # Entry point
│   ├── configs/     # YAML technique definitions
│   │   └── techniques/  # discovery.yaml, execution.yaml, etc.
│   ├── internal/
│   │   ├── domain/      # Business logic (entities, services, interfaces)
│   │   ├── application/ # Use cases (service orchestration)
│   │   └── infrastructure/ # Adapters (HTTP, SQLite, WebSocket)
│   └── pkg/         # Shared utilities
├── agent/           # Rust agent
│   └── src/         # Client, executor, config, system info
├── dashboard/       # React frontend
│   └── src/
│       ├── components/  # MitreMatrix, RunExecutionModal, Layout, ProtectedRoute
│       ├── contexts/    # AuthContext (authentication state)
│       ├── hooks/       # useWebSocket
│       ├── pages/       # Dashboard, Agents, Techniques, Matrix, Scenarios, Executions, ExecutionDetails, Settings, Login, Analytics, Scheduler, Admin/*
│       └── lib/         # API client
├── docs/            # MkDocs documentation
└── scripts/         # Build and deployment scripts
    └── mitre-import/    # MITRE ATT&CK import tool (Go CLI)
```

## Tech Stack

| Component | Technology |
|-----------|------------|
| **Server** | Go 1.21+, Gin, SQLite, gorilla/websocket |
| **Agent** | Rust 1.75+, tokio, tokio-tungstenite |
| **Dashboard** | React 18, TypeScript, TailwindCSS, Vite, TanStack Query, Chart.js |
| **Documentation** | MkDocs Material |

## Key Patterns

### Server (Hexagonal Architecture)
- **Domain Layer**: Pure business logic, no external dependencies
- **Application Layer**: Use case orchestration
- **Infrastructure Layer**: External adapters (HTTP, persistence, WebSocket)
- Dependencies flow INWARD toward domain

### Agent (Async Rust)
- Exponential backoff for reconnection (1s → 60s max)
- Platform-specific command execution (PowerShell/cmd/bash/sh)
- JSON-based WebSocket protocol
- Heartbeat every 30 seconds

### Dashboard (React)
- TanStack Query for server state
- WebSocket for real-time updates
- Conditional polling as fallback

## Development Commands

```bash
# Quick Start (recommended)
make run        # Build and start everything on https://localhost:8443
make agent      # Connect an agent
make stop       # Stop all services
make logs       # View server logs

# Server
cd server && go build ./...
cd server && go test ./... -cover

# Agent
cd agent && cargo build --release
cd agent && cargo test

# Dashboard
cd dashboard && npm install
cd dashboard && npm run build
cd dashboard && npm test -- --run
cd dashboard && npm run lint
cd dashboard && npm run type-check

# Full stack (Docker)
docker compose up --build

# Generate TLS certificates
make certs

# MITRE ATT&CK Import
make import-mitre       # Import techniques from STIX + Atomic Red Team
make import-mitre-safe  # Import only safe techniques
make import-mitre-dry   # Dry run (show stats without writing files)
```

## API Endpoints

Base URL: `https://localhost:8443/api/v1`

### Authentication (public routes)
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/auth/login` | POST | Login with username/password |
| `/auth/refresh` | POST | Refresh access token |
| `/auth/logout` | POST | Invalidate tokens |
| `/auth/me` | GET | Get current user info (requires token) |

### Core API (protected when auth enabled)
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check (returns `{"status": "ok", "auth_enabled": bool}`) |
| `/agents` | GET | List agents (`?all=true` for offline too) |
| `/agents/:paw` | GET | Get agent details |
| `/agents` | POST | Register agent |
| `/agents/:paw` | DELETE | Delete agent |
| `/agents/:paw/heartbeat` | POST | Update last_seen |
| `/techniques` | GET | List all techniques |
| `/techniques/:id` | GET | Get technique by ID |
| `/techniques/tactic/:tactic` | GET | Techniques by tactic |
| `/techniques/platform/:platform` | GET | Techniques by platform |
| `/techniques/coverage` | GET | MITRE coverage stats |
| `/techniques/:id/executors` | GET | List executors for a technique (`?platform=linux`) |
| `/techniques/import` | POST | Import from YAML |
| `/scenarios` | GET | List scenarios |
| `/scenarios/:id` | GET | Get scenario |
| `/scenarios/tag/:tag` | GET | Scenarios by tag |
| `/scenarios` | POST | Create scenario |
| `/scenarios/:id` | PUT | Update scenario |
| `/scenarios/:id` | DELETE | Delete scenario |
| `/executions` | GET | List executions (limit 50) |
| `/executions/:id` | GET | Get execution |
| `/executions` | POST | Start execution |
| `/executions/:id/results` | GET | Get results |
| `/executions/:id/stop` | POST | Stop execution |
| `/executions/:id/complete` | POST | Complete execution |

### Admin API (requires admin role)
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/admin/users` | GET | List all users |
| `/admin/users/:id` | GET | Get user by ID |
| `/admin/users` | POST | Create user |
| `/admin/users/:id` | PUT | Update user |
| `/admin/users/:id/role` | PUT | Update user role |
| `/admin/users/:id` | DELETE | Deactivate user |
| `/admin/users/:id/reactivate` | POST | Reactivate user |
| `/admin/users/:id/reset-password` | POST | Reset password |

### Schedules API
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/schedules` | GET | List all schedules |
| `/schedules/:id` | GET | Get schedule by ID |
| `/schedules` | POST | Create schedule |
| `/schedules/:id` | PUT | Update schedule |
| `/schedules/:id` | DELETE | Delete schedule |
| `/schedules/:id/pause` | POST | Pause schedule |
| `/schedules/:id/resume` | POST | Resume schedule |
| `/schedules/:id/run` | POST | Run schedule now |
| `/schedules/:id/runs` | GET | Get schedule run history |

### Notifications API
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/notifications` | GET | Get notifications |
| `/notifications/unread/count` | GET | Get unread count |
| `/notifications/:id/read` | POST | Mark as read |
| `/notifications/read-all` | POST | Mark all as read |
| `/notifications/settings` | GET/POST/PUT/DELETE | Manage notification settings |
| `/notifications/smtp` | GET | Get SMTP config (admin) |
| `/notifications/smtp/test` | POST | Test SMTP (admin) |

### Analytics API
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/analytics/period` | GET | Get period stats |
| `/analytics/comparison` | GET | Compare periods |
| `/analytics/trend` | GET | Get score trend |
| `/analytics/summary` | GET | Get execution summary |

### Permissions API
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/permissions/matrix` | GET | Get permission matrix |
| `/permissions/me` | GET | Get my permissions |
| `/permissions/check` | POST | Check permission |

## WebSocket Protocol

### Agent ↔ Server
Connection: `wss://server:8443/ws/agent`

```json
// Register (Agent → Server)
{"type": "register", "payload": {"paw": "...", "hostname": "...", "platform": "...", "executors": [...]}}

// Registered (Server → Agent)
{"type": "registered", "payload": {"status": "ok", "paw": "..."}}

// Heartbeat (Agent → Server, every 30s)
{"type": "heartbeat", "payload": {"paw": "..."}}

// Task (Server → Agent)
{"type": "task", "payload": {"id": "...", "technique_id": "...", "command": "...", "executor": "...", "timeout": 300}}

// Task Result (Agent → Server)
{"type": "task_result", "payload": {"task_id": "...", "technique_id": "...", "success": true, "output": "...", "exit_code": 0}}

// Task Ack (Server → Agent)
{"type": "task_ack", "payload": {"task_id": "...", "status": "received"}}
```

### Dashboard ↔ Server
Connection: `wss://server:8443/ws/dashboard`

```json
// Server broadcasts to all dashboards
{"type": "execution_started", "payload": {"execution_id": "...", "data": {...}}}
{"type": "execution_completed", "payload": {"execution_id": "...", "data": {...}}}
{"type": "execution_cancelled", "payload": {"execution_id": "...", "data": {...}}}

// Dashboard can send ping
{"type": "ping", "payload": {}}
// Server responds with pong
{"type": "pong", "payload": {}}
```

## Environment Variables

### Server (.env)
- `DATABASE_PATH` - SQLite database path (default: `./data/autostrike.db`)
- `DASHBOARD_PATH` - Path to dashboard dist folder (default: `../dashboard/dist`)
- `JWT_SECRET` - JWT signing key (optional - auth disabled if not set)
- `AGENT_SECRET` - Agent authentication secret
- `ENABLE_AUTH` - Explicit auth override (`true`/`false`)
- `ALLOWED_ORIGINS` - CORS origins (default: `localhost:3000,localhost:8443`)
- `LOG_LEVEL` - Logging level (debug, info, warn, error)

**Authentication behavior:**
- `JWT_SECRET` not set → Auth **disabled** (development mode)
- `JWT_SECRET` set → Auth **enabled** (production mode)

### Dashboard (.env) - Only for Vite dev server
- `VITE_SERVER_URL` - Backend server URL
- `VITE_API_BASE_URL` - API base path
- `VITE_WS_HOST` - WebSocket host override

## Available Techniques

48 built-in techniques ship with the project. Run `make import-mitre` to import additional techniques from MITRE ATT&CK STIX + Atomic Red Team (typically 200+ techniques).

| Tactic | Built-in | IDs (built-in) |
|--------|----------|-----------------|
| Reconnaissance | 2 | T1592.004, T1595.002 |
| Initial Access | 3 | T1078, T1133, T1190 |
| Execution | 5 | T1059.001, T1059.003, T1059.004, T1047, T1059.006 |
| Persistence | 4 | T1053.005, T1547.001, T1053.003, T1543.002 |
| Privilege Escalation | 4 | T1548.001, T1548.002, T1078.003, T1134.001 |
| Defense Evasion | 6 | T1070.004, T1562.001, T1027, T1070.001, T1036.005, T1218.011 |
| Credential Access | 4 | T1552.001, T1555.003, T1003.008, T1552.004 |
| Discovery | 9 | T1082, T1083, T1057, T1016, T1049, T1087, T1069, T1018, T1007 |
| Lateral Movement | 3 | T1021.001, T1021.002, T1021.004 |
| Collection | 4 | T1005, T1039, T1074.001, T1119 |
| Command and Control | 3 | T1071.001, T1105, T1572 |
| Exfiltration | 3 | T1048.003, T1041, T1567.002 |
| Impact | 3 | T1490, T1489, T1486 |

All built-in techniques are **safe mode compatible** (non-destructive). Imported techniques from Atomic Red Team may include unsafe techniques; use `make import-mitre-safe` to import only safe ones.

### Key Data Model Additions

**Executor** (extended fields):

| Field | Type | Description |
|-------|------|-------------|
| `name` | string (optional) | Executor name (distinguishes multiple executors per technique) |
| `platform` | string (optional) | Target platform (`windows`, `linux`, `macos`) |
| `elevation_required` | bool (optional) | Whether root/admin privileges are needed |

**Technique** (extended fields):

| Field | Type | Description |
|-------|------|-------------|
| `tactics` | []string (optional) | All MITRE tactics (multi-tactic support). Fallback: `[tactic]` |
| `references` | []string (optional) | MITRE ATT&CK reference URLs |

**TechniqueSelection** (scenario phases):

| Field | Type | Description |
|-------|------|-------------|
| `technique_id` | string | Technique ID |
| `executor_name` | string (optional) | Preferred executor name (empty = auto-select) |

Phase `techniques` field accepts both `[]string` (legacy) and `[]TechniqueSelection` (new format) via custom JSON unmarshaling.

## Important Notes

1. **Security**: This is a security testing tool. Use only on authorized systems.
2. **Safe Mode**: Always test with `safe_mode: true` first.
3. **mTLS**: Production deployments should use mutual TLS.
4. **Logging**: Server uses structured logging (zap), agent uses tracing.

## Testing

**Test coverage (Phase 3):**
- **Server**: 200+ tests
  - application: 83.0%
  - entity: 95.0%
  - service: 99.2%
  - handlers: 87.5%
  - websocket: 91.6%
  - middleware: 94.3%
  - rest/server: 87.9%
  - sqlite: 85.0%
- **Agent**: 67 unit tests (`cargo test`)
- **Dashboard**: 513 tests across 25 files (`npm test`)

```bash
# Run all tests
cd server && go test ./... -cover
cd agent && cargo test
cd dashboard && npm test -- --run
```

## Security Score Formula

```
score = (blocked * 100 + detected * 50) / (total * 100) * 100
```

| Result | Points | Meaning |
|--------|--------|---------|
| Blocked | 100 | Attack prevented |
| Detected | 50 | Attack seen but not stopped |
| Success | 0 | Attack succeeded undetected |

## Git Rules

### Branches
- NEVER push to main directly
- Naming: `feat/<scope>-<description>`, `fix/<scope>-<description>`, `test/<scope>-<description>`
- Scope: server, dashboard, agent, techniques, scenarios, docs, ci
- Always branch from latest main: `git checkout main && git pull && git checkout -b feat/...`

### Commits
- Format: `type(scope): description`
- Types: feat, fix, test, refactor, docs, ci, chore
- Examples: `feat(server): add PDF report endpoint`, `test(dashboard): add Reports page tests`
- NEVER mention Claude, AI, LLM, assistant, copilot, or Anthropic in commits
- NEVER add Co-Authored-By with Claude or Anthropic
- Messages in English, concise, descriptive
- Stage specific files (`git add <files>`) — never use `git add -A` or `git add .`

### Pull Requests
- Use `gh pr create` with title (<70 chars) and body
- NEVER mention Claude/AI/LLM in PRs
- PR body format:
  ```
  ## Summary
  - <change 1>
  - <change 2>

  ## Test plan
  - [ ] Unit tests added
  - [ ] Coverage >= 95%
  - [ ] Lint passes
  - [ ] No regressions
  ```

## Quality Standards

### Coverage Targets
- Server Go: 95%+ per package
- Dashboard React: 95%+ lines
- Agent Rust: 90%+

### Lint (must pass with zero warnings)
```bash
cd server && go vet ./... && go test ./... -cover
cd dashboard && npm run lint && npm run type-check && npm test -- --run
cd agent && cargo fmt --check && cargo clippy -- -D warnings && cargo test
```

### Code Rules
- No `any` in TypeScript unless justified with a comment
- No `unsafe` in Rust unless justified and documented
- No `panic` in Go handlers — always return proper HTTP errors
- Input validation on all public endpoints
- Error handling explicit — no unhandled promises, no ignored errors
- Hexagonal architecture: domain/ has ZERO imports from infrastructure/

## Contributing

1. Follow hexagonal architecture for server changes
2. Use `cargo fmt` and `cargo clippy` for Rust
3. Use `npm run lint` and `npm run type-check` for TypeScript
4. Add tests for new functionality
5. Update documentation in `docs/` for user-facing changes
