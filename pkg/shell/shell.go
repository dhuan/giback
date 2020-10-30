package shell

import (
	"fmt"
	"github.com/dhuan/giback/pkg/utils"
	"os"
	"os/exec"
)

func Run(dir string, command string, env map[string]string) ([]byte, error) {
	commandName, commandParameters, err := parseCommand(command)

	if err != nil {
		return nil, err
	}

	cmd := exec.Command(commandName, commandParameters...)

	if dir != "" {
		cmd.Dir = dir
	}

	cmd.Env = os.Environ()

	if len(env) > 0 {
		applyEnv(cmd, env)
	}

	output, err := cmd.CombinedOutput()

	if err != nil {
		return output, err
	}

	return output, nil
}

func parseCommand(command string) (string, []string, error) {
	commandArr, err := utils.SplitPreservingQuotes(command)

	if err != nil {
		return "", nil, err
	}

	commandName := commandArr[0]

	commandParameters := commandArr[1:]

	return commandName, commandParameters, nil
}

func applyEnv(cmd *exec.Cmd, env map[string]string) {
	for key, value := range env {
		cmd.Env = append(cmd.Env, toEnvVariable(key, value))
	}
}

func toEnvVariable(key string, value string) string {
	return fmt.Sprintf("%s=\"%s\"", key, value)
}
