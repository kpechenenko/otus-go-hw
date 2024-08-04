package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for k, v := range env {
		// Нужно удалить каждый элемент из списка. Часть элементов после удаления нужно присвоить заново
		if err := os.Unsetenv(k); err != nil {
			return 1
		}
		// у переменных, которые нужно присвоить, флаг в читающей функции записан в false
		// при этом по v.Value != "" сравнивать нельзя, т.к. в value может быть пустая строка
		if !v.NeedRemove {
			if err := os.Setenv(k, v.Value); err != nil {
				return 1
			}
		}
	}
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
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return 1
	}
	return 0
}
