# Backend (Go)

Le serveur de contrôle AutoStrike est développé en **Go 1.21** avec le framework **Gin** et une architecture **hexagonale**.

---

## Architecture Hexagonale

Le serveur suit le pattern **Ports & Adapters** pour une séparation claire des responsabilités :

```
server/
├── cmd/autostrike/
│   └── main.go                    # Point d'entrée, DI
├── config/
│   └── config.yaml                # Configuration YAML
├── internal/
│   ├── domain/                    # 🟢 Couche Métier (indépendante)
│   │   ├── entity/                # Entités: Agent, Technique, Scenario, Result
│   │   ├── repository/            # Interfaces (ports sortants)
│   │   ├── service/               # Services domaine
│   │   │   ├── orchestrator.go    # Orchestration des attaques
│   │   │   ├── validator.go       # Validation compatibilité
│   │   │   └── score_calculator.go # Calcul des scores
│   │   └── valueobject/           # Objets valeur
│   ├── application/               # 🟡 Cas d'utilisation
│   │   ├── agent_service.go       # CRUD agents, heartbeat
│   │   ├── execution_service.go   # Démarrage/suivi exécutions
│   │   ├── scenario_service.go    # Gestion scénarios
│   │   └── technique_service.go   # Catalogue techniques
│   └── infrastructure/            # 🔵 Adaptateurs externes
│       ├── api/rest/
│       │   └── server.go          # Serveur REST Gin
│       ├── http/
│       │   ├── handlers/          # Handlers HTTP
│       │   │   ├── agent_handler.go
│       │   │   ├── technique_handler.go
│       │   │   └── execution_handler.go
│       │   └── middleware/        # Auth JWT, Logging
│       │       ├── auth.go
│       │       └── logging.go
│       ├── persistence/sqlite/    # Implémentation SQLite
│       │   ├── schema.go
│       │   ├── agent_repository.go
│       │   ├── technique_repository.go
│       │   ├── scenario_repository.go
│       │   └── result_repository.go
│       └── websocket/             # Communication agents
│           ├── hub.go
│           └── client.go
├── go.mod
└── go.sum
```

---

## Flux de Dépendances

```
Infrastructure → Application → Domain
     ↓               ↓           ↓
  Handlers      Services     Entities
  Repositories              Interfaces
  WebSocket
```

**Règle** : Les dépendances pointent toujours vers le centre (Domain).

---

## Configuration

```yaml
# config/config.yaml
server:
  host: "0.0.0.0"
  port: 8443
  mode: "release"

database:
  driver: "sqlite3"
  path: "./data/autostrike.db"

security:
  jwt_secret: "${JWT_SECRET}"
  agent_secret: "${AGENT_SECRET}"
  tls:
    enabled: true
    cert_file: "./certs/server.crt"
    key_file: "./certs/server.key"
    ca_file: "./certs/ca.crt"
    mtls: true

agent:
  heartbeat_interval: 30
  stale_timeout: 120

execution:
  default_timeout: 300
  max_concurrent: 10
  safe_mode_default: true
```

---

## API REST

Base URL: `https://localhost:8443/api/v1`

### Agents
| Méthode | Endpoint | Description |
|---------|----------|-------------|
| `GET` | `/agents` | Liste tous les agents |
| `GET` | `/agents/:paw` | Détails d'un agent |
| `POST` | `/agents` | Enregistrer un agent |
| `DELETE` | `/agents/:paw` | Supprimer un agent |
| `POST` | `/agents/:paw/heartbeat` | Heartbeat |

### Techniques
| Méthode | Endpoint | Description |
|---------|----------|-------------|
| `GET` | `/techniques` | Liste techniques MITRE |
| `GET` | `/techniques/:id` | Détails technique |
| `GET` | `/techniques/tactic/:tactic` | Par tactique |
| `GET` | `/techniques/coverage` | Statistiques couverture |
| `POST` | `/techniques/import` | Import YAML |

### Exécutions
| Méthode | Endpoint | Description |
|---------|----------|-------------|
| `GET` | `/executions` | Exécutions récentes |
| `GET` | `/executions/:id` | Détails exécution |
| `GET` | `/executions/:id/results` | Résultats |
| `POST` | `/executions` | Démarrer exécution |

---

## WebSocket

Endpoint: `wss://localhost:8443/ws/agent`

### Protocole

```json
// Agent → Server : Enregistrement
{"type": "register", "payload": {"paw": "...", "hostname": "...", "platform": "...", "executors": [...]}}

// Agent → Server : Heartbeat
{"type": "heartbeat", "payload": {"paw": "..."}}

// Server → Agent : Tâche
{"type": "task", "payload": {"id": "...", "technique_id": "T1082", "command": "...", "timeout": 300}}

// Agent → Server : Résultat
{"type": "task_result", "payload": {"task_id": "...", "success": true, "output": "...", "exit_code": 0}}
```

---

## Score de Sécurité

**Formule** : `(blocked × 100 + detected × 50) / (total × 100) × 100`

| Statut | Points | Description |
|--------|--------|-------------|
| Blocked | 100 | Technique bloquée par les défenses |
| Detected | 50 | Technique détectée mais exécutée |
| Successful | 0 | Technique exécutée sans détection |

---

## Lancement

```bash
# Développement
go run ./cmd/autostrike

# Production
go build -o autostrike ./cmd/autostrike
./autostrike
```
