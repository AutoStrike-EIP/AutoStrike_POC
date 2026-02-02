# Démarrage rapide

Ce guide vous permet de lancer votre première simulation en 5 minutes.

---

## Prérequis

- Go 1.21+
- Node.js 18+ (pour le build)
- Rust 1.75+ (pour l'agent)

---

## 1. Installation

```bash
# Cloner le projet
git clone https://github.com/AutoStrike-EIP/AutoStrike.git
cd AutoStrike

# Installer les dépendances
make install
```

---

## 2. Démarrer AutoStrike

```bash
make run
```

Le serveur démarre sur **http://localhost:8443** et sert :

- `/` - Dashboard
- `/api/v1/*` - API REST
- `/ws/*` - WebSocket
- `/health` - Health check

!!! note "Authentification"
    L'authentification est **désactivée** par défaut en développement.
    Pour l'activer, configurez `JWT_SECRET` dans `server/.env`.

---

## 3. Vérifier les agents

Dans le menu **Agents**, vérifiez que vos agents sont connectés (statut "Online").

Pour connecter un agent :

```bash
make agent
```

---

## 4. Lancer un scénario

1. Allez dans **Scénarios**
2. Sélectionnez "Discovery - Basic"
3. Choisissez les agents cibles
4. Cliquez sur **Exécuter**

---

## 5. Analyser les résultats

Les résultats s'affichent en temps réel :

| Statut | Description |
|--------|-------------|
| **Blocked** | Technique bloquée par les défenses |
| **Detected** | Technique détectée mais non bloquée |
| **Missed** | Technique non détectée |

---

## 6. Consulter la matrice MITRE

La **Matrice MITRE ATT&CK** affiche votre couverture de détection avec un code couleur :

| Couleur | Signification |
|---------|---------------|
| Vert | Technique bloquée |
| Orange | Technique détectée |
| Rouge | Technique non détectée |
| Gris | Non testée |

---

## Commandes utiles

| Commande | Description |
|----------|-------------|
| `make run` | Démarrer le serveur |
| `make agent` | Connecter un agent |
| `make stop` | Arrêter les services |
| `make logs` | Voir les logs serveur |
| `make test` | Lancer les tests |

---

## Prochaines étapes

- [Créer un scénario personnalisé](../mitre/techniques.md)
- [Configurer les intégrations](../architecture/index.md)
- [Référence API](../api/index.md)
