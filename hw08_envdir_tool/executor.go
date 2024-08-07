package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if err := updateEnv(env); err != nil {
		return ReturnCodeErr
	}
	if err := executeCommand(cmd); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return ReturnCodeErr
	}
	return ReturnCodeOk
}

func updateEnv(env Environment) error {
	for k, v := range env {
		// Независимо от назначения элемента из списка его нужно удалить.
		// Часть элементов после удаления нужно присвоить заново.
		// Удаление происходит за 1 операцию, без предварительной проверки.
		if err := os.Unsetenv(k); err != nil {
			return err
		}
		if v.NeedSetToEnv() {
			if err := os.Setenv(k, v.Value); err != nil {
				return err
			}
		}
	}
	return nil
}

func executeCommand(cmd []string) error {
	var c *exec.Cmd
	if len(cmd) > 1 {
		c = exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	} else {
		c = exec.Command(cmd[0]) //nolint:gosec
	}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	if err := c.Run(); err != nil {
		return err
	}
	return nil
}
