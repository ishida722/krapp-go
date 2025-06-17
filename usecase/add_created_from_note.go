package usecase

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/ishida722/krapp-go/models"
)

var datePattern = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)

// AddCreatedFromNote parses the note's file name and content to
// find the first date string (YYYY-MM-DD) and sets it as the
// Created field in the frontmatter. If the field already exists,
// the function does nothing. It returns an error when no date
// string can be found.
func AddCreatedFromNote(n *models.Note) error {
	if n == nil {
		return fmt.Errorf("note is nil")
	}
	if n.FrontMatter == nil {
		n.FrontMatter = models.FrontMatter{}
	} else if _, ok := n.FrontMatter["Created"]; ok {
		return nil
	}

	if n.FilePath != "" {
		base := filepath.Base(n.FilePath)
		if m := datePattern.FindString(base); m != "" {
			n.FrontMatter["Created"] = m
			return nil
		}
	}

	if m := datePattern.FindString(n.Content); m != "" {
		n.FrontMatter["Created"] = m
		return nil
	}

	return fmt.Errorf("date string not found")
}
