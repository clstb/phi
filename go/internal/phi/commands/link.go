package commands

import (
	"fmt"
	"os/exec"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/clstb/phi/go/pkg/client"
)

func OpenLink(client *client.Client) tea.Cmd {
	return func() tea.Msg {
		link, err := client.GetLink()
		if err != nil {
			return err
		}
		switch runtime.GOOS {
		case "linux":
			err = exec.Command("xdg-open", link).Start()
		case "windows":
			err = exec.Command("rundll32", "url.dll,FileProtocolHandler", link).Start()
		case "darwin":
			err = exec.Command("open", link).Start()
		default:
			err = fmt.Errorf("unsupported platform")
		}
		if err != nil {
			return err
		}
		return nil
	}
}
