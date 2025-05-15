package request

type GenerateLivekitTokenRequest struct {
	Room     string `json:"room" binding:"required" example:"interview-id-uuid"`
	Identity string `json:"identity" binding:"required" example:"userId-uuid"`
}
