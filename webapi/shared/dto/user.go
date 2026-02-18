package dto

type UserDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	// #nosec G117
	Password string `json:"password"`
}
