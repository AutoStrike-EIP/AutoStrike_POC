import { useQuery } from '@tanstack/react-query';
import { Squares2X2Icon } from '@heroicons/react/24/outline';
import { api } from '../lib/api';
import { Technique, TacticCoverage } from '../types';
import { LoadingState } from '../components/LoadingState';
import { EmptyState } from '../components/EmptyState';
import { MitreMatrix } from '../components/MitreMatrix';

/**
 * Matrix page component.
 * Displays the MITRE ATT&CK matrix visualization with technique coverage.
 *
 * @returns The Matrix page component
 */
export default function Matrix() {
  const { data: techniques, isLoading: techniquesLoading } = useQuery<Technique[]>({
    queryKey: ['techniques'],
    queryFn: () => api.get('/techniques').then(res => res.data),
  });

  const { data: coverage } = useQuery<TacticCoverage>({
    queryKey: ['techniques', 'coverage'],
    queryFn: () => api.get('/techniques/coverage').then(res => res.data),
  });

  if (techniquesLoading) {
    return <LoadingState message="Loading MITRE ATT&CK matrix..." />;
  }

  // Calculate coverage stats
  const totalTechniques = techniques?.length || 0;
  const coveredTactics = coverage ? Object.keys(coverage).filter(k => coverage[k] > 0).length : 0;
  const totalTactics = 14;

  return (
    <div>
      <div className="flex justify-between items-start mb-8">
        <div>
          <h1 className="text-3xl font-bold">MITRE ATT&CK Matrix</h1>
          <p className="text-gray-500 mt-1">Interactive visualization of available attack techniques</p>
        </div>
        <div className="flex gap-4">
          <div className="card p-4 text-center">
            <p className="text-2xl font-bold text-primary-600">{totalTechniques}</p>
            <p className="text-xs text-gray-500">Techniques</p>
          </div>
          <div className="card p-4 text-center">
            <p className="text-2xl font-bold text-primary-600">{coveredTactics}/{totalTactics}</p>
            <p className="text-xs text-gray-500">Tactics Covered</p>
          </div>
        </div>
      </div>

      {techniques && techniques.length > 0 ? (
        <div className="card p-4">
          <MitreMatrix techniques={techniques} />
        </div>
      ) : (
        <EmptyState
          icon={Squares2X2Icon}
          title="No techniques loaded"
          description="Import techniques to see the MITRE ATT&CK matrix"
        />
      )}
    </div>
  );
}
