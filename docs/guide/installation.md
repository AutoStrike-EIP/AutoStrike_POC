# Installation

## Prérequis

### Serveur

- Go 1.21+
- Node.js 20+
- SQLite ou PostgreSQL

### Agent

- Windows 10+ ou Linux (Ubuntu 20+)
- Droits administrateur

---

## Installation du serveur

### 1. Cloner le repository

```bash
git clone https://github.com/AutoStrike-EIP/AutoStrike.git
cd AutoStrike
```

### 2. Backend (Go)

```bash
cd server
go mod download
go build -o autostrike-server ./cmd/autostrike
./autostrike-server
```

### 3. Frontend (React)

```bash
cd dashboard
npm install
npm run build
npm run preview
```

---

## Installation de l'agent

### Windows (PowerShell Admin)

```powershell
# Télécharger et installer l'agent
Invoke-WebRequest -Uri "https://server:8443/deploy/agent.exe" -OutFile "autostrike-agent.exe"
.\autostrike-agent.exe --server https://server:8443
```

### Linux

```bash
# Télécharger et installer l'agent
curl -o autostrike-agent https://server:8443/deploy/agent
chmod +x autostrike-agent
sudo ./autostrike-agent --server https://server:8443
```

---

## Configuration

Voir la section [Configuration](../architecture/backend.md) pour les options avancées.
