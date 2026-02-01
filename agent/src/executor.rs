//! Command execution with timeout support.

use std::process::Stdio;
use std::time::Duration;
use tokio::process::Command;
use tokio::time::timeout;
use tracing::{debug, error};

/// Result of a command execution.
pub struct ExecutionResult {
    /// Whether the command executed successfully.
    pub success: bool,
    /// Combined stdout and stderr output.
    pub output: String,
    /// Process exit code, if available.
    pub exit_code: Option<i32>,
}

/// Executes commands using platform-specific shells.
pub struct CommandExecutor;

impl CommandExecutor {
    /// Creates a new command executor instance.
    pub fn new() -> Self {
        Self
    }

    /// Executes a command with the specified executor and timeout.
    pub async fn execute(
        &self,
        executor_type: &str,
        command: &str,
        time_limit: Duration,
    ) -> ExecutionResult {
        debug!("Executing command with {}: {}", executor_type, command);

        let result = timeout(time_limit, self.run_command(executor_type, command)).await;

        match result {
            Ok(exec_result) => exec_result,
            Err(_) => ExecutionResult {
                success: false,
                output: "Command timed out".to_string(),
                exit_code: None,
            },
        }
    }

    async fn run_command(&self, executor_type: &str, command: &str) -> ExecutionResult {
        let mut cmd = self.build_command(executor_type, command);

        cmd.stdout(Stdio::piped()).stderr(Stdio::piped());

        match cmd.output().await {
            Ok(output) => {
                let stdout = String::from_utf8_lossy(&output.stdout);
                let stderr = String::from_utf8_lossy(&output.stderr);
                let combined = format!("{}{}", stdout, stderr);

                ExecutionResult {
                    success: output.status.success(),
                    output: combined.trim().to_string(),
                    exit_code: output.status.code(),
                }
            }
            Err(e) => {
                error!("Command execution failed: {}", e);
                ExecutionResult {
                    success: false,
                    output: format!("Execution error: {}", e),
                    exit_code: None,
                }
            }
        }
    }

    #[cfg(target_os = "windows")]
    fn build_command(&self, executor_type: &str, command: &str) -> Command {
        let cmd = match executor_type {
            "powershell" | "ps" => {
                let mut c = Command::new("powershell.exe");
                c.args(["-NoProfile", "-NonInteractive", "-Command", command]);
                c
            }
            "pwsh" | "powershell7" => {
                let mut c = Command::new("pwsh.exe");
                c.args(["-NoProfile", "-NonInteractive", "-Command", command]);
                c
            }
            "cmd" => {
                let mut c = Command::new("cmd.exe");
                c.args(["/C", command]);
                c
            }
            _ => {
                let mut c = Command::new("powershell.exe");
                c.args(["-NoProfile", "-NonInteractive", "-Command", command]);
                c
            }
        };
        cmd
    }

    #[cfg(not(target_os = "windows"))]
    fn build_command(&self, executor_type: &str, command: &str) -> Command {
        let shell = match executor_type {
            "bash" => "/bin/bash",
            "zsh" => "/bin/zsh",
            "sh" => "/bin/sh",
            _ => "/bin/sh",
        };

        let mut cmd = Command::new(shell);
        cmd.args(["-c", command]);
        cmd
    }
}

impl Default for CommandExecutor {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_executor_default() {
        let _executor = CommandExecutor;
    }

    #[test]
    fn test_executor_new() {
        let executor = CommandExecutor::new();
        let _ = executor;
    }

    #[test]
    fn test_execution_result_struct() {
        let result = ExecutionResult {
            success: true,
            output: "test output".to_string(),
            exit_code: Some(0),
        };
        assert!(result.success);
        assert_eq!(result.output, "test output");
        assert_eq!(result.exit_code, Some(0));
    }

    #[test]
    fn test_execution_result_failure() {
        let result = ExecutionResult {
            success: false,
            output: "error message".to_string(),
            exit_code: Some(1),
        };
        assert!(!result.success);
        assert_eq!(result.exit_code, Some(1));
    }

    #[test]
    fn test_execution_result_no_exit_code() {
        let result = ExecutionResult {
            success: false,
            output: "timed out".to_string(),
            exit_code: None,
        };
        assert!(!result.success);
        assert!(result.exit_code.is_none());
    }

    #[tokio::test]
    async fn test_simple_command() {
        let executor = CommandExecutor::new();

        #[cfg(not(target_os = "windows"))]
        let result = executor
            .execute("sh", "echo hello", Duration::from_secs(5))
            .await;

        #[cfg(target_os = "windows")]
        let result = executor
            .execute("cmd", "echo hello", Duration::from_secs(5))
            .await;

        assert!(result.success);
        assert!(result.output.contains("hello"));
    }

    #[tokio::test]
    async fn test_command_with_exit_code() {
        let executor = CommandExecutor::new();

        #[cfg(not(target_os = "windows"))]
        let result = executor
            .execute("sh", "exit 0", Duration::from_secs(5))
            .await;

        #[cfg(target_os = "windows")]
        let result = executor
            .execute("cmd", "exit /b 0", Duration::from_secs(5))
            .await;

        assert!(result.success);
        assert_eq!(result.exit_code, Some(0));
    }

    #[tokio::test]
    async fn test_failed_command() {
        let executor = CommandExecutor::new();

        #[cfg(not(target_os = "windows"))]
        let result = executor
            .execute("sh", "exit 1", Duration::from_secs(5))
            .await;

        #[cfg(target_os = "windows")]
        let result = executor
            .execute("cmd", "exit /b 1", Duration::from_secs(5))
            .await;

        assert!(!result.success);
        assert_eq!(result.exit_code, Some(1));
    }

    #[tokio::test]
    async fn test_command_timeout() {
        let executor = CommandExecutor::new();

        #[cfg(not(target_os = "windows"))]
        let result = executor
            .execute("sh", "sleep 10", Duration::from_millis(100))
            .await;

        #[cfg(target_os = "windows")]
        let result = executor
            .execute("cmd", "ping -n 10 127.0.0.1", Duration::from_millis(100))
            .await;

        assert!(!result.success);
        assert!(result.output.contains("timed out"));
        assert!(result.exit_code.is_none());
    }

    #[tokio::test]
    async fn test_bash_executor() {
        let executor = CommandExecutor::new();

        #[cfg(not(target_os = "windows"))]
        {
            let result = executor
                .execute("bash", "echo $SHELL", Duration::from_secs(5))
                .await;
            assert!(result.success);
        }
    }

    #[tokio::test]
    async fn test_default_executor_fallback() {
        let executor = CommandExecutor::new();

        #[cfg(not(target_os = "windows"))]
        {
            let result = executor
                .execute("unknown_executor", "echo fallback", Duration::from_secs(5))
                .await;
            assert!(result.success);
            assert!(result.output.contains("fallback"));
        }
    }

    #[tokio::test]
    async fn test_command_with_stderr() {
        let executor = CommandExecutor::new();

        #[cfg(not(target_os = "windows"))]
        let result = executor
            .execute("sh", "echo error >&2", Duration::from_secs(5))
            .await;

        #[cfg(target_os = "windows")]
        let result = executor
            .execute("cmd", "echo error 1>&2", Duration::from_secs(5))
            .await;

        assert!(result.output.contains("error"));
    }

    #[tokio::test]
    async fn test_multiline_output() {
        let executor = CommandExecutor::new();

        #[cfg(not(target_os = "windows"))]
        let result = executor
            .execute("sh", "echo line1; echo line2", Duration::from_secs(5))
            .await;

        #[cfg(target_os = "windows")]
        let result = executor
            .execute("cmd", "echo line1 & echo line2", Duration::from_secs(5))
            .await;

        assert!(result.success);
        assert!(result.output.contains("line1"));
        assert!(result.output.contains("line2"));
    }

    #[tokio::test]
    async fn test_zsh_executor() {
        let executor = CommandExecutor::new();

        #[cfg(not(target_os = "windows"))]
        {
            let result = executor
                .execute("zsh", "echo zsh_test", Duration::from_secs(5))
                .await;
            // This may succeed or fail depending on if zsh is installed
            let _ = result;
        }
    }
}
