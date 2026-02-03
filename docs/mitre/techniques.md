# MITRE ATT&CK Techniques

AutoStrike implements 15 MITRE ATT&CK techniques across 4 tactics.

---

## Overview

| Tactic | Count | Description |
|--------|-------|-------------|
| Discovery | 9 | Information gathering about the system and network |
| Execution | 3 | Running commands via shells |
| Persistence | 2 | Maintaining access |
| Defense Evasion | 1 | Avoiding detection |

**All techniques are Safe Mode compatible** (non-destructive).

---

## Discovery (9 techniques)

### T1082 - System Information Discovery

Collects information about the operating system and hardware.

=== "Windows"
    ```cmd
    systeminfo
    ```

=== "Linux"
    ```bash
    uname -a && cat /etc/os-release
    ```

**Platforms:** Windows, Linux
**Safe:** Yes
**Detection:** Process Creation (systeminfo.exe, uname)

---

### T1083 - File and Directory Discovery

Enumerates files and directories.

=== "Windows"
    ```cmd
    dir C:\Users
    ```

=== "Linux"
    ```bash
    ls -la /home
    ```

**Platforms:** Windows, Linux
**Safe:** Yes
**Detection:** Process Creation (dir.exe, ls)

---

### T1057 - Process Discovery

Lists running processes.

=== "Windows"
    ```cmd
    tasklist
    ```

=== "Linux"
    ```bash
    ps aux
    ```

**Platforms:** Windows, Linux
**Safe:** Yes
**Detection:** Process Creation (tasklist.exe, ps)

---

### T1016 - System Network Configuration Discovery

Gathers network configuration information.

=== "Windows"
    ```cmd
    ipconfig /all
    ```

=== "Linux"
    ```bash
    ip addr && cat /etc/resolv.conf
    ```

**Platforms:** Windows, Linux
**Safe:** Yes
**Detection:** Process Creation (ipconfig.exe, ip)

---

### T1049 - System Network Connections Discovery

Lists active network connections.

=== "Windows"
    ```cmd
    netstat -ano
    ```

=== "Linux"
    ```bash
    netstat -tunap || ss -tunap
    ```

**Platforms:** Windows, Linux
**Safe:** Yes
**Detection:** Process Creation (netstat.exe, ss)

---

### T1087 - Account Discovery

Enumerates user accounts.

=== "Windows"
    ```cmd
    net user
    ```

=== "Linux"
    ```bash
    cat /etc/passwd | cut -d: -f1
    ```

**Platforms:** Windows, Linux
**Safe:** Yes
**Detection:** Process Creation (net.exe, cat)

---

### T1069 - Permission Groups Discovery

Enumerates permission groups.

=== "Windows"
    ```cmd
    net localgroup
    ```

=== "Linux"
    ```bash
    cat /etc/group
    ```

**Platforms:** Windows, Linux
**Safe:** Yes
**Detection:** Process Creation (net.exe, cat)

---

### T1018 - Remote System Discovery

Discovers remote systems on the network.

=== "Windows"
    ```cmd
    net view
    ```

=== "Linux"
    ```bash
    arp -a
    ```

**Platforms:** Windows, Linux
**Safe:** Yes
**Detection:** Process Creation (net.exe, arp)

---

### T1007 - System Service Discovery

Lists running services.

=== "Windows"
    ```cmd
    sc query
    ```

=== "Linux"
    ```bash
    systemctl list-units --type=service --state=running
    ```

**Platforms:** Windows, Linux
**Safe:** Yes
**Detection:** Process Creation (sc.exe, systemctl)

---

## Execution (3 techniques)

### T1059.001 - PowerShell

Executes PowerShell commands.

```powershell
Get-Process | Select-Object -First 5
```

**Platforms:** Windows
**Safe:** Yes (current command)
**Detection:** Process Creation (powershell.exe), Script Block Logging

---

### T1059.003 - Windows Command Shell

Executes Windows CMD commands.

```cmd
echo %USERNAME% && echo %COMPUTERNAME%
```

**Platforms:** Windows
**Safe:** Yes
**Detection:** Process Creation (cmd.exe)

---

### T1059.004 - Unix Shell

Executes Unix shell commands.

=== "Bash"
    ```bash
    echo $USER && hostname
    ```

=== "sh"
    ```sh
    echo $USER && hostname
    ```

**Platforms:** Linux
**Safe:** Yes
**Detection:** Process Creation (bash, sh)

---

## Persistence (2 techniques)

### T1053.005 - Scheduled Task

Queries scheduled tasks.

```cmd
schtasks /query /fo LIST
```

**Platforms:** Windows
**Safe:** Yes (query only, doesn't create tasks)
**Detection:** Process Creation (schtasks.exe)

---

### T1547.001 - Registry Run Keys / Startup Folder

Queries registry Run keys for persistence mechanisms.

=== "CMD"
    ```cmd
    reg query HKCU\Software\Microsoft\Windows\CurrentVersion\Run
    ```

=== "PowerShell"
    ```powershell
    Get-ItemProperty -Path 'HKCU:\Software\Microsoft\Windows\CurrentVersion\Run'
    ```

**Platforms:** Windows
**Safe:** Yes (query only, doesn't modify registry)
**Detection:** Process Creation (reg.exe), Registry access

---

## Defense Evasion (1 technique)

### T1070.004 - File Deletion

Simulates file deletion (safe mode only echoes).

=== "Windows"
    ```cmd
    echo Test file deletion simulation
    ```

=== "Linux"
    ```bash
    echo Test file deletion simulation
    ```

**Platforms:** Windows, Linux
**Safe:** Yes (simulation only)
**Detection:** File System events, Process Creation

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
    platforms:
      - windows
      - linux
    is_safe: true
    executors:
      - type: cmd
        platform: windows
        command: "echo Hello"
        cleanup: ""
        timeout: 60
      - type: bash
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
| `tactic` | string | MITRE tactic (discovery, execution, etc.) |
| `platforms` | array | Supported platforms (windows, linux, darwin) |
| `is_safe` | boolean | Safe for production testing |
| `executors` | array | Command definitions per platform |
| `detection` | array | Expected detection indicators |

### Import Techniques

```bash
# Via API
curl -X POST https://localhost:8443/api/v1/techniques/import \
  -H "Content-Type: application/json" \
  -d '{"path": "/path/to/techniques.yaml"}'
```

Techniques in `server/configs/techniques/` are auto-imported at server startup.

---

## Technique Locations

| File | Tactic | Count |
|------|--------|-------|
| `discovery.yaml` | Discovery | 9 |
| `execution.yaml` | Execution | 3 |
| `persistence.yaml` | Persistence | 2 |
| `defense-evasion.yaml` | Defense Evasion | 1 |
