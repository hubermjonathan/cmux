package notify

import (
	"fmt"
	"os"

	"github.com/hubermjonathan/cmux/internal/tmux"
)

func BellBytes() []byte {
	return []byte("\a")
}

func OSC9Bytes(message string) []byte {
	return []byte(fmt.Sprintf("\x1b]9;%s\x1b\\", message))
}

func SendBell(paneID int) {
	writeToPane(paneID, BellBytes())
}

func SendOSC9(paneID int, message string) {
	writeToPane(paneID, OSC9Bytes(message))
}

func SendAll(paneID int, bell bool, osc9Message string) {
	tty, err := tmux.PaneTTY(fmt.Sprintf("%d", paneID))
	if err != nil {
		return
	}
	f, err := os.OpenFile(tty, os.O_WRONLY, 0)
	if err != nil {
		return
	}
	defer f.Close()
	if bell {
		f.Write(BellBytes())
	}
	if osc9Message != "" {
		f.Write(OSC9Bytes(osc9Message))
	}
}

func writeToPane(paneID int, data []byte) {
	tty, err := tmux.PaneTTY(fmt.Sprintf("%d", paneID))
	if err != nil {
		return
	}
	f, err := os.OpenFile(tty, os.O_WRONLY, 0)
	if err != nil {
		return
	}
	defer f.Close()
	f.Write(data)
}
