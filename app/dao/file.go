package dao

type NewFileRequest struct {
	Title string `json:"title"`
	Name  string `json:"name"`
}

type FileStatsResponse struct {
	Total uint64 `json:"total"`
	Free  uint64 `json:"free"`
}
