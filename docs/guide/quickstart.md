# Quick Start Guide

This guide will help you run your first attack simulation in 5 minutes.

---

## Prerequisites

| Component | Version | Purpose |
|-----------|---------|---------|
| Go | 1.24+ | Server compilation |
| Node.js | 18+ | Dashboard build |
| Rust | 1.75+ | Agent compilation |

---

## 1. Installation

```bash
# Clone the repository
git clone https://github.com/AutoStrike-EIP/AutoStrike_POC.git
cd AutoStrike_POC

# Install dependencies and build
make install
```

---

## 2. Start AutoStrike

```bash
make run
```

The server starts on **https://localhost:8443** and serves:

| Path | Description |
|------|-------------|
| `/` | Dashboard (React SPA) |
| `/api/v1/*` | REST API |
| `/ws/agent` | WebSocket for agents |
| `/ws/dashboard` | WebSocket for real-time updates |
| `/health` | Health check endpoint |

!!! note "Authentication"
    Authentication is **disabled** by default in development.
    To enable it, set `JWT_SECRET` in `server/.env`.

---

## 3. Connect an Agent

Open a new terminal and run:

```bash
make agent
```

This compiles and runs the Rust agent, which will:

1. Connect to `wss://localhost:8443/ws/agent`
2. Register with hostname, platform, and available executors
3. Start sending heartbeats every 30 seconds
4. Wait for tasks from the server

Verify the agent is connected in the **Agents** page (status should be "Online").

!!! tip "Agent Authentication"
    If `AGENT_SECRET` is set on the server, pass it to the agent:
    ```bash
    cd agent && cargo run --release -- -k "your-agent-secret"
    ```

---

## 4. Run a Scenario

1. Navigate to **Scenarios** in the dashboard
2. Click **Run** on "Discovery - Basic" (or any available scenario)
3. In the modal:
   - Select target agents (checkboxes)
   - Toggle **Safe Mode** (recommended for first run)
4. Click **Run Execution**

The execution will start and you'll be redirected to the Executions page.

---

## 5. Monitor Execution

The **Executions** page shows real-time progress:

| Status | Description |
|--------|-------------|
| **Pending** | Execution created, waiting to start |
| **Running** | Techniques being executed on agents |
| **Completed** | All techniques finished |
| **Cancelled** | Execution stopped by user |

Click on an execution row to see detailed results.

---

## 6. Analyze Results

In **Execution Details**, each technique shows:

| Status | Meaning | Score Impact |
|--------|---------|--------------|
| **Blocked** | Attack was prevented by security controls | +100 points |
| **Detected** | Attack succeeded but was detected | +50 points |
| **Success** | Attack succeeded undetected | +0 points |
| **Failed** | Technique execution failed | Not counted |

**Security Score Formula:**
```
Score = (blocked * 100 + detected * 50) / (total * 100) * 100%
```

Higher scores indicate better security posture.

---

## 7. Explore the MITRE Matrix

The **ATT&CK Matrix** page provides a visual overview:

| Color | Meaning |
|-------|---------|
| Green cell | Techniques available and safe |
| Red dot | Technique is unsafe (use with caution) |
| Click | View technique details and detection info |

Filter by platform (Windows/Linux) using the dropdown.

---

## Useful Commands

| Command | Description |
|---------|-------------|
| `make run` | Start the server |
| `make agent` | Connect a local agent |
| `make stop` | Stop all services |
| `make logs` | View server logs |
| `make test` | Run all tests |
| `make build` | Build all components |

---

## Available Techniques (48 total)

AutoStrike includes **48 MITRE ATT&CK techniques** across **13 tactics**:

### Reconnaissance (2 techniques)
- T1592.004 - Gather Victim Host Information: Client Configurations
- T1595.002 - Active Scanning: Vulnerability Scanning

### Initial Access (3 techniques)
- T1078 - Valid Accounts
- T1133 - External Remote Services
- T1190 - Exploit Public-Facing Application

### Execution (5 techniques)
- T1059.001 - PowerShell
- T1059.003 - Windows Command Shell
- T1059.004 - Unix Shell
- T1047 - Windows Management Instrumentation
- T1059.006 - Python

### Persistence (4 techniques)
- T1053.005 - Scheduled Task
- T1547.001 - Registry Run Keys
- T1053.003 - Cron
- T1543.002 - Systemd Service

### Privilege Escalation (4 techniques)
- T1548.001 - Setuid and Setgid
- T1548.002 - Bypass User Account Control
- T1078.003 - Local Accounts
- T1134.001 - Token Impersonation/Theft

### Defense Evasion (6 techniques)
- T1070.004 - File Deletion
- T1562.001 - Disable or Modify Tools
- T1027 - Obfuscated Files or Information
- T1070.001 - Clear Windows Event Logs
- T1036.005 - Match Legitimate Name or Location
- T1218.011 - Rundll32

### Credential Access (4 techniques)
- T1552.001 - Credentials In Files
- T1555.003 - Credentials from Web Browsers
- T1003.008 - /etc/passwd and /etc/shadow
- T1552.004 - Private Keys

### Discovery (9 techniques)
- T1082 - System Information Discovery
- T1083 - File and Directory Discovery
- T1057 - Process Discovery
- T1016 - System Network Configuration Discovery
- T1049 - System Network Connections Discovery
- T1087 - Account Discovery
- T1069 - Permission Groups Discovery
- T1018 - Remote System Discovery
- T1007 - System Service Discovery

### Lateral Movement (3 techniques)
- T1021.001 - Remote Desktop Protocol
- T1021.002 - SMB/Windows Admin Shares
- T1021.004 - SSH

### Collection (4 techniques)
- T1005 - Data from Local System
- T1039 - Data from Network Shared Drive
- T1074.001 - Local Data Staging
- T1119 - Automated Collection

### Command and Control (3 techniques)
- T1071.001 - Web Protocols
- T1105 - Ingress Tool Transfer
- T1572 - Protocol Tunneling

### Exfiltration (3 techniques)
- T1048.003 - Exfiltration Over Unencrypted Non-C2 Protocol
- T1041 - Exfiltration Over C2 Channel
- T1567.002 - Exfiltration to Cloud Storage

### Impact (3 techniques)
- T1490 - Inhibit System Recovery
- T1489 - Service Stop
- T1486 - Data Encrypted for Impact

All techniques support **Safe Mode** for production-safe testing.

---

## Troubleshooting

### Agent won't connect

1. Check server is running: `curl -k https://localhost:8443/health`
2. Verify WebSocket endpoint: `wss://localhost:8443/ws/agent`
3. Check agent logs for connection errors
4. Ensure no firewall blocking port 8443
5. If `AGENT_SECRET` is set, verify the agent passes `-k` with the correct secret

### Execution stuck in "Running"

1. Check if agent is still online in Agents page
2. Agent may have disconnected - reconnect with `make agent`
3. Click **Stop** to cancel the execution

### Dashboard not loading

1. Verify server is running on port 8443
2. Clear browser cache
3. Check browser console for errors
4. Try incognito/private window

---

## Next Steps

- [Create a custom scenario](../mitre/techniques.md)
- [API Reference](../api/reference.md)
- [Architecture Overview](../architecture/index.md)
- [Deployment Guide](installation.md)
