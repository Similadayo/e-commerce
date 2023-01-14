package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Suspension struct {
	gorm.Model
	UserID    uint
	StartTime time.Time
	EndTime   time.Time
	Reason    string
}
