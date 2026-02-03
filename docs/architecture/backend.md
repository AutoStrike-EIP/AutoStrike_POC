# Backend (Go)

The AutoStrike control server is built with **Go 1.21+** using the **Gin** framework and **hexagonal architecture**.

---

## Hexagonal Architecture

The server follows the **Ports & Adapters** pattern for clear separation of concerns:

```
server/
├── cmd/autostrike/
│   └── main.go                    # Entry point, DI, startup
├── configs/
│   └── techniques/                # YAML technique definitions
│       ├── discovery.yaml
│       ├── execution.yaml
│       ├── persistence.yaml
│       └── defense-evasion.yaml
├── internal/
│   ├── domain/                    # 🟢 Business Layer (independent)
│   │   ├── entity/                # Entities: Agent, Technique, Scenario, Execution, Result
│   │   ├── repository/            # Interfaces (outbound ports)
│   │   └── service/               # Domain services
│   │       ├── orchestrator.go    # Attack orchestration
│   │       ├── validator.go       # Compatibility validation
│   │       └── score_calculator.go # Security score calculation
│   ├── application/               # 🟡 Use Cases
│   │   ├── agent_service.go       # Agent CRUD, heartbeat
│   │   ├── execution_service.go   # Execution lifecycle
│   │   ├── scenario_service.go    # Scenario management
│   │   └── technique_service.go   # Technique catalog
│   └── infrastructure/            # 🔵 External Adapters
│       ├── api/rest/
│       │   └── server.go          # Gin REST server
│       ├── http/
│       │   ├── handlers/          # HTTP handlers
│       │   │   ├── agent_handler.go
│       │   │   ├── technique_handler.go
│       │   │   ├── scenario_handler.go
│       │   │   ├── execution_handler.go
│       │   │   └── websocket_handler.go
│       │   └── middleware/        # JWT Auth, Logging
│       │       ├── auth.go
│       │       └── logging.go
│       ├── persistence/sqlite/    # SQLite implementation
│       │   ├── schema.go
│       │   ├── agent_repository.go
│       │   ├── technique_repository.go
│       │   ├── scenario_repository.go
│       │   └── result_repository.go
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
```

**Rule**: Dependencies always point inward toward Domain.

---

## API Endpoints

Base URL: `https://localhost:8443/api/v1`

### Health
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Server health check |

### Agents
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/agents` | List agents (`?all=true` for offline) |
| `GET` | `/agents/:paw` | Get agent details |
| `POST` | `/agents` | Register agent |
| `DELETE` | `/agents/:paw` | Delete agent |
| `POST` | `/agents/:paw/heartbeat` | Update last_seen |

### Techniques
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/techniques` | List all techniques |
| `GET` | `/techniques/:id` | Get technique by ID |
| `GET` | `/techniques/tactic/:tactic` | By tactic |
| `GET` | `/techniques/platform/:platform` | By platform |
| `GET` | `/techniques/coverage` | Coverage statistics |
| `POST` | `/techniques/import` | Import from YAML |

### Scenarios
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/scenarios` | List all scenarios |
| `GET` | `/scenarios/:id` | Get scenario details |
| `GET` | `/scenarios/tag/:tag` | By tag |
| `POST` | `/scenarios` | Create scenario |
| `PUT` | `/scenarios/:id` | Update scenario |
| `DELETE` | `/scenarios/:id` | Delete scenario |

### Executions
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/executions` | Recent executions (limit 50) |
| `GET` | `/executions/:id` | Get execution details |
| `GET` | `/executions/:id/results` | Get results |
| `POST` | `/executions` | Start execution |
| `POST` | `/executions/:id/stop` | Stop execution |
| `POST` | `/executions/:id/complete` | Complete execution |

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
    Platform  string      // windows, linux, darwin
    Username  string
    Executors []string    // psh, cmd, bash, sh
    Status    AgentStatus // online, offline, busy
    LastSeen  time.Time
}
```

### Technique
```go
type Technique struct {
    ID          string
    Name        string
    Description string
    Tactic      TacticType
    Platforms   []string
    Executors   []Executor
    Detection   []Detection
    IsSafe      bool
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

---

## Configuration

Environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_PATH` | SQLite database path | `./data/autostrike.db` |
| `DASHBOARD_PATH` | Dashboard dist folder | `../dashboard/dist` |
| `JWT_SECRET` | JWT signing key | - (auth disabled if not set) |
| `AGENT_SECRET` | Agent authentication | - |
| `ALLOWED_ORIGINS` | CORS origins | `localhost:3000,localhost:8443` |
| `LOG_LEVEL` | Logging level | `info` |

---

## Testing

Test coverage:
- **application**: 100%
- **entity**: 100%
- **service**: 99.2%
- **handlers**: 97.2%
- **websocket**: 91.6%
- **middleware**: 100%

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

# With environment
JWT_SECRET=secret ./autostrike
```
