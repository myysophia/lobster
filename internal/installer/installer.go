package installer

import (
	"fmt"
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

func Run(info platform.Info, product products.Product, dryRun bool) (Result, error) {
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

	if len(plan.Exec) == 0 {
		return Result{
			Plan:       plan,
			DryRun:     false,
			Outcome:    OutcomeInstallFailed,
			PreStatus:  preStatus,
			PostStatus: preStatus,
		}, fmt.Errorf("安装计划为空")
	}

	cmd := exec.Command(plan.Exec[0], plan.Exec[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
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
