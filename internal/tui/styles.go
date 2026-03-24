package tui

import "github.com/charmbracelet/lipgloss"

type styles struct {
	page         lipgloss.Style
	hero         lipgloss.Style
	title        lipgloss.Style
	subtitle     lipgloss.Style
	tagLine      lipgloss.Style
	sectionTitle lipgloss.Style
	metaLine     lipgloss.Style
	card         lipgloss.Style
	cardFocused  lipgloss.Style
	cardMuted    lipgloss.Style
	cardTitle    lipgloss.Style
	cardBody     lipgloss.Style
	badgeActive  lipgloss.Style
	badgeMuted   lipgloss.Style
	infoPanel    lipgloss.Style
	tipPanel     lipgloss.Style
	noticePanel  lipgloss.Style
	panel        lipgloss.Style
	successPanel lipgloss.Style
	warnPanel    lipgloss.Style
	errorPanel   lipgloss.Style
	footer       lipgloss.Style
	hotkey       lipgloss.Style
	paragraph    lipgloss.Style
}

var uiStyles = newStyles()

func newStyles() styles {
	return styles{
		page: lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.HiddenBorder()).
			Background(lipgloss.Color("#010A16")).
			Foreground(lipgloss.Color("#E2E8F0")),
		hero: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#2563EB")).
			Background(lipgloss.Color("#0F172A")).
			Padding(1, 2).
			MarginBottom(1),
		title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F8FAFC")).
			MarginBottom(0),
		subtitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#94A3B8")).
			MarginBottom(1),
		tagLine: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A5B4FC")).
			Italic(true).
			MarginBottom(1),
		sectionTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#A7F3D0")).
			MarginBottom(1),
		metaLine: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#64748B")).
			MarginBottom(1),
		card: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#1E293B")).
			Background(lipgloss.Color("#020617")).
			Padding(1, 2).
			MarginTop(1),
		cardFocused: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#2563EB")).
			Background(lipgloss.Color("#0F172A")).
			Padding(1, 2).
			MarginTop(1),
		cardMuted: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#1E293B")).
			Foreground(lipgloss.Color("#94A3B8")).
			Background(lipgloss.Color("#020617")).
			Padding(1, 2).
			MarginTop(1),
		cardTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#E0F2FE")).
			MarginBottom(0),
		cardBody: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CBD5F5")).
			MarginTop(0),
		badgeActive: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#0F766E")).
			Background(lipgloss.Color("#D9F99D")).
			Padding(0, 1).
			MarginLeft(1),
		badgeMuted: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#A855F7")).
			Background(lipgloss.Color("#EDE9FE")).
			Padding(0, 1).
			MarginLeft(1),
		infoPanel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#94A3B8")).
			Background(lipgloss.Color("#020617")).
			Padding(1, 2).
			MarginTop(1),
		tipPanel: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#1E293B")).
			Background(lipgloss.Color("#0B1220")).
			Padding(1, 2).
			MarginTop(1),
		noticePanel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#94A3B8")).
			Background(lipgloss.Color("#020617")).
			Padding(1, 2).
			MarginTop(1),
		panel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#475569")).
			Padding(1, 2).
			MarginTop(1),
		successPanel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#16A34A")).
			Background(lipgloss.Color("#022B17")).
			Foreground(lipgloss.Color("#ECFDF5")).
			Padding(1, 2).
			MarginTop(1),
		warnPanel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#D97706")).
			Background(lipgloss.Color("#1A1207")).
			Foreground(lipgloss.Color("#FDE68A")).
			Padding(1, 2).
			MarginTop(1),
		errorPanel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#DC2626")).
			Background(lipgloss.Color("#2B0202")).
			Foreground(lipgloss.Color("#FECACA")).
			Padding(1, 2).
			MarginTop(1),
		footer: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#94A3B8")).
			MarginTop(1),
		hotkey: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#38BDF8")),
		paragraph: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CBD5F5")),
	}
}
