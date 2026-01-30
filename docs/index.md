# AutoStrike

## Plateforme de Breach and Attack Simulation (BAS)

**Projet EIP EPITECH - Promotion 2028**

---

## Qu'est-ce qu'AutoStrike ?

AutoStrike est une plateforme open-source de **validation continue des défenses de sécurité** par simulation d'attaques basées sur le framework **MITRE ATT&CK**.

### Fonctionnalités principales

- **Matrice MITRE ATT&CK** - Visualisation de la couverture de détection
- **Scénarios d'attaque** - Exécution automatisée de techniques
- **Agents multi-plateformes** - Windows et Linux
- **Dashboard temps réel** - Suivi des exécutions et résultats
- **Rapports** - Export PDF pour les audits

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Dashboard (React)                         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                  Control Server (Go)                         │
│                                                              │
│   • REST API          • WebSocket          • Orchestration  │
└─────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
        ┌──────────┐   ┌──────────┐   ┌──────────┐
        │  Agent   │   │  Agent   │   │  Agent   │
        │ (Rust)   │   │ (Rust)   │   │ (Rust)   │
        │ Windows  │   │  Linux   │   │  Linux   │
        └──────────┘   └──────────┘   └──────────┘
```

---

## Stack Technique

| Composant | Technologie |
|-----------|-------------|
| **Frontend** | React 18, TypeScript, TailwindCSS, D3.js |
| **Backend** | Go 1.21, Gin Framework |
| **Agent** | Rust |
| **Database** | SQLite (MVP) → PostgreSQL |
| **Communication** | REST API, WebSocket, mTLS |

---

## Démarrage rapide

Consultez le [Guide d'installation](guide/installation.md) pour commencer.

---

## Liens

- [GitHub Repository](https://github.com/AutoStrike-EIP/AutoStrike)
- [MITRE ATT&CK Framework](https://attack.mitre.org/)

