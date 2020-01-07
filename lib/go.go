package lib

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type BuildOption struct {
	Option
	OutputPath string
	BuildPath  string
	OS         string
	Arch       string
}

type Option struct {
	Dir string
	Env []string
}

func GoGet(getPath string) (string, error) {
	cmdArgs := []string{
		"get",
		"-d",
		getPath,
	}
	fmt.Println("go", strings.Join(cmdArgs, " "))
	stdout, err := execCommand("go", cmdArgs, &Option{})
	if err != nil {
		return "", err
	}
	return stdout, nil
}

func GoBuild(opt *BuildOption) (string, error) {
	var stderr, stdout bytes.Buffer
	cmdArgs := []string{
		"build",
		"-o", opt.OutputPath,
		opt.BuildPath,
	}
	fmt.Println("go", strings.Join(cmdArgs, " "))
	cmd := exec.Command("go", cmdArgs...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if opt.Env != nil {
		cmd.Env = opt.Env
	}

	if opt.OS != "" {
		cmd.Env = append(cmd.Env, "GOOS="+opt.OS)
	}

	if opt.Arch != "" {
		cmd.Env = append(cmd.Env, "GOARCH="+opt.Arch)
	}

	if opt.Dir != "" {
		cmd.Dir = opt.Dir
	}
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s: %w", stderr.String(), err)
	}
	return stdout.String(), nil
}

func execCommand(name string, args []string, opt *Option) (string, error) {
	var stderr, stdout bytes.Buffer
	cmd := exec.Command("go", args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if opt.Env != nil {
		cmd.Env = opt.Env
	}

	if opt.Dir != "" {
		cmd.Dir = opt.Dir
	}
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s: %w", stderr.String(), err)
	}
	return stdout.String(), nil
}
