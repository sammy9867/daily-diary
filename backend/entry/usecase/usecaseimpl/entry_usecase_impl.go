package usecaseimpl

import (
	"github.com/sammy9867/daily-diary/backend/domain"
	"github.com/sammy9867/daily-diary/backend/entry/repository"
	"github.com/sammy9867/daily-diary/backend/entry/usecase"
)

type entryUsecase struct {
	entryRepo repository.EntryRepository
}

// NewEntryUseCase will create an object that will implement EntryUseCase interface
// Note: Need to implement all the methods from the interface
func NewEntryUseCase(er repository.EntryRepository) usecase.EntryUseCase {
	return &entryUsecase{entryRepo: er}
}

func (entryUC *entryUsecase) CreateEntry(e *domain.Entry) (*domain.Entry, error) {
	createdEntry, err := entryUC.entryRepo.CreateEntry(e)
	return createdEntry, err
}

func (entryUC *entryUsecase) UpdateEntry(eid uint64, e *domain.Entry) (*domain.Entry, error) {
	updatedEntry, err := entryUC.entryRepo.UpdateEntry(eid, e)
	return updatedEntry, err
}

func (entryUC *entryUsecase) DeleteEntry(eid, uid uint64) (int64, error) {
	deletedEntry, err := entryUC.entryRepo.DeleteEntry(eid, uid)
	return deletedEntry, err
}

func (entryUC *entryUsecase) GetEntryOfUserByID(eid, uid uint64) (*domain.Entry, error) {
	entry, err := entryUC.entryRepo.GetEntryOfUserByID(eid, uid)
	return entry, err
}

func (entryUC *entryUsecase) GetAllEntriesOfUser(uid uint64, limit, pageNumber, year1, year2 uint32, sort string) (*[]domain.Entry, error) {
	entries, err := entryUC.entryRepo.GetAllEntriesOfUser(uid, limit, pageNumber, year1, year2, sort)
	return entries, err
}
