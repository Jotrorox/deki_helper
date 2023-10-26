package main

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Command struct {
	ID       uint
	Trigger  string
	Response string
}

func connectToSQLite(cfg *Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DB_PATH), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createEntryTable(db *gorm.DB) error {
	err := db.AutoMigrate(&Command{})
	if err != nil {
		return err
	}
	return nil
}

func addEntry(db *gorm.DB, cmd Command) error {
	entry := cmd
	result := db.Create(&entry)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateRow(db *gorm.DB, keyword string, newData Command) error {
	return db.Model(&Command{}).Where("trigger = ?", keyword).Updates(newData).Error
}

func RemoveRow(db *gorm.DB, keyword string) error {
	return db.Where("trigger = ?", keyword).Delete(&Command{}).Error
}

func queryEntries(db *gorm.DB) ([]Command, error) {
	var entries []Command
	result := db.Find(&entries)
	if result.Error != nil {
		return nil, result.Error
	}
	return entries, nil
}
