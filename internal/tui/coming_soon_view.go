package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) viewComingSoon() string {
	header := lipgloss.JoinVertical(
		lipgloss.Left,
		uiStyles.title.Render("正在筹备中"),
		uiStyles.subtitle.Render("产品消费牌已经登记，正式安装逻辑下一轮上线。"),
	)

	body := uiStyles.warnPanel.Render(strings.Join([]string{
		uiStyles.sectionTitle.Render(m.selectedProduct.DisplayName),
		uiStyles.paragraph.Render("On The Way · 安装能力正在打磨中，第一版暂不开放。"),
		uiStyles.paragraph.Render("建议先使用 WorkBuddy 完成安装流程，以获取最流畅的体验。"),
	}, "\n"))

	hint := uiStyles.tipPanel.Render("选择 WorkBuddy 即可进入真实安装向导，其他产品会在后续版本逐步接入。")

	return uiStyles.page.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		body,
		hint,
		renderFooter("Esc 返回", "Enter 返回", "q 退出"),
	))
}
