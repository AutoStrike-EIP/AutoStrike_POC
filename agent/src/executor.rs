//! Command execution with timeout support.

use std::process::Stdio;
use std::sync::atomic::{AtomicUsize, Ordering};
use std::sync::Arc;
use std::time::Duration;
use tokio::io::AsyncReadExt;
use tokio::process::Command;
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

/// Maximum output size in bytes (1 MB) to prevent memory exhaustion.
const MAX_OUTPUT_SIZE: usize = 1_048_576;

/// Executes commands using platform-specific shells.
pub struct CommandExecutor;

impl CommandExecutor {
    /// Creates a new command executor instance.
    pub fn new() -> Self {
        Self
    }

    /// Executes a command with the specified executor and timeout.
    /// On timeout, the child process is actively killed.
    pub async fn execute(
        &self,
        executor_type: &str,
        command: &str,
        time_limit: Duration,
    ) -> ExecutionResult {
        debug!("Executing command with {}: {}", executor_type, command);

        let mut cmd = self.build_command(executor_type, command);
        cmd.stdout(Stdio::piped()).stderr(Stdio::piped());

        let mut child = match cmd.spawn() {
            Ok(child) => child,
            Err(e) => {
                error!("Failed to spawn command: {}", e);
                return ExecutionResult {
                    success: false,
                    output: format!("Execution error: {}", e),
                    exit_code: None,
                };
            }
        };

        // Take ownership of stdout/stderr for concurrent reads
        let stdout = child.stdout.take().expect("stdout piped");
        let stderr = child.stderr.take().expect("stderr piped");

        // Shared byte budget to cap total output across both streams
        let budget = Arc::new(AtomicUsize::new(MAX_OUTPUT_SIZE));

        // Drain stdout concurrently
        let stdout_fut = {
            let budget = budget.clone();
            async move { drain_stream(stdout, &budget).await }
        };

        // Drain stderr concurrently
        let stderr_fut = {
            let budget = budget.clone();
            async move { drain_stream(stderr, &budget).await }
        };

        // Both streams are polled concurrently via join!, preventing pipe deadlocks
        let read_output = async {
            let (stdout_buf, stderr_buf) = tokio::join!(stdout_fut, stderr_fut);
            let truncated = budget.load(Ordering::Relaxed) == 0;
            let stdout_str = String::from_utf8_lossy(&stdout_buf);
            let stderr_str = String::from_utf8_lossy(&stderr_buf);
            let combined = format!("{}{}", stdout_str, stderr_str);
            let mut output = combined.trim().to_string();
            if truncated {
                // Safe UTF-8 truncation
                let safe_boundary = find_char_boundary(&output, MAX_OUTPUT_SIZE);
                output.truncate(safe_boundary);
                output.push_str("\n... [output truncated]");
            }
            output
        };

        // Race output collection against timeout, kill child on timeout
        tokio::select! {
            output = read_output => {
                // Output collected, now wait for child to exit
                match child.wait().await {
                    Ok(status) => ExecutionResult {
                        success: status.success(),
                        output,
                        exit_code: status.code(),
                    },
                    Err(e) => {
                        error!("Failed to wait for child: {}", e);
                        ExecutionResult {
                            success: false,
                            output,
                            exit_code: None,
                        }
                    }
                }
            }
            _ = tokio::time::sleep(time_limit) => {
                // Timeout: kill the child process
                let _ = child.kill().await;
                let _ = child.wait().await; // Reap the zombie
                ExecutionResult {
                    success: false,
                    output: "Command timed out".to_string(),
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

/// Drains an async reader into a Vec, claiming bytes from a shared atomic budget.
/// Returns the collected bytes. Stops when the stream is exhausted or the budget is depleted.
async fn drain_stream<R: tokio::io::AsyncRead + Unpin>(
    mut stream: R,
    budget: &AtomicUsize,
) -> Vec<u8> {
    let mut buf = Vec::new();
    let mut chunk = [0u8; 8192];
    loop {
        if budget.load(Ordering::Relaxed) == 0 {
            break;
        }
        match stream.read(&mut chunk).await {
            Ok(0) => break,
            Ok(n) => {
                let claimed = claim_budget(budget, n);
                if claimed > 0 {
                    buf.extend_from_slice(&chunk[..claimed]);
                }
                if claimed < n {
                    break; // Budget exhausted
                }
            }
            Err(_) => break,
        }
    }
    buf
}

/// Atomically claims up to `want` bytes from the shared budget.
/// Returns the number of bytes actually claimed.
fn claim_budget(budget: &AtomicUsize, want: usize) -> usize {
    loop {
        let current = budget.load(Ordering::Relaxed);
        if current == 0 {
            return 0;
        }
        let claim = want.min(current);
        match budget.compare_exchange_weak(
            current,
            current - claim,
            Ordering::Relaxed,
            Ordering::Relaxed,
        ) {
            Ok(_) => return claim,
            Err(_) => continue,
        }
    }
}

/// Finds the largest valid UTF-8 char boundary at or before `max` bytes.
/// Prevents panics when slicing multi-byte characters.
fn find_char_boundary(s: &str, max: usize) -> usize {
    if max >= s.len() {
        return s.len();
    }
    let mut boundary = max;
    while boundary > 0 && !s.is_char_boundary(boundary) {
        boundary -= 1;
    }
    boundary
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
