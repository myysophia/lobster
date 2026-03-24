package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) viewWorkBuddyInstalling() string {
	lines := []string{
		"正在尽量保留在向导内执行官方安装逻辑，并实时等待结果。",
		"如果安装器需要权限或交互输入，请直接在当前终端完成。",
	}
	if m.installResult.Plan.Summary != "" {
		lines = append(lines, "安装策略："+m.installResult.Plan.Summary)
	}

	panel := uiStyles.warnPanel.Render(strings.Join(lines, "\n"))
	hint := uiStyles.tipPanel.Render(strings.Join([]string{
		uiStyles.sectionTitle.Render("安装提示"),
		"准备阶段会验证平台信息与权限，若需要交互，请按官方提示操作。",
		"安装结束后会自动切换到结果页，并保留最近一次安装输出摘要。",
	}, "\n"))
	notice := renderNoticePanel(m.notice, m.err)

	return uiStyles.page.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		uiStyles.title.Render("WorkBuddy 安装中"),
		uiStyles.subtitle.Render("Lobster 会在安装后自动校验，并指导下一步。"),
		panel,
		hint,
		notice,
		renderFooter("q 退出"),
	))
}
