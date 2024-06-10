package repo

import (
	"database/sql"
	"hanzotg/internal/app/models"
	"hanzotg/internal/app/utils"
)

type UserRepo struct {
	db *sql.DB
}

func NewUser(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (u *UserRepo) Create(id, username string) error {
	_, err := u.db.Exec(
		`INSERT INTO public."users" (id, username, password) VALUES ($1, $2, null)`,
		id,
		utils.NewNullString(username),
	)
	return err
}

func (u *UserRepo) GetById(id string) (*models.User, error) {
	var user models.User
	err := u.db.QueryRow(
		`SELECT id, username, password FROM public."users" WHERE id = $1`,
		id,
	).Scan(&user.Id, &user.Username, &user.Password)
	return &user, err
}

func (u *UserRepo) InsertPassword(userId, password string) error {
	_, err := u.db.Exec(
		`UPDATE public."users" SET password = $1 WHERE id = $2`,
		password,
		userId,
	)
	return err
}
