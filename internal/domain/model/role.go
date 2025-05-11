package model

import (
	"gorm.io/gorm"
	"time"
)

type Role struct {
	ID          string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Active      bool           `gorm:"default:true"`
	IsSystem    bool           `gorm:"default:false" json:"is_system"`
	Permissions []*Permission  `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
