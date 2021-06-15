package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitDb() error {
	config := &gorm.Config{NamingStrategy: &schema.NamingStrategy{SingularTable: true}}
	var err error
	DB, err = gorm.Open(sqlite.Open("db/jzsgic.db"), config)
	if err != nil {
		return err
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Ping()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(10)
	return autoMigrateTable()
}

func autoMigrateTable() error {
	err := DB.AutoMigrate(&TService{})
	if err != nil {
		return err
	}
	err = DB.AutoMigrate(&Block{})
	if err != nil {
		return err
	}
	err = DB.AutoMigrate(&Tx{})
	if err != nil {
		return err
	}
	return DB.AutoMigrate(&TPolicy{})
}
