package termites_web

import (
	"fmt"
	"github.com/RoelofRuis/termites/termites"
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
		CommandIn: builder.InPort("Command", struct{}{}),

		url: url,
	}

	builder.OnRun(manager.Run)

	return manager
}

func (b *BrowserManager) Run(c termites.NodeControl) error {
	for range b.CommandIn.Receive() {
		err := RunBrowser(runtime.GOOS, b.url)
		if err != nil {
			c.LogError("Error when opening browser", err)
		}
	}
	return nil
}

func RunBrowser(os string, url string) (err error) {
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}

	// TODO: error is not properly returned
	switch os {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("cmd", "/c", "start", "chrome", url).Start()
	default:
		err = fmt.Errorf("unsupported platform [%s]", os)
	}
	return
}
