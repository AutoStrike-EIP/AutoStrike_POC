# Backend (Go)

Le serveur de contrôle AutoStrike est développé en **Go 1.21** avec le framework **Gin**.

---

## Structure

```
server/
├── cmd/
│   └── autostrike/
│       └── main.go              # Point d'entrée
├── internal/
│   ├── domain/                  # Entités métier
│   │   ├── agent.go
│   │   ├── scenario.go
│   │   ├── technique.go
│   │   └── result.go
│   ├── service/                 # Logique métier
│   │   ├── orchestrator.go
│   │   ├── validator.go
│   │   └── score.go
│   ├── repository/              # Accès données
│   │   ├── sqlite/
│   │   └── postgres/
│   ├── api/                     # Handlers HTTP
│   │   ├── routes.go
│   │   ├── agents.go
│   │   ├── scenarios.go
│   │   └── websocket.go
│   └── config/
│       └── config.go
├── go.mod
└── go.sum
```

---

## Configuration

```yaml
# config.yaml
server:
  host: 0.0.0.0
  port: 8443
  tls:
    cert: /etc/autostrike/server.crt
    key: /etc/autostrike/server.key

database:
  driver: sqlite
  dsn: /var/lib/autostrike/autostrike.db

jwt:
  secret: ${JWT_SECRET}
  expiry: 24h

agents:
  heartbeat_interval: 30s
  timeout: 90s
```

---

## API Endpoints

| Méthode | Endpoint | Description |
|---------|----------|-------------|
| `GET` | `/api/v1/agents` | Liste des agents |
| `GET` | `/api/v1/scenarios` | Liste des scénarios |
| `POST` | `/api/v1/scenarios/{id}/execute` | Exécuter un scénario |
| `GET` | `/api/v1/results` | Résultats des exécutions |
| `WS` | `/ws` | WebSocket temps réel |
