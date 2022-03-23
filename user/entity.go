package user

import "time"

type User struct {
	ID           int
	Name         string
	Email        string
	Password     string
	Avatar       string
	DateRegister time.Time
	DateUpdated  time.Time
}
