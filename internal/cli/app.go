package cli

import (
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
	defaultProduct string
	registry       *products.Registry
	stdout         io.Writer
	stderr         io.Writer
	runTUI         func(defaultProduct string) error
}

func New(defaultProduct string) *App {
	return &App{
		defaultProduct: defaultProduct,
		registry:       products.NewRegistry(),
		stdout:         os.Stdout,
		stderr:         os.Stderr,
		runTUI:         tui.Run,
	}
}

func (a *App) Run(args []string) int {
	if len(args) == 0 {
		a.printUsage()
		return 0
	}

	cmd := args[0]
	switch cmd {
	case "help", "-h", "--help":
		a.printUsage()
		return 0
	case "list":
		a.handleList()
		return 0
	case "tui":
		if err := a.handleTUI(args[1:]); err != nil {
			return a.fail(err)
		}
		return 0
	case "install":
		if err := a.handleInstall(args[1:]); err != nil {
			return a.fail(err)
		}
		return 0
	case "status":
		if err := a.handleStatus(args[1:]); err != nil {
			return a.fail(err)
		}
		return 0
	case "open":
		if err := a.handleOpen(args[1:]); err != nil {
			return a.fail(err)
		}
		return 0
	case "doctor":
		if err := a.handleDoctor(args[1:]); err != nil {
			return a.fail(err)
		}
		return 0
	case "next":
		if err := a.handleNext(args[1:]); err != nil {
			return a.fail(err)
		}
		return 0
	default:
		return a.fail(fmt.Errorf("未知命令: %s", cmd))
	}
}

func (a *App) handleList() {
	fmt.Fprintln(a.stdout, "当前支持的产品：")
	for _, key := range a.registry.Keys() {
		product, _ := a.registry.Get(key)
		fmt.Fprintf(a.stdout, "- %s: %s\n", product.Key(), product.Summary())
	}
}

func (a *App) handleTUI(args []string) error {
	if len(args) > 0 {
		return errors.New("tui 命令暂不接受额外参数")
	}
	return a.runTUI(a.defaultProduct)
}

func (a *App) handleInstall(args []string) error {
	dryRun := false
	productArgs := make([]string, 0, len(args))
	for _, arg := range args {
		if arg == "--dry-run" {
			dryRun = true
			continue
		}
		productArgs = append(productArgs, arg)
	}

	product, err := a.resolveProduct(productArgs)
	if err != nil {
		return err
	}

	info := platform.Detect()
	fmt.Fprintf(a.stdout, "平台：%s\n", info.String())
	fmt.Fprintf(a.stdout, "目标产品：%s\n", product.DisplayName())

	result, err := installer.Run(info, product, dryRun)
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

func (a *App) handleStatus(args []string) error {
	product, err := a.resolveProduct(args)
	if err != nil {
		return err
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

func (a *App) handleOpen(args []string) error {
	product, err := a.resolveProduct(args)
	if err != nil {
		return err
	}

	result, err := launcher.Open(platform.Detect(), product)
	if err != nil {
		return err
	}

	fmt.Fprintf(a.stdout, "已尝试打开 %s：%s\n", product.DisplayName(), strings.Join(result.Method, " "))
	return nil
}

func (a *App) handleDoctor(args []string) error {
	product, err := a.resolveProduct(args)
	if err != nil {
		return err
	}

	info := platform.Detect()
	status := detector.Check(info, product)
	for _, line := range advisor.DoctorSummary(product, info, status) {
		fmt.Fprintln(a.stdout, line)
	}
	return nil
}

func (a *App) handleNext(args []string) error {
	product, err := a.resolveProduct(args)
	if err != nil {
		return err
	}

	status := detector.Check(platform.Detect(), product)
	for _, line := range advisor.NextSummary(product, status) {
		fmt.Fprintln(a.stdout, line)
	}
	return nil
}

func (a *App) resolveProduct(args []string) (products.Product, error) {
	if len(args) > 1 {
		return nil, errors.New("参数过多，请只传一个产品名")
	}

	key := a.defaultProduct
	if len(args) == 1 {
		key = args[0]
	}

	if key == "" {
		return nil, errors.New("请指定产品名，例如：lobster install workbuddy")
	}

	return a.registry.Get(key)
}

func (a *App) printUsage() {
	if a.defaultProduct == "workbuddy" {
		fmt.Fprintln(a.stdout, "用法：")
		fmt.Fprintln(a.stdout, "  wb tui")
		fmt.Fprintln(a.stdout, "  wb install [--dry-run]")
		fmt.Fprintln(a.stdout, "  wb status")
		fmt.Fprintln(a.stdout, "  wb open")
		fmt.Fprintln(a.stdout, "  wb doctor")
		fmt.Fprintln(a.stdout, "  wb next")
		return
	}

	fmt.Fprintln(a.stdout, "用法：")
	fmt.Fprintln(a.stdout, "  lobster tui")
	fmt.Fprintln(a.stdout, "  lobster list")
	fmt.Fprintln(a.stdout, "  lobster install <product> [--dry-run]")
	fmt.Fprintln(a.stdout, "  lobster status <product>")
	fmt.Fprintln(a.stdout, "  lobster open <product>")
	fmt.Fprintln(a.stdout, "  lobster doctor <product>")
	fmt.Fprintln(a.stdout, "  lobster next <product>")
}

func (a *App) fail(err error) int {
	fmt.Fprintf(a.stderr, "错误：%s\n", err)
	if a.defaultProduct == "workbuddy" {
		fmt.Fprintln(a.stderr, "试试 `wb help` 查看可用命令。")
	} else {
		fmt.Fprintln(a.stderr, "试试 `lobster help` 查看可用命令。")
	}
	return 1
}
