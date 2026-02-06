# AutoStrike Server

Backend Go du projet AutoStrike - Plateforme BAS (Breach and Attack Simulation).

## Architecture

Le serveur suit une **architecture hexagonale** (Ports & Adapters) :

```
server/
├── cmd/autostrike/          # Point d'entrée
│   └── main.go
├── configs/
│   └── techniques/          # Définitions YAML (13 fichiers)
│       ├── reconnaissance.yaml
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
│   ├── domain/              # Couche métier (indépendante)
│   │   ├── entity/          # Entités: Agent, Technique, Scenario, Execution, Result, User, Notification, Schedule, Permission
│   │   ├── repository/      # Interfaces (ports sortants)
│   │   └── service/         # Services domaine: Orchestrator, Validator, ScoreCalculator
│   ├── application/         # Cas d'utilisation
│   │   ├── agent_service.go
│   │   ├── auth_service.go
│   │   ├── execution_service.go
│   │   ├── scenario_service.go
│   │   ├── technique_service.go
│   │   ├── notification_service.go
│   │   ├── schedule_service.go
│   │   ├── analytics_service.go
│   │   └── token_blacklist.go
│   └── infrastructure/      # Adaptateurs externes
│       ├── api/rest/        # Serveur REST (Gin)
│       ├── http/
│       │   ├── handlers/    # Handlers HTTP (11 fichiers)
│       │   └── middleware/  # Auth JWT, Security Headers, Rate Limiting, Logging
│       ├── persistence/
│       │   └── sqlite/      # Implémentation SQLite
│       └── websocket/       # Communication agents
└── go.mod
```

## Prérequis

- Go 1.24+
- SQLite3
- OpenSSL (pour les certificats)

## Installation

```bash
# Installer les dépendances
go mod download

# Générer les certificats TLS
cd .. && ./scripts/generate-certs.sh ./certs

# Créer le dossier data
mkdir -p data
```

## Configuration

Variables d'environnement (fichier `.env`) :

```env
# Base de données
DATABASE_PATH=./data/autostrike.db

# Dashboard
DASHBOARD_PATH=../dashboard/dist

# Authentification JWT (optionnel - auth désactivée si non défini)
JWT_SECRET=your-secure-jwt-secret-32-characters
ENABLE_AUTH=true

# Agent
AGENT_SECRET=your-secure-agent-secret

# Admin
DEFAULT_ADMIN_PASSWORD=your-admin-password

# CORS
ALLOWED_ORIGINS=localhost:3000,localhost:8443

# Logging
LOG_LEVEL=info

# SMTP (optionnel - pour notifications email)
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=notifications@example.com
SMTP_PASSWORD=smtp-password
SMTP_FROM=noreply@example.com
SMTP_USE_TLS=true
DASHBOARD_URL=https://localhost:8443
```

## Lancement

```bash
# Mode développement
go run ./cmd/autostrike

# Mode production
go build -o autostrike ./cmd/autostrike
JWT_SECRET=secret AGENT_SECRET=agent-key ./autostrike
```

## API REST

Base URL: `https://localhost:8443/api/v1`

### Authentication (routes publiques, rate-limited)
| Méthode | Endpoint | Description |
|---------|----------|-------------|
| POST | `/auth/login` | Login (5 tentatives/min par IP) |
| POST | `/auth/refresh` | Rafraîchir le token (10 tentatives/min par IP) |
| POST | `/auth/logout` | Invalider les tokens |
| GET | `/auth/me` | Infos utilisateur courant |

### Agents
| Méthode | Endpoint | Description |
|---------|----------|-------------|
| GET | `/agents` | Liste des agents (`?all=true` pour offline) |
| GET | `/agents/:paw` | Détails d'un agent |
| POST | `/agents` | Enregistrer un agent |
| DELETE | `/agents/:paw` | Supprimer un agent |
| POST | `/agents/:paw/heartbeat` | Heartbeat agent |

### Techniques
| Méthode | Endpoint | Description |
|---------|----------|-------------|
| GET | `/techniques` | Liste des 48 techniques |
| GET | `/techniques/:id` | Détails d'une technique |
| GET | `/techniques/tactic/:tactic` | Techniques par tactique MITRE |
| GET | `/techniques/platform/:platform` | Techniques par plateforme |
| GET | `/techniques/coverage` | Statistiques de couverture (13 tactiques) |
| POST | `/techniques/import` | Importer depuis YAML |

### Scénarios
| Méthode | Endpoint | Description |
|---------|----------|-------------|
| GET | `/scenarios` | Liste des scénarios |
| GET | `/scenarios/:id` | Détails d'un scénario |
| GET | `/scenarios/tag/:tag` | Scénarios par tag |
| GET | `/scenarios/export` | Exporter tous les scénarios |
| GET | `/scenarios/:id/export` | Exporter un scénario |
| POST | `/scenarios` | Créer un scénario |
| POST | `/scenarios/import` | Importer des scénarios |
| PUT | `/scenarios/:id` | Modifier un scénario |
| DELETE | `/scenarios/:id` | Supprimer un scénario |

### Exécutions
| Méthode | Endpoint | Description |
|---------|----------|-------------|
| GET | `/executions` | Liste des exécutions récentes |
| GET | `/executions/:id` | Détails d'une exécution |
| GET | `/executions/:id/results` | Résultats d'une exécution |
| POST | `/executions` | Démarrer une exécution |
| POST | `/executions/:id/stop` | Arrêter une exécution |
| POST | `/executions/:id/complete` | Terminer une exécution |

### Schedules (Planification)
| Méthode | Endpoint | Description |
|---------|----------|-------------|
| GET | `/schedules` | Liste des planifications |
| GET | `/schedules/:id` | Détails d'une planification |
| GET | `/schedules/:id/runs` | Historique des exécutions |
| POST | `/schedules` | Créer une planification |
| PUT | `/schedules/:id` | Modifier une planification |
| DELETE | `/schedules/:id` | Supprimer une planification |
| POST | `/schedules/:id/pause` | Mettre en pause |
| POST | `/schedules/:id/resume` | Reprendre |
| POST | `/schedules/:id/run` | Exécuter maintenant |

### Notifications
| Méthode | Endpoint | Description |
|---------|----------|-------------|
| GET | `/notifications` | Liste des notifications |
| GET | `/notifications/unread/count` | Nombre de non-lues |
| POST | `/notifications/:id/read` | Marquer comme lue |
| POST | `/notifications/read-all` | Tout marquer comme lu |
| GET | `/notifications/settings` | Paramètres de notification |
| POST | `/notifications/settings` | Créer des paramètres |
| PUT | `/notifications/settings/:id` | Modifier les paramètres |
| DELETE | `/notifications/settings/:id` | Supprimer les paramètres |
| GET | `/notifications/smtp` | Configuration SMTP (admin) |
| POST | `/notifications/smtp/test` | Tester SMTP (admin) |

### Analytics
| Méthode | Endpoint | Description |
|---------|----------|-------------|
| GET | `/analytics/period` | Statistiques par période |
| GET | `/analytics/comparison` | Comparaison de périodes |
| GET | `/analytics/trend` | Tendance du score |
| GET | `/analytics/summary` | Résumé des exécutions |

### Admin (Utilisateurs) - rôle admin requis
| Méthode | Endpoint | Description |
|---------|----------|-------------|
| GET | `/admin/users` | Liste des utilisateurs |
| GET | `/admin/users/:id` | Détails d'un utilisateur |
| POST | `/admin/users` | Créer un utilisateur |
| PUT | `/admin/users/:id` | Modifier un utilisateur |
| PUT | `/admin/users/:id/role` | Changer le rôle |
| DELETE | `/admin/users/:id` | Désactiver un utilisateur |
| POST | `/admin/users/:id/reactivate` | Réactiver un utilisateur |
| POST | `/admin/users/:id/reset-password` | Réinitialiser le mot de passe |

### Permissions
| Méthode | Endpoint | Description |
|---------|----------|-------------|
| GET | `/permissions/matrix` | Matrice des permissions |
| GET | `/permissions/me` | Mes permissions |

## WebSocket

| Endpoint | Description |
|----------|-------------|
| `wss://localhost:8443/ws/agent` | Connexion agents |
| `wss://localhost:8443/ws/dashboard` | Mises à jour temps réel |

### Messages Agent → Server
```json
{"type": "register", "payload": {"paw": "...", "hostname": "...", "platform": "...", "executors": [...]}}
{"type": "heartbeat", "payload": {"paw": "..."}}
{"type": "task_result", "payload": {"task_id": "...", "success": true, "output": "..."}}
```

### Messages Server → Agent
```json
{"type": "registered", "payload": {"status": "ok", "paw": "..."}}
{"type": "task", "payload": {"id": "...", "technique_id": "T1082", "command": "...", "timeout": 300}}
{"type": "task_ack", "payload": {"task_id": "...", "status": "received"}}
```

## Score de Sécurité

Formule : `(blocked × 100 + detected × 50) / (total × 100) × 100`

| Statut | Points |
|--------|--------|
| Blocked | 100 |
| Detected | 50 |
| Successful | 0 |

## Tests

200+ tests avec couverture complète :

```bash
go test -v ./...
go test ./... -cover
```

## Docker

```bash
docker build -t autostrike-server .
docker run -p 8443:8443 -v $(pwd)/certs:/app/certs autostrike-server
```
