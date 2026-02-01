# AutoStrike Agent

Agent Rust déployé sur les machines cibles pour exécuter les techniques MITRE ATT&CK.

## Architecture

```
agent/
├── src/
│   ├── main.rs          # Point d'entrée, CLI parsing
│   ├── config.rs        # Gestion configuration
│   ├── client.rs        # Client WebSocket, communication serveur
│   ├── executor.rs      # Exécution des commandes
│   └── system.rs        # Détection système (OS, executors)
├── Cargo.toml
└── Dockerfile
```

## Fonctionnalités

- **Connexion WebSocket** avec reconnexion automatique
- **Détection automatique** de la plateforme et des executors disponibles
- **Exécution de commandes** avec timeout et capture de sortie
- **Heartbeat** périodique pour maintenir la connexion
- **Cross-compilation** pour Windows, Linux, macOS (x64/ARM64)

## Prérequis

- Rust 1.75+
- OpenSSL (pour le build)

## Installation

```bash
# Installer les dépendances
cargo fetch

# Build debug
cargo build

# Build release (optimisé)
cargo build --release
```

## Utilisation

```bash
# Avec paramètres CLI
./autostrike-agent --server https://server:8443 --paw my-agent-id

# Avec fichier de configuration
./autostrike-agent --config agent.yaml

# Mode debug
./autostrike-agent --server https://server:8443 --debug
```

### Options CLI

| Option | Description | Défaut |
|--------|-------------|--------|
| `-s, --server` | URL du serveur | `https://localhost:8443` |
| `-p, --paw` | Identifiant unique de l'agent | UUID généré |
| `-c, --config` | Chemin du fichier de configuration | `agent.yaml` |
| `-d, --debug` | Activer les logs de debug | `false` |

## Configuration

Fichier `agent.yaml` :

```yaml
server_url: "https://server:8443"
paw: "agent-001"
heartbeat_interval: 30

tls:
  cert_file: "./certs/agent.crt"
  key_file: "./certs/agent.key"
  ca_file: "./certs/ca.crt"
  verify: true
```

## Détection Système

L'agent détecte automatiquement :

| Information | Source |
|-------------|--------|
| Hostname | `sysinfo` |
| Username | `whoami` |
| Platform | `cfg!(target_os)` |
| OS Version | `sysinfo` |
| Architecture | `std::env::consts::ARCH` |

### Executors Détectés

| Windows | Linux/macOS |
|---------|-------------|
| powershell | bash |
| pwsh | sh |
| cmd | zsh |
| | python3 |

## Protocole WebSocket

### Enregistrement
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

### Réception de tâche
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

### Envoi du résultat
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

## Cross-Compilation

Utiliser le script fourni :

```bash
cd ..
./scripts/build-agent.sh v0.1.0
```

Targets supportés :
- `x86_64-unknown-linux-gnu`
- `x86_64-unknown-linux-musl`
- `x86_64-pc-windows-gnu`
- `x86_64-apple-darwin`
- `aarch64-unknown-linux-gnu`
- `aarch64-apple-darwin`

## Tests

```bash
cargo test
```

## Docker

```bash
docker build -t autostrike-agent .
docker run autostrike-agent --server https://host.docker.internal:8443
```

## Sécurité

- Communication TLS/mTLS avec le serveur
- Pas de stockage de credentials en dur
- Exécution en tant qu'utilisateur non-root recommandée
- Cleanup automatique après exécution des techniques
