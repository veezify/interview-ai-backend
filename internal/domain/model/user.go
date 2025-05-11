package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID         string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	FirstName  string         `json:"first_name"`
	LastName   string         `json:"last_name"`
	Email      string         `json:"email"`
	Password   string         `json:"password"`
	Phone      string         `json:"phone"`
	Active     bool           `gorm:"default:true"`
	Avatar     string         `json:"avatar"`
	Language   string         `json:"language"`
	Interviews []Interview    `gorm:"foreignKey:UserID" json:"interviews"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
