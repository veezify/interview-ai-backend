package service

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/veezify/interview-ai-backend/internal/domain/model"
	"github.com/veezify/interview-ai-backend/internal/repository"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type InterviewService struct {
	interviewRepo repository.InterviewRepository
	uploadDir     string
}

// NewInterviewService
func NewInterviewService(interviewRepo repository.InterviewRepository, uploadDir string) *InterviewService {

	if uploadDir == "" {
		uploadDir = "uploads/cv"
	}

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	return &InterviewService{
		interviewRepo: interviewRepo,
		uploadDir:     uploadDir,
	}
}

// CreateInterview
func (s *InterviewService) CreateInterview(interview model.Interview) (*model.Interview, error) {
	// Validare simplă
	if interview.UserID == "" {
		return nil, errors.New("user ID is required")
	}

	interview.Level = interview.Level
	interview.JobType = interview.JobType
	interview.InterviewLanguage = interview.InterviewLanguage
	interview.Country = interview.Country
	interview.InterviewType = interview.InterviewType
	interview.Stage = interview.Stage

	if interview.Mode == "" {
		interview.Mode = "voice"
	}

	return s.interviewRepo.CreateInterview(interview)
}

// UploadCV încarcă fișierul CV și returnează calea relativă
func (s *InterviewService) UploadCV(file *multipart.FileHeader, userID string) (string, error) {
	if file == nil {
		return "", nil // Niciun fișier încărcat, nu este o eroare
	}

	// Verificați tipul de fișier (permiteți doar PDF, DOCX, DOC)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".pdf" && ext != ".docx" && ext != ".doc" {
		return "", errors.New("only PDF, DOCX, and DOC files are allowed")
	}

	// Generați un nume unic pentru fișier
	fileName := s.generateUniqueFileName(file.Filename, userID)
	filePath := filepath.Join(s.uploadDir, fileName)

	// Deschideți fișierul sursă
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Creați fișierul destinație
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copiați conținutul
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	// Returnați calea relativă
	relativePath := fmt.Sprintf("/uploads/cv/%s", fileName)
	return relativePath, nil
}

func (s *InterviewService) generateUniqueFileName(originalName, userID string) string {
	// Extrageți extensia
	ext := filepath.Ext(originalName)

	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	randomString := hex.EncodeToString(randomBytes)

	// Construiți numele fișierului: timestamp_userid_random.extension
	timestamp := time.Now().Unix()
	fileName := fmt.Sprintf("%d_%s_%s%s", timestamp, userID, randomString, ext)

	return fileName
}

// GetInterviewByID preia un interviu după ID
func (s *InterviewService) GetInterviewByID(id string, userID string) (*model.Interview, error) {
	interview, err := s.interviewRepo.GetInterviewByID(id)
	if err != nil {
		return nil, err
	}

	// Verificați dacă utilizatorul are acces la acest interviu
	if interview.UserID != userID {
		return nil, errors.New("unauthorized access to interview")
	}

	return interview, nil
}

// ListInterviewsByUserID listează toate interviurile unui utilizator
func (s *InterviewService) ListInterviewsByUserID(userID string) ([]model.Interview, error) {
	return s.interviewRepo.ListInterviewsByUserID(userID)
}
