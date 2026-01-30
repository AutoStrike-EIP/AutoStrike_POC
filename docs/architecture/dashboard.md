# Dashboard (React)

Le dashboard AutoStrike est développé en **React 18** avec **TypeScript**.

---

## Stack Frontend

- **React 18** - Framework UI
- **TypeScript** - Typage statique
- **TailwindCSS** - Styling
- **D3.js** - Visualisation matrice MITRE
- **React Query** - Gestion des données
- **WebSocket** - Temps réel

---

## Structure

```
dashboard/
├── src/
│   ├── components/
│   │   ├── Layout/
│   │   ├── MitreMatrix/
│   │   ├── Agents/
│   │   ├── Scenarios/
│   │   └── Results/
│   ├── pages/
│   │   ├── Dashboard.tsx
│   │   ├── Agents.tsx
│   │   ├── Scenarios.tsx
│   │   └── Reports.tsx
│   ├── hooks/
│   │   ├── useAgents.ts
│   │   ├── useWebSocket.ts
│   │   └── useScenarios.ts
│   ├── services/
│   │   └── api.ts
│   ├── App.tsx
│   └── main.tsx
├── package.json
└── tsconfig.json
```

---

## Fonctionnalités

### Matrice MITRE ATT&CK

Visualisation interactive de la couverture de détection avec code couleur.

### Monitoring temps réel

WebSocket pour afficher les résultats des exécutions en direct.

### Rapports

Export PDF des résultats pour les audits de sécurité.

---

## Développement

```bash
cd dashboard
npm install
npm run dev    # Mode développement
npm run build  # Build production
npm run test   # Tests
```
