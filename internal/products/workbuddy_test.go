package products

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"lobster/internal/platform"
)

func TestWorkBuddyWindowsInstallPlanInjectsGitBashEnv(t *testing.T) {
	subject := NewWorkBuddy()
	info := platform.Info{OS: platform.Windows, Arch: "amd64"}

	tempDir := t.TempDir()
	bashPath := filepath.Join(tempDir, "bash.exe")
	if err := os.WriteFile(bashPath, []byte("fake"), 0o644); err != nil {
		t.Fatalf("准备 bash.exe 失败: %v", err)
	}

	originalEnv := os.Getenv("CODEBUDDY_CODE_GIT_BASH_PATH")
	t.Cleanup(func() {
		_ = os.Setenv("CODEBUDDY_CODE_GIT_BASH_PATH", originalEnv)
		workBuddyLookPath = exec.LookPath
		workBuddyStat = os.Stat
	})

	if err := os.Setenv("CODEBUDDY_CODE_GIT_BASH_PATH", bashPath); err != nil {
		t.Fatalf("设置环境变量失败: %v", err)
	}
	workBuddyLookPath = func(string) (string, error) {
		return "", os.ErrNotExist
	}
	workBuddyStat = os.Stat

	plan, err := subject.InstallPlan(info)
	if err != nil {
		t.Fatalf("InstallPlan 返回错误: %v", err)
	}

	if got := plan.Env["CODEBUDDY_CODE_GIT_BASH_PATH"]; got != bashPath {
		t.Fatalf("期望注入 Git Bash 环境变量 %q，实际 %q", bashPath, got)
	}

	command := strings.Join(plan.Exec, " ")
	if !strings.Contains(command, "$ProgressPreference='SilentlyContinue'") {
		t.Fatalf("Windows 安装命令应关闭 PowerShell 进度条，实际：%s", command)
	}
	if !strings.Contains(command, "UTF8Encoding") {
		t.Fatalf("Windows 安装命令应设置 UTF-8 输出，实际：%s", command)
	}
}

func TestWorkBuddyValidateInstallRequiresGitBashOnWindows(t *testing.T) {
	subject := NewWorkBuddy()
	info := platform.Info{OS: platform.Windows, Arch: "amd64"}

	originalEnv := os.Getenv("CODEBUDDY_CODE_GIT_BASH_PATH")
	t.Cleanup(func() {
		_ = os.Setenv("CODEBUDDY_CODE_GIT_BASH_PATH", originalEnv)
		workBuddyLookPath = exec.LookPath
		workBuddyStat = os.Stat
	})

	if err := os.Unsetenv("CODEBUDDY_CODE_GIT_BASH_PATH"); err != nil {
		t.Fatalf("清理环境变量失败: %v", err)
	}
	workBuddyLookPath = func(string) (string, error) {
		return "", os.ErrNotExist
	}
	workBuddyStat = func(string) (os.FileInfo, error) {
		return nil, os.ErrNotExist
	}

	err := subject.ValidateInstall(info)
	if err == nil {
		t.Fatalf("缺少 Git Bash 时应返回错误")
	}
	if !strings.Contains(err.Error(), "Git Bash") {
		t.Fatalf("错误信息应明确提示 Git Bash，实际：%v", err)
	}
}

func TestWorkBuddyDetectPlanIncludesWindowsUserInstallPaths(t *testing.T) {
	subject := NewWorkBuddy()
	info := platform.Info{OS: platform.Windows, Arch: "amd64"}

	originalLocalAppData := os.Getenv("LOCALAPPDATA")
	t.Cleanup(func() {
		_ = os.Setenv("LOCALAPPDATA", originalLocalAppData)
	})

	localAppData := `C:\Users\tester\AppData\Local`
	if err := os.Setenv("LOCALAPPDATA", localAppData); err != nil {
		t.Fatalf("设置 LOCALAPPDATA 失败: %v", err)
	}

	plan := subject.DetectPlan(info)
	joined := strings.Join(plan.Paths, "\n")

	if !strings.Contains(joined, `C:\Users\tester\AppData\Local`) || !strings.Contains(joined, "codebuddy/bin/codebuddy.exe") {
		t.Fatalf("Windows 检测路径应包含用户目录中的 codebuddy.exe，实际：%v", plan.Paths)
	}
	if !strings.Contains(joined, `C:\Users\tester\AppData\Local`) || !strings.Contains(joined, "codebuddy/bin") {
		t.Fatalf("Windows 检测路径应包含用户目录中的 bin 目录，实际：%v", plan.Paths)
	}
}
