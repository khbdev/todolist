package mysql

import (
	"database/sql"
	"todolist/internal/domain"
)

type UserRepo struct{
	db *sql.DB
}


func NewUserRepo(db *sql.DB) domain.UserRepository {
    return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(user *domain.User) error {
	query := `INSERT INTO users (username, email, password) VALUES (?, ?, ?)`
	result, err := r.db.Exec(query, user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(lastID)
	return nil
}


func (r *UserRepo) GetUserByEmail(email string) (*domain.User, error) {
    query := `SELECT id, username, email, password FROM users WHERE email = ?`
    row := r.db.QueryRow(query, email)

    var user domain.User
    err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
    if err != nil {
        return nil, err
    }

    return &user, nil
}