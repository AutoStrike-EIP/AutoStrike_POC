# Référence API

Base URL: `https://localhost:8443/api/v1`

---

## Authentification

### Tokens JWT

Toutes les requêtes API nécessitent un token JWT dans le header Authorization :

```http
Authorization: Bearer <token>
```

Le token JWT est signé avec le secret défini dans `JWT_SECRET` et contient :
- `sub` : ID de l'utilisateur
- `role` : Rôle de l'utilisateur (admin, operator, viewer)
- `exp` : Date d'expiration du token

### Authentification Agent

Les agents utilisent un header spécifique :

```http
X-Agent-Key: <agent_secret>
```

Le secret agent est défini dans la variable d'environnement `AGENT_SECRET`.

### Génération de token (développement)

Pour les tests, générer un token JWT avec :

```bash
# Générer un token avec openssl et jq
SECRET="your-jwt-secret"
HEADER=$(echo -n '{"alg":"HS256","typ":"JWT"}' | base64 -w0 | tr '/+' '_-' | tr -d '=')
PAYLOAD=$(echo -n '{"sub":"admin","role":"admin","exp":'$(($(date +%s) + 86400))'}' | base64 -w0 | tr '/+' '_-' | tr -d '=')
SIGNATURE=$(echo -n "${HEADER}.${PAYLOAD}" | openssl dgst -sha256 -hmac "${SECRET}" -binary | base64 -w0 | tr '/+' '_-' | tr -d '=')
echo "${HEADER}.${PAYLOAD}.${SIGNATURE}"
```

---

## Agents

### Lister les agents

```http
GET /api/v1/agents
```

**Réponse :**

```json
[
  {
    "paw": "agent-001",
    "hostname": "WORKSTATION-01",
    "username": "admin",
    "platform": "windows",
    "executors": ["powershell", "cmd"],
    "status": "online",
    "last_seen": "2024-01-01T12:00:00Z",
    "created_at": "2024-01-01T10:00:00Z"
  }
]
```

### Obtenir un agent

```http
GET /api/v1/agents/:paw
```

### Enregistrer un agent

```http
POST /api/v1/agents
```

**Body :**

```json
{
  "paw": "agent-001",
  "hostname": "WORKSTATION-01",
  "username": "admin",
  "platform": "windows",
  "executors": ["powershell", "cmd"]
}
```

### Supprimer un agent

```http
DELETE /api/v1/agents/:paw
```

### Heartbeat

```http
POST /api/v1/agents/:paw/heartbeat
```

---

## Techniques

### Lister les techniques

```http
GET /api/v1/techniques
```

**Réponse :**

```json
[
  {
    "id": "T1082",
    "name": "System Information Discovery",
    "description": "Adversaries may attempt to get detailed information...",
    "tactic": "discovery",
    "platforms": ["windows", "linux", "darwin"],
    "executors": [
      {
        "type": "cmd",
        "command": "systeminfo",
        "cleanup": "",
        "timeout": 60
      }
    ],
    "detection": [
      {
        "source": "Process Creation",
        "indicator": "systeminfo.exe execution"
      }
    ],
    "is_safe": true
  }
]
```

### Obtenir une technique

```http
GET /api/v1/techniques/:id
```

### Techniques par tactique

```http
GET /api/v1/techniques/tactic/:tactic
```

Tactiques MITRE disponibles :
- `reconnaissance`
- `resource-development`
- `initial-access`
- `execution`
- `persistence`
- `privilege-escalation`
- `defense-evasion`
- `credential-access`
- `discovery`
- `lateral-movement`
- `collection`
- `command-and-control`
- `exfiltration`
- `impact`

### Techniques par plateforme

```http
GET /api/v1/techniques/platform/:platform
```

Plateformes : `windows`, `linux`, `darwin`

### Couverture MITRE

```http
GET /api/v1/techniques/coverage
```

**Réponse :**

```json
{
  "discovery": 15,
  "execution": 8,
  "persistence": 12,
  "defense-evasion": 20
}
```

### Importer des techniques

```http
POST /api/v1/techniques/import
```

**Body :**

```json
{
  "path": "/path/to/techniques.yaml"
}
```

---

## Exécutions

### Lister les exécutions récentes

```http
GET /api/v1/executions
```

**Réponse :**

```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "scenario_id": "scenario-001",
    "status": "completed",
    "start_time": "2024-01-01T12:00:00Z",
    "end_time": "2024-01-01T12:05:00Z",
    "safe_mode": true,
    "score": {
      "overall": 75.0,
      "blocked": 3,
      "detected": 2,
      "successful": 1,
      "total": 6
    }
  }
]
```

### Obtenir une exécution

```http
GET /api/v1/executions/:id
```

### Résultats d'une exécution

```http
GET /api/v1/executions/:id/results
```

**Réponse :**

```json
[
  {
    "id": "result-uuid",
    "execution_id": "550e8400-e29b-41d4-a716-446655440000",
    "technique_id": "T1082",
    "agent_paw": "agent-001",
    "status": "detected",
    "output": "Host Name: WORKSTATION-01...",
    "detected": true,
    "start_time": "2024-01-01T12:00:05Z",
    "end_time": "2024-01-01T12:00:10Z"
  }
]
```

### Démarrer une exécution

```http
POST /api/v1/executions
```

**Body :**

```json
{
  "scenario_id": "scenario-001",
  "agent_paws": ["agent-001", "agent-002"],
  "safe_mode": true
}
```

**Réponse :**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "scenario_id": "scenario-001",
  "status": "running",
  "start_time": "2024-01-01T12:00:00Z",
  "safe_mode": true
}
```

### Terminer une exécution

```http
POST /api/v1/executions/:id/complete
```

---

## WebSocket (Agents)

### Connexion

```
wss://localhost:8443/ws/agent
```

### Messages Agent → Server

**Enregistrement :**
```json
{
  "type": "register",
  "payload": {
    "paw": "agent-001",
    "hostname": "WORKSTATION-01",
    "username": "admin",
    "platform": "windows",
    "executors": ["powershell", "cmd"]
  }
}
```

**Heartbeat :**
```json
{
  "type": "heartbeat",
  "payload": {
    "paw": "agent-001"
  }
}
```

**Résultat de tâche :**
```json
{
  "type": "task_result",
  "payload": {
    "task_id": "task-uuid",
    "technique_id": "T1082",
    "success": true,
    "output": "Host Name: WORKSTATION-01...",
    "exit_code": 0
  }
}
```

### Messages Server → Agent

**Tâche à exécuter :**
```json
{
  "type": "task",
  "payload": {
    "id": "task-uuid",
    "technique_id": "T1082",
    "command": "systeminfo",
    "executor": "cmd",
    "timeout": 300,
    "cleanup": ""
  }
}
```

**Ping :**
```json
{
  "type": "ping",
  "payload": {}
}
```

---

## Codes d'erreur

| Code | Description |
|------|-------------|
| 400 | Requête invalide |
| 401 | Non authentifié |
| 403 | Accès refusé |
| 404 | Ressource non trouvée |
| 500 | Erreur serveur |

**Format d'erreur :**

```json
{
  "error": "description de l'erreur"
}
```
