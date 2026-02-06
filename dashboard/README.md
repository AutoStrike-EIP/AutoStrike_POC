# AutoStrike Dashboard

Interface web React pour la gestion et le monitoring de la plateforme AutoStrike.

## Architecture

```
dashboard/
├── src/
│   ├── main.tsx              # Point d'entrée React + AuthProvider
│   ├── App.tsx               # Routes principales
│   ├── index.css             # Styles globaux + Tailwind
│   ├── components/
│   │   ├── Layout.tsx        # Layout avec sidebar + logout
│   │   ├── ProtectedRoute.tsx # Protection des routes
│   │   ├── MitreMatrix.tsx   # Matrice MITRE ATT&CK interactive
│   │   ├── RunExecutionModal.tsx # Modal de lancement d'exécution
│   │   ├── SecurityScore.tsx # Visualisation du score de sécurité
│   │   ├── CoverageReport.tsx # Rapport de couverture MITRE
│   │   ├── Modal.tsx         # Composant modal réutilisable
│   │   ├── Table.tsx         # Composant table réutilisable
│   │   ├── ThemeToggle.tsx   # Toggle thème sombre/clair
│   │   ├── LoadingState.tsx  # Spinner de chargement
│   │   ├── EmptyState.tsx    # État vide
│   │   └── ErrorBoundary.tsx # Gestion d'erreurs
│   ├── contexts/
│   │   └── AuthContext.tsx   # État d'authentification
│   ├── hooks/
│   │   └── useWebSocket.ts  # Hook WebSocket temps réel
│   ├── pages/
│   │   ├── Dashboard.tsx     # Vue d'ensemble, graphiques
│   │   ├── Agents.tsx        # Gestion des agents
│   │   ├── Techniques.tsx    # Catalogue MITRE ATT&CK
│   │   ├── Matrix.tsx        # Matrice MITRE ATT&CK
│   │   ├── Scenarios.tsx     # Scénarios d'attaque
│   │   ├── Executions.tsx    # Historique des exécutions
│   │   ├── ExecutionDetails.tsx # Détails d'une exécution
│   │   ├── Analytics.tsx     # Tendances et rapports
│   │   ├── Scheduler.tsx     # Planification des exécutions
│   │   ├── Settings.tsx      # Configuration
│   │   ├── Login.tsx         # Page de connexion
│   │   └── Admin/            # Pages admin
│   │       ├── Users.tsx     # Gestion des utilisateurs
│   │       └── Permissions.tsx # Matrice des permissions
│   └── lib/
│       └── api.ts            # Client Axios + authApi
├── package.json
├── vite.config.ts
├── tailwind.config.js
└── tsconfig.json
```

## Stack Technique

| Technologie | Version | Usage |
|-------------|---------|-------|
| React | 18.2 | Framework UI |
| TypeScript | 5.3 | Typage statique |
| Vite | 7.3 | Build tool |
| TailwindCSS | 3.4 | Styling |
| TanStack Query | 5.17 | Data fetching |
| React Router | 6.21 | Routing |
| Chart.js | 4.4 | Graphiques |
| Axios | 1.6 | HTTP client |
| Zustand | 4.4 | State management |
| Headless UI | 1.7 | Composants accessibles |
| Heroicons | 2.1 | Icônes SVG |
| date-fns | 3.2 | Formatage de dates |
| react-hot-toast | 2.4 | Notifications toast |
| Vitest | 4.0 | Tests unitaires |

## Prérequis

- Node.js 18+
- npm ou yarn

## Installation

```bash
# Installer les dépendances
npm install

# Lancer en développement
npm run dev

# Build production
npm run build
```

## Scripts Disponibles

| Script | Description |
|--------|-------------|
| `npm run dev` | Serveur de développement (port 3000) |
| `npm run build` | Build de production |
| `npm run preview` | Preview du build |
| `npm run lint` | Vérification ESLint |
| `npm run type-check` | Vérification TypeScript |
| `npm test` | Tests Vitest (513 tests) |

## Configuration

### Proxy API (développement)

Le fichier `vite.config.ts` configure le proxy vers le backend :

```typescript
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
}
```

### Alias d'import

```typescript
// Utiliser @/ pour importer depuis src/
import { api } from '@/lib/api';
```

## Pages (13 pages)

### Dashboard (`/dashboard`)
- Statistiques agents online/total
- Score de sécurité global (composant SecurityScore)
- Graphique donut des résultats (blocked/detected/successful)
- Activité récente

### Agents (`/agents`)
- Liste des agents avec statut
- Informations: hostname, platform, executors
- Last seen avec formatage relatif

### Techniques (`/techniques`)
- Tableau des 48 techniques MITRE ATT&CK
- Filtrage par tactique et plateforme
- Badge safe/unsafe
- Import depuis fichiers YAML

### Matrix (`/matrix`)
- Matrice MITRE ATT&CK interactive (14 colonnes)
- Statistiques de couverture (composant CoverageReport)
- Filtrage par plateforme

### Scenarios (`/scenarios`)
- Cartes de scénarios avec phases
- Tags et nombre de techniques
- Bouton Run → RunExecutionModal
- Import/export

### Executions (`/executions`)
- Historique des exécutions
- Score, statut, mode (safe/full)
- Mises à jour temps réel via WebSocket

### ExecutionDetails (`/executions/:id`)
- Détails des résultats par technique
- Breakdown blocked/detected/successful
- Output extensible

### Analytics (`/analytics`)
- Graphiques de tendance du score (7j, 30j, 90j)
- Comparaison de périodes
- Résumé des exécutions
- Score par tactique

### Scheduler (`/scheduler`)
- Création/modification de schedules
- Fréquences (once, hourly, daily, weekly, monthly, cron)
- Historique des exécutions par schedule
- Pause/reprise

### Settings (`/settings`)
- Configuration serveur
- Mode safe par défaut
- Paramètres de notification

### Login (`/login`)
- Formulaire de connexion

### Admin/Users (`/admin/users`)
- Gestion des utilisateurs (5 rôles: admin, rssi, operator, analyst, viewer)
- Activation/désactivation

### Admin/Permissions (`/admin/permissions`)
- Matrice des 28 permissions par rôle

## Composants (12)

| Composant | Description |
|-----------|-------------|
| Layout | Sidebar navigation + logout |
| ProtectedRoute | Protection des routes |
| MitreMatrix | Matrice MITRE ATT&CK interactive |
| RunExecutionModal | Modal de configuration d'exécution |
| SecurityScore | Visualisation du score |
| CoverageReport | Rapport de couverture MITRE |
| Modal | Modal réutilisable |
| Table | Table réutilisable |
| ThemeToggle | Toggle thème sombre/clair |
| LoadingState | Spinner de chargement |
| EmptyState | État vide |
| ErrorBoundary | Gestion d'erreurs |

## API Client

Le client API (`src/lib/api.ts`) :
- Base URL: `/api/v1`
- Intercepteur pour token JWT
- Gestion du 401 (redirection login)

```typescript
import { api } from '@/lib/api';

// GET
const agents = await api.get('/agents');

// POST
await api.post('/executions', { scenario_id, agent_paws, safe_mode });
```

## Data Fetching

Utilisation de TanStack Query :

```typescript
const { data: agents, isLoading } = useQuery({
  queryKey: ['agents'],
  queryFn: () => api.get('/agents').then(res => res.data),
});
```

## Docker

```bash
docker build -t autostrike-dashboard .
docker run -p 3000:80 autostrike-dashboard
```

Le conteneur utilise Nginx pour servir les fichiers statiques avec :
- Proxy vers le backend (`/api/`, `/ws/`)
- Fallback SPA (`try_files`)
- Cache des assets statiques

## Tests

513 tests across 25 fichiers :

```bash
npm test          # Mode watch
npm test -- --run # Une seule exécution
npm test -- --coverage # Avec couverture
```

## Linting

```bash
npm run lint
npm run type-check
```

Configuration ESLint stricte avec `max-warnings: 0`.
