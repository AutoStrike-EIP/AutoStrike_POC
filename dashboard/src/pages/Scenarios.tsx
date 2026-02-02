import { useState } from 'react';
import { useQuery, useMutation } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { DocumentTextIcon } from '@heroicons/react/24/outline';
import { api, executionApi } from '../lib/api';
import { Scenario } from '../types';
import { LoadingState } from '../components/LoadingState';
import { EmptyState } from '../components/EmptyState';
import { RunExecutionModal } from '../components/RunExecutionModal';
import toast from 'react-hot-toast';

/**
 * Scenarios page component.
 * Displays a grid of attack scenarios with their phases and configuration.
 *
 * @returns The Scenarios page component
 */
export default function Scenarios() {
  const navigate = useNavigate();
  const [scenarioToRun, setScenarioToRun] = useState<Scenario | null>(null);

  const { data: scenarios, isLoading } = useQuery<Scenario[]>({
    queryKey: ['scenarios'],
    queryFn: () => api.get('/scenarios').then(res => res.data),
  });

  const startMutation = useMutation({
    mutationFn: ({ scenarioId, agentPaws, safeMode }: { scenarioId: string; agentPaws: string[]; safeMode: boolean }) =>
      executionApi.start(scenarioId, agentPaws, safeMode),
    onSuccess: () => {
      toast.success('Execution started successfully');
      setScenarioToRun(null);
      navigate('/executions');
    },
    onError: (error: { response?: { data?: { error?: string } } }) => {
      const message = error.response?.data?.error || 'Failed to start execution';
      toast.error(message);
    },
  });

  const handleRunClick = (scenario: Scenario) => {
    setScenarioToRun(scenario);
  };

  const handleConfirmRun = (agentPaws: string[], safeMode: boolean) => {
    if (scenarioToRun) {
      startMutation.mutate({
        scenarioId: scenarioToRun.id,
        agentPaws,
        safeMode,
      });
    }
  };

  const handleCancelRun = () => {
    setScenarioToRun(null);
  };

  if (isLoading) {
    return <LoadingState message="Loading scenarios..." />;
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold">Scenarios</h1>
        <button className="btn-primary">Create Scenario</button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {scenarios?.map((scenario) => (
          <div key={scenario.id} className="card">
            <div className="flex items-start justify-between">
              <div>
                <h3 className="font-semibold text-lg">{scenario.name}</h3>
                <p className="text-sm text-gray-500 mt-1">{scenario.description}</p>
              </div>
              <button
                onClick={() => handleRunClick(scenario)}
                className="btn-primary text-sm"
              >
                Run
              </button>
            </div>

            <div className="mt-4">
              <p className="text-sm text-gray-500 mb-2">Phases</p>
              <div className="space-y-2">
                {scenario.phases.map((phase, idx) => (
                  <div key={`${phase.name}-${idx}`} className="flex items-center gap-2">
                    <span className="w-6 h-6 rounded-full bg-primary-100 text-primary-600 flex items-center justify-center text-xs font-medium">
                      {idx + 1}
                    </span>
                    <span className="text-sm">{phase.name}</span>
                    <span className="text-xs text-gray-400">
                      ({phase.techniques.length} techniques)
                    </span>
                  </div>
                ))}
              </div>
            </div>

            <div className="mt-4 pt-4 border-t border-gray-100">
              <div className="flex gap-1 flex-wrap">
                {scenario.tags?.map((tag) => (
                  <span key={tag} className="badge bg-gray-100 text-gray-700">
                    {tag}
                  </span>
                ))}
              </div>
            </div>
          </div>
        ))}
      </div>

      {scenarios?.length === 0 && (
        <EmptyState
          icon={DocumentTextIcon}
          title="No scenarios created"
          description="Create an attack scenario to test your defenses"
        />
      )}

      {/* Run Execution Modal */}
      {scenarioToRun && (
        <RunExecutionModal
          scenario={scenarioToRun}
          onConfirm={handleConfirmRun}
          onCancel={handleCancelRun}
          isLoading={startMutation.isPending}
        />
      )}
    </div>
  );
}
