package dao

type NewOrderItemRequest struct {
	ID     string `json:"id"`
	Amount int    `json:"amount"`
}

type NewOrderRequest struct {
	CaptchaID     string                `json:"captchaId"`
	CaptchaAnswer string                `json:"captchaAnswer"`
	Name          string                `json:"name"`
	Phone         string                `json:"phone"`
	Comment       string                `json:"comment"`
	Items         []NewOrderItemRequest `json:"items"`
}
