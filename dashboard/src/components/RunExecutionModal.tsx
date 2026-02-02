import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { PlayIcon, ComputerDesktopIcon } from '@heroicons/react/24/outline';
import { api } from '../lib/api';
import { Agent, Scenario } from '../types';

interface RunExecutionModalProps {
  readonly scenario: Scenario;
  readonly onConfirm: (agentPaws: string[], safeMode: boolean) => void;
  readonly onCancel: () => void;
  readonly isLoading: boolean;
}

/**
 * Modal component for running an execution.
 * Allows selecting target agents and safe mode toggle.
 */
export function RunExecutionModal({ scenario, onConfirm, onCancel, isLoading }: Readonly<RunExecutionModalProps>) {
  const [selectedAgents, setSelectedAgents] = useState<string[]>([]);
  const [safeMode, setSafeMode] = useState(true);

  const { data: agents, isLoading: agentsLoading } = useQuery<Agent[]>({
    queryKey: ['agents'],
    queryFn: () => api.get('/agents').then(res => res.data),
  });

  const onlineAgents = agents?.filter(a => a.status === 'online') || [];

  const handleAgentToggle = (paw: string) => {
    setSelectedAgents(prev =>
      prev.includes(paw)
        ? prev.filter(p => p !== paw)
        : [...prev, paw]
    );
  };

  const handleSelectAll = () => {
    if (selectedAgents.length === onlineAgents.length) {
      setSelectedAgents([]);
    } else {
      setSelectedAgents(onlineAgents.map(a => a.paw));
    }
  };

  const handleSubmit = () => {
    if (selectedAgents.length > 0) {
      onConfirm(selectedAgents, safeMode);
    }
  };

  const techniqueCount = scenario.phases.reduce((acc, phase) => acc + phase.techniques.length, 0);

  return (
    <dialog open className="fixed inset-0 z-50 overflow-y-auto bg-transparent" aria-labelledby="modal-title" aria-modal="true">
      <div className="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        {/* Background overlay */}
        <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true" onClick={onCancel}></div>

        {/* Modal panel */}
        <div className="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <div className="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <div className="sm:flex sm:items-start">
              <div className="mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-primary-100 sm:mx-0 sm:h-10 sm:w-10">
                <PlayIcon className="h-6 w-6 text-primary-600" aria-hidden="true" />
              </div>
              <div className="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left flex-1">
                <h3 className="text-lg leading-6 font-medium text-gray-900" id="modal-title">
                  Run Scenario
                </h3>
                <div className="mt-2">
                  <div className="p-3 bg-gray-50 rounded-md">
                    <p className="text-sm font-medium text-gray-900">{scenario.name}</p>
                    <p className="text-sm text-gray-500">{scenario.phases.length} phases, {techniqueCount} techniques</p>
                  </div>
                </div>

                {/* Agent Selection */}
                <div className="mt-4">
                  <div className="flex justify-between items-center mb-2">
                    <span className="text-sm font-medium text-gray-700">Select Target Agents</span>
                    {onlineAgents.length > 0 && (
                      <button
                        type="button"
                        onClick={handleSelectAll}
                        className="text-xs text-primary-600 hover:text-primary-800"
                      >
                        {selectedAgents.length === onlineAgents.length ? 'Deselect All' : 'Select All'}
                      </button>
                    )}
                  </div>

                  {agentsLoading && (
                    <div className="text-sm text-gray-500">Loading agents...</div>
                  )}
                  {!agentsLoading && onlineAgents.length === 0 && (
                    <div className="text-sm text-amber-600 p-3 bg-amber-50 rounded-md">
                      No online agents available. Please ensure at least one agent is connected.
                    </div>
                  )}
                  {!agentsLoading && onlineAgents.length > 0 && (
                    <div className="space-y-2 max-h-48 overflow-y-auto border rounded-md p-2">
                      {onlineAgents.map(agent => (
                        <label
                          key={agent.paw}
                          className={`flex items-center p-2 rounded cursor-pointer hover:bg-gray-50 ${
                            selectedAgents.includes(agent.paw) ? 'bg-primary-50 border border-primary-200' : ''
                          }`}
                        >
                          <input
                            type="checkbox"
                            checked={selectedAgents.includes(agent.paw)}
                            onChange={() => handleAgentToggle(agent.paw)}
                            className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
                          />
                          <ComputerDesktopIcon className="h-5 w-5 text-gray-400 ml-3" />
                          <div className="ml-2">
                            <p className="text-sm font-medium text-gray-900">{agent.hostname}</p>
                            <p className="text-xs text-gray-500">{agent.platform} - {agent.paw.slice(0, 8)}...</p>
                          </div>
                        </label>
                      ))}
                    </div>
                  )}
                </div>

                {/* Safe Mode Toggle */}
                <div className="mt-4">
                  <label className="flex items-center cursor-pointer">
                    <input
                      type="checkbox"
                      checked={safeMode}
                      onChange={(e) => setSafeMode(e.target.checked)}
                      className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
                    />
                    <span className="ml-2 text-sm font-medium text-gray-700">Safe Mode</span>
                    <span className="ml-2 text-xs text-gray-500">(Skip potentially destructive techniques)</span>
                  </label>
                </div>

                {selectedAgents.length === 0 && onlineAgents.length > 0 && (
                  <p className="mt-3 text-sm text-amber-600">
                    Please select at least one agent to run the scenario.
                  </p>
                )}
              </div>
            </div>
          </div>
          <div className="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
            <button
              type="button"
              disabled={isLoading || selectedAgents.length === 0}
              onClick={handleSubmit}
              className="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-primary-600 text-base font-medium text-white hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isLoading && 'Starting...'}
              {!isLoading && `Run on ${selectedAgents.length} agent${selectedAgents.length === 1 ? '' : 's'}`}
            </button>
            <button
              type="button"
              disabled={isLoading}
              onClick={onCancel}
              className="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50"
            >
              Cancel
            </button>
          </div>
        </div>
      </div>
    </dialog>
  );
}
