package dao

type CaptchaVerifyRequest struct {
	ID     string `json:"id"`
	Answer string `json:"answer"`
}

type CaptchaGenerateResponse struct {
	ID    string `json:"id"`
	Image string `json:"image"`
}
