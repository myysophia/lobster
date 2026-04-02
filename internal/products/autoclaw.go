package products

import (
	"fmt"
	"path/filepath"

	"lobster/internal/platform"
)

type AutoClaw struct{}

func NewAutoClaw() AutoClaw {
	return AutoClaw{}
}

func (AutoClaw) Key() string {
	return "autoclaw"
}

func (AutoClaw) DisplayName() string {
	return "AutoClaw"
}

func (AutoClaw) Summary() string {
	return "智谱 AutoClaw 安装入口，支持按平台拉起官方下载。"
}

func (AutoClaw) InstallPlan(info platform.Info) (InstallPlan, error) {
	const summary = "打开 AutoClaw 官方安装包下载链接"

	switch info.OS {
	case platform.Windows:
		if info.Arch != "amd64" {
			return InstallPlan{}, fmt.Errorf("当前平台暂不支持 AutoClaw 安装: %s。官网当前确认支持 Windows 10 / 11（64 位）。", info.String())
		}
		return urlInstallPlan(info, summary, "https://autoglm.aminer.cn/autoclaw/updates/autoclaw-0.2.25-setup.exe")
	case platform.Darwin:
		switch info.Arch {
		case "arm64":
			return urlInstallPlan(info, summary, "https://autoglm.aminer.cn/autoclaw/updates/autoclaw-0.2.25.dmg")
		case "amd64":
			return urlInstallPlan(info, summary, "https://autoglm.aminer.cn/autoclaw/updates/x64/autoclaw-0.2.25-x64.dmg")
		default:
			return InstallPlan{}, fmt.Errorf("当前平台暂不支持 AutoClaw 安装: %s。官网当前确认支持 macOS Apple Silicon 与 macOS Intel。", info.String())
		}
	default:
		return InstallPlan{}, fmt.Errorf("当前平台暂不支持 AutoClaw 安装: %s。官网当前确认支持 Windows、macOS Apple Silicon、macOS Intel。", info.String())
	}
}

func (AutoClaw) DetectPlan(info platform.Info) DetectPlan {
	paths := []string{}

	switch info.OS {
	case platform.Darwin:
		paths = append(paths,
			"/Applications/AutoClaw.app",
			filepath.Join(userHomeDir(), "Applications", "AutoClaw.app"),
		)
	case platform.Windows:
		paths = append(paths,
			`C:\Program Files\AutoClaw`,
			`C:\Program Files\AutoClaw\AutoClaw.exe`,
		)
		paths = append(paths, windowsProgramInstallPaths("AutoClaw")...)
	}

	return DetectPlan{
		Commands: []string{"autoclaw"},
		Paths:    paths,
		Notes: []string{
			"AutoClaw 官网当前可确认 Windows、macOS Apple Silicon、macOS Intel 下载入口。",
			"如果刚完成安装但命令仍不可见，可先重开终端再执行 status。",
		},
	}
}

func (AutoClaw) LaunchPlan(info platform.Info) LaunchPlan {
	switch info.OS {
	case platform.Darwin:
		return LaunchPlan{
			ExecCandidates: [][]string{
				{"open", "-a", "AutoClaw"},
				{"open", "/Applications/AutoClaw.app"},
			},
			Notes: []string{"macOS 下优先尝试通过 open 打开 AutoClaw。"},
		}
	case platform.Windows:
		return LaunchPlan{
			ExecCandidates: [][]string{
				{"cmd", "/c", "start", "", "AutoClaw"},
				{"powershell", "-NoProfile", "-Command", "Start-Process AutoClaw"},
			},
			Notes: []string{"Windows 下优先尝试通过 Start-Process 启动 AutoClaw。"},
		}
	default:
		return LaunchPlan{
			Notes: []string{"当前平台未定义 AutoClaw 的自动启动方式。"},
		}
	}
}
