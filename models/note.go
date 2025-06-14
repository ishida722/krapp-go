package models

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Note struct {
	FrontMatter FrontMatter
	Content     string
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
