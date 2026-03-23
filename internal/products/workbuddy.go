package products

import (
	"fmt"
	"os"
	"path/filepath"

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
		return InstallPlan{
			Summary: "调用腾讯官方 PowerShell 安装器",
			Exec: []string{
				"powershell",
				"-NoProfile",
				"-ExecutionPolicy",
				"Bypass",
				"-Command",
				"irm https://copilot.tencent.com/cli/install.ps1 | iex",
			},
		}, nil
	default:
		return InstallPlan{}, fmt.Errorf("当前平台暂不支持 WorkBuddy 安装: %s", info.String())
	}
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
