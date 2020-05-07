package usecaseimpl

import (
	"github.com/sammy9867/daily-diary/backend/entry/repository"
	"github.com/sammy9867/daily-diary/backend/entry/usecase"
	"github.com/sammy9867/daily-diary/backend/model"
)

type entryUsecase struct {
	entryRepo repository.EntryRepository
}

// NewEntryUseCase will create an object that will implement EntryUseCase interface
// Note: Need to implement all the methods from the interface
func NewEntryUseCase(er repository.EntryRepository) usecase.EntryUseCase {
	return &entryUsecase{entryRepo: er}
}

func (entryUC *entryUsecase) CreateEntry(e *model.Entry) (*model.Entry, error) {
	createdEntry, err := entryUC.entryRepo.CreateEntry(e)
	return createdEntry, err
}

func (entryUC *entryUsecase) UpdateEntry(eid uint64, e *model.Entry) (*model.Entry, error) {
	updatedEntry, err := entryUC.entryRepo.UpdateEntry(eid, e)
	return updatedEntry, err
}

func (entryUC *entryUsecase) DeleteEntry(eid, uid uint64) (int64, error) {
	deletedEntry, err := entryUC.entryRepo.DeleteEntry(eid, uid)
	return deletedEntry, err
}

func (entryUC *entryUsecase) GetEntryOfUserByID(eid, uid uint64) (*model.Entry, error) {
	entry, err := entryUC.entryRepo.GetEntryOfUserByID(eid, uid)
	return entry, err
}

func (entryUC *entryUsecase) GetAllEntriesOfUser(uid uint64) (*[]model.Entry, error) {
	entries, err := entryUC.entryRepo.GetAllEntriesOfUser(uid)
	return entries, err
}
