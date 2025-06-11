package usecase

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Label struct {
	Key  string
	Name string
	Path string
}

type LabelConfig struct {
	Labels []Label
}

// LoadLabels loads label definitions from a yaml file (labels.yaml)
func LoadLabelsYaml(path string) ([]Label, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var labels []Label
	var current Label
	var inLabels bool
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "labels:" {
			inLabels = true
			continue
		}
		if !inLabels || line == "" || strings.HasPrefix(line, "#") {
			continue









































































































































}	return ""	}		}			return strings.TrimSpace(strings.TrimPrefix(line, "label:"))		if strings.HasPrefix(line, "label:") {	for _, line := range lines {	lines := strings.Split(content, "\n")	}		return ""	if !strings.HasPrefix(content, "---\nlabel:") {func extractLabel(content string) string {}	return nil	}		}			return err		if err := os.Rename(file, destPath); err != nil {		}			return err		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {		}			continue		if destPath == "" {		}			}				break				destPath = filepath.Join(dir, l.Path, filepath.Base(file))			if l.Name == label {		for _, l := range labels {		var destPath string		}			continue		if label == "" {		label := extractLabel(string(content))		}			return err		if err != nil {		content, err := os.ReadFile(file)	for _, file := range files {	files = append(files, files2...)	}		return err	if err != nil {	files2, err := filepath.Glob(filepath.Join(dir, "*.txt"))	}		return err	if err != nil {	files, err := filepath.Glob(filepath.Join(dir, "*.md"))func SortTextFiles(dir string, labels []Label) error {// SortTextFiles moves files to label-specific folders}	return strings.HasPrefix(content, "---\nlabel:")func hasLabelFrontMatter(content string) bool {}	return nil	}		}			return err		if err := os.WriteFile(file, []byte(newContent), 0644); err != nil {		newContent := fmt.Sprintf("---\nlabel: %s\n---\n%s", labelName, string(content))		}			continue			fmt.Println("無効な入力。スキップします。")		if labelName == "" {		}			}				break				labelName = l.Name			if l.Key == input {		for _, l := range labels {		var labelName string		input = strings.TrimSpace(input)		input, _ := stdin.ReadString('\n')		fmt.Print("\nラベル番号を入力: ")		}			fmt.Printf("[%s] %s  ", l.Key, l.Name)		for _, l := range labels {		fmt.Println(preview)		}			preview = preview[:200] + "..."		if len(preview) > 200 {		preview := string(content)		fmt.Printf("\n==== %s ===\n", filepath.Base(file))		}			continue // skip already labeled		if hasLabelFrontMatter(string(content)) {		}			return err		if err != nil {		content, err := os.ReadFile(file)	for _, file := range files {	stdin := bufio.NewReader(os.Stdin)	}		return errors.New("no text files found")	if len(files) == 0 {	files = append(files, files2...)	}		return err	if err != nil {	files2, err := filepath.Glob(filepath.Join(dir, "*.txt"))	}		return err	if err != nil {	files, err := filepath.Glob(filepath.Join(dir, "*.md"))func LabelTextFiles(dir string, labels []Label) error {// LabelTextFiles applies labels to text files in dir, interactively}	return labels, nil	}		return nil, errors.New("no labels found")	if len(labels) == 0 {	}		labels = append(labels, current)	if current.Key != "" {	}		}			current.Path = strings.TrimSpace(strings.TrimPrefix(line, "path:"))		if strings.HasPrefix(line, "path:") {		}			current.Name = strings.TrimSpace(strings.TrimPrefix(line, "name:"))		if strings.HasPrefix(line, "name:") {		}			current.Key = strings.TrimSpace(strings.TrimPrefix(line, "key:"))		if strings.HasPrefix(line, "key:") {		}			continue			current = Label{}			}				labels = append(labels, current)			if current.Key != "" {		if strings.HasPrefix(line, "-") {		}		}
		if strings.HasPrefix(line, "-") {
			if current.Key != "" {
				labels = append(labels, current)
			}
			current = Label{}
			continue
		}
		if strings.HasPrefix(line, "key:") {
			current.Key = strings.TrimSpace(strings.TrimPrefix(line, "key:"))
		}
		if strings.HasPrefix(line, "name:") {
			current.Name = strings.TrimSpace(strings.TrimPrefix(line, "name:"))
		}
		if strings.HasPrefix(line, "path:") {
			current.Path = strings.TrimSpace(strings.TrimPrefix(line, "path:"))
		}
	}
	if current.Key != "" {
		labels = append(labels, current)
	}
	if len(labels) == 0 {
		return nil, errors.New("no labels found")
	}
	return labels, nil
}

// LabelTextFiles applies labels to text files in dir, interactively
func LabelTextFiles(dir string, labels []Label) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return err
	}
	files2, err := filepath.Glob(filepath.Join(dir, "*.txt"))
	if err != nil {
		return err
	}
	files = append(files, files2...)
	if len(files) == 0 {
		return errors.New("no text files found")
	}
	stdin := bufio.NewReader(os.Stdin)
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		if hasLabelFrontMatter(string(content)) {
			continue // skip already labeled
		}
		fmt.Printf("\n==== %s ===\n", filepath.Base(file))
		preview := string(content)
		if len(preview) > 200 {
			preview = preview[:200] + "..."
		}
		fmt.Println(preview)
		for _, l := range labels {
			fmt.Printf("[%s] %s  ", l.Key, l.Name)
		}
		fmt.Print("\nラベル番号を入力: ")
		input, _ := stdin.ReadString('\n')
		input = strings.TrimSpace(input)
		var labelName string
		for _, l := range labels {
			if l.Key == input {
				labelName = l.Name
				break
			}
		}
		if labelName == "" {
			fmt.Println("無効な入力。スキップします。")
			continue
		}
		newContent := fmt.Sprintf("---\nlabel: %s\n---\n%s", labelName, string(content))
		if err := os.WriteFile(file, []byte(newContent), 0644); err != nil {
			return err
		}
	}
	return nil
}

func hasLabelFrontMatter(content string) bool {
	return strings.HasPrefix(content, "---\nlabel:")
}

// SortTextFiles moves files to label-specific folders
func SortTextFiles(dir string, labels []Label) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return err
	}
	files2, err := filepath.Glob(filepath.Join(dir, "*.txt"))
	if err != nil {
		return err
	}
	files = append(files, files2...)
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		label := extractLabel(string(content))
		if label == "" {
			continue
		}
		var destPath string
		for _, l := range labels {
			if l.Name == label {
				destPath = filepath.Join(dir, l.Path, filepath.Base(file))
				break
			}
		}
		if destPath == "" {
			continue
		}
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}
		if err := os.Rename(file, destPath); err != nil {
			return err
		}
	}
	return nil
}

func extractLabel(content string) string {
	if !strings.HasPrefix(content, "---\nlabel:") {
		return ""
	}
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "label:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "label:"))
		}
	}
	return ""
}
