package dao

type NoDataResponse struct {
	Error bool   `json:"error"`
	Msg   string `json:"msg"`
}

type SuccessResponse[T any] struct {
	Error bool `json:"error"`
	Data  T    `json:"data"`
}

type IdDao[T any] struct {
	ID T `json:"id"`
}
