package validation

type CreateUser struct {
	Name     string `json:"name" validate:"required,max=50" example:"fake name"`
	Email    string `json:"email" validate:"required,email,max=50" example:"fake@example.com"`
	Password string `json:"password" validate:"required,min=8,max=20,password" example:"password1"`
	Role     string `json:"role" validate:"required,oneof=user admin,max=50" example:"user"`
}

type UpdateUser struct {
	Name     string `json:"name,omitempty" validate:"omitempty,max=50" example:"fake name"`
	Email    string `json:"email" validate:"omitempty,email,max=50" example:"fake@example.com"`
	Password string `json:"password,omitempty" validate:"omitempty,min=8,max=20,password" example:"password1"`
}

type UpdatePassOrVerify struct {
	Password      string `json:"password,omitempty" validate:"omitempty,min=8,max=20,password" example:"password1"`
	VerifiedEmail bool   `json:"verified_email" swaggerignore:"true" validate:"omitempty,boolean"`
}

type QueryUser struct {
	Page   int    `validate:"omitempty,number,max=50"`
	Limit  int    `validate:"omitempty,number,max=50"`
	Search string `validate:"omitempty,max=50"`
}
