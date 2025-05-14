package request

import (
	"mime/multipart"
	"time"
)

type CreateInterviewRequest struct {
	Level                string                `json:"level" form:"level" binding:"required" example:"junior"`
	JobType              string                `json:"jobType" form:"jobType" binding:"required" example:"backend"`
	ProgrammingLanguages []string              `json:"programmingLanguages" form:"programmingLanguages" binding:"required" example:"[\"go\",\"java\"]"`
	InterviewLanguage    string                `json:"interviewLanguage" form:"interviewLanguage" binding:"required" example:"en"`
	Country              string                `json:"country" form:"country" binding:"required" example:"US"`
	InterviewType        string                `json:"interviewType" form:"interviewType" binding:"required" example:"technical"`
	Stage                string                `json:"stage" form:"stage" binding:"required" example:"screening"`
	JobDescription       string                `json:"jobDescription" form:"jobDescription" binding:"required" example:"We are looking for a backend developer..."`
	Mode                 string                `json:"mode" form:"mode" binding:"required" example:"voice"`
	CV                   *multipart.FileHeader `form:"cv" swaggerignore:"true"`
}

type UpdateInterviewRequest struct {
	ID                   string     `json:"id"`
	EndedAt              *time.Time `json:"ended_at,omitempty"`
	Duration             *int       `json:"duration,omitempty"`
	TotalQuestions       *int       `json:"total_questions,omitempty"`
	CompletedQuestions   *int       `json:"completed_questions,omitempty"`
	RecordingURL         *string    `json:"recording_url,omitempty"`
	Level                *string    `json:"level,omitempty"`
	JobType              *string    `json:"job_type,omitempty"`
	ProgrammingLanguages []string   `json:"programming_languages,omitempty"`
	InterviewLanguage    *string    `json:"interview_language,omitempty"`
	Country              *string    `json:"country,omitempty"`
	InterviewType        *string    `json:"interview_type,omitempty"`
	Stage                *string    `json:"stage,omitempty"`
	JobDescription       *string    `json:"job_description,omitempty"`
	CvURL                *string    `json:"cv_url,omitempty"`
	Mode                 *string    `json:"mode,omitempty"`
}
