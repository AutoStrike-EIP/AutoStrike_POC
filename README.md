# AutoStrike

[![CI Pipeline](https://github.com/AutoStrike-EIP/AutoStrike_POC/actions/workflows/ci.yml/badge.svg)](https://github.com/AutoStrike-EIP/AutoStrike_POC/actions/workflows/ci.yml)
[![SonarCloud](https://github.com/AutoStrike-EIP/AutoStrike_POC/actions/workflows/sonar.yml/badge.svg)](https://github.com/AutoStrike-EIP/AutoStrike_POC/actions/workflows/sonar.yml)

**Plateforme de Breach and Attack Simulation (BAS) basee sur MITRE ATT&CK**

> Projet EIP EPITECH - Promotion 2028

---

## Description

AutoStrike est une plateforme open-source de validation continue des defenses de securite par simulation d'attaques. Elle permet aux equipes SOC et RSSI de tester automatiquement leurs capacites de detection en utilisant le framework MITRE ATT&CK.

## Fonctionnalites

- **294 techniques MITRE ATT&CK** couvrant 12 tactiques (importees via MITRE STIX + Atomic Red Team)
- **Matrice MITRE interactive** - Visualisation de la couverture de detection
- **Scenarios d'attaque** - Execution automatisee de techniques par phases
- **Agents multi-plateformes** - Windows, Linux et macOS (Rust)
- **Dashboard temps reel** - Mises a jour WebSocket instantanees
- **Score de securite** - Evaluation automatique des defenses (0-100)
- **Authentification complete** - JWT, 5 roles, 28 permissions granulaires
- **Scheduling** - Executions planifiees (cron, daily, weekly, monthly)
- **Notifications** - Email SMTP + webhooks automatiques
- **Analytics** - Tendances, comparaisons de periodes, graphiques
- **Safe Mode** - Classification de securite per-executor (220 safe, 74 unsafe) avec detection de patterns dangereux
- **Security hardening** - Rate limiting, security headers, CSP, HSTS
- **Docker ready** - docker-compose prod + dev + 3 Dockerfiles multi-stage

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    Dashboard (React + TypeScript)                │
│                         https://server:8443                      │
└─────────────────────────────────────────────────────────────────┘
                                │
                    HTTPS + WebSocket (port 8443)
                                │
┌─────────────────────────────────────────────────────────────────┐
│                     Control Server (Go + Gin)                    │
│   ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐    │
│   │  REST API   │  │  WebSocket  │  │  Static Dashboard   │    │
│   │  /api/v1/*  │  │  /ws/*      │  │  /                  │    │
│   └─────────────┘  └─────────────┘  └─────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
                                │
                          WebSocket (TLS)
                                │
            ┌───────────────────┼───────────────────┐
            ▼                   ▼                   ▼
      ┌──────────┐        ┌──────────┐        ┌──────────┐
      │  Agent   │        │  Agent   │        │  Agent   │
      │ Windows  │        │  Linux   │        │  macOS   │
      │  (Rust)  │        │  (Rust)  │        │  (Rust)  │
      └──────────┘        └──────────┘        └──────────┘
```

## Stack Technique

| Composant | Technologie |
|-----------|-------------|
| Frontend | React 18, TypeScript, TailwindCSS, TanStack Query, Chart.js |
| Backend | Go 1.24+, Gin, gorilla/websocket, SQLite |
| Agent | Rust 1.75+, tokio, tokio-tungstenite |
| Communication | WebSocket (TLS), REST API |
| CI/CD | GitHub Actions, SonarCloud, Docker |

## Demarrage Rapide

```bash
# Cloner le projet
git clone https://github.com/AutoStrike-EIP/AutoStrike_POC.git
cd AutoStrike_POC

# Installer les dependances et build
make install

# Lancer AutoStrike (API + Dashboard sur port 8443)
make run

# Dans un autre terminal, connecter un agent
make agent
```

Ouvrir **https://localhost:8443** (accepter le certificat auto-signe)

### Commandes utiles

| Commande | Description |
|----------|-------------|
| `make install` | Installer dependances + build complet |
| `make run` | Demarrer le serveur (API + Dashboard) |
| `make agent` | Connecter un agent local |
| `make stop` | Arreter les services |
| `make test` | Lancer tous les tests |
| `make logs` | Voir les logs serveur |
| `make certs` | Generer certificats TLS (CA + server + agent) |
| `make docker-up` | Lancer via Docker Compose |

## Techniques MITRE ATT&CK

294 techniques importees couvrant 12 tactiques (apres `make import-mitre`) :

| Tactique | Count |
|----------|-------|
| **Initial Access** | 4 |
| **Execution** | 22 |
| **Persistence** | 44 |
| **Privilege Escalation** | 18 |
| **Defense Evasion** | 89 |
| **Credential Access** | 34 |
| **Discovery** | 30 |
| **Lateral Movement** | 8 |
| **Collection** | 16 |
| **Command and Control** | 13 |
| **Exfiltration** | 8 |
| **Impact** | 8 |

## Tests

```bash
# Tous les tests
make test

# Par composant
cd server && go test ./... -cover   # 200+ tests, 95%+ coverage domaine
cd agent && cargo test              # 67 tests
cd dashboard && npm test -- --run   # 1004 tests
```

**1270+ tests au total** couvrant server, agent et dashboard.

## Documentation

Voir la [documentation complete](https://autostrike-eip.github.io/AutoStrike/).

## Licence

Ce projet est sous licence MIT - voir le fichier [LICENSE](LICENSE) pour plus de details.
