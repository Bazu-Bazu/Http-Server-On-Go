package user

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(u User) (User, error) {
	const query = `
		INSERT INTO users (nickname, name, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id`

	err := r.db.QueryRow(query, u.Nickname, u.Name, u.PasswordHash).Scan(&u.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "users_nickname_lower_uidx" {
			return User{}, ErrNicknameTaken
		}

		return User{}, fmt.Errorf("Insert user: %w", err)
	}

	return u, nil
}

func (r *PostgresRepository) GetByID(id int64) (User, error) {
	const query = `
		SELECT id, nickname, name, password_hash
		FROM users
		WHERE id = $1`

	var u User
	err := r.db.QueryRow(query, id).Scan(&u.ID, &u.Nickname, &u.Name, &u.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrNotFound
	}

	if err != nil {
		return User{}, fmt.Errorf("Select user by id: %w", err)
	}

	return u, nil
}

func (r *PostgresRepository) Delete(id int64) error {
	result, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("Delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Get deleted rows count: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
