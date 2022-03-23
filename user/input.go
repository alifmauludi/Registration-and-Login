package user

import "time"

type RegisterUserInput struct {
	Name         string `json:"name" binding:"required" validate:"required"`
	Email        string `json:"email" binding:"required,email" validate:"required,email"`
	Password     string `json:"password" binding:"required,min=8" validate:"required,min=8"`
	DateRegister time.Time
}

type LoginInput struct {
	Email    string `json:"email" binding:"required" validate:"required"`
	Password string `json:"password" binding:"required" validate:"required"`
}
