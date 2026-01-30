# MITRE ATT&CK

AutoStrike utilise le framework **MITRE ATT&CK** pour structurer les techniques d'attaque.

---

## Qu'est-ce que MITRE ATT&CK ?

MITRE ATT&CK (Adversarial Tactics, Techniques, and Common Knowledge) est une base de connaissances des tactiques et techniques utilisées par les attaquants.

---

## Tactiques supportées

| ID | Tactique | Description |
|----|----------|-------------|
| TA0001 | Initial Access | Accès initial au système |
| TA0002 | Execution | Exécution de code |
| TA0003 | Persistence | Maintien de l'accès |
| TA0004 | Privilege Escalation | Élévation de privilèges |
| TA0005 | Defense Evasion | Contournement des défenses |
| TA0006 | Credential Access | Accès aux identifiants |
| TA0007 | Discovery | Reconnaissance interne |
| TA0008 | Lateral Movement | Mouvement latéral |
| TA0009 | Collection | Collecte de données |
| TA0010 | Exfiltration | Exfiltration de données |
| TA0011 | Command and Control | Communication C2 |

---

## Couverture MVP

Le MVP d'AutoStrike couvre environ **40%** des techniques MITRE ATT&CK, avec un focus sur :

- **Discovery** (T1082, T1083, T1057, T1012, T1016)
- **Execution** (T1059, T1053)
- **Defense Evasion** (T1070, T1027)
- **Credential Access** (T1003, T1552)

---

## Voir aussi

- [Liste des techniques](techniques.md)
- [MITRE ATT&CK Official](https://attack.mitre.org/)
