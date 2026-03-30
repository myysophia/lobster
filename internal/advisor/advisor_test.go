package advisor

import (
	"strings"
	"testing"

	"lobster/internal/detector"
	"lobster/internal/installer"
	"lobster/internal/platform"
	"lobster/internal/products"
)

func TestInstallSummaryForInstalledResult(t *testing.T) {
	product := products.NewWorkBuddy()
	result := installer.Result{
		Outcome: installer.OutcomeInstalled,
		PostStatus: detector.Status{
			Installed:        true,
			CommandAvailable: true,
			CommandPath:      "/tmp/codebuddy",
		},
	}

	lines := InstallSummary(product, result)
	joined := strings.Join(lines, "\n")

	if !strings.Contains(joined, "安装完成。") {
		t.Fatalf("期望包含安装成功提示，实际输出：%s", joined)
	}
	if !strings.Contains(joined, "可执行命令路径：/tmp/codebuddy") {
		t.Fatalf("期望包含命令路径，实际输出：%s", joined)
	}
	if !strings.Contains(joined, "下一步：执行 `lobster workbuddy open`") {
		t.Fatalf("期望包含下一步建议，实际输出：%s", joined)
	}
}

func TestInstallSummaryForVerifyFailedResult(t *testing.T) {
	product := products.NewWorkBuddy()
	result := installer.Result{
		Outcome:    installer.OutcomeVerifyFailed,
		PostStatus: detector.Status{Installed: false},
	}

	lines := InstallSummary(product, result)
	joined := strings.Join(lines, "\n")

	if !strings.Contains(joined, "安装后仍未检测到明确结果") {
		t.Fatalf("期望包含校验失败提示，实际输出：%s", joined)
	}
	if !strings.Contains(joined, "lobster workbuddy doctor") {
		t.Fatalf("期望包含诊断建议，实际输出：%s", joined)
	}
}

func TestInstallSummaryForInstalledButPathNotReady(t *testing.T) {
	product := products.NewWorkBuddy()
	result := installer.Result{
		Outcome: installer.OutcomeInstalled,
		PostStatus: detector.Status{
			Installed:       true,
			HasPathEvidence: true,
			FoundPaths:      []string{`C:\Users\tester\AppData\Local\codebuddy\bin\codebuddy.exe`},
			Warnings:        []string{"检测到安装痕迹，但当前终端还没有识别到可执行命令，可能需要重新打开终端。"},
		},
	}

	lines := InstallSummary(product, result)
	joined := strings.Join(lines, "\n")

	if !strings.Contains(joined, "下一步：重开终端后执行 `lobster workbuddy status`") {
		t.Fatalf("期望包含 PATH 生效建议，实际输出：%s", joined)
	}
	if strings.Contains(joined, "lobster workbuddy open") {
		t.Fatalf("命令尚不可用时不应直接建议 open，实际输出：%s", joined)
	}
}

func TestDoctorSummaryForCommandAvailable(t *testing.T) {
	product := products.NewWorkBuddy()
	info := platform.Info{OS: platform.Darwin, Arch: "arm64", HasDesktop: true}
	status := detector.Status{
		Installed:        true,
		CommandAvailable: true,
		CommandPath:      "/tmp/codebuddy",
		FoundPaths:       []string{"/Applications/WorkBuddy.app"},
		Notes:            []string{"优先检测 codebuddy / workbuddy 命令是否已进入 PATH。"},
	}

	lines := DoctorSummary(product, info, status)
	joined := strings.Join(lines, "\n")

	if !strings.Contains(joined, "已检测到可用命令") {
		t.Fatalf("期望包含可用命令结论，实际输出：%s", joined)
	}
	if !strings.Contains(joined, "命中路径：/Applications/WorkBuddy.app") {
		t.Fatalf("期望包含命中路径，实际输出：%s", joined)
	}
	if !strings.Contains(joined, "建议：可以直接执行 `lobster workbuddy open`") {
		t.Fatalf("期望包含下一步建议，实际输出：%s", joined)
	}
}

func TestDoctorSummaryForNoDesktop(t *testing.T) {
	product := products.NewWorkBuddy()
	info := platform.Info{OS: platform.Linux, Arch: "amd64", HasDesktop: false}
	status := detector.Status{
		Installed:       true,
		HasPathEvidence: true,
		FoundPaths:      []string{"/tmp/workbuddy.desktop"},
	}

	lines := DoctorSummary(product, info, status)
	joined := strings.Join(lines, "\n")

	if !strings.Contains(joined, "当前会话未检测到桌面环境") {
		t.Fatalf("期望包含桌面环境提示，实际输出：%s", joined)
	}
	if !strings.Contains(joined, "当前命令还不可用") {
		t.Fatalf("期望包含命令不可用结论，实际输出：%s", joined)
	}
	if !strings.Contains(joined, "优先重新打开终端或刷新 PATH") {
		t.Fatalf("期望包含 PATH 建议，实际输出：%s", joined)
	}
}
