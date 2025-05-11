package service

import (
	"errors"
	"github.com/veezify/interview-ai-backend/internal/domain/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InterviewService gestionează operațiunile legate de interviuri
type InterviewService struct {
	DB *gorm.DB
}

// CreateInterviewRequest este structura pentru datele de intrare inițiale ale interviului
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

// UpdateInterviewRequest este structura pentru actualizarea datelor interviului
type UpdateInterviewRequest struct {
	ID                 string     `json:"id"`
	EndedAt            *time.Time `json:"ended_at,omitempty"`
	Duration           *int       `json:"duration,omitempty"`
	TotalQuestions     *int       `json:"total_questions,omitempty"`
	CompletedQuestions *int       `json:"completed_questions,omitempty"`
	RecordingURL       *string    `json:"recording_url,omitempty"`
	// Câmpuri opționale pentru actualizare
	Level                *string  `json:"level,omitempty"`
	JobType              *string  `json:"job_type,omitempty"`
	ProgrammingLanguages []string `json:"programming_languages,omitempty"`
	InterviewLanguage    *string  `json:"interview_language,omitempty"`
	Country              *string  `json:"country,omitempty"`
	InterviewType        *string  `json:"interview_type,omitempty"`
	Stage                *string  `json:"stage,omitempty"`
	JobDescription       *string  `json:"job_description,omitempty"`
	CvURL                *string  `json:"cv_url,omitempty"`
	Mode                 *string  `json:"mode,omitempty"`
}

// CreateInterview creează un nou interviu cu datele inițiale
func (s *InterviewService) CreateInterview(req CreateInterviewRequest) (*model.Interview, error) {
	// Validare date de intrare
	if err := s.validateInterviewRequest(req); err != nil {
		return nil, err
	}

	// Verificare dacă utilizatorul există
	var user model.User
	if err := s.DB.Where("id = ?", req.UserID).First(&user).Error; err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Creare interviu cu câmpurile inițiale
	interview := &model.Interview{
		ID:                   uuid.New().String(),
		Level:                req.Level,
		JobType:              req.JobType,
		ProgrammingLanguages: req.ProgrammingLanguages,
		InterviewLanguage:    req.InterviewLanguage,
		Country:              req.Country,
		InterviewType:        req.InterviewType,
		Stage:                req.Stage,
		JobDescription:       req.JobDescription,
		CvURL:                req.CvURL,
		UserID:               req.UserID,
		Mode:                 req.Mode,
		StartedAt:            nil, // Va fi setat la StartInterview
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	// Salvare interviu
	if err := s.DB.Create(interview).Error; err != nil {
		return nil, err
	}

	return interview, nil
}

// UpdateInterview actualizează câmpurile unui interviu existent
func (s *InterviewService) UpdateInterview(req UpdateInterviewRequest) (*model.Interview, error) {
	// Verifică dacă interviul există
	var interview model.Interview
	if err := s.DB.Where("id = ?", req.ID).First(&interview).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("interview not found")
		}
		return nil, err
	}

	// Pregătește un map cu câmpurile ce trebuie actualizate
	updates := make(map[string]interface{})
	updates["updated_at"] = time.Now()

	// Adaugă doar câmpurile non-nil în map
	if req.EndedAt != nil {
		updates["ended_at"] = req.EndedAt
	}
	if req.Duration != nil {
		updates["duration"] = req.Duration
	}
	if req.TotalQuestions != nil {
		updates["total_questions"] = req.TotalQuestions
	}
	if req.CompletedQuestions != nil {
		updates["completed_questions"] = req.CompletedQuestions
	}
	if req.RecordingURL != nil {
		updates["recording_url"] = req.RecordingURL
	}
	if req.Level != nil {
		updates["level"] = req.Level
	}
	if req.JobType != nil {
		updates["job_type"] = req.JobType
	}
	if len(req.ProgrammingLanguages) > 0 {
		updates["programming_languages"] = req.ProgrammingLanguages
	}
	if req.InterviewLanguage != nil {
		updates["interview_language"] = req.InterviewLanguage
	}
	if req.Country != nil {
		updates["country"] = req.Country
	}
	if req.InterviewType != nil {
		updates["interview_type"] = req.InterviewType
	}
	if req.Stage != nil {
		updates["stage"] = req.Stage
	}
	if req.JobDescription != nil {
		updates["job_description"] = req.JobDescription
	}
	if req.CvURL != nil {
		updates["cv_url"] = req.CvURL
	}
	if req.Mode != nil {
		updates["mode"] = req.Mode
	}

	// Dacă nu avem ce actualiza, returnăm direct interviul
	if len(updates) <= 1 { // doar updated_at
		return &interview, nil
	}

	// Actualizează interviul
	if err := s.DB.Model(&interview).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Reîncarcă interviul cu toate relațiile
	if err := s.DB.Preload("User").
		Preload("Questions").
		Preload("Responses").
		Preload("Feedback").
		Where("id = ?", req.ID).First(&interview).Error; err != nil {
		return nil, err
	}

	return &interview, nil
}

// StartInterview marchează un interviu ca fiind început
func (s *InterviewService) StartInterview(id string) error {
	now := time.Now()
	return s.DB.Model(&model.Interview{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"started_at": now,
			"updated_at": now,
		}).Error
}

// EndInterview marchează un interviu ca fiind finalizat
func (s *InterviewService) EndInterview(id string) error {
	now := time.Now()

	// Găsește interviul
	var interview model.Interview
	if err := s.DB.Where("id = ?", id).First(&interview).Error; err != nil {
		return err
	}

	// Calculează durata în secunde
	var duration int
	if interview.StartedAt != nil {
		duration = int(now.Sub(*interview.StartedAt).Seconds())
	}

	// Actualizează interviul
	return s.DB.Model(&model.Interview{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"ended_at":   now,
			"duration":   duration,
			"updated_at": now,
		}).Error
}

// GetInterviewByID obține un interviu după ID
func (s *InterviewService) GetInterviewByID(id string) (*model.Interview, error) {
	var interview model.Interview
	err := s.DB.Preload("User").
		Preload("Questions").
		Preload("Responses").
		Preload("Feedback").
		Where("id = ?", id).
		First(&interview).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("interview not found")
		}
		return nil, err
	}

	return &interview, nil
}

// GetInterviewsByUserID obține toate interviurile unui utilizator
func (s *InterviewService) GetInterviewsByUserID(userID string) ([]model.Interview, error) {
	var interviews []model.Interview
	err := s.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&interviews).Error

	if err != nil {
		return nil, err
	}

	return interviews, nil
}

// validateInterviewRequest validează datele de intrare pentru crearea unui interviu
func (s *InterviewService) validateInterviewRequest(req CreateInterviewRequest) error {
	// Validare câmpuri obligatorii
	if req.UserID == "" {
		return errors.New("user ID is required")
	}
	if req.Level == "" {
		return errors.New("level is required")
	}
	if req.JobType == "" {
		return errors.New("job type is required")
	}
	if req.InterviewLanguage == "" {
		return errors.New("interview language is required")
	}
	if req.Country == "" {
		return errors.New("country is required")
	}
	if req.InterviewType == "" {
		return errors.New("interview type is required")
	}
	if req.Stage == "" {
		return errors.New("stage is required")
	}
	if req.Mode == "" {
		return errors.New("mode is required")
	}

	// Validare level
	validLevels := []string{"JUNIOR", "MID_LEVEL", "SENIOR", "LEAD", "ARCHITECT"}
	if !contains(validLevels, req.Level) {
		return errors.New("invalid experience level")
	}

	// Validare mod interviu
	validModes := []string{"ONLINE", "IN_PERSON", "PHONE"}
	if !contains(validModes, req.Mode) {
		return errors.New("invalid interview mode")
	}

	// Validare tip interviu
	validTypes := []string{"TECHNICAL", "BEHAVIORAL", "MIXED"}
	if !contains(validTypes, req.InterviewType) {
		return errors.New("invalid interview type")
	}

	// Validare etapă interviu
	validStages := []string{"SCREENING", "TECHNICAL", "FINAL", "HR"}
	if !contains(validStages, req.Stage) {
		return errors.New("invalid interview stage")
	}

	// Validare limba interviului
	validLanguages := []string{"ENGLISH", "ROMANIAN", "RUSSIAN", "FRENCH", "GERMAN"}
	if !contains(validLanguages, req.InterviewLanguage) {
		return errors.New("invalid interview language")
	}

	return nil
}

// contains verifică dacă un slice conține un anumit element
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
