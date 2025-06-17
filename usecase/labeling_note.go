package usecase

import (
	"github.com/ishida722/krapp-go/models"
)

type NoteLabeler struct {
	Notes        []*models.Note
	Current      *models.Note
	CurrentIndex int
}

func (self *NoteLabeler) Next() error {
	self.Current = self.Notes[0]
	return nil
}
