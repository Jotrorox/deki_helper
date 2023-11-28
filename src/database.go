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
	hpe(err)
	return db, nil
}

func createEntryTable(db *gorm.DB) {
	err := db.AutoMigrate(&Command{})
	hpe(err)
}

func addEntry(db *gorm.DB, cmd Command) {
	entry := cmd
	result := db.Create(&entry)
	hpe(result.Error)
}

func UpdateRow(db *gorm.DB, keyword string, newData Command) {
	hpe(db.Model(&Command{}).Where("trigger = ?", keyword).Updates(newData).Error)
}

func RemoveRow(db *gorm.DB, keyword string) {
	hpe(db.Where("trigger = ?", keyword).Delete(&Command{}).Error)
}

func queryEntries(db *gorm.DB) []Command {
	var entries []Command
	result := db.Find(&entries)
	hpe(result.Error)
	return entries
}
