package detector

import (
	"os"
	"os/exec"

	"lobster/internal/platform"
	"lobster/internal/products"
)

var (
	lookPath = exec.LookPath
	statPath = os.Stat
)

type Status struct {
	Installed        bool
	CommandAvailable bool
	MatchedCommand   string
	CommandPath      string
	HasPathEvidence  bool
	FoundPaths       []string
	Notes            []string
	Warnings         []string
}

func Check(info platform.Info, product products.Product) Status {
	plan := product.DetectPlan(info)
	status := Status{
		Notes: append([]string{}, plan.Notes...),
	}

	for _, cmdName := range plan.Commands {
		path, err := lookPath(cmdName)
		if err == nil {
			status.Installed = true
			status.CommandAvailable = true
			status.MatchedCommand = cmdName
			status.CommandPath = path
			break
		}
	}

	for _, candidate := range plan.Paths {
		if candidate == "" {
			continue
		}
		if _, err := statPath(candidate); err == nil {
			status.Installed = true
			status.HasPathEvidence = true
			status.FoundPaths = append(status.FoundPaths, candidate)
		}
	}

	if status.Installed && !status.CommandAvailable {
		status.Warnings = append(status.Warnings, "检测到安装痕迹，但当前终端还没有识别到可执行命令，可能需要重新打开终端。")
	}

	if !status.Installed {
		status.Warnings = append(status.Warnings, "尚未检测到 WorkBuddy 的安装痕迹。")
	}

	return status
}
