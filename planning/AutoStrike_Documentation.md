# 🎯 AutoStrike

## Plateforme de Breach and Attack Simulation (BAS)

> **Projet EIP EPITECH - Promotion 2028**
> 
> Validation continue des défenses de sécurité par simulation d'attaques basées sur le framework MITRE ATT&CK

---

## 📋 Table des Matières

**Validation Compétences EIP** (🔴 NON VALIDÉES - Mars 2026)
- [C1 - Recenser les besoins client/utilisateurs](#c1---recenser-les-besoins-du-client-et-des-utilisateurs)
- [C2 - Audit technique, fonctionnel et sécurité](#c2---réaliser-un-audit-technique-fonctionnel-et-de-sécurité)
- [C3 - Spécifications techniques et fonctionnelles](#c3---rédiger-les-spécifications-techniques-et-fonctionnelles)
- [C4 - Chiffrage et benchmark](#c4---chiffrer-le-projet-et-réaliser-un-benchmark)
- [C5 - Impacts et mitigation](#c5---prévoir-les-impacts-et-sécuriser-des-pistes-de-mitigation)

**Documentation Technique**
1. [Vision du Projet](#1-vision-du-projet)
2. [Analyse de Marché](#2-analyse-de-marché)
3. [Architecture Technique](#3-architecture-technique)
4. [Composants Détaillés](#4-composants-détaillés)
5. [Stack Technologique](#5-stack-technologique)
6. [Modèle de Données](#6-modèle-de-données)
7. [Protocoles de Communication](#7-protocoles-de-communication)
8. [Techniques MITRE ATT&CK & EBIOS RM](#8-techniques-mitre-attck)
   - [8.1 Techniques Prioritaires (MVP)](#81-techniques-prioritaires-mvp)
   - [8.2 Implémentation Type](#82-implémentation-type)
   - [8.3 Scénarios Prédéfinis](#83-scénarios-prédéfinis)
   - [8.4 Alignement EBIOS RM (ANSSI)](#84-alignement-ebios-rm-méthode-anssi)
9. [Interface Utilisateur](#9-interface-utilisateur)
10. [Sécurité](#10-sécurité)
11. [Roadmap](#11-roadmap)
12. [Organisation de l'Équipe](#12-organisation-de-léquipe)
13. [Ressources et Références](#13-ressources-et-références)

---

## 0. Validation des Compétences EIP

> **Cette section documente la méthodologie appliquée pour valider les compétences C1 à C5 du référentiel EPITECH.**

### C1 - Recenser les besoins du client et des utilisateurs

> **Définition officielle:** Recenser les besoins du client et des utilisateurs en observant et en échangeant avec les parties prenantes afin de cerner les usages prévus, notamment pour les personnes en situation de handicap.
>
> **Statut: 🔴 NON VALIDÉ** (validation prévue: Mars 2026)

#### 1.1 Méthodologie de Recueil des Besoins

##### Parties Prenantes Identifiées

| Partie Prenante | Rôle | Besoins Principaux | Mode de Consultation |
|-----------------|------|-------------------|----------------------|
| **RSSI / CISO** | Décideur | Dashboard exécutif, ROI sécurité, conformité | Interviews, questionnaires |
| **Blue Team / SOC** | Utilisateur principal | Validation détections, alertes, techniques | Ateliers, observations terrain |
| **Red Team / Pentesters** | Utilisateur avancé | Automatisation, personnalisation scénarios | Focus groups |
| **Administrateurs IT** | Support technique | Déploiement simple, impact minimal | Entretiens techniques |
| **Équipes Conformité** | Contrôle | Rapports audit, traçabilité | Questionnaires |
| **Direction Générale** | Sponsor | Coût, valeur business, risques | Présentations exécutives |

##### Techniques de Recueil Utilisées

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    PROCESSUS DE RECUEIL DES BESOINS                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Phase 1: DÉCOUVERTE                                                        │
│  ┌─────────────────┐   ┌─────────────────┐   ┌─────────────────┐           │
│  │   Interviews    │   │  Observations   │   │  Questionnaires │           │
│  │   individuelles │   │    terrain      │   │     en ligne    │           │
│  │   (30-60 min)   │   │   (SOC réel)    │   │  (Google Forms) │           │
│  └────────┬────────┘   └────────┬────────┘   └────────┬────────┘           │
│           │                     │                     │                     │
│           └─────────────────────┼─────────────────────┘                     │
│                                 ▼                                           │
│  Phase 2: ANALYSE                                                           │
│  ┌─────────────────────────────────────────────────────────────┐           │
│  │              Synthèse et Priorisation (MoSCoW)              │           │
│  │  Must Have │ Should Have │ Could Have │ Won't Have          │           │
│  └─────────────────────────────────────────────────────────────┘           │
│                                 │                                           │
│                                 ▼                                           │
│  Phase 3: VALIDATION                                                        │
│  ┌─────────────────┐   ┌─────────────────┐   ┌─────────────────┐           │
│  │    Ateliers     │   │   Prototypes    │   │   Validation    │           │
│  │  collaboratifs  │   │   interactifs   │   │    formelle     │           │
│  └─────────────────┘   └─────────────────┘   └─────────────────┘           │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

##### Personas Utilisateurs

**Persona 1: Sophie - Analyste SOC**
```
┌─────────────────────────────────────────────────────────────────┐
│  👤 Sophie Martin - Analyste SOC Niveau 2                       │
├─────────────────────────────────────────────────────────────────┤
│  Âge: 28 ans                                                    │
│  Expérience: 4 ans en cybersécurité                            │
│  Outils quotidiens: Splunk, CrowdStrike, TheHive               │
│                                                                 │
│  OBJECTIFS:                                                     │
│  • Valider que les règles de détection fonctionnent            │
│  • Réduire les faux positifs                                   │
│  • Documenter la couverture MITRE ATT&CK                       │
│                                                                 │
│  FRUSTRATIONS:                                                  │
│  • Tests manuels chronophages                                  │
│  • Pas de visibilité sur les lacunes de détection              │
│  • Difficile de justifier les investissements sécurité         │
│                                                                 │
│  CITATION:                                                      │
│  "Je veux savoir si mes alertes se déclenchent vraiment        │
│   quand une attaque se produit, pas juste l'espérer."          │
└─────────────────────────────────────────────────────────────────┘
```

**Persona 2: Thomas - RSSI**
```
┌─────────────────────────────────────────────────────────────────┐
│  👤 Thomas Dubois - RSSI PME (200 employés)                     │
├─────────────────────────────────────────────────────────────────┤
│  Âge: 42 ans                                                    │
│  Expérience: 15 ans IT, 8 ans sécurité                         │
│  Budget annuel sécurité: 150K€                                 │
│                                                                 │
│  OBJECTIFS:                                                     │
│  • Justifier le budget sécurité au COMEX                       │
│  • Conformité réglementaire (NIS2, ISO 27001)                  │
│  • Réduire le risque cyber mesurable                           │
│                                                                 │
│  FRUSTRATIONS:                                                  │
│  • Solutions BAS trop chères (>100K€)                          │
│  • Rapports techniques incompréhensibles pour la direction     │
│  • Manque de temps pour des pentests réguliers                 │
│                                                                 │
│  CITATION:                                                      │
│  "J'ai besoin de prouver à mon DG que nos 150K€ de sécu       │
│   servent vraiment à quelque chose."                           │
└─────────────────────────────────────────────────────────────────┘
```

**Persona 3: Marc - Pentester**
```
┌─────────────────────────────────────────────────────────────────┐
│  👤 Marc Leroy - Pentester / Red Teamer                         │
├─────────────────────────────────────────────────────────────────┤
│  Âge: 32 ans                                                    │
│  Certifications: OSCP, CRTO                                    │
│  Outils: Cobalt Strike, Metasploit, custom C2                  │
│                                                                 │
│  OBJECTIFS:                                                     │
│  • Automatiser les phases répétitives                          │
│  • Personnaliser les scénarios d'attaque                       │
│  • Tester les défenses avant les vraies attaques               │
│                                                                 │
│  FRUSTRATIONS:                                                  │
│  • Refaire les mêmes tests manuellement                        │
│  • Outils BAS pas assez flexibles                              │
│  • Pas de mode "stealth" pour tester l'évasion                 │
│                                                                 │
│  CITATION:                                                      │
│  "Je veux pouvoir créer mes propres scénarios basés sur       │
│   les APTs que je vois dans mes missions."                     │
└─────────────────────────────────────────────────────────────────┘
```

##### Accessibilité (RGAA / WCAG)

| Critère | Implémentation AutoStrike | Conformité |
|---------|---------------------------|------------|
| **Contraste couleurs** | Ratio 4.5:1 minimum, mode daltonien | WCAG AA |
| **Navigation clavier** | Toutes fonctions accessibles Tab/Enter | WCAG AA |
| **Lecteur d'écran** | Labels ARIA, alt-text | WCAG AA |
| **Taille texte** | Zoom 200% sans perte fonctionnelle | WCAG AA |
| **Alternatives couleurs** | Icônes + texte (pas couleur seule) | WCAG AA |

```tsx
// Exemple: Composant accessible pour la matrice
<TechniqueCell
  technique={tech}
  status={status}
  // Pas seulement la couleur - aussi icône et texte
  aria-label={`${tech.name}: ${statusLabels[status]}`}
  role="gridcell"
  tabIndex={0}
  onKeyDown={(e) => e.key === 'Enter' && onClick()}
>
  <StatusIcon status={status} /> {/* ✓ ⚠ ✗ ○ */}
  <span className="sr-only">{statusLabels[status]}</span>
</TechniqueCell>
```

##### User Stories Priorisées (MoSCoW)

**Must Have (MVP)**
| ID | User Story | Critères d'acceptation |
|----|------------|------------------------|
| US-001 | En tant qu'analyste SOC, je veux voir la couverture MITRE ATT&CK de mon SI | Matrice colorée avec légende |
| US-002 | En tant que RSSI, je veux un score de sécurité global | Score 0-100 avec tendance |
| US-003 | En tant qu'opérateur, je veux lancer un scénario prédéfini | Exécution en 1 clic |
| US-004 | En tant qu'admin, je veux déployer un agent facilement | Script/binaire one-liner |
| US-005 | En tant qu'utilisateur, je veux voir les résultats en temps réel | WebSocket live updates |
| US-006 | En tant qu'opérateur, je veux voir la liste des agents connectés | Liste avec statut online/offline |
| US-007 | En tant qu'utilisateur, je veux me connecter au dashboard | Authentification JWT |
| US-008 | En tant qu'administrateur, je veux configurer le serveur | Fichier config YAML |
| US-009 | En tant qu'opérateur, je veux stopper une exécution en cours | Bouton stop + confirmation |

**Should Have**
| ID | User Story |
|----|------------|
| US-010 | En tant que pentester, je veux créer des scénarios custom |
| US-011 | En tant que RSSI, je veux exporter un rapport PDF |
| US-012 | En tant qu'analyste, je veux filtrer par tactic/technique |
| US-013 | En tant qu'analyste SOC, je veux voir l'historique des exécutions |
| US-014 | En tant que RSSI, je veux comparer les scores entre périodes |
| US-015 | En tant qu'opérateur, je veux planifier des exécutions automatiques |
| US-016 | En tant qu'administrateur, je veux gérer les utilisateurs |
| US-017 | En tant qu'analyste SOC, je veux voir les détails d'une technique |
| US-018 | En tant que pentester, je veux importer et exporter des scénarios |
| US-019 | En tant qu'utilisateur, je veux recevoir des notifications email |

**Could Have**
| ID | User Story |
|----|------------|
| US-020 | En tant qu'utilisateur, je veux des profils adversaires APT |
| US-021 | En tant qu'admin, je veux gérer les permissions utilisateurs |

**Won't Have (V1)**
| ID | User Story |
|----|------------|
| US-030 | Intégration SIEM automatique |
| US-031 | Mode SaaS multi-tenant |

---

### C2 - Réaliser un audit technique, fonctionnel et de sécurité

> **Définition officielle:** Réaliser un audit technique, fonctionnel et de sécurité de l'environnement dans lequel s'inscrit le projet (infrastructure, système d'information, ressources humaines, ...) afin de proposer les solutions les plus adaptées au contexte, en analysant les solutions déjà en place et leurs effets.
>
> **Statut: 🔴 NON VALIDÉ** (validation prévue: Mars 2026)

#### 2.1 Méthodologie d'Audit

##### Périmètre d'Audit

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         PÉRIMÈTRE D'AUDIT AUTOSTRIKE                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    ENVIRONNEMENT TECHNIQUE                          │   │
│  │  • Infrastructure existante (serveurs, réseau, cloud)               │   │
│  │  • Outils de sécurité en place (EDR, SIEM, Firewall)               │   │
│  │  • Systèmes d'exploitation et versions                              │   │
│  │  • Politiques de sécurité (GPO, restrictions)                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    ENVIRONNEMENT FONCTIONNEL                        │   │
│  │  • Processus métier impactés                                        │   │
│  │  • Workflows sécurité existants                                     │   │
│  │  • Intégrations nécessaires (ticketing, SIEM, reporting)           │   │
│  │  • Contraintes opérationnelles (maintenance windows, SLA)          │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    ENVIRONNEMENT HUMAIN                             │   │
│  │  • Compétences équipes (SOC, IT, Red Team)                         │   │
│  │  • Disponibilité et charge de travail                               │   │
│  │  • Besoins en formation                                             │   │
│  │  • Résistance au changement                                         │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                    ENVIRONNEMENT SÉCURITÉ                           │   │
│  │  • Politique de sécurité (PSSI)                                     │   │
│  │  • Conformité réglementaire (NIS2, ISO 27001, RGPD)                │   │
│  │  • Gestion des accès et privilèges                                  │   │
│  │  • Procédures d'incident                                            │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

##### Grille d'Audit Technique

| Domaine | Éléments Audités | Méthode | Livrables |
|---------|------------------|---------|-----------|
| **Infrastructure** | Serveurs, réseau, cloud | Inventaire, schéma | Cartographie SI |
| **Sécurité Endpoint** | EDR, AV, HIPS | Tests, interviews | Matrice outils |
| **Sécurité Réseau** | Firewall, IDS/IPS, proxy | Config review | Gaps analysis |
| **SIEM/Logs** | Centralisation, rétention | Audit config | Coverage map |
| **IAM** | Droits, privilèges, MFA | Review policies | Rapport IAM |
| **Vulnérabilités** | Scan, patch management | Nessus/Qualys | État des lieux |

##### Analyse des Solutions Existantes

**Matrice Comparative Solutions de Validation Sécurité**

| Critère | Tests Manuels | Pentests Externes | Caldera | Atomic Red Team | AutoStrike |
|---------|---------------|-------------------|---------|-----------------|------------|
| **Coût annuel** | ~20K€ (temps) | 30-80K€ | Gratuit | Gratuit | Gratuit/Freemium |
| **Fréquence** | Ponctuel | 1-2x/an | Continue | Ponctuel | Continue |
| **Couverture MITRE** | ~10% | ~30% | ~60% | ~80% | ~40% (MVP) |
| **Facilité d'usage** | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐ | ⭐⭐⭐⭐ |
| **Automatisation** | ❌ | ❌ | ✅ | Partiel | ✅ |
| **Dashboard** | ❌ | PDF | Basique | ❌ | ✅ |
| **Personnalisation** | ✅ | Limitée | ✅ | ✅ | ✅ |
| **Support** | Interne | Inclus | Communauté | Communauté | Communauté |

##### Matrice SWOT de l'Environnement Type

```
┌─────────────────────────────────────┬─────────────────────────────────────┐
│            FORCES                   │           FAIBLESSES                │
├─────────────────────────────────────┼─────────────────────────────────────┤
│ • EDR déployé sur 90%+ endpoints    │ • Règles SIEM non optimisées        │
│ • Équipe SOC compétente             │ • Pas de visibilité couverture      │
│ • Budget sécurité existant          │ • Tests manuels chronophages        │
│ • Volonté d'amélioration            │ • Documentation obsolète            │
│ • Infrastructure moderne            │ • Turnover équipe IT                │
├─────────────────────────────────────┼─────────────────────────────────────┤
│          OPPORTUNITÉS               │            MENACES                  │
├─────────────────────────────────────┼─────────────────────────────────────┤
│ • Conformité NIS2 à venir           │ • Ransomwares en augmentation       │
│ • Automatisation possible           │ • Budget serré                      │
│ • Solutions open-source matures     │ • Complexité croissante SI          │
│ • Sensibilisation direction         │ • Pénurie talents cyber             │
│ • Cloud hybride flexible            │ • Shadow IT non contrôlé            │
└─────────────────────────────────────┴─────────────────────────────────────┘
```

##### Recommandations Issues de l'Audit

| Priorité | Recommandation | Effort | Impact | Quick Win |
|----------|----------------|--------|--------|-----------|
| 🔴 Haute | Déployer validation continue MITRE | Moyen | Fort | ❌ |
| 🔴 Haute | Optimiser règles SIEM Discovery | Faible | Fort | ✅ |
| 🟡 Moyenne | Former équipe SOC aux techniques | Moyen | Moyen | ❌ |
| 🟡 Moyenne | Documenter baseline sécurité | Moyen | Moyen | ❌ |
| 🟢 Basse | Intégrer reporting automatique | Faible | Faible | ✅ |

---

### C3 - Rédiger les spécifications techniques et fonctionnelles

> **Définition officielle:** Rédiger les spécifications techniques et fonctionnelles à partir des résultats de l'audit, afin de couvrir tous les besoins clients, en décrivant précisément tous les aspects techniques (spécifications techniques) et humains (spécifications fonctionnelles).
>
> **Statut: 🔴 NON VALIDÉ** (validation prévue: Mars 2026)

#### 3.1 Document de Spécifications Fonctionnelles (SFD)

##### Vue d'Ensemble Fonctionnelle

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    DIAGRAMME DE CONTEXTE AUTOSTRIKE                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│                              ┌─────────────┐                                │
│                              │  AutoStrike │                                │
│                              │   Platform  │                                │
│                              └──────┬──────┘                                │
│                                     │                                       │
│         ┌───────────────────────────┼───────────────────────────┐          │
│         │                           │                           │          │
│         ▼                           ▼                           ▼          │
│  ┌─────────────┐           ┌─────────────┐           ┌─────────────┐       │
│  │   Analyste  │           │    RSSI     │           │  Pentester  │       │
│  │     SOC     │           │             │           │             │       │
│  └─────────────┘           └─────────────┘           └─────────────┘       │
│  • Lance scénarios         • Consulte dashboard      • Crée scénarios      │
│  • Analyse résultats       • Exporte rapports        • Personnalise        │
│  • Valide détections       • Suit tendances          • Teste évasion       │
│                                                                             │
│                           SYSTÈMES EXTERNES                                 │
│         ┌───────────────────────────┼───────────────────────────┐          │
│         │                           │                           │          │
│         ▼                           ▼                           ▼          │
│  ┌─────────────┐           ┌─────────────┐           ┌─────────────┐       │
│  │     EDR     │           │    SIEM     │           │   Active    │       │
│  │ (détection) │           │   (logs)    │           │  Directory  │       │
│  └─────────────┘           └─────────────┘           └─────────────┘       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

##### Cas d'Utilisation Détaillés

**UC-001: Exécuter un Scénario d'Attaque**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  UC-001: Exécuter un Scénario d'Attaque                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Acteur principal: Analyste SOC                                             │
│  Préconditions:                                                             │
│    • Utilisateur authentifié avec rôle "Operator" minimum                  │
│    • Au moins 1 agent online                                               │
│    • Scénario existant et valide                                           │
│                                                                             │
│  Scénario nominal:                                                          │
│    1. L'utilisateur accède à la page "Scénarios"                           │
│    2. L'utilisateur sélectionne un scénario                                │
│    3. Le système affiche les détails et agents compatibles                 │
│    4. L'utilisateur sélectionne les agents cibles                          │
│    5. L'utilisateur clique sur "Exécuter"                                  │
│    6. Le système demande confirmation                                       │
│    7. L'utilisateur confirme                                               │
│    8. Le système planifie l'exécution                                      │
│    9. Le système affiche le monitoring temps réel                          │
│   10. Les résultats s'affichent au fur et à mesure                         │
│   11. Le système calcule et affiche le score final                         │
│                                                                             │
│  Scénarios alternatifs:                                                     │
│    4a. Aucun agent compatible → Message d'erreur                           │
│    8a. Agent offline pendant exécution → Skip technique, continuer        │
│    9a. Utilisateur annule → Arrêt propre, résultats partiels              │
│                                                                             │
│  Postconditions:                                                            │
│    • Résultats stockés en base                                             │
│    • Matrice MITRE mise à jour                                             │
│    • Score de sécurité recalculé                                           │
│    • Logs générés pour audit                                               │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

**UC-002: Déployer un Agent**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  UC-002: Déployer un Agent                                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Acteur principal: Administrateur IT                                        │
│  Préconditions:                                                             │
│    • Utilisateur authentifié avec rôle "Admin"                             │
│    • Serveur AutoStrike accessible depuis la cible                         │
│    • Droits administrateur sur la machine cible                            │
│                                                                             │
│  Scénario nominal:                                                          │
│    1. L'utilisateur accède à "Agents" > "Déployer"                         │
│    2. Le système affiche les options (Windows/Linux)                       │
│    3. L'utilisateur sélectionne la plateforme                              │
│    4. Le système génère une commande one-liner                             │
│    5. L'utilisateur copie la commande                                      │
│    6. L'utilisateur exécute sur la machine cible                           │
│    7. L'agent se télécharge et s'installe                                  │
│    8. L'agent s'enregistre auprès du serveur                               │
│    9. Le nouvel agent apparaît dans le dashboard                           │
│                                                                             │
│  Commande générée (exemple Windows):                                        │
│  ```                                                                        │
│  powershell -c "IEX(New-Object Net.WebClient).DownloadString(             │
│    'https://server:8443/deploy/agent.ps1')" -Server https://server:8443   │
│  ```                                                                        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

##### Règles de Gestion

| ID | Règle | Description |
|----|-------|-------------|
| RG-001 | Authentification obligatoire | Toute action requiert une session valide |
| RG-002 | Timeout session | Session expire après 30 min d'inactivité |
| RG-003 | Agent heartbeat | Agent considéré offline après 3 beacons manqués |
| RG-004 | Technique safe only | Mode production = techniques `is_safe: true` uniquement |
| RG-005 | Cleanup obligatoire | Toute persistence créée doit être supprimée après test |
| RG-006 | Logs immuables | Résultats ne peuvent pas être supprimés (audit) |
| RG-007 | Score calculation | Score = (blocked×100 + detected×50) / (total×100) |

#### 3.2 Document de Spécifications Techniques (STD)

##### Architecture Applicative

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      ARCHITECTURE APPLICATIVE N-TIERS                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                     COUCHE PRÉSENTATION                              │   │
│  │  ┌─────────────────────────────────────────────────────────────┐    │   │
│  │  │  React 18 + TypeScript + TailwindCSS + D3.js                │    │   │
│  │  │  • SPA (Single Page Application)                             │    │   │
│  │  │  • Responsive design (mobile-first)                          │    │   │
│  │  │  • WebSocket pour temps réel                                 │    │   │
│  │  └─────────────────────────────────────────────────────────────┘    │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                            HTTPS / WSS                                      │
│                                    ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      COUCHE SERVICE (API)                           │   │
│  │  ┌─────────────────────────────────────────────────────────────┐    │   │
│  │  │  Go 1.21 + Gin Framework                                     │    │   │
│  │  │  • REST API (JSON)                                           │    │   │
│  │  │  • WebSocket server                                          │    │   │
│  │  │  • JWT Authentication                                        │    │   │
│  │  │  • Rate limiting                                             │    │   │
│  │  └─────────────────────────────────────────────────────────────┘    │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                                    ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      COUCHE MÉTIER (DOMAIN)                         │   │
│  │  ┌─────────────────────────────────────────────────────────────┐    │   │
│  │  │  Architecture Hexagonale                                     │    │   │
│  │  │  • Entities: Agent, Scenario, Technique, Result              │    │   │
│  │  │  • Services: Orchestrator, Validator, ScoreCalculator        │    │   │
│  │  │  • Ports: Repositories (interfaces)                          │    │   │
│  │  └─────────────────────────────────────────────────────────────┘    │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                    │                                        │
│                                    ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      COUCHE DONNÉES                                 │   │
│  │  ┌───────────────────┐    ┌───────────────────┐                    │   │
│  │  │  SQLite (MVP)     │    │  File System      │                    │   │
│  │  │  → PostgreSQL     │    │  (Techniques YAML)│                    │   │
│  │  └───────────────────┘    └───────────────────┘                    │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

##### Spécifications des Interfaces (API)

**OpenAPI 3.0 - Extrait**

```yaml
openapi: 3.0.3
info:
  title: AutoStrike API
  version: 1.0.0
  description: API de la plateforme BAS AutoStrike

servers:
  - url: https://localhost:8443/api/v1

paths:
  /scenarios/{id}/execute:
    post:
      summary: Exécuter un scénario
      tags: [Executions]
      security:
        - bearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ExecuteRequest'
      responses:
        '202':
          description: Exécution démarrée
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Execution'
        '400':
          description: Requête invalide
        '404':
          description: Scénario non trouvé

components:
  schemas:
    ExecuteRequest:
      type: object
      required:
        - agent_paws
      properties:
        agent_paws:
          type: array
          items:
            type: string
          minItems: 1
          example: ["agent-001", "agent-002"]
        safe_mode:
          type: boolean
          default: true
          description: Exécuter uniquement les techniques safe

    Execution:
      type: object
      properties:
        id:
          type: string
          format: uuid
        scenario_id:
          type: string
        status:
          type: string
          enum: [pending, running, completed, failed, cancelled]
        started_at:
          type: string
          format: date-time
        progress:
          type: object
          properties:
            current:
              type: integer
            total:
              type: integer

  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
```

##### Matrice de Traçabilité Exigences

| Exigence | User Story | Composant | Test | Status |
|----------|------------|-----------|------|--------|
| EX-001 | US-001 | MitreMatrix.tsx | TC-001 | 🟢 |
| EX-002 | US-002 | SecurityScore.tsx | TC-002 | 🟢 |
| EX-003 | US-003 | ScenarioService.go | TC-003 | 🟡 |
| EX-004 | US-004 | deploy.sh | TC-004 | 🟡 |
| EX-005 | US-005 | websocket/handler.go | TC-005 | 🔴 |

---

### C4 - Chiffrer le projet et réaliser un benchmark

> **Définition officielle:** Chiffrer le projet en calculant les éléments financiers de la solution technique et en réalisant un benchmark des solutions existantes afin de cadrer les prévisions budgétaires.
>
> **Statut: 🔴 NON VALIDÉ** (validation prévue: Mars 2026)

#### 4.1 Benchmark Solutions Existantes

##### Analyse Tarifaire du Marché

| Solution | Type | Prix Annuel | Licence | Support |
|----------|------|-------------|---------|---------|
| **Pentera** | SaaS | 150-300K€ | Per-asset | Premium inclus |
| **AttackIQ** | SaaS/On-prem | 100-250K€ | Per-agent | Niveaux |
| **SafeBreach** | SaaS | 120-200K€ | Per-simulator | Premium inclus |
| **Cymulate** | SaaS | 80-150K€ | Per-vector | Niveaux |
| **Picus** | SaaS | 60-120K€ | Per-integration | Inclus |
| **MITRE Caldera** | Open-source | 0€ | Apache 2.0 | Communauté |
| **Atomic Red Team** | Open-source | 0€ | MIT | Communauté |
| **AutoStrike** | Open-source | 0€ | MIT | Communauté |

##### Grille de Scoring Fonctionnel

| Fonctionnalité | Poids | Pentera | AttackIQ | Caldera | AutoStrike |
|----------------|-------|---------|----------|---------|------------|
| Facilité déploiement | 15% | 4/5 | 3/5 | 2/5 | 4/5 |
| Couverture MITRE | 20% | 5/5 | 5/5 | 4/5 | 3/5 |
| Dashboard UX | 15% | 5/5 | 4/5 | 2/5 | 4/5 |
| Personnalisation | 15% | 3/5 | 4/5 | 5/5 | 4/5 |
| Rapports | 10% | 5/5 | 4/5 | 2/5 | 3/5 |
| Intégrations | 10% | 5/5 | 5/5 | 3/5 | 2/5 |
| Prix | 15% | 1/5 | 2/5 | 5/5 | 5/5 |
| **Score Total** | 100% | **3.8** | **3.85** | **3.35** | **3.65** |

#### 4.2 Chiffrage du Projet AutoStrike

##### Coûts de Développement (Ressources Humaines)

| Profil | Nb | Durée | TJM Marché | Jours Travaillés | Coût Total |
|--------|-----|-------|------------|------------------|------------|
| Tech Lead / Architecte | 1 | 12 mois | 600€ | 220 j | 132 000€ |
| Dev Backend Go | 1 | 12 mois | 450€ | 220 j | 99 000€ |
| Dev Frontend React | 1 | 10 mois | 400€ | 183 j | 73 200€ |
| Dev Rust (Agent) | 1 | 8 mois | 500€ | 147 j | 73 500€ |
| DevOps | 0.5 | 12 mois | 450€ | 110 j | 49 500€ |
| **Total RH** | | | | | **427 200€** |

*Note: Dans le contexte EIP, ce coût représente le temps passé par l'équipe projet.*

##### Coûts d'Infrastructure (Environnement Dev/Prod)

| Ressource | Spécification | Coût Mensuel | Coût Annuel |
|-----------|---------------|--------------|-------------|
| Serveur Dev | 4 vCPU, 8GB RAM, 100GB SSD | 50€ | 600€ |
| Serveur Prod | 8 vCPU, 16GB RAM, 200GB SSD | 150€ | 1 800€ |
| CI/CD (GitHub Actions) | 3000 min/mois | 40€ | 480€ |
| Domaine + SSL | autostrike.io | 50€ | 50€/an |
| Backup S3 | 100GB | 10€ | 120€ |
| **Total Infra** | | | **3 050€** |

##### Coûts Outillage et Licences

| Outil | Usage | Coût Annuel |
|-------|-------|-------------|
| GitHub Team | Repos privés, CI/CD | 0€ (éducation) |
| Figma | Design UI/UX | 0€ (éducation) |
| JetBrains | IDE (GoLand, CLion) | 0€ (éducation) |
| Notion | Documentation | 0€ (éducation) |
| **Total Outils** | | **0€** |

##### Budget Total Projet

| Catégorie | Coût |
|-----------|------|
| Ressources Humaines (valorisé) | 427 200€ |
| Infrastructure | 3 050€ |
| Outillage | 0€ |
| **Total Valorisé** | **430 250€** |
| **Coût Réel (infrastructure seule)** | **3 050€** |

##### Analyse ROI pour Client Type

**Contexte: PME 200 employés, budget sécurité 150K€/an**

| Approche | Coût Annuel | Fréquence Tests | Couverture |
|----------|-------------|-----------------|------------|
| Pentest externe | 30 000€ | 2x/an | ~30% MITRE |
| Solution BAS commerciale | 80 000€ | Continue | ~70% MITRE |
| **AutoStrike** | 3 000€* | Continue | ~40% MITRE |

*Coût d'hébergement et maintenance interne*

**ROI AutoStrike vs Pentest:**
- Économie: 27 000€/an
- + Tests continus vs ponctuels
- + Visibilité permanente couverture

---

### C5 - Prévoir les impacts et sécuriser des pistes de mitigation

> **Définition officielle:** Prévoir les impacts techniques et fonctionnels de la solution préconisée, afin de sécuriser des pistes de mitigation le cas échéant, en s'assurant de sa bonne intégration dans l'environnement d'exploitation du client.
>
> **Statut: 🔴 NON VALIDÉ** (validation prévue: Mars 2026)

#### 5.1 Analyse des Impacts

##### Matrice d'Impact

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         MATRICE D'IMPACT AUTOSTRIKE                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Impact Technique                                                           │
│  ├─ 🟡 Charge réseau (beaconing agents)                                    │
│  ├─ 🟡 Ressources CPU/RAM sur endpoints                                    │
│  ├─ 🟢 Compatibilité OS (Windows 10+, Ubuntu 20+)                          │
│  ├─ 🔴 Faux positifs EDR (détection agent)                                 │
│  └─ 🟡 Ouverture ports firewall (8443 inbound)                             │
│                                                                             │
│  Impact Fonctionnel                                                         │
│  ├─ 🟢 Intégration workflow SOC (non bloquant)                             │
│  ├─ 🟡 Formation utilisateurs (2-4h)                                       │
│  ├─ 🟢 Processus existants (complémentaire)                                │
│  └─ 🟡 Maintenance (mises à jour techniques MITRE)                         │
│                                                                             │
│  Impact Organisationnel                                                     │
│  ├─ 🟢 Pas de changement structure équipe                                  │
│  ├─ 🟡 Nouvelle responsabilité (ownership plateforme)                      │
│  └─ 🟢 Montée en compétences équipe                                        │
│                                                                             │
│  Impact Sécurité                                                            │
│  ├─ 🔴 Risque si agent compromis (C2 légitime)                             │
│  ├─ 🟡 Exposition nouveau service (serveur AutoStrike)                     │
│  ├─ 🟢 Amélioration posture sécurité globale                               │
│  └─ 🟡 Conformité (logs, traçabilité)                                      │
│                                                                             │
│  Légende: 🟢 Faible  🟡 Moyen  🔴 Élevé                                    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### 5.2 Plan de Mitigation des Risques

| Risque | Probabilité | Impact | Mitigation | Responsable |
|--------|-------------|--------|------------|-------------|
| **R1: Agent détecté par EDR** | Haute | Moyen | Whitelist agent, signature connue | Admin IT |
| **R2: Agent compromis/détourné** | Basse | Critique | mTLS, révocation certificat, audit logs | RSSI |
| **R3: Surcharge réseau** | Moyenne | Faible | Rate limiting, beaconing ajustable | Ops |
| **R4: Technique cause dommage** | Basse | Élevé | Mode safe par défaut, review techniques | Dev |
| **R5: Indisponibilité serveur** | Moyenne | Moyen | HA optionnel, backup config agents | Ops |
| **R6: Fuite données résultats** | Basse | Élevé | Chiffrement DB, accès RBAC | RSSI |

##### Plan de Mitigation Détaillé - R1: Agent détecté par EDR

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  RISQUE R1: Agent AutoStrike détecté comme menace par EDR                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  CONTEXTE:                                                                  │
│  L'agent AutoStrike exécute des techniques qui ressemblent à des attaques. │
│  Les EDR modernes peuvent le détecter et le bloquer.                       │
│                                                                             │
│  MITIGATIONS:                                                               │
│                                                                             │
│  1. Whitelist par Hash (Préventif)                                         │
│     ┌─────────────────────────────────────────────────────────────────┐    │
│     │  • Calculer SHA256 de l'agent compilé                          │    │
│     │  • Ajouter à la liste d'exclusion EDR                          │    │
│     │  • Documenter dans runbook déploiement                         │    │
│     └─────────────────────────────────────────────────────────────────┘    │
│                                                                             │
│  2. Whitelist par Certificat (Recommandé)                                  │
│     ┌─────────────────────────────────────────────────────────────────┐    │
│     │  • Signer l'agent avec certificat code signing                  │    │
│     │  • Configurer EDR pour faire confiance au certificat           │    │
│     │  • Renouvellement annuel certificat                            │    │
│     └─────────────────────────────────────────────────────────────────┘    │
│                                                                             │
│  3. Whitelist par Chemin (Fallback)                                        │
│     ┌─────────────────────────────────────────────────────────────────┐    │
│     │  • Installer agent dans répertoire dédié                        │    │
│     │  • Ex: C:\Program Files\AutoStrike\agent.exe                   │    │
│     │  • Exclure ce chemin dans EDR                                  │    │
│     └─────────────────────────────────────────────────────────────────┘    │
│                                                                             │
│  PROCÉDURE D'URGENCE:                                                       │
│  Si agent bloqué en production:                                             │
│  1. Vérifier logs EDR pour identifier la règle déclenchée                  │
│  2. Ajouter exclusion temporaire                                           │
│  3. Analyser si comportement attendu ou bug                                │
│  4. Mettre à jour documentation                                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### 5.3 Plan d'Intégration Environnement Client

##### Prérequis Techniques

| Catégorie | Prérequis | Vérification |
|-----------|-----------|--------------|
| **Réseau** | Port 8443/TCP ouvert vers serveur | `telnet server 8443` |
| **Réseau** | DNS résolution serveur | `nslookup autostrike.local` |
| **Endpoint** | Windows 10+ ou Ubuntu 20+ | `winver` / `lsb_release -a` |
| **Endpoint** | Droits admin pour installation | `whoami /groups` |
| **Sécurité** | Exclusion EDR configurée | Test manuel |
| **Sécurité** | Certificat CA déployé | `certutil -verify` |

##### Checklist Déploiement

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      CHECKLIST DÉPLOIEMENT AUTOSTRIKE                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Phase 1: Préparation (J-7)                                                 │
│  ☐ Valider prérequis techniques                                            │
│  ☐ Obtenir approbations (RSSI, IT)                                         │
│  ☐ Planifier fenêtre de déploiement                                        │
│  ☐ Préparer exclusions EDR                                                 │
│  ☐ Générer certificats mTLS                                                │
│                                                                             │
│  Phase 2: Installation Serveur (J-1)                                        │
│  ☐ Déployer serveur AutoStrike                                             │
│  ☐ Configurer HTTPS/certificats                                            │
│  ☐ Tester accès dashboard                                                  │
│  ☐ Créer comptes utilisateurs                                              │
│  ☐ Importer techniques MITRE                                               │
│                                                                             │
│  Phase 3: Déploiement Agents (J0)                                          │
│  ☐ Déployer sur 1 machine test                                             │
│  ☐ Valider enregistrement agent                                            │
│  ☐ Exécuter scénario test                                                  │
│  ☐ Vérifier pas de blocage EDR                                             │
│  ☐ Déployer sur scope complet                                              │
│                                                                             │
│  Phase 4: Validation (J+1)                                                  │
│  ☐ Vérifier tous agents online                                             │
│  ☐ Exécuter scénario complet                                               │
│  ☐ Valider résultats cohérents                                             │
│  ☐ Former utilisateurs clés                                                │
│  ☐ Documenter configuration                                                │
│                                                                             │
│  Phase 5: Hypercare (J+7)                                                   │
│  ☐ Monitoring quotidien                                                    │
│  ☐ Support utilisateurs                                                    │
│  ☐ Ajustements configuration                                               │
│  ☐ Handover équipe interne                                                 │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

##### Plan de Rollback

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           PROCÉDURE DE ROLLBACK                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  DÉCLENCHEURS:                                                              │
│  • Impact production détecté (CPU >80%, réseau saturé)                     │
│  • Incident sécurité lié à AutoStrike                                      │
│  • Demande explicite du RSSI                                               │
│                                                                             │
│  ÉTAPES ROLLBACK:                                                           │
│                                                                             │
│  1. Arrêt immédiat des exécutions                                          │
│     POST /api/v1/executions/stop-all                                       │
│                                                                             │
│  2. Désactivation agents (ne pas désinstaller)                             │
│     • Agents continuent à beacon mais n'exécutent plus                     │
│     • Permet réactivation rapide si fausse alerte                         │
│                                                                             │
│  3. Si nécessaire: désinstallation agents                                  │
│     Windows: C:\Program Files\AutoStrike\uninstall.bat                     │
│     Linux: /opt/autostrike/uninstall.sh                                    │
│                                                                             │
│  4. Conservation données                                                    │
│     • Backup base de données                                               │
│     • Export logs pour analyse post-mortem                                 │
│                                                                             │
│  5. Communication                                                           │
│     • Informer parties prenantes                                           │
│     • Documenter incident                                                  │
│     • Planifier analyse root cause                                         │
│                                                                             │
│  TEMPS ESTIMÉ: < 30 minutes pour rollback complet                          │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 1. Vision du Projet

### 1.1 Problématique

Les entreprises investissent des sommes considérables dans leurs infrastructures de sécurité (EDR, SIEM, Firewalls, etc.), mais n'ont souvent aucun moyen de valider l'efficacité réelle de ces défenses face aux techniques d'attaque actuelles.

**Questions sans réponse pour la plupart des organisations :**
- Mes 200K€ d'outils de sécurité servent-ils vraiment ?
- Mon EDR détecte-t-il les techniques utilisées par les ransomwares actuels ?
- Quelles sont mes lacunes de détection dans la matrice MITRE ATT&CK ?
- Si un attaquant pénètre mon réseau, jusqu'où peut-il aller ?

### 1.2 Solution : AutoStrike

AutoStrike est une plateforme de **Breach and Attack Simulation (BAS)** qui permet de :

- ✅ **Simuler** des attaques réelles de manière sécurisée en production
- ✅ **Valider** l'efficacité des contrôles de sécurité existants
- ✅ **Visualiser** la couverture de détection via la matrice MITRE ATT&CK
- ✅ **Prioriser** les remédiations basées sur des preuves concrètes
- ✅ **Automatiser** les tests de sécurité de manière continue

### 1.3 Proposition de Valeur

| Pour qui | Valeur apportée |
|----------|-----------------|
| **Blue Teams / SOC** | Valider que les alertes se déclenchent correctement |
| **RSSI / CISO** | Dashboard exécutif de posture de sécurité |
| **Pentesters** | Automatiser les phases de reconnaissance et validation |
| **Écoles / Formations** | Plateforme pédagogique pour apprendre les TTPs |
| **PME** | Alternative accessible aux solutions enterprise (Pentera, AttackIQ) |

### 1.4 Différenciation

| Aspect | Solutions Enterprise (Pentera, AttackIQ) | AutoStrike |
|--------|------------------------------------------|------------|
| **Prix** | 100-200K€/an | Open-source / Freemium |
| **Complexité** | Déploiement enterprise | Léger, déploiement rapide |
| **Focus** | Couverture maximale | Pédagogie + Essentiel |
| **Personnalisation** | Limitée | Scénarios custom, code ouvert |
| **Cible** | Grandes entreprises | PME, écoles, Blue Teams |

---

## 2. Analyse de Marché

### 2.1 Marché BAS Global

| Métrique | Valeur |
|----------|--------|
| Taille marché 2025 | $7.37 milliards |
| Projection 2030 | $14.66 milliards |
| CAGR | 12.2% |
| Croissance principale | Cloud security validation |

### 2.2 Acteurs Principaux

#### Solutions Enterprise
| Vendor | Positionnement | Prix indicatif |
|--------|----------------|----------------|
| **Pentera** | Automated pentesting + BAS | 150-300K€/an |
| **AttackIQ** | BAS + Purple teaming | 100-250K€/an |
| **SafeBreach** | BAS + Threat intelligence | 120-200K€/an |
| **Cymulate** | BAS + Attack Surface Management | 80-150K€/an |
| **Picus Security** | BAS + Security Control Validation | 60-120K€/an |

#### Solutions Open-Source
| Projet | Description | Limitations |
|--------|-------------|-------------|
| **MITRE Caldera** | Adversary emulation platform | Complexe, orienté red team |
| **Atomic Red Team** | Bibliothèque de tests ATT&CK | Pas de dashboard, manuel |
| **OpenBAS (Filigran)** | BAS open-source récent | Écosystème jeune |
| **Infection Monkey** | Breach simulation | Focus réseau uniquement |

### 2.3 Positionnement AutoStrike

```
                    COMPLEXITÉ
                        ▲
                        │
          Pentera ●     │     ● AttackIQ
                        │
                        │
    ────────────────────┼────────────────────► PRIX
                        │
         AutoStrike ●   │     ● Picus
                        │
      Atomic Red Team ● │
                        │
```

**AutoStrike se positionne comme :**
- Plus accessible que les solutions enterprise
- Plus intégré et user-friendly que Caldera/Atomic Red Team
- Focus sur la pédagogie et la compréhension des résultats

---

## 3. Architecture Technique

### 3.1 Vue d'Ensemble

```
┌─────────────────────────────────────────────────────────────────────────┐
│                                                                         │
│                         🌐 DASHBOARD (React)                            │
│                                                                         │
│   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐   │
│   │   Matrice   │  │  Scénarios  │  │   Agents    │  │  Rapports   │   │
│   │   ATT&CK    │  │   Builder   │  │   Manager   │  │   Export    │   │
│   └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘   │
│                                                                         │
└─────────────────────────────────────┬───────────────────────────────────┘
                                      │
                                      │ REST API + WebSocket
                                      ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                                                                         │
│                      🖥️ CONTROL SERVER (Go)                            │
│                                                                         │
│   ┌─────────────────────────────────────────────────────────────────┐   │
│   │                    Architecture Hexagonale                       │   │
│   │  ┌───────────┐    ┌───────────────┐    ┌───────────────────┐   │   │
│   │  │  Adapters │───▶│    Domain     │◀───│     Adapters      │   │   │
│   │  │  (API)    │    │ (Orchestrator)│    │ (DB, Agent Comm)  │   │   │
│   │  └───────────┘    └───────────────┘    └───────────────────┘   │   │
│   └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────┬───────────────────────────────────┘
                                      │
                                      │ mTLS / HTTPS (Beaconing)
                                      ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                                                                         │
│                          🔧 AGENTS (Rust)                               │
│                                                                         │
│   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐   │
│   │   Agent     │  │   Agent     │  │   Agent     │  │   Agent     │   │
│   │  Windows    │  │   Linux     │  │  Windows    │  │   Linux     │   │
│   │  PC-001     │  │  SRV-001    │  │  PC-002     │  │  SRV-002    │   │
│   └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘   │
│                                                                         │
│                    Exécution des Techniques MITRE                       │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Architecture Hexagonale (Control Server)

```
                         DRIVING ADAPTERS
                    ┌─────────────────────┐
                    │    REST API (Gin)   │
                    │   WebSocket Handler │
                    │    CLI Interface    │
                    └──────────┬──────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │    INBOUND PORTS    │
                    │  (Use Case Interfaces)
                    │                     │
                    │ • ScenarioService   │
                    │ • AgentService      │
                    │ • ResultService     │
                    │ • TechniqueService  │
                    └──────────┬──────────┘
                               │
                               ▼
┌──────────────────────────────────────────────────────────────────┐
│                                                                  │
│                        DOMAIN CORE                               │
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │                      ENTITIES                              │  │
│  │  • Agent (paw, hostname, platform, status)                 │  │
│  │  • Scenario (name, phases, techniques)                     │  │
│  │  • Technique (id, name, tactic, commands)                  │  │
│  │  • ExecutionResult (status, output, detected_by)           │  │
│  │  • AdversaryProfile (name, description, techniques)        │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │                   DOMAIN SERVICES                          │  │
│  │  • AttackOrchestrator (planification, ordonnancement)      │  │
│  │  • TechniqueValidator (compatibilité agent/technique)      │  │
│  │  • ScoreCalculator (calcul score sécurité)                 │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │                    VALUE OBJECTS                           │  │
│  │  • TacticType, ResultStatus, AgentStatus                   │  │
│  │  • TechniqueID, ExecutionPlan                              │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │   OUTBOUND PORTS    │
                    │ (Repository Interfaces)
                    │                     │
                    │ • AgentRepository   │
                    │ • ScenarioRepository│
                    │ • ResultRepository  │
                    │ • AgentCommunicator │
                    └──────────┬──────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │   DRIVEN ADAPTERS   │
                    │                     │
                    │ • SQLite/PostgreSQL │
                    │ • HTTP Agent Comm   │
                    │ • File System       │
                    │ • MITRE ATT&CK API  │
                    └─────────────────────┘
```

### 3.3 Structure des Répertoires

```
autostrike/
│
├── 📁 server/                          # Control Server (Go)
│   ├── cmd/
│   │   └── autostrike/
│   │       └── main.go                 # Point d'entrée, DI
│   │
│   ├── internal/
│   │   ├── domain/                     # 🎯 Cœur métier
│   │   │   ├── entity/
│   │   │   │   ├── agent.go
│   │   │   │   ├── scenario.go
│   │   │   │   ├── technique.go
│   │   │   │   └── result.go
│   │   │   ├── service/
│   │   │   │   ├── orchestrator.go
│   │   │   │   ├── validator.go
│   │   │   │   └── score_calculator.go
│   │   │   ├── repository/             # Interfaces outbound
│   │   │   │   ├── agent_repo.go
│   │   │   │   ├── scenario_repo.go
│   │   │   │   └── result_repo.go
│   │   │   └── valueobject/
│   │   │       ├── tactic.go
│   │   │       └── status.go
│   │   │
│   │   ├── application/                # Use cases
│   │   │   ├── scenario_service.go
│   │   │   ├── agent_service.go
│   │   │   ├── result_service.go
│   │   │   └── dto/
│   │   │       ├── requests.go
│   │   │       └── responses.go
│   │   │
│   │   └── infrastructure/             # Adapters
│   │       ├── persistence/
│   │       │   ├── sqlite/
│   │       │   │   ├── connection.go
│   │       │   │   ├── agent_repo.go
│   │       │   │   ├── scenario_repo.go
│   │       │   │   └── migrations/
│   │       │   └── memory/             # Pour tests
│   │       │       └── agent_repo.go
│   │       ├── communication/
│   │       │   └── http/
│   │       │       └── agent_comm.go
│   │       └── api/
│   │           ├── rest/
│   │           │   ├── router.go
│   │           │   ├── handlers/
│   │           │   │   ├── agent_handler.go
│   │           │   │   ├── scenario_handler.go
│   │           │   │   └── result_handler.go
│   │           │   └── middleware/
│   │           │       ├── auth.go
│   │           │       └── logging.go
│   │           └── websocket/
│   │               └── handler.go
│   │
│   ├── pkg/                            # Code partageable
│   │   ├── mitre/
│   │   │   ├── parser.go
│   │   │   └── navigator.go
│   │   └── crypto/
│   │       └── tls.go
│   │
│   ├── configs/
│   │   ├── config.yaml
│   │   └── techniques/                 # Définitions YAML
│   │       ├── T1059.001.yaml
│   │       ├── T1082.yaml
│   │       └── ...
│   │
│   ├── go.mod
│   ├── go.sum
│   └── Makefile
│
├── 📁 agent/                           # Agent (Rust)
│   ├── src/
│   │   ├── main.rs
│   │   ├── config.rs
│   │   ├── beacon/
│   │   │   ├── mod.rs
│   │   │   ├── client.rs
│   │   │   └── protocol.rs
│   │   ├── executor/
│   │   │   ├── mod.rs
│   │   │   ├── powershell.rs
│   │   │   ├── cmd.rs
│   │   │   ├── bash.rs
│   │   │   └── traits.rs
│   │   ├── techniques/
│   │   │   ├── mod.rs
│   │   │   ├── discovery/
│   │   │   │   ├── mod.rs
│   │   │   │   ├── t1082_system_info.rs
│   │   │   │   ├── t1083_file_discovery.rs
│   │   │   │   └── t1057_process_discovery.rs
│   │   │   ├── execution/
│   │   │   │   ├── mod.rs
│   │   │   │   └── t1059_command_scripting.rs
│   │   │   ├── persistence/
│   │   │   │   ├── mod.rs
│   │   │   │   ├── t1053_scheduled_task.rs
│   │   │   │   └── t1547_registry_run.rs
│   │   │   └── defense_evasion/
│   │   │       ├── mod.rs
│   │   │       └── t1070_indicator_removal.rs
│   │   └── utils/
│   │       ├── mod.rs
│   │       ├── system.rs
│   │       └── crypto.rs
│   │
│   ├── Cargo.toml
│   ├── Cargo.lock
│   └── build.rs                        # Cross-compilation
│
├── 📁 dashboard/                       # Frontend (React)
│   ├── src/
│   │   ├── components/
│   │   │   ├── Layout/
│   │   │   │   ├── Sidebar.tsx
│   │   │   │   ├── Header.tsx
│   │   │   │   └── Layout.tsx
│   │   │   ├── AttackMatrix/
│   │   │   │   ├── MitreMatrix.tsx
│   │   │   │   ├── TacticColumn.tsx
│   │   │   │   ├── TechniqueCell.tsx
│   │   │   │   ├── HeatmapLegend.tsx
│   │   │   │   └── TechniqueDetails.tsx
│   │   │   ├── Scenarios/
│   │   │   │   ├── ScenarioList.tsx
│   │   │   │   ├── ScenarioBuilder.tsx
│   │   │   │   ├── PhaseEditor.tsx
│   │   │   │   └── AdversarySelector.tsx
│   │   │   ├── Agents/
│   │   │   │   ├── AgentList.tsx
│   │   │   │   ├── AgentCard.tsx
│   │   │   │   ├── AgentStatus.tsx
│   │   │   │   └── DeployInstructions.tsx
│   │   │   ├── Execution/
│   │   │   │   ├── ExecutionMonitor.tsx
│   │   │   │   ├── LiveLogs.tsx
│   │   │   │   └── ProgressBar.tsx
│   │   │   ├── Reports/
│   │   │   │   ├── SecurityScore.tsx
│   │   │   │   ├── CoverageReport.tsx
│   │   │   │   ├── TechniqueReport.tsx
│   │   │   │   └── ExportPDF.tsx
│   │   │   └── common/
│   │   │       ├── Button.tsx
│   │   │       ├── Modal.tsx
│   │   │       ├── Card.tsx
│   │   │       └── Loading.tsx
│   │   │
│   │   ├── hooks/
│   │   │   ├── useWebSocket.ts
│   │   │   ├── useAgents.ts
│   │   │   ├── useScenarios.ts
│   │   │   └── useResults.ts
│   │   │
│   │   ├── services/
│   │   │   ├── api.ts
│   │   │   └── websocket.ts
│   │   │
│   │   ├── store/                      # State management
│   │   │   ├── index.ts
│   │   │   ├── agentSlice.ts
│   │   │   ├── scenarioSlice.ts
│   │   │   └── resultSlice.ts
│   │   │
│   │   ├── types/
│   │   │   ├── agent.ts
│   │   │   ├── scenario.ts
│   │   │   ├── technique.ts
│   │   │   └── mitre.ts
│   │   │
│   │   ├── utils/
│   │   │   ├── colors.ts
│   │   │   └── mitre.ts
│   │   │
│   │   ├── pages/
│   │   │   ├── Dashboard.tsx
│   │   │   ├── Matrix.tsx
│   │   │   ├── Scenarios.tsx
│   │   │   ├── Agents.tsx
│   │   │   ├── Execution.tsx
│   │   │   └── Reports.tsx
│   │   │
│   │   ├── App.tsx
│   │   ├── index.tsx
│   │   └── index.css
│   │
│   ├── public/
│   │   └── index.html
│   │
│   ├── package.json
│   ├── tsconfig.json
│   ├── tailwind.config.js
│   └── vite.config.ts
│
├── 📁 docs/                            # Documentation
│   ├── architecture/
│   │   ├── overview.md
│   │   ├── hexagonal.md
│   │   └── diagrams/
│   ├── api/
│   │   └── openapi.yaml
│   ├── techniques/
│   │   └── implementation-guide.md
│   └── deployment/
│       ├── docker.md
│       └── kubernetes.md
│
├── 📁 deployments/                     # Configs déploiement
│   ├── docker/
│   │   ├── Dockerfile.server
│   │   ├── Dockerfile.dashboard
│   │   └── docker-compose.yml
│   └── kubernetes/
│       ├── server-deployment.yaml
│       └── dashboard-deployment.yaml
│
├── 📁 scripts/                         # Scripts utilitaires
│   ├── build-agent.sh
│   ├── generate-certs.sh
│   └── import-mitre.sh
│
├── .gitignore
├── README.md
├── LICENSE
└── Makefile
```

---

## 4. Composants Détaillés

### 4.1 Control Server (Go)

#### 4.1.1 Responsabilités

| Responsabilité | Description |
|----------------|-------------|
| **Orchestration** | Planifier et coordonner l'exécution des techniques sur les agents |
| **Gestion Agents** | Enregistrer, monitorer et communiquer avec les agents |
| **Gestion Scénarios** | CRUD des scénarios et profils adversaires |
| **Collecte Résultats** | Recevoir et stocker les résultats d'exécution |
| **API REST** | Exposer les fonctionnalités au dashboard |
| **WebSocket** | Notifications temps réel |

#### 4.1.2 Entités du Domain

```go
// entity/agent.go
package entity

import "time"

type AgentStatus string

const (
    AgentOnline   AgentStatus = "online"
    AgentOffline  AgentStatus = "offline"
    AgentBusy     AgentStatus = "busy"
    AgentUntrusted AgentStatus = "untrusted"
)

type Agent struct {
    Paw         string            // Identifiant unique
    Hostname    string
    Platform    string            // "windows", "linux", "darwin"
    Username    string
    Executors   []string          // ["psh", "cmd", "bash"]
    Status      AgentStatus
    LastSeen    time.Time
    IPAddress   string
    OSVersion   string
    Metadata    map[string]string
}

func (a *Agent) IsCompatible(technique *Technique) bool {
    // Vérifie si l'agent peut exécuter cette technique
    for _, platform := range technique.Platforms {
        if platform == a.Platform {
            for _, executor := range technique.Executors {
                for _, agentExec := range a.Executors {
                    if executor.Type == agentExec {
                        return true
                    }
                }
            }
        }
    }
    return false
}
```

```go
// entity/technique.go
package entity

type TacticType string

const (
    Reconnaissance       TacticType = "reconnaissance"
    ResourceDevelopment  TacticType = "resource-development"
    InitialAccess        TacticType = "initial-access"
    Execution            TacticType = "execution"
    Persistence          TacticType = "persistence"
    PrivilegeEscalation  TacticType = "privilege-escalation"
    DefenseEvasion       TacticType = "defense-evasion"
    CredentialAccess     TacticType = "credential-access"
    Discovery            TacticType = "discovery"
    LateralMovement      TacticType = "lateral-movement"
    Collection           TacticType = "collection"
    CommandAndControl    TacticType = "command-and-control"
    Exfiltration         TacticType = "exfiltration"
    Impact               TacticType = "impact"
)

type Technique struct {
    ID          string      // "T1059.001"
    Name        string      // "PowerShell"
    Tactic      TacticType
    Description string
    Platforms   []string    // ["windows"]
    Executors   []Executor
    Detection   []Detection
    References  []string
    IsSafe      bool        // Ne cause pas de dommages
}

type Executor struct {
    Type    string // "psh", "cmd", "bash"
    Command string
    Cleanup string // Commande de nettoyage (optionnel)
    Timeout int    // Secondes
}

type Detection struct {
    Source    string // "Process Creation", "File Creation"
    Indicator string // Pattern de détection attendu
}
```

```go
// entity/scenario.go
package entity

import "time"

type Scenario struct {
    ID          string
    Name        string
    Description string
    Phases      []Phase
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type Phase struct {
    Name        string
    Description string
    Techniques  []string // IDs des techniques
    Order       int
}

type AdversaryProfile struct {
    ID          string
    Name        string   // "APT29 - Cozy Bear"
    Description string
    Techniques  []string
    References  []string
}
```

```go
// entity/result.go
package entity

import "time"

type ResultStatus string

const (
    StatusSuccess  ResultStatus = "success"   // Exécutée, non détectée
    StatusBlocked  ResultStatus = "blocked"   // Bloquée par défense
    StatusDetected ResultStatus = "detected"  // Exécutée mais alertée
    StatusFailed   ResultStatus = "failed"    // Erreur technique
    StatusSkipped  ResultStatus = "skipped"   // Non exécutée
)

type ExecutionResult struct {
    ID           string
    ExecutionID  string       // ID de l'exécution globale
    TechniqueID  string
    AgentPaw     string
    Status       ResultStatus
    Output       string       // Output encodé base64
    Stderr       string
    ExitCode     int
    DetectedBy   string       // "Windows Defender", "CrowdStrike"
    StartedAt    time.Time
    CompletedAt  time.Time
    Duration     time.Duration
}

type Execution struct {
    ID          string
    ScenarioID  string
    StartedAt   time.Time
    CompletedAt *time.Time
    Status      string // "running", "completed", "failed"
    Results     []ExecutionResult
}
```

#### 4.1.3 Domain Services

```go
// service/orchestrator.go
package service

type AttackOrchestrator struct {
    agentRepo    repository.AgentRepository
    techRepo     repository.TechniqueRepository
    validator    *TechniqueValidator
}

type ExecutionPlan struct {
    ID     string
    Tasks  []PlannedTask
}

type PlannedTask struct {
    TechniqueID string
    AgentPaw    string
    Phase       string
    Order       int
    Command     string
    Timeout     int
}

func (o *AttackOrchestrator) PlanExecution(
    scenario *entity.Scenario,
    targetAgents []*entity.Agent,
) (*ExecutionPlan, error) {
    
    plan := &ExecutionPlan{
        ID:    uuid.New().String(),
        Tasks: make([]PlannedTask, 0),
    }
    
    taskOrder := 0
    
    for _, phase := range scenario.Phases {
        for _, techID := range phase.Techniques {
            // Récupérer la technique
            technique, err := o.techRepo.FindByID(techID)
            if err != nil {
                continue
            }
            
            // Trouver un agent compatible
            var selectedAgent *entity.Agent
            for _, agent := range targetAgents {
                if agent.IsCompatible(technique) && agent.Status == entity.AgentOnline {
                    selectedAgent = agent
                    break
                }
            }
            
            if selectedAgent == nil {
                // Log: aucun agent compatible
                continue
            }
            
            // Sélectionner l'executor approprié
            executor := o.selectExecutor(technique, selectedAgent)
            
            plan.Tasks = append(plan.Tasks, PlannedTask{
                TechniqueID: techID,
                AgentPaw:    selectedAgent.Paw,
                Phase:       phase.Name,
                Order:       taskOrder,
                Command:     executor.Command,
                Timeout:     executor.Timeout,
            })
            
            taskOrder++
        }
    }
    
    return plan, nil
}

func (o *AttackOrchestrator) selectExecutor(
    tech *entity.Technique,
    agent *entity.Agent,
) *entity.Executor {
    // Priorité: psh > cmd > bash selon plateforme
    priority := map[string]int{"psh": 1, "cmd": 2, "bash": 3}
    
    var bestExecutor *entity.Executor
    bestPriority := 999
    
    for i, exec := range tech.Executors {
        for _, agentExec := range agent.Executors {
            if exec.Type == agentExec {
                if p, ok := priority[exec.Type]; ok && p < bestPriority {
                    bestPriority = p
                    bestExecutor = &tech.Executors[i]
                }
            }
        }
    }
    
    return bestExecutor
}
```

#### 4.1.4 API REST Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| **Agents** |
| GET | `/api/v1/agents` | Liste tous les agents |
| GET | `/api/v1/agents/:paw` | Détails d'un agent |
| DELETE | `/api/v1/agents/:paw` | Supprimer un agent |
| POST | `/api/v1/agents/:paw/beacon` | Endpoint beaconing (agents) |
| **Techniques** |
| GET | `/api/v1/techniques` | Liste toutes les techniques |
| GET | `/api/v1/techniques/:id` | Détails d'une technique |
| GET | `/api/v1/techniques/tactic/:tactic` | Techniques par tactic |
| **Scénarios** |
| GET | `/api/v1/scenarios` | Liste tous les scénarios |
| POST | `/api/v1/scenarios` | Créer un scénario |
| GET | `/api/v1/scenarios/:id` | Détails d'un scénario |
| PUT | `/api/v1/scenarios/:id` | Modifier un scénario |
| DELETE | `/api/v1/scenarios/:id` | Supprimer un scénario |
| **Exécutions** |
| POST | `/api/v1/executions` | Lancer une exécution |
| GET | `/api/v1/executions/:id` | Status d'une exécution |
| GET | `/api/v1/executions/:id/results` | Résultats d'une exécution |
| POST | `/api/v1/executions/:id/stop` | Arrêter une exécution |
| **Profils Adversaires** |
| GET | `/api/v1/adversaries` | Liste des profils |
| GET | `/api/v1/adversaries/:id` | Détails d'un profil |
| **Rapports** |
| GET | `/api/v1/reports/coverage` | Couverture MITRE |
| GET | `/api/v1/reports/score` | Score de sécurité |
| GET | `/api/v1/reports/export/:format` | Export PDF/JSON |
| **WebSocket** |
| WS | `/ws/live` | Notifications temps réel |

---

### 4.2 Agent (Rust)

#### 4.2.1 Responsabilités

| Responsabilité | Description |
|----------------|-------------|
| **Beaconing** | Contacter périodiquement le serveur pour récupérer des tâches |
| **Exécution** | Exécuter les techniques MITRE assignées |
| **Reporting** | Renvoyer les résultats au serveur |
| **Discrétion** | Rester léger et peu détectable |
| **Multi-plateforme** | Fonctionner sur Windows et Linux |

#### 4.2.2 Structure du Code

```rust
// src/main.rs
use std::time::Duration;
use tokio::time::sleep;

mod config;
mod beacon;
mod executor;
mod techniques;
mod utils;

use config::Config;
use beacon::BeaconClient;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Charger la configuration
    let config = Config::load()?;
    
    // Créer le client beacon
    let mut client = BeaconClient::new(&config)?;
    
    // Boucle principale
    loop {
        match client.beacon().await {
            Ok(instructions) => {
                for instruction in instructions {
                    let result = client.execute(&instruction).await;
                    client.report_result(result).await?;
                }
            }
            Err(e) => {
                eprintln!("Beacon error: {}", e);
            }
        }
        
        // Sleep avec jitter
        let jitter = rand::random::<u64>() % config.jitter;
        sleep(Duration::from_secs(config.sleep + jitter)).await;
    }
}
```

```rust
// src/config.rs
use serde::Deserialize;

#[derive(Deserialize)]
pub struct Config {
    pub server_url: String,
    pub sleep: u64,        // Intervalle en secondes
    pub jitter: u64,       // Variation aléatoire
    pub paw: Option<String>,
}

impl Config {
    pub fn load() -> Result<Self, Box<dyn std::error::Error>> {
        // Priorité: args CLI > env vars > fichier config > defaults
        Ok(Config {
            server_url: std::env::var("AUTOSTRIKE_SERVER")
                .unwrap_or_else(|_| "https://localhost:8443".to_string()),
            sleep: std::env::var("AUTOSTRIKE_SLEEP")
                .unwrap_or_else(|_| "5".to_string())
                .parse()
                .unwrap_or(5),
            jitter: std::env::var("AUTOSTRIKE_JITTER")
                .unwrap_or_else(|_| "3".to_string())
                .parse()
                .unwrap_or(3),
            paw: None,
        })
    }
}
```

```rust
// src/beacon/client.rs
use reqwest::Client;
use serde::{Deserialize, Serialize};

pub struct BeaconClient {
    client: Client,
    config: Config,
    paw: String,
}

#[derive(Serialize)]
struct BeaconRequest {
    paw: String,
    hostname: String,
    platform: String,
    username: String,
    executors: Vec<String>,
    results: Vec<TaskResult>,
}

#[derive(Deserialize)]
pub struct BeaconResponse {
    pub paw: String,
    pub sleep: u64,
    pub instructions: Vec<Instruction>,
}

#[derive(Deserialize)]
pub struct Instruction {
    pub id: String,
    pub technique_id: String,
    pub executor: String,
    pub command: String,
    pub timeout: u64,
}

#[derive(Serialize)]
pub struct TaskResult {
    pub id: String,
    pub technique_id: String,
    pub status: String,
    pub output: String,
    pub stderr: String,
    pub exit_code: i32,
    pub pid: u32,
}

impl BeaconClient {
    pub fn new(config: &Config) -> Result<Self, Box<dyn std::error::Error>> {
        let client = Client::builder()
            .danger_accept_invalid_certs(true) // Pour dev, à enlever en prod
            .build()?;
        
        let paw = config.paw.clone().unwrap_or_else(|| uuid::Uuid::new_v4().to_string());
        
        Ok(Self {
            client,
            config: config.clone(),
            paw,
        })
    }
    
    pub async fn beacon(&mut self) -> Result<Vec<Instruction>, Box<dyn std::error::Error>> {
        let request = BeaconRequest {
            paw: self.paw.clone(),
            hostname: hostname::get()?.to_string_lossy().to_string(),
            platform: std::env::consts::OS.to_string(),
            username: whoami::username(),
            executors: self.get_available_executors(),
            results: vec![],
        };
        
        let response: BeaconResponse = self.client
            .post(&format!("{}/api/v1/agents/{}/beacon", self.config.server_url, self.paw))
            .json(&request)
            .send()
            .await?
            .json()
            .await?;
        
        self.paw = response.paw;
        
        Ok(response.instructions)
    }
    
    fn get_available_executors(&self) -> Vec<String> {
        let mut executors = vec![];
        
        #[cfg(target_os = "windows")]
        {
            executors.push("psh".to_string());
            executors.push("cmd".to_string());
        }
        
        #[cfg(target_os = "linux")]
        {
            executors.push("bash".to_string());
            executors.push("sh".to_string());
        }
        
        #[cfg(target_os = "macos")]
        {
            executors.push("bash".to_string());
            executors.push("zsh".to_string());
        }
        
        executors
    }
}
```

```rust
// src/executor/mod.rs
pub mod powershell;
pub mod cmd;
pub mod bash;

use std::process::Output;
use std::time::Duration;

pub trait Executor {
    fn execute(&self, command: &str, timeout: Duration) -> Result<Output, ExecutorError>;
    fn name(&self) -> &str;
}

#[derive(Debug)]
pub enum ExecutorError {
    Timeout,
    ExecutionFailed(String),
    NotSupported,
}
```

```rust
// src/executor/powershell.rs
use std::process::{Command, Output, Stdio};
use std::time::Duration;
use super::{Executor, ExecutorError};

pub struct PowerShellExecutor;

impl Executor for PowerShellExecutor {
    fn execute(&self, command: &str, timeout: Duration) -> Result<Output, ExecutorError> {
        let output = Command::new("powershell.exe")
            .args(&[
                "-NoProfile",
                "-NonInteractive",
                "-ExecutionPolicy", "Bypass",
                "-Command", command
            ])
            .stdout(Stdio::piped())
            .stderr(Stdio::piped())
            .output()
            .map_err(|e| ExecutorError::ExecutionFailed(e.to_string()))?;
        
        Ok(output)
    }
    
    fn name(&self) -> &str {
        "psh"
    }
}
```

```rust
// src/techniques/discovery/t1082_system_info.rs
use crate::executor::Executor;
use crate::beacon::TaskResult;

pub struct T1082SystemInfo;

impl T1082SystemInfo {
    pub fn execute<E: Executor>(executor: &E, task_id: &str) -> TaskResult {
        let command = if cfg!(target_os = "windows") {
            "systeminfo; hostname; whoami /all"
        } else {
            "uname -a; hostname; id; cat /etc/os-release 2>/dev/null"
        };
        
        match executor.execute(command, std::time::Duration::from_secs(30)) {
            Ok(output) => {
                let stdout = String::from_utf8_lossy(&output.stdout).to_string();
                let stderr = String::from_utf8_lossy(&output.stderr).to_string();
                
                TaskResult {
                    id: task_id.to_string(),
                    technique_id: "T1082".to_string(),
                    status: if output.status.success() { "success" } else { "failed" }.to_string(),
                    output: base64::encode(&stdout),
                    stderr: base64::encode(&stderr),
                    exit_code: output.status.code().unwrap_or(-1),
                    pid: std::process::id(),
                }
            }
            Err(e) => {
                TaskResult {
                    id: task_id.to_string(),
                    technique_id: "T1082".to_string(),
                    status: "failed".to_string(),
                    output: String::new(),
                    stderr: format!("{:?}", e),
                    exit_code: -1,
                    pid: std::process::id(),
                }
            }
        }
    }
}
```

#### 4.2.3 Compilation Multi-Plateforme

```toml
# Cargo.toml
[package]
name = "autostrike-agent"
version = "0.1.0"
edition = "2021"

[dependencies]
tokio = { version = "1", features = ["full"] }
reqwest = { version = "0.11", features = ["json", "rustls-tls"] }
serde = { version = "1", features = ["derive"] }
serde_json = "1"
uuid = { version = "1", features = ["v4"] }
base64 = "0.21"
hostname = "0.3"
whoami = "1"
rand = "0.8"

[target.'cfg(windows)'.dependencies]
windows = { version = "0.48", features = ["Win32_System_Threading"] }

[profile.release]
opt-level = "z"     # Optimiser taille
lto = true          # Link Time Optimization
codegen-units = 1   # Meilleure optimisation
strip = true        # Strip symbols
```

```bash
# scripts/build-agent.sh
#!/bin/bash

# Windows x64
cargo build --release --target x86_64-pc-windows-gnu

# Linux x64
cargo build --release --target x86_64-unknown-linux-gnu

# Linux ARM64
cargo build --release --target aarch64-unknown-linux-gnu

echo "Agents compilés dans target/*/release/"
```

---

### 4.3 Dashboard (React)

#### 4.3.1 Composants Principaux

##### Matrice MITRE ATT&CK

```tsx
// src/components/AttackMatrix/MitreMatrix.tsx
import React, { useMemo } from 'react';
import * as d3 from 'd3';
import { TacticColumn } from './TacticColumn';
import { HeatmapLegend } from './HeatmapLegend';
import { Technique, ExecutionResult } from '../../types';

interface MitreMatrixProps {
  techniques: Technique[];
  results: ExecutionResult[];
  onTechniqueClick: (technique: Technique) => void;
}

const TACTICS_ORDER = [
  'reconnaissance',
  'resource-development',
  'initial-access',
  'execution',
  'persistence',
  'privilege-escalation',
  'defense-evasion',
  'credential-access',
  'discovery',
  'lateral-movement',
  'collection',
  'command-and-control',
  'exfiltration',
  'impact',
];

export const MitreMatrix: React.FC<MitreMatrixProps> = ({
  techniques,
  results,
  onTechniqueClick,
}) => {
  // Grouper techniques par tactic
  const techniquesByTactic = useMemo(() => {
    const grouped: Record<string, Technique[]> = {};
    
    TACTICS_ORDER.forEach(tactic => {
      grouped[tactic] = techniques.filter(t => t.tactic === tactic);
    });
    
    return grouped;
  }, [techniques]);
  
  // Calculer les couleurs basées sur les résultats
  const colorScale = useMemo(() => {
    return d3.scaleOrdinal<string>()
      .domain(['success', 'blocked', 'detected', 'untested'])
      .range(['#ef4444', '#22c55e', '#f59e0b', '#6b7280']);
  }, []);
  
  // Mapper les résultats par technique
  const resultsByTechnique = useMemo(() => {
    const map: Record<string, ExecutionResult[]> = {};
    results.forEach(r => {
      if (!map[r.techniqueId]) map[r.techniqueId] = [];
      map[r.techniqueId].push(r);
    });
    return map;
  }, [results]);
  
  const getTechniqueStatus = (techniqueId: string): string => {
    const techResults = resultsByTechnique[techniqueId];
    if (!techResults || techResults.length === 0) return 'untested';
    
    // Priorité: success > detected > blocked
    if (techResults.some(r => r.status === 'success')) return 'success';
    if (techResults.some(r => r.status === 'detected')) return 'detected';
    if (techResults.some(r => r.status === 'blocked')) return 'blocked';
    
    return 'untested';
  };
  
  return (
    <div className="overflow-x-auto">
      <div className="flex gap-1 min-w-max p-4">
        {TACTICS_ORDER.map(tactic => (
          <TacticColumn
            key={tactic}
            tactic={tactic}
            techniques={techniquesByTactic[tactic] || []}
            getStatus={getTechniqueStatus}
            colorScale={colorScale}
            onTechniqueClick={onTechniqueClick}
          />
        ))}
      </div>
      <HeatmapLegend colorScale={colorScale} />
    </div>
  );
};
```

```tsx
// src/components/AttackMatrix/TechniqueCell.tsx
import React from 'react';
import { Technique } from '../../types';

interface TechniqueCellProps {
  technique: Technique;
  status: string;
  color: string;
  onClick: () => void;
}

export const TechniqueCell: React.FC<TechniqueCellProps> = ({
  technique,
  status,
  color,
  onClick,
}) => {
  return (
    <div
      className="p-2 rounded cursor-pointer transition-all hover:scale-105 hover:shadow-lg"
      style={{ backgroundColor: color }}
      onClick={onClick}
      title={`${technique.id}: ${technique.name}\nStatus: ${status}`}
    >
      <div className="text-xs font-mono text-white opacity-75">
        {technique.id}
      </div>
      <div className="text-sm text-white font-medium truncate">
        {technique.name}
      </div>
    </div>
  );
};
```

##### Score de Sécurité

```tsx
// src/components/Reports/SecurityScore.tsx
import React from 'react';
import { useMemo } from 'react';
import { ExecutionResult } from '../../types';

interface SecurityScoreProps {
  results: ExecutionResult[];
  totalTechniques: number;
}

export const SecurityScore: React.FC<SecurityScoreProps> = ({
  results,
  totalTechniques,
}) => {
  const score = useMemo(() => {
    if (results.length === 0) return 0;
    
    const blocked = results.filter(r => r.status === 'blocked').length;
    const detected = results.filter(r => r.status === 'detected').length;
    const success = results.filter(r => r.status === 'success').length;
    
    // Blocked = 100%, Detected = 50%, Success = 0%
    const points = (blocked * 100) + (detected * 50);
    const maxPoints = results.length * 100;
    
    return Math.round((points / maxPoints) * 100);
  }, [results]);
  
  const getScoreColor = (score: number): string => {
    if (score >= 80) return 'text-green-500';
    if (score >= 60) return 'text-yellow-500';
    if (score >= 40) return 'text-orange-500';
    return 'text-red-500';
  };
  
  const getScoreLabel = (score: number): string => {
    if (score >= 80) return 'Excellent';
    if (score >= 60) return 'Bon';
    if (score >= 40) return 'À améliorer';
    return 'Critique';
  };
  
  return (
    <div className="bg-white rounded-xl shadow-lg p-6">
      <h2 className="text-lg font-semibold text-gray-700 mb-4">
        Score de Sécurité
      </h2>
      
      <div className="flex items-center justify-center">
        <div className="relative">
          {/* Cercle de progression */}
          <svg className="w-32 h-32 transform -rotate-90">
            <circle
              cx="64"
              cy="64"
              r="56"
              stroke="#e5e7eb"
              strokeWidth="12"
              fill="none"
            />
            <circle
              cx="64"
              cy="64"
              r="56"
              stroke="currentColor"
              strokeWidth="12"
              fill="none"
              className={getScoreColor(score)}
              strokeDasharray={`${(score / 100) * 352} 352`}
              strokeLinecap="round"
            />
          </svg>
          
          {/* Score au centre */}
          <div className="absolute inset-0 flex flex-col items-center justify-center">
            <span className={`text-3xl font-bold ${getScoreColor(score)}`}>
              {score}%
            </span>
            <span className="text-sm text-gray-500">
              {getScoreLabel(score)}
            </span>
          </div>
        </div>
      </div>
      
      {/* Détails */}
      <div className="mt-6 grid grid-cols-3 gap-4 text-center">
        <div>
          <div className="text-2xl font-bold text-green-500">
            {results.filter(r => r.status === 'blocked').length}
          </div>
          <div className="text-xs text-gray-500">Bloquées</div>
        </div>
        <div>
          <div className="text-2xl font-bold text-yellow-500">
            {results.filter(r => r.status === 'detected').length}
          </div>
          <div className="text-xs text-gray-500">Détectées</div>
        </div>
        <div>
          <div className="text-2xl font-bold text-red-500">
            {results.filter(r => r.status === 'success').length}
          </div>
          <div className="text-xs text-gray-500">Réussies</div>
        </div>
      </div>
      
      <div className="mt-4 text-center text-sm text-gray-400">
        {results.length} / {totalTechniques} techniques testées
      </div>
    </div>
  );
};
```

##### Monitoring Temps Réel

```tsx
// src/components/Execution/ExecutionMonitor.tsx
import React, { useEffect, useState } from 'react';
import { useWebSocket } from '../../hooks/useWebSocket';
import { LiveLogs } from './LiveLogs';
import { ProgressBar } from './ProgressBar';

interface ExecutionMonitorProps {
  executionId: string;
}

interface ExecutionEvent {
  type: 'started' | 'technique_complete' | 'completed' | 'error';
  techniqueId?: string;
  status?: string;
  output?: string;
  timestamp: string;
}

export const ExecutionMonitor: React.FC<ExecutionMonitorProps> = ({
  executionId,
}) => {
  const [events, setEvents] = useState<ExecutionEvent[]>([]);
  const [progress, setProgress] = useState({ current: 0, total: 0 });
  
  const { lastMessage, connectionStatus } = useWebSocket(
    `wss://localhost:8443/ws/live?execution=${executionId}`
  );
  
  useEffect(() => {
    if (lastMessage) {
      const event: ExecutionEvent = JSON.parse(lastMessage.data);
      setEvents(prev => [...prev, event]);
      
      if (event.type === 'technique_complete') {
        setProgress(prev => ({ ...prev, current: prev.current + 1 }));
      }
      
      if (event.type === 'started' && event.total) {
        setProgress({ current: 0, total: event.total });
      }
    }
  }, [lastMessage]);
  
  return (
    <div className="bg-white rounded-xl shadow-lg p-6">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold text-gray-700">
          Exécution en cours
        </h2>
        <span className={`px-2 py-1 rounded text-sm ${
          connectionStatus === 'connected' 
            ? 'bg-green-100 text-green-700'
            : 'bg-red-100 text-red-700'
        }`}>
          {connectionStatus === 'connected' ? '● Live' : '○ Déconnecté'}
        </span>
      </div>
      
      <ProgressBar 
        current={progress.current} 
        total={progress.total} 
      />
      
      <LiveLogs events={events} />
    </div>
  );
};
```

---

## 5. Stack Technologique

### 5.1 Tableau Récapitulatif

| Composant | Technologie | Version | Justification |
|-----------|-------------|---------|---------------|
| **Server** | Go | 1.21+ | Performance, concurrence, simplicité |
| **Server Framework** | Gin | 1.9+ | Rapide, bien documenté |
| **Agent** | Rust | 1.75+ | Sécurité mémoire, performance, évasion AV |
| **Dashboard** | React | 18+ | Écosystème riche, composants réutilisables |
| **Dashboard Build** | Vite | 5+ | Bundler rapide |
| **TypeScript** | TypeScript | 5+ | Typage statique |
| **Styling** | TailwindCSS | 3+ | Utility-first, rapide à développer |
| **Visualisation** | D3.js | 7+ | Flexibilité pour la matrice ATT&CK |
| **State Management** | Zustand | 4+ | Simple, léger |
| **BDD** | SQLite | 3+ | MVP simple, puis PostgreSQL |
| **Communication** | REST + WebSocket | - | API + temps réel |
| **Sécurité** | mTLS | TLS 1.3 | Authentification mutuelle |
| **Config Techniques** | YAML | - | Lisible, facile à maintenir |
| **Container** | Docker | 24+ | Déploiement standardisé |

### 5.2 Dépendances Go (Server)

```go
// go.mod
module autostrike

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/gorilla/websocket v1.5.0
    github.com/google/uuid v1.4.0
    github.com/mattn/go-sqlite3 v1.14.18
    github.com/spf13/viper v1.17.0
    gopkg.in/yaml.v3 v3.0.1
    go.uber.org/zap v1.26.0
)
```

### 5.3 Dépendances React (Dashboard)

```json
// package.json
{
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-router-dom": "^6.18.0",
    "d3": "^7.8.5",
    "zustand": "^4.4.6",
    "axios": "^1.6.0",
    "@tanstack/react-query": "^5.8.0"
  },
  "devDependencies": {
    "typescript": "^5.2.2",
    "vite": "^5.0.0",
    "@types/react": "^18.2.37",
    "@types/d3": "^7.4.2",
    "tailwindcss": "^3.3.5",
    "autoprefixer": "^10.4.16",
    "postcss": "^8.4.31"
  }
}
```

---

## 6. Modèle de Données

### 6.1 Schéma SQLite

```sql
-- migrations/001_initial.sql

-- Agents
CREATE TABLE agents (
    paw TEXT PRIMARY KEY,
    hostname TEXT NOT NULL,
    platform TEXT NOT NULL,
    username TEXT,
    ip_address TEXT,
    os_version TEXT,
    executors TEXT, -- JSON array
    status TEXT DEFAULT 'offline',
    last_seen DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_agents_status ON agents(status);

-- Techniques
CREATE TABLE techniques (
    id TEXT PRIMARY KEY,        -- "T1059.001"
    name TEXT NOT NULL,
    tactic TEXT NOT NULL,
    description TEXT,
    platforms TEXT,             -- JSON array
    executors TEXT,             -- JSON array
    detection TEXT,             -- JSON array
    is_safe BOOLEAN DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_techniques_tactic ON techniques(tactic);

-- Scénarios
CREATE TABLE scenarios (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    phases TEXT NOT NULL,       -- JSON array
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Profils Adversaires
CREATE TABLE adversary_profiles (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    techniques TEXT NOT NULL,   -- JSON array
    references TEXT,            -- JSON array
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Exécutions
CREATE TABLE executions (
    id TEXT PRIMARY KEY,
    scenario_id TEXT REFERENCES scenarios(id),
    status TEXT DEFAULT 'pending',
    started_at DATETIME,
    completed_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_executions_status ON executions(status);

-- Résultats
CREATE TABLE execution_results (
    id TEXT PRIMARY KEY,
    execution_id TEXT REFERENCES executions(id),
    technique_id TEXT REFERENCES techniques(id),
    agent_paw TEXT REFERENCES agents(paw),
    status TEXT NOT NULL,       -- success, blocked, detected, failed
    output TEXT,                -- base64
    stderr TEXT,
    exit_code INTEGER,
    detected_by TEXT,
    started_at DATETIME,
    completed_at DATETIME,
    duration_ms INTEGER
);

CREATE INDEX idx_results_execution ON execution_results(execution_id);
CREATE INDEX idx_results_technique ON execution_results(technique_id);
CREATE INDEX idx_results_status ON execution_results(status);

-- Tâches en attente
CREATE TABLE pending_tasks (
    id TEXT PRIMARY KEY,
    agent_paw TEXT REFERENCES agents(paw),
    execution_id TEXT REFERENCES executions(id),
    technique_id TEXT REFERENCES techniques(id),
    command TEXT NOT NULL,
    timeout INTEGER DEFAULT 30,
    status TEXT DEFAULT 'pending',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_pending_agent ON pending_tasks(agent_paw, status);
```

### 6.2 Structures JSON

#### Technique YAML

```yaml
# configs/techniques/T1059.001.yaml
id: "T1059.001"
name: "PowerShell"
tactic: "execution"
description: |
  Adversaries may abuse PowerShell commands and scripts for execution.
  PowerShell is a powerful interactive command-line interface and 
  scripting environment included in the Windows operating system.

platforms:
  - windows

executors:
  - type: psh
    command: |
      $ExecutionContext.SessionState.LanguageMode
      Get-Process | Select-Object -First 5
    cleanup: null
    timeout: 30
    
  - type: cmd
    command: |
      powershell.exe -NoProfile -Command "Get-Process | Select -First 5"
    cleanup: null
    timeout: 30

detection:
  - source: "Process Creation"
    indicator: "powershell.exe with suspicious arguments"
  - source: "Script Block Logging"
    indicator: "Event ID 4104"

is_safe: true

references:
  - https://attack.mitre.org/techniques/T1059/001/
```

#### Scénario JSON

```json
{
  "id": "scenario-001",
  "name": "APT29 - Initial Compromise",
  "description": "Emulation des TTPs APT29 pour la phase initiale",
  "phases": [
    {
      "name": "Reconnaissance",
      "description": "Collecte d'informations système",
      "order": 1,
      "techniques": ["T1082", "T1083", "T1057"]
    },
    {
      "name": "Execution",
      "description": "Exécution de commandes",
      "order": 2,
      "techniques": ["T1059.001", "T1059.003"]
    },
    {
      "name": "Persistence",
      "description": "Installation de persistance",
      "order": 3,
      "techniques": ["T1053.005", "T1547.001"]
    }
  ],
  "created_at": "2026-09-15T10:00:00Z"  // Date exemple
}
```

---

## 7. Protocoles de Communication

### 7.1 Protocole Beaconing Agent ↔ Server

```
┌─────────────────┐                          ┌─────────────────┐
│                 │                          │                 │
│     AGENT       │                          │     SERVER      │
│                 │                          │                 │
└────────┬────────┘                          └────────┬────────┘
         │                                            │
         │  1. POST /api/v1/agents/{paw}/beacon       │
         │  ─────────────────────────────────────────►│
         │  {                                         │
         │    "paw": "agent-uuid",                    │
         │    "hostname": "PC-001",                   │
         │    "platform": "windows",                  │
         │    "executors": ["psh", "cmd"],            │
         │    "results": []                           │
         │  }                                         │
         │                                            │
         │  2. Response avec instructions             │
         │  ◄─────────────────────────────────────────│
         │  {                                         │
         │    "paw": "agent-uuid",                    │
         │    "sleep": 5,                             │
         │    "instructions": [                       │
         │      {                                     │
         │        "id": "task-001",                   │
         │        "technique_id": "T1082",            │
         │        "executor": "psh",                  │
         │        "command": "systeminfo",            │
         │        "timeout": 30                       │
         │      }                                     │
         │    ]                                       │
         │  }                                         │
         │                                            │
         │  [Agent exécute les instructions]          │
         │                                            │
         │  3. POST /api/v1/agents/{paw}/beacon       │
         │  ─────────────────────────────────────────►│
         │  {                                         │
         │    "paw": "agent-uuid",                    │
         │    ...                                     │
         │    "results": [                            │
         │      {                                     │
         │        "id": "task-001",                   │
         │        "technique_id": "T1082",            │
         │        "status": "success",                │
         │        "output": "base64...",              │
         │        "exit_code": 0                      │
         │      }                                     │
         │    ]                                       │
         │  }                                         │
         │                                            │
         │  [Sleep avec jitter]                       │
         │                                            │
         │  4. Répéter...                             │
         │                                            │
```

### 7.2 WebSocket Events (Server → Dashboard)

```typescript
// Types d'événements WebSocket
interface WSEvent {
  type: EventType;
  timestamp: string;
  payload: unknown;
}

type EventType = 
  | 'agent_connected'
  | 'agent_disconnected'
  | 'execution_started'
  | 'technique_started'
  | 'technique_completed'
  | 'execution_completed'
  | 'error';

// Exemples d'événements
{
  "type": "agent_connected",
  "timestamp": "2026-01-29T10:30:00Z",
  "payload": {
    "paw": "agent-001",
    "hostname": "PC-TARGET-01",
    "platform": "windows"
  }
}

{
  "type": "technique_completed",
  "timestamp": "2026-01-29T10:30:15Z",
  "payload": {
    "execution_id": "exec-001",
    "technique_id": "T1082",
    "agent_paw": "agent-001",
    "status": "success",
    "duration_ms": 1523
  }
}
```

### 7.3 Sécurité des Communications

```
┌──────────────────────────────────────────────────────────────────┐
│                        SÉCURITÉ mTLS                             │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────┐           TLS 1.3           ┌─────────────┐    │
│  │             │  ◄──────────────────────►   │             │    │
│  │   AGENT     │                             │   SERVER    │    │
│  │             │                             │             │    │
│  └──────┬──────┘                             └──────┬──────┘    │
│         │                                           │           │
│         │  Certificat Client                        │           │
│         │  (généré lors du déploiement)             │           │
│         │                                           │           │
│         │  Certificat Serveur                       │           │
│         │  (CA interne AutoStrike)                  │           │
│                                                                  │
├──────────────────────────────────────────────────────────────────┤
│  Génération des certificats:                                     │
│                                                                  │
│  1. CA racine AutoStrike (auto-signée)                          │
│  2. Certificat serveur signé par CA                             │
│  3. Certificat agent signé par CA (unique par agent)            │
│                                                                  │
│  Le serveur vérifie:                                             │
│  - Certificat client valide                                      │
│  - Signé par la CA AutoStrike                                   │
│  - Non révoqué                                                   │
│                                                                  │
│  L'agent vérifie:                                                │
│  - Certificat serveur valide                                     │
│  - Signé par la CA AutoStrike                                   │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

---

## 8. Techniques MITRE ATT&CK

### 8.1 Techniques Prioritaires (MVP)

| ID | Nom | Tactic | Complexité | Priorité |
|----|-----|--------|------------|----------|
| **Discovery** |
| T1082 | System Information Discovery | Discovery | Faible | P0 |
| T1083 | File and Directory Discovery | Discovery | Faible | P0 |
| T1057 | Process Discovery | Discovery | Faible | P0 |
| T1016 | System Network Configuration | Discovery | Faible | P1 |
| T1069 | Permission Groups Discovery | Discovery | Faible | P1 |
| T1087 | Account Discovery | Discovery | Faible | P1 |
| T1018 | Remote System Discovery | Discovery | Faible | P1 |
| T1007 | System Service Discovery | Discovery | Faible | P1 |
| T1049 | System Network Connections | Discovery | Faible | P1 |
| **Execution** |
| T1059.001 | PowerShell | Execution | Faible | P0 |
| T1059.003 | Windows Command Shell | Execution | Faible | P0 |
| T1059.004 | Unix Shell | Execution | Faible | P0 |
| **Persistence** |
| T1053.005 | Scheduled Task | Persistence | Moyenne | P1 |
| T1547.001 | Registry Run Keys | Persistence | Moyenne | P1 |
| T1136.001 | Local Account | Persistence | Moyenne | P2 |
| **Defense Evasion** |
| T1070.004 | File Deletion | Defense Evasion | Faible | P1 |
| T1027 | Obfuscated Files | Defense Evasion | Moyenne | P2 |
| **Credential Access** |
| T1003.001 | LSASS Memory (simulation) | Credential Access | Haute | P2 |

### 8.2 Implémentation Type

```yaml
# T1082 - System Information Discovery
id: "T1082"
name: "System Information Discovery"
tactic: "discovery"
description: |
  An adversary may attempt to get detailed information about the 
  operating system and hardware.

platforms:
  - windows
  - linux
  - macos

executors:
  # Windows - PowerShell
  - type: psh
    command: |
      systeminfo
      Get-ComputerInfo | Select-Object WindowsVersion, OsHardwareAbstractionLayer
      Get-WmiObject Win32_OperatingSystem | Select-Object Caption, Version, BuildNumber
    timeout: 60
    
  # Windows - CMD
  - type: cmd
    command: |
      systeminfo
      hostname
      ver
    timeout: 60
    
  # Linux
  - type: bash
    command: |
      uname -a
      cat /etc/os-release
      hostnamectl
      lscpu
    timeout: 30
    
  # macOS
  - type: zsh
    command: |
      uname -a
      sw_vers
      system_profiler SPHardwareDataType
    timeout: 30

detection:
  - source: "Process Creation"
    indicator: "systeminfo.exe execution"
  - source: "Command Line"
    indicator: "Get-ComputerInfo or Get-WmiObject"

is_safe: true

references:
  - https://attack.mitre.org/techniques/T1082/
```

### 8.3 Scénarios Prédéfinis

#### Scénario 1: Reconnaissance Basique

```yaml
name: "Basic Reconnaissance"
description: "Collecte d'informations système standard"
phases:
  - name: "System Discovery"
    techniques:
      - T1082  # System Information
      - T1083  # File Discovery
      - T1057  # Process Discovery
      - T1016  # Network Configuration
```

#### Scénario 2: APT29 (Cozy Bear)

```yaml
name: "APT29 - Cozy Bear Emulation"
description: "Émulation des TTPs du groupe APT29"
phases:
  - name: "Initial Reconnaissance"
    techniques:
      - T1082
      - T1083
      - T1057
      
  - name: "Execution"
    techniques:
      - T1059.001  # PowerShell
      
  - name: "Persistence"
    techniques:
      - T1053.005  # Scheduled Task
      - T1547.001  # Registry Run Keys
      
  - name: "Defense Evasion"
    techniques:
      - T1070.004  # File Deletion
```

#### Scénario 3: Ransomware Simulation

```yaml
name: "Ransomware Behavior"
description: "Simulation du comportement typique d'un ransomware (safe)"
phases:
  - name: "Discovery"
    techniques:
      - T1082
      - T1083
      - T1135  # Network Share Discovery
      
  - name: "Collection"
    techniques:
      - T1005  # Data from Local System (list only)
      
  - name: "Impact Simulation"
    techniques:
      - T1486  # Data Encrypted (simulation - crée fichiers test)
```

### 8.4 Alignement EBIOS RM (Méthode ANSSI)

AutoStrike s'aligne sur **deux frameworks complémentaires** pour une couverture internationale et française :

| Framework | Origine | Usage | Granularité |
|-----------|---------|-------|-------------|
| **MITRE ATT&CK** | USA (MITRE Corp) | Standard international | Techniques détaillées (T1082, T1059...) |
| **EBIOS RM** | France (ANSSI) | Conformité française | Phases d'attaque (CRTE) |

#### Méthodologie CRTE

La méthode EBIOS Risk Manager de l'ANSSI structure les scénarios d'attaque en 4 phases **CRTE** :

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  CONNAÎTRE  │───>│   RENTRER   │───>│   TROUVER   │───>│  EXPLOITER  │
│             │    │             │    │             │    │             │
│ Reconnaissance│   │ Accès initial│   │ Exploration │    │ Exécution   │
│ de la cible │    │ au système  │    │ interne     │    │ de l'attaque│
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
```

#### Mapping CRTE ↔ MITRE ATT&CK

| Phase CRTE | Description | Tactiques MITRE | Couverture AutoStrike |
|------------|-------------|-----------------|----------------------|
| **C - Connaitre** | Reconnaissance de la cible | Reconnaissance, Discovery | ✅ 11 techniques (2 Recon + 9 Discovery) |
| **R - Rentrer** | Acces initial au systeme | Initial Access, Privilege Escalation | ✅ 7 techniques (3 IA + 4 PrivEsc) |
| **T - Trouver** | Exploration interne, mouvement lateral | Lateral Movement, Credential Access, Collection | ✅ 11 techniques (3 LM + 4 Cred + 4 Coll) |
| **E - Exploiter** | Execution de l'attaque, impact | Execution, Persistence, Defense Evasion, C2, Exfiltration, Impact | ✅ 19 techniques |

**Couverture actuelle :** 4/4 phases CRTE couvertes (48 techniques, 13 tactiques)

#### Pourquoi ce double alignement ?

| Contexte | Framework privilégié |
|----------|---------------------|
| **Client international** | MITRE ATT&CK (standard mondial) |
| **Client français / ANSSI** | EBIOS RM + MITRE ATT&CK |
| **Audit de conformité** | EBIOS RM (méthode homologuée) |
| **Équipe SOC technique** | MITRE ATT&CK (granularité) |

> **Argument clé :** "AutoStrike parle les deux langages - MITRE pour la technique, EBIOS RM pour la conformité française."

---

## 9. Interface Utilisateur

### 9.1 Wireframes

#### Dashboard Principal

```
┌──────────────────────────────────────────────────────────────────────────┐
│  🎯 AutoStrike                                    [User] [Settings] [?]  │
├────────────┬─────────────────────────────────────────────────────────────┤
│            │                                                             │
│  Dashboard │   ┌─────────────────────────────────────────────────────┐   │
│            │   │              Security Score                         │   │
│  ○ Matrix  │   │                                                     │   │
│            │   │                   72%                                │   │
│  ○ Agents  │   │              ┌─────────┐                            │   │
│            │   │              │   😐    │                            │   │
│  ○ Scenar. │   │              └─────────┘                            │   │
│            │   │           "Needs Improvement"                       │   │
│  ○ Reports │   │                                                     │   │
│            │   │  Blocked: 45  Detected: 12  Success: 8              │   │
│            │   └─────────────────────────────────────────────────────┘   │
│            │                                                             │
│            │   ┌──────────────────────┐  ┌──────────────────────┐       │
│            │   │ Active Agents: 5     │  │ Last Execution       │       │
│            │   │ ● PC-001 (Windows)   │  │ APT29 Simulation     │       │
│            │   │ ● PC-002 (Windows)   │  │ 2h ago - 65 techniques│       │
│            │   │ ● SRV-01 (Linux)     │  │ Score: 72%           │       │
│            │   │ ○ PC-003 (Offline)   │  │ [View Details]       │       │
│            │   └──────────────────────┘  └──────────────────────┘       │
│            │                                                             │
└────────────┴─────────────────────────────────────────────────────────────┘
```

#### Matrice ATT&CK

```
┌──────────────────────────────────────────────────────────────────────────┐
│  🎯 AutoStrike > MITRE ATT&CK Matrix          [Filter ▼] [Export]       │
├────────────┬─────────────────────────────────────────────────────────────┤
│            │                                                             │
│  Dashboard │   Legend: 🟢 Blocked  🟡 Detected  🔴 Success  ⚪ Untested  │
│            │                                                             │
│  ● Matrix  │   ┌────────┬────────┬────────┬────────┬────────┬────────┐  │
│            │   │ Recon  │ Exec   │ Persist│ Priv   │ Defense│ Discov │  │
│  ○ Agents  │   ├────────┼────────┼────────┼────────┼────────┼────────┤  │
│            │   │🟢T1595 │🔴T1059 │🟢T1053 │⚪T1548 │🟡T1070 │🟢T1082 │  │
│  ○ Scenar. │   │        │.001    │.005    │        │.004    │        │  │
│            │   ├────────┼────────┼────────┼────────┼────────┼────────┤  │
│  ○ Reports │   │⚪T1592 │🟢T1059 │🔴T1547 │⚪T1134 │⚪T1027 │🟢T1083 │  │
│            │   │        │.003    │.001    │        │        │        │  │
│            │   ├────────┼────────┼────────┼────────┼────────┼────────┤  │
│            │   │⚪T1589 │⚪T1106 │⚪T1136 │⚪T1068 │⚪T1140 │🟡T1057 │  │
│            │   │        │        │        │        │        │        │  │
│            │   └────────┴────────┴────────┴────────┴────────┴────────┘  │
│            │                                                             │
│            │   [Click technique for details]                             │
│            │                                                             │
└────────────┴─────────────────────────────────────────────────────────────┘
```

#### Détail Technique

```
┌──────────────────────────────────────────────────────────────────────────┐
│  T1059.001 - PowerShell                                          [X]    │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  Tactic: Execution                    Status: 🔴 SUCCESS                 │
│  Platforms: Windows                   Last Tested: 2h ago               │
│                                                                          │
│  ┌────────────────────────────────────────────────────────────────────┐  │
│  │ Description                                                        │  │
│  │                                                                    │  │
│  │ Adversaries may abuse PowerShell commands and scripts for         │  │
│  │ execution. PowerShell is a powerful interactive command-line      │  │
│  │ interface and scripting environment.                              │  │
│  └────────────────────────────────────────────────────────────────────┘  │
│                                                                          │
│  ┌────────────────────────────────────────────────────────────────────┐  │
│  │ Execution History                                                  │  │
│  │                                                                    │  │
│  │  Date       │ Agent   │ Status  │ Duration │ Detected By          │  │
│  │  ─────────────────────────────────────────────────────────────     │  │
│  │  Jan 29     │ PC-001  │ 🔴      │ 1.2s     │ -                    │  │
│  │  Jan 28     │ PC-002  │ 🟡      │ 0.8s     │ Windows Defender     │  │
│  │  Jan 27     │ PC-001  │ 🟢      │ -        │ CrowdStrike          │  │
│  └────────────────────────────────────────────────────────────────────┘  │
│                                                                          │
│  ┌────────────────────────────────────────────────────────────────────┐  │
│  │ Detection Guidance                                                 │  │
│  │                                                                    │  │
│  │ • Monitor for powershell.exe with suspicious arguments            │  │
│  │ • Enable Script Block Logging (Event ID 4104)                     │  │
│  │ • Watch for encoded commands (-enc, -e)                           │  │
│  └────────────────────────────────────────────────────────────────────┘  │
│                                                                          │
│  [Run on Selected Agents]  [Add to Scenario]  [View on MITRE]           │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

### 9.2 Palette de Couleurs

```css
/* Couleurs principales */
--color-success: #22c55e;    /* Blocked - Vert */
--color-warning: #f59e0b;    /* Detected - Orange */
--color-danger: #ef4444;     /* Success (attaque) - Rouge */
--color-neutral: #6b7280;    /* Untested - Gris */

/* Couleurs secondaires */
--color-primary: #3b82f6;    /* Bleu principal */
--color-secondary: #8b5cf6;  /* Violet accent */

/* Backgrounds */
--bg-dark: #1f2937;
--bg-card: #ffffff;
--bg-hover: #f3f4f6;

/* Texte */
--text-primary: #111827;
--text-secondary: #6b7280;
```

---

## 10. Sécurité

### 10.1 Considérations de Sécurité

| Risque | Mitigation |
|--------|------------|
| **Agent malveillant** | Authentification mTLS, validation certificat |
| **Interception trafic** | TLS 1.3, certificate pinning |
| **Exécution non autorisée** | Mode "safe" par défaut, techniques validées |
| **Accès dashboard non autorisé** | Authentification, RBAC |
| **Stockage credentials** | Pas de stockage, tokens temporaires |
| **Techniques destructives** | Flag `is_safe`, cleanup automatique |

### 10.2 Mode Safe

Toutes les techniques sont classifiées :

```yaml
# Technique SAFE - Peut s'exécuter en production
is_safe: true
# Actions: lecture seule, pas de modification système

# Technique UNSAFE - Environnement lab uniquement  
is_safe: false
# Actions: modification registre, création fichiers, etc.
```

### 10.3 Authentification

```
┌─────────────────────────────────────────────────────────────────┐
│                    Authentification                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Dashboard → Server:                                             │
│  • JWT avec expiration courte (15 min)                          │
│  • Refresh token (7 jours)                                      │
│  • HTTPS obligatoire                                            │
│                                                                  │
│  Agent → Server:                                                 │
│  • mTLS (certificat client unique par agent)                    │
│  • Certificat révocable                                         │
│  • IP whitelist optionnel                                       │
│                                                                  │
│  Roles (Dashboard):                                              │
│  • Admin: toutes permissions                                     │
│  • Operator: exécuter scénarios, voir résultats                 │
│  • Viewer: lecture seule                                         │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 11. Roadmap

### 11.1 Vue d'Ensemble

```
2026                                    2027
 │                                       │
 ▼                                       ▼
 ┌───────────┬───────────┬───────────┬───────────┐
 │   M1-M3   │   M4-M6   │   M7-M9   │  M10-M12  │
 │   Alpha   │   Beta    │   MVP     │   V1.0    │
 └───────────┴───────────┴───────────┴───────────┘
      │            │            │            │
      ▼            ▼            ▼            ▼
   Fondations  Core Features  Polish     Production
```

### 11.2 Année 1 - Détail

#### Phase 1: Fondations (M1-M3) ✅ COMPLÈTE

**Objectifs:**
- [x] Architecture hexagonale serveur (Go)
- [x] Agent basique Windows (Rust)
- [x] Protocole de communication mTLS
- [x] 5 techniques Discovery de base
- [x] API REST CRUD basique
- [x] Dashboard squelette React

**Livrables:**
- ✅ Serveur Go fonctionnel avec API REST
- ✅ Agent Rust qui beacon et exécute des commandes
- ✅ Communication sécurisée mTLS
- ✅ UI basique pour voir les agents

**Techniques implémentées:**
- T1082 (System Information Discovery)
- T1083 (File and Directory Discovery)
- T1057 (Process Discovery)
- T1059.001 (PowerShell)
- T1059.003 (Windows Command Shell)

#### Phase 2: Core Features (M4-M6) ✅ COMPLÈTE

**Objectifs:**
- [x] Gestion des scénarios (CRUD)
- [x] Orchestrateur d'attaques
- [x] 10 techniques supplémentaires
- [x] Matrice ATT&CK interactive (CSS Grid)
- [x] WebSocket temps réel
- [x] Système de résultats

**Livrables:**
- ✅ Créer et exécuter des scénarios
- ✅ Visualisation matrice avec couleurs
- ✅ Monitoring temps réel des exécutions
- ✅ Score de sécurité basique

**Métriques atteintes (Phase 1+2) :**
| Métrique | Valeur |
|----------|--------|
| Tests | 780+ (200+ server + 513 dashboard + 67 agent) |
| Coverage | 95%+ sur le domaine |
| Techniques MITRE | 48 (13 tactiques sur 14) |
| Issues GitHub | 170+ fermées |
| Lignes de code | ~18,000 |

**Techniques ajoutées:**
- T1016 (System Network Configuration)
- T1069 (Permission Groups Discovery)
- T1087 (Account Discovery)
- T1053.005 (Scheduled Task)
- T1547.001 (Registry Run Keys)
- T1070.004 (File Deletion)
- T1059.004 (Unix Shell)
- T1018 (Remote System Discovery)
- T1049 (System Network Connections)
- T1007 (System Service Discovery)

#### Phase 3: Polish (M7-M9) ✅ COMPLÈTE

**Objectifs:**
- [x] Agent Linux (cross-compilation)
- [x] Authentification complète (JWT, 5 rôles, 28 permissions, token blacklist)
- [x] Security hardening (rate limiting, security headers, CSP, HSTS)
- [x] Amélioration UX dashboard (12+ pages)
- [x] Documentation technique (MkDocs)
- [x] Tests unitaires et intégration (780+ tests)

**Livrables:**
- ✅ Support multi-plateforme (Windows + Linux)
- ✅ Auth complète avec RBAC granulaire
- ✅ Security headers, rate limiting, audit logging
- ✅ Documentation complète

#### Phase 4: V1.0 (M10-M12) ✅ COMPLÈTE

**Objectifs:**
- [x] 48 techniques MITRE (13 tactiques sur 14)
- [x] Mode déploiement Docker (docker-compose prod + dev)
- [x] Hardening sécurité (rate limiting, CSP, HSTS, mTLS)
- [x] Scheduling (cron, daily, weekly, monthly)
- [x] Notifications (email SMTP + webhooks)
- [x] Analytics (comparaison périodes, tendances, Security Score)
- [x] Import/Export scénarios (YAML/JSON)
- [x] CI/CD complet (GitHub Actions, SonarCloud)

**Livrables:**
- ✅ Version 1.0 stable
- ✅ Documentation utilisateur (MkDocs sur GitHub Pages)
- ✅ Guide de déploiement (Docker)
- ✅ Démo fonctionnelle complète

#### Phase 5: Features Avancées (en cours)

**Priorité haute:**
- [ ] Profils adversaires APT (APT29, Ransomware, Insider Threat, Full Kill Chain)
- [ ] Export rapports PDF (handler backend + générateur + page Reports)

**Priorité moyenne:**
- [ ] ScenarioBuilder visuel (drag & drop)
- [ ] Cleanup automatique post-exécution
- [ ] Agent auto-deploy (scripts bash/PowerShell)
- [ ] LiveLogs (page logs temps réel via WebSocket)

### 11.3 Année 2 (Extension)

**Stretch Goals - Features Caldera-like:**
- **Planners intelligents** : séquentiel (arrêt si échec), conditionnel (décision dynamique), buckets (randomisation par tactique)
- **Facts / Data Exchange** : passer des données entre techniques (ex: users découverts → cibles brute force)
- **Obfuscation des commandes** : Base64, concaténation, substitution de variables
- **Recommandations de remédiation** : mapping mitigations ATT&CK automatique post-exécution
- **Multiple agent types** : agent Python léger, agent reverse shell

**Fonctionnalités potentielles:**
- Agent macOS
- Cloud testing (AWS, Azure)
- Intégration SIEM (Splunk, ELK)
- API publique pour intégrations
- Marketplace de techniques communautaires
- Mode SaaS multi-tenant

---

## 12. Organisation de l'Équipe

### 12.1 Membres

| Rôle | Responsabilités |
|------|-----------------|
| **Project Lead** | Architecture globale, coordination |
| **Security Lead** | Techniques MITRE, agent Rust, tests sécu |
| **Backend Lead** | Control Server Go, API, DB |
| **Frontend Lead** | Dashboard React, UX/UI |
| **DevOps / QA** | CI/CD, tests, documentation |

### 12.2 Répartition des Tâches

```
┌──────────────────────────────────────────────────────────────────┐
│                    Répartition par Composant                      │
├──────────────────────────────────────────────────────────────────┤
│                                                                   │
│  AGENT (Rust)          │  SERVER (Go)         │  DASHBOARD (React)│
│  ──────────────────    │  ──────────────────  │  ─────────────────│
│  Security Lead         │  Backend Lead        │  Frontend Lead    │
│  • Communication       │  • API REST          │  • Matrice ATT&CK │
│  • Executors           │  • WebSocket         │  • Scénarios UI   │
│  • Techniques MITRE    │  • Orchestrateur     │  • Agents Manager │
│  • Évasion             │  • Persistence       │  • Rapports       │
│                        │                      │                   │
│  Project Lead (support)│  Project Lead (review)│ Project Lead (review)
│                                                                   │
├──────────────────────────────────────────────────────────────────┤
│                                                                   │
│  TRANSVERSE            │                                          │
│  ──────────────────    │                                          │
│  Project Lead          │                                          │
│  • Architecture        │                                          │
│  • Documentation       │                                          │
│  • Gestion projet      │                                          │
│                        │                                          │
│  [DevOps]              │                                          │
│  • CI/CD               │                                          │
│  • Docker              │                                          │
│  • Tests E2E           │                                          │
│                                                                   │
└──────────────────────────────────────────────────────────────────┘
```

### 12.3 Outils de Collaboration

| Outil | Usage |
|-------|-------|
| **GitHub** | Code, issues, PR, projects |
| **Discord** | Communication quotidienne |
| **Notion/Confluence** | Documentation, wiki |
| **Figma** | Design UI/UX |
| **Linear/Jira** | Gestion tickets (optionnel) |

### 12.4 Méthodologie

- **Sprints de 2 semaines**
- **Daily standups** (async sur Discord)
- **Weekly sync** (visio 1h)
- **Code review** obligatoire
- **Branche par feature** (GitFlow simplifié)

---

## 13. Ressources et Références

### 13.1 Documentation Sécurité

| Ressource | Origine | URL |
|-----------|---------|-----|
| MITRE ATT&CK | USA | https://attack.mitre.org/ |
| ATT&CK Navigator | USA | https://mitre-attack.github.io/attack-navigator/ |
| **EBIOS RM** | **France (ANSSI)** | **https://cyber.gouv.fr/la-methode-ebios-risk-manager** |
| MITRE Caldera | USA | https://github.com/mitre/caldera |
| Atomic Red Team | USA | https://github.com/redcanaryco/atomic-red-team |

### 13.2 Projets Open-Source de Référence

| Projet | Description | Utilité |
|--------|-------------|---------|
| **MITRE Caldera** | Adversary emulation platform | Architecture de référence |
| **Atomic Red Team** | Bibliothèque de tests | Définitions techniques |
| **OpenBAS** | BAS open-source | Inspiration fonctionnelle |
| **Sliver** | C2 framework | Architecture agent/server |
| **Havoc** | C2 framework moderne | UI/UX inspiration |

### 13.3 Ressources Techniques

| Sujet | Ressource |
|-------|-----------|
| **Go** | https://go.dev/doc/ |
| **Rust** | https://doc.rust-lang.org/book/ |
| **React** | https://react.dev/ |
| **D3.js** | https://d3js.org/ |
| **Architecture Hexagonale** | https://alistair.cockburn.us/hexagonal-architecture/ |

### 13.4 Livres Recommandés

- *The Art of Attack* - Maxie Reynolds
- *Red Team Development and Operations* - Joe Vest
- *Practical Malware Analysis* - Michael Sikorski
- *Clean Architecture* - Robert C. Martin

### 13.5 Documents Complémentaires

| Document | Contenu | Usage |
|----------|---------|-------|
| [ROADMAP.md](./ROADMAP.md) | Issues GitHub, effort estimé, timeline 2026-2028 | Planification opérationnelle |
| [PRESENTATION.md](./PRESENTATION.md) | Slides, arguments clés, FAQ | Présentation équipe |
| [VISION_V2.md](./VISION_V2.md) | Decision Engine, Blackbox, Agent propagation | Vision autonome V2 |
| [CLAUDE.md](../CLAUDE.md) | Contexte technique pour IA | Assistance développement |

---

## Annexes

### A. Glossaire

| Terme | Définition |
|-------|------------|
| **BAS** | Breach and Attack Simulation |
| **CRTE** | Connaître, Rentrer, Trouver, Exploiter - Phases d'attaque EBIOS RM |
| **EBIOS RM** | Expression des Besoins et Identification des Objectifs de Sécurité - Risk Manager (méthode ANSSI) |
| **TTP** | Tactics, Techniques, and Procedures |
| **EDR** | Endpoint Detection and Response |
| **SIEM** | Security Information and Event Management |
| **C2/C&C** | Command and Control |
| **APT** | Advanced Persistent Threat |
| **IOC** | Indicator of Compromise |
| **mTLS** | Mutual TLS (authentification bidirectionnelle) |
| **Beaconing** | Communication périodique agent → serveur |
| **Paw** | Identifiant unique d'un agent (terminologie Caldera) |

### B. Commandes Utiles

```bash
# Démarrer le serveur (dev)
cd server && go run cmd/autostrike/main.go

# Compiler l'agent Windows
cd agent && cargo build --release --target x86_64-pc-windows-gnu

# Démarrer le dashboard (dev)
cd dashboard && npm run dev

# Docker compose (tout)
docker-compose up -d

# Générer certificats mTLS
./scripts/generate-certs.sh

# Importer techniques MITRE
./scripts/import-mitre.sh
```

### C. Variables d'Environnement

```bash
# Server
AUTOSTRIKE_PORT=8443
AUTOSTRIKE_DB_PATH=./data/autostrike.db
AUTOSTRIKE_CERT_PATH=./certs/server.crt
AUTOSTRIKE_KEY_PATH=./certs/server.key
AUTOSTRIKE_CA_PATH=./certs/ca.crt
AUTOSTRIKE_LOG_LEVEL=info

# Agent
AUTOSTRIKE_SERVER=https://server:8443
AUTOSTRIKE_SLEEP=5
AUTOSTRIKE_JITTER=3

# Dashboard
VITE_API_URL=https://localhost:8443
VITE_WS_URL=wss://localhost:8443
```

---

## Licence

Ce projet est développé dans le cadre de l'EIP EPITECH.

**Auteurs:** Équipe EIP AutoStrike

**Année:** 2026-2028 (Promotion 2028)

---

## Historique des Versions

| Version | Date | Changements |
|---------|------|-------------|
| 1.0.0 | Jan 2026 | Version initiale |
| 2.0.0 | Fév 2026 | Ajout EBIOS RM (CRTE), mise à jour métriques Phase 1+2, cross-références |

---

*Document mis à jour le 2026-02-03*
*Version: 2.0.0*
