package termites_dbg

import (
	"fmt"
	"os/exec"
)

func Open(file string, line int) error {
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
