package model

import "gorm.io/gorm"

type Day struct {
	gorm.Model `json:"-"`
	Caption    string `json:"caption"`
	Pairs      []Pair `json:"pairs" gorm:"many2many:foreignKey:day_pairs;"`
}
