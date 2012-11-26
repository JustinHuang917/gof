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
	out, err := run(goFilePath, buidingArgs...)
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

func run(path string, buidingArgs ...string) ([]byte, error) {
	var buf bytes.Buffer
	args := make([]string, 2, 10)
	args[0] = "build"
	args[1] = path
	args = append(args, buidingArgs...)
	//cmd := exec.Command("go", "build", path, buidingArgs...)
	cmd := exec.Command("go", args...)
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
