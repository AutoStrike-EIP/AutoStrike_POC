import { useQuery } from '@tanstack/react-query';
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js';
import { Doughnut } from 'react-chartjs-2';
import { api } from '../lib/api';
import { Agent, Execution } from '../types';

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

  const onlineAgents = agents?.filter((a) => a.status === 'online').length || 0;
  const totalAgents = agents?.length || 0;

  const latestExecution = executions?.[0];

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
      <h1 className="text-3xl font-bold mb-8">Dashboard</h1>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div className="card">
          <p className="text-sm text-gray-500 mb-1">Agents Online</p>
          <p className="text-3xl font-bold text-primary-600">{onlineAgents}</p>
          <p className="text-sm text-gray-400">of {totalAgents} total</p>
        </div>

        <div className="card">
          <p className="text-sm text-gray-500 mb-1">Security Score</p>
          <p className="text-3xl font-bold text-success-600">
            {latestExecution?.score?.overall?.toFixed(1) || 0}%
          </p>
          <p className="text-sm text-gray-400">Latest execution</p>
        </div>

        <div className="card">
          <p className="text-sm text-gray-500 mb-1">Techniques Tested</p>
          <p className="text-3xl font-bold">{latestExecution?.score?.total || 0}</p>
          <p className="text-sm text-gray-400">In latest run</p>
        </div>

        <div className="card">
          <p className="text-sm text-gray-500 mb-1">Executions Today</p>
          <p className="text-3xl font-bold">{executions?.length || 0}</p>
          <p className="text-sm text-gray-400">Total runs</p>
        </div>
      </div>

      {/* Charts */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="card">
          <h2 className="text-lg font-semibold mb-4">Detection Results</h2>
          <div className="h-64 flex items-center justify-center">
            <Doughnut data={scoreData} options={{ maintainAspectRatio: false }} />
          </div>
        </div>

        <div className="card">
          <h2 className="text-lg font-semibold mb-4">Recent Activity</h2>
          <div className="space-y-4">
            {executions?.slice(0, 5).map((exec) => (
              <div key={exec.id} className="flex items-center justify-between py-2 border-b border-gray-100">
                <div>
                  <p className="font-medium">{exec.scenario_id}</p>
                  <p className="text-sm text-gray-500">{new Date(exec.started_at).toLocaleString()}</p>
                </div>
                <span className={`badge ${exec.status === 'completed' ? 'badge-success' : 'badge-warning'}`}>
                  {exec.status}
                </span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
