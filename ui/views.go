package ui

import (
	"time"

	"luna/config"
	"luna/git"

	tea "github.com/charmbracelet/bubbletea"
)

func loadAvailableFiles(cfg config.Config) tea.Cmd {
	return func() tea.Msg {
		files, err := git.GetAvailableFiles()
		if err != nil {
			return availableFilesMsg{err: err}
		}

		var filtered []git.FileStatus
		for _, f := range files {
			if !git.ShouldIgnoreFile(f.Path, cfg) {
				filtered = append(filtered, f)
			}
		}

		return availableFilesMsg{files: filtered}
	}
}

func handleSelectionInput(m CommitUI, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {

	case "up":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down":
		if m.cursor < len(m.available)-1 {
			m.cursor++
		}

	case " ":
		path := m.available[m.cursor].Path
		m.selected[path] = !m.selected[path]

	case "a":
		selectedCount := 0
		for _, f := range m.available {
			if m.selected[f.Path] {
				selectedCount++
			}
		}
		selectAll := selectedCount != len(m.available)
		for _, f := range m.available {
			m.selected[f.Path] = selectAll
		}

	case "enter":
		var chosen []string
		for _, f := range m.available {
			if m.selected[f.Path] {
				chosen = append(chosen, f.Path)
			}
		}
		if len(chosen) == 0 {
			return m, nil
		}
		m.files = chosen
		m.currentFile = 0
		m.state = stateStaging
		return m, stageFiles(chosen)

	case "q":
		return m, tea.Quit
	}

	return m, nil
}

func stageFiles(files []string) tea.Cmd {
	return func() tea.Msg {
		for _, f := range files {
			if err := git.StageFile(f); err != nil {
				return stagedMsg{err: err}
			}
		}
		return stagedMsg{}
	}
}

func handleReviewInput(m CommitUI, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {

	case "c":
		file := m.files[m.currentFile]
		commitMsg := m.commitMsgs[file]

		out, err := git.CommitFile(file, commitMsg)
		if err != nil {
			m.state = stateError
			m.err = err
			return m, nil
		}

		m.commitResults[file] = out
		m.currentFile++

		if m.currentFile >= len(m.files) {
			m.state = stateComplete
			return m, func() tea.Msg {
				time.Sleep(2 * time.Second)
				return tea.QuitMsg{}
			}
		}

		m.state = stateProcessing
		return m, processNextFile(m)

	case "r":
		file := m.files[m.currentFile]
		diff, err := git.GetFileDiff(file)
		if err != nil {
			m.state = stateError
			m.err = err
			return m, nil
		}

		newMsg := git.GenerateCommitMessage(
			m.cfg.ApiKey,
			diff,
			file,
			m.cfg,
			m.includeEmoji,
		)

		m.commitMsgs[file] = newMsg
		return m, nil

	case "q":
		return m, tea.Quit
	}

	return m, nil
}
