package model

import (
	"gorm.io/gorm"
	"time"
)

type Feedback struct {
	ID                 string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OverallScore       float64        `json:"overall_score"`
	Strengths          []string       `gorm:"type:text[]" json:"strengths"`
	AreasToImprove     []string       `gorm:"type:text[]" json:"areas_to_improve"`
	Summary            string         `gorm:"type:text" json:"summary"`
	ContentScore       *float64       `json:"content_score"`
	CommunicationScore *float64       `json:"communication_score"`
	BehaviorScore      *float64       `json:"behavior_score"`
	ConfidenceScore    *float64       `json:"confidence_score"`
	ClarityScore       *float64       `json:"clarity_score"`
	RelevanceScore     *float64       `json:"relevance_score"`
	InterviewID        string         `gorm:"type:uuid;uniqueIndex" json:"interview_id"`
	Interview          Interview      `gorm:"foreignKey:InterviewID" json:"interview"`
	UserID             string         `gorm:"type:uuid;index" json:"user_id"`
	User               User           `gorm:"foreignKey:UserID" json:"user"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}
