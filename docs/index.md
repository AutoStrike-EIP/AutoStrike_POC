# AutoStrike

## Breach and Attack Simulation (BAS) Platform

**EIP Project - EPITECH Promotion 2028**

---

## What is AutoStrike?

AutoStrike is an open-source platform for **continuous security defense validation** through attack simulations based on the **MITRE ATT&CK** framework.

### Key Features

| Feature | Description |
|---------|-------------|
| **294 MITRE ATT&CK Techniques** | 12 tactics from Initial Access to Impact (imported from MITRE STIX + Atomic Red Team) |
| **Interactive MITRE Matrix** | Visual detection coverage with 12 tactic columns |
| **Attack Scenarios** | Automated technique execution with phases |
| **Multi-platform Agents** | Windows, Linux, and macOS support (Rust) |
| **Real-time Dashboard** | Live execution monitoring via WebSocket |
| **Security Scoring** | Quantified defense effectiveness (0-100) |
| **Authentication & RBAC** | JWT tokens, 5 roles, 28 granular permissions |
| **Scheduling** | Automated executions (cron, daily, weekly, monthly) |
| **Notifications** | Email SMTP + webhook alerts |
| **Analytics** | Score trends, period comparisons, charts |
| **Safe Mode** | Per-executor safety classification (220 safe, 74 unsafe) with dangerous pattern detection |
| **Security Hardening** | Rate limiting, security headers, CSP, HSTS |

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│              Control Server (Go) - Port 8443                │
│                                                             │
│   ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌───────────┐  │
│   │Dashboard │  │ REST API │  │WebSocket │  │Orchestrat.│  │
│   │ (React)  │  │/api/v1/* │  │  /ws/*   │  │           │  │
│   └──────────┘  └──────────┘  └──────────┘  └───────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
        ┌──────────┐   ┌──────────┐   ┌──────────┐
        │  Agent   │   │  Agent   │   │  Agent   │
        │ (Rust)   │   │ (Rust)   │   │ (Rust)   │
        │ Windows  │   │  Linux   │   │  macOS   │
        └──────────┘   └──────────┘   └──────────┘
```

A single server on **port 8443** serves the Dashboard, REST API, and WebSocket connections.

---

## Tech Stack

| Component | Technology |
|-----------|------------|
| **Frontend** | React 18, TypeScript, TailwindCSS, TanStack Query, Chart.js |
| **Backend** | Go 1.24+, Gin Framework, SQLite |
| **Agent** | Rust 1.75+, Tokio, tokio-tungstenite |
| **Communication** | REST API, WebSocket (real-time), TLS |
| **CI/CD** | GitHub Actions, SonarCloud, Docker |

---

## Quick Start

```bash
# Clone and install
git clone https://github.com/AutoStrike-EIP/AutoStrike_POC.git
cd AutoStrike_POC
make install

# Start server and dashboard
make run

# Connect an agent (in another terminal)
make agent
```

Access the dashboard at **https://localhost:8443**

See the [Quick Start Guide](guide/quickstart.md) for detailed instructions.

---

## Dashboard Pages

| Page | Description |
|------|-------------|
| **Dashboard** | Overview with stats, security score, recent activity |
| **Agents** | Connected agents with status and deployment commands |
| **Techniques** | Browse and import MITRE ATT&CK techniques |
| **ATT&CK Matrix** | Interactive 12-column MITRE matrix visualization |
| **Scenarios** | Attack scenarios with phases, import/export, Run |
| **Executions** | Execution history with real-time WebSocket updates |
| **Analytics** | Score trends, period comparisons, charts |
| **Scheduler** | Schedule automated scenario executions (cron, daily, etc.) |
| **Settings** | Notifications, SMTP, execution defaults |
| **Admin/Users** | User management with 5 roles (admin only) |
| **Admin/Permissions** | Role-based permission matrix (admin only) |

---

## Current Techniques (294 total)

After running `make import-mitre`, AutoStrike provides **294 techniques** across 12 tactics:

| Tactic | Count |
|--------|-------|
| Initial Access | 4 |
| Execution | 22 |
| Persistence | 44 |
| Privilege Escalation | 18 |
| Defense Evasion | 89 |
| Credential Access | 34 |
| Discovery | 30 |
| Lateral Movement | 8 |
| Collection | 16 |
| Command and Control | 13 |
| Exfiltration | 8 |
| Impact | 8 |

**Safety classification:** 220 safe techniques, 74 unsafe. Per-executor safety based on elevation requirements and dangerous command pattern detection. Use `make import-mitre-safe` to import only safe techniques.

---

## API Overview

### Authentication
| Endpoint | Description |
|----------|-------------|
| `POST /api/v1/auth/login` | Login (username/password) |
| `POST /api/v1/auth/refresh` | Refresh access token |
| `POST /api/v1/auth/logout` | Invalidate tokens |
| `GET /api/v1/auth/me` | Get current user info |

### Core API
| Endpoint | Description |
|----------|-------------|
| `GET /api/v1/agents` | List connected agents |
| `GET /api/v1/techniques` | List MITRE techniques |
| `GET /api/v1/techniques/coverage` | MITRE coverage statistics |
| `GET /api/v1/scenarios` | List attack scenarios |
| `POST /api/v1/executions` | Start an execution |
| `GET /api/v1/executions/:id/results` | Get execution results |
| `GET /api/v1/analytics/trend` | Score trend over time |
| `GET /api/v1/schedules` | List scheduled executions |
| `GET /api/v1/notifications` | Get notifications |
| `GET /api/v1/admin/users` | User management (admin) |
| `GET /api/v1/permissions/matrix` | Permission matrix |

See the complete [API Reference](api/reference.md).

---

## WebSocket Events

Real-time updates via `wss://localhost:8443/ws/dashboard`:

- `execution_started` - Execution began
- `execution_completed` - Execution finished
- `execution_cancelled` - Execution stopped

---

## Security Score

Measures how well your defenses perform:

| Result | Points | Meaning |
|--------|--------|---------|
| Blocked | 100 | Attack prevented |
| Detected | 50 | Attack seen but not stopped |
| Success | 0 | Attack succeeded undetected |

**Formula:** `(blocked*100 + detected*50) / (total*100) * 100%`

---

## Links

- [GitHub Repository](https://github.com/AutoStrike-EIP/AutoStrike_POC)
- [MITRE ATT&CK Framework](https://attack.mitre.org/)
- [API Reference](api/reference.md)
- [Quick Start Guide](guide/quickstart.md)
