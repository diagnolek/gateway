package http

type CreateUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	LastName string `json:"lastname"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	Id uint `json:"id"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

type GetUserResponse struct {
	Name     string `json:"name"`
	LastName string `json:"lastname"`
	Email    string `json:"email"`
}

type ErrorOccurredModel struct {
	Message string `json:"message"`
}
