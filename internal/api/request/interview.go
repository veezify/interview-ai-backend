package request

import "time"

type CreateInterviewRequest struct {
	Level                string   `json:"level"`
	JobType              string   `json:"job_type"`
	ProgrammingLanguages []string `json:"programming_languages"`
	InterviewLanguage    string   `json:"interview_language"`
	Country              string   `json:"country"`
	InterviewType        string   `json:"interview_type"`
	Stage                string   `json:"stage"`
	JobDescription       string   `json:"job_description"`
	CvURL                string   `json:"cv_url"`
	UserID               string   `json:"user_id"`
	Mode                 string   `json:"mode"`
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
