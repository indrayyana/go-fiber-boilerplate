package example

type Unauthorized struct {
	Code    int    `json:"code" example:"401"`
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Please authenticate"`
}

type FailedLogin struct {
	Code    int    `json:"code" example:"401"`
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Invalid email or password"`
}

type FailedResetPassword struct {
	Code    int    `json:"code" example:"401"`
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Password reset failed"`
}

type FailedVerifyEmail struct {
	Code    int    `json:"code" example:"401"`
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Verify email failed"`
}

type Forbidden struct {
	Code    int    `json:"code" example:"403"`
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"You don't have permission to access this resource"`
}

type NotFound struct {
	Code    int    `json:"code" example:"404"`
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Not found"`
}

type DuplicateEmail struct {
	Code    int    `json:"code" example:"409"`
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Email already taken"`
}
