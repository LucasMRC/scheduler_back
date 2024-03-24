package database

import "errors"

var tokens []string

func SaveToken(token string) error {
	tokens = append(tokens, token)
	return nil
}

func DeleteToken(token string) error {
	for i, t := range tokens {
		if t == token {
			tokens = append(tokens[:i], tokens[i+1:]...)
			return nil
		}
	}
	return errors.New("Token not found")
}
