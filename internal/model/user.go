package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `json:"name" gorm:"index;unique;name"`
	Phone     string         `json:"phone" gorm:"index;phone"`
	Email     string         `json:"email" gorm:"index;email"`
	Password  string         `json:"password" gorm:"password"`
	Avatar    string         `json:"avatar" gorm:"avatar"`
	Age       int            `json:"age" gorm:"age"`
	ID        uint           `gorm:"primaryKey;autoIncrement"`
}
