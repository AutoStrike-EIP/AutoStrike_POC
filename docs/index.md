# AutoStrike

## Breach and Attack Simulation (BAS) Platform

**EIP Project - EPITECH Promotion 2028**

---

## What is AutoStrike?

AutoStrike is an open-source platform for **continuous security defense validation** through attack simulations based on the **MITRE ATT&CK** framework.

### Key Features

| Feature | Description |
|---------|-------------|
| **48 MITRE ATT&CK Techniques** | 13 tactics covering Reconnaissance to Impact |
| **Interactive MITRE Matrix** | Visual detection coverage with 14 tactic columns |
| **Attack Scenarios** | Automated technique execution with phases |
| **Multi-platform Agents** | Windows, Linux, and macOS support (Rust) |
| **Real-time Dashboard** | Live execution monitoring via WebSocket |
| **Security Scoring** | Quantified defense effectiveness (0-100) |
| **Authentication & RBAC** | JWT tokens, 5 roles, 28 granular permissions |
| **Scheduling** | Automated executions (cron, daily, weekly, monthly) |
| **Notifications** | Email SMTP + webhook alerts |
| **Analytics** | Score trends, period comparisons, charts |
| **Safe Mode** | All techniques are non-destructive |
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
| **ATT&CK Matrix** | Interactive 14-column MITRE matrix visualization |
| **Scenarios** | Attack scenarios with phases, import/export, Run |
| **Executions** | Execution history with real-time WebSocket updates |
| **Analytics** | Score trends, period comparisons, charts |
| **Scheduler** | Schedule automated scenario executions (cron, daily, etc.) |
| **Settings** | Notifications, SMTP, execution defaults |
| **Admin/Users** | User management with 5 roles (admin only) |
| **Admin/Permissions** | Role-based permission matrix (admin only) |

---

## Current Techniques (48 total)

### Reconnaissance (2)
T1592.004, T1595.002

### Initial Access (3)
T1078, T1133, T1190

### Execution (5)
T1059.001, T1059.003, T1059.004, T1047, T1059.006

### Persistence (4)
T1053.005, T1547.001, T1053.003, T1543.002

### Privilege Escalation (4)
T1548.001, T1548.002, T1078.003, T1134.001

### Defense Evasion (6)
T1070.004, T1562.001, T1027, T1070.001, T1036.005, T1218.011

### Credential Access (4)
T1552.001, T1555.003, T1003.008, T1552.004

### Discovery (9)
T1082, T1083, T1057, T1016, T1049, T1087, T1069, T1018, T1007

### Lateral Movement (3)
T1021.001, T1021.002, T1021.004

### Collection (4)
T1005, T1039, T1074.001, T1119

### Command and Control (3)
T1071.001, T1105, T1572

### Exfiltration (3)
T1048.003, T1041, T1567.002

### Impact (3)
T1490, T1489, T1486

All techniques support **Safe Mode** for production-safe testing.

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
