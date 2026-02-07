# Backend (Go)

The AutoStrike control server is built with **Go 1.24+** using the **Gin** framework and **hexagonal architecture**.

---

## Qu'est-ce que l'Architecture Hexagonale ?

L'architecture hexagonale (aussi appelée **Ports & Adapters**) est un pattern qui **isole la logique métier** des détails techniques (base de données, HTTP, fichiers, etc.).

### Le Principe

```
┌─────────────────────────────────────────────────────────────┐
│                    INFRASTRUCTURE                            │
│    (HTTP handlers, SQLite, WebSocket, fichiers YAML)        │
│  ┌───────────────────────────────────────────────────────┐  │
│  │                   APPLICATION                         │  │
│  │           (Orchestration des use cases)               │  │
│  │  ┌─────────────────────────────────────────────────┐  │  │
│  │  │                   DOMAIN                        │  │  │
│  │  │                                                 │  │  │
│  │  │   Entités : Agent, Technique, Execution, User   │  │  │
│  │  │             Notification, Schedule, Permission   │  │  │
│  │  │   Services : Orchestrator, Validator, Score     │  │  │
│  │  │   Interfaces : Repository (ports)               │  │  │
│  │  │                                                 │  │  │
│  │  │        ⚠️ AUCUNE DÉPENDANCE EXTERNE             │  │  │
│  │  └─────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘

Les dépendances vont TOUJOURS vers l'intérieur (→ Domain)
```

### Pourquoi c'est utile ?

| Avantage | Exemple concret dans AutoStrike |
|----------|--------------------------------|
| **Changer de BDD facilement** | SQLite → PostgreSQL = modifier uniquement `infrastructure/persistence/` |
| **Changer de framework HTTP** | Gin → Echo = modifier uniquement `infrastructure/http/` |
| **Tests unitaires simples** | Tester `ExecutionService` sans base de données (mocks) |
| **Code métier stable** | Ajouter une API GraphQL sans toucher au domain |

### Les 3 couches

| Couche | Responsabilité | Dépend de |
|--------|----------------|-----------|
| **Domain** | Logique métier pure, entités, interfaces | Rien (indépendant) |
| **Application** | Orchestration des use cases | Domain uniquement |
| **Infrastructure** | Adapters externes (HTTP, DB, WS) | Application + Domain |

### Exemple : Ajouter un agent

```
1. HTTP Handler reçoit POST /api/v1/agents
   └─> infrastructure/http/handlers/agent_handler.go

2. Handler appelle AgentService.Register()
   └─> application/agent_service.go

3. Service valide et appelle Repository.Save()
   └─> domain/repository/agent_repository.go (interface)

4. SQLite Repository implémente l'interface
   └─> infrastructure/persistence/sqlite/agent_repository.go
```

**Le Domain ne sait pas** que c'est SQLite ou HTTP. Il définit juste des interfaces.

---

## Structure des dossiers

```
server/
├── cmd/autostrike/
│   └── main.go                    # Entry point, DI, startup
├── configs/
│   └── techniques/                # YAML technique definitions (auto-loaded via os.ReadDir)
│       ├── reconnaissance.yaml    # 13 built-in files + any imported via make import-mitre
│       ├── initial-access.yaml
│       ├── execution.yaml
│       ├── persistence.yaml
│       ├── privilege-escalation.yaml
│       ├── defense-evasion.yaml
│       ├── credential-access.yaml
│       ├── discovery.yaml
│       ├── lateral-movement.yaml
│       ├── collection.yaml
│       ├── command-and-control.yaml
│       ├── exfiltration.yaml
│       └── impact.yaml
├── internal/
│   ├── domain/                    # 🟢 Business Layer (independent)
│   │   ├── entity/                # Entities
│   │   │   ├── agent.go           # Agent, AgentStatus
│   │   │   ├── technique.go       # Technique, Executor, Detection
│   │   │   ├── scenario.go        # Scenario, Phase
│   │   │   ├── execution.go       # Execution, SecurityScore
│   │   │   ├── result.go          # ExecutionResult, ResultStatus
│   │   │   ├── user.go            # User, UserRole
│   │   │   ├── notification.go    # Notification, NotificationSettings, SMTPConfig
│   │   │   ├── schedule.go        # Schedule, ScheduleRun, ScheduleFrequency
│   │   │   └── permission.go      # Permission, PermissionMatrix
│   │   ├── repository/            # Interfaces (outbound ports)
│   │   └── service/               # Domain services
│   │       ├── orchestrator.go    # Attack orchestration
│   │       ├── validator.go       # Compatibility validation
│   │       └── score_calculator.go # Security score calculation
│   ├── application/               # 🟡 Use Cases
│   │   ├── agent_service.go       # Agent CRUD, heartbeat
│   │   ├── auth_service.go        # Authentication (login, tokens, JWT)
│   │   ├── execution_service.go   # Execution lifecycle
│   │   ├── scenario_service.go    # Scenario management
│   │   ├── technique_service.go   # Technique catalog
│   │   ├── notification_service.go # Notification management, SMTP
│   │   ├── schedule_service.go    # Schedule management, cron
│   │   ├── analytics_service.go   # Analytics, trends, comparisons
│   │   └── token_blacklist.go     # JWT token blacklist for logout
│   └── infrastructure/            # 🔵 External Adapters
│       ├── api/rest/
│       │   └── server.go          # Gin REST server, route registration
│       ├── http/
│       │   ├── handlers/          # HTTP handlers
│       │   │   ├── agent_handler.go
│       │   │   ├── auth_handler.go
│       │   │   ├── technique_handler.go
│       │   │   ├── scenario_handler.go
│       │   │   ├── execution_handler.go
│       │   │   ├── admin_handler.go        # User management (admin)
│       │   │   ├── analytics_handler.go    # Analytics endpoints
│       │   │   ├── notification_handler.go # Notification endpoints
│       │   │   ├── schedule_handler.go     # Schedule endpoints
│       │   │   ├── permission_handler.go   # Permission endpoints
│       │   │   └── websocket_handler.go
│       │   └── middleware/
│       │       ├── auth.go        # JWT auth, agent auth, roles, permissions
│       │       ├── security.go    # Security headers (HSTS, CSP, etc.)
│       │       ├── ratelimit.go   # Per-IP rate limiting
│       │       └── logging.go     # Request logging, panic recovery
│       ├── persistence/sqlite/    # SQLite implementation
│       │   ├── schema.go
│       │   ├── agent_repository.go
│       │   ├── user_repository.go
│       │   ├── technique_repository.go
│       │   ├── scenario_repository.go
│       │   ├── result_repository.go
│       │   ├── notification_repository.go
│       │   └── schedule_repository.go
│       └── websocket/             # Agent communication
│           ├── hub.go             # Connection management
│           └── client.go          # Client handling
├── go.mod
└── go.sum
```

---

## Dependency Flow

```
Infrastructure → Application → Domain
     ↓               ↓           ↓
  Handlers      Services     Entities
  Repositories              Interfaces
  WebSocket
  Middleware
```

**Rule**: Dependencies always point inward toward Domain.

---

## API Endpoints

Base URL: `https://localhost:8443/api/v1`

### Health
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Server health check |

### Authentication (public, rate-limited)
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/auth/login` | Login (5 attempts/min per IP) |
| `POST` | `/auth/refresh` | Refresh token (10 attempts/min per IP) |
| `POST` | `/auth/logout` | Invalidate tokens |
| `GET` | `/auth/me` | Get current user info |

### Agents
| Method | Endpoint | Permission | Description |
|--------|----------|------------|-------------|
| `GET` | `/agents` | `agents:view` | List agents (`?all=true` for offline) |
| `GET` | `/agents/:paw` | `agents:view` | Get agent details |
| `POST` | `/agents` | `agents:create` | Register agent |
| `DELETE` | `/agents/:paw` | `agents:delete` | Delete agent |
| `POST` | `/agents/:paw/heartbeat` | `agents:view` | Update last_seen |

### Techniques
| Method | Endpoint | Permission | Description |
|--------|----------|------------|-------------|
| `GET` | `/techniques` | `techniques:view` | List all techniques |
| `GET` | `/techniques/:id` | `techniques:view` | Get technique by ID |
| `GET` | `/techniques/tactic/:tactic` | `techniques:view` | By tactic |
| `GET` | `/techniques/platform/:platform` | `techniques:view` | By platform |
| `GET` | `/techniques/coverage` | `techniques:view` | Coverage statistics |
| `GET` | `/techniques/:id/executors` | `techniques:view` | List executors (`?platform=`) |
| `POST` | `/techniques/import` | `techniques:import` | Import from YAML |

### Scenarios
| Method | Endpoint | Permission | Description |
|--------|----------|------------|-------------|
| `GET` | `/scenarios` | `scenarios:view` | List all scenarios |
| `GET` | `/scenarios/:id` | `scenarios:view` | Get scenario details |
| `GET` | `/scenarios/tag/:tag` | `scenarios:view` | By tag |
| `GET` | `/scenarios/export` | `scenarios:export` | Export all scenarios |
| `GET` | `/scenarios/:id/export` | `scenarios:export` | Export single scenario |
| `POST` | `/scenarios` | `scenarios:create` | Create scenario |
| `POST` | `/scenarios/import` | `scenarios:import` | Import scenarios |
| `PUT` | `/scenarios/:id` | `scenarios:edit` | Update scenario |
| `DELETE` | `/scenarios/:id` | `scenarios:delete` | Delete scenario |

### Executions
| Method | Endpoint | Permission | Description |
|--------|----------|------------|-------------|
| `GET` | `/executions` | `executions:view` | Recent executions (limit 50) |
| `GET` | `/executions/:id` | `executions:view` | Get execution details |
| `GET` | `/executions/:id/results` | `executions:view` | Get results |
| `POST` | `/executions` | `executions:start` | Start execution |
| `POST` | `/executions/:id/stop` | `executions:stop` | Stop execution |
| `POST` | `/executions/:id/complete` | `executions:view` | Complete execution |

### Analytics
| Method | Endpoint | Permission | Description |
|--------|----------|------------|-------------|
| `GET` | `/analytics/period` | `analytics:view` | Period stats |
| `GET` | `/analytics/comparison` | `analytics:compare` | Compare periods |
| `GET` | `/analytics/trend` | `analytics:view` | Score trend |
| `GET` | `/analytics/summary` | `analytics:view` | Execution summary |

### Notifications
| Method | Endpoint | Permission | Description |
|--------|----------|------------|-------------|
| `GET` | `/notifications` | authenticated | Get notifications |
| `GET` | `/notifications/unread/count` | authenticated | Unread count |
| `POST` | `/notifications/:id/read` | authenticated | Mark as read |
| `POST` | `/notifications/read-all` | authenticated | Mark all as read |
| `GET` | `/notifications/settings` | authenticated | Get settings |
| `POST` | `/notifications/settings` | authenticated | Create settings |
| `PUT` | `/notifications/settings/:id` | authenticated | Update settings |
| `DELETE` | `/notifications/settings/:id` | authenticated | Delete settings |
| `GET` | `/notifications/smtp` | admin | Get SMTP config |
| `POST` | `/notifications/smtp/test` | admin | Test SMTP connection |

### Schedules
| Method | Endpoint | Permission | Description |
|--------|----------|------------|-------------|
| `GET` | `/schedules` | `scheduler:view` | List schedules |
| `GET` | `/schedules/:id` | `scheduler:view` | Get schedule |
| `GET` | `/schedules/:id/runs` | `scheduler:view` | Get run history |
| `POST` | `/schedules` | `scheduler:create` | Create schedule |
| `PUT` | `/schedules/:id` | `scheduler:edit` | Update schedule |
| `DELETE` | `/schedules/:id` | `scheduler:delete` | Delete schedule |
| `POST` | `/schedules/:id/pause` | `scheduler:edit` | Pause schedule |
| `POST` | `/schedules/:id/resume` | `scheduler:edit` | Resume schedule |
| `POST` | `/schedules/:id/run` | `executions:start` | Run schedule now |

### Permissions
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/permissions/matrix` | Permission matrix for all roles |
| `GET` | `/permissions/me` | Current user permissions |

### Admin (requires admin role)
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/admin/users` | List all users |
| `GET` | `/admin/users/:id` | Get user by ID |
| `POST` | `/admin/users` | Create user |
| `PUT` | `/admin/users/:id` | Update user |
| `PUT` | `/admin/users/:id/role` | Update user role |
| `DELETE` | `/admin/users/:id` | Deactivate user |
| `POST` | `/admin/users/:id/reactivate` | Reactivate user |
| `POST` | `/admin/users/:id/reset-password` | Reset user password |

---

## Middleware

### Authentication (`auth.go`)
| Middleware | Description |
|-----------|-------------|
| `NoAuthMiddleware()` | Dev mode: sets anonymous user with admin role |
| `AuthMiddleware(config)` | JWT token validation and user context |
| `AgentAuthMiddleware(config)` | Agent auth via `X-Agent-Key` header |
| `RoleMiddleware(roles...)` | Role-based access control |
| `PermissionMiddleware(perms...)` | Permission check (requires ALL) |
| `RequireAnyPermission(perms...)` | Permission check (requires ANY) |

### Security Headers (`security.go`)
Adds production security headers:
- `Strict-Transport-Security` (HSTS)
- `Content-Security-Policy` (CSP)
- `X-Frame-Options`
- `X-Content-Type-Options`
- `X-XSS-Protection`
- `Referrer-Policy`
- `Permissions-Policy`

### Rate Limiting (`ratelimit.go`)
Per-IP rate limiting with automatic cleanup every 5 minutes:
- Login: 5 attempts/minute
- Token refresh: 10 attempts/minute
- Returns 429 Too Many Requests when exceeded

### Logging (`logging.go`)
- Structured request/response logging with zap
- Panic recovery middleware

---

## WebSocket Protocol

### Agent Connection
Endpoint: `wss://localhost:8443/ws/agent`

```json
// Agent → Server: Registration
{"type": "register", "payload": {"paw": "...", "hostname": "...", "platform": "...", "executors": [...]}}

// Server → Agent: Registered
{"type": "registered", "payload": {"status": "ok", "paw": "..."}}

// Agent → Server: Heartbeat (every 30s)
{"type": "heartbeat", "payload": {"paw": "..."}}

// Server → Agent: Task
{"type": "task", "payload": {"id": "...", "technique_id": "T1082", "command": "...", "executor": "cmd", "timeout": 300}}

// Agent → Server: Result
{"type": "task_result", "payload": {"task_id": "...", "technique_id": "...", "success": true, "output": "...", "exit_code": 0}}

// Server → Agent: Acknowledgment
{"type": "task_ack", "payload": {"task_id": "...", "status": "received"}}
```

### Dashboard Connection
Endpoint: `wss://localhost:8443/ws/dashboard`

```json
// Server broadcasts to all dashboards
{"type": "execution_started", "payload": {"execution_id": "...", "data": {...}}}
{"type": "execution_completed", "payload": {"execution_id": "...", "data": {...}}}
{"type": "execution_cancelled", "payload": {"execution_id": "...", "data": {...}}}

// Dashboard → Server: Ping
{"type": "ping", "payload": {}}
// Server → Dashboard: Pong
{"type": "pong", "payload": {}}
```

### Connection Parameters
| Parameter | Value | Description |
|-----------|-------|-------------|
| Ping interval | 54 seconds | Server sends ping frames |
| Pong timeout | 60 seconds | Connection closed if no pong |
| Max message size | 512 KB | Maximum frame size |
| Write timeout | 10 seconds | Maximum write duration |

---

## WebSocket Hub

The Hub manages all WebSocket connections:

```go
type Hub struct {
    clients   map[*Client]bool     // All connected clients
    agents    map[string]*Client   // Agents indexed by PAW
    broadcast chan []byte          // Broadcast channel
    register  chan *Client         // Registration channel
    unregister chan *Client        // Unregistration channel
}

// Key methods
func (h *Hub) SendToAgent(paw string, message []byte) bool
func (h *Hub) Broadcast(message []byte)
func (h *Hub) IsAgentConnected(paw string) bool
func (h *Hub) GetConnectedAgents() []string
func (h *Hub) SetOnAgentDisconnect(callback func(paw string))
```

---

## Security Score

**Formula**: `(blocked × 100 + detected × 50) / (total × 100) × 100`

| Status | Points | Description |
|--------|--------|-------------|
| Blocked | 100 | Attack prevented by defenses |
| Detected | 50 | Attack executed but detected |
| Success | 0 | Attack executed without detection |

Example: 5 techniques, 2 blocked, 2 detected, 1 successful
```
Score = (2×100 + 2×50) / (5×100) × 100 = 60%
```

---

## Domain Entities

### Agent
```go
type Agent struct {
    Paw       string
    Hostname  string
    Platform  string            // windows, linux, darwin
    Username  string
    Executors []string          // psh, cmd, bash, sh
    Status    AgentStatus       // online, offline, busy, untrusted
    LastSeen  time.Time
    IPAddress string
    OSVersion string
    Metadata  map[string]string
    CreatedAt time.Time
}
```

### Technique
```go
type Technique struct {
    ID          string
    Name        string
    Description string
    Tactic      TacticType     // Primary tactic (retro-compatible)
    Tactics     []TacticType   // All tactics (multi-tactic support)
    Platforms   []string
    Executors   []Executor
    Detection   []Detection
    References  []string       // MITRE ATT&CK URLs
    IsSafe      bool
}

type Executor struct {
    Name              string // Executor display name (optional)
    Type              string // cmd, powershell, bash, sh
    Platform          string // windows, linux, macos (optional)
    Command           string
    Cleanup           string
    Timeout           int
    ElevationRequired bool   // Needs admin/root privileges (optional)
}

type TechniqueSelection struct {
    TechniqueID  string // Technique ID
    ExecutorName string // Preferred executor (empty = auto-select)
}
```

### Execution
```go
type Execution struct {
    ID          string
    ScenarioID  string
    Status      ExecutionStatus // pending, running, completed, failed, cancelled
    StartedAt   time.Time
    CompletedAt *time.Time
    SafeMode    bool
    Score       *SecurityScore
}
```

### ExecutionResult
```go
type ExecutionResult struct {
    ID          string
    ExecutionID string
    TechniqueID string
    AgentPaw    string
    Status      ResultStatus // pending, success, blocked, detected, failed
    Output      string
    ExitCode    int
    StartedAt   time.Time
    CompletedAt *time.Time
}
```

### User
```go
type User struct {
    ID           string
    Username     string
    Email        string
    PasswordHash string      // Never exposed in JSON
    Role         UserRole    // admin, rssi, operator, analyst, viewer
    IsActive     bool
    LastLoginAt  *time.Time
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

### Notification
```go
type Notification struct {
    ID        string
    UserID    string
    Type      NotificationType // execution_started, execution_completed, execution_failed, score_alert, agent_offline
    Title     string
    Message   string
    Data      map[string]any
    Read      bool
    SentAt    *time.Time
    CreatedAt time.Time
}

type NotificationSettings struct {
    ID                   string
    UserID               string
    Channel              NotificationChannel // email, webhook
    Enabled              bool
    EmailAddress         string
    WebhookURL           string
    NotifyOnStart        bool
    NotifyOnComplete     bool
    NotifyOnFailure      bool
    NotifyOnScoreAlert   bool
    ScoreAlertThreshold  float64
    NotifyOnAgentOffline bool
    CreatedAt            time.Time
    UpdatedAt            time.Time
}
```

### Schedule
```go
type Schedule struct {
    ID          string
    Name        string
    Description string
    ScenarioID  string
    AgentPaw    string              // empty = any available
    Frequency   ScheduleFrequency   // once, hourly, daily, weekly, monthly, cron
    CronExpr    string              // only for cron frequency
    SafeMode    bool
    Status      ScheduleStatus      // active, paused, disabled
    NextRunAt   *time.Time
    LastRunAt   *time.Time
    LastRunID   string
    CreatedBy   string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type ScheduleRun struct {
    ID          string
    ScheduleID  string
    ExecutionID string
    StartedAt   time.Time
    CompletedAt *time.Time
    Status      string    // pending, running, completed, failed
    Error       string
}
```

### Permission
```go
type Permission string
// 28 permissions across 10 categories:
// users:view, users:create, users:edit, users:delete
// agents:view, agents:create, agents:delete
// techniques:view, techniques:import
// scenarios:view, scenarios:create, scenarios:edit, scenarios:delete, scenarios:import, scenarios:export
// executions:view, executions:start, executions:stop
// analytics:view, analytics:compare, analytics:export
// settings:view, settings:edit
// scheduler:view, scheduler:create, scheduler:edit, scheduler:delete
```

---

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_PATH` | SQLite database path | `./data/autostrike.db` |
| `DASHBOARD_PATH` | Dashboard dist folder | `../dashboard/dist` |
| `JWT_SECRET` | JWT signing key (enables auth when set) | - (auth disabled) |
| `ENABLE_AUTH` | Explicit auth override (`true`/`false`) | - |
| `AGENT_SECRET` | Agent authentication secret | - |
| `DEFAULT_ADMIN_PASSWORD` | Initial admin password | Random |
| `ALLOWED_ORIGINS` | CORS origins | `localhost:3000,localhost:8443` |
| `LOG_LEVEL` | Logging level | `info` |

### SMTP Configuration (optional)

| Variable | Description | Default |
|----------|-------------|---------|
| `SMTP_HOST` | Mail server hostname | - |
| `SMTP_PORT` | Mail server port | `587` |
| `SMTP_USERNAME` | SMTP username | - |
| `SMTP_PASSWORD` | SMTP password | - |
| `SMTP_FROM` | Sender email address | - |
| `SMTP_USE_TLS` | Use TLS | `false` |
| `DASHBOARD_URL` | Dashboard URL for email links | `https://localhost:8443` |

### Authentication Behavior

| Configuration | Auth Status |
|--------------|-------------|
| `JWT_SECRET` not set | Auth **disabled** (development mode) |
| `JWT_SECRET` set | Auth **enabled** (production mode) |
| `ENABLE_AUTH=false` | Auth **disabled** (explicit override) |
| `ENABLE_AUTH=true` | Auth **enabled** (explicit override) |

---

## Testing

200+ tests with comprehensive coverage:

| Package | Coverage |
|---------|----------|
| **application** | 83.0% |
| **entity** | 95.0% |
| **service** | 99.2% |
| **handlers** | 87.5% |
| **websocket** | 91.6% |
| **middleware** | 94.3% |
| **rest/server** | 87.9% |
| **sqlite** | 85.0% |

```bash
cd server
go test ./...              # Run all tests
go test ./... -cover       # With coverage
go test ./... -v           # Verbose output
```

---

## Running

```bash
# Development
go run ./cmd/autostrike

# Production build
go build -o autostrike ./cmd/autostrike
./autostrike

# With authentication
JWT_SECRET=secret ./autostrike

# With full configuration
JWT_SECRET=secret AGENT_SECRET=agent-key SMTP_HOST=mail.example.com ./autostrike
```
