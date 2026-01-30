# Démarrage rapide

Ce guide vous permet de lancer votre première simulation en 5 minutes.

---

## 1. Accéder au Dashboard

Ouvrez votre navigateur et accédez à :

```
https://localhost:8443
```

Connectez-vous avec les identifiants par défaut :

- **Email** : admin@autostrike.local
- **Mot de passe** : changeme

!!! warning "Sécurité"
    Changez le mot de passe par défaut immédiatement après la première connexion.

---

## 2. Vérifier les agents

Dans le menu **Agents**, vérifiez que vos agents sont connectés (statut "Online").

---

## 3. Lancer un scénario

1. Allez dans **Scénarios**
2. Sélectionnez "Discovery - Basic"
3. Choisissez les agents cibles
4. Cliquez sur **Exécuter**

---

## 4. Analyser les résultats

Les résultats s'affichent en temps réel :

- **Blocked** - Technique bloquée par les défenses
- **Detected** - Technique détectée mais non bloquée
- **Missed** - Technique non détectée

---

## 5. Consulter la matrice MITRE

La **Matrice MITRE ATT&CK** affiche votre couverture de détection avec un code couleur :

| Couleur | Signification |
|---------|---------------|
| Vert | Technique bloquée |
| Orange | Technique détectée |
| Rouge | Technique non détectée |
| Gris | Non testée |

---

## Prochaines étapes

- [Créer un scénario personnalisé](../mitre/techniques.md)
- [Configurer les intégrations](../architecture/backend.md)
- [Générer un rapport](../api/reference.md)
