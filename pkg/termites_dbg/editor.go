package termites_dbg

import (
	"fmt"
	"os/exec"
)

type CodeEditor interface {
	Open(file string, line int) error
}

type Goland struct{}

func (e Goland) Open(file string, line int) error {
	err := exec.Command(
		"goland",
		"--line",
		fmt.Sprintf("%d", line),
		file,
	).Start()
	if err != nil {
		return err
	}
	return nil
}

type VisualStudioCode struct{}

func (e VisualStudioCode) Open(file string, line int) error {
	err := exec.Command(
		"code",
		"--goto",
		fmt.Sprintf("%s:%d:0", file, line),
	).Start()
	if err != nil {
		return err
	}
	return nil
}
