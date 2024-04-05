package database

import (
	"errors"
	"strings"
)

var tokens = make(map[string]User)

func SaveToken(token string, user User) error {
	if _, ok := tokens[token]; ok {
		return errors.New("Token already exists")
	}
	tokens[token] = user
	return nil
}

func DeleteToken(token string) error {
	_, ok := tokens[token]
	if ok {
		delete(tokens, token)
		return nil
	}
	return errors.New("Token not found")
}

func GetUserFromToken(tokenHeader string) (User, error) {
	token := strings.Replace(tokenHeader, "Bearer ", "", 1)
	user, ok := tokens[token]
	if ok {
		return user, nil
	}
	return User{}, errors.New("Token not found")
}
