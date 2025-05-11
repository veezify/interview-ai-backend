package model

import (
	"gorm.io/gorm"
	"time"
)

type Question struct {
	ID          string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Content     string         `gorm:"type:text" json:"content"`
	Type        string         `json:"type"`
	OrderIndex  int            `json:"order_index"`
	Difficulty  string         `json:"difficulty"`
	InterviewID string         `gorm:"type:uuid" json:"interview_id"`
	Interview   Interview      `gorm:"foreignKey:InterviewID" json:"interview"`
	Response    *Response      `gorm:"foreignKey:QuestionID" json:"response"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
