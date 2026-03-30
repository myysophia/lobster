package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunWithoutArgsPrintsGlobalUsage(t *testing.T) {
	app := New()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	app.stdout = &stdout
	app.stderr = &stderr

	code := app.Run(nil)
	if code != 0 {
		t.Fatalf("无参数运行应返回 0，实际：%d", code)
	}
	if !strings.Contains(stdout.String(), "lobster <product> install [--dry-run]") {
		t.Fatalf("应输出新的产品子命令用法，实际：%s", stdout.String())
	}
	if strings.Contains(stdout.String(), "wb ") {
		t.Fatalf("帮助文案不应再出现 wb，实际：%s", stdout.String())
	}
}

func TestRunUnknownCommandReturnsProductError(t *testing.T) {
	app := New()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	app.stdout = &stdout
	app.stderr = &stderr

	code := app.Run([]string{"unknown"})
	if code != 1 {
		t.Fatalf("未知产品应返回 1，实际：%d", code)
	}
	if !strings.Contains(stderr.String(), "不支持的产品: unknown") {
		t.Fatalf("应输出未知产品提示，实际：%s", stderr.String())
	}
}

func TestOldSyntaxShowsMigrationHint(t *testing.T) {
	app := New()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	app.stdout = &stdout
	app.stderr = &stderr

	code := app.Run([]string{"install", "workbuddy"})
	if code != 1 {
		t.Fatalf("旧语法应返回 1，实际：%d", code)
	}
	if !strings.Contains(stderr.String(), "旧语法已移除，请改用 `lobster workbuddy install`") {
		t.Fatalf("应提示新语法，实际：%s", stderr.String())
	}
}

func TestProductHelp(t *testing.T) {
	app := New()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	app.stdout = &stdout
	app.stderr = &stderr

	code := app.Run([]string{"workbuddy", "help"})
	if code != 0 {
		t.Fatalf("产品 help 应返回 0，实际：%d", code)
	}
	if !strings.Contains(stdout.String(), "lobster workbuddy install [--dry-run]") {
		t.Fatalf("应输出 workbuddy 子命令用法，实际：%s", stdout.String())
	}
	if strings.Contains(stdout.String(), "wb ") {
		t.Fatalf("产品帮助中不应再出现 wb，实际：%s", stdout.String())
	}
}

func TestRunGlobalTUICommand(t *testing.T) {
	app := New()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	app.stdout = &stdout
	app.stderr = &stderr

	calledWith := ""
	app.runTUI = func(productKey string) error {
		calledWith = productKey
		return nil
	}

	code := app.Run([]string{"tui"})
	if code != 0 {
		t.Fatalf("全局 tui 命令应返回 0，实际：%d", code)
	}
	if calledWith != "" {
		t.Fatalf("lobster tui 应传入空产品 key，实际：%q", calledWith)
	}
}

func TestRunProductTUICommand(t *testing.T) {
	app := New()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	app.stdout = &stdout
	app.stderr = &stderr

	calledWith := ""
	app.runTUI = func(productKey string) error {
		calledWith = productKey
		return nil
	}

	code := app.Run([]string{"workbuddy", "tui"})
	if code != 0 {
		t.Fatalf("产品 tui 命令应返回 0，实际：%d", code)
	}
	if calledWith != "workbuddy" {
		t.Fatalf("lobster workbuddy tui 应传入 workbuddy，实际：%q", calledWith)
	}
}
