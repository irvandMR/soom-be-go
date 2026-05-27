package domain

type RegisterRequest struct {
	Username string  `json:"username"  binding:"required"`
	Email    string  `json:"email"     binding:"required,email"`
	Password string  `json:"password"  binding:"required,min=6"`
	TenantId *string `json:"tenant_id"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Role         string `json:"role"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
