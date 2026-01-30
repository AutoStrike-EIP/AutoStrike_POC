# Agent (Rust)

L'agent AutoStrike est développé en **Rust** pour des raisons de performance et de sécurité.

---

## Fonctionnalités

- Exécution des techniques MITRE ATT&CK
- Communication sécurisée (mTLS)
- Multi-plateformes (Windows, Linux)
- Mode safe (techniques non destructives)
- Auto-cleanup après exécution

---

## Structure

```
agent/
├── src/
│   ├── main.rs
│   ├── config.rs
│   ├── comms/
│   │   ├── mod.rs
│   │   ├── beacon.rs
│   │   └── tls.rs
│   ├── executor/
│   │   ├── mod.rs
│   │   ├── windows/
│   │   └── linux/
│   └── techniques/
│       ├── mod.rs
│       ├── discovery.rs
│       ├── execution.rs
│       └── persistence.rs
├── Cargo.toml
└── Cargo.lock
```

---

## Déploiement

### Windows

```powershell
.\autostrike-agent.exe --server https://server:8443 --paw AGENT_001
```

### Linux

```bash
./autostrike-agent --server https://server:8443 --paw AGENT_001
```

---

## Options

| Option | Description | Défaut |
|--------|-------------|--------|
| `--server` | URL du serveur | (requis) |
| `--paw` | Identifiant unique | (auto-généré) |
| `--interval` | Intervalle beacon | 30s |
| `--safe-mode` | Mode safe uniquement | true |
