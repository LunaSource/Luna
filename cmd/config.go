
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"luna/config"
	"luna/style"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage project configuration.",
	Aliases: []string{"cfg"},
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string {"init", "show", "edit"},
	
	Run: func(cmd *cobra.Command, args []string) {
		
		config.ManageConfig(args[0])
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.SetUsageFunc(func(cmd *cobra.Command) error {
	fmt.Println(style.Accent("Usage:"))
	fmt.Println("  luna config [init|show|edit]")
	fmt.Println()
	fmt.Println(style.Accent("Arguments:"))
	fmt.Println("  init   Create the configuration file")
	fmt.Println("  show   Display the current configuration")
	fmt.Println("  edit   Edit the project configuration")
	fmt.Println()
	fmt.Println(style.Accent("Config Priority:"))
	fmt.Println("  " + style.IconAngleRight + "  API Key & Model: Global > Project > Default")
	fmt.Println("  " + style.IconAngleRight + "  Other settings: Project > Default")
	fmt.Println()
	// fmt.Println("Flags:")
	// cmd.Flags().PrintDefaults()
	return nil
})
}
