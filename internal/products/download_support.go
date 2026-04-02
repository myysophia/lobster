package products

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"lobster/internal/platform"
)

func urlInstallPlan(info platform.Info, summary string, url string) (InstallPlan, error) {
	switch info.OS {
	case platform.Darwin:
		return InstallPlan{
			Summary:    summary,
			Exec:       []string{"open", url},
			SkipVerify: true,
		}, nil
	case platform.Windows:
		return InstallPlan{
			Summary:    summary,
			Exec:       []string{"cmd", "/c", "start", "", url},
			SkipVerify: true,
		}, nil
	case platform.Linux:
		return InstallPlan{
			Summary:    summary,
			Exec:       []string{"xdg-open", url},
			SkipVerify: true,
		}, nil
	default:
		return InstallPlan{}, fmt.Errorf("当前平台暂不支持下载入口拉起: %s", info.String())
	}
}

func windowsProgramInstallPaths(appName string) []string {
	localAppData := strings.TrimSpace(osLocalAppData())
	if localAppData == "" {
		localAppData = filepath.Join(userHomeDir(), "AppData", "Local")
	}
	if strings.TrimSpace(localAppData) == "" {
		return nil
	}

	return []string{
		filepath.Join(localAppData, "Programs", appName),
		filepath.Join(localAppData, "Programs", appName, appName+".exe"),
	}
}

var osLocalAppData = func() string {
	return strings.TrimSpace(os.Getenv("LOCALAPPDATA"))
}
