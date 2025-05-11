package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"time"
)

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, j)
}

type Response struct {
	ID              string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	TextResponse    string         `gorm:"type:text" json:"text_response"`
	AudioURL        string         `json:"audio_url"`
	InterviewID     string         `gorm:"type:uuid;index" json:"interview_id"`
	Interview       Interview      `gorm:"foreignKey:InterviewID" json:"interview"`
	QuestionID      string         `gorm:"type:uuid;uniqueIndex" json:"question_id"`
	Question        Question       `gorm:"foreignKey:QuestionID" json:"question"`
	ContentAnalysis JSONB          `gorm:"type:jsonb" json:"content_analysis"`
	VoiceAnalysis   JSONB          `gorm:"type:jsonb" json:"voice_analysis"`
	VideoAnalysis   JSONB          `gorm:"type:jsonb" json:"video_analysis"`
	ResponseScore   *float64       `json:"response_score"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
