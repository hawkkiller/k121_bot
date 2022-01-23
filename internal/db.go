package internal

import (
	"fmt"
	"github.com/hawkkiller/k121_bot/internal/model"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func SetDB(db *gorm.DB) {
	DB = db
	err := DB.AutoMigrate(&model.Pair{})
	err = DB.AutoMigrate(&model.Day{})
	err = DB.AutoMigrate(&model.Time{})
	err = DB.AutoMigrate(&model.Schedule{})

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
}
