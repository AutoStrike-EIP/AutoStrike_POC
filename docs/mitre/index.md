# MITRE ATT&CK

AutoStrike utilise le framework **MITRE ATT&CK** pour structurer les techniques d'attaque.

---

## Qu'est-ce que MITRE ATT&CK ?

MITRE ATT&CK (Adversarial Tactics, Techniques, and Common Knowledge) est une base de connaissances des tactiques et techniques utilisées par les attaquants.

---

## Tactiques implémentées

AutoStrike implémente actuellement **4 tactiques** avec **15 techniques** :

| ID | Tactique | Techniques | Description |
|----|----------|------------|-------------|
| TA0007 | Discovery | 9 | Reconnaissance interne du système et réseau |
| TA0002 | Execution | 3 | Exécution de commandes via shells |
| TA0003 | Persistence | 2 | Maintien de l'accès |
| TA0005 | Defense Evasion | 1 | Contournement des défenses |

---

## Techniques par tactique

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

### Execution (3 techniques)
| ID | Nom |
|----|-----|
| T1059.001 | PowerShell |
| T1059.003 | Windows Command Shell |
| T1059.004 | Unix Shell |

### Persistence (2 techniques)
| ID | Nom |
|----|-----|
| T1053.005 | Scheduled Task |
| T1547.001 | Registry Run Keys |

### Defense Evasion (1 technique)
| ID | Nom |
|----|-----|
| T1070.004 | File Deletion |

---

## Safe Mode

Toutes les techniques sont **Safe Mode compatible** (non-destructives). Elles se limitent à :
- Lecture d'informations système
- Requêtes read-only (registry, services, processes)
- Simulations sans modification

---

## Voir aussi

- [Liste détaillée des techniques](techniques.md)
- [MITRE ATT&CK Official](https://attack.mitre.org/)
