package dto

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type OrderResponse struct {
	ID      string  `json:"id"`
	UserID  string  `json:"user_id"`
	Product string  `json:"product"`
	Amount  float64 `json:"amount"`
	Status  string  `json:"status"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ListResponse struct {
	Items interface{} `json:"items"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}