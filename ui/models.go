package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"luna/config"
	"luna/git"
	"luna/style"
)

var (
	asciiArt = `
 ██▓     █    ██  ███▄    █  ▄▄▄
▓██▒     ██  ▓██▒ ██ ▀█   █ ▒████▄
▒██░    ▓██  ▒██░▓██  ▀█ ██▒▒██  ▀█▄
▒██░    ▓▓█  ░██░▓██▒  ▐▌██▒░██▄▄▄▄██
░██████▒▒▒█████▓ ▒██░   ▓██░ ▓█   ▓██▒
░ ▒░▓  ░░▒▓▒ ▒ ▒ ░ ▒░   ▒ ▒  ▒▒   ▓▒█░
░ ░ ▒  ░░░▒░ ░ ░ ░ ░░   ░ ▒░  ▒   ▒▒ ░
  ░ ░    ░░░ ░ ░    ░   ░ ░   ░   ▒
    ░  ░   ░              ░       ░  ░


made by: hax (github.com/iuoz)
`
)

type uiState int

const (
	stateLoading uiState = iota
	stateSelecting
	stateStaging
	stateProcessing
	stateReview
	stateComplete
	stateError
)

type CommitUI struct {
	cfg           config.Config
	includeEmoji  bool
	available     []git.FileStatus
	selected      map[string]bool
	cursor        int
	files         []string
	currentFile   int
	commitMsgs    map[string]string
	commitResults map[string]string
	spinner       spinner.Model
	progress      progress.Model
	viewport      viewport.Model
	state         uiState
	err           error
}

type fileProcessedMsg struct {
	filename string
	result   string
	err      error
}

type availableFilesMsg struct {
	files []git.FileStatus
	err   error
}

type stagedMsg struct {
	err error
}

func InitializeCommitUI(cfg config.Config, includeEmoji bool) CommitUI {
	sp := spinner.New()
	sp.Spinner = spinner.Points
	sp.Style = lipgloss.NewStyle().Foreground(style.ColorAccent)

	prog := progress.New(
		progress.WithSolidFill(string(style.ColorPrimary)),
		progress.WithGradient("#7209b7", "#4361ee"),
	)

	vp := viewport.New(80, 20)

	return CommitUI{
		cfg:           cfg,
		includeEmoji:  includeEmoji,
		selected:      make(map[string]bool),
		commitMsgs:    make(map[string]string),
		commitResults: make(map[string]string),
		spinner:       sp,
		progress:      prog,
		viewport:      vp,
		state:         stateLoading,
	}
}

func (m CommitUI) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		loadAvailableFiles(m.cfg),
	)
}

func (m CommitUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		switch m.state {
		case stateSelecting:
			return handleSelectionInput(m, msg)
		case stateReview:
			return handleReviewInput(m, msg)
		}

	case availableFilesMsg:
		if msg.err != nil {
			m.err = msg.err
			m.state = stateError
			return m, nil
		}
		if len(msg.files) == 0 {
			m.state = stateComplete
			return m, tea.Quit
		}
		m.available = msg.files
		m.selected = make(map[string]bool)
		m.cursor = 0
		m.state = stateSelecting
		return m, nil

	case stagedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.state = stateError
			return m, nil
		}
		m.state = stateProcessing
		return m, processNextFile(m)

	case fileProcessedMsg:
		if msg.err != nil {
			m.commitResults[msg.filename] = "Error: " + msg.err.Error()
		} else {
			m.commitResults[msg.filename] = msg.result
		}

		m.commitMsgs[msg.filename] = m.commitMsgs[msg.filename]

		m.state = stateReview
		return m, nil
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func processNextFile(m CommitUI) tea.Cmd {
	if m.currentFile >= len(m.files) {
		return nil
	}

	file := m.files[m.currentFile]

	return func() tea.Msg {
		diff, err := git.GetFileDiff(file)
		if err != nil {
			return fileProcessedMsg{filename: file, result: "", err: err}
		}

		commitMsg := git.GenerateCommitMessage(
			m.cfg.ApiKey,
			diff,
			file,
			m.cfg,
			m.includeEmoji,
		)

		m.commitMsgs[file] = commitMsg

		return fileProcessedMsg{
			filename: file,
			result:   "Commit message generated",
			err:      nil,
		}
	}
}

func (m CommitUI) View() string {
	var b strings.Builder

	now := time.Now().Format("15:04:05")
	header := style.TitleStyle.Render(style.IconMoon + "  Luna Commits  " + style.IconAngleRight + "  " + now)

	switch m.state {

	case stateLoading:
		b.WriteString(header + "\n")
		b.WriteString(style.BoxStyle.Render(m.spinner.View() + "  Loading available files..."))

	case stateSelecting:
		var list strings.Builder
		for i, f := range m.available {
			cursor := "  "
			if i == m.cursor {
				cursor = style.Accent(style.IconAngleRight + " ")
			}

			box := style.Muted(style.IconSquare)
			if m.selected[f.Path] {
				box = style.SuccessStyle.Render(style.IconCheckSquare)
			}

			badge := strings.TrimSpace(f.Status)
			if badge == "" {
				badge = "??"
			}

			list.WriteString(fmt.Sprintf("%s%s  %s  %s\n", cursor, box, style.Muted(badge), f.Path))
		}

		selectedCount := 0
		for _, v := range m.selected {
			if v {
				selectedCount++
			}
		}

		help := style.Muted(fmt.Sprintf(
			"%s%s move   space select   a all   enter stage & commit (%d selected)   q quit",
			style.IconArrowUp, style.IconArrowDown, selectedCount,
		))

		body := fmt.Sprintf(
			"%s\n%s\n\n%s",
			style.Muted("Select files to stage and commit"),
			list.String(),
			help,
		)

		b.WriteString(header + "\n")
		b.WriteString(style.BoxStyle.Render(body))

	case stateStaging:
		b.WriteString(header + "\n")
		b.WriteString(style.BoxStyle.Render(m.spinner.View() + "  Staging selected files..."))

	case stateProcessing:
		progressVal := float64(m.currentFile) / float64(len(m.files))
		box := fmt.Sprintf(
			"%s  Processing\n\n%s\n%s\n\n%s\n%s",
			m.spinner.View(),
			style.Muted("Progress"),
			m.progress.ViewAs(progressVal),
			style.Muted("Current file"),
			style.AccentStyle.Render(m.files[m.currentFile]),
		)

		b.WriteString(header + "\n")
		b.WriteString(style.BoxStyle.Render(box))

	case stateReview:
		file := m.files[m.currentFile]
		msg := m.commitMsgs[file]
		result := m.commitResults[file]

		icon := style.SuccessStyle.Render(style.IconCheck)
		if strings.Contains(result, "Error") {
			icon = style.ErrorStyle.Render(style.IconCross)
		}

		hints := style.Muted(fmt.Sprintf(
			"%s c  confirm    %s r  retry    %s q  quit",
			style.IconCheck, style.IconRefresh, style.IconCross,
		))

		body := fmt.Sprintf(
			"%s\n%s\n\n%s  %s\n\nCommit message:\n%s\n\nStatus: %s  %s\n\n%s",
			asciiArt,
			style.TitleStyle.Render(style.IconMoon+"  github.com/LunaSource/Luna"),
			style.Muted("File"),
			style.AccentStyle.Render(file),
			msg,
			icon, result,
			hints,
		)

		b.WriteString(header + "\n")
		b.WriteString(style.BoxStyle.Render(body))

	case stateComplete:
		b.WriteString(header + "\n")
		b.WriteString(style.SuccessStyle.Render(style.IconRocket + "  All commits completed."))

	case stateError:
		b.WriteString(header + "\n")
		b.WriteString(style.Err(fmt.Sprintf("Error: %v", m.err)))
	}

	return b.String()
}
