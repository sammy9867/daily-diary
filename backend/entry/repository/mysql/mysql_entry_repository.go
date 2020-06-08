package mysql

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sammy9867/daily-diary/backend/domain"
	"github.com/sammy9867/daily-diary/backend/entry/repository"
)

type mysqlEntryRepository struct {
	DB *gorm.DB
}

// NewMysqlEntryRepository will create an object that will implement EntryRepository interface
// Note: Need to implement all the methods from the interface
func NewMysqlEntryRepository(DB *gorm.DB) repository.EntryRepository {
	return &mysqlEntryRepository{DB}
}

func (mysqlEntryRepo *mysqlEntryRepository) CreateEntry(entry *domain.Entry) (*domain.Entry, error) {
	var err error

	err = mysqlEntryRepo.DB.Debug().Create(&entry).Error
	if err != nil {
		return &domain.Entry{}, err
	}

	return entry, nil
}

func (mysqlEntryRepo *mysqlEntryRepository) UpdateEntry(eid uint64, entry *domain.Entry) (*domain.Entry, error) {
	var err error

	if entry.EntryImages != nil {

		for i := range entry.EntryImages {

			// Updates Method automatically updated UpdatedAt with the new timestamp
			db := mysqlEntryRepo.DB.Debug().Model(&domain.EntryImage{}).Where("id = ? AND entry_id = ?", entry.EntryImages[i].ID, entry.EntryImages[i].EntryID).Updates(
				domain.EntryImage{
					URL: entry.EntryImages[i].URL,
				},
			)
			if db.Error != nil {
				return &domain.Entry{}, err
			}

		}

	}

	db := mysqlEntryRepo.DB.Debug().Model(&domain.Entry{}).Where("id = ?", eid).UpdateColumns(
		map[string]interface{}{
			"title":       entry.Title,
			"description": entry.Description,
			"updated_at":  time.Now(),
		},
	)
	if db.Error != nil {
		return &domain.Entry{}, err
	}

	return entry, nil
}

func (mysqlEntryRepo *mysqlEntryRepository) DeleteEntry(eid uint64, uid uint64) (int64, error) {

	db := mysqlEntryRepo.DB.Debug().Model(&domain.Entry{}).Where("id = ? and owner_id = ?", eid, uid).Take(&domain.Entry{}).Delete(&domain.Entry{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Entry not found")
		}
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

func (mysqlEntryRepo *mysqlEntryRepository) GetEntryOfUserByID(eid uint64, uid uint64) (*domain.Entry, error) {

	defer timeTrack(time.Now(), "GetEntryOfUserByID")
	var err error
	entry := domain.Entry{}
	err = mysqlEntryRepo.DB.Debug().Model(&domain.Entry{}).Where("id = ? and owner_id = ?", eid, uid).Take(&entry).Error
	if err != nil {
		return &domain.Entry{}, err
	}

	if entry.ID != 0 {
		entryImages := []domain.EntryImage{}
		if err := mysqlEntryRepo.DB.Raw("CALL GetAllEntryImagesOfEntry(?)", entry.ID).Scan(&entryImages).Error; err != nil {
			return &domain.Entry{}, err
		}
		entry.EntryImages = entryImages

	}
	return &entry, nil
}

func (mysqlEntryRepo *mysqlEntryRepository) GetAllEntriesOfUser(uid uint64, limit, pageNumber, year1, year2 uint32, sort string) (*[]domain.Entry, error) {

	var err error
	entries := []domain.Entry{}

	order := orderBy(sort)

	if year1 == 0 && year2 == 0 { // If no year is provided
		err = mysqlEntryRepo.DB.Debug().Model(&domain.Entry{}).Where("owner_id = ?", uid).Order(order).Limit(limit).Offset(limit * (pageNumber - 1)).Find(&entries).Error

	} else {
		var yearStart, yearEnd uint32

		if yearStart = year1; year1 == 0 {
			yearStart = year2
		}

		if yearEnd = year2; year2 == 0 {
			yearEnd = year1
		}

		err = mysqlEntryRepo.DB.Debug().Model(&domain.Entry{}).Where("owner_id = ? AND created_at BETWEEN STR_TO_DATE(?, '%Y') AND STR_TO_DATE(?, '%Y')", uid, yearStart, yearEnd+1).Order(order).Limit(limit).Offset(limit * (pageNumber - 1)).Find(&entries).Error

	}

	if err != nil {
		return &[]domain.Entry{}, err
	}

	if len(entries) > 0 {
		for i := range entries {
			entryImages := []domain.EntryImage{}
			if err := mysqlEntryRepo.DB.Raw("CALL GetAllEntryImagesOfEntry(?)", entries[i].ID).Scan(&entryImages).Error; err != nil {
				return &[]domain.Entry{}, err
			}
			entries[i].EntryImages = entryImages

		}
	}
	return &entries, nil
}

// get all the orders()
func orderBy(sort string) string {
	var order string

	filters := strings.Split(sort, ",")

	for _, filter := range filters {
		if filter[0] == '-' {
			order += filter[1:len(filter)] + " desc, "
		} else {
			order += filter + ", "
		}
	}

	order += "id desc"
	return order
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("Time Taken by %s is %s", name, elapsed)
}
