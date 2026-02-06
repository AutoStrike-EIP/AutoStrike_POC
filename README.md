# AutoStrike

[![CI Pipeline](https://github.com/AutoStrike-EIP/AutoStrike_POC/actions/workflows/ci.yml/badge.svg)](https://github.com/AutoStrike-EIP/AutoStrike_POC/actions/workflows/ci.yml)
[![SonarCloud](https://github.com/AutoStrike-EIP/AutoStrike_POC/actions/workflows/sonar.yml/badge.svg)](https://github.com/AutoStrike-EIP/AutoStrike_POC/actions/workflows/sonar.yml)

**Plateforme de Breach and Attack Simulation (BAS) basee sur MITRE ATT&CK**

> Projet EIP EPITECH - Promotion 2028

---

## Description

AutoStrike est une plateforme open-source de validation continue des defenses de securite par simulation d'attaques. Elle permet aux equipes SOC et RSSI de tester automatiquement leurs capacites de detection en utilisant le framework MITRE ATT&CK.

## Fonctionnalites

- **48 techniques MITRE ATT&CK** couvrant 13 tactiques sur 14 (de Reconnaissance a Impact)
- **Matrice MITRE interactive** - Visualisation de la couverture de detection
- **Scenarios d'attaque** - Execution automatisee de techniques par phases
- **Agents multi-plateformes** - Windows, Linux et macOS (Rust)
- **Dashboard temps reel** - Mises a jour WebSocket instantanees
- **Score de securite** - Evaluation automatique des defenses (0-100)
- **Authentification complete** - JWT, 5 roles, 28 permissions granulaires
- **Scheduling** - Executions planifiees (cron, daily, weekly, monthly)
- **Notifications** - Email SMTP + webhooks automatiques
- **Analytics** - Tendances, comparaisons de periodes, graphiques
- **Safe Mode** - Toutes les techniques sont non-destructives
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

48 techniques implementees couvrant 13 tactiques :

| Tactique | Techniques |
|----------|------------|
| **Reconnaissance** (2) | T1592.004, T1595.002 |
| **Initial Access** (3) | T1078, T1133, T1190 |
| **Execution** (5) | T1059.001, T1059.003, T1059.004, T1047, T1059.006 |
| **Persistence** (4) | T1053.005, T1547.001, T1053.003, T1543.002 |
| **Privilege Escalation** (4) | T1548.001, T1548.002, T1078.003, T1134.001 |
| **Defense Evasion** (6) | T1070.004, T1562.001, T1027, T1070.001, T1036.005, T1218.011 |
| **Credential Access** (4) | T1552.001, T1555.003, T1003.008, T1552.004 |
| **Discovery** (9) | T1082, T1083, T1057, T1016, T1049, T1087, T1069, T1018, T1007 |
| **Lateral Movement** (3) | T1021.001, T1021.002, T1021.004 |
| **Collection** (4) | T1005, T1039, T1074.001, T1119 |
| **Command and Control** (3) | T1071.001, T1105, T1572 |
| **Exfiltration** (3) | T1048.003, T1041, T1567.002 |
| **Impact** (3) | T1490, T1489, T1486 |

## Tests

```bash
# Tous les tests
make test

# Par composant
cd server && go test ./... -cover   # 200+ tests, 95%+ coverage domaine
cd agent && cargo test              # 67 tests
cd dashboard && npm test -- --run   # 513 tests
```

**780+ tests au total** couvrant server, agent et dashboard.

## Documentation

Voir la [documentation complete](https://autostrike-eip.github.io/AutoStrike/).

## Licence

Ce projet est sous licence MIT - voir le fichier [LICENSE](LICENSE) pour plus de details.
