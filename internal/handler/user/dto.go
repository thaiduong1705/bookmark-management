package user

import "bookmark-management/internal/handler/response"

// RegisterRequest is the request body for registering a user.
type RegisterRequest struct {
	DisplayName string `json:"display_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required"`
	Username    string `json:"username" binding:"required"`
}

// RegisterResponse is the body returned after a user registers.
type RegisterResponse struct {
	response.Message
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
