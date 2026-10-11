package services

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"

	"github.com/moon9t/svcmgr/internal/config"
	"github.com/rs/zerolog/log"
	"golang.org/x/term"
)

type SSHHandler struct{}
type MySQLHandler struct{}

func (h *SSHHandler) Connect(s config.Service, pwd string) error {
	var cmd *exec.Cmd
	args := []string{
		fmt.Sprintf("%s@%s", s.Username, s.Host),
		"-p", fmt.Sprintf("%d", s.Port),
	}

	if term.IsTerminal(int(syscall.Stdin)) {
		args = append(args, "-t")
	}

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("putty", append([]string{"-ssh"}, args...)...)
	default:
		cmd = exec.Command("ssh", args...)
	}

	log.Debug().Str("command", cmd.String()).Msg("Executing SSH command")
	return executeCommand(cmd, pwd)
}

func (h *MySQLHandler) Connect(s config.Service, pwd string) error {
	cmd := exec.Command("mysql",
		"-h", s.Host,
		"-u", s.Username,
		"-P", fmt.Sprintf("%d", s.Port),
	)

	log.Debug().Str("command", cmd.String()).Msg("Executing MySQL command")
	if pwd != "" {
		cmd.Env = append(os.Environ(), "MYSQL_PWD="+pwd)
	}
	return executeCommand(cmd, "")
}

// Universal command executor with password injection
func executeCommand(cmd *exec.Cmd, password string) error {
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		defer stdin.Close()
		if password != "" {
			stdin.Write([]byte(password + "\n"))
		}
	}()

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
