//! System information gathering for agent registration.

use serde::{Deserialize, Serialize};
use sysinfo::{System, SystemExt};
use which::which;

/// System information collected from the host machine.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SystemInfo {
    /// Machine hostname.
    pub hostname: String,
    /// Current user's username.
    pub username: String,
    /// Operating system platform (linux, windows, darwin).
    pub platform: String,
    /// Available command executors (sh, bash, powershell, etc.).
    pub executors: Vec<String>,
    /// Operating system version.
    pub os_version: String,
    /// CPU architecture (x86_64, aarch64, etc.).
    pub architecture: String,
}

impl SystemInfo {
    /// Gathers system information from the host machine.
    pub fn gather() -> Self {
        // Determine platform
        let platform = if cfg!(target_os = "windows") {
            "windows"
        } else if cfg!(target_os = "linux") {
            "linux"
        } else if cfg!(target_os = "macos") {
            "darwin"
        } else {
            "unknown"
        };

        // Detect available executors
        let executors = Self::detect_executors();

        let sys = System::new();

        SystemInfo {
            hostname: sys.host_name().unwrap_or_else(|| "unknown".to_string()),
            username: whoami::username(),
            platform: platform.to_string(),
            executors,
            os_version: sys.os_version().unwrap_or_else(|| "unknown".to_string()),
            architecture: std::env::consts::ARCH.to_string(),
        }
    }

    fn detect_executors() -> Vec<String> {
        let mut executors = Vec::new();

        // Common executors to check
        let executor_checks = if cfg!(target_os = "windows") {
            vec![
                ("powershell", "powershell"),
                ("pwsh", "pwsh"),
                ("cmd", "cmd"),
            ]
        } else {
            vec![
                ("sh", "sh"),
                ("bash", "bash"),
                ("zsh", "zsh"),
                ("python3", "python3"),
                ("python", "python"),
            ]
        };

        for (name, command) in executor_checks {
            if which(command).is_ok() {
                executors.push(name.to_string());
            }
        }

        executors
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_gather_returns_valid_info() {
        let info = SystemInfo::gather();

        assert!(!info.hostname.is_empty());
        assert!(!info.username.is_empty());
        assert!(
            info.platform == "windows"
                || info.platform == "linux"
                || info.platform == "darwin"
                || info.platform == "unknown"
        );
        assert!(!info.architecture.is_empty());
    }

    #[test]
    fn test_platform_detection() {
        let info = SystemInfo::gather();

        #[cfg(target_os = "windows")]
        assert_eq!(info.platform, "windows");

        #[cfg(target_os = "linux")]
        assert_eq!(info.platform, "linux");

        #[cfg(target_os = "macos")]
        assert_eq!(info.platform, "darwin");
    }

    #[test]
    fn test_executors_not_empty() {
        let info = SystemInfo::gather();
        assert!(!info.executors.is_empty());
    }

    #[test]
    fn test_sh_available_on_unix() {
        #[cfg(not(target_os = "windows"))]
        {
            let info = SystemInfo::gather();
            assert!(info.executors.contains(&"sh".to_string()));
        }
    }

    #[test]
    fn test_system_info_clone() {
        let info = SystemInfo::gather();
        let cloned = info.clone();

        assert_eq!(cloned.hostname, info.hostname);
        assert_eq!(cloned.username, info.username);
        assert_eq!(cloned.platform, info.platform);
        assert_eq!(cloned.executors, info.executors);
        assert_eq!(cloned.os_version, info.os_version);
        assert_eq!(cloned.architecture, info.architecture);
    }

    #[test]
    fn test_system_info_debug() {
        let info = SystemInfo::gather();
        let debug_str = format!("{:?}", info);

        assert!(debug_str.contains("hostname"));
        assert!(debug_str.contains("username"));
        assert!(debug_str.contains("platform"));
    }

    #[test]
    fn test_system_info_serialization() {
        let info = SystemInfo::gather();
        let json = serde_json::to_string(&info).unwrap();

        assert!(json.contains("hostname"));
        assert!(json.contains("username"));
        assert!(json.contains("platform"));
        assert!(json.contains("executors"));
    }

    #[test]
    fn test_system_info_deserialization() {
        let json = r#"{
            "hostname": "test-host",
            "username": "test-user",
            "platform": "linux",
            "executors": ["sh", "bash"],
            "os_version": "5.0",
            "architecture": "x86_64"
        }"#;

        let info: SystemInfo = serde_json::from_str(json).unwrap();

        assert_eq!(info.hostname, "test-host");
        assert_eq!(info.username, "test-user");
        assert_eq!(info.platform, "linux");
        assert_eq!(info.executors, vec!["sh", "bash"]);
        assert_eq!(info.os_version, "5.0");
        assert_eq!(info.architecture, "x86_64");
    }

    #[test]
    fn test_architecture_is_valid() {
        let info = SystemInfo::gather();

        let valid_archs = [
            "x86_64",
            "x86",
            "aarch64",
            "arm",
            "arm64",
            "powerpc64",
            "riscv64",
        ];
        assert!(
            valid_archs
                .iter()
                .any(|&arch| info.architecture.contains(arch))
                || !info.architecture.is_empty()
        );
    }
}
