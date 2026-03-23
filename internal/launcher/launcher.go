package launcher

import (
	"fmt"
	"os"
	"os/exec"

	"lobster/internal/platform"
	"lobster/internal/products"
)

type Result struct {
	Launched bool
	Method   []string
}

func Open(info platform.Info, product products.Product) (Result, error) {
	if !info.HasDesktop {
		return Result{}, fmt.Errorf("当前环境缺少桌面会话，暂时无法自动打开应用")
	}

	plan := product.LaunchPlan(info)
	for _, candidate := range plan.ExecCandidates {
		if len(candidate) == 0 {
			continue
		}

		cmd := exec.Command(candidate[0], candidate[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err == nil {
			return Result{
				Launched: true,
				Method:   candidate,
			}, nil
		}
	}

	return Result{}, fmt.Errorf("未找到可用的启动方式，请手动打开 %s", product.DisplayName())
}
