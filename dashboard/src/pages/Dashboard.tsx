import { useQuery } from '@tanstack/react-query';
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js';
import { Doughnut } from 'react-chartjs-2';
import { api, techniqueApi } from '../lib/api';
import { Agent, Execution, Technique } from '../types';
import { SecurityScore } from '../components/SecurityScore';
import { CoverageReport } from '../components/CoverageReport';

ChartJS.register(ArcElement, Tooltip, Legend);

/**
 * Dashboard page component.
 * Displays an overview of agents, executions, and security scores.
 *
 * @returns The Dashboard page component
 */
export default function Dashboard() {
  const { data: agents } = useQuery<Agent[]>({
    queryKey: ['agents'],
    queryFn: () => api.get('/agents').then(res => res.data),
  });

  const { data: executions } = useQuery<Execution[]>({
    queryKey: ['executions'],
    queryFn: () => api.get('/executions').then(res => res.data),
  });

  const { data: coverage } = useQuery<Record<string, number>>({
    queryKey: ['techniques', 'coverage'],
    queryFn: () => techniqueApi.getCoverage().then(res => res.data),
  });

  const { data: techniques } = useQuery<Technique[]>({
    queryKey: ['techniques'],
    queryFn: () => api.get('/techniques').then(res => res.data),
  });

  const onlineAgents = agents?.filter((a) => a.status === 'online').length || 0;
  const totalAgents = agents?.length || 0;

  const latestExecution = executions?.[0];
  const totalTechniques = techniques?.length || 0;

  const scoreData = {
    labels: ['Blocked', 'Detected', 'Successful'],
    datasets: [
      {
        data: [
          latestExecution?.score?.blocked || 0,
          latestExecution?.score?.detected || 0,
          latestExecution?.score?.successful || 0,
        ],
        backgroundColor: ['#22c55e', '#f59e0b', '#ef4444'],
        borderWidth: 0,
      },
    ],
  };

  return (
    <div>
      <h1 className="text-3xl font-bold mb-8 text-gray-900 dark:text-gray-100">Dashboard</h1>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div className="card">
          <p className="text-sm text-gray-500 dark:text-gray-400 mb-1">Agents Online</p>
          <p className="text-3xl font-bold text-primary-600">{onlineAgents}</p>
          <p className="text-sm text-gray-400 dark:text-gray-500">of {totalAgents} total</p>
        </div>

        <div className="card">
          <p className="text-sm text-gray-500 dark:text-gray-400 mb-1">Security Score</p>
          <SecurityScore
            score={latestExecution?.score?.overall || 0}
            size="sm"
            animated
          />
        </div>

        <div className="card">
          <p className="text-sm text-gray-500 dark:text-gray-400 mb-1">Techniques Tested</p>
          <p className="text-3xl font-bold text-gray-900 dark:text-gray-100">{latestExecution?.score?.total || 0}</p>
          <p className="text-sm text-gray-400 dark:text-gray-500">In latest run</p>
        </div>

        <div className="card">
          <p className="text-sm text-gray-500 dark:text-gray-400 mb-1">Executions Today</p>
          <p className="text-3xl font-bold text-gray-900 dark:text-gray-100">{executions?.length || 0}</p>
          <p className="text-sm text-gray-400 dark:text-gray-500">Total runs</p>
        </div>
      </div>

      {/* Main content grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
        {/* Security Score with breakdown */}
        <div className="card">
          <h2 className="text-lg font-semibold mb-4 text-gray-900 dark:text-gray-100">Security Score</h2>
          <SecurityScore
            score={latestExecution?.score?.overall || 0}
            breakdown={latestExecution?.score}
            size="md"
            animated
          />
        </div>

        {/* Detection Results Chart */}
        <div className="card">
          <h2 className="text-lg font-semibold mb-4 text-gray-900 dark:text-gray-100">Detection Results</h2>
          <div className="h-64 flex items-center justify-center">
            <Doughnut data={scoreData} options={{ maintainAspectRatio: false }} />
          </div>
        </div>

        {/* MITRE Coverage */}
        <div className="card">
          <h2 className="text-lg font-semibold mb-4 text-gray-900 dark:text-gray-100">MITRE Coverage</h2>
          {coverage ? (
            <CoverageReport
              coverage={coverage}
              totalTechniques={totalTechniques}
              variant="compact"
            />
          ) : (
            <div className="text-center py-8 text-gray-500 dark:text-gray-400">
              Loading coverage data...
            </div>
          )}
        </div>
      </div>

      {/* Recent Activity */}
      <div className="card">
        <h2 className="text-lg font-semibold mb-4 text-gray-900 dark:text-gray-100">Recent Activity</h2>
        <div className="space-y-4">
          {executions?.slice(0, 5).map((exec) => (
            <div key={exec.id} className="flex items-center justify-between py-2 border-b border-gray-100 dark:border-gray-700">
              <div>
                <p className="font-medium text-gray-900 dark:text-gray-100">{exec.scenario_id}</p>
                <p className="text-sm text-gray-500 dark:text-gray-400">{new Date(exec.started_at).toLocaleString()}</p>
              </div>
              <span className={`badge ${exec.status === 'completed' ? 'badge-success' : 'badge-warning'}`}>
                {exec.status}
              </span>
            </div>
          ))}
          {(!executions || executions.length === 0) && (
            <p className="text-center py-4 text-gray-500 dark:text-gray-400">No recent executions</p>
          )}
        </div>
      </div>
    </div>
  );
}
