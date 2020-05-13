package mysql

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sammy9867/daily-diary/backend/entry/repository"
	"github.com/sammy9867/daily-diary/backend/model"
)

type mysqlEntryRepository struct {
	DB *gorm.DB
}

// NewMysqlEntryRepository will create an object that will implement EntryRepository interface
// Note: Need to implement all the methods from the interface
func NewMysqlEntryRepository(DB *gorm.DB) repository.EntryRepository {
	return &mysqlEntryRepository{DB}
}

func (mysqlEntryRepo *mysqlEntryRepository) CreateEntry(entry *model.Entry) (*model.Entry, error) {
	var err error

	err = mysqlEntryRepo.DB.Debug().Create(&entry).Error
	if err != nil {
		return &model.Entry{}, err
	}

	if entry.ID != 0 {
		err = mysqlEntryRepo.DB.Debug().Model(&model.User{}).Where("id = ?", entry.OwnerID).Take(&entry.Owner).Error
		if err != nil {
			return &model.Entry{}, err
		}
	}
	return entry, nil
}

func (mysqlEntryRepo *mysqlEntryRepository) UpdateEntry(eid uint64, entry *model.Entry) (*model.Entry, error) {
	var err error

	db := mysqlEntryRepo.DB.Debug().Model(&model.Entry{}).Where("id = ?", eid).UpdateColumns(
		map[string]interface{}{
			"title":       entry.Title,
			"description": entry.Description,
			"images":      entry.EntryImages,
			"updated_at":  time.Now(),
		},
	)
	if db.Error != nil {
		return &model.Entry{}, err
	}
	if entry.ID != 0 {
		err = mysqlEntryRepo.DB.Debug().Model(&model.Entry{}).Where("id = ?", entry.OwnerID).Take(&entry.Owner).Error
		if err != nil {
			return &model.Entry{}, err
		}
	}
	return entry, nil
}

func (mysqlEntryRepo *mysqlEntryRepository) DeleteEntry(eid uint64, uid uint64) (int64, error) {
	db := mysqlEntryRepo.DB.Debug().Model(&model.Entry{}).Where("id = ? and owner_id = ?", eid, uid).Take(&model.Entry{}).Delete(&model.Entry{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Entry not found")
		}
		return 0, db.Error
	}

	_, err := mysqlEntryRepo.DeleteEntryImages(eid)
	if err != nil {
		return 0, err
	}

	return db.RowsAffected, nil
}

func (mysqlEntryRepo *mysqlEntryRepository) GetEntryOfUserByID(eid uint64, uid uint64) (*model.Entry, error) {
	var err error
	entry := model.Entry{}
	err = mysqlEntryRepo.DB.Debug().Model(&model.Entry{}).Where("id = ? and owner_id = ?", eid, uid).Take(&entry).Error
	if err != nil {
		return &model.Entry{}, err
	}
	if entry.ID != 0 {
		err = mysqlEntryRepo.DB.Debug().Model(&model.Entry{}).Where("id = ?", entry.OwnerID).Take(&entry.Owner).Error
		if err != nil {
			return &model.Entry{}, err
		}
		entryImages, err := mysqlEntryRepo.GetAllEntryImagesOfEntry(entry.ID)
		if err != nil {
			return &model.Entry{}, err
		}
		entry.EntryImages = *entryImages
	}
	return &entry, nil
}

func (mysqlEntryRepo *mysqlEntryRepository) GetAllEntriesOfUser(uid uint64) (*[]model.Entry, error) {

	var err error
	entries := []model.Entry{}
	err = mysqlEntryRepo.DB.Debug().Model(&model.Entry{}).Where("owner_id = ?", uid).Limit(100).Find(&entries).Error
	if err != nil {
		return &[]model.Entry{}, err
	}
	if len(entries) > 0 {
		for i := range entries {
			err := mysqlEntryRepo.DB.Debug().Model(&model.Entry{}).Where("id = ?", entries[i].OwnerID).Take(&entries[i].Owner).Error
			if err != nil {
				return &[]model.Entry{}, err
			}

			entryImages, err := mysqlEntryRepo.GetAllEntryImagesOfEntry(entries[i].ID)
			if err != nil {
				return &[]model.Entry{}, err
			}
			entries[i].EntryImages = *entryImages

		}
	}
	return &entries, nil
}

// DeleteEntryImages will delete all images of an entry if the entry is deleted
func (mysqlEntryRepo *mysqlEntryRepository) DeleteEntryImages(eid uint64) (int64, error) {
	db := mysqlEntryRepo.DB.Debug().Model(&model.EntryImage{}).Where("entry_id = ?", eid).Find(&model.EntryImage{}).Delete(&model.EntryImage{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, nil
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// GetAllEntryImagesOfEntry will get all images of an entry
func (mysqlEntryRepo *mysqlEntryRepository) GetAllEntryImagesOfEntry(eid uint64) (*[]model.EntryImage, error) {
	var err error
	entryImages := []model.EntryImage{}

	err = mysqlEntryRepo.DB.Debug().Model(model.EntryImage{}).Where("entry_id = ?", eid).Limit(100).Find(&entryImages).Error
	fmt.Println(len(entryImages))
	if err != nil {
		return &[]model.EntryImage{}, err
	}

	return &entryImages, nil
}
