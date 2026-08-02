package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"luna/config"
	"luna/style"
	"luna/ui"
)

var includeEmoji bool
var cfg = config.LoadConfig()

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generates commit messages using OpenRouter AI.",
	Aliases: []string{"c"},
	Run: func(cmd *cobra.Command, args []string) {

		if !includeEmoji {
			cfg = config.LoadConfig()

		}

		if cfg.ApiKey == "" {
			fmt.Println(style.Err("Set your API key first: luna apikey <KEY>"))
			return
		}

		if cfg.Model == "" {
			fmt.Println(style.Err("Set your model first: luna model <MODEL_ID>"))
			return
		}

		var p *tea.Program

		p = tea.NewProgram(ui.InitializeCommitUI(cfg, includeEmoji))

		if _, err := p.Run(); err != nil {
			fmt.Println(style.Err(fmt.Sprintf("Error running UI: %v", err)))
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)

	commitCmd.Flags().BoolVarP(&includeEmoji, "emoji", "e", false, "Include emoji to the commit message.")
	
	
}
