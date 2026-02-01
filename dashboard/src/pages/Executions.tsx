import { useQuery } from '@tanstack/react-query';
import { PlayIcon } from '@heroicons/react/24/outline';
import { api } from '../lib/api';
import { formatDistanceToNow } from 'date-fns';
import { Execution } from '../types';
import { LoadingState } from '../components/LoadingState';
import { EmptyState } from '../components/EmptyState';

/**
 * Returns the appropriate badge class for an execution status.
 */
function getStatusBadgeClass(status: string): string {
  switch (status) {
    case 'completed':
      return 'badge-success';
    case 'running':
      return 'badge-warning';
    default:
      return 'badge-danger';
  }
}

/**
 * Executions page component.
 * Displays a table of scenario executions with their results and scores.
 *
 * @returns The Executions page component
 */
export default function Executions() {
  const { data: executions, isLoading } = useQuery<Execution[]>({
    queryKey: ['executions'],
    queryFn: () => api.get('/executions').then(res => res.data),
  });

  if (isLoading) {
    return <LoadingState message="Loading executions..." />;
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold">Executions</h1>
        <button className="btn-primary">New Execution</button>
      </div>

      <div className="card overflow-hidden">
        <table className="w-full">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Scenario
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Status
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Score
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Results
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Started
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Mode
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {executions?.map((execution) => (
              <tr key={execution.id} className="hover:bg-gray-50 cursor-pointer">
                <td className="px-6 py-4 whitespace-nowrap">
                  <p className="font-medium">{execution.scenario_id}</p>
                  <p className="text-xs text-gray-400">{execution.id.slice(0, 8)}...</p>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className={`badge ${getStatusBadgeClass(execution.status)}`}>
                    {execution.status}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className="text-2xl font-bold">
                    {execution.score?.overall?.toFixed(1) || '-'}%
                  </span>
                </td>
                <td className="px-6 py-4">
                  <div className="flex gap-2 text-sm">
                    <span className="text-success-600">
                      {execution.score?.blocked || 0} blocked
                    </span>
                    <span className="text-warning-600">
                      {execution.score?.detected || 0} detected
                    </span>
                    <span className="text-danger-600">
                      {execution.score?.successful || 0} success
                    </span>
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {formatDistanceToNow(new Date(execution.started_at), { addSuffix: true })}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className={`badge ${execution.safe_mode ? 'badge-success' : 'badge-danger'}`}>
                    {execution.safe_mode ? 'Safe' : 'Full'}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {executions?.length === 0 && (
        <EmptyState
          icon={PlayIcon}
          title="No executions yet"
          description="Run a scenario to see results here"
        />
      )}
    </div>
  );
}
