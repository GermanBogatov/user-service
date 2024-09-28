package model

type SignUpRequest struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	SignInRequest
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	Email       string  `json:"email"`
	CreatedDate string  `json:"createdDate"`
	UpdatedDate *string `json:"updatedDate"`
	JWT         *JWT    `json:"jwt"`
}

type JWT struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}
