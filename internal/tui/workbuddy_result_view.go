package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"lobster/internal/installer"
)

func (m model) viewWorkBuddyResult() string {
	outcomeTitle := "安装结果"
	panelStyle := uiStyles.panel

	switch m.installResult.Outcome {
	case installer.OutcomeAlreadyInstalled:
		outcomeTitle = "已安装 · 跳过重复流程"
		panelStyle = uiStyles.successPanel
	case installer.OutcomeInstalled:
		outcomeTitle = "安装完成 · 欢迎使用"
		panelStyle = uiStyles.successPanel
	case installer.OutcomeVerifyFailed:
		outcomeTitle = "校验未通过 · 需手动排查"
		panelStyle = uiStyles.warnPanel
	case installer.OutcomeInstallFailed:
		outcomeTitle = "安装失败 · 需要重试"
		panelStyle = uiStyles.errorPanel
	}

	lines := []string{outcomeTitle}
	if m.notice != "" {
		lines = append(lines, "", m.notice)
	}

	mainPanel := panelStyle.Render(strings.Join(lines, "\n"))

	nextPanel := ""
	if len(m.nextLines) > 0 {
		nextPanel = uiStyles.infoPanel.Render(strings.Join(append([]string{
			uiStyles.sectionTitle.Render("下一步建议"),
		}, m.nextLines...), "\n"))
	}

	doctorPanel := ""
	if m.showDoctor {
		doctorPanel = uiStyles.panel.Render(strings.Join(append([]string{
			uiStyles.sectionTitle.Render("诊断详情"),
		}, m.doctorLines...), "\n"))
	}

	outputPanel := ""
	if shouldShowInstallOutput(m.installResult, m.err) {
		output := compactOutput(m.installOutput)
		if output != "" {
			outputPanel = uiStyles.tipPanel.Render(strings.Join([]string{
				uiStyles.sectionTitle.Render("最近安装输出"),
				output,
			}, "\n"))
		}
	}

	notice := renderNoticePanel("", m.err)

	return uiStyles.page.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		uiStyles.title.Render("WorkBuddy 向导结果"),
		uiStyles.subtitle.Render("安装完成后将自动校验并展示下一步建议。"),
		mainPanel,
		outputPanel,
		nextPanel,
		doctorPanel,
		notice,
		renderFooter("r 重新检查", "o 打开应用", "d 诊断详情", "i 再装一次", "Esc 返回", "q 退出"),
	))
}

func shouldShowInstallOutput(result installer.Result, err error) bool {
	if err != nil {
		return true
	}
	switch result.Outcome {
	case installer.OutcomeInstallFailed, installer.OutcomeVerifyFailed:
		return true
	default:
		return false
	}
}
