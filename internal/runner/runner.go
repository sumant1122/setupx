package runner

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

type Runner struct {
	DryRun   bool
	MaxLines int
}

func (r *Runner) Run(cmd []string) error {
	if r.DryRun {
		return nil
	}

	fmt.Printf("[Executing] %s\n", format(cmd))
	c := exec.Command(cmd[0], cmd[1:]...)
	
	if r.MaxLines > 0 {
		return r.runLimited(c)
	}

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	return c.Run()
}

func (r *Runner) RunOutput(cmd []string) (string, error) {
	if r.DryRun {
		return "", nil
	}

	c := exec.Command(cmd[0], cmd[1:]...)
	out, err := c.CombinedOutput()
	return string(out), err
}

func (r *Runner) Check(cmd []string) bool {
	c := exec.Command(cmd[0], cmd[1:]...)
	err := c.Run()
	return err == nil
}

func (r *Runner) runLimited(c *exec.Cmd) error {
	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}
	c.Stderr = os.Stderr

	if err := c.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	lineCount := 0
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		lineCount++
		if lineCount >= r.MaxLines {
			fmt.Printf("\n[Truncated to top %d results]\n", r.MaxLines)
			break
		}
	}

	// We don't necessarily want to wait for the command if we truncated it,
	// but for some package managers it might be better to signal to stop.
	// However, for simplicity, we'll just return. 
	// Note: The process might still be running in background until it finishes or c.Wait is called.
	return nil 
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
