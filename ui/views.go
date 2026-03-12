package ui

import (
	"time"

	"luna/git"

	tea "github.com/charmbracelet/bubbletea"
)


func loadFiles() tea.Msg {
	files, err := git.GetStagedFiles()
	return loadedFilesMsg{files: files, err: err}
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
