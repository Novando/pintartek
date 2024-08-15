package user

type RegisterRequest struct {
	FullName        string `json:"fullName" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

type RegisterResponse struct {
	PrivateKey string `json:"privateKey"`
}
