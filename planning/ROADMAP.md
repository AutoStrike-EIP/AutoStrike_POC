# AutoStrike - Roadmap & Tâches Restantes

> Généré le 2026-02-03 | 49 issues ouvertes sur GitHub

**Documents liés :**
- [AutoStrike_Documentation.md](./AutoStrike_Documentation.md) - Vision stratégique 3 ans, architecture, EBIOS RM
- [PRESENTATION.md](./PRESENTATION.md) - Slides pour présentation équipe

---

## Vue d'ensemble

| Catégorie | Issues | Priorité |
|-----------|--------|----------|
| Authentification (Auth) | 10 | 🔴 Haute |
| User Stories Avancées | 8 | 🟡 Moyenne |
| Security Audit | 6 | 🔴 Haute |
| Docker & Déploiement | 3 | 🟡 Moyenne |
| Frontend Features | 10 | 🟢 Basse |
| Backend Features | 8 | 🟡 Moyenne |
| Documentation | 2 | 🟢 Basse |

---

## 🔴 Phase 3 - Authentification (Priorité Haute)

### Backend Auth
| Issue | Titre | Effort | Assigné |
|-------|-------|--------|---------|
| #51 | Middleware JWT | ✅ Fait | - |
| #52 | Handler `/api/v1/auth` (login, refresh, logout) | 4h | |
| #135 | Service API auth (frontend) | 2h | |

### Frontend Auth
| Issue | Titre | Effort | Assigné |
|-------|-------|--------|---------|
| #142 | Page Login | 4h | |
| #174 | Protected routes | 2h | |
| #175 | Token storage (localStorage/cookies) | 2h | |
| #136-141 | Zustand stores (optionnel - on utilise TanStack Query) | - | Skipper |

**Dépendances:** #52 → #135 → #142 → #174 → #175

---

## 🔴 Phase 3 - Security Audit (Priorité Haute)

| Issue | Titre | Description | Effort |
|-------|-------|-------------|--------|
| #209 | Audit sécurité Go | Revue OWASP, injection SQL, auth bypass | 4h |
| #210 | Audit sécurité Rust | Memory safety, command injection | 4h |
| #211 | Audit sécurité React | XSS, CSRF, token exposure | 4h |
| #212 | Rate limiting | Implémenter middleware rate limit | 2h |
| #213 | Audit logging | Logger toutes les actions sensibles | 4h |
| #214 | Review mTLS | Vérifier config certificats prod | 2h |

**Livrable:** Rapport de sécurité + corrections

---

## 🟡 Phase 4 - Docker Production

| Issue | Titre | Description | Effort |
|-------|-------|-------------|--------|
| #206 | docker-compose.yml | Stack complète (server + dashboard + db) | 2h |
| #207 | docker-compose.dev.yml | Config développement avec hot reload | 1h |
| #208 | Script génération certificats | Auto-génération certs pour Docker | 1h |

**Commandes cibles:**
```bash
# Production
docker compose up -d

# Développement
docker compose -f docker-compose.dev.yml up
```

---

## 🟡 Phase 4 - Features Avancées

### User Stories
| Issue | Titre | Description | Effort |
|-------|-------|-------------|--------|
| #12 | Login UI | Page connexion avec JWT | 4h |
| #13 | Créer scénarios custom | ScenarioBuilder drag & drop | 8h |
| #14 | Export PDF rapports | Génération PDF avec résultats | 6h |
| #16 | Profils adversaires APT | APT28, APT29 predefined scenarios | 4h |
| #218 | Comparer scores | Graphique évolution temporelle | 4h |
| #219 | Planifier exécutions | Cron-like scheduling | 6h |
| #220 | Gestion utilisateurs | CRUD users, roles, permissions | 8h |
| #222 | Import/Export scénarios | YAML/JSON import/export | 4h |
| #223 | Notifications email | Alertes par email | 6h |
| #224 | Permissions granulaires | RBAC complet | 8h |

### Backend
| Issue | Titre | Effort |
|-------|-------|--------|
| #40 | Migrations BDD versionnées | 4h |
| #50 | Handler `/api/v1/reports` | 4h |
| #54 | Middleware rate limiting | 2h |
| #61 | Système notifications | 4h |
| #78 | Endpoint script déploiement | 2h |
| #79 | Export PDF backend | 4h |
| #80-81 | Profils APT28/APT29 | 4h |

---

## 🟢 Phase 5 - Frontend Avancé

| Issue | Titre | Description | Effort |
|-------|-------|-------------|--------|
| #151 | ScenarioBuilder | Éditeur visuel de scénarios | 8h |
| #158-159 | D3.js | ❌ Non nécessaire (CSS Grid suffit) | Skipper |
| #163 | LiveLogs | Logs temps réel WebSocket | 4h |
| #170 | Page Reports | Dashboard rapports | 6h |
| #171 | SecurityScore component | Widget score réutilisable | 2h |
| #172 | CoverageReport | Rapport couverture MITRE | 4h |
| #173 | ExportPDF component | Bouton export PDF | 4h |
| #176 | Thème sombre/clair | Toggle dark mode | 4h |

---

## 🟢 Documentation

| Issue | Titre | Effort |
|-------|-------|--------|
| #201 | Guide déploiement agent | 2h |
| #203 | Changelog | 1h |

---

## Répartition par Profil

### 👨‍💻 Développeur Backend (Go)
- #52 Handler auth
- #50 Handler reports
- #54 Rate limiting
- #61 Notifications
- #79 Export PDF
- #40 Migrations

### 👩‍💻 Développeur Frontend (React)
- #142 Page Login
- #151 ScenarioBuilder
- #170-173 Reports & Export
- #176 Dark mode
- #163 LiveLogs

### 🔒 Security Engineer
- #209-214 Audits sécurité

### 🐳 DevOps
- #206-208 Docker compose
- #201 Guide déploiement

### 📝 N'importe qui
- #203 Changelog
- Ajouter techniques YAML
- Tests manuels

---

## Timeline Suggérée

```
Semaine 1-2: Phase 3 Auth
├── Backend auth handler (#52)
├── Frontend login (#142, #174, #175)
└── Service API auth (#135)

Semaine 3: Security
├── Audits Go/Rust/React (#209-211)
├── Rate limiting (#54, #212)
└── Audit logging (#213)

Semaine 4: Docker & Deploy
├── docker-compose.yml (#206-208)
└── Guide déploiement (#201)

Semaine 5+: Features avancées
├── ScenarioBuilder (#151)
├── Export PDF (#14, #79, #173)
├── Reports (#50, #170-172)
└── Scheduling (#219)
```

---

## Issues à Skipper

Ces issues ne sont plus pertinentes (architecture changée):

| Issue | Raison |
|-------|--------|
| #136-141 | Zustand stores → On utilise TanStack Query |
| #158-159 | D3.js → CSS Grid suffit pour la matrice |

---

## Commandes Utiles

```bash
# Voir les issues ouvertes
gh issue list --state open

# Assigner une issue
gh issue edit <numero> --add-assignee <username>

# Fermer une issue
gh issue close <numero> --comment "✅ Implémenté dans <commit>"

# Créer une branche pour une issue
git checkout -b feat/issue-<numero>-description

# Lier un commit à une issue
git commit -m "feat: description (#<numero>)"
```

---

## Métriques Actuelles

| Métrique | Valeur |
|----------|--------|
| Issues fermées | 170 |
| Issues ouvertes | 49 |
| Tests | 447 (tous passent) |
| Coverage | 97%+ domaine |
| Techniques MITRE | 15 |

---

*Dernière mise à jour: 2026-02-03*
