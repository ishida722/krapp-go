package usecase

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ImportNotes copies .txt and .md files from src to dst recursively.
// .txt files are converted to .md and all files are saved as UTF-8.
func ImportNotes(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		fmt.Println("Processing:", path)
		ext := strings.ToLower(filepath.Ext(d.Name()))
		if ext != ".txt" && ext != ".md" {
			return nil
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		newName := d.Name()
		if ext == ".txt" {
			newName = strings.TrimSuffix(d.Name(), ".txt") + ".md"
		}
		dstPath := filepath.Join(dst, filepath.Dir(rel), newName)
		if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
			return err
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		out, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer out.Close()
		// 読み込み→UTF-8で書き込み（GoのstringはUTF-8なのでそのままコピーでOK）
		if _, err := io.Copy(out, in); err != nil {
			return err
		}
		return nil
	})
}
