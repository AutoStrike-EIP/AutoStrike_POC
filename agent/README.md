# AutoStrike Agent

Agent Rust déployé sur les machines cibles pour exécuter les techniques MITRE ATT&CK.

## Architecture

```
agent/
├── src/
│   ├── main.rs          # Point d'entrée, CLI parsing (clap)
│   ├── config.rs        # Gestion configuration YAML
│   ├── client.rs        # Client WebSocket, communication serveur
│   ├── executor.rs      # Exécution des commandes avec timeout
│   └── system.rs        # Détection système (OS, executors)
├── Cargo.toml
└── Dockerfile
```

## Fonctionnalités

- **Connexion WebSocket** avec reconnexion automatique (backoff exponentiel 1s → 60s)
- **Détection automatique** de la plateforme et des executors disponibles
- **Exécution de commandes** avec timeout et capture de sortie
- **Heartbeat** périodique pour maintenir la connexion (30 secondes)
- **Authentication agent** via header `X-Agent-Key`
- **Troncature de sortie** à 1 MB max (limite UTF-8 safe)
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

# Avec authentification agent
./autostrike-agent --server https://server:8443 --paw my-agent-id -k "your-agent-secret"

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
| `-k, --agent-secret` | Secret d'authentification agent (header `X-Agent-Key`) | - |

## Configuration

Fichier `agent.yaml` :

```yaml
server_url: "https://server:8443"
paw: "agent-001"
heartbeat_interval: 30
agent_secret: "your-agent-secret"  # optionnel

tls:
  cert_file: "./certs/agent.crt"
  key_file: "./certs/agent.key"
  ca_file: "./certs/ca.crt"
  verify: true
```

**Priorité :** Arguments CLI > Fichier de configuration > Défauts

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

## Exécution de Commandes

### Timeout
- Timeout configurable par commande (défaut: 300 secondes)
- En cas de timeout: `success: false`, `output: "Command timed out"`

### Troncature de Sortie
- Taille max: **1 MB** (1,048,576 octets)
- Troncature à une frontière UTF-8 valide
- Message `"\n... [output truncated]"` ajouté si tronqué

### Capture de Sortie
- stdout et stderr capturés séparément puis combinés
- Décodage UTF-8 avec conversion lossy
- Whitespace en début/fin supprimé

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

```bash
# Linux x64
cargo build --release --target x86_64-unknown-linux-gnu

# Linux ARM64
cargo build --release --target aarch64-unknown-linux-gnu

# Windows x64
cargo build --release --target x86_64-pc-windows-gnu

# macOS x64
cargo build --release --target x86_64-apple-darwin

# macOS ARM64 (M1/M2)
cargo build --release --target aarch64-apple-darwin
```

Targets supportés :
- `x86_64-unknown-linux-gnu`
- `x86_64-unknown-linux-musl`
- `x86_64-pc-windows-gnu`
- `x86_64-apple-darwin`
- `aarch64-unknown-linux-gnu`
- `aarch64-apple-darwin`

## Tests

67 tests unitaires :

```bash
cargo test
cargo test -- --nocapture  # Avec sortie
```

## Docker

```bash
docker build -t autostrike-agent .
docker run autostrike-agent --server https://host.docker.internal:8443
```

## Sécurité

- Communication TLS/mTLS avec le serveur
- Authentification agent via header `X-Agent-Key`
- Pas de stockage de credentials en dur
- Exécution en tant qu'utilisateur non-root recommandée
- Cleanup automatique après exécution des techniques
- Protection timeout contre les commandes bloquées
- Troncature de sortie pour éviter l'épuisement mémoire
