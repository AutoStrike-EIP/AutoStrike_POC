# Installation

## Prérequis

### Serveur

| Composant | Version minimale | Recommandé |
|-----------|------------------|------------|
| Go | 1.21+ | 1.22+ |
| Node.js | 20+ | 22 LTS |
| SQLite | 3.35+ | 3.40+ |
| PostgreSQL | 13+ (optionnel) | 16+ |
| Docker | 24+ (optionnel) | 25+ |
| Docker Compose | 2.20+ (optionnel) | 2.24+ |

### Agent

| OS | Version minimale | Architecture |
|-----|------------------|--------------|
| Windows | 10 (1809+) | x64, ARM64 |
| Linux | Ubuntu 20.04+ / Debian 11+ | x64, ARM64 |
| macOS | 12 (Monterey)+ | x64, ARM64 |

**Privilèges requis :**
- Windows : Droits administrateur (pour certaines techniques)
- Linux : Accès root ou sudo (pour certaines techniques)
- macOS : Droits administrateur (pour certaines techniques)

---

## Installation du serveur

### Option 1 : Installation manuelle

#### 1. Cloner le repository

```bash
git clone https://github.com/AutoStrike-EIP/AutoStrike.git
cd AutoStrike
```

#### 2. Configuration des variables d'environnement

Créez un fichier `.env` à la racine du projet :

```bash
# Copier le fichier d'exemple
cp .env.example .env

# Éditer les variables
nano .env
```

Variables requises :

```env
# JWT Configuration
JWT_SECRET=votre-secret-jwt-securise-32-caracteres

# Agent Authentication
AGENT_SECRET=votre-secret-agent-securise

# Database Configuration
DATABASE_URL=./data/autostrike.db

# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8443

# TLS Configuration (production)
TLS_CERT_PATH=/path/to/cert.pem
TLS_KEY_PATH=/path/to/key.pem

# Log Level (debug, info, warn, error)
LOG_LEVEL=info
```

#### 3. Backend (Go)

```bash
# Aller dans le dossier serveur
cd server

# Télécharger les dépendances
go mod download

# Compiler le serveur
go build -o autostrike-server ./cmd/autostrike

# Lancer le serveur
./autostrike-server
```

Le serveur démarre sur `https://localhost:8443` par défaut.

#### 4. Frontend (React)

```bash
# Aller dans le dossier dashboard
cd dashboard

# Installer les dépendances
npm install

# Construire le dashboard
npm run build

# Lancer en mode preview
npm run preview
```

Le dashboard est accessible sur `http://localhost:4173`.

### Option 2 : Installation avec Docker

#### Pré-requis Docker

Assurez-vous que Docker et Docker Compose sont installés :

```bash
docker --version
docker compose version
```

#### Lancement avec Docker Compose

```bash
# Construire et lancer les conteneurs
docker compose up -d

# Voir les logs
docker compose logs -f

# Arrêter les conteneurs
docker compose down
```

Les services sont disponibles sur :
- API : `https://localhost:8443`
- Dashboard : `http://localhost:3000`

### Vérification de l'installation

Testez que le serveur répond :

```bash
# Test de santé du serveur
curl -k https://localhost:8443/health

# Réponse attendue
{"status": "healthy"}
```

---

## Installation de l'agent

### Compilation depuis les sources

#### Prérequis

- Rust 1.75+ avec Cargo
- OpenSSL development libraries (Linux)

```bash
# Installer Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# Aller dans le dossier agent
cd agent

# Compiler en mode release
cargo build --release

# L'exécutable se trouve dans target/release/
```

### Windows (PowerShell Admin)

```powershell
# Télécharger l'agent
Invoke-WebRequest -Uri "https://server:8443/deploy/agent.exe" -OutFile "autostrike-agent.exe"

# Vérifier le hash (optionnel mais recommandé)
Get-FileHash -Path "autostrike-agent.exe" -Algorithm SHA256

# Lancer l'agent
.\autostrike-agent.exe --server https://server:8443 --paw "agent-001"
```

#### Installation en tant que service Windows

```powershell
# Créer le service
New-Service -Name "AutoStrike-Agent" -BinaryPathName "C:\Path\To\autostrike-agent.exe --server https://server:8443" -DisplayName "AutoStrike Agent" -StartupType Automatic

# Démarrer le service
Start-Service -Name "AutoStrike-Agent"

# Vérifier le statut
Get-Service -Name "AutoStrike-Agent"
```

### Linux

```bash
# Télécharger l'agent
curl -o autostrike-agent https://server:8443/deploy/agent
chmod +x autostrike-agent

# Vérifier le hash (optionnel mais recommandé)
sha256sum autostrike-agent

# Lancer l'agent
sudo ./autostrike-agent --server https://server:8443 --paw "agent-001"
```

#### Installation en tant que service systemd

Créez le fichier `/etc/systemd/system/autostrike-agent.service` :

```ini
[Unit]
Description=AutoStrike Agent
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/autostrike
ExecStart=/opt/autostrike/autostrike-agent --server https://server:8443
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Puis activez et démarrez le service :

```bash
# Recharger systemd
sudo systemctl daemon-reload

# Activer le service au démarrage
sudo systemctl enable autostrike-agent

# Démarrer le service
sudo systemctl start autostrike-agent

# Vérifier le statut
sudo systemctl status autostrike-agent
```

### macOS

```bash
# Télécharger l'agent
curl -o autostrike-agent https://server:8443/deploy/agent-darwin
chmod +x autostrike-agent

# Lancer l'agent
sudo ./autostrike-agent --server https://server:8443 --paw "agent-001"
```

### Options de l'agent

| Option | Description | Défaut |
|--------|-------------|--------|
| `--server` | URL du serveur AutoStrike | Requis |
| `--paw` | Identifiant unique de l'agent | Généré automatiquement |
| `--interval` | Intervalle de heartbeat (secondes) | 30 |
| `--executors` | Exécuteurs à activer | Tous disponibles |

---

## Mise à jour

### Serveur

```bash
# Arrêter le serveur
docker compose down

# Télécharger les mises à jour
git pull origin main

# Reconstruire et relancer
docker compose up -d --build
```

### Agent

```bash
# Télécharger la nouvelle version
curl -o autostrike-agent-new https://server:8443/deploy/agent

# Remplacer l'ancien exécutable
sudo systemctl stop autostrike-agent
sudo mv autostrike-agent-new /opt/autostrike/autostrike-agent
chmod +x /opt/autostrike/autostrike-agent
sudo systemctl start autostrike-agent
```

---

## Dépannage

### Le serveur ne démarre pas

1. Vérifiez les variables d'environnement :
   ```bash
   cat .env
   ```

2. Vérifiez les permissions sur le fichier de base de données :
   ```bash
   ls -la data/
   ```

3. Consultez les logs :
   ```bash
   docker compose logs server
   ```

### L'agent ne se connecte pas

1. Vérifiez la connectivité réseau :
   ```bash
   curl -k https://server:8443/health
   ```

2. Vérifiez le certificat TLS :
   ```bash
   openssl s_client -connect server:8443
   ```

3. Vérifiez les logs de l'agent :
   ```bash
   journalctl -u autostrike-agent -f
   ```

### Erreur "connection refused"

- Assurez-vous que le serveur écoute sur la bonne interface
- Vérifiez les règles de firewall
- Confirmez que le port 8443 est ouvert

---

## Configuration avancée

Consultez les sections suivantes pour plus de détails :

- [Architecture Backend](../architecture/backend.md)
- [Guide de Déploiement](./deployment.md)
- [Référence API](../api/reference.md)
