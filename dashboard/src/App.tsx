import { Routes, Route, Navigate } from 'react-router-dom';
import Layout from './components/Layout';
import { ErrorBoundary } from './components/ErrorBoundary';
import Dashboard from './pages/Dashboard';
import Agents from './pages/Agents';
import Techniques from './pages/Techniques';
import Scenarios from './pages/Scenarios';
import Executions from './pages/Executions';
import Settings from './pages/Settings';

/**
 * Root application component.
 * Sets up routing and global error handling.
 *
 * @returns The App component with routes
 */
function App() {
  return (
    <ErrorBoundary>
      <Layout>
        <Routes>
          <Route path="/" element={<Navigate to="/dashboard" replace />} />
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/agents" element={<Agents />} />
          <Route path="/techniques" element={<Techniques />} />
          <Route path="/scenarios" element={<Scenarios />} />
          <Route path="/executions" element={<Executions />} />
          <Route path="/settings" element={<Settings />} />
        </Routes>
      </Layout>
    </ErrorBoundary>
  );
}

export default App;
