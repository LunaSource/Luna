// Package style is Luna's single source of truth for colors, icons and
// status rendering, shared by the TUI and every CLI command so the whole
// interface looks and behaves the same way.
//
// Icons are Nerd Font glyphs (Private Use Area codepoints). A Nerd Font
// (e.g. FiraCode Nerd Font, JetBrainsMono NF) must be active in the
// terminal, otherwise they render as blank boxes.
package style

import "github.com/charmbracelet/lipgloss"

const (
	IconCheck       = "" // nf-fa-check
	IconCross       = "" // nf-fa-times
	IconWarning     = "" // nf-fa-exclamation_triangle
	IconInfo        = "" // nf-fa-info_circle
	IconAngleRight  = "" // nf-fa-angle_right (cursor / prompt)
	IconArrowUp     = "" // nf-fa-arrow_up
	IconArrowDown   = "" // nf-fa-arrow_down
	IconKey         = "" // nf-fa-key
	IconChip        = "" // nf-fa-microchip (AI model)
	IconMarker      = "" // nf-fa-map_marker (path / location)
	IconFolderOpen  = "" // nf-fa-folder_open
	IconStar        = "" // nf-fa-star (selected)
	IconStarEmpty   = "" // nf-fa-star_o (unselected)
	IconSquare      = "" // nf-fa-square_o (unchecked checkbox)
	IconCheckSquare = "" // nf-fa-check_square (checked checkbox)
	IconSignOut     = "" // nf-fa-sign_out (exit / cancel)
	IconRocket      = "" // nf-fa-rocket (done)
	IconRefresh     = "" // nf-fa-refresh (retry)
	IconMoon        = "" // nf-fa-moon_o (Luna branding)
)

var (
	ColorPrimary   = lipgloss.Color("#9d4edd")
	ColorSecondary = lipgloss.Color("#5a189a")
	ColorAccent    = lipgloss.Color("#00f5d4")
	ColorSuccess   = lipgloss.Color("#80ed99")
	ColorError     = lipgloss.Color("#ff4d6d")
	ColorWarning   = lipgloss.Color("#f9c74f")
	ColorMuted     = lipgloss.Color("#7b7b8b")
)

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent).
			Background(ColorSecondary).
			Padding(0, 2)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSecondary).
			Padding(1, 2).
			Margin(1)

	SuccessStyle = lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true)
	ErrorStyle   = lipgloss.NewStyle().Foreground(ColorError).Bold(true)
	WarningStyle = lipgloss.NewStyle().Foreground(ColorWarning).Bold(true)
	InfoStyle    = lipgloss.NewStyle().Foreground(ColorAccent)
	MutedStyle   = lipgloss.NewStyle().Foreground(ColorMuted)
	AccentStyle  = lipgloss.NewStyle().Foreground(ColorAccent).Bold(true)
)

func Success(msg string) string { return SuccessStyle.Render(IconCheck + "  " + msg) }
func Err(msg string) string     { return ErrorStyle.Render(IconCross + "  " + msg) }
func Warn(msg string) string    { return WarningStyle.Render(IconWarning + "  " + msg) }
func Info(msg string) string    { return InfoStyle.Render(IconInfo + "  " + msg) }
func Muted(msg string) string   { return MutedStyle.Render(msg) }
func Accent(msg string) string  { return AccentStyle.Render(msg) }
