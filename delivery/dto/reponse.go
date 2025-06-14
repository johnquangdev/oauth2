package dto

type UserResponse struct {
	Gmail   string `json:"gmail" validate:"omitempty"`
	Name    string `json:"name" validate:"omitempty"`
	Picture string `json:"picture"`
}
