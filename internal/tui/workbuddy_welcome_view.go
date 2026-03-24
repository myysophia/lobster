package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) viewWorkBuddyWelcome() string {
	statusTitle := "正在检查当前安装状态..."
	statusLines := []string{
		fmt.Sprintf("当前平台：%s", m.platformInfo.String()),
		fmt.Sprintf("当前产品：%s", m.selectedProduct.DisplayName),
	}

	if m.statusLoaded {
		switch {
		case m.status.CommandAvailable:
			statusTitle = "当前已安装，可直接打开或重新执行安装。"
		case m.status.Installed:
			statusTitle = "检测到安装痕迹，命令尚未可用（可能 PATH 仍未刷新）。"
		default:
			statusTitle = "尚未检测到可用安装结果，可以开始安装。"
		}
	}

	statusLines = append(statusLines, statusTitle)
	if m.status.CommandPath != "" {
		statusLines = append(statusLines, fmt.Sprintf("命令路径：%s", m.status.CommandPath))
	}

	statusPanel := uiStyles.infoPanel.Render(strings.Join(statusLines, "\n"))

	guidancePanel := uiStyles.tipPanel.Render(strings.Join([]string{
		uiStyles.sectionTitle.Render("下一步建议"),
		"1. 确认当前平台与命令可执行权限。",
		"2. 按 Enter 触发官方安装器，全程由官方逻辑执行。",
		"3. 安装完成后 Lobster 会自动进入结果与建议阶段。",
	}, "\n"))

	notice := renderNoticePanel(m.notice, m.err)

	return uiStyles.page.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		uiStyles.title.Render("WorkBuddy 安装向导"),
		uiStyles.subtitle.Render("先检查，再安装，再给出下一步建议。"),
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			statusPanel,
			lipgloss.NewStyle().MarginLeft(2).Render(guidancePanel),
		),
		notice,
		renderFooter("Enter 开始安装", "r 重新检查", welcomeBackHint(m.defaultProduct), "q 退出"),
	))
}

func welcomeBackHint(defaultProduct string) string {
	if defaultProduct == "workbuddy" {
		return "Esc 退出"
	}
	return "Esc 返回产品列表"
}
