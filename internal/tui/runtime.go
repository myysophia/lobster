package tui

import (
	"bytes"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"lobster/internal/advisor"
	"lobster/internal/detector"
	"lobster/internal/installer"
	"lobster/internal/launcher"
	"lobster/internal/platform"
	"lobster/internal/products"
)

func resolveProduct(item productItem) (products.Product, error) {
	return products.NewRegistry().Get(item.Key)
}

func detectProductCmd(item productItem) tea.Cmd {
	return func() tea.Msg {
		info := platform.Detect()
		product, err := resolveProduct(item)
		if err != nil {
			return statusCheckedMsg{info: info, err: err}
		}
		status := detector.Check(info, product)
		return statusCheckedMsg{
			info:   info,
			status: status,
		}
	}
}

func installProductCmd(item productItem, preInfo platform.Info) tea.Cmd {
	return func() tea.Msg {
		product, err := resolveProduct(item)
		if err != nil {
			return installFinishedMsg{info: preInfo, err: err}
		}

		info := preInfo
		if info.OS == "" {
			info = platform.Detect()
		}

		var output bytes.Buffer
		result, runErr := installer.RunWithIO(info, product, false, installer.ExecIO{
			Stdin:  os.Stdin,
			Stdout: &output,
			Stderr: &output,
		})

		return installFinishedMsg{
			info:   platform.Detect(),
			result: result,
			output: output.String(),
			err:    runErr,
		}
	}
}

func openProductCmd(item productItem) tea.Cmd {
	return func() tea.Msg {
		info := platform.Detect()
		product, err := resolveProduct(item)
		if err != nil {
			return openFinishedMsg{err: err}
		}
		result, err := launcher.Open(info, product)
		return openFinishedMsg{
			result: result,
			err:    err,
		}
	}
}

func buildAdvice(item productItem, info platform.Info, status detector.Status) ([]string, []string, error) {
	product, err := resolveProduct(item)
	if err != nil {
		return nil, nil, err
	}
	return advisor.NextSummary(product, status), advisor.DoctorSummary(product, info, status), nil
}
