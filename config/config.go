package config

import (
	"bufio"
	"encoding/json"
	"fmt"

	// "io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"luna/style"
)

const (
	GLOBAL_CONFIG_FILE  = ".lunarc"
	PROJECT_CONFIG_FILE = ".lunacfg"
)

type Config struct {
	IgnoredFiles    []string `json:"ignoredFiles"`
	IgnoredPatterns []string `json:"ignoredPatterns"`
	CommitPrefixes  []string `json:"commitPrefixes"`
	MaxCommitLength int      `json:"maxCommitLength"`
	DefaultEmoji    bool     `json:"defaultEmoji"`
	ApiKey          string   `json:"apiKey"`
	Model           string   `json:"model"`
}

type configItem struct {
	title string
	desc  string
	value string
}

func (i configItem) Title() string       { return i.title }
func (i configItem) Description() string { return i.desc }
func (i configItem) FilterValue() string { return i.value }

func GetDefaultConfig() Config {
	return Config{
		IgnoredFiles:    []string{},
		IgnoredPatterns: []string{"*.exe", "*.dll", "*.png", "*.jpg", "*.jpeg", "*.gif", "*.bin"},
		CommitPrefixes:  []string{"chore:", "refactor:", "feat:", "fix:", "docs:", "test:"},
		MaxCommitLength: 72,
		DefaultEmoji:    false,
		ApiKey:          "",
		Model:           "",
	}
}

func LoadConfig() Config {
	projectCfg := LoadProjectConfig()
	globalCfg := LoadGlobalConfig()

	if projectCfg.ApiKey == "" && globalCfg.ApiKey != "" {
		projectCfg.ApiKey = globalCfg.ApiKey
	}

	if projectCfg.Model == "" && globalCfg.Model != "" {
		projectCfg.Model = globalCfg.Model
	}

	return projectCfg
}

func SaveGlobalApiKey(apiKey string) error {
	globalCfg := LoadGlobalConfig()
	globalCfg.ApiKey = apiKey
	return SaveGlobalConfig(globalCfg)
}

func SaveGlobalModel(model string) error {
	globalCfg := LoadGlobalConfig()
	globalCfg.Model = model
	return SaveGlobalConfig(globalCfg)
}

func ManageConfig(subcmd string) {
	switch subcmd {
	case "init":
		initConfigInteractive()
	case "show":
		showConfig()
	case "edit":
		editConfig()
	default:
		fmt.Println(style.Err("Unknown config command"))
	}
}

func initConfigInteractive() {
	items := []list.Item{
		configItem{title: ".lunacfg", desc: "Project configuration (current directory)", value: "project"},
		configItem{title: ".lunarc", desc: "Global configuration (home directory)", value: "global"},
		configItem{title: "Both files", desc: "Create both configurations", value: "both"},
		configItem{title: "Exit", desc: "Exit without creating any files", value: "exit"},
	}

	l := list.New(items, list.NewDefaultDelegate(), 60, 14)
	l.Title = style.IconFolderOpen + "  Select configuration files to create"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	m := model{list: l}
	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Println(style.Err(fmt.Sprintf("Error running config menu: %v", err)))
		return
	}

	m = finalModel.(model)

	if contains(m.selected, "exit") || len(m.selected) == 0 {
		fmt.Println(style.Muted(style.IconSignOut + "  Exit without creating any files"))
		return
	}

	if contains(m.selected, "both") {
		m.selected = []string{"project", "global"}
	}

	createSelectedConfigs(m.selected)

	for _, sel := range m.selected {
		switch sel {
		case "project":
			apiKey := askForAPIKey("project")
			model := askForModel("project")
			cfg := LoadProjectConfig()
			if apiKey != "" {
				cfg.ApiKey = apiKey
			}
			if model != "" {
				cfg.Model = model
			}
			SaveProjectConfig(cfg)
		case "global":
			apiKey := askForAPIKey("global")
			model := askForModel("global")
			cfg := LoadGlobalConfig()
			if apiKey != "" {
				cfg.ApiKey = apiKey
			}
			if model != "" {
				cfg.Model = model
			}
			SaveGlobalConfig(cfg)
		}
	}
}

type model struct {
	list     list.Model
	selected []string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.String(); key {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case " ":
			selectedItem := m.list.SelectedItem().(configItem)
			if contains(m.selected, selectedItem.value) {
				m.selected = remove(m.selected, selectedItem.value)
			} else {
				m.selected = append(m.selected, selectedItem.value)
			}
		case "enter":
			m.quitting = true
			return m, tea.Quit
		case "a":
			m.selected = []string{"project", "global"}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var view strings.Builder
	view.WriteString(style.Accent(m.list.Title) + "\n\n")

	if len(m.selected) > 0 {
		view.WriteString(style.Accent("Selected: "))
		for i, sel := range m.selected {
			if i > 0 {
				view.WriteString(", ")
			}
			switch sel {
			case "project":
				view.WriteString(".lunacfg")
			case "global":
				view.WriteString(".lunarc")
			}
		}
		view.WriteString("\n\n")
	}

	for i, item := range m.list.Items() {
		cfgItem := item.(configItem)
		if i == m.list.Index() {
			view.WriteString(style.Accent(style.IconAngleRight + " "))
		} else {
			view.WriteString("  ")
		}

		if contains(m.selected, cfgItem.value) ||
			(cfgItem.value == "both" && len(m.selected) == 2) ||
			(cfgItem.value == "exit" && len(m.selected) == 0) {
			view.WriteString(style.Accent("(" + style.IconStar + ") "))
		} else {
			view.WriteString(style.Muted("(" + style.IconStarEmpty + ") "))
		}

		view.WriteString(cfgItem.title + " - " + cfgItem.desc + "\n")
	}

	help := fmt.Sprintf(
		"%s%s Navigate    Space Select    Enter Confirm    A Select All    Q Quit",
		style.IconArrowUp, style.IconArrowDown,
	)
	view.WriteString("\n" + style.Muted(help) + "\n")
	return view.String()
}

func createSelectedConfigs(selected []string) {
	if contains(selected, "project") {
		createProjectConfig()
	}
	if contains(selected, "global") {
		createGlobalConfig()
	}
}

func createProjectConfig() {
	path := filepath.Join(".", PROJECT_CONFIG_FILE)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		cfg := GetDefaultConfig()
		if err := SaveProjectConfig(cfg); err != nil {
			fmt.Println(style.Err(fmt.Sprintf("Error creating .lunacfg: %v", err)))
			return
		}
		fmt.Println(style.Success("Created .lunacfg in current directory"))
	} else {
		fmt.Println(style.Info(".lunacfg already exists in current directory"))
	}
}

func createGlobalConfig() {
	path, _ := getGlobalConfigPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		cfg := GetDefaultConfig()
		if err := SaveGlobalConfig(cfg); err != nil {
			fmt.Println(style.Err(fmt.Sprintf("Error creating .lunarc: %v", err)))
			return
		}
		fmt.Println(style.Success("Created .lunarc in home directory"))
	} else {
		fmt.Println(style.Info(".lunarc already exists in home directory"))
	}
}

func askForAPIKey(configType string) string {
	fmt.Println(style.Accent(fmt.Sprintf("%s  Set API key for %s? (Enter to skip or paste)", style.IconAngleRight, configType)))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		clip, _ := clipboard.ReadAll()
		return strings.TrimSpace(clip)
	}
	return input
}

func askForModel(configType string) string {
	fmt.Println(style.Accent(fmt.Sprintf("%s  Set OpenRouter model for %s? (e.g. openai/gpt-4o-mini, Enter to skip)", style.IconAngleRight, configType)))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func remove(slice []string, item string) []string {
	var out []string
	for _, s := range slice {
		if s != item {
			out = append(out, s)
		}
	}
	return out
}

func showConfig() {
	cfg := LoadConfig()
	data, _ := json.MarshalIndent(cfg, "", "  ")
	fmt.Println(string(data))
}

func editConfig() {
	fmt.Println(style.Warn("Edit is not implemented yet"))
}
