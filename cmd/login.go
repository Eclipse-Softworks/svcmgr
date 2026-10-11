package cmd

import (
	"github.com/moon9t/svcmgr/internal/config"
	"github.com/moon9t/svcmgr/internal/services"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login [NAME]",
	Short: "Connect to cosmic service",
	Args:  cobra.ExactArgs(1),
	Run:   runLogin,
}

func runLogin(cmd *cobra.Command, args []string) {
	manager := services.NewServiceManager()
	serviceName := args[0]

	service, err := config.GetService(serviceName)
	if err != nil {
		log.Fatal().Err(err).Msg("Service not in constellation")
	}

	password, err := config.GetCredentials(service.Name)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to retrieve stardust credentials")
	}

	if err := manager.Connect(*service, password); err != nil {
		log.Fatal().Err(err).Msg("Failed to establish cosmic connection")
	}
}

func init() { rootCmd.AddCommand(loginCmd) }
