package launcher

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"matu7/pkg/models"
)

// LaunchOfflineTool 启动离线工具
func LaunchOfflineTool(tool models.OfflineTool) error {
	// 检查路径是否存在
	if _, err := os.Stat(tool.Path); os.IsNotExist(err) {
		return fmt.Errorf("工具路径不存在: %s", tool.Path)
	}

	// 如果有指定命令，则使用该命令启动
	if tool.Command != "" {
		cmd := exec.Command("sh", "-c", tool.Command)
		cmd.Dir = tool.Path
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin // 添加标准输入，允许交互

		fmt.Printf("正在工具目录 %s 中执行命令: %s\n", tool.Path, tool.Command)
		return cmd.Start() // 使用Start替代Run，使命令在后台运行
	}

	// 检查目录中是否有jar文件，如果有则尝试运行第一个jar文件
	files, err := os.ReadDir(tool.Path)
	if err == nil {
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".jar") {
				jarPath := filepath.Join(tool.Path, file.Name())
				fmt.Printf("找到JAR文件，尝试运行: %s\n", jarPath)

				cmd := exec.Command("java", "-jar", jarPath)
				cmd.Dir = tool.Path
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Stdin = os.Stdin // 添加标准输入，允许交互
				return cmd.Start()   // 使用Start替代Run，使命令在后台运行
			}
		}
	}

	// 检查目录中是否有可执行文件
	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(tool.Path, file.Name())
			fileInfo, err := os.Stat(filePath)
			if err == nil {
				// 检查文件是否有可执行权限
				if fileInfo.Mode()&0111 != 0 {
					fmt.Printf("找到可执行文件，尝试运行: %s\n", filePath)

					cmd := exec.Command(filePath)
					cmd.Dir = tool.Path
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					cmd.Stdin = os.Stdin // 添加标准输入，允许交互
					return cmd.Start()   // 使用Start替代Run，使命令在后台运行
				}
			}
		}
	}

	// 如果没有找到可执行文件，则打开终端进入工具目录
	fmt.Printf("未找到可执行文件或命令，将打开终端并进入目录: %s\n", tool.Path)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		// 使用AppleScript打开终端并进入指定目录
		script := fmt.Sprintf(`
			tell application "Terminal"
				activate
				do script "cd '%s' && echo '您已进入工具目录：%s' && echo '输入 exit 可退出当前会话' && $SHELL"
			end tell
		`, tool.Path, tool.Path)
		cmd = exec.Command("osascript", "-e", script)
	case "linux":
		// 检查是否有常见的终端模拟器
		terminals := []string{"gnome-terminal", "konsole", "xterm"}
		for _, terminal := range terminals {
			if _, err := exec.LookPath(terminal); err == nil {
				if terminal == "gnome-terminal" {
					cmd = exec.Command(terminal, "--", "bash", "-c", fmt.Sprintf("cd '%s' && echo '您已进入工具目录：%s' && echo '输入 exit 可退出当前会话' && bash", tool.Path, tool.Path))
				} else {
					cmd = exec.Command(terminal, "-e", fmt.Sprintf("cd '%s' && echo '您已进入工具目录：%s' && echo '输入 exit 可退出当前会话' && bash", tool.Path, tool.Path))
				}
				break
			}
		}
		// 如果未找到终端模拟器，则使用xdg-open打开文件管理器
		if cmd == nil {
			cmd = exec.Command("xdg-open", tool.Path)
		}
	case "windows":
		cmd = exec.Command("cmd", "/C", "start", "cmd.exe", "/K", fmt.Sprintf("cd /d \"%s\" && echo 您已进入工具目录：%s && echo 输入 exit 可退出当前会话", tool.Path, tool.Path))
	default:
		// 默认行为，直接打开目录
		cmd = exec.Command("open", tool.Path)
	}

	return cmd.Start() // 使用Start替代Run，使命令在后台运行
}
