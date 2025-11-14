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

func (repo *TeamRepoImpl) AddTeam(
	ctx context.Context,
	t team.AddTeamRequest,
) (team.AddTeamResponse, error) {
	var resp team.AddTeamResponse
	if err := t.Team.Validate(); err != nil {
		return resp, err
	}

	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return resp, err
	}

	insertTeamQuery := `INSERT INTO teams (name) VALUES ($1) RETURNING name;`
	_, err = tx.ExecContext(ctx, insertTeamQuery, t.Team.Name)
	if err != nil {
		tx.Rollback()
		slog.Error(err.Error())
		return resp, models.ErrTeamExists
	}

	addMemberQuery := `UPDATE users SET team_name=$2 WHERE id=$1;`
	for _, v := range t.Team.Members {
		_, err := tx.ExecContext(ctx, addMemberQuery, v.Id, t.Team.Name)
		if err != nil {
			tx.Rollback()
			slog.Error(err.Error())
			return resp, models.ErrNotAssigned
		}
	}

	if err := tx.Commit(); err != nil {
		slog.Error(err.Error())
		return resp, models.ErrTeamExists
	}
	getTeamResp, err := repo.GetTeamByName(ctx, t.Team.Name)
	if err != nil {
		return resp, err
	}
	return team.AddTeamResponse{Team: team.Team(getTeamResp)}, nil
}

func (repo *TeamRepoImpl) GetTeamByName(
	ctx context.Context,
	name string,
) (team.GetTeamResponse, error) {
	resp := team.GetTeamResponse{Name: name}
	teamExistsQuery := `SELECT EXISTS(SELECT 1 FROM teams WHERE name=$1);`
	var exists bool
	err := repo.db.QueryRowContext(ctx, teamExistsQuery, name).Scan(&exists)
	if err != nil {
		return resp, err
	}
	if !exists {
		return resp, models.ErrNotFound
	}

	users, err := repo.userRepo.GetUsersByTeamName(ctx, name)
	if err != nil {
		return resp, err
	}
	members := []team.TeamMember{}
	for _, v := range users {
		member := team.TeamMember{Id: v.Id, Name: v.Name, IsActive: v.IsActive}
		members = append(members, member)
	}
	resp.Members = members
	return resp, nil
}
