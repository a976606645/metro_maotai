package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSqlite() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("mdl.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库文件加载错误")
	}

	return db
}
