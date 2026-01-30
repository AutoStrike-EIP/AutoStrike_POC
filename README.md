# AutoStrike

**Plateforme de Breach and Attack Simulation (BAS) basée sur MITRE ATT&CK**

> Projet EIP EPITECH - Promotion 2028

---

## Description

AutoStrike est une plateforme open-source de validation continue des défenses de sécurité par simulation d'attaques. Elle permet aux équipes SOC et RSSI de tester automatiquement leurs capacités de détection en utilisant le framework MITRE ATT&CK.

## Fonctionnalités

- **Matrice MITRE ATT&CK** - Visualisation de la couverture de détection
- **Scénarios d'attaque** - Exécution automatisée de techniques
- **Agents multi-plateformes** - Windows et Linux
- **Dashboard temps réel** - Suivi des exécutions et résultats
- **Rapports** - Export PDF pour les audits

## Architecture

```
┌─────────────────────────────────────────┐
│           Dashboard (React)              │
└─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│         Control Server (Go)              │
└─────────────────────────────────────────┘
                    │
        ┌───────────┼───────────┐
        ▼           ▼           ▼
   ┌─────────┐ ┌─────────┐ ┌─────────┐
   │  Agent  │ │  Agent  │ │  Agent  │
   │ (Rust)  │ │ (Rust)  │ │ (Rust)  │
   └─────────┘ └─────────┘ └─────────┘
```

## Stack Technique

| Composant | Technologie |
|-----------|-------------|
| Frontend | React 18, TypeScript, TailwindCSS |
| Backend | Go 1.21, Gin Framework |
| Agent | Rust |
| Database | SQLite / PostgreSQL |

## Installation

Voir la [documentation complète](https://autostrike-eip.github.io/AutoStrike/).

## Équipe

Projet EIP EPITECH - AutoStrike-EIP

## Licence

Ce projet est sous licence MIT - voir le fichier [LICENSE](LICENSE) pour plus de détails.
