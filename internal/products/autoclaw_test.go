package products

import (
	"strings"
	"testing"

	"lobster/internal/platform"
)

func TestAutoClawInstallPlanUsesVerifiedDirectDownloads(t *testing.T) {
	subject := NewAutoClaw()

	plan, err := subject.InstallPlan(platform.Info{OS: platform.Windows, Arch: "amd64"})
	if err != nil {
		t.Fatalf("Windows 安装计划不应报错，实际：%v", err)
	}
	if !plan.SkipVerify {
		t.Fatalf("AutoClaw 下载交接模式应跳过安装后即时校验")
	}
	if got := strings.Join(plan.Exec, " "); !strings.Contains(got, "autoclaw-0.2.25-setup.exe") {
		t.Fatalf("Windows 应使用已验证的直链下载地址，实际：%s", got)
	}

	plan, err = subject.InstallPlan(platform.Info{OS: platform.Darwin, Arch: "arm64"})
	if err != nil {
		t.Fatalf("Mac Apple Silicon 安装计划不应报错，实际：%v", err)
	}
	if got := strings.Join(plan.Exec, " "); !strings.Contains(got, "autoclaw-0.2.25.dmg") {
		t.Fatalf("Mac Apple Silicon 应使用 arm 版下载地址，实际：%s", got)
	}
	if strings.Contains(strings.Join(plan.Exec, " "), "-x64.dmg") {
		t.Fatalf("Apple Silicon 不应误用 Intel 安装包，实际：%v", plan.Exec)
	}

	plan, err = subject.InstallPlan(platform.Info{OS: platform.Darwin, Arch: "amd64"})
	if err != nil {
		t.Fatalf("Mac Intel 安装计划不应报错，实际：%v", err)
	}
	if got := strings.Join(plan.Exec, " "); !strings.Contains(got, "autoclaw-0.2.25-x64.dmg") {
		t.Fatalf("Mac Intel 应使用 x64 版下载地址，实际：%s", got)
	}
}

func TestAutoClawInstallPlanRejectsUnsupportedPlatform(t *testing.T) {
	subject := NewAutoClaw()

	_, err := subject.InstallPlan(platform.Info{OS: platform.Linux, Arch: "amd64"})
	if err == nil {
		t.Fatalf("Linux 当前应返回不支持")
	}
	if !strings.Contains(err.Error(), "Windows、macOS Apple Silicon、macOS Intel") {
		t.Fatalf("错误信息应明确官网支持范围，实际：%v", err)
	}
}
