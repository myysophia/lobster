package launcher

import (
	"fmt"
	"os"
	"os/exec"

	"lobster/internal/detector"
	"lobster/internal/platform"
	"lobster/internal/products"
)

type Result struct {
	Launched bool
	Method   []string
}

type commandRunner interface {
	Start() error
	SetStdout(*os.File)
	SetStderr(*os.File)
}

type execCmd struct {
	cmd *exec.Cmd
}

func (e execCmd) Start() error {
	return e.cmd.Start()
}

func (e execCmd) SetStdout(file *os.File) {
	e.cmd.Stdout = file
}

func (e execCmd) SetStderr(file *os.File) {
	e.cmd.Stderr = file
}

var execCommand = func(name string, args ...string) commandRunner {
	return execCmd{cmd: exec.Command(name, args...)}
}

func Open(info platform.Info, product products.Product) (Result, error) {
	status := detector.Check(info, product)
	return OpenWithStatus(info, product, status)
}

func OpenWithStatus(info platform.Info, product products.Product, status detector.Status) (Result, error) {
	if !info.HasDesktop {
		return Result{}, fmt.Errorf("当前环境缺少桌面会话，暂时无法自动打开应用")
	}

	if status.CommandAvailable && status.CommandPath != "" {
		cmd := execCommand(status.CommandPath)
		cmd.SetStdout(os.Stdout)
		cmd.SetStderr(os.Stderr)
		if err := cmd.Start(); err == nil {
			return Result{
				Launched: true,
				Method:   []string{status.CommandPath},
			}, nil
		}
	}

	plan := product.LaunchPlan(info)
	for _, candidate := range plan.ExecCandidates {
		if len(candidate) == 0 {
			continue
		}

		cmd := execCommand(candidate[0], candidate[1:]...)
		cmd.SetStdout(os.Stdout)
		cmd.SetStderr(os.Stderr)
		if err := cmd.Start(); err == nil {
			return Result{
				Launched: true,
				Method:   candidate,
			}, nil
		}
	}

	return Result{}, fmt.Errorf("未找到可用的启动方式，请手动打开 %s", product.DisplayName())
}
