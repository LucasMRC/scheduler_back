package users

import (
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

func (u *User) BeforeSave() error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.Password = string(hashedPassword)
    u.Username = html.EscapeString(strings.TrimSpace(u.Username))
    return nil
}
