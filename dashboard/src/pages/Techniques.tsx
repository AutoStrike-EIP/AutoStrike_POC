import { useQuery } from '@tanstack/react-query';
import { ShieldExclamationIcon } from '@heroicons/react/24/outline';
import { api } from '../lib/api';
import { Technique } from '../types';
import { LoadingState } from '../components/LoadingState';
import { EmptyState } from '../components/EmptyState';

/**
 * Color mapping for MITRE ATT&CK tactics.
 * Maps each tactic to its corresponding Tailwind CSS classes.
 */
const tacticColors: Record<string, string> = {
  reconnaissance: 'bg-purple-100 text-purple-700',
  resource_development: 'bg-blue-100 text-blue-700',
  initial_access: 'bg-red-100 text-red-700',
  execution: 'bg-orange-100 text-orange-700',
  persistence: 'bg-yellow-100 text-yellow-700',
  privilege_escalation: 'bg-pink-100 text-pink-700',
  defense_evasion: 'bg-green-100 text-green-700',
  credential_access: 'bg-indigo-100 text-indigo-700',
  discovery: 'bg-cyan-100 text-cyan-700',
  lateral_movement: 'bg-teal-100 text-teal-700',
  collection: 'bg-lime-100 text-lime-700',
  command_and_control: 'bg-amber-100 text-amber-700',
  exfiltration: 'bg-rose-100 text-rose-700',
  impact: 'bg-red-100 text-red-700',
};

/**
 * Techniques page component.
 * Displays a table of available MITRE ATT&CK techniques.
 *
 * @returns The Techniques page component
 */
export default function Techniques() {
  const { data: techniques, isLoading } = useQuery<Technique[]>({
    queryKey: ['techniques'],
    queryFn: () => api.get('/techniques').then(res => res.data),
  });

  if (isLoading) {
    return <LoadingState message="Loading techniques..." />;
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold">Techniques</h1>
        <button className="btn-primary">Import Techniques</button>
      </div>

      <div className="card overflow-hidden">
        <table className="w-full">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                ID
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Name
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Tactic
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Platforms
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Safe
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {techniques?.map((technique) => (
              <tr key={technique.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className="font-mono text-sm">{technique.id}</span>
                </td>
                <td className="px-6 py-4">
                  <p className="font-medium">{technique.name}</p>
                  <p className="text-sm text-gray-500 truncate max-w-md">
                    {technique.description}
                  </p>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span
                    className={`badge ${tacticColors[technique.tactic] || 'bg-gray-100 text-gray-700'}`}
                  >
                    {String(technique.tactic).replaceAll('_', ' ')}
                  </span>
                </td>
                <td className="px-6 py-4">
                  <div className="flex gap-1">
                    {technique.platforms.map((platform) => (
                      <span key={platform} className="badge bg-gray-100 text-gray-700">
                        {platform}
                      </span>
                    ))}
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className={`badge ${technique.is_safe ? 'badge-success' : 'badge-danger'}`}>
                    {technique.is_safe ? 'Safe' : 'Unsafe'}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {techniques?.length === 0 && (
        <EmptyState
          icon={ShieldExclamationIcon}
          title="No techniques loaded"
          description="Import techniques from Atomic Red Team"
        />
      )}
    </div>
  );
}
