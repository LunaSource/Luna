
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"luna/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage project configuration.",
	Aliases: []string{"cfg"},
	Args: cobra.OnlyValidArgs,
	ValidArgs: []string {"init", "show", "edit"},
	
	Run: func(cmd *cobra.Command, args []string) {
		
		config.ManageConfig(args[0])
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.SetUsageFunc(func(cmd *cobra.Command) error {
	fmt.Println("Usage:")
	fmt.Println("  luna config [init|show|edit]")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  init   Create the configuration file")
	fmt.Println("  show   Display the current configuration")
	fmt.Println("  edit   Edit the project configuration")
	fmt.Println()
	fmt.Println("Config Priority:")
	fmt.Println("  • API Key: Global > Project > Default")
	fmt.Println("  • Other settings: Project > Default")
	fmt.Println()
	// fmt.Println("Flags:")
	// cmd.Flags().PrintDefaults()
	return nil
})
}
