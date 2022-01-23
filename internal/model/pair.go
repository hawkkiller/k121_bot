package model

import "gorm.io/gorm"

type Pair struct {
	gorm.Model     `json:"-"`
	Title          string `json:"title"`
	Auditory       string `json:"auditory"`
	AdditionalInfo string `json:"additional_info"`
}
