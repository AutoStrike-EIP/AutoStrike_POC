# MITRE ATT&CK

AutoStrike utilise le framework **MITRE ATT&CK** pour structurer les techniques d'attaque.

---

## Qu'est-ce que MITRE ATT&CK ?

MITRE ATT&CK (Adversarial Tactics, Techniques, and Common Knowledge) est une base de connaissances des tactiques et techniques utilisees par les attaquants.

---

## Tactiques implementees

AutoStrike implemente actuellement **13 tactiques** avec **48 techniques** :

| ID | Tactique | Techniques | Description |
|----|----------|------------|-------------|
| TA0043 | Reconnaissance | 2 | Collecte d'informations sur la cible |
| TA0001 | Initial Access | 3 | Acces initial au reseau |
| TA0002 | Execution | 5 | Execution de commandes via shells |
| TA0003 | Persistence | 4 | Maintien de l'acces |
| TA0004 | Privilege Escalation | 4 | Elevation de privileges |
| TA0005 | Defense Evasion | 6 | Contournement des defenses |
| TA0006 | Credential Access | 4 | Vol d'identifiants |
| TA0007 | Discovery | 9 | Reconnaissance interne du systeme et reseau |
| TA0008 | Lateral Movement | 3 | Deplacement lateral dans le reseau |
| TA0009 | Collection | 4 | Collecte de donnees sur la cible |
| TA0011 | Command and Control | 3 | Communication avec les systemes compromis |
| TA0010 | Exfiltration | 3 | Exfiltration de donnees |
| TA0040 | Impact | 3 | Manipulation ou destruction de systemes |

---

## Techniques par tactique

### Reconnaissance (2 techniques)
| ID | Nom |
|----|-----|
| T1592.004 | Gather Victim Host Information: Client Configurations |
| T1595.002 | Active Scanning: Vulnerability Scanning |

### Initial Access (3 techniques)
| ID | Nom |
|----|-----|
| T1078 | Valid Accounts |
| T1133 | External Remote Services |
| T1190 | Exploit Public-Facing Application |

### Execution (5 techniques)
| ID | Nom |
|----|-----|
| T1059.001 | PowerShell |
| T1059.003 | Windows Command Shell |
| T1059.004 | Unix Shell |
| T1047 | Windows Management Instrumentation |
| T1059.006 | Python |

### Persistence (4 techniques)
| ID | Nom |
|----|-----|
| T1053.005 | Scheduled Task |
| T1547.001 | Registry Run Keys |
| T1053.003 | Cron |
| T1543.002 | Systemd Service |

### Privilege Escalation (4 techniques)
| ID | Nom |
|----|-----|
| T1548.001 | Setuid and Setgid |
| T1548.002 | Bypass User Account Control |
| T1078.003 | Local Accounts |
| T1134.001 | Token Impersonation/Theft |

### Defense Evasion (6 techniques)
| ID | Nom |
|----|-----|
| T1070.004 | File Deletion |
| T1562.001 | Disable or Modify Tools |
| T1027 | Obfuscated Files or Information |
| T1070.001 | Clear Windows Event Logs |
| T1036.005 | Match Legitimate Name or Location |
| T1218.011 | Rundll32 |

### Credential Access (4 techniques)
| ID | Nom |
|----|-----|
| T1552.001 | Credentials In Files |
| T1555.003 | Credentials from Web Browsers |
| T1003.008 | /etc/passwd and /etc/shadow |
| T1552.004 | Private Keys |

### Discovery (9 techniques)
| ID | Nom |
|----|-----|
| T1082 | System Information Discovery |
| T1083 | File and Directory Discovery |
| T1057 | Process Discovery |
| T1016 | System Network Configuration Discovery |
| T1049 | System Network Connections Discovery |
| T1087 | Account Discovery |
| T1069 | Permission Groups Discovery |
| T1018 | Remote System Discovery |
| T1007 | System Service Discovery |

### Lateral Movement (3 techniques)
| ID | Nom |
|----|-----|
| T1021.001 | Remote Desktop Protocol |
| T1021.002 | SMB/Windows Admin Shares |
| T1021.004 | SSH |

### Collection (4 techniques)
| ID | Nom |
|----|-----|
| T1005 | Data from Local System |
| T1039 | Data from Network Shared Drive |
| T1074.001 | Local Data Staging |
| T1119 | Automated Collection |

### Command and Control (3 techniques)
| ID | Nom |
|----|-----|
| T1071.001 | Web Protocols |
| T1105 | Ingress Tool Transfer |
| T1572 | Protocol Tunneling |

### Exfiltration (3 techniques)
| ID | Nom |
|----|-----|
| T1048.003 | Exfiltration Over Unencrypted Non-C2 Protocol |
| T1041 | Exfiltration Over C2 Channel |
| T1567.002 | Exfiltration to Cloud Storage |

### Impact (3 techniques)
| ID | Nom |
|----|-----|
| T1490 | Inhibit System Recovery |
| T1489 | Service Stop |
| T1486 | Data Encrypted for Impact |

---

## Safe Mode

Toutes les techniques sont **Safe Mode compatible** (non-destructives). Elles se limitent a :
- Lecture d'informations systeme
- Requetes read-only (registry, services, processes)
- Simulations sans modification
- Enumeration et detection

---

## Import automatique des techniques

AutoStrike inclut un outil d'import qui fusionne les donnees de **MITRE ATT&CK STIX 2.1** (metadonnees) et **Atomic Red Team** (commandes d'execution) pour generer automatiquement des techniques YAML.

### Utilisation

```bash
# Import complet (telecharge STIX + clone Atomic Red Team)
make import-mitre

# Import seulement les techniques safe (discovery, reconnaissance)
make import-mitre-safe

# Dry run : affiche les stats sans ecrire de fichiers
make import-mitre-dry
```

### Fonctionnement

1. **Telecharge** le fichier STIX `enterprise-attack.json` depuis GitHub (cache dans `~/.cache/autostrike/`)
2. **Clone** le repo Atomic Red Team en shallow clone (cache dans `~/.cache/autostrike/`)
3. **Inner join** : seules les techniques presentes dans les DEUX sources sont importees
4. **Genere** les fichiers YAML dans `server/configs/techniques/`, groupes par tactique

### Options du CLI

| Flag | Description | Defaut |
|------|-------------|--------|
| `--stix-path` | Chemin local vers enterprise-attack.json | Telecharge si absent |
| `--atomics-path` | Chemin local vers le repo Atomic Red Team | Clone si absent |
| `--output-dir` | Repertoire de sortie des YAML | `../../server/configs/techniques` |
| `--cache-dir` | Repertoire de cache | `~/.cache/autostrike` |
| `--dry-run` | Affiche les stats sans ecrire | `false` |
| `--safe-only` | N'importe que les techniques safe | `false` |
| `--force-download` | Re-telecharge meme si le cache existe | `false` |
| `--verbose` | Logs detailles | `false` |

### Heuristique Safe Mode

- Tactiques `discovery` et `reconnaissance` : `is_safe = true` par defaut
- Toutes les autres tactiques : `is_safe = false`
- Force `is_safe = false` si la commande contient des patterns dangereux (`rm -rf`, `del /f`, `format`, `kill`, `shutdown`, etc.)
- Force `is_safe = false` si `elevation_required: true`

### Mapping des types d'executor

| Atomic Red Team | AutoStrike | Plateforme |
|-----------------|------------|------------|
| `command_prompt` | `cmd` | Windows |
| `powershell` | `powershell` | Windows |
| `bash` | `bash` | Linux/macOS |
| `sh` | `sh` | Linux/macOS |
| `manual` | **Ignore** | - |

---

## Voir aussi

- [Liste detaillee des techniques](techniques.md)
- [MITRE ATT&CK Official](https://attack.mitre.org/)
