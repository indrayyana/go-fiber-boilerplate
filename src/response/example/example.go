package example

type RegisterResponse struct {
	Code    int    `json:"code" example:"201"`
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Register successfully"`
	User    User   `json:"user"`
	Tokens  Tokens `json:"tokens"`
}

type LoginResponse struct {
	Code    int    `json:"code" example:"200"`
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Login successfully"`
	User    User   `json:"user"`
	Tokens  Tokens `json:"tokens"`
}

type GoogleLoginResponse struct {
	Code    int        `json:"code" example:"200"`
	Status  string     `json:"status" example:"success"`
	Message string     `json:"message" example:"Login successfully"`
	User    GoogleUser `json:"user"`
	Tokens  Tokens     `json:"tokens"`
}

type LogoutResponse struct {
	Code    int    `json:"code" example:"200"`
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Logout successfully"`
}

type RefreshTokenResponse struct {
	Code   int    `json:"code" example:"200"`
	Status string `json:"status" example:"success"`
	Tokens Tokens `json:"tokens"`
}

type ForgotPasswordResponse struct {
	Code    int    `json:"code" example:"200"`
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"A password reset link has been sent to your email address."`
}

type ResetPasswordResponse struct {
	Code    int    `json:"code" example:"200"`
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Update password successfully"`
}

type SendVerificationEmailResponse struct {
	Code    int    `json:"code" example:"200"`
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Please check your email for a link to verify your account"`
}

type VerifyEmailResponse struct {
	Code    int    `json:"code" example:"200"`
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Verify email successfully"`
}

type GetAllUserResponse struct {
	Code         int    `json:"code" example:"200"`
	Status       string `json:"status" example:"success"`
	Message      string `json:"message" example:"Get all users successfully"`
	Results      []User `json:"results"`
	Page         int    `json:"page" example:"1"`
	Limit        int    `json:"limit" example:"10"`
	TotalPages   int64  `json:"total_pages" example:"1"`
	TotalResults int64  `json:"total_results" example:"1"`
}

type GetUserResponse struct {
	Code    int    `json:"code" example:"200"`
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Get user successfully"`
	User    User   `json:"user"`
}

type CreateUserResponse struct {
	Code    int    `json:"code" example:"201"`
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Create user successfully"`
	User    User   `json:"user"`
}

type UpdateUserResponse struct {
	Code    int    `json:"code" example:"200"`
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Update user successfully"`
	User    User   `json:"user"`
}

type DeleteUserResponse struct {
	Code    int    `json:"code" example:"200"`
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Delete user successfully"`
}
