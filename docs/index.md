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
│              Control Server (Go) - Port 8443                 │
│                                                              │
│   ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌───────────┐  │
│   │Dashboard │  │ REST API │  │WebSocket │  │Orchestrat.│  │
│   │ (Static) │  │/api/v1/* │  │  /ws/*   │  │           │  │
│   └──────────┘  └──────────┘  └──────────┘  └───────────┘  │
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

Un seul serveur sur le **port 8443** sert le Dashboard, l'API REST et les WebSockets.

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

```bash
make run    # Démarre sur http://localhost:8443
make agent  # Connecte un agent
make stop   # Arrête les services
```

Consultez le [Guide de démarrage](guide/quickstart.md) pour plus de détails.

---

## Liens

- [GitHub Repository](https://github.com/AutoStrike-EIP/AutoStrike)
- [MITRE ATT&CK Framework](https://attack.mitre.org/)

