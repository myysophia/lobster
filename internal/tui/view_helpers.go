package tui

import (
	"bytes"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
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
	trimmed := strings.TrimSpace(normalizeInstallOutput(output))
	if trimmed == "" {
		return ""
	}

	lines := strings.Split(trimmed, "\n")
	if len(lines) > 8 {
		lines = append(lines[:8], "...")
	}

	return strings.Join(lines, "\n")
}

func normalizeInstallOutput(output string) string {
	if output == "" {
		return ""
	}

	raw := []byte(strings.ReplaceAll(output, "\r\n", "\n"))
	raw = bytes.ReplaceAll(raw, []byte{'\r'}, []byte{'\n'})

	if utf8.Valid(raw) {
		return string(raw)
	}

	decoded, err := simplifiedchinese.GB18030.NewDecoder().Bytes(raw)
	if err == nil && utf8.Valid(decoded) {
		return string(decoded)
	}

	return string(raw)
}
