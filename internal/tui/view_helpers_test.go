package tui

import (
	"strings"
	"testing"

	"golang.org/x/text/encoding/simplifiedchinese"
)

func TestNormalizeInstallOutputDecodesGB18030(t *testing.T) {
	raw, err := simplifiedchinese.GB18030.NewEncoder().Bytes([]byte("位置"))
	if err != nil {
		t.Fatalf("编码测试数据失败: %v", err)
	}

	got := normalizeInstallOutput(string(raw))
	if got != "位置" {
		t.Fatalf("应正确解码 Windows 中文输出，实际：%q", got)
	}
}

func TestCompactOutputNormalizesCRLF(t *testing.T) {
	output := "line1\r\nline2\r\n"

	got := compactOutput(output)
	if strings.Contains(got, "\r") {
		t.Fatalf("输出中不应保留 \\r，实际：%q", got)
	}
	if got != "line1\nline2" {
		t.Fatalf("应保留规范化后的换行，实际：%q", got)
	}
}
