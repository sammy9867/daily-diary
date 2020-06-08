package usecase

import "github.com/sammy9867/daily-diary/backend/domain"

// EntryUseCase represents the entries usecase
type EntryUseCase interface {
	CreateEntry(*domain.Entry) (*domain.Entry, error)
	UpdateEntry(uint64, *domain.Entry) (*domain.Entry, error)
	DeleteEntry(eid, uid uint64) (int64, error)
	GetEntryOfUserByID(eid, uid uint64) (*domain.Entry, error)
	GetAllEntriesOfUser(uid uint64, limit, pageNumber, year1, year2 uint32, sort string) (*[]domain.Entry, error)
}
