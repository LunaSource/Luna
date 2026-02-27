package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"luna/config"
	"luna/ui"
)

var includeEmoji bool
var cfg = config.LoadConfig()

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generates commit messages using Gemini AI.",
	Aliases: []string{"c"},
	Run: func(cmd *cobra.Command, args []string) {

		if !includeEmoji {
			cfg = config.LoadConfig()

		}

		if cfg.ApiKey == "" {
			fmt.Println("Error: Set API key using 'luna apikey' first")
			return
		}

		var p *tea.Program

		p = tea.NewProgram(ui.InitializeCommitUI(cfg, includeEmoji))

		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running UI: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)

	commitCmd.Flags().BoolVarP(&includeEmoji, "emoji", "e", false, "Include emoji to the commit message.")
	
	
}
