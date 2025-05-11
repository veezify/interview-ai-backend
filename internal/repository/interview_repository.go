package repository

import (
	"errors"
	"time"

	"github.com/veezify/interview-ai-backend/internal/domain/model"
	"gorm.io/gorm"
)

type InterviewRepository struct {
	db *gorm.DB
}

func NewInterviewRepository(db *gorm.DB) *InterviewRepository {
	return &InterviewRepository{db: db}
}

func (r *InterviewRepository) Create(interview *model.Interview) error {
	return r.db.Create(interview).Error
}

func (r *InterviewRepository) FindByID(id string) (*model.Interview, error) {
	var interview model.Interview
	result := r.db.Where("id = ?", id).First(&interview)
	if result.Error != nil {
		return nil, result.Error
	}
	return &interview, nil
}

func (r *InterviewRepository) FindByIDWithRelations(id string) (*model.Interview, error) {
	var interview model.Interview
	result := r.db.
		Preload("User").
		Preload("Questions").
		Preload("Responses").
		Preload("Feedback").
		Where("id = ?", id).
		First(&interview)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("interview not found")
		}
		return nil, result.Error
	}

	return &interview, nil
}

func (r *InterviewRepository) FindByUserID(userID string) ([]model.Interview, error) {
	var interviews []model.Interview
	result := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&interviews)

	if result.Error != nil {
		return nil, result.Error
	}

	return interviews, nil
}

func (r *InterviewRepository) Update(interview *model.Interview) error {
	return r.db.Save(interview).Error
}

func (r *InterviewRepository) UpdateFields(id string, fields map[string]interface{}) error {
	// Adaugă câmpul updated_at
	fields["updated_at"] = time.Now()

	return r.db.Model(&model.Interview{}).
		Where("id = ?", id).
		Updates(fields).Error
}

func (r *InterviewRepository) Delete(id string) error {
	return r.db.Delete(&model.Interview{}, id).Error
}

func (r *InterviewRepository) Exists(id string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Interview{}).
		Where("id = ?", id).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *InterviewRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
