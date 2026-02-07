import { useState, useEffect } from 'react';
import { Technique, TacticType } from '../types';

/**
 * Ordered list of MITRE ATT&CK tactics for matrix display.
 */
const TACTICS: { id: TacticType; name: string }[] = [
  { id: 'reconnaissance', name: 'Reconnaissance' },
  { id: 'resource_development', name: 'Resource Development' },
  { id: 'initial_access', name: 'Initial Access' },
  { id: 'execution', name: 'Execution' },
  { id: 'persistence', name: 'Persistence' },
  { id: 'privilege_escalation', name: 'Privilege Escalation' },
  { id: 'defense_evasion', name: 'Defense Evasion' },
  { id: 'credential_access', name: 'Credential Access' },
  { id: 'discovery', name: 'Discovery' },
  { id: 'lateral_movement', name: 'Lateral Movement' },
  { id: 'collection', name: 'Collection' },
  { id: 'command_and_control', name: 'C2' },
  { id: 'exfiltration', name: 'Exfiltration' },
  { id: 'impact', name: 'Impact' },
];

/**
 * Tactic header background colors.
 */
const tacticHeaderColors: Record<TacticType, string> = {
  reconnaissance: 'bg-purple-600',
  resource_development: 'bg-blue-600',
  initial_access: 'bg-red-600',
  execution: 'bg-orange-600',
  persistence: 'bg-yellow-600',
  privilege_escalation: 'bg-pink-600',
  defense_evasion: 'bg-green-600',
  credential_access: 'bg-indigo-600',
  discovery: 'bg-cyan-600',
  lateral_movement: 'bg-teal-600',
  collection: 'bg-lime-600',
  command_and_control: 'bg-amber-600',
  exfiltration: 'bg-rose-600',
  impact: 'bg-red-700',
};

/**
 * Tactic cell background colors (lighter for light mode, darker for dark mode).
 */
const tacticCellColors: Record<TacticType, string> = {
  reconnaissance: 'bg-purple-50 hover:bg-purple-100 border-purple-200 dark:bg-purple-900/30 dark:hover:bg-purple-900/50 dark:border-purple-700',
  resource_development: 'bg-blue-50 hover:bg-blue-100 border-blue-200 dark:bg-blue-900/30 dark:hover:bg-blue-900/50 dark:border-blue-700',
  initial_access: 'bg-red-50 hover:bg-red-100 border-red-200 dark:bg-red-900/30 dark:hover:bg-red-900/50 dark:border-red-700',
  execution: 'bg-orange-50 hover:bg-orange-100 border-orange-200 dark:bg-orange-900/30 dark:hover:bg-orange-900/50 dark:border-orange-700',
  persistence: 'bg-yellow-50 hover:bg-yellow-100 border-yellow-200 dark:bg-yellow-900/30 dark:hover:bg-yellow-900/50 dark:border-yellow-700',
  privilege_escalation: 'bg-pink-50 hover:bg-pink-100 border-pink-200 dark:bg-pink-900/30 dark:hover:bg-pink-900/50 dark:border-pink-700',
  defense_evasion: 'bg-green-50 hover:bg-green-100 border-green-200 dark:bg-green-900/30 dark:hover:bg-green-900/50 dark:border-green-700',
  credential_access: 'bg-indigo-50 hover:bg-indigo-100 border-indigo-200 dark:bg-indigo-900/30 dark:hover:bg-indigo-900/50 dark:border-indigo-700',
  discovery: 'bg-cyan-50 hover:bg-cyan-100 border-cyan-200 dark:bg-cyan-900/30 dark:hover:bg-cyan-900/50 dark:border-cyan-700',
  lateral_movement: 'bg-teal-50 hover:bg-teal-100 border-teal-200 dark:bg-teal-900/30 dark:hover:bg-teal-900/50 dark:border-teal-700',
  collection: 'bg-lime-50 hover:bg-lime-100 border-lime-200 dark:bg-lime-900/30 dark:hover:bg-lime-900/50 dark:border-lime-700',
  command_and_control: 'bg-amber-50 hover:bg-amber-100 border-amber-200 dark:bg-amber-900/30 dark:hover:bg-amber-900/50 dark:border-amber-700',
  exfiltration: 'bg-rose-50 hover:bg-rose-100 border-rose-200 dark:bg-rose-900/30 dark:hover:bg-rose-900/50 dark:border-rose-700',
  impact: 'bg-red-50 hover:bg-red-100 border-red-200 dark:bg-red-900/30 dark:hover:bg-red-900/50 dark:border-red-700',
};

interface MitreMatrixProps {
  readonly techniques: Technique[];
  readonly onTechniqueClick?: (technique: Technique) => void;
}

/**
 * MITRE ATT&CK Matrix visualization component.
 * Displays techniques organized by tactic in a grid layout.
 */
export function MitreMatrix({ techniques, onTechniqueClick }: Readonly<MitreMatrixProps>) {
  const [platformFilter, setPlatformFilter] = useState<string>('all');
  const [selectedTechnique, setSelectedTechnique] = useState<Technique | null>(null);

  // Get all tactics for a technique (multi-tactic support with fallback)
  const getTactics = (t: Technique): TacticType[] => {
    if (t.tactics && t.tactics.length > 0) {
      return t.tactics.map(tac => String(tac).replaceAll('-', '_') as TacticType);
    }
    return [String(t.tactic).replaceAll('-', '_') as TacticType];
  };

  // Group techniques by tactic (a technique can appear in multiple columns)
  const techniquesByTactic = TACTICS.reduce((acc, tactic) => {
    acc[tactic.id] = techniques
      .filter(t => {
        const platformMatches = platformFilter === 'all' || t.platforms.includes(platformFilter);
        if (!platformMatches) return false;
        return getTactics(t).includes(tactic.id);
      })
      .sort((a, b) => a.id.localeCompare(b.id));
    return acc;
  }, {} as Record<TacticType, Technique[]>);

  const handleTechniqueClick = (technique: Technique) => {
    setSelectedTechnique(technique);
    onTechniqueClick?.(technique);
  };

  const handleCloseDetail = () => {
    setSelectedTechnique(null);
  };

  // Handle Escape key to close modal at document level
  useEffect(() => {
    if (!selectedTechnique) return;

    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        handleCloseDetail();
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [selectedTechnique]);

  // Get unique platforms from techniques
  const platforms = [...new Set(techniques.flatMap(t => t.platforms))].sort((a, b) => a.localeCompare(b));

  return (
    <div className="space-y-4">
      {/* Filters */}
      <div className="flex items-center gap-4">
        <label htmlFor="platform-filter" className="text-sm font-medium text-gray-700 dark:text-gray-300">Platform:</label>
        <select
          id="platform-filter"
          value={platformFilter}
          onChange={(e) => setPlatformFilter(e.target.value)}
          className="input py-1 px-2 text-sm"
        >
          <option value="all">All Platforms</option>
          {platforms.map(p => (
            <option key={p} value={p}>{p}</option>
          ))}
        </select>
        <span className="text-sm text-gray-500 dark:text-gray-400">
          {techniques.length} techniques loaded
        </span>
      </div>

      {/* Matrix Grid */}
      <div className="overflow-x-auto">
        <div className="inline-grid" style={{ gridTemplateColumns: `repeat(${TACTICS.length}, minmax(120px, 1fr))` }}>
          {/* Header Row */}
          {TACTICS.map(tactic => (
            <div
              key={tactic.id}
              className={`${tacticHeaderColors[tactic.id]} text-white text-xs font-semibold p-2 text-center border-r border-white/20`}
              title={tactic.name}
            >
              <div className="truncate">{tactic.name}</div>
              <div className="text-white/70 text-xs mt-1">
                {techniquesByTactic[tactic.id]?.length || 0}
              </div>
            </div>
          ))}

          {/* Technique Cells */}
          {TACTICS.map(tactic => (
            <div key={`col-${tactic.id}`} className="flex flex-col gap-1 p-1 bg-gray-50 dark:bg-gray-800 min-h-[200px]">
              {techniquesByTactic[tactic.id]?.map(technique => (
                <button
                  key={technique.id}
                  onClick={() => handleTechniqueClick(technique)}
                  className={`${tacticCellColors[tactic.id]} p-2 rounded border text-left transition-all cursor-pointer`}
                  title={`${technique.id}: ${technique.name}`}
                >
                  <div className="text-xs font-mono text-gray-600 dark:text-gray-400">{technique.id}</div>
                  <div className="text-xs font-medium text-gray-900 dark:text-gray-100 truncate">{technique.name}</div>
                  <div className="flex gap-1 mt-1">
                    {technique.is_safe ? (
                      <span className="w-2 h-2 rounded-full bg-green-500" title="Safe" />
                    ) : (
                      <span className="w-2 h-2 rounded-full bg-red-500" title="Unsafe" />
                    )}
                  </div>
                </button>
              ))}
              {techniquesByTactic[tactic.id]?.length === 0 && (
                <div className="text-xs text-gray-400 dark:text-gray-500 text-center py-4">
                  No techniques
                </div>
              )}
            </div>
          ))}
        </div>
      </div>

      {/* Technique Detail Panel */}
      {selectedTechnique && (
        <div className="fixed inset-0 z-50 overflow-y-auto">
          <div className="flex items-center justify-center min-h-screen p-4">
            <button
              type="button"
              className="fixed inset-0 bg-gray-500 bg-opacity-75 dark:bg-black dark:bg-opacity-60 cursor-default border-none"
              onClick={handleCloseDetail}
              aria-label="Close modal"
            />
            <div className="relative bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-lg w-full p-6">
              <button
                onClick={handleCloseDetail}
                className="absolute top-4 right-4 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
              >
                <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
              <div className="flex items-start gap-3">
                <div className={`${tacticHeaderColors[String(selectedTechnique.tactic).replaceAll('-', '_') as TacticType] ?? 'bg-gray-600'} text-white px-2 py-1 rounded text-xs font-semibold`}>
                  {selectedTechnique.id}
                </div>
                <div>
                  <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100">{selectedTechnique.name}</h3>
                  <p className="text-sm text-gray-500 dark:text-gray-400 capitalize">
                    {getTactics(selectedTechnique).map(t => t.replaceAll('_', ' ')).join(', ')}
                  </p>
                </div>
              </div>
              <p className="mt-4 text-sm text-gray-600 dark:text-gray-300">{selectedTechnique.description}</p>
              <div className="mt-4 flex flex-wrap gap-2">
                <span className={`badge ${selectedTechnique.is_safe ? 'badge-success' : 'badge-danger'}`}>
                  {selectedTechnique.is_safe ? 'Safe' : 'Unsafe'}
                </span>
                {selectedTechnique.platforms.map(p => (
                  <span key={p} className="badge bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300">{p}</span>
                ))}
              </div>
              {selectedTechnique.detection && selectedTechnique.detection.length > 0 && (
                <div className="mt-4">
                  <h4 className="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2">Detection</h4>
                  <ul className="text-sm text-gray-600 dark:text-gray-400 space-y-1">
                    {selectedTechnique.detection.map((d) => (
                      <li key={`${d.source}-${d.indicator}`} className="flex gap-2">
                        <span className="font-medium">{d.source}:</span>
                        <span>{d.indicator}</span>
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Legend */}
      <div className="flex items-center gap-4 text-xs text-gray-500 dark:text-gray-400">
        <div className="flex items-center gap-1">
          <span className="w-3 h-3 rounded-full bg-green-500"></span>
          <span>Safe technique</span>
        </div>
        <div className="flex items-center gap-1">
          <span className="w-3 h-3 rounded-full bg-red-500"></span>
          <span>Potentially unsafe</span>
        </div>
      </div>
    </div>
  );
}
