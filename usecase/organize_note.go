package usecase

import (
	"path/filepath"

	"github.com/ishida722/krapp-go/models"
)

// ラベルとディレクトリのマッピング
type LabelDirectoryMap map[string]string

func OrganizeNotesByCreated(notes []models.Note, baseDirectory string) {
	// Sort notes by CreatedAt in descending order
	sort.Slice(notes, func(i, j int) bool {
		createdI, errI := notes[i].FrontMatter.Created()
		createdJ, errJ := notes[j].FrontMatter.Created()
		if errI != nil || errJ != nil {
			return false // Treat invalid dates as equal
		}
		return createdI.After(createdJ) // Descending order
	})
	for _, note := range notes {
		created, err := note.FrontMatter.Created()
		if err != nil {
			continue // Skip notes without a valid created date
		}
		path := filepath.Join(created.Format("2006"), created.Format("01"))
		fullPath := filepath.Join(baseDirectory, path)
		err = os.MkdirAll(fullPath, 0755)
		if err != nil {
			continue // Skip notes if the directory cannot be created
		}
		err = note.MoveFile(fullPath)
		if err != nil {
			continue // Skip notes that cannot be moved
		}
	}
}

func OrganizeNotesByLabel(notes []models.Note, baseDirectory string, labelDirectoryMap LabelDirectoryMap) {
	// Sort notes by label
	for _, note := range notes {
		label, err := note.FrontMatter.Label()
		if err != nil {
			continue // Skip notes without a valid label
		}
		if label == "" {
			continue // Skip notes without a label
		}
		path, ok := labelDirectoryMap[label]
		if !ok {
			continue // Skip notes with labels not in the map
		}
		fullPath := filepath.Join(baseDirectory, path)
		err = note.MoveFile(fullPath)
		if err != nil {
			continue // Skip notes that cannot be moved
		}
	}
}
