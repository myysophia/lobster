package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) viewProductSelect() string {
	hero := uiStyles.hero.Render(strings.Join([]string{
		uiStyles.title.Render("Lobster · 统一安装入口"),
		uiStyles.subtitle.Render("WorkBuddy 已可安装，AutoClaw 与 QoderWork 已接入产品入口。"),
		uiStyles.tagLine.Render("一次启动，后续同一命令即可更新、启动与诊断。"),
	}, "\n"))

	sectionHeader := uiStyles.sectionTitle.Render("可选产品")
	cards := make([]string, 0, len(m.products))
	for index, item := range m.products {
		style := uiStyles.cardMuted
		if item.Available {
			style = uiStyles.card
		}
		if index == m.selectedIndex {
			style = uiStyles.cardFocused
		}

		badgeStyle := uiStyles.badgeMuted
		if item.Available {
			badgeStyle = uiStyles.badgeActive
		}

		card := style.Render(lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				uiStyles.cardTitle.Render(item.DisplayName),
				badgeStyle.Render(item.Badge),
			),
			uiStyles.cardBody.Render(item.Summary),
		))
		cards = append(cards, card)
	}

	infoPanel := uiStyles.infoPanel.Render(strings.Join([]string{
		"当前流程：先确认目标产品，再由 WorkBuddy 向导逐步执行安装和校验。",
		"方向提示：WorkBuddy 优先支持桌面化场景，暂不改变官方安装脚本。",
	}, "\n"))

	footer := renderFooter("↑/↓ 切换产品", "Enter 进入向导", "q 退出")
	return uiStyles.page.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		hero,
		sectionHeader,
		strings.Join(cards, "\n"),
		infoPanel,
		footer,
	))
}
