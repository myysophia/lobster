package products

import (
	"fmt"
	"path/filepath"

	"lobster/internal/platform"
)

type QoderWork struct{}

func NewQoderWork() QoderWork {
	return QoderWork{}
}

func (QoderWork) Key() string {
	return "qoderwork"
}

func (QoderWork) DisplayName() string {
	return "QoderWork"
}

func (QoderWork) Summary() string {
	return "QoderWork 安装入口，支持按平台拉起官方下载。"
}

func (QoderWork) InstallPlan(info platform.Info) (InstallPlan, error) {
	switch info.OS {
	case platform.Windows:
		if info.Arch != "amd64" {
			return InstallPlan{}, fmt.Errorf("当前平台暂不支持 QoderWork 安装: %s。官网当前页面仅明确 Windows 10+。", info.String())
		}
		return urlInstallPlan(info, "打开 QoderWork 官方安装包下载链接", "https://download.qoder.com.cn/qoder-work/releases/latest/QoderWork-Setup-User-x64.exe")
	case platform.Darwin:
		switch info.Arch {
		case "arm64":
			return urlInstallPlan(info, "打开 QoderWork 官方安装包下载链接", "https://download.qoder.com.cn/qoder-work/releases/latest/QoderWork-arm64.dmg")
		case "amd64":
			return urlInstallPlan(info, "打开 QoderWork 官方安装包下载链接", "https://download.qoder.com.cn/qoder-work/releases/latest/QoderWork-x64.dmg")
		default:
			return InstallPlan{}, fmt.Errorf("当前平台暂不支持 QoderWork 安装: %s。官网当前页面仅明确 macOS 14+。", info.String())
		}
	default:
		return InstallPlan{}, fmt.Errorf("当前平台暂不支持 QoderWork 安装: %s。官网当前页面仅明确 macOS 14+ 与 Windows 10+。", info.String())
	}
}

func (QoderWork) DetectPlan(info platform.Info) DetectPlan {
	paths := []string{}

	switch info.OS {
	case platform.Darwin:
		paths = append(paths,
			"/Applications/QoderWork.app",
			filepath.Join(userHomeDir(), "Applications", "QoderWork.app"),
		)
	case platform.Windows:
		paths = append(paths,
			`C:\Program Files\QoderWork`,
			`C:\Program Files\QoderWork\QoderWork.exe`,
		)
		paths = append(paths, windowsProgramInstallPaths("QoderWork")...)
	}

	return DetectPlan{
		Commands: []string{"qoderwork"},
		Paths:    paths,
		Notes: []string{
			"QoderWork 官网当前页面可确认 macOS 14+ 与 Windows 10+ 下载入口。",
			"当前版本默认拉起官网页面，再由用户完成官方下载与安装。",
		},
	}
}

func (QoderWork) LaunchPlan(info platform.Info) LaunchPlan {
	switch info.OS {
	case platform.Darwin:
		return LaunchPlan{
			ExecCandidates: [][]string{
				{"open", "-a", "QoderWork"},
				{"open", "/Applications/QoderWork.app"},
			},
			Notes: []string{"macOS 下优先尝试通过 open 打开 QoderWork。"},
		}
	case platform.Windows:
		return LaunchPlan{
			ExecCandidates: [][]string{
				{"cmd", "/c", "start", "", "QoderWork"},
				{"powershell", "-NoProfile", "-Command", "Start-Process QoderWork"},
			},
			Notes: []string{"Windows 下优先尝试通过 Start-Process 启动 QoderWork。"},
		}
	default:
		return LaunchPlan{
			Notes: []string{"当前平台未定义 QoderWork 的自动启动方式。"},
		}
	}
}
