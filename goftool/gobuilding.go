package goftool

import (
	"bytes"
	"fmt"
	"os/exec"
	"sync"
)

type BuildResult struct {
	Err error
	Out []byte
}

func Building(goFilePath string, buidingArgs ...string) *BuildResult {
	cmd1 := gobuild(goFilePath, buidingArgs...)
	cmd2 := gofmt(goFilePath)
	out, err := run(cmd1)
	if err == nil {
		out, err = run(cmd2)
	}
	return &BuildResult{err, out}
}

var running struct {
	sync.Mutex
	cmd *exec.Cmd
}

func stopRun() {
	running.Lock()
	if running.cmd != nil {
		running.cmd.Process.Kill()
		running.cmd = nil
	}
	running.Unlock()
}

func kill() {
	stopRun()
}

func gofmt(path string) *exec.Cmd {
	args := make([]string, 2, 10)
	args[0] = "-w"
	args[1] = path
	cmd := exec.Command("gofmt", args...)
	return cmd
}

func gobuild(path string, buidingArgs ...string) *exec.Cmd {
	args := make([]string, 2, 10)
	args[0] = "build"
	args[1] = path
	args = append(args, buidingArgs...)
	cmd := exec.Command("go", args...)
	return cmd
}

func run(cmd *exec.Cmd) ([]byte, error) {
	var buf bytes.Buffer
	// args := make([]string, 2, 10)
	// args[0] = "build"
	// args[1] = path
	// args = append(args, buidingArgs...)
	// //cmd := exec.Command("go", "build", path, buidingArgs...)
	// cmd := exec.Command("go", args...)
	cmd.Stdout = &buf
	cmd.Stderr = cmd.Stdout

	// Start command and leave in 'running'.
	running.Lock()
	if running.cmd != nil {
		defer running.Unlock()
		return nil, fmt.Errorf("already running %s", running.cmd.Path)
	}
	if err := cmd.Start(); err != nil {
		running.Unlock()
		return nil, err
	}
	running.cmd = cmd
	running.Unlock()

	// Wait for the command.  Clean up,
	err := cmd.Wait()
	running.Lock()
	if running.cmd == cmd {
		running.cmd = nil
	}
	running.Unlock()
	return buf.Bytes(), err
}
