package termites_dbg

import (
	"fmt"
	"os/exec"
)

type CodeEditor func(file string, line int) error

var EditorGoland CodeEditor = func(file string, line int) error {
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

var EditorVSCode CodeEditor = func(file string, line int) error {
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
