package git

import (
	"bytes"
	"fmt"
	"luna/config"
	"os/exec"
	"path/filepath"
	"strings"
)

// runGit runs a git command and, on failure, returns git's own stderr
// message (e.g. "not a git repository", "dubious ownership") instead of
// just the bare exit status.
func runGit(args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return nil, fmt.Errorf("%s", msg)
	}
	return out, nil
}

// FileStatus is a file reported by `git status`, eligible to be staged.
type FileStatus struct {
	Path   string
	Status string // raw two-letter porcelain status, e.g. "M ", "??", " D"
}

// GetAvailableFiles lists every file git status reports as changed
// (staged, unstaged or untracked), deduplicated, in porcelain order.
func GetAvailableFiles() ([]FileStatus, error) {
	output, err := runGit("status", "--porcelain")
	if err != nil {
		return nil, fmt.Errorf("error running git status: %v", err)
	}

	var files []FileStatus
	seen := make(map[string]bool)

	lines := strings.Split(strings.TrimRight(string(output), "\n"), "\n")
	for _, line := range lines {
		if len(line) < 4 {
			continue
		}

		status := line[:2]
		path := strings.TrimSpace(line[3:])
		if idx := strings.Index(path, " -> "); idx != -1 {
			path = path[idx+4:]
		}
		path = strings.Trim(path, "\"")

		if path == "" || seen[path] {
			continue
		}
		seen[path] = true
		files = append(files, FileStatus{Path: path, Status: status})
	}

	return files, nil
}

// StageFile runs `git add` for a single file.
func StageFile(filename string) error {
	out, err := exec.Command("git", "add", "--", filename).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error staging %s: %s", filename, strings.TrimSpace(string(out)))
	}
	return nil
}

func GetStagedFiles() ([]string, error) {
	output, err := runGit("diff", "--cached", "--name-only")
	if err != nil {
		return nil, fmt.Errorf("error running git diff --cached: %v", err)
	}

	var files []string
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			files = append(files, line)
		}
	}

	return files, nil
}

func GetFileDiff(filename string) (string, error) {
	diff, err := runGit("diff", "--cached", "--", filename)
	if err != nil {
		return "", fmt.Errorf("error getting diff for %s: %v", filename, err)
	}
	return string(diff), nil
}

func ShouldIgnoreFile(filename string, cfg config.Config) bool {
	for _, ignoredFile := range cfg.IgnoredFiles {
		if filename == ignoredFile {
			return true
		}
		if strings.HasSuffix(ignoredFile, "/") && strings.HasPrefix(filename, ignoredFile) {
			return true
		}
	}

	for _, pattern := range cfg.IgnoredPatterns {
		matched, _ := filepath.Match(pattern, filepath.Base(filename))
		if matched {
			return true
		}
		fullPattern := filepath.FromSlash(pattern)
		matched, _ = filepath.Match(fullPattern, filename)
		if matched {
			return true
		}
	}

	return false
}
