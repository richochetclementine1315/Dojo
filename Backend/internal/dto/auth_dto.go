package dto

// RegisterRequestrepresents the user registration request payload
type RegisterRequest struct {
	Email              string `json:"email" validate:"required,email"`
	Username           string `json:"username" validate:"required,min=3,max=50"`
	Password           string `json:"password" validate:"required,min=8"`
	LeetcodeUsername   string `json:"leetcode_username"`
	CodeforcesUsername string `json:"codeforces_username"`
	CodechefUsername   string `json:"codechef_username"`
	GFGUsername        string `json:"gfg_username"`
}

// LoginRequest represents the user login request payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// OAuthCallbackRequest represents the OAuth callback request payload
type OAuthCallbackRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state"`
}

// TokenResponse represents the JWT token response payload
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
