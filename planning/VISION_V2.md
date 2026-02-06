# AutoStrike V2 - Autonomous Threat Emulation

> Vision strategique pour faire evoluer AutoStrike d'un outil BAS classique vers un **pentesteur artificiel autonome**.

**Documents lies :**
- [ROADMAP.md](./ROADMAP.md) - Taches operationnelles et phases
- [PRESENTATION.md](./PRESENTATION.md) - Guide de presentation equipe
- [AutoStrike_Documentation.md](./AutoStrike_Documentation.md) - Documentation technique complete

---

## 1. Le Constat

AutoStrike V1 est un outil BAS **reactif** : l'utilisateur cree un scenario, selectionne des techniques, lance l'execution. L'agent execute betement une liste de taches.

**Le probleme :** Une vraie menace (APT, ransomware) ne suit pas une liste. Elle s'adapte, pivote, exploite les failles de maniere dynamique.

**L'objectif V2 :** Passer d'un **test statique** a un **test dynamique**. L'outil decouvre, analyse et decide seul de sa strategie d'attaque.

> "On ne change pas le but (tester la defense), on change la qualite du test."

---

## 2. Inspirations

| Reference | Ce qu'on en retient |
|-----------|-------------------|
| **Pentera** | Validation continue autonome, attack paths automatiques |
| **Devin (Cognition)** | Agent IA qui planifie, execute, analyse les erreurs et s'adapte |
| **MITRE Caldera** | Architecture agent/server, planners, facts system |
| **DeepWiki** | Comprehension semantique d'un environnement complexe |

---

## 3. Architecture V2 : Le "Cerveau"

### 3.1 Boucle de Decision (OODA Loop)

```
┌──────────────────────────────────────────────────────┐
│                   OODA LOOP                           │
│                                                       │
│  ┌──────────┐   ┌──────────┐   ┌──────────┐         │
│  │ OBSERVE  │──>│  ORIENT  │──>│  DECIDE  │──┐      │
│  │ Agent    │   │ Knowledge│   │ Planner  │  │      │
│  │ remonte  │   │ DB mise  │   │ choisit  │  │      │
│  │ les data │   │ a jour   │   │ l'action │  │      │
│  └──────────┘   └──────────┘   └──────────┘  │      │
│       ^                                       v      │
│  ┌──────────┐                          ┌──────────┐  │
│  │ RESULTS  │<─────────────────────────│   ACT    │  │
│  │ Output,  │                          │ Commande │  │
│  │ exit code│                          │ envoyee  │  │
│  └──────────┘                          └──────────┘  │
└──────────────────────────────────────────────────────┘
```

### 3.2 Target Knowledge DB (La memoire)

Au lieu de stocker juste des logs, le serveur maintient un **etat de la connaissance** sur la cible.

```go
// Concept - server/internal/domain/entity/knowledge.go
type TargetKnowledge struct {
    Hosts       []HostInfo       // Machines decouvertes
    OpenPorts   map[string][]int // IP -> ports ouverts
    OS          map[string]string // IP -> OS detecte
    Services    map[string][]Service // IP -> services
    Credentials []Credential     // Identifiants trouves
    Vulns       []Vulnerability  // Vulnerabilites identifiees
    AttackPaths []AttackPath     // Chemins d'attaque possibles
}
```

### 3.3 Autonomous Planner (Le cerveau)

Remplace l'orchestrateur lineaire par une logique conditionnelle :

```go
// Concept - server/internal/domain/service/planner.go
func (p *AutonomousPlanner) DecideNextStep(knowledge TargetKnowledge) []PlannedTask {
    // Si on a des ports ouverts non explores -> scanner les services
    // Si on a detecte SMBv1 -> tenter EternalBlue
    // Si on a des credentials -> tenter mouvement lateral
    // Si on a un acces admin -> deployer un nouvel agent
    // Si rien ne marche -> tenter une autre approche
}
```

**Approche incrementale :**
1. D'abord des regles `if/else` simples (machine a etats)
2. Ensuite un moteur de regles configurable
3. Plus tard, potentiellement un LLM leger pour la decision

### 3.4 Schema Architecture V2

```
┌─────────────────────────────────────────────────────────────┐
│                    Server (Go) - Port 8443                    │
│  ┌───────────┐  ┌──────────────┐  ┌────────────────┐        │
│  │ Dashboard  │  │   REST API   │  │   WebSocket    │        │
│  │ (React)   │  │  /api/v1/*   │  │  /ws/*         │        │
│  └───────────┘  └──────────────┘  └────────────────┘        │
│                                                               │
│  ┌───────────────────────────────────────────────────┐       │
│  │              Domain Core                           │       │
│  │  ┌──────────────┐  ┌───────────────────────┐      │       │
│  │  │  Autonomous   │  │  Target Knowledge DB  │      │       │
│  │  │  Planner      │──│  (etat de la cible)   │      │       │
│  │  └──────────────┘  └───────────────────────┘      │       │
│  │  ┌──────────────┐  ┌───────────────────────┐      │       │
│  │  │  Attack       │  │  Score Calculator     │      │       │
│  │  │  Orchestrator │  │  + Remediation        │      │       │
│  │  └──────────────┘  └───────────────────────┘      │       │
│  └───────────────────────────────────────────────────┘       │
└─────────────────────────────────────────────────────────────┘
                              ▲
                              │ WebSocket (mTLS)
                    ┌─────────┴─────────┐
                    │   Agent (Rust)     │
                    │   + Pivot/Proxy    │
                    │   + Self-Deploy    │
                    └─────────┬─────────┘
                              │ SMB/SSH
                    ┌─────────┴─────────┐
                    │  Agent #2 (auto)   │
                    └───────────────────┘
```

---

## 4. Capacites Agent V2 (Rust)

### 4.1 Self-Propagation (Worm-like)

L'agent peut se copier sur une machine voisine si des identifiants sont trouves :

1. Agent A decouvre des credentials (via technique T1552 ou T1003)
2. Le serveur decide de deployer un nouvel agent
3. Agent A copie son binaire sur la machine B via SMB ou SSH
4. Agent A execute le binaire a distance (psexec, WMI, ou SSH)
5. Agent B se connecte au serveur et rejoint la flotte

**Interet :**
- **Resilience** : Si un agent tombe, les autres continuent
- **Pivoting** : Contourner la segmentation VLAN
- **Blast Radius** : Mesurer la vitesse de propagation d'une infection

### 4.2 Network Pivoting

L'agent agit comme un proxy pour permettre au serveur de scanner le reseau interne profond :
- Scanner le sous-reseau local
- Remonter la topologie au serveur
- Servir de relai pour les commandes vers des zones non directement accessibles

### 4.3 Safety Rails (Garde-fous)

L'autonomie ne veut pas dire chaos. L'agent a des limites strictes :
- **Interdiction** des commandes destructrices (`rm -rf`, chiffrement reel, wipe)
- **Mode simulation** : Prouve qu'il POURRAIT le faire sans le faire (flag "Fichier Crypte")
- **Kill switch global** : Bouton "Panic" pour deconnecter tous les agents instantanement
- **Approvals** : Les actions critiques (exploit kernel, propagation) necessitent une validation humaine

---

## 5. Dashboard V2 (React)

### 5.1 Attack Graph / Attack Map

Graphe de noeuds montrant la propagation en temps reel :
- **Noeuds** = machines (icones serveur/PC/IoT)
- **Liens** = vecteurs d'attaque utilises (fleches animees)
- **Couleurs** : bleu (inconnu) → orange (en cours) → rouge (compromis) → vert (protege)
- Au clic sur un lien : technique MITRE utilisee pour le passage

**Techno suggeree :** React Flow ou D3.js

### 5.2 Live Terminal IA

Flux d'activite montrant la "pensee" du Decision Engine :
```
14:02 - RECON   : Port 445 detecte sur 192.168.1.15
14:02 - DECIDE  : Tentative enumeration SMB via Agent #1
14:03 - ALERT   : Acces refuse. Tentative bruteforce compte 'Admin'
14:05 - SUCCESS : Acces obtenu. Deploiement Agent #2
14:06 - PIVOT   : Agent #2 connecte. Scan du sous-reseau 10.0.0.0/24
```

### 5.3 Blast Radius View

Vue specifique pour les environnements segmentes :
- Schema du reseau avec les differentes zones (DMZ, LAN, serveurs, postes de travail)
- Indicateur de danger par zone
- Impact business affiche (ex: "Serveur AD : COMPROMIS")

### 5.4 Achievements / Gamification

Scoring avance au-dela du Security Score actuel :
- **Scores 1-10** par technique executee
- **Badges** : "Domain Dominator", "Ghost in the Shell" (3 pivots sans detection), etc.
- **Leaderboard** : Comparer la resilience de plusieurs reseaux/zones

### 5.5 Dashboard Overview ameliore

- **Cards de vulnerabilites** par severite (Critical / High / Medium / Low)
- **Device Discovery** : Breakdown par OS, services, etat
- **Host Grid** : Grille des hotes avec IPs, statut, agent deploye

### 5.6 Details Panel enrichi

Panel lateral (slide-over) avec :
- **Insight** : Explication du risque dans le contexte
- **Remediation** : Recommandations concretes (mapping MITRE Mitigations)
- **Simulate Fix** : Clic pour simuler que la faille est patchee → la case MITRE passe du rouge au vert
- **Parametres de l'attaque** : Commande executee, output, duree

### 5.7 Approvals System

Workflow de validation pour les actions critiques :
- Notification "Pending Approval" sur le dashboard avant exploit/propagation
- Approbation par host ou par type d'action
- Mode "Full Auto" (aucune approbation) vs "Supervised" (approbation requise)
- Bouton "Panic" / Kill Switch global

### 5.8 Timeline & Report

- **Timeline chronologique** : Flux vertical de toutes les actions (Recon → Scan → Exploit → Pivot)
- **Export PDF** : Rapport avec logo, score global, kill chain, remediations (deja planifie)

---

## 6. Objectif Demo (Soutenance 2028)

### Le scenario "Wahou"

1. L'operateur entre une IP cible (ou un range) dans le dashboard
2. AutoStrike lance le mode Blackbox
3. Le jury voit en temps reel :
   - Le graphe de propagation s'etendre
   - Les logs de decision de l'IA defiler
   - La matrice MITRE se remplir
4. En quelques minutes : Scan → Exploit → Root → Rapport genere

### Objectif CTF

Capacite a resoudre une box HackTheBox (niveau Easy) en autonomie :
- Scan de ports → Detection du service vulnerable
- Choix automatique de l'exploit
- Obtention d'un shell
- Elevation de privileges
- Rapport complet genere

---

## 7. Repartition equipe V2

| Role | Responsabilites V2 |
|------|-------------------|
| **Architecture & Agent Lead** | Architecture, Decision Engine (Go), Agent avance (Rust) |
| **Frontend Lead** | Dashboard V2 (Attack Graph, Live Terminal, Blast Radius, Approvals) |
| **Backend Dev 1** | Backend Decision Engine (Go), Target Knowledge DB, Facts System |
| **Backend Dev 2** | Backend Planners (Go), API Reports, Integration tests |

---

## 8. Phases d'implementation (2026-2028)

> Projet EIP - Promotion 2028. Rendu final : fin d'annee scolaire 2028.

### S1 2026 (fev-juil) : Quick Wins + Fondations V2
- [ ] Profils APT (scenarios YAML)
- [ ] Export PDF rapports (backend + frontend)
- [ ] Cleanup techniques + auto-deploy scripts
- [ ] ScenarioBuilder visuel + LiveLogs
- [ ] **Target Knowledge DB** (entite + repository)
- [ ] **Moteur de regles basique** (if/else dans le planner)

### S2 2026 (sept-dec) : Dashboard V2 + Decision Engine
- [ ] Attack Graph / Attack Map (React Flow / D3.js)
- [ ] Live Terminal IA (WebSocket)
- [ ] Details Panel enrichi (Insight + Remediation + Simulate Fix)
- [ ] Approvals System (workflow validation)
- [ ] **Feedback loop** : le serveur reagit aux outputs de l'agent
- [ ] **Mode Blackbox API** (POST /api/v1/executions avec mode: "blackbox")
- [ ] Planners intelligents (sequential, conditional, buckets)

### S1 2027 (jan-juil) : Agent Avance + Caldera-like
- [ ] Network scan local (decouverte des voisins)
- [ ] Self-propagation via SMB/SSH
- [ ] Safety Rails (blocage commandes destructrices)
- [ ] Kill Switch global
- [ ] Facts / Data Exchange system
- [ ] Obfuscation des commandes
- [ ] Recommandations de remediation (mapping MITRE Mitigations)
- [ ] Dashboard Overview ameliore (vulnerability cards, host grid)

### S2 2027 (sept-dec) : Polish + Features avancees
- [ ] Blast Radius View (environnements segmentes)
- [ ] Achievements / Gamification
- [ ] Timeline chronologique + integration PDF
- [ ] Multiple agent types (Python leger, reverse shell)

### S1 2028 (jan-juin) : Autonomie Avancee + Soutenance
- [ ] Decision Engine v2 (moteur de regles avance, apprentissage des patterns)
- [ ] CTF autonome (objectif HackTheBox Easy)
- [ ] Mode Blackbox complet (Scan → Exploit → Root → Rapport)
- [ ] Tests E2E complets sur environnements lab
- [ ] Documentation finale + preparation soutenance EIP
- [ ] Demo "wahou" pour le jury

---

## 9. Ce qui ne change PAS

- L'architecture hexagonale du serveur Go
- Le protocole WebSocket agent ↔ serveur
- Le dashboard React/TypeScript/Tailwind
- L'agent Rust (on ajoute des capacites, on ne reecrit pas)
- Les 48 techniques MITRE deja implementees
- L'auth JWT/RBAC, le scheduling, les notifications
- Le mode "manuel" classique (scenarios predefinies) reste disponible

> **"On ne casse rien. On ajoute un cerveau par-dessus les fondations solides."**

---

*Derniere mise a jour: 2026-02-06*
