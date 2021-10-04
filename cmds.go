package cmd

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func output(args ...string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	stderr := bytes.Buffer{}
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return "", errors.New(stderr.String())
	}
	return string(out), err
}

func currentContext() (string, error) {
	out, err := output("kubectl", "config", "current-context")
	return strings.TrimSpace(out), err
}

func connect() error {
	out, err := output("telepresence", "connect")
	if err != nil {
		return err
	}
	if !strings.Contains(out, "Connected to context") {
		return errors.New(out)
	}
	return nil
}

func uninstall() error {
	out, err := output("telepresence", "uninstall", "-e")
	if err != nil {
		return err
	}
	if !strings.Contains(out, "Telepresence Root Daemon quitting... done") {
		return errors.New(out)
	}
	return nil
}
