/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"luna/config"
	"luna/style"
)

// modelCmd represents the model command
var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "Register the OpenRouter model used to generate commit messages.",
	Aliases: []string{"m"},
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		var modelStr string = args[0]

		var err error = config.SaveGlobalModel(modelStr)

		if err != nil {
			fmt.Fprintln(os.Stderr, style.Err(fmt.Sprintf("Error saving model: %v", err)))
			return
		}
		fmt.Println(style.Success("Model saved to global configuration " + style.IconChip))
		fmt.Println(style.Muted(style.IconMarker + "  ~/.lunarc"))

	},
}

func init() {
	rootCmd.AddCommand(modelCmd)
}
