package termites_dbg

import (
	"fmt"
	"github.com/RoelofRuis/termites/pkg/termites"
	"os"
)

type ManagedTempDirectory struct {
	Dir string
}

func NewManagedTempDirectory(prefix string) *ManagedTempDirectory {
	dir, err := os.MkdirTemp("", prefix)
	if err != nil {
		panic(err)
	}

	return &ManagedTempDirectory{
		Dir: dir,
	}
}

func (m *ManagedTempDirectory) SetEventBus(b termites.EventBus) {
	b.Send(termites.Event{
		Type: termites.RegisterTeardown,
		Data: termites.RegisterTeardownEvent{Name: "temp dir", F: m.Teardown},
	})
}

func (m *ManagedTempDirectory) Teardown(control termites.TeardownControl) error {
	control.LogInfo(fmt.Sprintf("Cleaning up temp dir [%s]\n", m.Dir))
	return os.RemoveAll(m.Dir)
}
