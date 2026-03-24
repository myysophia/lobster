package tui

import tea "github.com/charmbracelet/bubbletea"

func Run(defaultProduct string) error {
	program := tea.NewProgram(newModel(defaultProduct), tea.WithAltScreen())
	_, err := program.Run()
	return err
}
