package runner

import (
	"fmt"
	"os"
	"os/exec"
)

type Runner struct {
	DryRun bool
}

func (r *Runner) Run(cmd []string) error {
	if r.DryRun {
		return nil
	}

	fmt.Printf("[Executing] %s\n", format(cmd))
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	return c.Run()
}

func format(cmd []string) string {
	var s string
	for i, v := range cmd {
		if i > 0 {
			s += " "
		}
		s += v
	}
	return s
}
