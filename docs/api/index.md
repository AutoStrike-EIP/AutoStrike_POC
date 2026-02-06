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

Quand `JWT_SECRET` est défini, utilisez l'endpoint `/auth/login` pour obtenir un token :

```bash
# Login
curl -X POST https://localhost:8443/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'

# Utiliser le token
curl https://localhost:8443/api/v1/agents \
  -H "Authorization: Bearer <access_token>"
```

**Identifiants par défaut**: `admin / admin123`

---

## Endpoints d'authentification

| Méthode | Endpoint | Description |
|---------|----------|-------------|
| POST | `/auth/login` | Login (retourne access_token + refresh_token) |
| POST | `/auth/refresh` | Rafraîchir le token d'accès |
| POST | `/auth/logout` | Invalider les tokens |
| GET | `/auth/me` | Infos utilisateur courant |

---

## Endpoints principaux

| Méthode | Endpoint | Description |
|---------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/agents` | Liste des agents |
| GET | `/techniques` | Liste des techniques MITRE |
| GET | `/techniques/coverage` | Statistiques de couverture MITRE |
| GET | `/scenarios` | Liste des scénarios |
| POST | `/executions` | Lancer une exécution |
| GET | `/executions/:id` | Détails d'une exécution |
| GET | `/executions/:id/results` | Résultats d'une exécution |

---

## Analytics

| Méthode | Endpoint | Description |
|---------|----------|-------------|
| GET | `/analytics/period` | Statistiques par période |
| GET | `/analytics/comparison` | Comparaison entre périodes |
| GET | `/analytics/trend` | Tendance du score |
| GET | `/analytics/summary` | Résumé des exécutions |

---

## Planification (Schedules)

| Méthode | Endpoint | Description |
|---------|----------|-------------|
| GET | `/schedules` | Liste des planifications |
| POST | `/schedules` | Créer une planification |
| PUT | `/schedules/:id` | Modifier une planification |
| DELETE | `/schedules/:id` | Supprimer une planification |
| POST | `/schedules/:id/pause` | Mettre en pause |
| POST | `/schedules/:id/resume` | Reprendre |
| POST | `/schedules/:id/run` | Exécuter maintenant |

---

## Notifications

| Méthode | Endpoint | Description |
|---------|----------|-------------|
| GET | `/notifications` | Liste des notifications |
| GET | `/notifications/unread/count` | Nombre de non-lues |
| POST | `/notifications/:id/read` | Marquer comme lue |
| POST | `/notifications/read-all` | Tout marquer comme lu |
| GET | `/notifications/settings` | Paramètres de notification |
| GET | `/notifications/smtp` | Configuration SMTP (admin) |

---

## Administration

| Méthode | Endpoint | Description |
|---------|----------|-------------|
| GET | `/admin/users` | Liste des utilisateurs |
| POST | `/admin/users` | Créer un utilisateur |
| PUT | `/admin/users/:id` | Modifier un utilisateur |
| PUT | `/admin/users/:id/role` | Changer le rôle |
| DELETE | `/admin/users/:id` | Désactiver un utilisateur |
| POST | `/admin/users/:id/reactivate` | Réactiver un utilisateur |
| POST | `/admin/users/:id/reset-password` | Réinitialiser le mot de passe |

---

## Permissions

| Méthode | Endpoint | Description |
|---------|----------|-------------|
| GET | `/permissions/matrix` | Matrice des permissions par rôle |
| GET | `/permissions/me` | Mes permissions |

---

## WebSocket

| Path | Description |
|------|-------------|
| `/ws/agent` | Connexion agent |
| `/ws/dashboard` | Mises à jour temps réel |

---

## Documentation complète

Voir la [Référence API](reference.md) pour tous les endpoints, paramètres et exemples.
