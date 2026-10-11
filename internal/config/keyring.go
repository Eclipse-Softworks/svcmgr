package config

import (
	"github.com/rs/zerolog/log"
	"github.com/zalando/go-keyring"
)

func GetService(serviceName string) (*Service, error) {
	return LoadService(serviceName)
}

const (
	keyringService = "svcmgr-cosmic-vault"
)

func StoreCredentials(serviceName, username, password string) error {
	err := keyring.Set(keyringService, serviceName+"-"+username, password)
	if err != nil {
		log.Error().Err(err).
			Str("service", serviceName).
			Msg("Failed to store cosmic credentials")
	}
	return err
}

func GetCredentials(serviceName string) (string, error) {
	service, err := GetService(serviceName)
	if err != nil {
		return "", err
	}

	password, err := keyring.Get(keyringService, serviceName+"-"+service.Username)
	if err != nil {
		log.Error().Err(err).
			Str("service", serviceName).
			Msg("Failed to retrieve cosmic credentials")
	}
	return password, err
}

func DeleteCredentials(serviceName string) error {
	service, err := GetService(serviceName)
	if err != nil {
		return err
	}

	return keyring.Delete(keyringService, serviceName+"-"+service.Username)
}
