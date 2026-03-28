package launcher

import (
	"errors"
	"os"
	"testing"

	"lobster/internal/detector"
	"lobster/internal/platform"
	"lobster/internal/products"
)

type fakeProduct struct {
	launchPlan products.LaunchPlan
}

func (f fakeProduct) Key() string { return "fake" }

func (f fakeProduct) DisplayName() string { return "Fake" }

func (f fakeProduct) Summary() string { return "fake summary" }

func (f fakeProduct) InstallPlan(platform.Info) (products.InstallPlan, error) {
	return products.InstallPlan{}, nil
}

func (f fakeProduct) DetectPlan(platform.Info) products.DetectPlan {
	return products.DetectPlan{}
}

func (f fakeProduct) LaunchPlan(platform.Info) products.LaunchPlan {
	return f.launchPlan
}

func TestOpenWithStatusPrefersDetectedCommandPath(t *testing.T) {
	originalCommand := execCommand
	t.Cleanup(func() {
		execCommand = originalCommand
	})

	calls := make([][]string, 0, 2)
	execCommand = func(name string, args ...string) commandRunner {
		call := append([]string{name}, args...)
		calls = append(calls, call)
		return fakeCommandRunner{}
	}

	info := platform.Info{OS: platform.Windows, Arch: "amd64", HasDesktop: true}
	product := fakeProduct{
		launchPlan: products.LaunchPlan{
			ExecCandidates: [][]string{{"cmd", "/c", "start", "", "WorkBuddy"}},
		},
	}
	status := detector.Status{
		CommandAvailable: true,
		CommandPath:      `C:\Users\tester\AppData\Local\codebuddy\bin\codebuddy.exe`,
	}

	result, err := OpenWithStatus(info, product, status)
	if err != nil {
		t.Fatalf("OpenWithStatus 返回错误: %v", err)
	}

	if len(calls) != 1 {
		t.Fatalf("期望只执行一次已检测到的命令路径，实际调用次数：%d", len(calls))
	}
	if calls[0][0] != status.CommandPath {
		t.Fatalf("应优先调用检测到的命令路径，实际：%v", calls[0])
	}
	if len(result.Method) != 1 || result.Method[0] != status.CommandPath {
		t.Fatalf("返回的启动方式不符合预期，实际：%v", result.Method)
	}
}

func TestOpenWithStatusFallsBackWhenDetectedCommandFails(t *testing.T) {
	originalCommand := execCommand
	t.Cleanup(func() {
		execCommand = originalCommand
	})

	calls := make([][]string, 0, 3)
	execCommand = func(name string, args ...string) commandRunner {
		call := append([]string{name}, args...)
		calls = append(calls, call)
		if len(calls) == 1 {
			return fakeCommandRunner{startErr: errors.New("boom")}
		}
		return fakeCommandRunner{}
	}

	info := platform.Info{OS: platform.Windows, Arch: "amd64", HasDesktop: true}
	product := fakeProduct{
		launchPlan: products.LaunchPlan{
			ExecCandidates: [][]string{{"cmd", "/c", "start", "", "WorkBuddy"}},
		},
	}
	status := detector.Status{
		CommandAvailable: true,
		CommandPath:      `C:\Users\tester\AppData\Local\codebuddy\bin\codebuddy.exe`,
	}

	result, err := OpenWithStatus(info, product, status)
	if err != nil {
		t.Fatalf("回退到 LaunchPlan 后不应报错: %v", err)
	}

	if len(calls) != 2 {
		t.Fatalf("应先尝试命令路径，再回退 LaunchPlan，实际调用次数：%d", len(calls))
	}
	if calls[1][0] != "cmd" {
		t.Fatalf("第二次应回退到 LaunchPlan，实际：%v", calls[1])
	}
	if result.Method[0] != "cmd" {
		t.Fatalf("返回的启动方式应来自 LaunchPlan，实际：%v", result.Method)
	}
}

type fakeCommandRunner struct {
	startErr error
}

func (f fakeCommandRunner) Start() error {
	return f.startErr
}

func (fakeCommandRunner) SetStdout(*os.File) {}

func (fakeCommandRunner) SetStderr(*os.File) {}
