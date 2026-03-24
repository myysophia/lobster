package tui

import (
	"lobster/internal/detector"
	"lobster/internal/installer"
	"lobster/internal/launcher"
	"lobster/internal/platform"
)

type statusCheckedMsg struct {
	info   platform.Info
	status detector.Status
	err    error
}

type installFinishedMsg struct {
	info   platform.Info
	result installer.Result
	output string
	err    error
}

type openFinishedMsg struct {
	result launcher.Result
	err    error
}
