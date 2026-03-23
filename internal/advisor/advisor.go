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
		lines = append(lines, fmt.Sprintf("%s 还未真正执行安装，本次只展示了安装计划。", product.DisplayName()))
	case installer.OutcomeAlreadyInstalled:
		lines = append(lines, fmt.Sprintf("%s 已经处于可用状态，本次跳过重复安装。", product.DisplayName()))
	case installer.OutcomeInstalled:
		lines = append(lines, fmt.Sprintf("%s 安装完成，并已通过安装后检测。", product.DisplayName()))
	case installer.OutcomeVerifyFailed:
		lines = append(lines, fmt.Sprintf("%s 的安装命令已执行，但安装后仍未检测到明确结果。", product.DisplayName()))
	default:
		lines = append(lines, fmt.Sprintf("%s 安装未完成。", product.DisplayName()))
	}

	status := result.PostStatus
	if status.Installed {
		if status.CommandPath != "" {
			lines = append(lines, fmt.Sprintf("可执行命令路径：%s", status.CommandPath))
		}
		lines = append(lines, fmt.Sprintf("建议下一步：运行 %s，然后在 %s 内完成微信或企业渠道绑定。", openCommandHint(product), product.DisplayName()))
	} else {
		lines = append(lines, fmt.Sprintf("建议下一步：先重新打开终端，再执行 %s 或 %s。", statusCommandHint(product), doctorCommandHint(product)))
	}

	lines = append(lines, status.Warnings...)
	return lines
}

func NextSummary(product products.Product, status detector.Status) []string {
	lines := []string{}

	if status.CommandAvailable {
		lines = append(lines, fmt.Sprintf("%s 已检测到安装。", product.DisplayName()))
		if status.CommandPath != "" {
			lines = append(lines, fmt.Sprintf("可执行命令路径：%s", status.CommandPath))
		}
		lines = append(lines, fmt.Sprintf("建议下一步：运行 %s，然后在 %s 内完成微信或企业渠道绑定。", openCommandHint(product), product.DisplayName()))
	} else if status.Installed {
		lines = append(lines, fmt.Sprintf("%s 已检测到安装痕迹，但当前命令还不可用。", product.DisplayName()))
		lines = append(lines, fmt.Sprintf("建议下一步：先重新打开终端，让 PATH 生效后再执行 %s。", statusCommandHint(product)))
	} else {
		lines = append(lines, fmt.Sprintf("%s 安装后仍未被检测到。", product.DisplayName()))
		lines = append(lines, fmt.Sprintf("建议下一步：先重新打开终端，再执行 %s 或 %s。", statusCommandHint(product), doctorCommandHint(product)))
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
	if product.Key() == "workbuddy" {
		return "`lobster open workbuddy` 或 `wb open`"
	}
	return fmt.Sprintf("`lobster open %s`", product.Key())
}

func statusCommandHint(product products.Product) string {
	if product.Key() == "workbuddy" {
		return "`lobster status workbuddy` 或 `wb status`"
	}
	return fmt.Sprintf("`lobster status %s`", product.Key())
}

func doctorCommandHint(product products.Product) string {
	if product.Key() == "workbuddy" {
		return "`lobster doctor workbuddy` 或 `wb doctor`"
	}
	return fmt.Sprintf("`lobster doctor %s`", product.Key())
}

func installCommandHint(product products.Product) string {
	if product.Key() == "workbuddy" {
		return "`lobster install workbuddy` 或 `wb install`"
	}
	return fmt.Sprintf("`lobster install %s`", product.Key())
}
