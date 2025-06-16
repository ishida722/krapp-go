package models

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

type Note struct {
	FrontMatter FrontMatter
	Content     string
	FilePath    string // ファイルパス
}

func (note Note) ToString() (string, error) {
	if note.FrontMatter == nil {
		return note.Content, nil
	}
	fmStr, err := note.FrontMatter.ToYAML()
	if err != nil {
		return "", err
	}
	return fmStr + "\n" + note.Content, nil
}

func ParseNote(raw string) (*Note, error) {
	// frontmatterがあるか判定
	if !strings.HasPrefix(raw, "---\n") {
		// なければbodyと空のフロントマターを返す
		return &Note{
			FrontMatter: FrontMatter{},
			Content:     strings.TrimSpace(raw),
		}, nil
	}

	// --- で分割
	parts := strings.SplitN(raw, "---", 3)
	if len(parts) < 3 {
		// フロントマターのフォーマットが不正
		return nil, errors.New("invalid format")
	}

	var fm FrontMatter
	err := yaml.Unmarshal([]byte(parts[1]), &fm)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	return &Note{
		FrontMatter: fm,
		Content:     strings.TrimSpace(parts[2]),
	}, nil
}

func LoadNoteFromFile(filePath string) (*Note, error) {
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	note, err := ParseNote(string(raw))
	if err != nil {
		return nil, fmt.Errorf("failed to parse note from file %s: %w", filePath, err)
	}
	note.FilePath = filePath
	return note, nil
}

func (note Note) SaveToFile() error {
	content, err := note.ToString()
	if err != nil {
		return fmt.Errorf("failed to convert note to string: %w", err)
	}

	if note.FilePath == "" {
		return errors.New("note file path is empty")
	}

	err = os.WriteFile(note.FilePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write note to file %s: %w", note.FilePath, err)
	}
	return nil
}

func (note *Note) MoveFile(newDir string) error {
	if note.FilePath == "" {
		return errors.New("note file path is empty")
	}
	newPath := path.Join(newDir, path.Base(note.FilePath))

	err := os.Rename(note.FilePath, newPath)
	if err != nil {
		return fmt.Errorf("failed to move note file from %s to %s: %w", note.FilePath, newPath, err)
	}

	note.FilePath = newPath
	return nil
}
