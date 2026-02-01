import { useQuery } from '@tanstack/react-query';
import { ComputerDesktopIcon } from '@heroicons/react/24/outline';
import { api } from '../lib/api';
import { formatDistanceToNow } from 'date-fns';
import { Agent } from '../types';
import { LoadingState } from '../components/LoadingState';
import { EmptyState } from '../components/EmptyState';

/**
 * Agents page component.
 * Displays a grid of connected agents with their status and details.
 *
 * @returns The Agents page component
 */
export default function Agents() {
  const { data: agents, isLoading } = useQuery<Agent[]>({
    queryKey: ['agents'],
    queryFn: () => api.get('/agents').then(res => res.data),
  });

  if (isLoading) {
    return <LoadingState message="Loading agents..." />;
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold">Agents</h1>
        <button className="btn-primary">Add Agent</button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {agents?.map((agent) => (
          <div key={agent.paw} className="card">
            <div className="flex items-start gap-4">
              <div className="p-3 bg-gray-100 rounded-lg">
                <ComputerDesktopIcon className="h-6 w-6 text-gray-600" />
              </div>
              <div className="flex-1">
                <div className="flex items-center justify-between">
                  <h3 className="font-semibold">{agent.hostname}</h3>
                  <span
                    className={`badge ${
                      agent.status === 'online' ? 'badge-success' : 'badge-danger'
                    }`}
                  >
                    {agent.status}
                  </span>
                </div>
                <p className="text-sm text-gray-500 mt-1">{agent.username}</p>
                <p className="text-sm text-gray-400 mt-1">PAW: {agent.paw.slice(0, 8)}...</p>
              </div>
            </div>

            <div className="mt-4 pt-4 border-t border-gray-100">
              <div className="grid grid-cols-2 gap-4 text-sm">
                <div>
                  <p className="text-gray-500">Platform</p>
                  <p className="font-medium">{agent.platform}</p>
                </div>
                <div>
                  <p className="text-gray-500">Last Seen</p>
                  <p className="font-medium">
                    {formatDistanceToNow(new Date(agent.last_seen), { addSuffix: true })}
                  </p>
                </div>
              </div>

              <div className="mt-3">
                <p className="text-gray-500 text-sm mb-1">Executors</p>
                <div className="flex gap-1 flex-wrap">
                  {agent.executors.map((exec) => (
                    <span key={exec} className="badge bg-gray-100 text-gray-700">
                      {exec}
                    </span>
                  ))}
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>

      {agents?.length === 0 && (
        <EmptyState
          icon={ComputerDesktopIcon}
          title="No agents connected"
          description="Deploy an agent to get started"
        />
      )}
    </div>
  );
}
