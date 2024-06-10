package utils

import (
	"database/sql"
	"math/rand"
	"strings"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func NewNullString(s string) sql.NullString {
	if len(strings.TrimSpace(s)) == 0 {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
