/**
 * Centralized type definitions for the AutoStrike dashboard.
 * These types mirror the API response structures from the Go backend.
 */

/**
 * Agent status enumeration.
 */
export type AgentStatus = 'online' | 'offline' | 'unknown';

/**
 * Represents a connected agent in the AutoStrike network.
 */
export interface Agent {
  /** Unique identifier for the agent */
  paw: string;
  /** Hostname of the machine running the agent */
  hostname: string;
  /** Username under which the agent is running */
  username: string;
  /** Operating system platform (windows, linux, darwin) */
  platform: string;
  /** Available command executors (powershell, cmd, bash, etc.) */
  executors: string[];
  /** Current connection status */
  status: AgentStatus;
  /** ISO timestamp of last heartbeat */
  last_seen: string;
  /** ISO timestamp of agent registration */
  created_at?: string;
}

/**
 * MITRE ATT&CK tactic types.
 */
export type TacticType =
  | 'reconnaissance'
  | 'resource_development'
  | 'initial_access'
  | 'execution'
  | 'persistence'
  | 'privilege_escalation'
  | 'defense_evasion'
  | 'credential_access'
  | 'discovery'
  | 'lateral_movement'
  | 'collection'
  | 'command_and_control'
  | 'exfiltration'
  | 'impact';

/**
 * Represents a MITRE ATT&CK technique.
 */
export interface Technique {
  /** MITRE technique ID (e.g., T1059) */
  id: string;
  /** Human-readable technique name */
  name: string;
  /** Detailed description of the technique */
  description: string;
  /** MITRE tactic this technique belongs to (primary) */
  tactic: TacticType;
  /** All MITRE tactics this technique belongs to (multi-tactic support) */
  tactics?: TacticType[];
  /** Supported operating system platforms */
  platforms: string[];
  /** Whether the technique is safe for production testing */
  is_safe: boolean;
  /** Available executors for this technique */
  executors?: TechniqueExecutor[];
  /** Detection indicators for this technique */
  detection?: DetectionIndicator[];
}

/**
 * Executor configuration for a technique.
 */
export interface TechniqueExecutor {
  /** Executor name (for distinguishing multiple executors) */
  name?: string;
  /** Executor type (cmd, powershell, bash, etc.) */
  type: string;
  /** Target platform (windows, linux, macos) */
  platform?: string;
  /** Command to execute */
  command: string;
  /** Cleanup command to run after execution */
  cleanup?: string;
  /** Execution timeout in seconds */
  timeout: number;
  /** Whether elevated privileges are required */
  elevation_required?: boolean;
}

/**
 * Technique selection with optional executor preference.
 */
export interface TechniqueSelection {
  /** Technique ID */
  technique_id: string;
  /** Preferred executor name (empty = auto-select best match) */
  executor_name?: string;
}

/**
 * Detection indicator for a technique.
 */
export interface DetectionIndicator {
  /** Detection data source */
  source: string;
  /** Indicator description */
  indicator: string;
}

/**
 * Execution status enumeration.
 */
export type ExecutionStatus = 'pending' | 'running' | 'completed' | 'failed' | 'cancelled';

/**
 * Security score breakdown from an execution.
 */
export interface ExecutionScore {
  /** Overall security score (0-100) */
  overall: number;
  /** Number of blocked techniques */
  blocked: number;
  /** Number of detected techniques */
  detected: number;
  /** Number of successful (undetected) techniques */
  successful: number;
  /** Total number of techniques tested */
  total: number;
}

/**
 * Represents a scenario execution.
 */
export interface Execution {
  /** Unique execution identifier */
  id: string;
  /** ID of the scenario being executed */
  scenario_id: string;
  /** Current execution status */
  status: ExecutionStatus;
  /** ISO timestamp when execution started */
  started_at: string;
  /** ISO timestamp when execution completed */
  completed_at?: string;
  /** Whether execution ran in safe mode */
  safe_mode: boolean;
  /** Security score results */
  score?: ExecutionScore;
}

/**
 * Represents a phase within a scenario.
 */
export interface ScenarioPhase {
  /** Phase name */
  name: string;
  /** List of technique selections or IDs to execute in this phase */
  techniques: TechniqueSelection[] | string[];
}

/**
 * Represents an attack scenario.
 */
export interface Scenario {
  /** Unique scenario identifier */
  id: string;
  /** Scenario name */
  name: string;
  /** Detailed scenario description */
  description: string;
  /** Ordered list of execution phases */
  phases: ScenarioPhase[];
  /** Tags for categorization */
  tags: string[];
}

/**
 * Result of a single technique execution.
 */
export interface ExecutionResult {
  /** Unique result identifier */
  id: string;
  /** Parent execution ID */
  execution_id: string;
  /** Technique that was executed */
  technique_id: string;
  /** Agent that executed the technique */
  agent_paw: string;
  /** Result status */
  status: 'blocked' | 'detected' | 'successful' | 'failed' | 'skipped';
  /** Command output */
  output: string;
  /** Whether the technique was detected */
  detected: boolean;
  /** ISO timestamp when execution started */
  start_time: string;
  /** ISO timestamp when execution ended */
  end_time: string;
}

/**
 * Coverage statistics by tactic.
 */
export type TacticCoverage = Record<string, number>;

/**
 * API error response structure.
 */
export interface ApiError {
  /** Error message */
  error: string;
  /** HTTP status code */
  status?: number;
}
