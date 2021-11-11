package termites_dbg

import (
	"fmt"
	"github.com/RoelofRuis/termites/pkg/termites"
	"os/exec"
)

func open(f termites.FunctionInfo, editor CodeEditor) error {
	if f.File == "" {
		return nil
	}
	return editor(f.File, f.Line)
}

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
