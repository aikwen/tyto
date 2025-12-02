package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Sync(repoURL, localDir string) (string, error) {
	// 1. 检查是否存在 .git 目录
	gitDir := filepath.Join(localDir, ".git")
	_, err := os.Stat(gitDir)

	if os.IsNotExist(err) {
		// .git 不存在，说明不是仓库，需要 Clone
		return Clone(repoURL, localDir)
	}

	// .git 存在，说明已经是仓库，直接 Pull
	return Pull(localDir)
}



func Clone(repoURL, dir string) (string, error) {
	fmt.Printf("[Git] Cloning %s to %s...\n", repoURL, dir)

	// 确保父目录存在 (例如 /opt/tyto)
	parentDir := filepath.Dir(dir)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create parent dir: %v", err)
	}

	// 执行命令: git clone --depth 1 <url> <dir>
	cmd := exec.Command("git", "clone", "--depth", "1", repoURL, dir)

	outputBytes, err := cmd.CombinedOutput()
	output := strings.TrimSpace(string(outputBytes))

	if err != nil {
		return output, fmt.Errorf("git clone failed: %w, output: %s", err, output)
	}

	return output, nil
}

func Pull(path string) (string, error) {
	cmd := exec.Command("git", "-C", path, "pull", "--depth", "1") // 建议加上 depth 1 保持一致
	outputBytes, err := cmd.CombinedOutput()
	// ... (保持原有的错误处理)
    output := strings.TrimSpace(string(outputBytes))
	if err != nil {
		return output, fmt.Errorf("git pull failed: %w, output: %s", err, output)
	}
	return output, nil
}