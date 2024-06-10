package models

import "database/sql"

type User struct {
	Id       string
	Username sql.NullString
	Password sql.NullString
}
