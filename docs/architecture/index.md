# Architecture

AutoStrike utilise une architecture N-Tiers avec trois composants principaux.

---

## Vue d'ensemble

```
┌─────────────────────────────────────────────────────────────────────┐
│                     COUCHE PRÉSENTATION                              │
│                     Dashboard (React + TypeScript)                   │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                            HTTPS / WebSocket
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      COUCHE SERVICE (API)                           │
│                     Control Server (Go + Gin)                        │
│                                                                      │
│   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                 │
│   │  REST API   │  │  WebSocket  │  │ Orchestrator│                 │
│   └─────────────┘  └─────────────┘  └─────────────┘                 │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                              mTLS / gRPC
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                          AGENTS (Rust)                              │
│                                                                      │
│   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                 │
│   │   Windows   │  │    Linux    │  │    Linux    │                 │
│   │   Agent     │  │    Agent    │  │    Agent    │                 │
│   └─────────────┘  └─────────────┘  └─────────────┘                 │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Composants

| Composant | Langage | Rôle |
|-----------|---------|------|
| [Dashboard](dashboard.md) | React/TypeScript | Interface utilisateur |
| [Backend](backend.md) | Go | API, orchestration, stockage |
| [Agent](agent.md) | Rust | Exécution des techniques sur les endpoints |

---

## Communication

- **Dashboard ↔ Backend** : HTTPS + WebSocket (temps réel)
- **Backend ↔ Agents** : mTLS (authentification mutuelle)
