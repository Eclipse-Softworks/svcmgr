package services

import (
	"github.com/moon9t/svcmgr/internal/config"
	"github.com/moon9t/svcmgr/internal/utils"
	"github.com/rs/zerolog/log"
	// Remove the import for the connection tracker
)

func NewFileAuditLogger() AuditLogger {
	// Implement the function to return an instance of a type that satisfies the AuditLogger interface

	return &FileAuditLogger{}
}

type FileAuditLogger struct{}

func (f *FileAuditLogger) LogConnection(service config.Service) {
	// Implement the logging logic here
}

type ServiceManager struct {
	auditLog   AuditLogger
	connection ConnectionTracker
	// Remove the connection tracker field
}

type AuditLogger interface {
	LogConnection(service config.Service)
}

type ConnectionTracker interface {
	Track(service config.Service)
}

type Handler interface {
	Connect(service config.Service, password string) error
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		auditLog: NewFileAuditLogger(),
	}
}

func (sm *ServiceManager) Connect(service config.Service, password string) error {
	log.Info().
		Str("service", service.Name).
		Str("type", service.Type).
		Msg("Initiating cosmic connection")

	handler := getHandler(service.Type)
	if handler == nil {
		return utils.ErrUnsupportedService
	}

	sm.auditLog.LogConnection(service)
	return handler.Connect(service, password)
}

// Remove the NewConnectionTracker function

func getHandler(serviceType string) Handler {
	switch serviceType {
	case "postgres":
		return &PostgresHandler{}
	case "redis":
		return &RedisHandler{}
	case "http":
		return &HTTPHandler{}
	case "ssh":
		return &SSHHandler{}
	case "mysql":
		return &MySQLHandler{}
	// ... existing handlers ...
	default:
		return nil
	}
}

type PostgresHandler struct{}

func (p *PostgresHandler) Connect(service config.Service, password string) error {
	return p.connect(service, password)
}
