package domain

import (
	"time"
)

type Base struct {
	ID        uint64    `gorm:"primary_key" sql:"AUTO_INCREMENT" json:"id"`
	CreatedAt time.Time `sql:"not null;type:date"  json:"created_at"`
	UpdatedAt time.Time `sql:"not null;type:date"  json:"updated_at"`
}
