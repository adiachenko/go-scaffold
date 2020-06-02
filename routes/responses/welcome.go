package responses

type WelcomeResponse struct {
	Status int    `json:"status"`
	Data   string `json:"data"`
}

func NewWelcomeResponse(msg string) WelcomeResponse {
	return WelcomeResponse{
		Status: 200,
		Data:   msg,
	}
}
