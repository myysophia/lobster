package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunWithoutArgsPrintsUsage(t *testing.T) {
	app := New("")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	app.stdout = &stdout
	app.stderr = &stderr

	code := app.Run(nil)
	if code != 0 {
		t.Fatalf("无参数运行应返回 0，实际：%d", code)
	}
	if !strings.Contains(stdout.String(), "lobster install <product>") {
		t.Fatalf("应输出 lobster 用法，实际：%s", stdout.String())
	}
	if !strings.Contains(stdout.String(), "lobster tui") {
		t.Fatalf("应输出 lobster tui 用法，实际：%s", stdout.String())
	}
}

func TestRunUnknownCommandReturnsError(t *testing.T) {
	app := New("")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	app.stdout = &stdout
	app.stderr = &stderr

	code := app.Run([]string{"unknown"})
	if code != 1 {
		t.Fatalf("未知命令应返回 1，实际：%d", code)
	}
	if !strings.Contains(stderr.String(), "未知命令: unknown") {
		t.Fatalf("应输出未知命令提示，实际：%s", stderr.String())
	}
	if !strings.Contains(stderr.String(), "lobster help") {
		t.Fatalf("应输出 help 引导，实际：%s", stderr.String())
	}
}

func TestRunHelpForWorkBuddyAlias(t *testing.T) {
	app := New("workbuddy")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	app.stdout = &stdout
	app.stderr = &stderr

	code := app.Run([]string{"help"})
	if code != 0 {
		t.Fatalf("help 应返回 0，实际：%d", code)
	}
	if !strings.Contains(stdout.String(), "wb install [--dry-run]") {
		t.Fatalf("wb 模式应输出 wb 用法，实际：%s", stdout.String())
	}
	if !strings.Contains(stdout.String(), "wb tui") {
		t.Fatalf("wb 模式应输出 wb tui 用法，实际：%s", stdout.String())
	}
}

func TestResolveProduct(t *testing.T) {
	app := New("workbuddy")

	product, err := app.resolveProduct(nil)
	if err != nil {
		t.Fatalf("默认产品存在时不应报错，实际：%v", err)
	}
	if product.Key() != "workbuddy" {
		t.Fatalf("默认产品解析错误，实际：%s", product.Key())
	}

	if _, err := app.resolveProduct([]string{"workbuddy", "extra"}); err == nil {
		t.Fatalf("参数过多时应返回错误")
	}

	emptyApp := New("")
	if _, err := emptyApp.resolveProduct(nil); err == nil {
		t.Fatalf("没有默认产品且无参数时应返回错误")
	}
}

func TestRunTUICommand(t *testing.T) {
	app := New("")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	app.stdout = &stdout
	app.stderr = &stderr

	calledWith := ""
	app.runTUI = func(defaultProduct string) error {
		calledWith = defaultProduct
		return nil
	}

	code := app.Run([]string{"tui"})
	if code != 0 {
		t.Fatalf("tui 命令应返回 0，实际：%d", code)
	}
	if calledWith != "" {
		t.Fatalf("lobster tui 应传入空默认产品，实际：%q", calledWith)
	}
}

func TestRunTUICommandForWB(t *testing.T) {
	app := New("workbuddy")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	app.stdout = &stdout
	app.stderr = &stderr

	calledWith := ""
	app.runTUI = func(defaultProduct string) error {
		calledWith = defaultProduct
		return nil
	}

	code := app.Run([]string{"tui"})
	if code != 0 {
		t.Fatalf("wb tui 命令应返回 0，实际：%d", code)
	}
	if calledWith != "workbuddy" {
		t.Fatalf("wb tui 应传入 workbuddy，实际：%q", calledWith)
	}
}
