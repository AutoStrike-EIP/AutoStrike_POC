# Référence API

## Agents

### Lister les agents

```http
GET /api/v1/agents
```

**Réponse :**

```json
{
  "agents": [
    {
      "paw": "AGENT_001",
      "hostname": "WORKSTATION-01",
      "platform": "windows",
      "status": "online",
      "last_seen": "2024-01-01T12:00:00Z"
    }
  ]
}
```

---

## Scénarios

### Lister les scénarios

```http
GET /api/v1/scenarios
```

### Exécuter un scénario

```http
POST /api/v1/scenarios/{id}/execute
```

**Body :**

```json
{
  "agent_paws": ["AGENT_001", "AGENT_002"],
  "safe_mode": true
}
```

**Réponse :**

```json
{
  "execution_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "running",
  "started_at": "2024-01-01T12:00:00Z"
}
```

---

## Résultats

### Obtenir les résultats

```http
GET /api/v1/executions/{id}/results
```

**Réponse :**

```json
{
  "execution_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "completed",
  "score": 75,
  "results": [
    {
      "technique_id": "T1082",
      "technique_name": "System Information Discovery",
      "status": "detected",
      "agent_paw": "AGENT_001",
      "executed_at": "2024-01-01T12:00:05Z"
    }
  ]
}
```

---

## WebSocket

### Connexion

```javascript
const ws = new WebSocket('wss://server:8443/ws');

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log(data);
};
```

### Événements

| Type | Description |
|------|-------------|
| `agent.status` | Changement de statut d'un agent |
| `execution.progress` | Progression d'une exécution |
| `execution.complete` | Fin d'une exécution |
| `technique.result` | Résultat d'une technique |
