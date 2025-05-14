package repository

import (
	"github.com/veezify/interview-ai-backend/internal/domain/model"
	"gorm.io/gorm"
)

type InterviewRepository interface {
	CreateInterview(interview model.Interview) (*model.Interview, error)
	GetInterviewByID(id string) (*model.Interview, error)
	ListInterviewsByUserID(userID string) ([]model.Interview, error)
	UpdateInterview(interview model.Interview) (*model.Interview, error)
	DeleteInterview(id string) error
}

type interviewRepository struct {
	db *gorm.DB
}

// NewInterviewRepository creează o nouă instanță a repository-ului
func NewInterviewRepository(db *gorm.DB) InterviewRepository {
	return &interviewRepository{db: db}
}

// CreateInterview creează un nou interviu în baza de date
func (r *interviewRepository) CreateInterview(interview model.Interview) (*model.Interview, error) {
	if err := r.db.Create(&interview).Error; err != nil {
		return nil, err
	}
	return &interview, nil
}

// GetInterviewByID preia un interviu după ID
func (r *interviewRepository) GetInterviewByID(id string) (*model.Interview, error) {
	var interview model.Interview
	if err := r.db.First(&interview, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &interview, nil
}

// ListInterviewsByUserID listează toate interviurile unui utilizator
func (r *interviewRepository) ListInterviewsByUserID(userID string) ([]model.Interview, error) {
	var interviews []model.Interview
	if err := r.db.Where("user_id = ?", userID).Find(&interviews).Error; err != nil {
		return nil, err
	}
	return interviews, nil
}

// UpdateInterview actualizează un interviu
func (r *interviewRepository) UpdateInterview(interview model.Interview) (*model.Interview, error) {
	if err := r.db.Save(&interview).Error; err != nil {
		return nil, err
	}
	return &interview, nil
}

// DeleteInterview șterge un interviu (soft delete cu GORM)
func (r *interviewRepository) DeleteInterview(id string) error {
	return r.db.Delete(&model.Interview{}, "id = ?", id).Error
}
