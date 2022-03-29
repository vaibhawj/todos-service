package model

type ErrorResponse struct {
	Error string `json:"error"`
}

type Todo struct {
	Id   string `json:"id" bson:"_id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type TodoRequest struct {
	Text string `json:"text" binding:"required"`
	Done bool   `json:"done"`
}

type NotFoundError struct {
}

func (err NotFoundError) Error() string {
	return "Todo with specified id not found"
}
