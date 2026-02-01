# Agent (Rust)

L'agent AutoStrike est développé en **Rust** pour des raisons de performance, sécurité mémoire et portabilité.

---

## Fonctionnalités

- **Exécution des techniques MITRE ATT&CK** via différents executors
- **Communication WebSocket** sécurisée avec reconnexion automatique
- **Multi-plateformes** : Windows, Linux, macOS (x64 et ARM64)
- **Détection automatique** de la plateforme et des executors disponibles
- **Auto-cleanup** après exécution des techniques
- **Heartbeat** périodique pour maintenir la connexion

---

## Structure

```
agent/
├── src/
│   ├── main.rs          # Point d'entrée, CLI (clap)
│   ├── config.rs        # Gestion configuration YAML
│   ├── client.rs        # Client WebSocket, protocole de communication
│   ├── executor.rs      # Exécution des commandes avec timeout
│   └── system.rs        # Détection système (OS, hostname, executors)
├── Cargo.toml           # Dépendances Rust
├── Cargo.lock
└── Dockerfile           # Build multi-stage
```

---

## Dépendances Principales

| Crate | Usage |
|-------|-------|
| `tokio` | Runtime async |
| `tokio-tungstenite` | Client WebSocket |
| `reqwest` | Client HTTP (avec rustls) |
| `serde` / `serde_json` | Sérialisation |
| `tracing` | Logging structuré |
| `sysinfo` | Informations système |
| `clap` | Parsing CLI |
| `anyhow` / `thiserror` | Gestion d'erreurs |

---

## Déploiement

### Windows

```powershell
.\autostrike-agent.exe --server https://server:8443 --paw AGENT_WIN_001
```

### Linux / macOS

```bash
./autostrike-agent --server https://server:8443 --paw AGENT_LIN_001
```

---

## Options CLI

| Option | Description | Défaut |
|--------|-------------|--------|
| `-s, --server` | URL du serveur AutoStrike | `https://localhost:8443` |
| `-p, --paw` | Identifiant unique de l'agent | UUID auto-généré |
| `-c, --config` | Chemin fichier de configuration | `agent.yaml` |
| `-d, --debug` | Activer les logs de debug | `false` |

---

## Configuration

Fichier `agent.yaml` :

```yaml
server_url: "https://server:8443"
paw: "agent-001"
heartbeat_interval: 30  # secondes

tls:
  cert_file: "./certs/agent.crt"
  key_file: "./certs/agent.key"
  ca_file: "./certs/ca.crt"
  verify: true
```

---

## Détection Système

L'agent détecte automatiquement au démarrage :

| Information | Méthode |
|-------------|---------|
| Hostname | `sysinfo::System::host_name()` |
| Username | `whoami::username()` |
| Platform | `cfg!(target_os)` |
| OS Version | `sysinfo::System::os_version()` |
| Architecture | `std::env::consts::ARCH` |

### Executors Détectés

| Windows | Linux/macOS |
|---------|-------------|
| `powershell` | `bash` |
| `pwsh` | `sh` |
| `cmd` | `zsh` |
| | `python3` |

---

## Protocole WebSocket

### Enregistrement (Agent → Server)
```json
{
  "type": "register",
  "payload": {
    "paw": "agent-001",
    "hostname": "DESKTOP-ABC",
    "username": "admin",
    "platform": "windows",
    "executors": ["powershell", "cmd"]
  }
}
```

### Réception Tâche (Server → Agent)
```json
{
  "type": "task",
  "payload": {
    "id": "task-uuid",
    "technique_id": "T1082",
    "command": "systeminfo",
    "executor": "cmd",
    "timeout": 300,
    "cleanup": "del /f output.txt"
  }
}
```

### Envoi Résultat (Agent → Server)
```json
{
  "type": "task_result",
  "payload": {
    "task_id": "task-uuid",
    "technique_id": "T1082",
    "success": true,
    "output": "Host Name: DESKTOP-ABC...",
    "exit_code": 0
  }
}
```

---

## Cross-Compilation

Script fourni pour compiler sur toutes les plateformes :

```bash
./scripts/build-agent.sh v0.1.0
```

### Targets Supportés

| Target | OS | Architecture |
|--------|-----|--------------|
| `x86_64-unknown-linux-gnu` | Linux | x64 |
| `x86_64-unknown-linux-musl` | Linux (statique) | x64 |
| `x86_64-pc-windows-gnu` | Windows | x64 |
| `x86_64-apple-darwin` | macOS | x64 |
| `aarch64-unknown-linux-gnu` | Linux | ARM64 |
| `aarch64-apple-darwin` | macOS M1/M2 | ARM64 |

---

## Sécurité

- **TLS/mTLS** : Communication chiffrée avec certificat client
- **Pas de credentials en dur** : Configuration via fichier ou CLI
- **Cleanup automatique** : Nettoyage après chaque technique
- **Mode non-root** : Exécution recommandée sans privilèges élevés
