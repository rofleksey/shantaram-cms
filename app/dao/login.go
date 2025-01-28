package dao

type LoginRequest struct {
	Pass string `json:"pass"`
}

type LoginResponse struct {
	ID       int64    `json:"id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	Token    string   `json:"token"`
}
