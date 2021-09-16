package termites

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type closeOnTeardown struct {
	name   string
	closer io.Closer
	bus    EventBus
}

func (c closeOnTeardown) SetEventBus(b EventBus) {
	b.Send(Event{
		Type: RegisterTeardown,
		Data: RegisterTeardownEvent{Name: c.name, F: c.Teardown},
	})
	c.bus = b
}

func (c closeOnTeardown) Teardown(control TeardownControl) error {
	if err := c.closer.Close(); err != nil {
		control.LogError(fmt.Sprintf("error closing resource [%s]", c.name), err)
		return err
	}
	return nil
}

type ManagedTempDirectory struct {
	Dir string
}

func NewManagedTempDirectory(prefix string) *ManagedTempDirectory {
	dir, err := ioutil.TempDir("", prefix)
	if err != nil {
		panic(err)
	}

	return &ManagedTempDirectory{
		Dir: dir,
	}
}

func (m *ManagedTempDirectory) SetEventBus(b EventBus) {
	b.Send(Event{
		Type: RegisterTeardown,
		Data: RegisterTeardownEvent{Name: "temp dir", F: m.Teardown},
	})
}

func (m *ManagedTempDirectory) Teardown(control TeardownControl) error {
	control.LogInfo(fmt.Sprintf("Cleaning up temp dir [%s]\n", m.Dir))
	return os.RemoveAll(m.Dir)
}
