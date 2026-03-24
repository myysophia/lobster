package tui

import (
	"strings"
)

func renderFooter(parts ...string) string {
	segments := make([]string, 0, len(parts))
	for _, part := range parts {
		segments = append(segments, uiStyles.hotkey.Render(part))
	}
	return uiStyles.footer.Render(strings.Join(segments, "  ·  "))
}

func renderNoticePanel(notice string, err error) string {
	if err != nil {
		return uiStyles.errorPanel.Render("错误：" + err.Error())
	}
	if notice != "" {
		return uiStyles.successPanel.Render(notice)
	}
	return ""
}

func compactOutput(output string) string {
	trimmed := strings.TrimSpace(output)
	if trimmed == "" {
		return ""
	}

	lines := strings.Split(trimmed, "\n")
	if len(lines) > 8 {
		lines = append(lines[:8], "...")
	}

	return strings.Join(lines, "\n")
}
