# MITRE ATT&CK Techniques

AutoStrike implements 294 MITRE ATT&CK techniques across 12 tactics, all imported from MITRE ATT&CK STIX + Atomic Red Team via `make import-mitre`.

---

## Overview

| Tactic | Count | Description |
|--------|-------|-------------|
| Initial Access | 4 | Gaining initial access to the network |
| Execution | 22 | Running commands via shells and interpreters |
| Persistence | 44 | Maintaining access to the system |
| Privilege Escalation | 18 | Gaining higher permissions |
| Defense Evasion | 89 | Avoiding detection by security tools |
| Credential Access | 34 | Stealing credentials and secrets |
| Discovery | 30 | Information gathering about the system and network |
| Lateral Movement | 8 | Moving through the network |
| Collection | 16 | Gathering target data |
| Command and Control | 13 | Communicating with compromised systems |
| Exfiltration | 8 | Stealing data from the target |
| Impact | 8 | Manipulating or destroying systems |

**Safety classification:** 220 safe, 74 unsafe. Per-executor safety based on elevation requirements and dangerous command pattern detection.

---

## Safety System

Safety is determined **per-executor** using the formula:

```
is_safe = !elevation_required && !hasDangerousPattern(command) && !hasDangerousPattern(cleanup)
```

### Dangerous Patterns

Commands containing any of the following are automatically marked unsafe:

| Pattern | Description |
|---------|-------------|
| `rm -rf`, `rm -fr` | Recursive forced deletion (Linux) |
| `del /f`, `rd /s`, `rmdir /s` | Forced deletion (Windows) |
| `dd of=/dev/` | Direct disk write |
| `mkfs`, `fdisk`, `format C:` | Filesystem/partition modification |
| `shutdown`, `reboot`, `init 0` | System shutdown/restart |
| `taskkill /f`, `kill -9`, `killall`, `pkill` | Force process termination |
| `systemctl stop/disable` | Service disruption |
| `chmod 000` | Permission removal |
| `iptables -F` | Firewall rule flush |

### Safe Mode

When safe mode is enabled, only executors with `is_safe: true` are used. Use `make import-mitre-safe` to import only techniques that have at least one safe executor.

---

## Adding Custom Techniques

Techniques are defined in YAML files under `server/configs/techniques/`:

```yaml
# server/configs/techniques/custom.yaml
techniques:
  - id: T1234
    name: Custom Technique
    description: Description of the technique
    tactic: discovery
    tactics:                     # Optional: all tactics (multi-tactic)
      - discovery
      - collection
    platforms:
      - windows
      - linux
    is_safe: true
    references:                  # Optional: MITRE ATT&CK URLs
      - https://attack.mitre.org/techniques/T1234
    executors:
      - name: "System info via cmd"  # Optional: executor name
        type: cmd
        platform: windows
        command: "echo Hello"
        cleanup: ""
        timeout: 60
        elevation_required: false    # Optional: needs admin/root
      - name: "System info via bash"
        type: bash
        platform: linux
        command: "echo Hello"
        cleanup: ""
        timeout: 60
    detection:
      - source: Process Creation
        indicator: echo execution
```

### YAML Structure

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | MITRE ATT&CK ID (e.g., T1082) |
| `name` | string | Technique name |
| `description` | string | Detailed description |
| `tactic` | string | Primary MITRE tactic (discovery, execution, etc.) |
| `tactics` | array (optional) | All MITRE tactics for multi-tactic techniques |
| `platforms` | array | Supported platforms (windows, linux, macos) |
| `is_safe` | boolean | Safe for production testing |
| `references` | array (optional) | MITRE ATT&CK reference URLs |
| `executors` | array | Command definitions per platform |
| `detection` | array | Expected detection indicators |

### Executor Fields

| Field | Type | Description |
|-------|------|-------------|
| `name` | string (optional) | Executor display name (distinguishes multiple executors) |
| `type` | string | Executor type: `cmd`, `powershell`, `bash`, `sh` |
| `platform` | string (optional) | Target platform: `windows`, `linux`, `macos` |
| `command` | string | Command to execute |
| `cleanup` | string (optional) | Cleanup command run after execution |
| `timeout` | int | Execution timeout in seconds |
| `elevation_required` | bool (optional) | Whether root/admin privileges are needed |
| `is_safe` | bool (optional) | Per-executor safety (auto-computed: `!elevation_required && no dangerous patterns in command/cleanup`) |

### Import Techniques

```bash
# Automatic import from MITRE ATT&CK + Atomic Red Team
make import-mitre           # Full import
make import-mitre-safe      # Safe techniques only
make import-mitre-dry       # Dry run (stats only)

# Via API (YAML file on server)
curl -X POST https://localhost:8443/api/v1/techniques/import \
  -H "Content-Type: application/json" \
  -d '{"path": "/path/to/techniques.yaml"}'

# List executors for a technique (with optional platform filter)
curl https://localhost:8443/api/v1/techniques/T1082/executors?platform=linux
```

Techniques in `server/configs/techniques/` are auto-imported at server startup (all `*.yaml` and `*.yml` files are loaded dynamically via `os.ReadDir`).

---

## Technique Locations

| File | Tactic | Count |
|------|--------|-------|
| `initial-access.yaml` | Initial Access | 4 |
| `execution.yaml` | Execution | 22 |
| `persistence.yaml` | Persistence | 44 |
| `privilege-escalation.yaml` | Privilege Escalation | 18 |
| `defense-evasion.yaml` | Defense Evasion | 89 |
| `credential-access.yaml` | Credential Access | 34 |
| `discovery.yaml` | Discovery | 30 |
| `lateral-movement.yaml` | Lateral Movement | 8 |
| `collection.yaml` | Collection | 16 |
| `command-and-control.yaml` | Command and Control | 13 |
| `exfiltration.yaml` | Exfiltration | 8 |
| `impact.yaml` | Impact | 8 |
