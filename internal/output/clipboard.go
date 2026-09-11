package output

import (
	"bytes"
	"os/exec"
	"runtime"
)

func CopyToClipboard(data []byte) error {
	commandName, args := clipboardCommand()

	command := exec.Command(
		commandName,
		args...,
	)

	command.Stdin = bytes.NewReader(data)

	return command.Run()
}

func clipboardCommand() (string, []string) {
	switch runtime.GOOS {
	case "darwin":
		return "pbcopy", nil

	case "windows":
		return "clip", nil

	default:
		return "xclip", []string{
			"-selection",
			"clipboard",
		}
	}
}
