package dto

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" validate:"omitempty,min=2,max=100"`
	Email string `json:"email" validate:"omitempty,email"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
