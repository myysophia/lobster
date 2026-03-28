package installer

import (
	"errors"
	"os"
	"testing"

	"lobster/internal/platform"
	"lobster/internal/products"
)

type fakeProduct struct {
	installPlan products.InstallPlan
}

func (f fakeProduct) Key() string         { return "fake" }
func (f fakeProduct) DisplayName() string { return "假装产品" }
func (f fakeProduct) Summary() string     { return "不要真实执行" }
func (f fakeProduct) InstallPlan(platform.Info) (products.InstallPlan, error) {
	return f.installPlan, nil
}
func (fakeProduct) DetectPlan(platform.Info) products.DetectPlan {
	return products.DetectPlan{}
}
func (fakeProduct) LaunchPlan(platform.Info) products.LaunchPlan {
	return products.LaunchPlan{}
}

type fakeValidatedProduct struct {
	fakeProduct
	validateErr error
}

func (f fakeValidatedProduct) ValidateInstall(platform.Info) error {
	return f.validateErr
}

func TestRunWithIO_UsesInjectedStreams(t *testing.T) {
	info := platform.Info{OS: platform.Linux, Arch: "amd64"}
	plan := products.InstallPlan{Exec: []string{"fake-binary"}}
	prod := fakeProduct{installPlan: plan}

	t.Cleanup(func() { executePlan = defaultExecutePlan })

	var capturedPlan products.InstallPlan
	var capturedIO ExecIO
	executePlan = func(p products.InstallPlan, streams ExecIO) error {
		capturedPlan = p
		capturedIO = streams
		return nil
	}

	res, err := RunWithIO(info, prod, false, ExecIO{Stdout: os.Stdout})
	if err != nil {
		t.Fatalf("RunWithIO 返回错误: %v", err)
	}

	if res.Executed != true {
		t.Fatalf("期望 Executed=true, got %v", res.Executed)
	}

	if res.Outcome != OutcomeVerifyFailed {
		t.Fatalf("期望 Outcome=%s, got %s", OutcomeVerifyFailed, res.Outcome)
	}

	if capturedPlan.Exec[0] != plan.Exec[0] {
		t.Fatalf("传入计划不一致, got %v", capturedPlan.Exec)
	}

	if capturedIO.Stdout != os.Stdout {
		t.Fatalf("期望 Stdout 被注入, got %v", capturedIO.Stdout)
	}

	if capturedIO.Stderr != os.Stderr {
		t.Fatalf("默认 Stderr 没填充, got %v", capturedIO.Stderr)
	}

	if capturedIO.Stdin != os.Stdin {
		t.Fatalf("默认 Stdin 没填充, got %v", capturedIO.Stdin)
	}
}

func TestRunWithIO_ErrorFromExecutor(t *testing.T) {
	info := platform.Info{OS: platform.Linux, Arch: "amd64"}
	plan := products.InstallPlan{Exec: []string{"fake-binary"}}
	prod := fakeProduct{installPlan: plan}

	t.Cleanup(func() { executePlan = defaultExecutePlan })

	executePlan = func(products.InstallPlan, ExecIO) error {
		return errors.New("boom")
	}

	_, err := RunWithIO(info, prod, false, ExecIO{})
	if err == nil {
		t.Fatalf("期望 RunWithIO 返回错误，实际为 nil")
	}
}

func TestRunWithIO_StopsWhenValidationFails(t *testing.T) {
	info := platform.Info{OS: platform.Windows, Arch: "amd64"}
	prod := fakeValidatedProduct{
		fakeProduct: fakeProduct{
			installPlan: products.InstallPlan{Exec: []string{"fake-binary"}},
		},
		validateErr: errors.New("missing git bash"),
	}

	t.Cleanup(func() { executePlan = defaultExecutePlan })

	called := false
	executePlan = func(products.InstallPlan, ExecIO) error {
		called = true
		return nil
	}

	result, err := RunWithIO(info, prod, false, ExecIO{})
	if err == nil {
		t.Fatalf("校验失败时应返回错误")
	}
	if called {
		t.Fatalf("前置校验失败后不应继续执行安装器")
	}
	if result.Outcome != OutcomeInstallFailed {
		t.Fatalf("应标记为安装失败，实际：%s", result.Outcome)
	}
	if result.Executed {
		t.Fatalf("前置校验失败时不应标记为已执行")
	}
}
