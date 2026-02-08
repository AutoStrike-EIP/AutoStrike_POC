//! Post-execution output file capture.
//!
//! Many Atomic Red Team commands redirect their output to files (e.g. `ps >> /tmp/loot.txt`)
//! instead of stdout. This module detects such redirections, reads the output files after
//! command execution, and returns their content as the task output.

use regex::Regex;
use std::path::{Path, PathBuf};
use tracing::debug;

/// Threshold in bytes below which we attempt file-based output capture.
const OUTPUT_THRESHOLD: usize = 50;

/// Maximum total bytes to read from output files (1 MB).
const MAX_FILE_READ_SIZE: u64 = 1_048_576;

/// Collects unique regex group-1 matches from a command string.
fn collect_matches(command: &str, patterns: &[&str]) -> Vec<String> {
    let mut paths = Vec::new();
    for pattern in patterns {
        let re = Regex::new(pattern).unwrap();
        for cap in re.captures_iter(command) {
            if let Some(m) = cap.get(1) {
                let path = m.as_str().to_string();
                if !paths.contains(&path) {
                    paths.push(path);
                }
            }
        }
    }
    paths
}

/// Extracts output file paths from a command string based on the executor type.
///
/// Only captures redirections to safe directories (`/tmp/` on Unix, `%TEMP%`/`$env:TEMP` on Windows)
/// to avoid reading arbitrary files.
pub fn extract_output_paths(command: &str, executor_type: &str) -> Vec<String> {
    match executor_type {
        "sh" | "bash" | "zsh" => collect_matches(command, &[r#">{1,2}\s*(/tmp/[^\s;|&"']+)"#]),
        "cmd" => collect_matches(
            command,
            &[r#"(?i)>{1,2}\s*(%temp%\\[^\s;|&"']+|%userprofile%\\[^\s;|&"']+)"#],
        ),
        "powershell" | "ps" | "pwsh" | "powershell7" => collect_matches(
            command,
            &[
                r#"(?i)>{1,2}\s*(\$env:TEMP\\[^\s;|&"']+|\$env:USERPROFILE\\[^\s;|&"']+)"#,
                r#"(?i)Out-File\s+(?:-FilePath\s+)?(\$env:TEMP\\[^\s;|&"']+|\$env:USERPROFILE\\[^\s;|&"']+)"#,
                r#"(?i)Set-Content\s+(?:-Path\s+)?(\$env:TEMP\\[^\s;|&"']+|\$env:USERPROFILE\\[^\s;|&"']+)"#,
                r#">{1,2}\s*(/tmp/[^\s;|&"']+)"#,
            ],
        ),
        _ => Vec::new(),
    }
}

/// Looks up the system temp directory from environment variables.
fn lookup_temp_dir() -> Option<String> {
    std::env::var("TEMP")
        .or_else(|_| std::env::var("TMP"))
        .or_else(|_| std::env::var("TMPDIR"))
        .ok()
}

/// Looks up the user home directory from environment variables.
fn lookup_home_dir() -> Option<String> {
    std::env::var("USERPROFILE")
        .or_else(|_| std::env::var("HOME"))
        .ok()
}

/// Normalizes a path by resolving `.` and `..` components logically (no filesystem access).
/// Prevents path traversal attacks like `/tmp/../../etc/passwd` → `/etc/passwd`.
fn normalize_path(path: &Path) -> PathBuf {
    use std::path::Component;
    let mut components: Vec<Component> = Vec::new();
    for component in path.components() {
        match component {
            Component::ParentDir => {
                // Only pop Normal components — never pop past RootDir or Prefix
                if matches!(components.last(), Some(Component::Normal(_))) {
                    components.pop();
                }
            }
            Component::CurDir => {}
            c => components.push(c),
        }
    }
    components.iter().collect()
}

/// Resolves a suffix relative to the system temp directory, normalizing to prevent traversal.
/// Returns `None` if the resolved path escapes the temp directory.
fn resolve_temp_suffix(suffix: &str) -> Option<PathBuf> {
    let base = lookup_temp_dir().unwrap_or_else(|| "/tmp".to_string());
    let resolved = normalize_path(&PathBuf::from(base).join(suffix));
    let temp = normalize_path(&std::env::temp_dir());
    if resolved.starts_with(&temp) || resolved.starts_with("/tmp") {
        return Some(resolved);
    }
    None
}

/// Resolves a suffix relative to the user home directory.
fn resolve_home_suffix(suffix: &str) -> Option<PathBuf> {
    lookup_home_dir().map(|home| PathBuf::from(home).join(suffix))
}

/// Resolves a raw path string (potentially containing environment variables) to an absolute path.
/// Unix paths are normalized and validated to stay within `/tmp/`.
/// Environment variable paths are normalized to prevent traversal via `..` in suffixes.
pub fn resolve_path(raw_path: &str) -> Option<PathBuf> {
    // Unix absolute paths — normalize and verify /tmp/ containment
    if raw_path.starts_with('/') {
        let normalized = normalize_path(&PathBuf::from(raw_path));
        return if normalized.starts_with("/tmp") {
            Some(normalized)
        } else {
            None
        };
    }

    // Windows absolute paths (e.g. C:\Users\...)
    if raw_path.len() >= 3
        && raw_path.as_bytes()[1] == b':'
        && raw_path.as_bytes()[0].is_ascii_alphabetic()
    {
        return Some(PathBuf::from(raw_path));
    }

    let lower = raw_path.to_lowercase();

    // PowerShell $env:TEMP\... or Windows cmd %temp%\...
    if lower.starts_with("$env:temp\\") || lower.starts_with("$env:temp/") {
        return resolve_temp_suffix(&raw_path[10..]);
    }
    if lower.starts_with("%temp%\\") || lower.starts_with("%temp%/") {
        return resolve_temp_suffix(&raw_path[7..]);
    }

    // PowerShell $env:USERPROFILE\... or Windows cmd %userprofile%\...
    if lower.starts_with("$env:userprofile\\") || lower.starts_with("$env:userprofile/") {
        return resolve_home_suffix(&raw_path[17..]);
    }
    if lower.starts_with("%userprofile%\\") || lower.starts_with("%userprofile%/") {
        return resolve_home_suffix(&raw_path[14..]);
    }

    None
}

/// Checks that a resolved path is inside a safe directory (temp dir) after
/// canonicalization, preventing path traversal attacks (e.g. `/tmp/../../etc/passwd`).
fn is_safe_path(resolved: &Path) -> bool {
    let canonical = match std::fs::canonicalize(resolved) {
        Ok(p) => p,
        Err(_) => return false, // file doesn't exist or can't be resolved
    };

    let safe_dirs = [std::env::temp_dir(), PathBuf::from("/tmp")];

    safe_dirs.iter().any(|dir| {
        if let Ok(canon_dir) = std::fs::canonicalize(dir) {
            canonical.starts_with(&canon_dir)
        } else {
            false
        }
    })
}

/// Reads a single file up to `budget` bytes, returning its content and bytes consumed.
/// Only reads the needed amount from disk to avoid OOM on large files.
async fn read_single_file(resolved: &Path, budget: u64) -> Option<(String, u64)> {
    use tokio::io::AsyncReadExt;

    let metadata = tokio::fs::metadata(resolved).await.ok()?;
    if !metadata.is_file() {
        return None;
    }

    let to_read = metadata.len().min(budget) as usize;
    if to_read == 0 {
        return None;
    }

    let mut file = tokio::fs::File::open(resolved).await.ok()?;
    let mut buf = vec![0u8; to_read];
    let n = file.read(&mut buf).await.ok()?;
    buf.truncate(n);

    let content = String::from_utf8_lossy(&buf).into_owned();
    Some((content, n as u64))
}

/// Reads output files asynchronously, respecting a total byte budget.
/// Returns the combined content of all readable files, or None if nothing was read.
pub async fn read_output_files(paths: &[String]) -> Option<String> {
    let mut combined = String::new();
    let mut budget_remaining = MAX_FILE_READ_SIZE;

    for raw_path in paths {
        let resolved = match resolve_path(raw_path) {
            Some(p) => p,
            None => {
                debug!("Could not resolve path: {}", raw_path);
                continue;
            }
        };

        if !is_safe_path(&resolved) {
            debug!("Path traversal blocked: {}", resolved.display());
            continue;
        }

        match read_single_file(&resolved, budget_remaining).await {
            Some((content, bytes_read)) => {
                if !combined.is_empty() {
                    combined.push_str("\n--- ");
                    combined.push_str(&resolved.display().to_string());
                    combined.push_str(" ---\n");
                }
                combined.push_str(&content);
                budget_remaining = budget_remaining.saturating_sub(bytes_read);
                if budget_remaining == 0 {
                    combined.push_str("\n... [file output truncated]");
                    break;
                }
            }
            None => {
                debug!("Cannot read {}", resolved.display());
            }
        }
    }

    if combined.is_empty() {
        None
    } else {
        Some(combined)
    }
}

/// Enriches command output by reading output files when stdout is too small.
///
/// If the original output is >= OUTPUT_THRESHOLD bytes (trimmed), it is returned as-is.
/// Otherwise, extracts file paths from the command, reads the files, and returns their content.
pub async fn enrich_output(command: &str, executor_type: &str, original_output: &str) -> String {
    if original_output.trim().len() >= OUTPUT_THRESHOLD {
        return original_output.to_string();
    }

    let paths = extract_output_paths(command, executor_type);
    if paths.is_empty() {
        return original_output.to_string();
    }

    debug!(
        "Output below threshold ({}B), checking {} output file(s)",
        original_output.trim().len(),
        paths.len()
    );

    match read_output_files(&paths).await {
        Some(file_content) => file_content,
        None => original_output.to_string(),
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    // --- extract_output_paths tests ---

    #[test]
    fn test_extract_unix_simple_redirect() {
        let paths = extract_output_paths("ps aux > /tmp/loot.txt", "bash");
        assert_eq!(paths, vec!["/tmp/loot.txt"]);
    }

    #[test]
    fn test_extract_unix_append_redirect() {
        let paths = extract_output_paths("ps >> /tmp/loot.txt", "sh");
        assert_eq!(paths, vec!["/tmp/loot.txt"]);
    }

    #[test]
    fn test_extract_unix_multiple_redirects() {
        let paths = extract_output_paths(
            "ls > /tmp/files.txt; cat /etc/passwd >> /tmp/loot.txt",
            "bash",
        );
        assert_eq!(paths, vec!["/tmp/files.txt", "/tmp/loot.txt"]);
    }

    #[test]
    fn test_extract_unix_dedup() {
        let paths =
            extract_output_paths("echo a >> /tmp/loot.txt; echo b >> /tmp/loot.txt", "bash");
        assert_eq!(paths, vec!["/tmp/loot.txt"]);
    }

    #[test]
    fn test_extract_unix_no_redirect() {
        let paths = extract_output_paths("whoami", "bash");
        assert!(paths.is_empty());
    }

    #[test]
    fn test_extract_unix_path_outside_tmp() {
        let paths = extract_output_paths("echo test > /var/log/output.txt", "bash");
        assert!(paths.is_empty());
    }

    #[test]
    fn test_extract_unix_zsh() {
        let paths = extract_output_paths("ls > /tmp/out.txt", "zsh");
        assert_eq!(paths, vec!["/tmp/out.txt"]);
    }

    #[test]
    fn test_extract_cmd_temp_redirect() {
        let paths = extract_output_paths("dir > %temp%\\output.txt", "cmd");
        assert_eq!(paths, vec!["%temp%\\output.txt"]);
    }

    #[test]
    fn test_extract_cmd_temp_append() {
        let paths = extract_output_paths("systeminfo >> %TEMP%\\sysinfo.txt", "cmd");
        assert_eq!(paths, vec!["%TEMP%\\sysinfo.txt"]);
    }

    #[test]
    fn test_extract_cmd_userprofile() {
        let paths = extract_output_paths("dir >> %userprofile%\\data.txt", "cmd");
        assert_eq!(paths, vec!["%userprofile%\\data.txt"]);
    }

    #[test]
    fn test_extract_powershell_redirect() {
        let paths = extract_output_paths("Get-Process >> $env:TEMP\\procs.txt", "powershell");
        assert_eq!(paths, vec!["$env:TEMP\\procs.txt"]);
    }

    #[test]
    fn test_extract_powershell_outfile() {
        let paths = extract_output_paths(
            "Get-Service | Out-File $env:TEMP\\services.txt",
            "powershell",
        );
        assert_eq!(paths, vec!["$env:TEMP\\services.txt"]);
    }

    #[test]
    fn test_extract_powershell_outfile_filepath() {
        let paths = extract_output_paths(
            "Get-Service | Out-File -FilePath $env:TEMP\\services.txt",
            "powershell",
        );
        assert_eq!(paths, vec!["$env:TEMP\\services.txt"]);
    }

    #[test]
    fn test_extract_powershell_set_content() {
        let paths = extract_output_paths(
            "Set-Content -Path $env:TEMP\\data.txt -Value 'test'",
            "powershell",
        );
        assert_eq!(paths, vec!["$env:TEMP\\data.txt"]);
    }

    #[test]
    fn test_extract_powershell_unix_tmp() {
        let paths = extract_output_paths("Get-Process > /tmp/procs.txt", "pwsh");
        assert_eq!(paths, vec!["/tmp/procs.txt"]);
    }

    #[test]
    fn test_extract_unknown_executor() {
        let paths = extract_output_paths("echo test > /tmp/file.txt", "unknown");
        assert!(paths.is_empty());
    }

    #[test]
    fn test_extract_empty_command() {
        let paths = extract_output_paths("", "bash");
        assert!(paths.is_empty());
    }

    // --- resolve_path tests ---

    #[test]
    fn test_resolve_unix_absolute() {
        let p = resolve_path("/tmp/loot.txt");
        assert_eq!(p, Some(PathBuf::from("/tmp/loot.txt")));
    }

    #[test]
    fn test_resolve_windows_absolute() {
        let p = resolve_path("C:\\Users\\test\\file.txt");
        assert_eq!(p, Some(PathBuf::from("C:\\Users\\test\\file.txt")));
    }

    #[test]
    fn test_resolve_powershell_env_temp() {
        let p = resolve_path("$env:TEMP\\test.txt");
        assert!(p.is_some());
    }

    #[test]
    fn test_resolve_cmd_temp() {
        let p = resolve_path("%temp%\\test.txt");
        assert!(p.is_some());
    }

    #[test]
    fn test_resolve_unknown_format() {
        let p = resolve_path("relative/path.txt");
        assert!(p.is_none());
    }

    #[test]
    fn test_resolve_empty() {
        let p = resolve_path("");
        assert!(p.is_none());
    }

    // --- is_safe_path tests ---

    #[test]
    fn test_safe_path_in_tmp() {
        let tmp = std::env::temp_dir();
        let test_file = tmp.join("autostrike_safe_test.txt");
        std::fs::write(&test_file, "test").unwrap();
        assert!(is_safe_path(&test_file));
        std::fs::remove_file(&test_file).unwrap();
    }

    #[test]
    fn test_unsafe_path_traversal() {
        // /tmp/../../etc/passwd resolves to /etc/passwd — outside /tmp/
        assert!(!is_safe_path(&PathBuf::from("/tmp/../../etc/passwd")));
    }

    #[test]
    fn test_unsafe_path_outside_tmp() {
        assert!(!is_safe_path(&PathBuf::from("/etc/hostname")));
    }

    #[test]
    fn test_safe_path_nonexistent() {
        assert!(!is_safe_path(&PathBuf::from("/tmp/nonexistent_99999.txt")));
    }

    #[tokio::test]
    async fn test_read_blocks_path_traversal() {
        // Even though resolve_path returns a path, read_output_files should block traversal
        let result = read_output_files(&["/tmp/../../etc/passwd".to_string()]).await;
        assert!(result.is_none());
    }

    // --- read_output_files tests ---

    #[tokio::test]
    async fn test_read_nonexistent_file() {
        let result =
            read_output_files(&["/tmp/nonexistent_autostrike_test_12345.txt".to_string()]).await;
        assert!(result.is_none());
    }

    #[tokio::test]
    async fn test_read_existing_file() {
        let tmp = std::env::temp_dir();
        let test_path = tmp.join("autostrike_output_test.txt");
        tokio::fs::write(&test_path, "test output content")
            .await
            .unwrap();

        let result = read_output_files(&[test_path.to_string_lossy().to_string()]).await;
        assert!(result.is_some());
        assert_eq!(result.unwrap(), "test output content");

        tokio::fs::remove_file(&test_path).await.unwrap();
    }

    #[tokio::test]
    async fn test_read_multiple_files() {
        let tmp = std::env::temp_dir();
        let path1 = tmp.join("autostrike_test1.txt");
        let path2 = tmp.join("autostrike_test2.txt");
        tokio::fs::write(&path1, "content1").await.unwrap();
        tokio::fs::write(&path2, "content2").await.unwrap();

        let result = read_output_files(&[
            path1.to_string_lossy().to_string(),
            path2.to_string_lossy().to_string(),
        ])
        .await;
        assert!(result.is_some());
        let content = result.unwrap();
        assert!(content.contains("content1"));
        assert!(content.contains("content2"));

        tokio::fs::remove_file(&path1).await.unwrap();
        tokio::fs::remove_file(&path2).await.unwrap();
    }

    #[tokio::test]
    async fn test_read_empty_paths() {
        let result = read_output_files(&[]).await;
        assert!(result.is_none());
    }

    // --- enrich_output tests ---

    #[tokio::test]
    async fn test_enrich_long_output_skipped() {
        let long_output = "x".repeat(100);
        let result = enrich_output("ps >> /tmp/loot.txt", "bash", &long_output).await;
        assert_eq!(result, long_output);
    }

    #[tokio::test]
    async fn test_enrich_no_paths_returns_original() {
        let result = enrich_output("whoami", "bash", "").await;
        assert_eq!(result, "");
    }

    #[tokio::test]
    async fn test_enrich_with_file() {
        let test_path = "/tmp/autostrike_enrich_test.txt";
        // Write the file to /tmp/ (this test only runs meaningfully on Unix
        // where /tmp/ exists; on Windows it will gracefully fall through)
        if tokio::fs::write(test_path, "enriched content here")
            .await
            .is_ok()
        {
            let result = enrich_output(&format!("some_cmd >> {}", test_path), "bash", "").await;
            assert_eq!(result, "enriched content here");
            let _ = tokio::fs::remove_file(test_path).await;
        }
    }

    #[tokio::test]
    async fn test_enrich_file_not_found_returns_original() {
        let result =
            enrich_output("cmd >> /tmp/nonexistent_file_98765.txt", "bash", "original").await;
        assert_eq!(result, "original");
    }

    #[tokio::test]
    async fn test_enrich_threshold_boundary() {
        // Exactly 49 bytes (trimmed) — should try enrichment
        let output_49 = "x".repeat(49);
        let paths = extract_output_paths(&format!("{} > /tmp/test.txt", output_49), "bash");
        // The output itself doesn't trigger skip
        assert!(output_49.trim().len() < OUTPUT_THRESHOLD);
        assert!(!paths.is_empty());
    }

    #[tokio::test]
    async fn test_enrich_threshold_exact() {
        // Exactly 50 bytes — should skip enrichment
        let output_50 = "x".repeat(50);
        let result = enrich_output("ps >> /tmp/loot.txt", "bash", &output_50).await;
        assert_eq!(result, output_50);
    }

    // --- Integration tests (Unix only — use sh executor and /tmp/) ---

    #[cfg(unix)]
    #[tokio::test]
    async fn test_integration_real_command_with_redirect() {
        let executor = crate::executor::CommandExecutor::new();
        let command = "echo 'hello from integration test' > /tmp/autostrike_integ_test.txt";
        let result = executor
            .execute("sh", command, std::time::Duration::from_secs(5))
            .await;
        assert!(result.success);
        assert!(result.output.trim().is_empty());

        let enriched = enrich_output(command, "sh", &result.output).await;
        assert!(
            enriched.contains("hello from integration test"),
            "Expected file content, got: {}",
            enriched
        );

        tokio::fs::remove_file("/tmp/autostrike_integ_test.txt")
            .await
            .unwrap();
    }

    #[cfg(unix)]
    #[tokio::test]
    async fn test_integration_append_redirect() {
        let executor = crate::executor::CommandExecutor::new();
        let command =
            "echo line1 >> /tmp/autostrike_append_test.txt; echo line2 >> /tmp/autostrike_append_test.txt";
        let result = executor
            .execute("sh", command, std::time::Duration::from_secs(5))
            .await;
        assert!(result.success);

        let enriched = enrich_output(command, "sh", &result.output).await;
        assert!(
            enriched.contains("line1"),
            "Expected line1, got: {}",
            enriched
        );
        assert!(
            enriched.contains("line2"),
            "Expected line2, got: {}",
            enriched
        );

        tokio::fs::remove_file("/tmp/autostrike_append_test.txt")
            .await
            .unwrap();
    }

    #[cfg(unix)]
    #[tokio::test]
    async fn test_integration_command_with_stdout_and_redirect() {
        let executor = crate::executor::CommandExecutor::new();
        let command = "echo 'This goes to stdout and is over fifty bytes for certain right now'; echo 'file only' > /tmp/autostrike_both_test.txt";
        let result = executor
            .execute("sh", command, std::time::Duration::from_secs(5))
            .await;
        assert!(result.success);

        let enriched = enrich_output(command, "sh", &result.output).await;
        assert_eq!(enriched, result.output);

        tokio::fs::remove_file("/tmp/autostrike_both_test.txt")
            .await
            .unwrap();
    }

    #[cfg(unix)]
    #[tokio::test]
    async fn test_integration_ps_redirect_like_t1057() {
        let executor = crate::executor::CommandExecutor::new();
        let command = "ps >> /tmp/autostrike_ps_test.txt";
        let result = executor
            .execute("sh", command, std::time::Duration::from_secs(5))
            .await;
        assert!(result.success);
        assert!(result.output.trim().is_empty());

        let enriched = enrich_output(command, "sh", &result.output).await;
        assert!(
            enriched.contains("PID"),
            "Expected process list, got: {}",
            &enriched[..enriched.len().min(200)]
        );

        tokio::fs::remove_file("/tmp/autostrike_ps_test.txt")
            .await
            .unwrap();
    }

    #[cfg(unix)]
    #[tokio::test]
    async fn test_integration_uname_redirect_like_t1082() {
        let executor = crate::executor::CommandExecutor::new();
        let command = "uname -a >> /tmp/autostrike_uname_test.txt";
        let result = executor
            .execute("sh", command, std::time::Duration::from_secs(5))
            .await;
        assert!(result.success);
        assert!(result.output.trim().is_empty());

        let enriched = enrich_output(command, "sh", &result.output).await;
        assert!(
            enriched.contains("Linux"),
            "Expected Linux info, got: {}",
            enriched
        );

        tokio::fs::remove_file("/tmp/autostrike_uname_test.txt")
            .await
            .unwrap();
    }
}
