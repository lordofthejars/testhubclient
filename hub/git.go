package hub

import (
	"bytes"
	"os/exec"
	"strings"
)

func getCurrentRevision() (string, error) {
	params := []string{"rev-parse"}

	params = append(params, "HEAD")
	return execCommand("git", params...)
}

func getCurrentBranch() (string, error) {
	return execCommand("git", "name-rev", "--name-only", "HEAD")
}

func execCommand(name string, arg ...string) (string, error) {
	var buff bytes.Buffer
	gitCmd := exec.Command(name, arg...)
	gitCmd.Stdout = &buff
	err := gitCmd.Run()
	return strings.TrimSpace(buff.String()), err
}
