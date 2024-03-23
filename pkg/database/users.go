package database

import (
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Alias    string `json:"alias"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func (u *User) BeforeSave() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	u.Alias = html.EscapeString(strings.TrimSpace(u.Alias))
	return nil
}
