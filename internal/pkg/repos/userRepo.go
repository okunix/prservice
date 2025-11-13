package repos

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/okunix/prservice/internal/pkg/models"
	"github.com/okunix/prservice/internal/pkg/user"
)

type UserRepoImpl struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) user.Repo {
	return &UserRepoImpl{db: db}
}

func (u *UserRepoImpl) GetUsersByTeamName(
	ctx context.Context,
	teamName string,
) ([]user.User, error) {
	q := `SELECT u.id, u.username, u.is_active FROM users u WHERE u.team_name = $1;`
	var users []user.User
	rows, err := u.db.QueryContext(ctx, q, teamName)
	if err != nil {
		return users, err
	}
	for rows.Next() {
		var user user.User
		if err := rows.Scan(&user.Id, &user.Name, &user.IsActive); err != nil {
			slog.Error(err.Error())
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

func (repo *UserRepoImpl) SetIsActive(
	ctx context.Context,
	userId string,
	isActive bool,
) (*user.User, error) {
	q := `UPDATE users SET is_active = $1 WHERE id = $2 RETURNING id, username, is_active, team_name;`
	var u user.User
	err := repo.db.QueryRowContext(ctx, q, isActive, userId).
		Scan(&u.Id, &u.Name, &u.IsActive, &u.TeamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}
