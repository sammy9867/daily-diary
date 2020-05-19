package repository

import "github.com/sammy9867/daily-diary/backend/model"

// EntryRepository represents the entries repository
type EntryRepository interface {
	CreateEntry(*model.Entry) (*model.Entry, error)
	UpdateEntry(uint64, *model.Entry) (*model.Entry, error)
	DeleteEntry(eid, uid uint64) (int64, error)
	GetEntryOfUserByID(eid uint64, uid uint64) (*model.Entry, error)
	GetAllEntriesOfUser(uid uint64) (*[]model.Entry, error)
}
