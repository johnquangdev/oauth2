package dto

type LoginOauth2 struct {
	Code string `json:"code" validate:"required"`
}
