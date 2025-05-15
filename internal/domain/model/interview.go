package model

import (
	"github.com/lib/pq"
	//swaggerignore:start
	//swaggerignore:end
	"gorm.io/gorm"
	"time"
)

type Interview struct {
	ID      string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Level   string `json:"level"`
	JobType string `json:"job_type"`
	//swaggerignore:start
	ProgrammingLanguages pq.StringArray `gorm:"type:text[]" json:"programming_languages"`
	//ProgrammingLanguages []string `gorm:"type:text[]" json:"programming_languages"`
	//swaggerignore:end

	//swagger:include
	// ProgrammingLanguages array of strings
	// swagger:allOf
	// ProgrammingLanguages []string `json:"programming_languages"`
	InterviewLanguage  string         `json:"interview_language"`
	Country            string         `json:"country"`
	InterviewType      string         `json:"interview_type"`
	Stage              string         `json:"stage"`
	JobDescription     string         `gorm:"type:text" json:"job_description"`
	CvURL              string         `json:"cv_url"`
	UserID             string         `gorm:"type:uuid" json:"user_id"`
	User               User           `gorm:"foreignKey:UserID" json:"user"`
	Mode               string         `json:"mode"`
	StartedAt          *time.Time     `json:"started_at"`
	EndedAt            *time.Time     `json:"ended_at"`
	Duration           *int           `json:"duration"`
	TotalQuestions     *int           `json:"total_questions"`
	CompletedQuestions *int           `json:"completed_questions"`
	RecordingURL       string         `json:"recording_url"`
	Questions          []Question     `gorm:"foreignKey:InterviewID" json:"questions"`
	Responses          []Response     `gorm:"foreignKey:InterviewID" json:"responses"`
	Feedback           *Feedback      `gorm:"foreignKey:InterviewID" json:"feedback"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}
