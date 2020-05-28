package mysql_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //mysql driver
	"github.com/joho/godotenv"

	"github.com/go-playground/assert/v2"
	"github.com/sammy9867/daily-diary/backend/domain"
	_entryRepo "github.com/sammy9867/daily-diary/backend/entry/repository/mysql"
)

type DatabaseConnection struct {
	DB *gorm.DB
}

var dbConn = DatabaseConnection{}

func TestCreateEntry(t *testing.T) {

	mockDB()
	var err error
	user := domain.User{
		ID:       1,
		Username: "Sammy",
		Email:    "sammy@gmail.com",
		Password: "password",
	}

	err = dbConn.DB.Model(&domain.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("Error saving the user: %v\n", err)
	}

	entry := domain.Entry{
		ID:          1,
		Title:       "title sam",
		Description: "description sam",
		OwnerID:     user.ID,
	}

	savedEntry, err := _entryRepo.NewMysqlEntryRepository(dbConn.DB).CreateEntry(&entry)
	if err != nil {
		t.Errorf("Error while saving the entry: %v\n", err)
		return
	}

	assert.Equal(t, entry.ID, savedEntry.ID)
	assert.Equal(t, entry.Title, savedEntry.Title)
	assert.Equal(t, entry.Description, savedEntry.Description)
	assert.Equal(t, entry.OwnerID, savedEntry.OwnerID)

}

func TestUpdateEntry(t *testing.T) {

	mockDB()
	var err error
	user := domain.User{
		ID:       1,
		Username: "Sammy",
		Email:    "sammy@gmail.com",
		Password: "password",
	}

	err = dbConn.DB.Model(&domain.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("Error saving the user: %v\n", err)
	}

	entry := domain.Entry{
		ID:          1,
		Title:       "title sam",
		Description: "description sam",
		OwnerID:     user.ID,
	}

	err = dbConn.DB.Model(&domain.Entry{}).Create(&entry).Error
	if err != nil {
		log.Fatalf("Error saving the entry: %v\n", err)
	}

	entryUpdate := domain.Entry{
		ID:          1,
		Title:       "entryUpdated Title",
		Description: "entryUpdated Desc",
		OwnerID:     entry.OwnerID,
	}

	updatedEntry, err := _entryRepo.NewMysqlEntryRepository(dbConn.DB).UpdateEntry(entryUpdate.ID, &entryUpdate)
	if err != nil {
		t.Errorf("Error while updating the entry: %v\n", err)
		return
	}

	assert.Equal(t, entryUpdate.ID, updatedEntry.ID)
	assert.Equal(t, entryUpdate.Title, updatedEntry.Title)
	assert.Equal(t, entryUpdate.Description, updatedEntry.Description)
	assert.Equal(t, entryUpdate.OwnerID, updatedEntry.OwnerID)

}

func TestDeleteEntryWithoutImage(t *testing.T) {

	mockDB()
	var err error
	user := domain.User{
		ID:       1,
		Username: "Sammy",
		Email:    "sammy@gmail.com",
		Password: "password",
	}

	err = dbConn.DB.Model(&domain.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("Error saving the user: %v\n", err)
	}

	entryWithoutImage := domain.Entry{
		ID:          1,
		Title:       "title sam",
		Description: "description sam",
		OwnerID:     user.ID,
	}

	err = dbConn.DB.Model(&domain.Entry{}).Create(&entryWithoutImage).Error
	if err != nil {
		log.Fatalf("Error saving the entry without image: %v\n", err)
	}

	isDeletedWithoutImage, err := _entryRepo.NewMysqlEntryRepository(dbConn.DB).DeleteEntry(entryWithoutImage.ID, user.ID)
	if err != nil {
		t.Errorf("Error while deleting the entry: %v\n", err)
		return
	}

	assert.Equal(t, isDeletedWithoutImage, int64(1))
}

func TestDeleteEntryWithImage(t *testing.T) {

	mockDB()
	var err error
	user := domain.User{
		ID:       1,
		Username: "Sammy",
		Email:    "sammy@gmail.com",
		Password: "password",
	}

	err = dbConn.DB.Model(&domain.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("Error saving the user: %v\n", err)
	}

	entryWithImage := domain.Entry{
		ID:          1,
		Title:       "title sam",
		Description: "description sam",
		EntryImages: []domain.EntryImage{
			{
				ID:      1,
				URL:     "image url",
				EntryID: 1,
			},
		},
		OwnerID: user.ID,
	}

	err = dbConn.DB.Model(&domain.Entry{}).Create(&entryWithImage).Error
	if err != nil {
		log.Fatalf("Error saving the entry with image: %v\n", err)
	}

	isDeletedWithImage, err := _entryRepo.NewMysqlEntryRepository(dbConn.DB).DeleteEntry(entryWithImage.ID, user.ID)
	if err != nil {
		t.Errorf("Error while deleting the entry: %v\n", err)
		return
	}

	assert.Equal(t, isDeletedWithImage, int64(1))
}

func TestGetEntryOfUserByID(t *testing.T) {

	mockDB()
	var err error
	user := domain.User{
		ID:       1,
		Username: "Sammy",
		Email:    "sammy@gmail.com",
		Password: "password",
	}

	err = dbConn.DB.Model(&domain.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("Error saving the user: %v\n", err)
	}

	entry := domain.Entry{
		ID:          1,
		Title:       "title sam",
		Description: "description sam",
		OwnerID:     user.ID,
	}

	err = dbConn.DB.Model(&domain.Entry{}).Create(&entry).Error
	if err != nil {
		log.Fatalf("Error saving the entry: %v\n", err)
	}

	entryFound, err := _entryRepo.NewMysqlEntryRepository(dbConn.DB).GetEntryOfUserByID(entry.ID, user.ID)
	if err != nil {
		t.Errorf("Error while finding the entry: %v\n", err)
		return
	}

	assert.Equal(t, entry.ID, entryFound.ID)
	assert.Equal(t, entry.Title, entryFound.Title)
	assert.Equal(t, entry.Description, entryFound.Description)
	assert.Equal(t, entry.OwnerID, entryFound.OwnerID)
}

func TestGetAllEntriesOfUserWithoutImage(t *testing.T) {

	mockDB()
	var err error
	user := domain.User{
		ID:       1,
		Username: "Sammy",
		Email:    "sammy@gmail.com",
		Password: "password",
	}

	err = dbConn.DB.Model(&domain.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("Error saving the user: %v\n", err)
	}

	entries := []domain.Entry{

		domain.Entry{
			ID:          1,
			Title:       "title sam",
			Description: "description sam",
			OwnerID:     user.ID,
		},

		domain.Entry{
			ID:          2,
			Title:       "title pam",
			Description: "description pam",
			OwnerID:     user.ID,
		},
	}

	for i := range entries {
		err := dbConn.DB.Model(&domain.Entry{}).Create(&entries[i]).Error
		if err != nil {
			log.Fatalf("Error saving the entry: %v\n", err)
		}
	}

	entriesFound, err := _entryRepo.NewMysqlEntryRepository(dbConn.DB).GetAllEntriesOfUser(user.ID)
	if err != nil {
		t.Errorf("Error while finding users: %v\n", err)
		return
	}

	assert.Equal(t, len(entries), len(*entriesFound))
}

func TestGetAllEntriesOfUserWithImage(t *testing.T) {

	mockDB()
	var err error
	user := domain.User{
		ID:       1,
		Username: "Sammy",
		Email:    "sammy@gmail.com",
		Password: "password",
	}

	err = dbConn.DB.Model(&domain.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("Error saving the user: %v\n", err)
	}

	entries := []domain.Entry{

		domain.Entry{
			ID:          1,
			Title:       "title sam",
			Description: "description sam",
			EntryImages: []domain.EntryImage{
				{
					ID:      1,
					URL:     "image url",
					EntryID: 1,
				},
				{
					ID:      2,
					URL:     "image url2",
					EntryID: 1,
				},
			},
			OwnerID: user.ID,
		},

		domain.Entry{
			ID:          2,
			Title:       "title pam",
			Description: "description pam",
			EntryImages: []domain.EntryImage{
				{
					ID:      3,
					URL:     "image url",
					EntryID: 2,
				},
				{
					ID:      4,
					URL:     "image url2",
					EntryID: 2,
				},
			},
			OwnerID: user.ID,
		},
	}

	for i := range entries {
		err := dbConn.DB.Model(&domain.Entry{}).Create(&entries[i]).Error
		if err != nil {
			log.Fatalf("Error saving the entry: %v\n", err)
		}
	}

	entriesFound, err := _entryRepo.NewMysqlEntryRepository(dbConn.DB).GetAllEntriesOfUser(user.ID)
	if err != nil {
		t.Errorf("Error while finding users: %v\n", err)
		return
	}

	assert.Equal(t, len(entries), len(*entriesFound))
}

func mockDB() {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../../.env"))

	if err != nil {
		log.Fatalf("Error opening env %v\n", err)
	}

	dbConn.InitializeDBTest(os.Getenv("DB_DRIVER_TEST"), os.Getenv("DB_USER_TEST"), os.Getenv("DB_PASSWORD_TEST"), os.Getenv("DB_PORT_TEST"),
		os.Getenv("DB_HOST_TEST"), os.Getenv("DB_NAME_TEST"))

	if err := dbConn.DB.Raw("CALL TrucateTables()").Scan(&domain.EntryImage{}).Scan(&domain.Entry{}).Scan(&domain.User{}).Error; err != nil {
		log.Printf("Error truncating tables: %v\n", err)
	}

	log.Printf("Successfully refreshed table")
}

func (dbConnec *DatabaseConnection) InitializeDBTest(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {

	var err error

	if Dbdriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
		dbConnec.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", Dbdriver)
		}
	}
}
