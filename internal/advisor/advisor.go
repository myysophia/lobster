package advisor

import (
	"fmt"
	"strings"

	"lobster/internal/detector"
	"lobster/internal/installer"
	"lobster/internal/platform"
	"lobster/internal/products"
)

func InstallSummary(product products.Product, result installer.Result) []string {
	lines := []string{}

	switch result.Outcome {
	case installer.OutcomeDryRun:
		lines = append(lines, fmt.Sprintf("%s 尚未真正安装，本次只展示安装计划。", product.DisplayName()))
	case installer.OutcomeAlreadyInstalled:
		lines = append(lines, fmt.Sprintf("%s 已可用，本次跳过安装。", product.DisplayName()))
	case installer.OutcomeInstalled:
		lines = append(lines, fmt.Sprintf("%s 安装完成。", product.DisplayName()))
	case installer.OutcomeActionRequired:
		lines = append(lines, fmt.Sprintf("%s 安装入口已打开，请按官方引导完成安装。", product.DisplayName()))
	case installer.OutcomeVerifyFailed:
		lines = append(lines, fmt.Sprintf("%s 安装后仍未检测到明确结果。", product.DisplayName()))
	default:
		lines = append(lines, fmt.Sprintf("%s 安装失败。", product.DisplayName()))
	}

	if result.Outcome == installer.OutcomeActionRequired {
		lines = append(lines, fmt.Sprintf("下一步：完成安装后执行 %s 或 %s。", statusCommandHint(product), doctorCommandHint(product)))
		return lines
	}

	status := result.PostStatus
	if status.CommandAvailable {
		if status.CommandPath != "" {
			lines = append(lines, fmt.Sprintf("可执行命令路径：%s", status.CommandPath))
		}
		lines = append(lines, fmt.Sprintf("下一步：执行 %s。", openCommandHint(product)))
	} else if status.Installed {
		lines = append(lines, fmt.Sprintf("下一步：重开终端后执行 %s。", statusCommandHint(product)))
	} else {
		lines = append(lines, fmt.Sprintf("下一步：执行 %s 或 %s。", statusCommandHint(product), doctorCommandHint(product)))
	}
	return lines
}

func NextSummary(product products.Product, status detector.Status) []string {
	lines := []string{}

	if status.CommandAvailable {
		lines = append(lines, fmt.Sprintf("%s 已检测到安装。", product.DisplayName()))
		if status.CommandPath != "" {
			lines = append(lines, fmt.Sprintf("可执行命令路径：%s", status.CommandPath))
		}
		lines = append(lines, fmt.Sprintf("建议下一步：执行 %s。", openCommandHint(product)))
	} else if status.Installed {
		lines = append(lines, fmt.Sprintf("%s 已检测到安装痕迹，但当前命令还不可用。", product.DisplayName()))
		lines = append(lines, fmt.Sprintf("建议下一步：重开终端后执行 %s。", statusCommandHint(product)))
	} else {
		lines = append(lines, fmt.Sprintf("%s 安装后仍未被检测到。", product.DisplayName()))
		lines = append(lines, fmt.Sprintf("建议下一步：执行 %s 或 %s。", statusCommandHint(product), doctorCommandHint(product)))
	}

	lines = append(lines, status.Warnings...)
	return lines
}

func DoctorSummary(product products.Product, info platform.Info, status detector.Status) []string {
	lines := []string{
		fmt.Sprintf("产品：%s", product.DisplayName()),
	}

	if status.CommandAvailable {
		lines = append(lines, "结论：已检测到可用命令，当前可认为安装可用。")
	} else if status.Installed {
		lines = append(lines, "结论：已检测到安装痕迹，但当前命令还不可用。")
	} else {
		lines = append(lines, "结论：尚未检测到明确安装结果。")
	}

	if status.CommandPath != "" {
		lines = append(lines, fmt.Sprintf("命令路径：%s", status.CommandPath))
	}
	if len(status.FoundPaths) > 0 {
		lines = append(lines, fmt.Sprintf("命中路径：%s", strings.Join(status.FoundPaths, ", ")))
	}
	if !info.HasDesktop {
		lines = append(lines, "环境提示：当前会话未检测到桌面环境，`open` 命令可能无法自动拉起应用。")
	}

	lines = append(lines, status.Notes...)
	lines = append(lines, status.Warnings...)

	if status.CommandAvailable {
		lines = append(lines, fmt.Sprintf("建议：可以直接执行 %s 进行下一步。", openCommandHint(product)))
	} else if status.HasPathEvidence {
		lines = append(lines, fmt.Sprintf("建议：优先重新打开终端或刷新 PATH，再执行 %s。", statusCommandHint(product)))
	} else {
		lines = append(lines, fmt.Sprintf("建议：先执行 %s，安装后再复查状态。", installCommandHint(product)))
	}
	return lines
}

func openCommandHint(product products.Product) string {
	return fmt.Sprintf("`lobster %s open`", product.Key())
}

func statusCommandHint(product products.Product) string {
	return fmt.Sprintf("`lobster %s status`", product.Key())
}

func doctorCommandHint(product products.Product) string {
	return fmt.Sprintf("`lobster %s doctor`", product.Key())
}

func installCommandHint(product products.Product) string {
	return fmt.Sprintf("`lobster %s install`", product.Key())
}
