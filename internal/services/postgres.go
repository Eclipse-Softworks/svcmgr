package services

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/moon9t/svcmgr/internal/config"
	"github.com/rs/zerolog/log"
)

func (h *PostgresHandler) connect(s config.Service, pwd string) error {
	var extra map[string]interface{}
	if err := json.Unmarshal([]byte(s.Extra), &extra); err != nil {
		return fmt.Errorf("failed to parse extra: %v", err)
	}
	database, ok := extra["database"].(string)
	if !ok {
		return fmt.Errorf("database key is not a string")
	}
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		s.Username, pwd, s.Host, s.Port, database)

	cmd := exec.Command("psql", connString)
	log.Debug().Str("command", cmd.String()).Msg("Launching PSQL")

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
