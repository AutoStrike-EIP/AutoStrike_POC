# MITRE ATT&CK Techniques

AutoStrike implements 48 MITRE ATT&CK techniques across 13 tactics.

---

## Overview

| Tactic | Count | Description |
|--------|-------|-------------|
| Reconnaissance | 2 | Information gathering about the target |
| Initial Access | 3 | Gaining initial access to the network |
| Execution | 5 | Running commands via shells and interpreters |
| Persistence | 4 | Maintaining access to the system |
| Privilege Escalation | 4 | Gaining higher permissions |
| Defense Evasion | 6 | Avoiding detection by security tools |
| Credential Access | 4 | Stealing credentials and secrets |
| Discovery | 9 | Information gathering about the system and network |
| Lateral Movement | 3 | Moving through the network |
| Collection | 4 | Gathering target data |
| Command and Control | 3 | Communicating with compromised systems |
| Exfiltration | 3 | Stealing data from the target |
| Impact | 3 | Manipulating or destroying systems |

**All techniques are Safe Mode compatible** (non-destructive).

---

## Reconnaissance (2 techniques)

### T1592.004 - Gather Victim Host Information: Client Configurations

Enumerates client software configurations.

=== "Windows"
    ```cmd
    systeminfo | findstr /B /C:"OS Name" /C:"OS Version" /C:"System Type"
    ```

=== "Linux"
    ```bash
    cat /etc/os-release && uname -m
    ```

**Platforms:** Windows, Linux
**Safe:** Yes
**Detection:** Process Creation (systeminfo.exe, cat)

---

### T1595.002 - Active Scanning: Vulnerability Scanning

Simulates vulnerability scanning by checking for common service ports.

=== "Windows"
    ```cmd
    netstat -an | findstr "LISTENING"
    ```

=== "Linux"
    ```bash
    ss -tlnp 2>/dev/null || netstat -tlnp 2>/dev/null
    ```

**Platforms:** Windows, Linux
**Safe:** Yes
**Detection:** Network Connection, Process Creation

---

## Initial Access (3 techniques)

### T1078 - Valid Accounts

Enumerates valid accounts that could be used for initial access.

=== "Windows"
    ```cmd
    net user && net localgroup administrators
    ```

=== "Linux"
    ```bash
    cat /etc/passwd | grep -v nologin | grep -v false
    ```

**Platforms:** Windows, Linux
**Safe:** Yes (enumeration only)
**Detection:** Process Creation (net.exe, cat)

---

### T1133 - External Remote Services

Checks for external remote services (RDP, SSH, VNC).

=== "Windows"
    ```cmd
    netstat -an | findstr ":3389 :22 :5900"
    ```

=== "Linux"
    ```bash
    ss -tlnp | grep -E ":(22|3389|5900)\s"
    ```

**Platforms:** Windows, Linux
**Safe:** Yes (enumeration only)
**Detection:** Process Creation, Network Connection

---

### T1190 - Exploit Public-Facing Application

Enumerates public-facing services and their versions.

=== "Windows"
    ```cmd
    netstat -ano | findstr "LISTENING"
    ```

=== "Linux"
    ```bash
    ss -tlnp
    ```

**Platforms:** Windows, Linux
**Safe:** Yes (enumeration only)
**Detection:** Process Creation

---

## Execution (5 techniques)

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

### T1047 - Windows Management Instrumentation

Enumerates system information via WMI.

```cmd
wmic os get caption,version,buildnumber /format:list
```

**Platforms:** Windows
**Safe:** Yes (read-only WMI query)
**Detection:** Process Creation (wmic.exe), WMI Activity

---

### T1059.006 - Python

Executes Python commands for system enumeration.

=== "Windows"
    ```cmd
    python -c "import platform; print(platform.platform())"
    ```

=== "Linux"
    ```bash
    python3 -c "import platform; print(platform.platform())"
    ```

**Platforms:** Windows, Linux
**Safe:** Yes
**Detection:** Process Creation (python.exe, python3)

---

## Persistence (4 techniques)

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

### T1053.003 - Cron

Enumerates cron jobs for persistence detection.

```bash
crontab -l 2>/dev/null; ls -la /etc/cron* 2>/dev/null
```

**Platforms:** Linux
**Safe:** Yes (read-only)
**Detection:** Process Creation (crontab), File Access

---

### T1543.002 - Systemd Service

Enumerates systemd services for persistence detection.

```bash
systemctl list-unit-files --type=service --state=enabled
```

**Platforms:** Linux
**Safe:** Yes (read-only)
**Detection:** Process Creation (systemctl)

---

## Privilege Escalation (4 techniques)

### T1548.001 - Setuid and Setgid

Searches for SUID/SGID binaries that could be exploited.

```bash
find / -perm -4000 -type f 2>/dev/null | head -20
```

**Platforms:** Linux
**Safe:** Yes (read-only search)
**Detection:** Process Creation (find), File Access

---

### T1548.002 - Bypass User Account Control

Checks UAC configuration status.

```cmd
reg query HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System /v EnableLUA
```

**Platforms:** Windows
**Safe:** Yes (read-only)
**Detection:** Process Creation (reg.exe), Registry Access

---

### T1078.003 - Local Accounts

Enumerates local accounts with elevated privileges.

=== "Windows"
    ```cmd
    net localgroup administrators
    ```

=== "Linux"
    ```bash
    grep -E "^(root|sudo)" /etc/group
    ```

**Platforms:** Windows, Linux
**Safe:** Yes (enumeration only)
**Detection:** Process Creation (net.exe, grep)

---

### T1134.001 - Token Impersonation/Theft

Enumerates process tokens and privileges.

```cmd
whoami /priv
```

**Platforms:** Windows
**Safe:** Yes (read-only)
**Detection:** Process Creation (whoami.exe)

---

## Defense Evasion (6 techniques)

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

### T1562.001 - Disable or Modify Tools

Checks security tool status without disabling them.

=== "Windows"
    ```cmd
    sc query WinDefend
    ```

=== "Linux"
    ```bash
    systemctl status apparmor 2>/dev/null || systemctl status selinux 2>/dev/null
    ```

**Platforms:** Windows, Linux
**Safe:** Yes (query only)
**Detection:** Process Creation, Service Query

---

### T1027 - Obfuscated Files or Information

Demonstrates obfuscation detection by encoding/decoding a test string.

=== "Windows"
    ```powershell
    [Convert]::ToBase64String([Text.Encoding]::UTF8.GetBytes("AutoStrike-Test"))
    ```

=== "Linux"
    ```bash
    echo "AutoStrike-Test" | base64
    ```

**Platforms:** Windows, Linux
**Safe:** Yes (harmless encoding)
**Detection:** Script Block Logging, Process Creation

---

### T1070.001 - Clear Windows Event Logs

Queries event log sizes without clearing them.

```cmd
wevtutil gli Security
```

**Platforms:** Windows
**Safe:** Yes (read-only query)
**Detection:** Process Creation (wevtutil.exe)

---

### T1036.005 - Match Legitimate Name or Location

Checks for processes running from unusual locations.

=== "Windows"
    ```cmd
    wmic process get name,executablepath /format:csv
    ```

=== "Linux"
    ```bash
    ps aux --sort=-%mem | head -20
    ```

**Platforms:** Windows, Linux
**Safe:** Yes (read-only)
**Detection:** Process Creation

---

### T1218.011 - Rundll32

Enumerates Rundll32 usage patterns.

```cmd
tasklist /FI "IMAGENAME eq rundll32.exe"
```

**Platforms:** Windows
**Safe:** Yes (query only)
**Detection:** Process Creation (tasklist.exe)

---

## Credential Access (4 techniques)

### T1552.001 - Credentials In Files

Searches for files that may contain credentials.

=== "Windows"
    ```cmd
    dir /s /b C:\Users\*.txt C:\Users\*.ini C:\Users\*.cfg 2>nul | findstr /i "pass config cred"
    ```

=== "Linux"
    ```bash
    find /home -name "*.conf" -o -name "*.cfg" -o -name "*.ini" 2>/dev/null | head -20
    ```

**Platforms:** Windows, Linux
**Safe:** Yes (file search only, no reading)
**Detection:** Process Creation, File Access

---

### T1555.003 - Credentials from Web Browsers

Checks for browser credential storage locations.

=== "Windows"
    ```cmd
    dir "%LOCALAPPDATA%\Google\Chrome\User Data\Default\Login Data" 2>nul
    ```

=== "Linux"
    ```bash
    ls -la ~/.config/google-chrome/Default/Login\ Data 2>/dev/null
    ```

**Platforms:** Windows, Linux
**Safe:** Yes (existence check only)
**Detection:** File Access

---

### T1003.008 - /etc/passwd and /etc/shadow

Reads password file information (non-sensitive parts).

```bash
cat /etc/passwd | cut -d: -f1,3,6,7
```

**Platforms:** Linux
**Safe:** Yes (reads public passwd info only)
**Detection:** Process Creation (cat), File Access

---

### T1552.004 - Private Keys

Searches for private key files.

=== "Windows"
    ```cmd
    dir /s /b C:\Users\*.pem C:\Users\*.key C:\Users\*.ppk 2>nul
    ```

=== "Linux"
    ```bash
    find /home -name "*.pem" -o -name "*.key" -o -name "id_rsa" 2>/dev/null | head -20
    ```

**Platforms:** Windows, Linux
**Safe:** Yes (file search only, no reading)
**Detection:** Process Creation, File Access

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

## Lateral Movement (3 techniques)

### T1021.001 - Remote Desktop Protocol

Checks for RDP availability and configuration.

```cmd
reg query "HKLM\SYSTEM\CurrentControlSet\Control\Terminal Server" /v fDenyTSConnections
```

**Platforms:** Windows
**Safe:** Yes (read-only query)
**Detection:** Process Creation (reg.exe), Registry Access

---

### T1021.002 - SMB/Windows Admin Shares

Enumerates SMB shares.

=== "Windows"
    ```cmd
    net share
    ```

=== "Linux"
    ```bash
    smbclient -L localhost -N 2>/dev/null || echo "SMB not available"
    ```

**Platforms:** Windows, Linux
**Safe:** Yes (enumeration only)
**Detection:** Process Creation, Network Connection

---

### T1021.004 - SSH

Checks for SSH availability and configuration.

=== "Windows"
    ```cmd
    sc query sshd 2>nul || echo SSH service not found
    ```

=== "Linux"
    ```bash
    systemctl status sshd 2>/dev/null || systemctl status ssh 2>/dev/null
    ```

**Platforms:** Windows, Linux
**Safe:** Yes (status check only)
**Detection:** Process Creation, Service Query

---

## Collection (4 techniques)

### T1005 - Data from Local System

Enumerates files of interest on the local system.

=== "Windows"
    ```cmd
    dir /s /b C:\Users\*.docx C:\Users\*.xlsx C:\Users\*.pdf 2>nul | findstr /v "AppData"
    ```

=== "Linux"
    ```bash
    find /home -name "*.pdf" -o -name "*.docx" -o -name "*.xlsx" 2>/dev/null | head -20
    ```

**Platforms:** Windows, Linux
**Safe:** Yes (file listing only)
**Detection:** Process Creation, File Access

---

### T1039 - Data from Network Shared Drive

Enumerates network shared drives.

=== "Windows"
    ```cmd
    net use
    ```

=== "Linux"
    ```bash
    mount | grep -E "(cifs|nfs|smb)"
    ```

**Platforms:** Windows, Linux
**Safe:** Yes (enumeration only)
**Detection:** Process Creation

---

### T1074.001 - Local Data Staging

Checks for common staging directories.

=== "Windows"
    ```cmd
    dir %TEMP% /O-D /T:W
    ```

=== "Linux"
    ```bash
    ls -lt /tmp/ | head -20
    ```

**Platforms:** Windows, Linux
**Safe:** Yes (read-only)
**Detection:** Process Creation, File Access

---

### T1119 - Automated Collection

Simulates automated file collection by listing recent files.

=== "Windows"
    ```cmd
    forfiles /P C:\Users /S /D -7 /C "cmd /c echo @path @fdate" 2>nul | findstr /i ".doc .xls .pdf"
    ```

=== "Linux"
    ```bash
    find /home -mtime -7 -name "*.pdf" -o -name "*.doc" -o -name "*.xls" 2>/dev/null | head -20
    ```

**Platforms:** Windows, Linux
**Safe:** Yes (file listing only)
**Detection:** Process Creation, File Access

---

## Command and Control (3 techniques)

### T1071.001 - Web Protocols

Tests HTTP/HTTPS connectivity to external services.

=== "Windows"
    ```powershell
    (Invoke-WebRequest -Uri "https://www.google.com" -UseBasicParsing -TimeoutSec 5).StatusCode
    ```

=== "Linux"
    ```bash
    curl -s -o /dev/null -w "%{http_code}" https://www.google.com --connect-timeout 5
    ```

**Platforms:** Windows, Linux
**Safe:** Yes (connectivity check only)
**Detection:** Network Connection, DNS Resolution

---

### T1105 - Ingress Tool Transfer

Checks download capabilities without actually downloading malicious content.

=== "Windows"
    ```powershell
    [Net.ServicePointManager]::SecurityProtocol; Write-Output "Download capabilities available"
    ```

=== "Linux"
    ```bash
    which wget curl 2>/dev/null && echo "Download tools available"
    ```

**Platforms:** Windows, Linux
**Safe:** Yes (capability check only)
**Detection:** Process Creation

---

### T1572 - Protocol Tunneling

Checks for tunneling tools availability.

=== "Windows"
    ```cmd
    netsh interface portproxy show all
    ```

=== "Linux"
    ```bash
    which ssh socat ncat 2>/dev/null; ss -tlnp | grep -E ":(1080|8080|3128)"
    ```

**Platforms:** Windows, Linux
**Safe:** Yes (enumeration only)
**Detection:** Process Creation, Network Connection

---

## Exfiltration (3 techniques)

### T1048.003 - Exfiltration Over Unencrypted Non-C2 Protocol

Checks for common exfiltration channels (DNS, HTTP).

=== "Windows"
    ```cmd
    nslookup test.example.com 2>nul && echo "DNS resolution available"
    ```

=== "Linux"
    ```bash
    dig +short test.example.com 2>/dev/null || nslookup test.example.com 2>/dev/null
    ```

**Platforms:** Windows, Linux
**Safe:** Yes (DNS query only)
**Detection:** DNS Resolution, Network Connection

---

### T1041 - Exfiltration Over C2 Channel

Checks outbound connectivity that could be used for exfiltration.

=== "Windows"
    ```powershell
    Test-NetConnection -ComputerName "8.8.8.8" -Port 443 -InformationLevel Quiet
    ```

=== "Linux"
    ```bash
    timeout 5 bash -c "echo >/dev/tcp/8.8.8.8/443" 2>/dev/null && echo "Port 443 reachable"
    ```

**Platforms:** Windows, Linux
**Safe:** Yes (connectivity test only)
**Detection:** Network Connection

---

### T1567.002 - Exfiltration to Cloud Storage

Checks for cloud storage CLI tools that could be used for exfiltration.

=== "Windows"
    ```cmd
    where rclone aws gsutil 2>nul
    ```

=== "Linux"
    ```bash
    which rclone aws gsutil 2>/dev/null
    ```

**Platforms:** Windows, Linux
**Safe:** Yes (tool detection only)
**Detection:** Process Creation

---

## Impact (3 techniques)

### T1490 - Inhibit System Recovery

Checks system recovery configuration status.

=== "Windows"
    ```cmd
    vssadmin list shadows
    ```

=== "Linux"
    ```bash
    ls -la /boot/grub/ 2>/dev/null; cat /proc/cmdline
    ```

**Platforms:** Windows, Linux
**Safe:** Yes (read-only query)
**Detection:** Process Creation (vssadmin.exe)

---

### T1489 - Service Stop

Lists critical services without stopping them.

=== "Windows"
    ```cmd
    sc query type= service state= all | findstr /i "SERVICE_NAME DISPLAY_NAME STATE"
    ```

=== "Linux"
    ```bash
    systemctl list-units --type=service --state=running --no-pager
    ```

**Platforms:** Windows, Linux
**Safe:** Yes (listing only, no services stopped)
**Detection:** Process Creation, Service Query

---

### T1486 - Data Encrypted for Impact

Simulates ransomware detection by checking encryption capabilities.

=== "Windows"
    ```powershell
    Get-Command -Name *crypt* -ErrorAction SilentlyContinue | Select-Object Name
    ```

=== "Linux"
    ```bash
    which openssl gpg 2>/dev/null && echo "Encryption tools available"
    ```

**Platforms:** Windows, Linux
**Safe:** Yes (capability check only, no encryption performed)
**Detection:** Process Creation

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
# Via API (YAML file on server)
curl -X POST https://localhost:8443/api/v1/techniques/import \
  -H "Content-Type: application/json" \
  -d '{"path": "/path/to/techniques.yaml"}'

# Via API (JSON array directly)
curl -X POST https://localhost:8443/api/v1/techniques/import/json \
  -H "Content-Type: application/json" \
  -d '[{"id": "T1234", "name": "Custom", ...}]'
```

Techniques in `server/configs/techniques/` are auto-imported at server startup.

---

## Technique Locations

| File | Tactic | Count |
|------|--------|-------|
| `reconnaissance.yaml` | Reconnaissance | 2 |
| `initial-access.yaml` | Initial Access | 3 |
| `execution.yaml` | Execution | 5 |
| `persistence.yaml` | Persistence | 4 |
| `privilege-escalation.yaml` | Privilege Escalation | 4 |
| `defense-evasion.yaml` | Defense Evasion | 6 |
| `credential-access.yaml` | Credential Access | 4 |
| `discovery.yaml` | Discovery | 9 |
| `lateral-movement.yaml` | Lateral Movement | 3 |
| `collection.yaml` | Collection | 4 |
| `command-and-control.yaml` | Command and Control | 3 |
| `exfiltration.yaml` | Exfiltration | 3 |
| `impact.yaml` | Impact | 3 |
