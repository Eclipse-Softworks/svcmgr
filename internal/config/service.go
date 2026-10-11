package config

import (
	"fmt"
)

// Service represents a service configuration.
type Service struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Category string `json:"category,omitempty"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Extra    string `json:"extra"` // Add the Extra field
}

// LoadService looks up a service by name in the encrypted store.
func LoadService(name string) (*Service, error) {
	services, err := LoadServices()
	if err != nil {
		return nil, err
	}
	for i := range services {
		if services[i].Name == name {
			return &services[i], nil
		}
	}
	return nil, fmt.Errorf("service %q not found", name)
}
