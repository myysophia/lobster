package installer

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"lobster/internal/detector"
	"lobster/internal/platform"
	"lobster/internal/products"
)

type Outcome string

const (
	OutcomeDryRun           Outcome = "dry_run"
	OutcomeAlreadyInstalled Outcome = "already_installed"
	OutcomeInstalled        Outcome = "installed"
	OutcomeActionRequired   Outcome = "action_required"
	OutcomeInstallFailed    Outcome = "install_failed"
	OutcomeVerifyFailed     Outcome = "verify_failed"
)

type Result struct {
	Plan          products.InstallPlan
	DryRun        bool
	Executed      bool
	Succeeded     bool
	Outcome       Outcome
	PreStatus     detector.Status
	PostStatus    detector.Status
	VerifyChecked bool
}

// ExecIO 用于在运行官方安装器时注入自定义的标准输入输出
type ExecIO struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

var executePlan = defaultExecutePlan

func (s ExecIO) withDefaults() ExecIO {
	if s.Stdin == nil {
		s.Stdin = os.Stdin
	}
	if s.Stdout == nil {
		s.Stdout = os.Stdout
	}
	if s.Stderr == nil {
		s.Stderr = os.Stderr
	}
	return s
}

func Run(info platform.Info, product products.Product, dryRun bool) (Result, error) {
	return RunWithIO(info, product, dryRun, ExecIO{})
}

func RunWithIO(info platform.Info, product products.Product, dryRun bool, streams ExecIO) (Result, error) {
	preStatus := detector.Check(info, product)

	plan, err := product.InstallPlan(info)
	if err != nil {
		return Result{
			DryRun:     dryRun,
			Outcome:    OutcomeInstallFailed,
			PreStatus:  preStatus,
			PostStatus: preStatus,
		}, err
	}

	if dryRun {
		return Result{
			Plan:          plan,
			DryRun:        true,
			Executed:      false,
			Succeeded:     false,
			Outcome:       OutcomeDryRun,
			PreStatus:     preStatus,
			PostStatus:    preStatus,
			VerifyChecked: false,
		}, nil
	}

	if preStatus.CommandAvailable {
		return Result{
			Plan:          plan,
			DryRun:        false,
			Executed:      false,
			Succeeded:     true,
			Outcome:       OutcomeAlreadyInstalled,
			PreStatus:     preStatus,
			PostStatus:    preStatus,
			VerifyChecked: true,
		}, nil
	}

	if validator, ok := product.(products.InstallValidator); ok {
		if err := validator.ValidateInstall(info); err != nil {
			return Result{
				Plan:          plan,
				DryRun:        false,
				Executed:      false,
				Succeeded:     false,
				Outcome:       OutcomeInstallFailed,
				PreStatus:     preStatus,
				PostStatus:    preStatus,
				VerifyChecked: false,
			}, err
		}
	}

	if len(plan.Exec) == 0 {
		return Result{
			Plan:       plan,
			DryRun:     false,
			Outcome:    OutcomeInstallFailed,
			PreStatus:  preStatus,
			PostStatus: preStatus,
		}, fmt.Errorf("安装计划为空")
	}

	streams = streams.withDefaults()
	if err := executePlan(plan, streams); err != nil {
		postStatus := detector.Check(info, product)
		return Result{
			Plan:          plan,
			DryRun:        false,
			Executed:      true,
			Succeeded:     false,
			Outcome:       OutcomeInstallFailed,
			PreStatus:     preStatus,
			PostStatus:    postStatus,
			VerifyChecked: true,
		}, fmt.Errorf("官方安装器执行失败: %w", err)
	}

	if plan.SkipVerify {
		postStatus := detector.Check(info, product)
		return Result{
			Plan:          plan,
			DryRun:        false,
			Executed:      true,
			Succeeded:     true,
			Outcome:       OutcomeActionRequired,
			PreStatus:     preStatus,
			PostStatus:    postStatus,
			VerifyChecked: false,
		}, nil
	}

	postStatus := detector.Check(info, product)
	outcome := OutcomeVerifyFailed
	succeeded := false
	if postStatus.Installed {
		outcome = OutcomeInstalled
		succeeded = true
	}

	return Result{
		Plan:          plan,
		DryRun:        false,
		Executed:      true,
		Succeeded:     succeeded,
		Outcome:       outcome,
		PreStatus:     preStatus,
		PostStatus:    postStatus,
		VerifyChecked: true,
	}, nil
}

func defaultExecutePlan(plan products.InstallPlan, streams ExecIO) error {
	cmd := exec.Command(plan.Exec[0], plan.Exec[1:]...)
	cmd.Stdout = streams.Stdout
	cmd.Stderr = streams.Stderr
	cmd.Stdin = streams.Stdin
	cmd.Env = append([]string{}, os.Environ()...)
	for key, value := range plan.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}
	return cmd.Run()
}
