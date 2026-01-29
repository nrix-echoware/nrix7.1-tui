package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	ColorPrimary   = lipgloss.Color("#FF79C6")
	ColorSecondary = lipgloss.Color("#8BE9FD")
	ColorAccent    = lipgloss.Color("#50FA7B")
	ColorWarning   = lipgloss.Color("#FFB86C")
	ColorError     = lipgloss.Color("#FF5555")
	ColorMuted     = lipgloss.Color("#6272A4")
	ColorBg        = lipgloss.Color("#282A36")
	ColorBgLight   = lipgloss.Color("#44475A")
	ColorFg        = lipgloss.Color("#F8F8F2")

	// Base styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			MarginBottom(1)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true).
			Padding(0, 1)

	LoadingStyle = lipgloss.NewStyle().
			Foreground(ColorWarning).
			Bold(true)

	SelectedStyle = lipgloss.NewStyle().
			Background(ColorBgLight).
			Foreground(ColorFg).
			Bold(true).
			Padding(0, 1)

	NormalStyle = lipgloss.NewStyle().
			Foreground(ColorFg).
			Padding(0, 1)

	PriceStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)

	BrandStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary)

	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			MarginTop(1)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)

	// Box styles
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorMuted).
			Padding(1, 2)

	HeaderBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorPrimary).
			Padding(0, 2).
			MarginBottom(1)

	// Notification styles
	NotificationSuccessStyle = lipgloss.NewStyle().
					Background(ColorAccent).
					Foreground(ColorBg).
					Bold(true).
					Padding(0, 2).
					MarginTop(1)

	NotificationErrorStyle = lipgloss.NewStyle().
				Background(ColorError).
				Foreground(ColorFg).
				Bold(true).
				Padding(0, 2).
				MarginTop(1)

	NotificationInfoStyle = lipgloss.NewStyle().
				Background(ColorSecondary).
				Foreground(ColorBg).
				Bold(true).
				Padding(0, 2).
				MarginTop(1)

	// Input styles
	InputStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorMuted).
			Padding(0, 1)

	InputFocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(ColorPrimary).
				Padding(0, 1)

	// Footer style
	FooterStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(ColorMuted).
			PaddingTop(1).
			MarginTop(1)

	// Divider style
	DividerStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	// Product card style
	ProductCardStyle = lipgloss.NewStyle().
				Padding(0, 1)

	ProductCardSelectedStyle = lipgloss.NewStyle().
					Background(ColorBgLight).
					Padding(0, 1)

	// Badge style
	BadgeStyle = lipgloss.NewStyle().
			Background(ColorPrimary).
			Foreground(ColorBg).
			Padding(0, 1).
			Bold(true)

	CartBadgeStyle = lipgloss.NewStyle().
			Background(ColorAccent).
			Foreground(ColorBg).
			Padding(0, 1).
			Bold(true)

	// Option row styles
	OptionLabelStyle = lipgloss.NewStyle().
				Foreground(ColorSecondary).
				Width(12)

	OptionValueStyle = lipgloss.NewStyle().
				Foreground(ColorFg)

	OptionValueSelectedStyle = lipgloss.NewStyle().
					Background(ColorPrimary).
					Foreground(ColorBg).
					Bold(true).
					Padding(0, 1)

	OptionValueUnselectedStyle = lipgloss.NewStyle().
					Foreground(ColorMuted).
					Padding(0, 1)

	OptionRowFocusedStyle = lipgloss.NewStyle().
				Background(ColorBgLight).
				Padding(0, 1)

	OptionRowStyle = lipgloss.NewStyle().
			Padding(0, 1)
)

const AsciiLogo = `
▖ ▖  ▘    ▄▖                  
▛▖▌▛▘▌▚▘  ▙▖▛▘▛▌▛▛▌▛▛▌█▌▛▘▛▘█▌
▌▝▌▌ ▌▞▖  ▙▖▙▖▙▌▌▌▌▌▌▌▙▖▌ ▙▖▙▖
	⚡ Terminal Shop ⚡
`

const SuccessArt = `
╔══════════════════════════════════════╗
    ✓ ✓ ✓  ORDER PLACED!  ✓ ✓ ✓     
╚══════════════════════════════════════╝
`

const LoadingFrames = "⣾⣽⣻⢿⡿⣟⣯⣷"
