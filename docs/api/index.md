# API

Documentation de l'API REST AutoStrike.

---

## Base URL

```
https://server:8443/api/v1
```

---

## Authentification

L'API utilise des tokens JWT.

```bash
curl -X POST https://server:8443/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password"}'
```

Réponse :

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_at": "2024-01-01T00:00:00Z"
}
```

Utilisez le token dans le header `Authorization` :

```bash
curl https://server:8443/api/v1/agents \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

---

## Endpoints

Voir la [Référence API](reference.md) pour la liste complète.
