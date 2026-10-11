package cmd

import (
	"context"
	"fmt"
	"time"

	"encoding/json"

	"github.com/AlecAivazis/survey/v2"
	"github.com/moon9t/svcmgr/internal/config"
	"github.com/moon9t/svcmgr/internal/storage"
	"github.com/moon9t/svcmgr/internal/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	forceBackup  bool
	backupBucket string
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Secure cosmic configuration backup",
	Long:  "Store celestial configuration in the cosmic vault",
	Run:   runBackup,
}

func init() {
	backupCmd.Flags().BoolVarP(&forceBackup, "force", "f", false,
		"Bypass cosmic confirmation")
	backupCmd.Flags().StringVarP(&backupBucket, "bucket", "b", "",
		"Celestial storage bucket")
	rootCmd.AddCommand(backupCmd)
}

func runBackup(cmd *cobra.Command, args []string) {
	if backupBucket == "" {
		log.Fatal().Msg("Pass a bucket with --bucket")
	}
	if !forceBackup {
		confirm := false
		prompt := &survey.Confirm{
			Message: "Initiate cosmic backup sequence?",
		}
		survey.AskOne(prompt, &confirm)
		if !confirm {
			fmt.Println("Backup aborted by cosmic operator")
			return
		}
	}

	vault, err := storage.NewCosmicVault(backupBucket)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to access cosmic vault")
	}

	services, err := config.LoadServices()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read cosmic configuration")
	}

	servicesData, err := json.Marshal(services)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to serialize cosmic configuration")
	}

	encryptedData, err := config.EncryptConfig(servicesData)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to encrypt stardust data")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := vault.BackupConfig(ctx, encryptedData); err != nil {
		log.Fatal().Err(err).Msg("Cosmic backup failed")
	}

	fmt.Println(utils.CosmicSuccess)
	log.Info().Msg("Cosmic configuration secured in the vault")
}
