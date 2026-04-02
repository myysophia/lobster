package products

import (
	"strings"
	"testing"

	"lobster/internal/platform"
)

func TestQoderWorkInstallPlanUsesOfficialDownloadPage(t *testing.T) {
	subject := NewQoderWork()

	plan, err := subject.InstallPlan(platform.Info{OS: platform.Windows, Arch: "amd64"})
	if err != nil {
		t.Fatalf("Windows 安装计划不应报错，实际：%v", err)
	}
	if !plan.SkipVerify {
		t.Fatalf("QoderWork 下载交接模式应跳过安装后即时校验")
	}
	if got := strings.Join(plan.Exec, " "); !strings.Contains(got, "QoderWork-Setup-User-x64.exe") {
		t.Fatalf("Windows 应使用已验证的用户级安装器直链，实际：%s", got)
	}

	plan, err = subject.InstallPlan(platform.Info{OS: platform.Darwin, Arch: "arm64"})
	if err != nil {
		t.Fatalf("macOS 14+ 当前应视为支持平台，实际：%v", err)
	}
	if got := strings.Join(plan.Exec, " "); !strings.Contains(got, "QoderWork-arm64.dmg") {
		t.Fatalf("Mac Apple Silicon 应使用 arm64 直链，实际：%s", got)
	}

	plan, err = subject.InstallPlan(platform.Info{OS: platform.Darwin, Arch: "amd64"})
	if err != nil {
		t.Fatalf("Mac Intel 当前应视为支持平台，实际：%v", err)
	}
	if got := strings.Join(plan.Exec, " "); !strings.Contains(got, "QoderWork-x64.dmg") {
		t.Fatalf("Mac Intel 应使用 x64 直链，实际：%s", got)
	}
}

func TestQoderWorkInstallPlanRejectsUnsupportedPlatform(t *testing.T) {
	subject := NewQoderWork()

	_, err := subject.InstallPlan(platform.Info{OS: platform.Linux, Arch: "amd64"})
	if err == nil {
		t.Fatalf("Linux 当前应返回不支持")
	}
	if !strings.Contains(err.Error(), "macOS 14+") || !strings.Contains(err.Error(), "Windows 10+") {
		t.Fatalf("错误信息应明确官网当前支持范围，实际：%v", err)
	}
}
