package notify

import (
	"fmt"
	"os"

	"github.com/jon-huber/cmux/internal/tmux"
)

func BellBytes() []byte {
	return []byte("\a")
}

func OSC9Bytes(message string) []byte {
	return []byte(fmt.Sprintf("\x1b]9;%s\x1b\\", message))
}

func SendBell(paneID int) {
	tty, err := tmux.PaneTTY(fmt.Sprintf("%d", paneID))
	if err != nil {
		return
	}
	f, err := os.OpenFile(tty, os.O_WRONLY, 0)
	if err != nil {
		return
	}
	defer f.Close()
	f.Write(BellBytes())
}

func SendOSC9(paneID int, message string) {
	tty, err := tmux.PaneTTY(fmt.Sprintf("%d", paneID))
	if err != nil {
		return
	}
	f, err := os.OpenFile(tty, os.O_WRONLY, 0)
	if err != nil {
		return
	}
	defer f.Close()
	f.Write(OSC9Bytes(message))
}
