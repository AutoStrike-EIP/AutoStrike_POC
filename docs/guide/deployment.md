# Déploiement

Guide de déploiement d'AutoStrike en production.

---

## Prérequis

- Docker et Docker Compose
- Certificats TLS (ou Let's Encrypt)
- Serveur Linux (Ubuntu 22.04+ recommandé)

---

## Configuration

### 1. Variables d'environnement

Créer un fichier `.env` à partir du template :

```bash
cp .env.example .env
```

Générer des secrets sécurisés :

```bash
# Générer JWT_SECRET
openssl rand -base64 32

# Générer AGENT_SECRET
openssl rand -base64 32
```

Éditer `.env` avec les valeurs générées :

```env
JWT_SECRET=<votre-jwt-secret-généré>
AGENT_SECRET=<votre-agent-secret-généré>
```

### 2. Certificats TLS

#### Option A : Certificats auto-signés (test)

```bash
mkdir -p certs
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout certs/server.key \
  -out certs/server.crt \
  -subj "/CN=autostrike.local"
```

#### Option B : Let's Encrypt (production)

```bash
# Installer certbot
apt install certbot

# Obtenir un certificat
certbot certonly --standalone -d autostrike.example.com

# Copier les certificats
cp /etc/letsencrypt/live/autostrike.example.com/fullchain.pem certs/server.crt
cp /etc/letsencrypt/live/autostrike.example.com/privkey.pem certs/server.key
```

---

## Déploiement Docker

### Lancement

```bash
# Construire et démarrer les services
docker compose up -d

# Vérifier l'état
docker compose ps

# Voir les logs
docker compose logs -f server
```

### Services

| Service | Port | Description |
|---------|------|-------------|
| server | 8443 | API REST + WebSocket |
| dashboard | 3000 | Interface web |

### Vérification

```bash
# Tester la connexion HTTPS
curl -k https://localhost:8443/health

# Tester l'API (nécessite un token JWT)
curl -k -H "Authorization: Bearer <token>" https://localhost:8443/api/v1/agents
```

---

## Déploiement des agents

### Windows

```powershell
# Télécharger l'agent
Invoke-WebRequest -Uri "https://server:8443/deploy/agent.exe" -OutFile "autostrike-agent.exe"

# Lancer avec le serveur configuré
.\autostrike-agent.exe --server https://server:8443 --paw agent-win-01
```

### Linux

```bash
# Télécharger l'agent
curl -k -o autostrike-agent https://server:8443/deploy/agent
chmod +x autostrike-agent

# Lancer l'agent
./autostrike-agent --server https://server:8443 --paw agent-linux-01
```

### Agent comme service (systemd)

```bash
# Créer le fichier service
cat > /etc/systemd/system/autostrike-agent.service << EOF
[Unit]
Description=AutoStrike BAS Agent
After=network.target

[Service]
Type=simple
ExecStart=/opt/autostrike/agent --server https://server:8443 --paw $(hostname)
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

# Activer et démarrer
systemctl enable autostrike-agent
systemctl start autostrike-agent
```

---

## Maintenance

### Sauvegarde

```bash
# Sauvegarder la base de données
docker compose exec server sqlite3 /app/data/autostrike.db ".backup /app/data/backup.db"
docker cp autostrike-server:/app/data/backup.db ./backup-$(date +%Y%m%d).db
```

### Mise à jour

```bash
# Arrêter les services
docker compose down

# Mettre à jour le code
git pull

# Reconstruire et redémarrer
docker compose up -d --build
```

### Logs

```bash
# Logs du serveur
docker compose logs -f server

# Logs du dashboard
docker compose logs -f dashboard
```

---

## Sécurité

### Recommandations

1. **Secrets** : Utiliser des secrets de 32+ caractères générés aléatoirement
2. **TLS** : Toujours utiliser HTTPS en production
3. **Firewall** : Restreindre l'accès aux ports 8443 et 3000
4. **Réseau** : Isoler les agents dans un réseau dédié
5. **Logs** : Activer la journalisation centralisée

### Ports à ouvrir

| Port | Direction | Description |
|------|-----------|-------------|
| 8443 | Entrant | API + WebSocket agents |
| 3000 | Entrant | Dashboard web |

---

## Troubleshooting

### L'agent ne se connecte pas

1. Vérifier la connectivité réseau vers le serveur
2. Vérifier les certificats TLS (utiliser `--insecure` pour les tests)
3. Vérifier les logs : `./agent --debug`

### Erreur 401 Unauthorized

1. Vérifier que le token JWT n'est pas expiré
2. Vérifier que `JWT_SECRET` est identique entre génération et serveur

### Base de données corrompue

```bash
# Sauvegarder l'état actuel
docker cp autostrike-server:/app/data/autostrike.db ./autostrike-corrupted.db

# Supprimer et recréer
docker compose down -v
docker compose up -d
```
