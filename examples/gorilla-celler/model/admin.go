package model

// Admin example
type Admin struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"admin name"`
}

// Message example
type Message struct {
	Message string `json:"message" example:"message"`
}
