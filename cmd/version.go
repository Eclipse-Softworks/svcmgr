package cmd

import (
	"fmt"

	"github.com/moon9t/svcmgr/internal/ui"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version command",
	Run:   func(cmd *cobra.Command, args []string) { fmt.Print(ui.Banner()) },
}

func init() { rootCmd.AddCommand(versionCmd) }
