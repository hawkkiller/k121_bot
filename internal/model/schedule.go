package model

import "gorm.io/gorm"

type Schedule struct {
	gorm.Model `json:"-"`
	ChatId     int64  `json:"-" gorm:"primaryKey;autoIncrement:false"`
	Days       []Day  `json:"days" gorm:"many2many:foreignKey:schedule_days;"`
	Times      []Time `json:"times" gorm:"many2many:foreignKey:schedule_times;"`
}
