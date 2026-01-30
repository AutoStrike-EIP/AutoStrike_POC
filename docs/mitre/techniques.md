# Techniques

Liste des techniques MITRE ATT&CK implémentées dans AutoStrike.

---

## Discovery

### T1082 - System Information Discovery

Collecte d'informations sur le système.

=== "Windows"

    ```powershell
    systeminfo
    Get-ComputerInfo
    ```

=== "Linux"

    ```bash
    uname -a
    cat /etc/os-release
    ```

**Safe:** Oui

---

### T1083 - File and Directory Discovery

Énumération des fichiers et répertoires.

=== "Windows"

    ```powershell
    dir C:\Users /s
    Get-ChildItem -Recurse
    ```

=== "Linux"

    ```bash
    find / -type f -name "*.conf"
    ls -laR /home
    ```

**Safe:** Oui

---

### T1057 - Process Discovery

Liste des processus en cours.

=== "Windows"

    ```powershell
    tasklist
    Get-Process
    ```

=== "Linux"

    ```bash
    ps aux
    ```

**Safe:** Oui

---

## Execution

### T1059.001 - PowerShell

Exécution de commandes PowerShell.

```powershell
powershell -ExecutionPolicy Bypass -Command "Write-Host 'Test'"
```

**Safe:** Dépend de la commande

---

### T1059.004 - Unix Shell

Exécution de commandes shell Unix.

```bash
bash -c "echo 'Test'"
```

**Safe:** Dépend de la commande

---

## Credential Access

### T1003.001 - LSASS Memory

!!! danger "Non-safe"
    Cette technique peut être détectée et bloquée par les EDR.

Dump de la mémoire LSASS pour récupérer des identifiants.

```powershell
# Nécessite des privilèges élevés
rundll32.exe comsvcs.dll MiniDump
```

**Safe:** Non

---

## Ajouter une technique

Les techniques sont définies en YAML :

```yaml
id: T1082
name: System Information Discovery
tactic: discovery
platforms:
  - windows
  - linux
is_safe: true
commands:
  windows:
    - cmd: systeminfo
      cleanup: null
  linux:
    - cmd: uname -a
      cleanup: null
```
