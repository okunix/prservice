package repos

import (
	"context"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/okunix/prservice/internal/pkg/models"
	"github.com/okunix/prservice/internal/pkg/team"
	"github.com/okunix/prservice/internal/pkg/user"
)

type TeamRepoImpl struct {
	db       *sqlx.DB
	userRepo user.Repo
}

func NewTeamRepo(db *sqlx.DB, userRepo user.Repo) team.Repo {
	return &TeamRepoImpl{
		db:       db,
		userRepo: userRepo,
	}
}

func (repo *TeamRepoImpl) AddTeam(ctx context.Context, t team.Team) (*team.Team, error) {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	insertTeamQuery := `INSERT INTO teams (name) VALUES ($1) RETURNING name;`
	_, err = tx.ExecContext(ctx, insertTeamQuery, t.Name)
	if err != nil {
		tx.Rollback()
		slog.Error(err.Error())
		return nil, err
	}

	addMemberQuery := `UPDATE users SET team_name=$2 WHERE id=$1;`
	for _, v := range t.Members {
		_, err := tx.ExecContext(ctx, addMemberQuery, v.Id, t.Name)
		if err != nil {
			tx.Rollback()
			slog.Error(err.Error())
			return nil, models.ErrNotAssigned
		}
	}

	if err := tx.Commit(); err != nil {
		slog.Error(err.Error())
		return nil, models.ErrTeamExists
	}
	return &team.Team{Name: t.Name}, nil
}

func (repo *TeamRepoImpl) GetTeamByName(ctx context.Context, name string) (*team.Team, error) {
	teamExistsQuery := `SELECT EXISTS(SELECT 1 FROM teams WHERE name=$1);`
	var exists bool
	err := repo.db.QueryRowContext(ctx, teamExistsQuery, name).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, models.ErrNotFound
	}

	users, err := repo.userRepo.GetUsersByTeamName(ctx, name)
	if err != nil {
		return nil, err
	}
	var members []team.TeamMember
	for _, v := range users {
		member := team.TeamMember{Id: v.Id, Name: v.Name, IsActive: v.IsActive}
		members = append(members, member)
	}
	return &team.Team{Name: name, Members: members}, nil
}
