package shell

import (
	"fmt"
	"github.com/dhuan/giback/pkg/utils"
	"log"
	"os"
	"os/exec"
)

func Run(dir string, command string, options RunOptions) ([]byte, error) {
	if options.Debug {
		dirDebug := dir

		if dirDebug == "" {
			pwd, _ := os.Getwd()

			dirDebug = pwd
		}

		log.Println(fmt.Sprintf("[%s] %s", dir, command))
	}

	commandName, commandParameters, err := parseCommand(command)

	if err != nil {
		return nil, err
	}

	cmd := exec.Command(commandName, commandParameters...)

	if dir != "" {
		cmd.Dir = dir
	}

	cmd.Env = os.Environ()

	if len(options.Env) > 0 {
		applyEnv(cmd, options.Env)
	}

	output, err := cmd.CombinedOutput()

	if err != nil {
		return output, err
	}

	return output, nil
}

type RunOptions struct {
	Debug bool
	Env   map[string]string
}

func RunOptionsDefault() RunOptions {
	return RunOptions{false, make(map[string]string)}
}

func RunMany(dir string, commands []string, output *[]byte, shellRunOptions RunOptions) error {
	if len(commands) == 0 {
		return nil
	}

	result, err := Run(dir, commands[0], shellRunOptions)

	if err != nil {
		return err
	}

	*output = append(*output, result...)

	return RunMany(dir, commands[1:], output, shellRunOptions)
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
	return fmt.Sprintf("%s=%s", key, value)
}
