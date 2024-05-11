package dto

type UserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
	Email    string `json:"email" validate:"required,email"`
}

type UserUpdateRequest struct {
	Username string `json:"username" valid:"Required"`
	Email    string `json:"email" valid:"Required,email"`
}

type UserResponse struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UpdatePasswordRequest struct {
	Password string `json:"password" validate:"required,min=6"`
}
