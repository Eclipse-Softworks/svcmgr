package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

const configDirName = "svcmgr"

func getConfigPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to find cosmic config directory")
	}
	return filepath.Join(configDir, configDirName, "services.enc")
}

func (s Service) Save() error {
	services, err := LoadServices()
	if err != nil {
		return err
	}

	// Update existing or add new
	found := false
	for i, existing := range services {
		if existing.Name == s.Name {
			services[i] = s
			found = true
			break
		}
	}
	if !found {
		services = append(services, s)
	}

	return saveServices(services)
}

func LoadServices() ([]Service, error) {
	configPath := getConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return []Service{}, nil
	}

	encryptedData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	decryptedData, err := DecryptConfig(encryptedData)
	if err != nil {
		return nil, err
	}

	var services []Service
	if err := json.Unmarshal(decryptedData, &services); err != nil {
		return nil, err
	}

	return services, nil
}

func saveServices(services []Service) error {
	data, err := json.MarshalIndent(services, "", "  ")
	if err != nil {
		return err
	}

	encryptedData, err := EncryptConfig(data)
	if err != nil {
		return err
	}

	configPath := getConfigPath()
	if err := os.MkdirAll(filepath.Dir(configPath), 0700); err != nil {
		return err
	}

	return os.WriteFile(configPath, encryptedData, 0600)
}

// RemoveService deletes a service by name from the encrypted store.
func RemoveService(name string) (bool, error) {
	services, err := LoadServices()
	if err != nil {
		return false, err
	}
	for i := range services {
		if services[i].Name == name {
			services = append(services[:i], services[i+1:]...)
			return true, saveServices(services)
		}
	}
	return false, nil
}
