package models

import "database/sql"

type User struct {
	Id       string
	Username string
	Password sql.NullString
}
