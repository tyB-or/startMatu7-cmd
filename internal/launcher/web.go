package launcher

import (
	"os/exec"
	"runtime"

	"matu7/pkg/models"
)

// LaunchWebTool 打开网页工具
func LaunchWebTool(tool models.WebTool) error {
	url := tool.URL
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default: // linux, freebsd, etc.
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Start()
}
