package playbook

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// PlaybookInfo はPlaybookの情報を保持する
type PlaybookInfo struct {
	Name        string
	Path        string
	Description string
	IsLocal     bool
}

// ListLocalPlaybooks はローカルのPlaybook一覧を取得する
// basePath: playbooks/containersなどのベースパス
func ListLocalPlaybooks(basePath string) ([]PlaybookInfo, error) {
	var playbooks []PlaybookInfo

	// ディレクトリが存在しない場合は空のリストを返す
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return playbooks, nil
	}

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read playbooks directory: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		playbookPath := filepath.Join(basePath, entry.Name(), "main.yml")
		if _, err := os.Stat(playbookPath); err == nil {
			description := readDescription(filepath.Join(basePath, entry.Name()))
			playbooks = append(playbooks, PlaybookInfo{
				Name:        entry.Name(),
				Path:        playbookPath,
				Description: description,
				IsLocal:     true,
			})
		}
	}

	return playbooks, nil
}

// readDescription はREADME.mdから説明を読み取る（最初の段落）
func readDescription(dir string) string {
	readmePath := filepath.Join(dir, "README.md")
	content, err := os.ReadFile(readmePath)
	if err != nil {
		return "説明なし"
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			if len(line) > 100 {
				return line[:100] + "..."
			}
			return line
		}
	}

	return "説明なし"
}

// DownloadFromGit はGitリポジトリからPlaybookをダウンロードする
// url: GitリポジトリのURL
// basePath: ダウンロード先のベースパス
// name: Playbookの名前（ディレクトリ名）
func DownloadFromGit(url, basePath, name string) error {
	targetDir := filepath.Join(basePath, name)

	// ディレクトリが既に存在する場合は削除
	if _, err := os.Stat(targetDir); err == nil {
		if err := os.RemoveAll(targetDir); err != nil {
			return fmt.Errorf("failed to remove existing directory: %v", err)
		}
	}

	// git cloneを実行
	cmd := exec.Command("git", "clone", url, targetDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone failed: %v, output: %s", err, string(output))
	}

	// main.ymlの存在確認
	playbookPath := filepath.Join(targetDir, "main.yml")
	if _, err := os.Stat(playbookPath); os.IsNotExist(err) {
		os.RemoveAll(targetDir)
		return fmt.Errorf("main.yml not found in repository")
	}

	return nil
}

// DownloadFromURL は指定されたURLからYAMLファイルをダウンロードする
// url: YAMLファイルのURL
// basePath: ダウンロード先のベースパス
// name: Playbookの名前（ディレクトリ名）
func DownloadFromURL(url, basePath, name string) error {
	targetDir := filepath.Join(basePath, name)

	// ディレクトリ作成
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// HTTPリクエスト
	resp, err := http.Get(url)
	if err != nil {
		os.RemoveAll(targetDir)
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		os.RemoveAll(targetDir)
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// main.ymlとして保存
	playbookPath := filepath.Join(targetDir, "main.yml")
	file, err := os.Create(playbookPath)
	if err != nil {
		os.RemoveAll(targetDir)
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		os.RemoveAll(targetDir)
		return fmt.Errorf("failed to save file: %v", err)
	}

	return nil
}

// DeletePlaybook はPlaybookを削除する
func DeletePlaybook(basePath, name string) error {
	targetDir := filepath.Join(basePath, name)

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return fmt.Errorf("playbook not found: %s", name)
	}

	if err := os.RemoveAll(targetDir); err != nil {
		return fmt.Errorf("failed to delete playbook: %v", err)
	}

	return nil
}

// GetPlaybookPath はPlaybook名からmain.ymlのパスを返す
func GetPlaybookPath(basePath, name string) string {
	return filepath.Join(basePath, name, "main.yml")
}

// ValidatePlaybookExists はPlaybookが存在するか確認する
func ValidatePlaybookExists(basePath, name string) error {
	playbookPath := GetPlaybookPath(basePath, name)
	if _, err := os.Stat(playbookPath); os.IsNotExist(err) {
		return fmt.Errorf("playbook not found: %s", name)
	}
	return nil
}
