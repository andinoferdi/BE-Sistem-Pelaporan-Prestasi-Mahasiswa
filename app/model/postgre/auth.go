package model

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginUserResponse struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	FullName    string   `json:"fullName"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

type LoginResponse struct {
	Status string `json:"status"`
	Data   struct {
		Token        string            `json:"token"`
		RefreshToken string            `json:"refreshToken"`
		User         LoginUserResponse `json:"user"`
	} `json:"data"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type RefreshTokenResponse struct {
	Status string `json:"status"`
	Data   struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refreshToken"`
	} `json:"data"`
}

type GetProfileResponse struct {
	Status string `json:"status"`
	Data   struct {
		UserID      string   `json:"user_id"`
		Username    string   `json:"username"`
		Email       string   `json:"email"`
		FullName    string   `json:"full_name"`
		RoleID      string   `json:"role_id"`
		Role        string   `json:"role"`
		Permissions []string `json:"permissions"`
	} `json:"data"`
}
