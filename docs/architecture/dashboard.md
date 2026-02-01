# Dashboard (React)

Le dashboard AutoStrike est développé en **React 18** avec **TypeScript** et **TailwindCSS**.

---

## Stack Technique

| Technologie | Version | Usage |
|-------------|---------|-------|
| React | 18.2 | Framework UI |
| TypeScript | 5.3 | Typage statique |
| Vite | 5.0 | Build tool |
| TailwindCSS | 3.4 | Styling utility-first |
| TanStack Query | 5.17 | Data fetching & caching |
| React Router | 6.21 | Navigation SPA |
| Chart.js | 4.4 | Graphiques (Doughnut) |
| Axios | 1.6 | Client HTTP |
| Zustand | 4.4 | State management |

---

## Structure

```
dashboard/
├── src/
│   ├── main.tsx              # Point d'entrée React + providers
│   ├── App.tsx               # Configuration des routes
│   ├── index.css             # Styles globaux + directives Tailwind
│   ├── components/
│   │   └── Layout.tsx        # Layout avec sidebar navigation
│   ├── pages/
│   │   ├── Dashboard.tsx     # Vue d'ensemble, KPIs, graphiques
│   │   ├── Agents.tsx        # Liste et gestion des agents
│   │   ├── Techniques.tsx    # Catalogue MITRE ATT&CK
│   │   ├── Scenarios.tsx     # Scénarios d'attaque
│   │   ├── Executions.tsx    # Historique des exécutions
│   │   └── Settings.tsx      # Configuration application
│   └── lib/
│       └── api.ts            # Client Axios avec intercepteurs
├── public/
├── package.json
├── vite.config.ts            # Config Vite + proxy API
├── tailwind.config.js        # Config Tailwind + couleurs custom
├── tsconfig.json
└── Dockerfile                # Build multi-stage + nginx
```

---

## Pages

### Dashboard (`/dashboard`)

Vue d'ensemble avec :
- **Agents Online** : Nombre d'agents connectés
- **Security Score** : Score global de la dernière exécution
- **Techniques Tested** : Nombre de techniques testées
- **Graphique Doughnut** : Répartition blocked/detected/successful
- **Activité Récente** : 5 dernières exécutions

### Agents (`/agents`)

Gestion des agents :
- Cartes avec statut (online/offline)
- Informations : hostname, username, platform
- Executors disponibles (badges)
- Last seen avec formatage relatif

### Techniques (`/techniques`)

Catalogue MITRE ATT&CK :
- Tableau avec ID, nom, description
- Filtrage par tactique et plateforme
- Badge safe/unsafe
- Import depuis fichiers YAML

### Scenarios (`/scenarios`)

Scénarios d'attaque :
- Cartes avec phases numérotées
- Nombre de techniques par phase
- Tags de catégorie
- Bouton Run

### Executions (`/executions`)

Historique :
- Tableau avec score, statut, mode
- Détails : blocked/detected/successful
- Temps d'exécution
- Lien vers les résultats détaillés

### Settings (`/settings`)

Configuration :
- URL du serveur
- Mode safe par défaut (toggle)
- Paramètres agents (heartbeat, timeout)
- Chemins certificats TLS

---

## Composants UI

### Classes CSS Custom (Tailwind)

```css
.btn-primary    /* Bouton principal bleu */
.btn-danger     /* Bouton danger rouge */
.card           /* Carte avec ombre et bordure */
.input          /* Champ de saisie stylé */
.badge          /* Badge inline */
.badge-success  /* Badge vert (online, safe) */
.badge-danger   /* Badge rouge (offline, unsafe) */
.badge-warning  /* Badge orange (running) */
```

### Palette Couleurs

```javascript
// tailwind.config.js
colors: {
  primary: { 50-950 },      // Bleu
  danger: { 50, 500, 600 }, // Rouge
  success: { 50, 500, 600 },// Vert
  warning: { 50, 500, 600 } // Orange
}
```

---

## Data Fetching

Utilisation de **TanStack Query** pour le fetching et le caching :

```typescript
import { useQuery } from '@tanstack/react-query';
import { api } from '@/lib/api';

const { data: agents, isLoading, isError } = useQuery({
  queryKey: ['agents'],
  queryFn: () => api.get('/agents').then(res => res.data),
});
```

---

## API Client

Le client API (`src/lib/api.ts`) inclut :
- Base URL configurable (`/api`)
- Intercepteur pour injecter le token JWT
- Gestion des erreurs 401 (redirection login)

```typescript
import { api } from '@/lib/api';

// GET
const response = await api.get('/agents');

// POST
await api.post('/executions', {
  scenario_id: 'scn-001',
  agent_paws: ['agent-001', 'agent-002'],
  safe_mode: true
});
```

---

## Configuration Vite

```typescript
// vite.config.ts
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: { '@': path.resolve(__dirname, './src') }
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'https://localhost:8443',
        changeOrigin: true,
        secure: false,
      },
      '/ws': {
        target: 'wss://localhost:8443',
        ws: true,
        secure: false,
      },
    },
  },
});
```

---

## Développement

```bash
cd dashboard
npm install       # Installer les dépendances
npm run dev       # Serveur de développement (port 3000)
npm run build     # Build production
npm run preview   # Preview du build
npm run lint      # Vérification ESLint
npm run test      # Tests Vitest
```

---

## Docker

Le Dockerfile utilise un build multi-stage :
1. **Builder** : Node 20 Alpine pour le build
2. **Runtime** : Nginx Alpine pour servir les fichiers

```bash
docker build -t autostrike-dashboard .
docker run -p 3000:80 autostrike-dashboard
```

Nginx est configuré pour :
- Proxy `/api/` vers le backend
- Proxy `/ws/` pour WebSocket
- Fallback SPA (`try_files`)
- Cache des assets statiques
