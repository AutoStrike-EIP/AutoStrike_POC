# Deployment

Production deployment guide for AutoStrike.

---

## Prerequisites

- Docker and Docker Compose (optional)
- TLS certificates (or Let's Encrypt)
- Linux server (Ubuntu 22.04+ recommended)
- Go 1.24+ and Node.js 18+ (for manual deployment)

---

## Configuration

### 1. Environment Variables

Create a `.env` file from the template:

```bash
cd server
cp .env.example .env
```

Generate secure secrets:

```bash
# Generate JWT_SECRET
openssl rand -base64 32

# Generate AGENT_SECRET
openssl rand -base64 32
```

Edit `.env` with the generated values:

```env
# Authentication (required for production)
JWT_SECRET=<your-generated-jwt-secret>
AGENT_SECRET=<your-generated-agent-secret>

# Database
DATABASE_PATH=./data/autostrike.db

# Dashboard (path to built React app)
DASHBOARD_PATH=../dashboard/dist

# CORS (adjust for your domain)
ALLOWED_ORIGINS=https://your-domain.com

# Logging
LOG_LEVEL=info

# SMTP (optional - for email notifications)
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=notifications@example.com
SMTP_PASSWORD=<smtp-password>
SMTP_FROM=noreply@example.com
SMTP_USE_TLS=true
DASHBOARD_URL=https://your-domain.com
```

### 2. TLS Certificates

#### Option A: Self-signed certificates (testing)

```bash
mkdir -p certs
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout certs/server.key \
  -out certs/server.crt \
  -subj "/CN=autostrike.local"
```

#### Option B: Let's Encrypt (production)

```bash
# Install certbot
apt install certbot

# Obtain certificate
certbot certonly --standalone -d autostrike.example.com

# Copy certificates
cp /etc/letsencrypt/live/autostrike.example.com/fullchain.pem certs/server.crt
cp /etc/letsencrypt/live/autostrike.example.com/privkey.pem certs/server.key
```

---

## Deployment Methods

### Option 1: Makefile (Recommended)

```bash
# Build everything
make install

# Start server (background)
make run

# Check logs
make logs

# Stop
make stop
```

### Option 2: Docker Compose

```bash
# Build and start
docker compose up -d

# Check status
docker compose ps

# View logs
docker compose logs -f

# Stop
docker compose down
```

### Option 3: Manual

```bash
# Build backend
cd server
go build -o autostrike ./cmd/autostrike

# Build frontend
cd ../dashboard
npm install && npm run build

# Start server
cd ../server
./autostrike
```

### Verification

```bash
# Test HTTPS connection
curl -k https://localhost:8443/health
# Expected: {"status":"ok","auth_enabled":true}

# Test API (with JWT)
curl -k https://localhost:8443/api/v1/agents \
  -H "Authorization: Bearer <token>"
```

---

## Single Port Architecture

AutoStrike serves everything on **port 8443**:

| Path | Description |
|------|-------------|
| `/` | Dashboard (React SPA) |
| `/api/v1/*` | REST API |
| `/ws/agent` | WebSocket for agents |
| `/ws/dashboard` | WebSocket for real-time updates |
| `/health` | Health check |

No separate ports needed for dashboard.

---

## Agent Deployment

### Build from Source

```bash
cd agent
cargo build --release
# Binary: target/release/autostrike-agent
```

### Windows

```powershell
# Copy binary to target machine
# Run with server URL and agent secret
.\autostrike-agent.exe --server https://server:8443 --paw agent-win-01 -k "your-agent-secret"
```

### Linux

```bash
# Copy binary to target machine
chmod +x autostrike-agent

# Run with server URL and agent secret
./autostrike-agent --server https://server:8443 --paw agent-linux-01 -k "your-agent-secret"
```

### Agent as systemd Service

```bash
# Create service file
cat > /etc/systemd/system/autostrike-agent.service << EOF
[Unit]
Description=AutoStrike BAS Agent
After=network.target

[Service]
Type=simple
ExecStart=/opt/autostrike/autostrike-agent --server https://server:8443 --paw $(hostname)
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

# Enable and start
systemctl daemon-reload
systemctl enable autostrike-agent
systemctl start autostrike-agent
```

---

## Maintenance

### Backup

```bash
# Backup SQLite database
cp server/data/autostrike.db backup-$(date +%Y%m%d).db

# With Docker
docker compose exec server sqlite3 /app/data/autostrike.db ".backup /app/data/backup.db"
docker cp autostrike-server:/app/data/backup.db ./backup-$(date +%Y%m%d).db
```

### Update

```bash
# Stop services
make stop  # or docker compose down

# Pull updates
git pull

# Rebuild and restart
make install
make run  # or docker compose up -d --build
```

### Logs

```bash
# With Makefile
make logs

# With Docker
docker compose logs -f

# Direct
tail -f server/logs/autostrike.log
```

---

## Security Recommendations

### Production Checklist

1. **Secrets**: Use 32+ character randomly generated secrets for `JWT_SECRET` and `AGENT_SECRET`
2. **TLS**: Always use HTTPS in production (Let's Encrypt or commercial cert)
3. **Firewall**: Restrict access to port 8443
4. **Network**: Isolate agents in a dedicated network/VLAN
5. **Logging**: Enable centralized logging
6. **Updates**: Keep Go, Rust, and dependencies updated
7. **Auth**: Ensure `JWT_SECRET` is set (auth enabled)
8. **Agent Auth**: Set `AGENT_SECRET` to authenticate agent connections

### Firewall Rules

```bash
# Allow HTTPS
ufw allow 8443/tcp

# Restrict to specific IPs (optional)
ufw allow from 10.0.0.0/8 to any port 8443
```

### Port Requirements

| Port | Direction | Description |
|------|-----------|-------------|
| 8443 | Inbound | API + WebSocket + Dashboard |

---

## Troubleshooting

### Agent won't connect

1. Verify network connectivity:
   ```bash
   curl -k https://server:8443/health
   ```
2. Check TLS certificates (use `--debug` flag)
3. Check agent logs: `journalctl -u autostrike-agent -f`
4. If `AGENT_SECRET` is set, verify the agent passes `-k` with the correct secret

### Error 401 Unauthorized

1. Verify JWT token is not expired
2. Verify `JWT_SECRET` matches between token generation and server
3. Check if auth is enabled (set `JWT_SECRET` in `.env`)

### Error 429 Too Many Requests

Rate limiting is active on authentication endpoints:
- Login: 5 attempts/minute per IP
- Token refresh: 10 attempts/minute per IP
- Wait for the rate limit window to expire

### Database issues

```bash
# Check database file
ls -la server/data/autostrike.db

# Reset database (WARNING: deletes all data)
rm server/data/autostrike.db
make run  # Will create new database
```

### WebSocket connection issues

1. Check firewall allows WebSocket upgrade
2. Verify no proxy is interfering with WebSocket
3. Check server logs for connection errors
