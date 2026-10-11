package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/moon9t/svcmgr/internal/config"
	"github.com/moon9t/svcmgr/internal/ui"
	"github.com/moon9t/svcmgr/internal/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new cosmic service",
	Run:   runAdd,
}

func runAdd(cmd *cobra.Command, args []string) {
	qs := []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "Service name:",
				Help:    "Unique identifier for this service",
			},
			Validate: survey.Required,
		},
		{
			Name: "type",
			Prompt: &survey.Select{
				Message: "Service type:",
				Options: []string{"ssh", "mysql", "http", "postgres", "redis"},
				Default: "ssh",
			},
		},
		{
			Name: "host",
			Prompt: &survey.Input{
				Message: "Host:",
				Default: "localhost",
			},
			Validate: survey.Required,
		},
		{
			Name: "port",
			Prompt: &survey.Input{
				Message: "Port:",
				Default: "22",
			},
			Validate: survey.Required,
		},
		{
			Name: "username",
			Prompt: &survey.Input{
				Message: "Username:",
			},
		},
	}

	var answers struct {
		Name     string
		Type     string
		Host     string
		Port     string
		Username string
	}

	if err := survey.Ask(qs, &answers); err != nil {
		log.Fatal().Err(err).Msg("Failed to collect stardust")
	}

	service := config.Service{
		Name: answers.Name,
		Type: answers.Type,
		Host: answers.Host,
		Port: func() int {
			port, err := utils.ParsePort(answers.Port)
			if err != nil {
				log.Fatal().Err(err).Msg("Invalid port number")
			}
			return port
		}(),
		Username: answers.Username,
	}

	if err := service.Save(); err != nil {
		log.Fatal().Err(err).Msg("Failed to save cosmic configuration")
	}

	fmt.Printf("\n%s Service '%s' added to constellation!\n", ui.SuccessIcon(), service.Name)
}

func init() { rootCmd.AddCommand(addCmd) }
