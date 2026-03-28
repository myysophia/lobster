package products

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"lobster/internal/platform"
)

type WorkBuddy struct{}

func NewWorkBuddy() WorkBuddy {
	return WorkBuddy{}
}

func (WorkBuddy) Key() string {
	return "workbuddy"
}

func (WorkBuddy) DisplayName() string {
	return "WorkBuddy"
}

func (WorkBuddy) Summary() string {
	return "腾讯 WorkBuddy 安装与启动入口"
}

func (WorkBuddy) InstallPlan(info platform.Info) (InstallPlan, error) {
	switch info.OS {
	case platform.Darwin, platform.Linux:
		return InstallPlan{
			Summary: "调用腾讯官方 shell 安装器",
			Exec: []string{
				"/bin/bash",
				"-lc",
				"curl -fsSL https://copilot.tencent.com/cli/install.sh | bash",
			},
		}, nil
	case platform.Windows:
		env := map[string]string{}
		if path := detectGitBashPath(); path != "" {
			env["CODEBUDDY_CODE_GIT_BASH_PATH"] = path
		}

		return InstallPlan{
			Summary: "调用腾讯官方 PowerShell 安装器",
			Exec: []string{
				"powershell",
				"-NoProfile",
				"-NonInteractive",
				"-ExecutionPolicy",
				"Bypass",
				"-Command",
				"$ProgressPreference='SilentlyContinue'; [Console]::OutputEncoding=[System.Text.UTF8Encoding]::new($false); $OutputEncoding=[System.Text.UTF8Encoding]::new($false); irm https://copilot.tencent.com/cli/install.ps1 | iex",
			},
			Env: env,
		}, nil
	default:
		return InstallPlan{}, fmt.Errorf("当前平台暂不支持 WorkBuddy 安装: %s", info.String())
	}
}

func (WorkBuddy) ValidateInstall(info platform.Info) error {
	if info.OS != platform.Windows {
		return nil
	}

	if path := detectGitBashPath(); path != "" {
		return nil
	}

	return fmt.Errorf("Windows 安装 WorkBuddy 之前需要先安装 Git Bash（https://git-scm.com/downloads/win），安装完成后请重新打开终端再执行；如果 Git Bash 已安装但不在 PATH，请设置环境变量 CODEBUDDY_CODE_GIT_BASH_PATH 指向 bash.exe")
}

func (WorkBuddy) DetectPlan(info platform.Info) DetectPlan {
	paths := []string{}

	switch info.OS {
	case platform.Darwin:
		paths = append(paths,
			"/Applications/WorkBuddy.app",
			filepath.Join(userHomeDir(), "Applications", "WorkBuddy.app"),
			filepath.Join(userHomeDir(), ".codebuddy"),
			filepath.Join(userHomeDir(), ".workbuddy"),
		)
	case platform.Windows:
		paths = append(paths,
			`C:\Program Files\WorkBuddy`,
			`C:\Program Files\Tencent\WorkBuddy`,
			`C:\Users\Public\Desktop\WorkBuddy.lnk`,
		)
		paths = append(paths, windowsUserInstallPaths()...)
	case platform.Linux:
		paths = append(paths,
			filepath.Join(userHomeDir(), ".local", "bin", "codebuddy"),
			filepath.Join(userHomeDir(), ".local", "share", "applications", "workbuddy.desktop"),
			filepath.Join(userHomeDir(), ".codebuddy"),
			filepath.Join(userHomeDir(), ".workbuddy"),
		)
	}

	return DetectPlan{
		Commands: []string{"codebuddy", "workbuddy"},
		Paths:    paths,
		Notes: []string{
			"优先检测 codebuddy / workbuddy 命令是否已进入 PATH。",
			"如果命令还不可见，也会尝试检查常见安装路径。",
		},
	}
}

func (WorkBuddy) LaunchPlan(info platform.Info) LaunchPlan {
	switch info.OS {
	case platform.Darwin:
		return LaunchPlan{
			ExecCandidates: [][]string{
				{"open", "-a", "WorkBuddy"},
				{"open", "/Applications/WorkBuddy.app"},
			},
			Notes: []string{"macOS 下优先尝试通过 open 打开 WorkBuddy。"},
		}
	case platform.Windows:
		return LaunchPlan{
			ExecCandidates: [][]string{
				{"cmd", "/c", "start", "", "WorkBuddy"},
				{"powershell", "-NoProfile", "-Command", "Start-Process WorkBuddy"},
			},
			Notes: []string{"Windows 下优先尝试通过 Start-Process 启动 WorkBuddy。"},
		}
	case platform.Linux:
		return LaunchPlan{
			ExecCandidates: [][]string{
				{"xdg-open", "workbuddy"},
				{"workbuddy"},
			},
			Notes: []string{"Linux 下优先尝试 xdg-open，其次尝试直接调用 workbuddy。"},
		}
	default:
		return LaunchPlan{
			Notes: []string{"当前平台未定义启动方案。"},
		}
	}
}

func userHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return home
}

var (
	workBuddyLookPath = exec.LookPath
	workBuddyStat     = os.Stat
)

func detectGitBashPath() string {
	if override := strings.TrimSpace(os.Getenv("CODEBUDDY_CODE_GIT_BASH_PATH")); override != "" {
		if fileExists(override) {
			return override
		}
	}

	if path, err := workBuddyLookPath("bash.exe"); err == nil && fileExists(path) {
		return path
	}

	for _, candidate := range gitBashCandidates() {
		if fileExists(candidate) {
			return candidate
		}
	}

	return ""
}

func gitBashCandidates() []string {
	programFiles := strings.TrimSpace(os.Getenv("ProgramFiles"))
	programFilesX86 := strings.TrimSpace(os.Getenv("ProgramFiles(x86)"))

	candidates := []string{
		`C:\Program Files\Git\bin\bash.exe`,
		`C:\Program Files\Git\usr\bin\bash.exe`,
		`C:\Program Files (x86)\Git\bin\bash.exe`,
		`C:\Program Files (x86)\Git\usr\bin\bash.exe`,
	}

	if programFiles != "" {
		candidates = append(candidates,
			filepath.Join(programFiles, "Git", "bin", "bash.exe"),
			filepath.Join(programFiles, "Git", "usr", "bin", "bash.exe"),
		)
	}

	if programFilesX86 != "" {
		candidates = append(candidates,
			filepath.Join(programFilesX86, "Git", "bin", "bash.exe"),
			filepath.Join(programFilesX86, "Git", "usr", "bin", "bash.exe"),
		)
	}

	return candidates
}

func fileExists(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}
	_, err := workBuddyStat(path)
	return err == nil
}

func windowsUserInstallPaths() []string {
	localAppData := strings.TrimSpace(os.Getenv("LOCALAPPDATA"))
	if localAppData == "" {
		localAppData = filepath.Join(userHomeDir(), "AppData", "Local")
	}

	if strings.TrimSpace(localAppData) == "" {
		return nil
	}

	return []string{
		filepath.Join(localAppData, "codebuddy"),
		filepath.Join(localAppData, "codebuddy", "bin"),
		filepath.Join(localAppData, "codebuddy", "bin", "codebuddy.exe"),
		filepath.Join(localAppData, "codebuddy", "bin", "workbuddy.exe"),
	}
}
