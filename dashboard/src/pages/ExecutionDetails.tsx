import { useParams, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { ArrowLeftIcon, CheckCircleIcon, XCircleIcon, ExclamationTriangleIcon, ClockIcon } from '@heroicons/react/24/outline';
import { executionApi } from '../lib/api';
import { Execution, ExecutionResult } from '../types';
import { LoadingState } from '../components/LoadingState';
import { formatDistanceToNow, format } from 'date-fns';

/**
 * Returns the appropriate badge class and icon for a result status.
 *
 * Status meanings:
 * - success/successful: technique executed successfully (attack worked - bad for security)
 * - failed: technique execution failed (technical error - neutral)
 * - blocked: technique was blocked by security controls (good for security)
 * - detected: technique was detected by security tools (good for security)
 */
function getResultStatusInfo(status: string): { badgeClass: string; Icon: React.ElementType } {
  switch (status) {
    case 'success':
    case 'successful':
      // Attack succeeded = security vulnerability = danger (red)
      return { badgeClass: 'badge-danger', Icon: ExclamationTriangleIcon };
    case 'failed':
      // Technique execution failed (technical error) = neutral warning
      return { badgeClass: 'badge-warning', Icon: XCircleIcon };
    case 'detected':
      // Attack detected = security working = success (green)
      return { badgeClass: 'badge-success', Icon: CheckCircleIcon };
    case 'blocked':
      // Attack blocked = security working = success (green)
      return { badgeClass: 'badge-success', Icon: CheckCircleIcon };
    case 'pending':
    case 'running':
      return { badgeClass: 'badge-warning', Icon: ClockIcon };
    default:
      return { badgeClass: 'bg-gray-100 text-gray-800', Icon: ClockIcon };
  }
}

/**
 * Returns human-readable status label.
 */
function getStatusLabel(status: string): string {
  switch (status) {
    case 'success':
    case 'successful':
      return 'Attack Succeeded';
    case 'failed':
      return 'Execution Failed';
    case 'detected':
      return 'Detected';
    case 'blocked':
      return 'Blocked';
    case 'pending':
      return 'Pending';
    case 'running':
      return 'Running';
    default:
      return status;
  }
}

/**
 * Execution Details page component.
 * Displays detailed results of a scenario execution.
 */
export default function ExecutionDetails() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const { data: execution, isLoading: loadingExecution } = useQuery<Execution>({
    queryKey: ['execution', id],
    queryFn: () => executionApi.get(id!).then(res => res.data),
    enabled: !!id,
    refetchInterval: ({ state }) => {
      // Only poll if execution is still running
      const data = state.data;
      return data?.status === 'running' || data?.status === 'pending' ? 2000 : false;
    },
  });

  const { data: results, isLoading: loadingResults } = useQuery<ExecutionResult[]>({
    queryKey: ['execution-results', id],
    queryFn: () => executionApi.getResults(id!).then(res => res.data),
    enabled: !!id,
    refetchInterval: () => {
      // Poll while execution is active
      return execution?.status === 'running' || execution?.status === 'pending' ? 2000 : false;
    },
  });

  if (loadingExecution || loadingResults) {
    return <LoadingState message="Loading execution details..." />;
  }

  if (!execution) {
    return (
      <div className="text-center py-12">
        <h2 className="text-xl font-semibold text-gray-700">Execution not found</h2>
        <button onClick={() => navigate('/executions')} className="mt-4 btn-primary">
          Back to Executions
        </button>
      </div>
    );
  }

  const getExecutionStatusBadge = (status: string) => {
    switch (status) {
      case 'completed':
        return 'badge-success';
      case 'running':
      case 'pending':
        return 'badge-warning';
      case 'cancelled':
      case 'failed':
        return 'badge-danger';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  return (
    <div>
      {/* Header */}
      <div className="mb-8">
        <button
          onClick={() => navigate('/executions')}
          className="flex items-center text-gray-600 hover:text-gray-900 mb-4"
        >
          <ArrowLeftIcon className="h-5 w-5 mr-2" />
          Back to Executions
        </button>

        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-3xl font-bold">Execution Details</h1>
            <p className="text-gray-500 font-mono text-sm mt-1">{execution.id}</p>
          </div>
          <span className={`badge ${getExecutionStatusBadge(execution.status)}`}>
            {execution.status}
          </span>
        </div>
      </div>

      {/* Execution Summary */}
      <div className="card mb-8">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
          <div>
            <p className="text-sm text-gray-500">Scenario</p>
            <p className="font-semibold">{execution.scenario_id}</p>
          </div>
          <div>
            <p className="text-sm text-gray-500">Mode</p>
            <span className={`badge ${execution.safe_mode ? 'badge-success' : 'badge-danger'}`}>
              {execution.safe_mode ? 'Safe Mode' : 'Full Mode'}
            </span>
          </div>
          <div>
            <p className="text-sm text-gray-500">Started</p>
            <p className="font-semibold">
              {format(new Date(execution.started_at), 'PPpp')}
            </p>
            <p className="text-xs text-gray-400">
              {formatDistanceToNow(new Date(execution.started_at), { addSuffix: true })}
            </p>
          </div>
          <div>
            <p className="text-sm text-gray-500">Overall Score</p>
            <p className="text-3xl font-bold">
              {execution.score?.overall?.toFixed(1) || '0'}%
            </p>
          </div>
        </div>

        {/* Score Breakdown */}
        {execution.score && (
          <div className="mt-6 pt-6 border-t border-gray-200">
            <h3 className="text-sm font-medium text-gray-500 mb-4">Security Score Breakdown</h3>
            <div className="grid grid-cols-4 gap-4">
              <div className="text-center p-4 bg-green-50 rounded-lg">
                <p className="text-2xl font-bold text-green-700">{execution.score.blocked}</p>
                <p className="text-sm text-green-600">Blocked</p>
              </div>
              <div className="text-center p-4 bg-yellow-50 rounded-lg">
                <p className="text-2xl font-bold text-yellow-700">{execution.score.detected}</p>
                <p className="text-sm text-yellow-600">Detected</p>
              </div>
              <div className="text-center p-4 bg-red-50 rounded-lg">
                <p className="text-2xl font-bold text-red-700">{execution.score.successful}</p>
                <p className="text-sm text-red-600">Successful Attacks</p>
              </div>
              <div className="text-center p-4 bg-gray-50 rounded-lg">
                <p className="text-2xl font-bold text-gray-700">{execution.score.total}</p>
                <p className="text-sm text-gray-600">Total Tests</p>
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Results Table */}
      <div className="card">
        <h2 className="text-xl font-semibold mb-4">Technique Results</h2>

        {results && results.length > 0 ? (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Technique
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Agent
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Status
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Output
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {results.map((result) => {
                  const { badgeClass, Icon } = getResultStatusInfo(result.status);
                  return (
                    <tr key={result.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap">
                        <p className="font-medium font-mono">{result.technique_id}</p>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <p className="text-sm text-gray-600">{result.agent_paw}</p>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className={`badge ${badgeClass} flex items-center gap-1 w-fit`}>
                          <Icon className="h-4 w-4" />
                          {getStatusLabel(result.status)}
                        </span>
                      </td>
                      <td className="px-6 py-4">
                        {result.output ? (
                          <details className="cursor-pointer">
                            <summary className="text-sm text-blue-600 hover:text-blue-800">
                              View output
                            </summary>
                            <pre className="mt-2 p-3 bg-gray-900 text-gray-100 rounded text-xs overflow-x-auto max-h-48 overflow-y-auto">
                              {result.output}
                            </pre>
                          </details>
                        ) : (
                          <span className="text-gray-400 text-sm">No output</span>
                        )}
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        ) : (
          <div className="text-center py-8 text-gray-500">
            <ClockIcon className="h-12 w-12 mx-auto mb-4 text-gray-300" />
            <p>No results yet</p>
            <p className="text-sm">Results will appear here as techniques are executed</p>
          </div>
        )}
      </div>
    </div>
  );
}
