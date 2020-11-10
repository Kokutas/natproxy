package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	version = "v1.0"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "NPC version",
	Long:  `NPC version.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("NPC Version: %v.\n", version)
	},
}
