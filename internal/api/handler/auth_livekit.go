package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/livekit/protocol/auth"
	"github.com/veezify/interview-ai-backend/internal/api/request"
	"github.com/veezify/interview-ai-backend/internal/api/response"
	"net/http"
	"time"
)

type AuthLivekitHandler struct {
	apiKey    string
	apiSecret string
}

func NewAuthLivekitHandler() *AuthLivekitHandler {
	return &AuthLivekitHandler{
		apiKey:    "APIBWQchjyts3d6",
		apiSecret: "gGXVGRBL3edk2CPcvwE2rk3vqxekkUfGCvCWlAwshsuA",
	}
}

// GenerateToken godoc
// @Summary Generate LiveKit token
// @Description Generate a LiveKit JWT token for connecting to a room
// @Tags livekit
// @Accept json
// @Produce json
// @Param request body request.GenerateLivekitTokenRequest true "Room and identity information"
// @Success 200 {object} response.GenerateLivekitTokenResponse "Returns a JWT token"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /livekit/token [post]
// @Security ApiKeyAuth
func (h *AuthLivekitHandler) GenerateLivekitToken(c *gin.Context) {
	var req request.GenerateLivekitTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Room == "" || req.Identity == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room and identity are required"})
		return
	}

	at := auth.NewAccessToken(h.apiKey, h.apiSecret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     req.Room,
	}

	at.AddGrant(grant).
		SetIdentity(req.Identity).
		SetValidFor(time.Hour)

	token, err := at.ToJWT()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, response.GenerateLivekitTokenResponse{
		Token: token,
	})
}
