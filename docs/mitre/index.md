# MITRE ATT&CK

AutoStrike utilise le framework **MITRE ATT&CK** pour structurer les techniques d'attaque.

---

## Qu'est-ce que MITRE ATT&CK ?

MITRE ATT&CK (Adversarial Tactics, Techniques, and Common Knowledge) est une base de connaissances des tactiques et techniques utilisees par les attaquants.

---

## Tactiques implementees

AutoStrike implemente actuellement **12 tactiques** avec **294 techniques** (importees via MITRE STIX + Atomic Red Team) :

| ID | Tactique | Techniques | Description |
|----|----------|------------|-------------|
| TA0001 | Initial Access | 4 | Acces initial au reseau |
| TA0002 | Execution | 22 | Execution de commandes via shells |
| TA0003 | Persistence | 44 | Maintien de l'acces |
| TA0004 | Privilege Escalation | 18 | Elevation de privileges |
| TA0005 | Defense Evasion | 89 | Contournement des defenses |
| TA0006 | Credential Access | 34 | Vol d'identifiants |
| TA0007 | Discovery | 30 | Reconnaissance interne du systeme et reseau |
| TA0008 | Lateral Movement | 8 | Deplacement lateral dans le reseau |
| TA0009 | Collection | 16 | Collecte de donnees sur la cible |
| TA0011 | Command and Control | 13 | Communication avec les systemes compromis |
| TA0010 | Exfiltration | 8 | Exfiltration de donnees |
| TA0040 | Impact | 8 | Manipulation ou destruction de systemes |

---

## Techniques

Les 294 techniques sont importees via `make import-mitre` (inner join MITRE STIX + Atomic Red Team). Il n'y a pas de techniques built-in — toutes proviennent de l'import. La liste complete est disponible via l'API (`GET /api/v1/techniques`), dans le dashboard, ou dans la [liste des techniques](techniques.md).

---

## Safe Mode

AutoStrike utilise un systeme de classification de securite a **deux niveaux** :

### Niveau technique
Le champ `is_safe` de la technique indique si au moins un executor est safe.

### Niveau executor
Chaque executor possede son propre champ `is_safe`, determine par :
- `!elevation_required` - ne necessite pas de privileges admin/root
- `!hasDangerousPattern(command)` - la commande ne contient pas de patterns dangereux
- `!hasDangerousPattern(cleanup)` - la commande de cleanup ne contient pas de patterns dangereux

### Patterns dangereux detectes
`rm -rf`, `del /f`, `rd /s`, `dd of=/dev/`, `mkfs`, `fdisk`, `format C:`, `shutdown`, `reboot`, `taskkill /f`, `kill -9`, `killall`, `pkill`, `systemctl stop/disable`, `chmod 000`, `iptables -F`

### Statistiques actuelles
- **220 techniques safe** (la majorite des techniques importees)
- **74 techniques unsafe** (necessitent elevation ou contiennent des commandes dangereuses)

---

## Import automatique des techniques

AutoStrike inclut un outil d'import qui fusionne les donnees de **MITRE ATT&CK STIX 2.1** (metadonnees) et **Atomic Red Team** (commandes d'execution) pour generer automatiquement des techniques YAML.

### Utilisation

```bash
# Import complet (telecharge STIX + clone Atomic Red Team)
make import-mitre

# Import seulement les techniques safe (toutes tactiques, filtre par executor)
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

La classification de securite est **per-executor** :

Pour chaque executor :
```
is_safe = !elevation_required && !hasDangerousPattern(command) && !hasDangerousPattern(cleanup)
```

Pour la technique globale :
- `is_safe = true` si au moins un executor est safe

Retro-compatibilite Safe Mode :
- Si `technique.IsSafe = true` mais qu'aucun executor n'a `is_safe` defini (format legacy), tous les executors sont consideres safe

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
