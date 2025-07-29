package usecase

import (
	"context"
	"errors"
	"github.com/apriliantocecep/ayo-football/services/player/internal/entity"
	"github.com/apriliantocecep/ayo-football/services/player/internal/model"
	"github.com/apriliantocecep/ayo-football/services/player/internal/repository"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type IPlayerUseCase interface {
	CreatePlayer(ctx context.Context, input model.CreatePlayerInput) (*entity.Player, error)
	GetByID(ctx context.Context, id string) (*entity.Player, error)
	UpdatePlayer(ctx context.Context, id string, input model.UpdatePlayerInput) (*entity.Player, error)
	DeletePlayer(ctx context.Context, id string) error
	ListByTeamID(ctx context.Context, teamID string, page, pageSize int) ([]*entity.Player, error)
}

type PlayerUseCase struct {
	DB               *gorm.DB
	PlayerRepository *repository.PlayerRepository
}

func (uc *PlayerUseCase) CreatePlayer(ctx context.Context, input model.CreatePlayerInput) (*entity.Player, error) {
	teamID, err := uuid.Parse(input.TeamID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid team ID")
	}

	used, err := uc.PlayerRepository.IsBackNumberUsed(ctx, uc.DB, teamID, input.BackNumber, "")
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "can't check back number: %v", err)
	}
	if used {
		return nil, status.Errorf(codes.AlreadyExists, "back number already used in this team")
	}

	player := &entity.Player{
		TeamID:     teamID,
		Name:       input.Name,
		Height:     input.Height,
		Weight:     input.Weight,
		Position:   input.Position,
		BackNumber: input.BackNumber,
	}

	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := uc.PlayerRepository.Create(ctx, tx, player); err != nil {
		return nil, status.Errorf(codes.Aborted, "can not create player")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Aborted, "can not commit create player")
	}

	return player, nil
}

func (uc *PlayerUseCase) GetByID(ctx context.Context, id string) (*entity.Player, error) {
	playerID, err := uuid.Parse(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid player ID")
	}

	player, err := uc.PlayerRepository.GetByID(ctx, uc.DB, playerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "player not found")
		}
		return nil, err
	}
	return player, nil
}

func (uc *PlayerUseCase) UpdatePlayer(ctx context.Context, id string, input model.UpdatePlayerInput) (*entity.Player, error) {
	playerID, err := uuid.Parse(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid player ID")
	}

	teamID, err := uuid.Parse(input.TeamID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid team ID")
	}

	used, err := uc.PlayerRepository.IsBackNumberUsed(ctx, uc.DB, teamID, input.BackNumber, id)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "can't check back number: %v", err)
	}
	if used {
		return nil, status.Errorf(codes.AlreadyExists, "back number already used in this team")
	}

	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	existing, err := uc.PlayerRepository.GetByID(ctx, tx, playerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "player not found")
		}
		return nil, err
	}

	existing.Name = input.Name
	existing.Height = input.Height
	existing.Weight = input.Weight
	existing.Position = input.Position
	existing.BackNumber = input.BackNumber

	if err = uc.PlayerRepository.Update(ctx, tx, existing); err != nil {
		return nil, status.Errorf(codes.Aborted, "can not update player")
	}

	if err = tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Aborted, "can not commit update player")
	}

	return existing, nil
}

func (uc *PlayerUseCase) DeletePlayer(ctx context.Context, id string) error {
	playerID, err := uuid.Parse(id)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid player ID")
	}

	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	_, err = uc.PlayerRepository.GetByID(ctx, uc.DB, playerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return status.Errorf(codes.NotFound, "player not found")
		}
		return err
	}

	err = uc.PlayerRepository.Delete(ctx, tx, playerID)
	if err != nil {
		return status.Errorf(codes.Aborted, "can not delete player: %v", err)
	}

	if err = tx.Commit().Error; err != nil {
		return status.Errorf(codes.Aborted, "can not commit delete player")
	}

	return nil
}

func (uc *PlayerUseCase) ListByTeamID(ctx context.Context, teamID string, page, pageSize int) ([]*entity.Player, error) {
	teamUUID, err := uuid.Parse(teamID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid team ID")
	}
	offset := (page - 1) * pageSize
	return uc.PlayerRepository.ListByTeamID(ctx, uc.DB, teamUUID, offset, pageSize)
}

var _ IPlayerUseCase = &PlayerUseCase{}

func NewPlayerUseCase(DB *gorm.DB, playerRepository *repository.PlayerRepository) *PlayerUseCase {
	return &PlayerUseCase{DB: DB, PlayerRepository: playerRepository}
}
