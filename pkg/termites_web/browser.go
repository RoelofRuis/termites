package termites_web

import (
	"fmt"
	"github.com/RoelofRuis/termites/pkg/termites"
	"os/exec"
	"runtime"
	"strings"
)

type BrowserManager struct {
	CommandIn *termites.InPort

	url string
}

func NewBrowserManager(url string) *BrowserManager {
	builder := termites.NewBuilder("Browser Manager")

	manager := &BrowserManager{
		CommandIn: termites.NewInPort[struct{}](builder),

		url: url,
	}

	builder.OnRun(manager.Run)

	return manager
}

func (b *BrowserManager) Run(c termites.NodeControl) error {
	for range b.CommandIn.Receive() {
		err := RunBrowser(b.url)
		if err != nil {
			c.LogError("Error when opening browser", err)
		}
	}
	return nil
}

func RunBrowser(url string) (err error) {
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}

	// TODO: error is not properly returned
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("cmd", "/c", "start", "chrome", url).Start()
	default:
		err = fmt.Errorf("unsupported platform [%s]", runtime.GOOS)
	}
	return
}
