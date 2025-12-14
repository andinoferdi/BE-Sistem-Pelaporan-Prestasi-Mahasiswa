package model

// #1 proses: struct untuk request login, butuh username dan password
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// #2 proses: struct untuk data user yang dikembalikan saat login
type LoginUserResponse struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	FullName    string   `json:"fullName"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

// #3 proses: struct response untuk login, berisi token, refresh token, dan data user
type LoginResponse struct {
	Status string `json:"status"`
	Data   struct {
		Token        string            `json:"token"`
		RefreshToken string            `json:"refreshToken"`
		User         LoginUserResponse `json:"user"`
	} `json:"data"`
}

// #4 proses: struct untuk request refresh token, kirim refresh token lama
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

// #5 proses: struct response untuk refresh token, dapatkan token dan refresh token baru
type RefreshTokenResponse struct {
	Status string `json:"status"`
	Data   struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refreshToken"`
	} `json:"data"`
}

// #6 proses: struct response untuk get profile user yang sedang login, lengkap dengan role dan permissions
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
