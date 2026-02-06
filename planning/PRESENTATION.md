# Guide de Présentation AutoStrike

**Documents liés :**
- [AutoStrike_Documentation.md](./AutoStrike_Documentation.md) - Vision stratégique 3 ans, EBIOS RM, compétences EIP
- [ROADMAP.md](./ROADMAP.md) - Issues GitHub, tâches opérationnelles

---

## 1. Pitch d'intro (30 sec)

> "AutoStrike c'est une plateforme de Breach and Attack Simulation. En gros, on simule des attaques basées sur le framework MITRE ATT&CK pour tester si les défenses d'un réseau détectent bien les menaces. C'est ce qu'utilisent les équipes SOC et les pentesters pour valider leurs détections."

---

## 2. Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    Server (Go) - Port 8443                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │  Dashboard  │  │  REST API   │  │  WebSocket  │              │
│  │  (Static)   │  │  /api/v1/*  │  │  /ws/*      │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
└─────────────────────────────────────────────────────────────────┘
                              ▲
                              │ WebSocket
                    ┌─────────┴─────────┐
                    │      Agent        │
                    │     (Rust)        │
                    └───────────────────┘
```

**Explication simple :**
> "Un seul serveur Go sur le port 8443 qui fait tout : il sert le dashboard React, expose l'API REST, et gère les WebSockets pour communiquer avec les agents. Les agents sont en Rust, ils se connectent au serveur et exécutent les commandes qu'on leur envoie."

---

## 3. Pourquoi ces technos ?

| Techno | Argument |
|--------|----------|
| **Go (backend)** | Performant, compile en un seul binaire, gestion native de la concurrence avec les goroutines, parfait pour gérer plein d'agents en parallèle |
| **Rust (agent)** | Sécurité mémoire sans garbage collector, cross-compile facilement pour Windows/Linux, difficile à reverse-engineer pour un outil de sécu |
| **React + TS** | Écosystème mature, TypeScript pour la maintenabilité, facile à recruter des devs |
| **WebSocket** | Communication bidirectionnelle en temps réel, l'agent peut recevoir des commandes et renvoyer les résultats instantanément |
| **SQLite** | Zéro config pour le MVP, on migre vers PostgreSQL en prod si besoin |

---

## 4. Architecture Hexagonale

```
server/internal/
├── domain/          # Logique métier pure (pas de dépendances externes)
│   ├── entity/      # Structures de données (Agent, Technique, Execution)
│   └── service/     # Logique business (TechniqueService, ExecutionService)
├── application/     # Orchestration des cas d'usage
└── infrastructure/  # Adapters vers l'extérieur
    ├── http/        # Handlers API REST
    ├── persistence/ # SQLite repositories
    └── websocket/   # Communication agents
```

**Explication simple :**
> "C'est un pattern d'architecture où les dépendances vont toujours vers l'intérieur. Le domain ne connaît pas la base de données ni HTTP. Ça permet de changer facilement SQLite par PostgreSQL ou Gin par un autre framework sans toucher à la logique métier."

---

## 5. Flux d'une exécution

```
1. User clique "Run" sur un scénario dans le Dashboard
2. Dashboard appelle POST /api/v1/executions
3. Server crée l'exécution en DB et notifie les agents via WebSocket
4. Agent reçoit les tâches, exécute les commandes
5. Agent renvoie les résultats via WebSocket
6. Server stocke les résultats et notifie le Dashboard
7. Dashboard affiche les résultats en temps réel
```

---

## 6. MITRE ATT&CK

> "MITRE ATT&CK c'est un framework qui catalogue toutes les techniques d'attaque connues. C'est organisé en tactiques (ce que l'attaquant veut faire : Discovery, Execution, Persistence...) et techniques (comment il le fait : T1082 System Information Discovery, T1059 Command Execution...). Nous on implémente ces techniques pour les simuler."

**Exemple concret :**
> "La technique T1082 c'est 'System Information Discovery'. Sur Windows ça exécute `systeminfo`, sur Linux `uname -a`. L'agent exécute ça et renvoie le résultat."

---

## 7. Démo live (script)

```bash
# Terminal 1 - Démarre le serveur
make run
# Montre les logs : "Server started on :8443"

# Navigateur - Ouvre le dashboard
# https://localhost:8443
# Montre : Agents (vide), Techniques, Scenarios, Matrix

# Terminal 2 - Lance un agent
make agent
# Montre dans les logs : "Agent registered: ..."

# Dashboard - Refresh
# L'agent apparaît dans la liste avec son hostname, OS, etc.

# Dashboard - Va dans Scenarios
# Montre un scénario, clique Run, sélectionne l'agent
# Montre l'exécution en temps réel
```

---

## 8. Structure des fichiers clés

| Fichier | Rôle |
|---------|------|
| `server/cmd/autostrike/main.go` | Point d'entrée, wire les dépendances |
| `server/internal/domain/entity/` | Toutes les structures (Agent, Technique, Execution) |
| `server/internal/domain/service/` | Logique métier |
| `server/internal/infrastructure/http/handlers/` | Endpoints API |
| `server/internal/infrastructure/websocket/` | Hub + Client WebSocket |
| `dashboard/src/pages/` | Pages React (Agents, Techniques, Executions, Matrix) |
| `dashboard/src/lib/api.ts` | Client API |
| `agent/src/main.rs` | Point d'entrée agent |
| `agent/src/executor.rs` | Exécution des commandes |

---

## 9. Pourquoi ne pas changer (arguments clés)

### Langages - Standards de l'industrie

| Choix | Argument imparable |
|-------|-------------------|
| **Go** | "Kubernetes, Docker, Terraform, Caldera sont en Go. C'est LE standard pour les outils d'infra et sécu." |
| **Rust** | "Les agents Python se font bloquer par Windows Defender. Rust compile en natif, passe sous les radars. C'est un choix de sécurité, pas de confort." |
| **React** | "Framework le plus populaire, tout le monde peut contribuer. TypeScript évite les bugs runtime." |

### Reponses aux objections

**"On aurait pu faire plus simple"**
> "Simple = Python partout. Mais un agent Python se fait detecter en 2 secondes par un AV. On a fait les choix d'un outil de secu **professionnel**."

**"Je connais pas Go/Rust"**
> "Pas besoin de tout maitriser. Tu peux contribuer sur React (dashboard), YAML (techniques), ou Markdown (docs). L'archi decouplée permet de bosser isolement."

**"Pourquoi pas repartir de zero ?"**
> "On a 780+ tests qui passent, 48 techniques MITRE, une auth complete avec RBAC, du scheduling, des notifications, Docker ready. Repartir = perdre des mois de travail."

**"Pourquoi pas utiliser Caldera directement ?"**
> "Caldera est en Python, plus lourd. Notre agent Rust est plus discret. Et on a des features que Caldera n'a pas : scoring, analytics, scheduling, RBAC granulaire. On peut meme aller plus loin avec les planners et facts."

### L'avantage de l'architecture hexagonale

> "Si demain on veut changer SQLite pour PostgreSQL, on touche **1 seul fichier**. Si on veut remplacer le framework HTTP, on touche **1 dossier**. Le code métier ne bouge pas."

### Phrase clé

> **"J'ai posé les fondations solides. Maintenant on construit dessus ensemble. Chacun peut bosser sur sa partie sans casser le reste."**

---

## 10. Tâches accessibles pour chacun

| Profil | Tâche | Techno |
|--------|-------|--------|
| Connaît React | Page Login, améliorer UI | React/TS |
| Veut apprendre | Ajouter techniques MITRE | YAML (copier/coller) |
| DevOps | Docker prod, monitoring | Docker/K8s |
| N'importe qui | Tests, documentation | Markdown |

---

## 11. Questions probables + réponses

### Q: "Pourquoi pas Python pour l'agent ?"
> "Python c'est facile à détecter et à bloquer par les antivirus. Rust compile en binaire natif, c'est plus discret et plus performant pour un outil de sécu."

### Q: "C'est quoi le Safe Mode ?"
> "Quand safe_mode est activé, l'agent n'exécute pas vraiment les commandes, il simule juste. C'est pour tester sans risquer de casser quelque chose."

### Q: "Comment on ajoute une nouvelle technique ?"
> "Tu crées un fichier YAML dans `server/configs/techniques/` avec l'ID MITRE, le nom, la tactique, et les commandes par OS. Au démarrage le serveur les importe."

### Q: "C'est sécurisé ?"
> "En prod on active JWT pour l'API et mTLS pour les agents. En dev c'est désactivé pour simplifier. Y'a aussi un AGENT_SECRET pour authentifier les agents."

### Q: "Comment tu as fait les tests ?"
> "J'ai utilisé Vitest pour le front avec des mocks des APIs. Pour le back c'est les tests Go classiques avec des mocks des repositories."

### Q: "Pourquoi un seul port ?"
> "Simplifie le déploiement et la config firewall. Le serveur Go sert tout : les fichiers statiques du dashboard, l'API REST, et les WebSockets."

---

## 12. Points forts a mentionner

- **780+ tests** qui passent (200+ server + 513 dashboard + 67 agent)
- **48 techniques MITRE** couvrant **13 tactiques sur 14** (de Reconnaissance a Impact)
- **Documentation complete** (MkDocs deploye sur GitHub Pages)
- **CI/CD** deja en place (GitHub Actions : build Go/Rust/Node, tests, Docker, SonarCloud)
- **Code coverage** 95%+ sur le domaine
- **Authentification complete** : JWT, 5 roles, 28 permissions, token blacklist
- **Scheduling** : executions planifiees (cron, daily, weekly, monthly)
- **Notifications** : email SMTP + webhooks
- **Security hardening** : rate limiting, security headers, CSP, HSTS
- **Docker ready** : docker-compose prod + dev + 3 Dockerfiles multi-stage
- Architecture **modulaire** et **decouplee** (hexagonale)
- Pret pour une **demo fonctionnelle**
- Couvre **~80% des features de MITRE Caldera** avec des avantages uniques (scoring, analytics, scheduling)

---

## 13. Ce qui reste a faire

### Priorite haute
- **Profils adversaires APT** : scenarios predefinis APT29, Ransomware, Insider Threat (specs deja ecrites)
- **Export PDF rapports** : handler backend + generateur PDF + page Reports

### Priorite moyenne
- **ScenarioBuilder visuel** : editeur drag & drop de phases/techniques
- **Cleanup automatique** : commandes de nettoyage post-execution (champ deja dans le protocole)
- **Agent auto-deploy** : endpoints pour generer les scripts d'installation
- **LiveLogs** : page logs temps reel via WebSocket

### Stretch Goals (features Caldera-like)
- **Planners intelligents** : enchainement conditionnel (si echec → skip, si succes → continuer)
- **Facts system** : passer des donnees entre techniques (ex: users decouverts → cibles brute force)
- **Obfuscation** : encoder les commandes (Base64, concatenation)
- **Recommandations remediation** : mapping mitigations ATT&CK automatique post-execution

---

## 14. Repartition du code

| Composant | Lignes | Notes |
|-----------|--------|-------|
| Server (Go) | ~14,000 | Code + 200+ tests |
| Dashboard (React) | ~16,400 | Code + 513 tests |
| Agent (Rust) | ~1,700 | 67 tests |
| Techniques YAML | ~1,200 | 48 techniques, 13 tactiques |
| Docs (MkDocs) | ~2,300 | |
| CI/CD, Docker, configs | ~2,000 | GitHub Actions, Dockerfiles, Makefile |

**Code source : ~18k lignes** (server + agent + dashboard, hors tests, docs et fichiers auto-generes)

## 15. Comparaison avec MITRE Caldera

> "AutoStrike c'est comme un Caldera custom. On couvre 80% de ses features core, mais on a des trucs qu'il n'a pas."

### Ce qu'on a que Caldera n'a pas
| Feature | Detail |
|---------|--------|
| **Scheduling** | Cron, daily, weekly, monthly (Caldera n'a pas ca nativement) |
| **Security Score** | Score 0-100 avec formule blocked/detected/success |
| **Analytics** | Comparaison entre periodes, tendances, graphiques |
| **RBAC** | 5 roles, 28 permissions granulaires |
| **Notifications** | Email SMTP + webhooks automatiques |

### Ce qui nous manque vs Caldera
| Feature | Description | Effort |
|---------|-------------|--------|
| **Planners intelligents** | Enchainement conditionnel des techniques | 22h |
| **Facts system** | Passer des donnees entre techniques | 18h |
| **Profils APT** | Scenarios predefinis (specs deja ecrites) | 5h |
| **Recommandations** | Mitigations ATT&CK post-execution | 13h |

### Reponse a "Pourquoi pas juste utiliser Caldera ?"
> "Caldera c'est Python, c'est plus lourd a deployer, et il n'a pas de scoring, pas de scheduling, pas d'analytics. Notre agent Rust est plus leger et plus discret qu'un agent Python. Et on a une architecture hexagonale propre qui permet d'evoluer facilement."

---

## 16. Vision V2 : Autonomous Threat Emulation

> Detail complet dans [VISION_V2.md](./VISION_V2.md)

### Le pivot strategique
> "On passe d'un test statique (l'utilisateur choisit tout) a un test dynamique (l'outil decide seul). On ne change pas le but (tester la defense), on change la qualite du test."

### Ce que ca change concrètement

| Avant (V1) | Apres (V2) |
|-------------|------------|
| L'utilisateur cree un scenario | L'utilisateur donne une IP cible |
| L'agent execute une liste fixe | L'agent s'adapte aux resultats |
| Resultat = rapport statique | Resultat = graphe de propagation live |
| Test ponctuel | Simulation d'une vraie APT |

### Les 3 piliers V2

1. **Decision Engine** : Le serveur Go maintient une "memoire" de la cible et decide de la prochaine action (boucle OODA)
2. **Agent Avance** : Self-propagation, pivoting reseau, safety rails (garde-fous)
3. **Dashboard V2** : Attack Graph, Live Terminal IA, Approvals System, Blast Radius

### Objectif demo (soutenance 2028)
> "On donne une IP a AutoStrike, et le jury voit en temps reel : le graphe de propagation s'etendre, les decisions de l'IA defiler, et la matrice MITRE se remplir. Scan → Exploit → Root → Rapport."

### Reponse a "C'est pas du malware ?"
> "C'est 100% dans l'esprit BAS. Une vraie menace (APT, ransomware) ne suit pas une liste, elle s'adapte. Si on simule des menaces realistes, notre outil doit s'adapter aussi. La difference avec un vrai virus : nos Safety Rails empechent toute action destructrice. L'agent prouve qu'il POURRAIT le faire sans le faire."
