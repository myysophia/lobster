package products

import "lobster/internal/platform"

type InstallPlan struct {
	Summary string
	Exec    []string
	Env     map[string]string
}

type DetectPlan struct {
	Commands []string
	Paths    []string
	Notes    []string
}

type LaunchPlan struct {
	ExecCandidates [][]string
	Notes          []string
}

type Product interface {
	Key() string
	DisplayName() string
	Summary() string
	InstallPlan(platform.Info) (InstallPlan, error)
	DetectPlan(platform.Info) DetectPlan
	LaunchPlan(platform.Info) LaunchPlan
}

type InstallValidator interface {
	ValidateInstall(platform.Info) error
}
