package model

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		User  User   `json:"user"`
		Token string `json:"token"`
	} `json:"data"`
}

type GetProfileResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		UserID   string `json:"user_id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		FullName string `json:"full_name"`
		RoleID   string `json:"role_id"`
	} `json:"data"`
}
