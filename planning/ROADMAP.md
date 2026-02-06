# AutoStrike - Roadmap & Taches Restantes

> Mis a jour le 2026-02-06

**Documents lies :**
- [AutoStrike_Documentation.md](./AutoStrike_Documentation.md) - Vision strategique 3 ans, architecture, EBIOS RM
- [PRESENTATION.md](./PRESENTATION.md) - Slides pour presentation equipe
- [VISION_V2.md](./VISION_V2.md) - Vision autonome (Decision Engine, Blackbox, Agent propagation)

---

## Metriques Actuelles

| Metrique | Valeur |
|----------|--------|
| Issues fermees | 170+ |
| Tests | 780+ (200+ server + 513 dashboard + 67 agent) |
| Coverage | 95%+ domaine |
| Techniques MITRE | 48 (13 tactiques sur 14) |
| Lignes de code | ~18,000 |

---

## Vue d'ensemble

| Categorie | Status | Priorite |
|-----------|--------|----------|
| Authentification (Auth) | ✅ **Termine** | - |
| Security Hardening | ✅ **Termine** | - |
| Docker & Deploiement | ✅ **Termine** | - |
| Techniques MITRE (48) | ✅ **Termine** | - |
| Scheduling & Notifications | ✅ **Termine** | - |
| RBAC & Permissions | ✅ **Termine** | - |
| Profils adversaires APT | ❌ A faire | 🔴 Haute |
| Export PDF rapports | ❌ A faire | 🔴 Haute |
| ScenarioBuilder visuel | ❌ A faire | 🟡 Moyenne |
| Cleanup techniques | ❌ A faire | 🟡 Moyenne |
| Agent auto-deploy | ❌ A faire | 🟡 Moyenne |
| Features Caldera-like | ❌ A faire | 🟢 Stretch |
| **Dashboard V2** (Attack Graph, Approvals) | ❌ A faire | 🟡 Moyenne |
| **Vision Autonome** (Decision Engine, Blackbox) | ❌ A faire | 🔵 V2 |

---

## ✅ Phase 3 - Authentification - TERMINEE

### Backend Auth
| Issue | Titre | Status |
|-------|-------|--------|
| #51 | Middleware JWT | ✅ Fait |
| #52 | Handler `/api/v1/auth` (login, refresh, logout) | ✅ Fait |
| #135 | Service API auth (frontend) | ✅ Fait |

### Frontend Auth
| Issue | Titre | Status |
|-------|-------|--------|
| #142 | Page Login | ✅ Fait |
| #174 | Protected routes | ✅ Fait |
| #175 | Token storage | ✅ Fait |

---

## ✅ Phase 3 - Security Hardening - TERMINEE

| Issue | Titre | Status |
|-------|-------|--------|
| #209 | Audit securite Go | ✅ Fait (security headers, input validation) |
| #210 | Audit securite Rust | ✅ Fait (output truncation, UTF-8 safety) |
| #211 | Audit securite React | ✅ Fait (CSP, XSS protection) |
| #212 | Rate limiting | ✅ Fait (per-IP, login 5/min, refresh 10/min) |
| #213 | Audit logging | ✅ Fait (structured logging zap) |
| #214 | Review mTLS | ✅ Fait (script generate-certs.sh) |
| - | Token blacklist (logout) | ✅ Fait (in-memory avec auto-cleanup) |
| - | Security headers middleware | ✅ Fait (CSP, HSTS, X-Frame-Options) |

---

## ✅ Phase 4 - Docker Production - TERMINEE

| Issue | Titre | Status |
|-------|-------|--------|
| #206 | docker-compose.yml | ✅ Fait |
| #207 | docker-compose.dev.yml | ✅ Fait (hot reload, health checks) |
| #208 | Script generation certificats | ✅ Fait (mTLS CA + server + agent) |

---

## ✅ Phase 4 - Features Avancees - TERMINEE

### User Stories Completees
| Issue | Titre | Status |
|-------|-------|--------|
| #12 | Login UI | ✅ Fait |
| #218 | Comparer scores | ✅ Fait (analytics avec comparaison periodes) |
| #219 | Planifier executions | ✅ Fait (cron, daily, weekly, monthly) |
| #220 | Gestion utilisateurs | ✅ Fait (CRUD admin complet) |
| #222 | Import/Export scenarios | ✅ Fait (YAML/JSON) |
| #223 | Notifications email | ✅ Fait (SMTP + webhooks) |
| #224 | Permissions granulaires | ✅ Fait (28 permissions, 5 roles) |
| #54 | Middleware rate limiting | ✅ Fait |
| #61 | Systeme notifications | ✅ Fait |

### Techniques MITRE - 48 techniques, 13 tactiques
| Tactique | Count | Status |
|----------|-------|--------|
| Reconnaissance | 2 | ✅ Nouveau |
| Initial Access | 3 | ✅ Nouveau |
| Execution | 5 | ✅ +2 |
| Persistence | 4 | ✅ |
| Privilege Escalation | 4 | ✅ Nouveau |
| Defense Evasion | 6 | ✅ +3 |
| Credential Access | 4 | ✅ |
| Discovery | 9 | ✅ |
| Lateral Movement | 3 | ✅ |
| Collection | 4 | ✅ |
| Command and Control | 3 | ✅ Nouveau |
| Exfiltration | 3 | ✅ Nouveau |
| Impact | 3 | ✅ Nouveau |

### Frontend Complete
| Issue | Titre | Status |
|-------|-------|--------|
| #171 | SecurityScore component | ✅ Fait |
| #172 | CoverageReport | ✅ Fait |
| #176 | Theme sombre/clair | ✅ Fait |

---

## 🔴 Phase 5 - Priorite Haute (a faire)

### Profils Adversaires APT
| Tache | Description | Effort |
|-------|-------------|--------|
| Scenario APT29 (Cozy Bear) | 4 phases : recon, exec, persistence, defense evasion | 1h |
| Scenario Ransomware | 3 phases : discovery, collection, impact | 1h |
| Scenario Basic Reconnaissance | Discovery complete | 30min |
| Scenario Insider Threat | Credential access + collection + exfiltration | 1h |
| Scenario Full Kill Chain | Toutes les tactiques enchainées | 2h |

**Spec deja ecrite dans AutoStrike_Documentation.md (l.2885-2942)**

### Export PDF Rapports
| Tache | Description | Effort |
|-------|-------------|--------|
| #50 | Backend handler `/api/v1/reports` | 4h |
| #79 | Generateur PDF (go-pdf ou wkhtmltopdf) | 6h |
| #170 | Page Reports dashboard | 4h |
| #173 | Bouton ExportPDF component | 2h |

---

## 🟡 Phase 5 - Priorite Moyenne (a faire)

### Cleanup Automatique (RG-005)
| Tache | Description | Effort |
|-------|-------------|--------|
| Ajouter cleanup aux techniques YAML | Commandes de nettoyage post-execution | 2h |
| T1053.005 cleanup | `schtasks /Delete /TN autostrike-test /F` | 15min |
| T1547.001 cleanup | `reg delete HKCU\...\Run /v autostrike-test /f` | 15min |
| Verifier execution cleanup dans l'agent | Le champ existe deja dans le protocole | 1h |

### Agent Auto-Deploy (BE-56)
| Tache | Description | Effort |
|-------|-------------|--------|
| Endpoint `GET /deploy/agent.sh` | Script bash avec URL serveur injectee | 2h |
| Endpoint `GET /deploy/agent.ps1` | Script PowerShell equivalent | 2h |
| Page dashboard "Deploy Agent" | UI avec commande copier/coller | 2h |

### ScenarioBuilder Visuel
| Tache | Description | Effort |
|-------|-------------|--------|
| #151 | Editeur drag & drop de phases/techniques | 8h |
| Selection par tactic/platform | Filtres dans le builder | 2h |
| Preview YAML | Voir le YAML genere avant sauvegarde | 1h |

### LiveLogs
| Tache | Description | Effort |
|-------|-------------|--------|
| #163 | Page logs temps reel via WebSocket | 4h |

---

## 🟢 Phase 6 - Stretch Goals (Caldera-like features)

> Ces features rapprocheraient AutoStrike de MITRE Caldera. Non planifiees initialement mais differenciantes pour l'EIP.

### Planners Intelligents
| Tache | Description | Effort |
|-------|-------------|--------|
| Planner "atomic" | Execute chaque technique independamment (actuel) | ✅ Fait |
| Planner "sequential" | Execute en sequence, arrete si echec | 4h |
| Planner "conditional" | Decision dynamique selon les resultats precedents | 12h |
| Planner "buckets" | Regroupe par tactic, randomise l'ordre | 6h |

**Impact :** Necessite refactor de `AttackOrchestrator` pour supporter differentes strategies d'enchainement.

### Facts / Data Exchange
| Tache | Description | Effort |
|-------|-------------|--------|
| Entite `Fact` (key/value) | Stocker les donnees decouvertes pendant l'execution | 4h |
| Collecte facts depuis output | Parser les resultats des techniques pour extraire des facts | 8h |
| Variables dans les commandes | Remplacer `#{username}` par le fact decouvert | 4h |
| Facts store persistant | Stocker les facts entre executions | 2h |

**Exemple :** T1087 (Account Discovery) decouvre les users → T1552.001 cherche les credentials de ces users specifiques.

### Obfuscation des Commandes
| Tache | Description | Effort |
|-------|-------------|--------|
| Base64 encoding | Encoder les commandes PowerShell en Base64 | 2h |
| String concatenation | Decouper les commandes en morceaux | 2h |
| Variable substitution | Utiliser des variables intermediaires | 2h |
| Plugin framework | Interface extensible pour les obfuscators | 4h |

### Recommandations de Remediation
| Tache | Description | Effort |
|-------|-------------|--------|
| Champ `mitigations` dans Technique | Liste des mitigations ATT&CK par technique | 2h |
| Endpoint `/api/v1/executions/:id/recommendations` | Generer les recommandations post-execution | 4h |
| Page Recommendations dashboard | Afficher les remediations par priorite | 4h |
| Mapping MITRE Mitigations | Associer M1036, M1049, etc. aux techniques | 3h |

### Multiple Agent Types
| Tache | Description | Effort |
|-------|-------------|--------|
| Agent Python (leger) | Pour les environnements ou Rust ne compile pas | 8h |
| Agent reverse shell | Connexion inverse pour les reseaux restreints | 6h |

---

## 🟡 Dashboard V2 - Ameliorations UI (a faire)

> Features inspirees de Pentera et des outils BAS modernes. Voir [VISION_V2.md](./VISION_V2.md) pour le detail.

### Attack Graph / Attack Map
| Tache | Description | Effort |
|-------|-------------|--------|
| Composant AttackGraph | Graphe de noeuds (React Flow / D3.js) montrant la propagation | 8h |
| Animations temps reel | Fleches animees, couleurs par severite | 4h |
| Clic sur lien → technique | Afficher la technique MITRE utilisee pour chaque passage | 2h |

### Details Panel enrichi
| Tache | Description | Effort |
|-------|-------------|--------|
| Panel Insight | Explication du risque contextuel (slide-over) | 3h |
| Panel Remediation | Recommandations concretes par technique | 4h |
| Simulate Fix | Simuler un patch → la case MITRE passe du rouge au vert | 4h |

### Approvals System
| Tache | Description | Effort |
|-------|-------------|--------|
| Workflow approbation | Notification "Pending Approval" avant actions critiques | 6h |
| Approbation par host | Valider les cibles individuellement | 3h |
| Mode Full Auto vs Supervised | Toggle dans les settings | 2h |
| Kill Switch global | Bouton "Panic" pour deconnecter tous les agents | 2h |

### Dashboard Overview ameliore
| Tache | Description | Effort |
|-------|-------------|--------|
| Vulnerability cards | Cards par severite (Critical/High/Medium/Low) | 3h |
| Device discovery | Breakdown par OS, services | 3h |
| Host grid | Grille des hotes avec IPs et statuts | 3h |

### Achievements / Gamification
| Tache | Description | Effort |
|-------|-------------|--------|
| Scores 1-10 par technique | Scoring granulaire au-dela du Security Score global | 4h |
| Badges | "Domain Dominator", "Ghost in the Shell", etc. | 3h |
| Leaderboard | Comparer la resilience de plusieurs reseaux | 3h |

### Timeline chronologique
| Tache | Description | Effort |
|-------|-------------|--------|
| Timeline verticale | Flux chronologique de toutes les actions | 4h |
| Integration avec Export PDF | Ajouter la timeline dans les rapports | 2h |

---

## 🔵 Phase V2 - Vision Autonome (2026-2028)

> Faire evoluer AutoStrike d'un outil BAS classique vers un pentesteur autonome.
> Detail complet dans [VISION_V2.md](./VISION_V2.md).

### Decision Engine ("The Brain")
| Tache | Description | Effort |
|-------|-------------|--------|
| Target Knowledge DB | Entite + repository pour stocker l'etat de la connaissance cible | 8h |
| Autonomous Planner | Moteur de regles (if/else) qui decide la prochaine action | 12h |
| Feedback Loop | Le serveur analyse les outputs pour mettre a jour la Knowledge DB | 6h |
| Mode Blackbox API | `POST /api/v1/executions` avec `mode: "blackbox"` | 4h |

### Agent Avance (Rust)
| Tache | Description | Effort |
|-------|-------------|--------|
| Network scan local | Decouverte des machines voisines | 6h |
| Self-propagation | Copie du binaire via SMB/SSH si credentials trouves | 10h |
| Safety Rails | Module bloquant les commandes destructrices | 4h |
| Kill Switch | Deconnexion globale sur ordre du serveur | 2h |

### Live Terminal IA
| Tache | Description | Effort |
|-------|-------------|--------|
| Log de decision | Afficher les decisions du planner en temps reel (WebSocket) | 4h |
| Format lisible | Icones par type (RECON, DECIDE, ALERT, SUCCESS, PIVOT) | 2h |

### Objectif CTF
| Tache | Description | Effort |
|-------|-------------|--------|
| Prototype Blackbox | Resoudre une box HackTheBox Easy en autonomie | 20h+ |

---

## Issues a Skipper

| Issue | Raison |
|-------|--------|
| #136-141 | Zustand stores → On utilise TanStack Query |
| #158-159 | D3.js → CSS Grid suffit pour la matrice |
| #40 | Migrations BDD versionnees → SQLite suffit pour le MVP |

---

## Repartition par Profil

### Architecture & Agent Lead
- Decision Engine / Autonomous Planner (V2)
- Agent avance Rust (self-propagation, pivoting, safety rails)
- Architecture globale et coordination

### Frontend Lead (React)
- Attack Graph / Attack Map (V2)
- Live Terminal IA (V2)
- Blast Radius View (V2)
- Approvals System (V2)
- Details Panel enrichi (Insight, Remediation, Simulate Fix)
- Dashboard Overview ameliore
- Achievements / Gamification
- Page Reports (#170)
- ScenarioBuilder (#151)
- LiveLogs (#163)

### Backend Dev 1 (Go)
- Target Knowledge DB (V2)
- Facts system (Caldera-like)
- Export PDF (#50, #79)
- Endpoint auto-deploy (BE-56)
- Planners intelligents

### Backend Dev 2 (Go)
- Feedback Loop / Planner logic (V2)
- Endpoint recommendations
- Integration tests
- Obfuscation plugins

### Security / MITRE (tous)
- Profils APT (YAML scenarios)
- Cleanup techniques (YAML)
- Mapping mitigations ATT&CK

### N'importe qui
- Tests manuels
- Documentation
- Ajouter techniques YAML

---

## Timeline (2026-2028)

> Projet EIP Promotion 2028. Rendu final : fin d'annee scolaire 2028.

```
S1 2026 (fev-juil) : Quick Wins + Fondations V2
├── Profils APT (scenarios YAML)
├── Export PDF (backend + frontend)
├── Cleanup techniques + auto-deploy scripts
├── ScenarioBuilder visuel + LiveLogs
├── Target Knowledge DB (fondation Decision Engine)
└── Moteur de regles basique (if/else)

S2 2026 (sept-dec) : Dashboard V2 + Decision Engine
├── Attack Graph / Attack Map (React Flow / D3.js)
├── Live Terminal IA + Details Panel enrichi
├── Approvals System + Kill Switch
├── Feedback Loop (serveur reagit aux outputs)
├── Mode Blackbox API
└── Planners intelligents (sequential, conditional, buckets)

S1 2027 (jan-juil) : Agent Avance + Caldera-like
├── Network scan local + self-propagation (SMB/SSH)
├── Safety Rails (garde-fous agent)
├── Facts / Data Exchange system
├── Obfuscation des commandes
├── Recommandations de remediation
└── Dashboard Overview ameliore (vuln cards, host grid)

S2 2027 (sept-dec) : Polish + Features avancees
├── Blast Radius View (environnements segmentes)
├── Achievements / Gamification
├── Timeline chronologique + integration PDF
└── Multiple agent types

S1 2028 (jan-juin) : Autonomie Avancee + Soutenance
├── Decision Engine v2 (patterns avances)
├── CTF autonome (HackTheBox Easy)
├── Mode Blackbox complet (Scan → Root → Rapport)
├── Tests E2E sur environnements lab
├── Documentation finale
└── Demo soutenance EIP
```

---

## Comparaison avec MITRE Caldera

> AutoStrike couvre ~80% des features core de Caldera avec des avantages uniques.

### Ce qu'AutoStrike a que Caldera n'a pas
| Feature | Detail |
|---------|--------|
| Scheduling automatique | Cron, daily, weekly, monthly |
| Security Score 0-100 | Formule blocked/detected/success |
| Analytics & tendances | Comparaison periodes, trends |
| RBAC granulaire | 5 roles, 28 permissions |
| Notifications email/webhook | Alertes automatiques |
| Dark mode | Interface moderne |

### Ce qui manque vs Caldera
| Feature | Status | Effort |
|---------|--------|--------|
| Planners intelligents | ❌ A faire | 22h total |
| Facts / data exchange | ❌ A faire | 18h total |
| Obfuscation plugins | ❌ A faire | 10h total |
| Multiple agent types | ❌ A faire | 14h total |
| Agent auto-deploy | ❌ A faire | 6h total |
| Cleanup automatique | ❌ A faire | 3h total |
| Recommandations remediation | ❌ A faire | 13h total |
| Profils APT predefinis | ❌ A faire | 5h total |

---

## Commandes Utiles

```bash
# Voir les issues ouvertes
gh issue list --state open

# Assigner une issue
gh issue edit <numero> --add-assignee <username>

# Fermer une issue
gh issue close <numero> --comment "Implemente dans <commit>"

# Creer une branche pour une issue
git checkout -b feat/issue-<numero>-description

# Lier un commit a une issue
git commit -m "feat: description (#<numero>)"
```

---

*Derniere mise a jour: 2026-02-06*
