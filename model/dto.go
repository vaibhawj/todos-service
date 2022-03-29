package model

type ErrorResponse struct {
	Error string `json:"error"`
}

type Todo struct {
	Id   string `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type TodoRequest struct {
	Text string `json:"text" binding:"required"`
	Done bool   `json:"done"`
}
