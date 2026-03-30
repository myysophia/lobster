package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"lobster/internal/advisor"
	"lobster/internal/detector"
	"lobster/internal/installer"
	"lobster/internal/launcher"
	"lobster/internal/platform"
	"lobster/internal/products"
	"lobster/internal/tui"
)

type App struct {
	registry *products.Registry
	stdout   io.Writer
	stderr   io.Writer
	runTUI   func(productKey string) error
}

func New() *App {
	return &App{
		registry: products.NewRegistry(),
		stdout:   os.Stdout,
		stderr:   os.Stderr,
		runTUI:   tui.Run,
	}
}

func (a *App) Run(args []string) int {
	if len(args) == 0 {
		a.printGlobalUsage()
		return 0
	}

	switch args[0] {
	case "help", "-h", "--help":
		a.printGlobalUsage()
		return 0
	case "list":
		if err := a.handleList(args[1:]); err != nil {
			return a.fail(err)
		}
		return 0
	case "tui":
		if err := a.handleGlobalTUI(args[1:]); err != nil {
			return a.fail(err)
		}
		return 0
	case "install", "status", "open", "doctor", "next":
		return a.fail(a.oldSyntaxError(args))
	}

	product, err := a.registry.Get(args[0])
	if err != nil {
		return a.fail(err)
	}

	if err := a.handleProductCommand(product, args[1:]); err != nil {
		return a.fail(err)
	}
	return 0
}

func (a *App) handleList(args []string) error {
	if len(args) > 0 {
		return errors.New("list 命令不接受额外参数")
	}

	fmt.Fprintln(a.stdout, "当前支持的产品：")
	for _, key := range a.registry.Keys() {
		product, _ := a.registry.Get(key)
		fmt.Fprintf(a.stdout, "- %s: %s\n", product.Key(), product.Summary())
	}
	return nil
}

func (a *App) handleGlobalTUI(args []string) error {
	if len(args) > 0 {
		return errors.New("tui 命令不接受额外参数")
	}
	return a.runTUI("")
}

func (a *App) handleProductCommand(product products.Product, args []string) error {
	if len(args) == 0 {
		a.printProductUsage(product)
		return nil
	}

	switch args[0] {
	case "help", "-h", "--help":
		a.printProductUsage(product)
		return nil
	case "tui":
		return a.handleProductTUI(product, args[1:])
	case "install":
		return a.handleInstall(product, args[1:])
	case "status":
		return a.handleStatus(product, args[1:])
	case "open":
		return a.handleOpen(product, args[1:])
	case "doctor":
		return a.handleDoctor(product, args[1:])
	case "next":
		return a.handleNext(product, args[1:])
	default:
		return fmt.Errorf("未知动作: %s", args[0])
	}
}

func (a *App) handleProductTUI(product products.Product, args []string) error {
	if len(args) > 0 {
		return errors.New("tui 命令不接受额外参数")
	}
	return a.runTUI(product.Key())
}

func (a *App) handleInstall(product products.Product, args []string) error {
	dryRun := false
	for _, arg := range args {
		if arg == "--dry-run" {
			dryRun = true
			continue
		}
		return fmt.Errorf("未知参数: %s", arg)
	}

	info := platform.Detect()
	fmt.Fprintf(a.stdout, "平台：%s\n", info.String())
	fmt.Fprintf(a.stdout, "目标产品：%s\n", product.DisplayName())

	var output bytes.Buffer
	result, err := installer.RunWithIO(info, product, dryRun, installer.ExecIO{
		Stdin:  os.Stdin,
		Stdout: &output,
		Stderr: &output,
	})
	if result.Plan.Summary != "" {
		fmt.Fprintf(a.stdout, "安装策略：%s\n", result.Plan.Summary)
	}
	if dryRun {
		if err != nil {
			return err
		}
		fmt.Fprintln(a.stdout, "安装命令：")
		fmt.Fprintf(a.stdout, "  %s\n", strings.Join(result.Plan.Exec, " "))
		return nil
	}

	for _, line := range advisor.InstallSummary(product, result) {
		fmt.Fprintln(a.stdout, line)
	}

	if err != nil {
		return err
	}
	if result.Outcome == installer.OutcomeVerifyFailed {
		return errors.New("安装命令已执行，但安装后校验未通过")
	}

	return nil
}

func (a *App) handleStatus(product products.Product, args []string) error {
	if len(args) > 0 {
		return errors.New("status 命令不接受额外参数")
	}

	status := detector.Check(platform.Detect(), product)
	if status.CommandAvailable {
		fmt.Fprintf(a.stdout, "%s：已检测到可用安装\n", product.DisplayName())
	} else if status.Installed {
		fmt.Fprintf(a.stdout, "%s：已检测到安装痕迹，但命令暂不可用\n", product.DisplayName())
	} else {
		fmt.Fprintf(a.stdout, "%s：未检测到安装\n", product.DisplayName())
	}
	if status.CommandPath != "" {
		fmt.Fprintf(a.stdout, "命令路径：%s\n", status.CommandPath)
	}
	for _, warning := range status.Warnings {
		fmt.Fprintf(a.stdout, "提示：%s\n", warning)
	}

	return nil
}

func (a *App) handleOpen(product products.Product, args []string) error {
	if len(args) > 0 {
		return errors.New("open 命令不接受额外参数")
	}

	result, err := launcher.Open(platform.Detect(), product)
	if err != nil {
		return err
	}

	fmt.Fprintf(a.stdout, "已尝试打开 %s：%s\n", product.DisplayName(), strings.Join(result.Method, " "))
	return nil
}

func (a *App) handleDoctor(product products.Product, args []string) error {
	if len(args) > 0 {
		return errors.New("doctor 命令不接受额外参数")
	}

	info := platform.Detect()
	status := detector.Check(info, product)
	for _, line := range advisor.DoctorSummary(product, info, status) {
		fmt.Fprintln(a.stdout, line)
	}
	return nil
}

func (a *App) handleNext(product products.Product, args []string) error {
	if len(args) > 0 {
		return errors.New("next 命令不接受额外参数")
	}

	status := detector.Check(platform.Detect(), product)
	for _, line := range advisor.NextSummary(product, status) {
		fmt.Fprintln(a.stdout, line)
	}
	return nil
}

func (a *App) printGlobalUsage() {
	fmt.Fprintln(a.stdout, "用法：")
	fmt.Fprintln(a.stdout, "  lobster help")
	fmt.Fprintln(a.stdout, "  lobster list")
	fmt.Fprintln(a.stdout, "  lobster tui")
	fmt.Fprintln(a.stdout, "  lobster <product> help")
	fmt.Fprintln(a.stdout, "  lobster <product> install [--dry-run]")
	fmt.Fprintln(a.stdout, "  lobster <product> status")
	fmt.Fprintln(a.stdout, "  lobster <product> open")
	fmt.Fprintln(a.stdout, "  lobster <product> doctor")
	fmt.Fprintln(a.stdout, "  lobster <product> next")
	fmt.Fprintln(a.stdout, "  lobster <product> tui")
}

func (a *App) printProductUsage(product products.Product) {
	fmt.Fprintln(a.stdout, "用法：")
	fmt.Fprintf(a.stdout, "  lobster %s help\n", product.Key())
	fmt.Fprintf(a.stdout, "  lobster %s install [--dry-run]\n", product.Key())
	fmt.Fprintf(a.stdout, "  lobster %s status\n", product.Key())
	fmt.Fprintf(a.stdout, "  lobster %s open\n", product.Key())
	fmt.Fprintf(a.stdout, "  lobster %s doctor\n", product.Key())
	fmt.Fprintf(a.stdout, "  lobster %s next\n", product.Key())
	fmt.Fprintf(a.stdout, "  lobster %s tui\n", product.Key())
}

func (a *App) oldSyntaxError(args []string) error {
	action := args[0]
	if len(args) > 1 {
		if _, err := a.registry.Get(args[1]); err == nil {
			return fmt.Errorf("旧语法已移除，请改用 `lobster %s %s`", args[1], action)
		}
	}
	return fmt.Errorf("旧语法已移除，请改用 `lobster <product> %s`", action)
}

func (a *App) fail(err error) int {
	fmt.Fprintf(a.stderr, "错误：%s\n", err)
	fmt.Fprintln(a.stderr, "试试 `lobster help` 查看可用命令。")
	return 1
}
