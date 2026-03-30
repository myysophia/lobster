package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"lobster/internal/detector"
	"lobster/internal/installer"
	"lobster/internal/platform"
)

func TestModelHandlesInstallFinishedMsg(t *testing.T) {
	subject := newModel("workbuddy")
	result := installer.Result{
		Outcome:    installer.OutcomeInstalled,
		PostStatus: detector.Status{Installed: true},
	}
	msg := installFinishedMsg{
		info:   platform.Info{OS: platform.Linux},
		result: result,
		output: "step 1\nstep 2",
	}

	updated, _ := subject.Update(msg)
	got, ok := updated.(model)
	if !ok {
		t.Fatalf("Update 未返回 model 类型")
	}

	if got.screen != screenWorkBuddyResult {
		t.Fatalf("安装完成后应切回结果页，实际：%s", got.screen)
	}
	if got.installResult.Outcome != result.Outcome {
		t.Fatalf("应记录最新安装结果 Outcome，实际：%s", got.installResult.Outcome)
	}
	if !got.statusLoaded || got.status.Installed != result.PostStatus.Installed {
		t.Fatalf("状态未同步到 result.PostStatus，实际：%#v", got.status)
	}
	if got.installOutput != "step 1\nstep 2" {
		t.Fatalf("应记录安装输出，实际：%q", got.installOutput)
	}
}

func TestResultScreenEscReturnsToWelcome(t *testing.T) {
	subject := newModel("workbuddy")
	subject.screen = screenWorkBuddyResult
	subject.selectedProduct = findProduct(subject.products, "workbuddy")

	updated, _ := subject.Update(keyMsg("esc"))
	got := updated.(model)

	if got.screen != screenWorkBuddyWelcome {
		t.Fatalf("Esc 应返回欢迎页，实际：%s", got.screen)
	}
}

func TestResultScreenIKeyStartsInstall(t *testing.T) {
	subject := newModel("workbuddy")
	subject.screen = screenWorkBuddyResult
	subject.platformInfo = platform.Info{OS: platform.Linux}
	subject.status = detector.Status{}
	subject.installOutput = "old output"

	updated, cmd := subject.Update(keyMsg("i"))
	got := updated.(model)

	if got.screen != screenWorkBuddyInstalling {
		t.Fatalf("按 i 后应进入安装页，实际：%s", got.screen)
	}
	if cmd == nil {
		t.Fatalf("应返回安装命令 cmd")
	}
	if got.installOutput != "" {
		t.Fatalf("重新安装前应清空旧输出，实际：%q", got.installOutput)
	}
}

func TestShouldShowInstallOutputOnlyForFailure(t *testing.T) {
	if shouldShowInstallOutput(installer.Result{Outcome: installer.OutcomeInstalled}, nil) {
		t.Fatalf("安装成功时不应默认展示安装输出")
	}
	if !shouldShowInstallOutput(installer.Result{Outcome: installer.OutcomeInstallFailed}, nil) {
		t.Fatalf("安装失败时应展示安装输出")
	}
	if !shouldShowInstallOutput(installer.Result{Outcome: installer.OutcomeInstalled}, assertErr{}) {
		t.Fatalf("存在错误时应展示安装输出")
	}
}

type assertErr struct{}

func (assertErr) Error() string { return "boom" }

func keyMsg(key string) tea.KeyMsg {
	switch key {
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEsc}
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}
	default:
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
	}
}
