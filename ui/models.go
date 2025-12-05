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
)

var (
	colorPrimary   = lipgloss.Color("#9d4edd")
	colorSecondary = lipgloss.Color("#5a189a")
	colorAccent    = lipgloss.Color("#00f5d4")
	colorSuccess   = lipgloss.Color("#80ed99")
	colorError     = lipgloss.Color("#ff4d6d")
	// colorWarning   = lipgloss.Color("#f9c74f")
	// colorText      = lipgloss.Color("#d8d8d8")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorAccent).
			Background(colorSecondary).
			Padding(0, 2)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorSecondary).
			Padding(1, 2).
			Margin(1)

	statusSuccess = lipgloss.NewStyle().Foreground(colorSuccess)
	statusError   = lipgloss.NewStyle().Foreground(colorError)
	// statusWarning = lipgloss.NewStyle().Foreground(colorWarning)
	// statusInfo    = lipgloss.NewStyle().Foreground(colorAccent)

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
                                      

made by: hax (github.com/i1lo)
`
)

type uiState int

const (
	stateLoading uiState = iota
	stateProcessing
	stateReview
	stateComplete
	stateError
)

type CommitUI struct {
	cfg           config.Config
	includeEmoji  bool
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

type loadedFilesMsg struct {
	files []string
	err   error
}

func InitializeCommitUI(cfg config.Config, includeEmoji bool) CommitUI {
	sp := spinner.New()
	sp.Spinner = spinner.Points
	sp.Style = lipgloss.NewStyle().Foreground(colorAccent)

	prog := progress.New(
		progress.WithSolidFill(string(colorPrimary)),
		progress.WithGradient("#7209b7", "#4361ee"),
	)

	vp := viewport.New(80, 20)

	return CommitUI{
		cfg:           cfg,
		includeEmoji:  includeEmoji,
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
		loadFiles,
	)
}

func (m CommitUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		if m.state == stateReview {
			return handleReviewInput(m, msg)
		}

	case loadedFilesMsg:
		if msg.err != nil {
			m.err = msg.err
			m.state = stateError
			return m, nil
		}
		if len(msg.files) == 0 {
			m.state = stateComplete
			return m, tea.Quit
		}
		m.files = msg.files
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
		if git.ShouldIgnoreFile(file, m.cfg) {
			return fileProcessedMsg{filename: file, result: "Ignored", err: nil}
		}

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
	header := titleStyle.Render(" Luna Commits  •  " + now)

	switch m.state {

	case stateLoading:
		b.WriteString(header + "\n")
		b.WriteString(boxStyle.Render(m.spinner.View() + " Loading files..."))

	case stateProcessing:
		progressVal := float64(m.currentFile) / float64(len(m.files))
		box := fmt.Sprintf(
			"%s Processing...\n\nProgress:\n%s\n\nCurrent file:\n%s",
			m.spinner.View(),
			m.progress.ViewAs(progressVal),
			lipgloss.NewStyle().Foreground(colorAccent).Render(m.files[m.currentFile]),
		)

		b.WriteString(header + "\n")
		b.WriteString(boxStyle.Render(box))

	case stateReview:
		file := m.files[m.currentFile]
		msg := m.commitMsgs[file]
		result := m.commitResults[file]

		icon := statusSuccess.Render("✓")
		if strings.Contains(result, "Error") {
			icon = statusError.Render("✗")
		}

		body := fmt.Sprintf(
			"%s\n%s\n\nFile: %s\n\nCommit message:\n%s\n\nStatus: %s\n\n[c] Confirm   [r] Retry   [q] Quit",
			asciiArt,
			titleStyle.Render("Repository: github.com/LunaSource/Luna"),
			lipgloss.NewStyle().Bold(true).Foreground(colorAccent).Render(file),
			msg,
			icon+" "+result,
		)

		b.WriteString(header + "\n")
		b.WriteString(boxStyle.Render(body))

	case stateComplete:
		b.WriteString(header + "\n")
		b.WriteString(statusSuccess.Render("All commits completed."))

	case stateError:
		b.WriteString(header + "\n")
		b.WriteString(statusError.Render(fmt.Sprintf("Error: %v", m.err)))
	}

	return b.String()
}

