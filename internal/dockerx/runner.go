package dockerx

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var ErrDockerUnavailable = errors.New("docker command unavailable")

type Mount struct {
	Source   string
	Target   string
	ReadOnly bool
}

type Command struct {
	Image   string
	Args    []string
	Mounts  []Mount
	Workdir string
}

type Result struct {
	ExitCode int
	Stdout   string
	Stderr   string
}

type Runner interface {
	Run(ctx context.Context, command Command) (Result, error)
}

type CLIRunner struct{}

func (CLIRunner) Run(ctx context.Context, command Command) (Result, error) {
	if _, err := exec.LookPath("docker"); err != nil {
		return Result{}, fmt.Errorf("%w: %v", ErrDockerUnavailable, err)
	}
	args := []string{"run", "--rm"}
	for _, mount := range command.Mounts {
		args = append(args, "--mount", mountSpec(mount))
	}
	if command.Workdir != "" {
		args = append(args, "--workdir", command.Workdir)
	}
	args = append(args, command.Image)
	args = append(args, command.Args...)

	cmd := exec.CommandContext(ctx, "docker", args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	result := Result{
		ExitCode: exitCode(err),
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
	}
	if err != nil {
		return result, fmt.Errorf("docker run failed with exit code %d", result.ExitCode)
	}
	return result, nil
}

func IsDockerUnavailable(err error) bool {
	return errors.Is(err, ErrDockerUnavailable)
}

func mountSpec(mount Mount) string {
	parts := []string{
		"type=bind",
		"source=" + mount.Source,
		"target=" + mount.Target,
	}
	if mount.ReadOnly {
		parts = append(parts, "readonly")
	}
	return strings.Join(parts, ",")
}

func exitCode(err error) int {
	if err == nil {
		return 0
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}
	return 1
}
