package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// GlobalDB - глобальная переменная, представляющая базу данных
	GlobalDB *gorm.DB
)

func InitDatabase() (*gorm.DB, error) {
	// Параметры подключения к PostgreSQL базе данных
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=blogdb sslmode=disable password=postgres")
	if err != nil {
		return nil, err
	}
	GlobalDB = db // Присваиваем глобальной переменной значение
	return db, nil
}
