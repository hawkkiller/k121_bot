package model

import "gorm.io/gorm"

type Time struct {
	gorm.Model `json:"-"`
	Start      string `json:"start"`
	End        string `json:"end"`
}
