package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) viewWorkBuddyInstalling() string {
	lines := []string{
		"正在安装 WorkBuddy，请稍候。",
		"如果官方安装器需要权限或交互，请直接在当前终端完成。",
	}
	if m.installResult.Plan.Summary != "" {
		lines = append(lines, "安装策略："+m.installResult.Plan.Summary)
	}

	panel := uiStyles.warnPanel.Render(strings.Join(lines, "\n"))
	hint := uiStyles.tipPanel.Render(strings.Join([]string{
		uiStyles.sectionTitle.Render("安装提示"),
		"安装结束后会自动校验结果。",
		"默认只展示简要结果，详细排查请查看诊断详情。",
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
