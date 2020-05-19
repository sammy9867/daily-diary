package mysql

import (
	"errors"
	"log"
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

	return db.RowsAffected, nil
}

func (mysqlEntryRepo *mysqlEntryRepository) GetEntryOfUserByID(eid uint64, uid uint64) (*model.Entry, error) {

	defer timeTrack(time.Now(), "GetEntryOfUserByID")
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
		entryImages := []model.EntryImage{}
		if err := mysqlEntryRepo.DB.Raw("CALL GetAllEntryImagesOfEntry(?)", entry.ID).Scan(&entryImages).Error; err != nil {
			return &model.Entry{}, err
		}
		entry.EntryImages = entryImages

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

			entryImages := []model.EntryImage{}
			if err := mysqlEntryRepo.DB.Raw("CALL GetAllEntryImagesOfEntry(?)", entries[i].ID).Scan(&entryImages).Error; err != nil {
				return &[]model.Entry{}, err
			}
			entries[i].EntryImages = entryImages

		}
	}
	return &entries, nil
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("Time Taken by %s is %s", name, elapsed)
}
