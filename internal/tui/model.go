package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"lobster/internal/detector"
	"lobster/internal/installer"
	"lobster/internal/platform"
)

type screen string

const (
	screenProductSelect       screen = "product_select"
	screenComingSoon          screen = "coming_soon"
	screenWorkBuddyWelcome    screen = "workbuddy_welcome"
	screenWorkBuddyInstalling screen = "workbuddy_installing"
	screenWorkBuddyResult     screen = "workbuddy_result"
)

type model struct {
	defaultProduct  string
	screen          screen
	width           int
	height          int
	products        []productItem
	selectedIndex   int
	selectedProduct productItem
	platformInfo    platform.Info
	status          detector.Status
	statusLoaded    bool
	installResult   installer.Result
	installOutput   string
	nextLines       []string
	doctorLines     []string
	loading         bool
	showDoctor      bool
	err             error
	notice          string
}

func newModel(defaultProduct string) model {
	items := defaultProducts()
	selected := items[0]
	currentScreen := screenProductSelect
	currentIndex := 0

	if defaultProduct != "" {
		selected = findProduct(items, defaultProduct)
		if selected.Key == "workbuddy" {
			currentScreen = screenWorkBuddyWelcome
		} else if selected.Key == defaultProduct {
			currentScreen = screenComingSoon
		}
	}

	for index, item := range items {
		if item.Key == selected.Key {
			currentIndex = index
			break
		}
	}

	return model{
		defaultProduct:  defaultProduct,
		screen:          currentScreen,
		products:        items,
		selectedIndex:   currentIndex,
		selectedProduct: selected,
	}
}

func (m model) Init() tea.Cmd {
	if m.screen == screenWorkBuddyWelcome {
		return detectProductCmd(m.selectedProduct)
	}
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case statusCheckedMsg:
		m.loading = false
		m.platformInfo = msg.info
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.statusLoaded = true
		m.status = msg.status
		m.refreshAdvice()
		return m, nil
	case installFinishedMsg:
		m.loading = false
		m.platformInfo = msg.info
		m.installResult = msg.result
		m.installOutput = msg.output
		m.status = msg.result.PostStatus
		m.statusLoaded = true
		m.err = msg.err
		m.refreshAdvice()
		m.screen = screenWorkBuddyResult
		return m, nil
	case openFinishedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.notice = fmt.Sprintf("已尝试打开 %s：%s", m.selectedProduct.DisplayName, strings.Join(msg.result.Method, " "))
		return m, nil
	case tea.KeyMsg:
		switch {
		case keyMatches(msg, "ctrl+c", "q"):
			return m, tea.Quit
		}

		switch m.screen {
		case screenProductSelect:
			return m.updateProductSelect(msg)
		case screenComingSoon:
			return m.updateComingSoon(msg)
		case screenWorkBuddyWelcome:
			return m.updateWorkBuddyWelcome(msg)
		case screenWorkBuddyInstalling:
			return m, nil
		case screenWorkBuddyResult:
			return m.updateWorkBuddyResult(msg)
		}
	}

	return m, nil
}

func (m model) View() string {
	switch m.screen {
	case screenProductSelect:
		return m.viewProductSelect()
	case screenComingSoon:
		return m.viewComingSoon()
	case screenWorkBuddyWelcome:
		return m.viewWorkBuddyWelcome()
	case screenWorkBuddyInstalling:
		return m.viewWorkBuddyInstalling()
	case screenWorkBuddyResult:
		return m.viewWorkBuddyResult()
	default:
		return ""
	}
}

func (m model) updateProductSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case keyMatches(msg, "up", "k"):
		if m.selectedIndex > 0 {
			m.selectedIndex--
		}
		m.selectedProduct = m.products[m.selectedIndex]
	case keyMatches(msg, "down", "j"):
		if m.selectedIndex < len(m.products)-1 {
			m.selectedIndex++
		}
		m.selectedProduct = m.products[m.selectedIndex]
	case keyMatches(msg, "enter"):
		m.selectedProduct = m.products[m.selectedIndex]
		m.notice = ""
		m.err = nil
		if m.selectedProduct.Available && m.selectedProduct.Key == "workbuddy" {
			m.screen = screenWorkBuddyWelcome
			m.loading = true
			return m, detectProductCmd(m.selectedProduct)
		}
		m.screen = screenComingSoon
	}
	return m, nil
}

func (m model) updateComingSoon(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case keyMatches(msg, "esc", "enter"):
		if m.defaultProduct != "" {
			return m, tea.Quit
		}
		m.screen = screenProductSelect
	}
	return m, nil
}

func (m model) updateWorkBuddyWelcome(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case keyMatches(msg, "esc"):
		if m.defaultProduct == "workbuddy" {
			return m, tea.Quit
		}
		m.screen = screenProductSelect
		return m, nil
	case keyMatches(msg, "r"):
		m.loading = true
		m.notice = ""
		m.err = nil
		m.installOutput = ""
		return m, detectProductCmd(m.selectedProduct)
	case keyMatches(msg, "enter"):
		m.loading = true
		m.notice = ""
		m.err = nil
		m.installOutput = ""
		m.showDoctor = false
		m.screen = screenWorkBuddyInstalling
		return m, installProductCmd(m.selectedProduct, m.platformInfo)
	}
	return m, nil
}

func (m model) updateWorkBuddyResult(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case keyMatches(msg, "esc"):
		m.screen = screenWorkBuddyWelcome
		m.showDoctor = false
		return m, nil
	case keyMatches(msg, "r"):
		m.loading = true
		m.notice = ""
		m.err = nil
		m.installOutput = ""
		return m, detectProductCmd(m.selectedProduct)
	case keyMatches(msg, "o"):
		m.loading = true
		m.notice = ""
		m.err = nil
		return m, openProductCmd(m.selectedProduct)
	case keyMatches(msg, "d"):
		m.showDoctor = !m.showDoctor
	case keyMatches(msg, "i"):
		m.loading = true
		m.notice = ""
		m.err = nil
		m.installOutput = ""
		m.showDoctor = false
		m.screen = screenWorkBuddyInstalling
		return m, installProductCmd(m.selectedProduct, m.platformInfo)
	}
	return m, nil
}

func (m *model) refreshAdvice() {
	nextLines, doctorLines, err := buildAdvice(m.selectedProduct, m.platformInfo, m.status)
	if err != nil {
		m.err = err
		return
	}
	m.nextLines = nextLines
	m.doctorLines = doctorLines
}

func keyMatches(msg tea.KeyMsg, keys ...string) bool {
	for _, key := range keys {
		if msg.String() == key {
			return true
		}
	}
	return false
}
