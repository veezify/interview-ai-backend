package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/veezify/interview-ai-backend/internal/api/request"
	"github.com/veezify/interview-ai-backend/internal/api/response"
	"github.com/veezify/interview-ai-backend/internal/domain/model"
	"github.com/veezify/interview-ai-backend/internal/service"
	"net/http"
	"time"
)

type InterviewHandler struct {
	interviewService *service.InterviewService
	userService      *service.UserService
}

func NewInterviewHandler(interviewService *service.InterviewService, userService *service.UserService) *InterviewHandler {
	return &InterviewHandler{
		interviewService: interviewService,
		userService:      userService,
	}
}

// GetSelectOptions godoc
// @Summary Get selection options
// @Description Returns lists of options for dropdowns including countries, levels, languages, and programming languages
// @Tags options
// @Accept json
// @Produce json
// @Success 200 {object} response.SelectOptionsResponse "List of options for dropdowns"
// @Router /interviews/options [get]
func (h *InterviewHandler) GetSelectOptions(c *gin.Context) {
	options := response.SelectOptionsResponse{
		Countries: []response.SelectOption{
			{Value: "md", Label: "Moldova"},
			{Value: "ro", Label: "România"},
			{Value: "us", Label: "USA"},
		},
		Levels: []response.SelectOption{
			{Value: "junior", Label: "Junior"},
			{Value: "middle", Label: "Middle"},
			{Value: "senior", Label: "Senior"},
			{Value: "lead", Label: "Team Lead"},
		},

		JobTypes: []response.SelectOption{
			{Value: "frontend", Label: "Frontend Developer"},
			{Value: "backend", Label: "Backend Developer"},
			{Value: "fullstack", Label: "Fullstack Developer"},
			{Value: "mobile", Label: "Mobile Developer"},
			{Value: "devops", Label: "DevOps Engineer"},
			{Value: "qa", Label: "QA Engineer"},
		},
		Languages: []response.SelectOption{
			{Value: "en", Label: "English"},
		},
		ProgrammingLanguages: []response.SelectOption{
			{Value: "javascript", Label: "JavaScript"},
			{Value: "typescript", Label: "TypeScript"},
			{Value: "python", Label: "Python"},
			{Value: "java", Label: "Java"},
			{Value: "csharp", Label: "C#"},
			{Value: "cpp", Label: "C++"},
			{Value: "ruby", Label: "Ruby"},
			{Value: "go", Label: "Go"},
			{Value: "php", Label: "PHP"},
			{Value: "swift", Label: "Swift"},
			{Value: "kotlin", Label: "Kotlin"},
		},
		InterviewTypes: []response.SelectOption{
			{Value: "technical", Label: "Technical (coding, architecture, system design)"},
			{Value: "behavioral", Label: "Behavioral (soft skills, team fit)"},
			{Value: "mixed", Label: "Mixed (combination of both)"},
		},
		Stages: []response.SelectOption{
			{Value: "screening", Label: "Initial Screening<"},
			{Value: "technical", Label: "Technical Interview"},
			{Value: "manager", Label: "Manager Interview"},
			{Value: "final", Label: "Final Round"},
			{Value: "onsite", Label: "Onsite Interview"},
		},
	}

	c.JSON(http.StatusOK, options)
}

// CreateInterview godoc
// @Summary Create a new interview
// @Description Create a new interview with the specified parameters
// @Tags interviews
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param level formData string true "Experience level (e.g., JUNIOR, MIDDLE, SENIOR)"
// @Param jobType formData string true "Job type (e.g., BACKEND, FRONTEND)"
// @Param programmingLanguages formData string true "Programming languages array as JSON string"
// @Param interviewLanguage formData string true "Interview language (e.g., ENGLISH)"
// @Param country formData string true "Country (e.g., US, UK)"
// @Param interviewType formData string true "Interview type (e.g., TECHNICAL, BEHAVIORAL)"
// @Param stage formData string true "Interview stage (e.g., SCREENING, FINAL)"
// @Param jobDescription formData string true "Job description text"
// @Param mode formData string true "Interview mode (e.g., TEXT_WITH_VOICE)"
// @Param cv formData file false "CV/Resume file (PDF or DOCX)"
// @Success 201 {object} response.CreateInterviewResponse "Successfully created interview"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /interviews [post]
func (h *InterviewHandler) CreateInterview(c *gin.Context) {

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var cvURL string
	file, err := c.FormFile("cv")
	if err == nil && file != nil {
		cvURL, err = h.interviewService.UploadCV(file, userID.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload CV: " + err.Error()})
			return
		}
	}

	var req request.CreateInterviewRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data: " + err.Error()})
		return
	}

	progLangsJSON := c.PostForm("programmingLanguages")
	if progLangsJSON == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Programming languages are required"})
		return
	}

	interview := model.Interview{
		Level:                req.Level,
		JobType:              req.JobType,
		ProgrammingLanguages: req.ProgrammingLanguages,
		InterviewLanguage:    req.InterviewLanguage,
		Country:              req.Country,
		InterviewType:        req.InterviewType,
		Stage:                req.Stage,
		JobDescription:       req.JobDescription,
		CvURL:                cvURL,
		UserID:               userID.(string),
		Mode:                 req.Mode,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	createdInterview, err := h.interviewService.CreateInterview(interview)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create interview: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response.CreateInterviewResponse{
		InterviewID: createdInterview.ID,
	})
}

// GetInterview godoc
// @Summary Get interview details
// @Description Get detailed information about an interview
// @Tags interviews
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Interview ID"
// @Success 200 {object} model.Interview "Interview details"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Interview not found"
// @Router /interviews/{id} [get]
func (h *InterviewHandler) GetInterview(c *gin.Context) {
	// Obțineți ID-ul utilizatorului autentificat
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Obțineți ID-ul interviului din URL
	interviewID := c.Param("id")
	if interviewID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Interview ID is required"})
		return
	}

	// Preluați interviul din baza de date
	interview, err := h.interviewService.GetInterviewByID(interviewID, userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Interview not found"})
		return
	}

	// Returnați interviul
	c.JSON(http.StatusOK, interview)
}

// ListInterviews godoc
// @Summary List user interviews
// @Description Get a list of all interviews for the current user
// @Tags interviews
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} model.Interview "List of interviews"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /interviews [get]
func (h *InterviewHandler) ListInterviews(c *gin.Context) {
	// Obțineți ID-ul utilizatorului autentificat
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Preluați toate interviurile utilizatorului
	interviews, err := h.interviewService.ListInterviewsByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch interviews"})
		return
	}

	// Returnați lista de interviuri
	c.JSON(http.StatusOK, interviews)
}
