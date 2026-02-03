# API

Documentation de l'API REST AutoStrike.

---

## Base URL

```
https://server:8443/api/v1
```

---

## Authentification

L'authentification JWT est **optionnelle** par défaut. Elle est activée uniquement si `JWT_SECRET` est défini dans `.env`.

### Sans authentification (développement)

```bash
curl https://localhost:8443/api/v1/agents
```

### Avec authentification (production)

Quand `JWT_SECRET` est défini, incluez le token dans le header `Authorization` :

```bash
curl https://localhost:8443/api/v1/agents \
  -H "Authorization: Bearer <your-jwt-token>"
```

> **Note**: Il n'y a pas d'endpoint `/auth/login`. Les tokens JWT doivent être générés par votre système d'authentification externe ou via les outils de développement.

---

## Endpoints principaux

| Méthode | Endpoint | Description |
|---------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/agents` | Liste des agents |
| GET | `/techniques` | Liste des techniques MITRE |
| GET | `/scenarios` | Liste des scénarios |
| POST | `/executions` | Lancer une exécution |
| GET | `/executions/:id` | Détails d'une exécution |

---

## WebSocket

| Path | Description |
|------|-------------|
| `/ws/agent` | Connexion agent |
| `/ws/dashboard` | Mises à jour temps réel |

---

## Documentation complète

Voir la [Référence API](reference.md) pour tous les endpoints, paramètres et exemples.
