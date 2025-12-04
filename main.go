package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"luna/config"
	"luna/ui"
)

func main() {
	if len(os.Args) < 2 {
		showUsage()
		return
	}

	var cmd string
	includeEmoji := false

	for _, arg := range os.Args[1:] {
		larg := strings.ToLower(arg)
		switch larg {
		case "-e":
			includeEmoji = true
		case "commit", "c":
			cmd = "commit"
		case "help", "h":
			cmd = "help"
		case "apikey", "k":
			cmd = "apikey"
		case "config", "cfg":
			cmd = "config"
		}
	}

	switch cmd {
	case "help":
		ui.ShowHelp()
	case "commit":
		runCommitGenerator(includeEmoji)
	case "apikey":
		setApiKey()
	case "config":
		manageConfig()
	default:
		fmt.Println("Unknown command. Use: luna help")
	}
}

func runCommitGenerator(includeEmoji bool) {
	cfg := config.LoadConfig()

	if !includeEmoji {
		includeEmoji = cfg.DefaultEmoji
	}

	if cfg.ApiKey == "" {
		fmt.Println("Error: Set API key using 'luna apikey' first")
		return
	}

	p := tea.NewProgram(ui.InitializeCommitUI(cfg, includeEmoji))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running UI: %v\n", err)
		os.Exit(1)
	}
}

func setApiKey() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: luna apikey YOUR_API_KEY")
		return
	}

	apiKey := os.Args[2]
	err := config.SaveGlobalApiKey(apiKey)
	if err != nil {
		fmt.Printf("Error saving API key: %v\n", err)
		return
	}
	fmt.Println("✅ API key saved successfully in global configuration!")
	fmt.Println("📍 Location: ~/.lunarc")
}

func manageConfig() {
	if len(os.Args) < 3 {
		ui.ShowConfigHelp()
		return
	}

	subcmd := os.Args[2]
	config.ManageConfig(subcmd)
}

func showUsage() {
	fmt.Println("Use: luna help to see available commands")
	fmt.Println("Available commands with aliases:")
	fmt.Println("help (h)")
	fmt.Println("commit (c)")
	fmt.Println("apikey (k)")
	fmt.Println("config (cfg)")
}
