package usecase

import (
	"context"
	"errors"
	"github.com/apriliantocecep/ayo-football/services/team/internal/entity"
	"github.com/apriliantocecep/ayo-football/services/team/internal/model"
	"github.com/apriliantocecep/ayo-football/services/team/internal/repository"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type ITeamUseCase interface {
	CreateTeam(ctx context.Context, input model.CreateTeamInput) (*entity.Team, error)
	GetTeamByID(ctx context.Context, id string) (*entity.Team, error)
	UpdateTeam(ctx context.Context, id string, input model.UpdateTeamInput) (*entity.Team, error)
	DeleteTeam(ctx context.Context, id string) error
	ListTeams(ctx context.Context, page, pageSize int) ([]entity.Team, error)
}

type TeamUseCase struct {
	DB             *gorm.DB
	TeamRepository *repository.TeamRepository
}

func (uc *TeamUseCase) CreateTeam(ctx context.Context, input model.CreateTeamInput) (*entity.Team, error) {
	team := &entity.Team{
		Name:      input.Name,
		Logo:      input.Logo,
		FoundedAt: input.FoundedAt,
		Address:   input.Address,
		City:      input.City,
	}

	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := uc.TeamRepository.Create(ctx, tx, team); err != nil {
		return nil, status.Errorf(codes.Aborted, "can not create team")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Aborted, "can not commit create team")
	}

	return team, nil
}

func (uc *TeamUseCase) GetTeamByID(ctx context.Context, id string) (*entity.Team, error) {
	teamID, err := uuid.Parse(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid team ID")
	}

	team, err := uc.TeamRepository.GetByID(ctx, uc.DB, teamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "team not found")
		}
		return nil, err
	}
	return team, nil
}

func (uc *TeamUseCase) UpdateTeam(ctx context.Context, id string, input model.UpdateTeamInput) (*entity.Team, error) {
	teamID, err := uuid.Parse(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid team ID")
	}

	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	existing, err := uc.TeamRepository.GetByID(ctx, uc.DB, teamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "team not found")
		}
		return nil, err
	}

	existing.Name = input.Name
	existing.Logo = input.Logo
	existing.FoundedAt = input.FoundedAt
	existing.Address = input.Address
	existing.City = input.City

	if err = uc.TeamRepository.Update(ctx, tx, existing); err != nil {
		return nil, status.Errorf(codes.Aborted, "can not update team")
	}

	if err = tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Aborted, "can not commit update team")
	}

	return existing, nil
}

func (uc *TeamUseCase) DeleteTeam(ctx context.Context, id string) error {
	teamID, err := uuid.Parse(id)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid team ID")
	}

	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	_, err = uc.TeamRepository.GetByID(ctx, uc.DB, teamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return status.Errorf(codes.NotFound, "team not found")
		}
		return err
	}

	err = uc.TeamRepository.Delete(ctx, tx, teamID)
	if err != nil {
		return status.Errorf(codes.Aborted, "can not delete team: %v", err)
	}

	if err = tx.Commit().Error; err != nil {
		return status.Errorf(codes.Aborted, "can not commit delete team")
	}

	return nil
}

func (uc *TeamUseCase) ListTeams(ctx context.Context, page, pageSize int) ([]entity.Team, error) {
	offset := (page - 1) * pageSize
	return uc.TeamRepository.ListAll(ctx, uc.DB, offset, pageSize)
}

var _ ITeamUseCase = &TeamUseCase{}

func NewTeamUseCase(DB *gorm.DB, teamRepository *repository.TeamRepository) *TeamUseCase {
	return &TeamUseCase{DB: DB, TeamRepository: teamRepository}
}
