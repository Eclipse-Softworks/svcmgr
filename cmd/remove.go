package cmd

import (
	"fmt"

	"github.com/moon9t/svcmgr/internal/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [NAME]",
	Short: "remove command",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		removed, err := config.RemoveService(args[0])
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to remove cosmic service")
		}
		if !removed {
			fmt.Printf("Service '%s' not found\n", args[0])
			return
		}
		fmt.Printf("Service '%s' removed\n", args[0])
	},
}

func init() { rootCmd.AddCommand(removeCmd) }
