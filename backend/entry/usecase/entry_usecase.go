package usecase

import "github.com/sammy9867/daily-diary/backend/model"

// EntryUseCase represents the entries usecase
type EntryUseCase interface {
	CreateEntry(*model.Entry) (*model.Entry, error)
	UpdateEntry(uint64, *model.Entry) (*model.Entry, error)
	DeleteEntry(eid, uid uint64) (int64, error)
	GetEntryOfUserByID(eid, uid uint64) (*model.Entry, error)
	GetAllEntriesOfUser(uid uint64) (*[]model.Entry, error)
}
