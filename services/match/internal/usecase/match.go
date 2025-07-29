package usecase

import (
	"context"
	"errors"
	"github.com/apriliantocecep/ayo-football/services/match/internal/entity"
	"github.com/apriliantocecep/ayo-football/services/match/internal/model"
	"github.com/apriliantocecep/ayo-football/services/match/internal/repository"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type IMatchUseCase interface {
	Create(ctx context.Context, input model.CreateMatchInput) (*entity.Match, error)
	GetByID(ctx context.Context, id string) (*entity.Match, error)
	Update(ctx context.Context, id string, input model.UpdateMatchInput) (*entity.Match, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, pageSize int) ([]*entity.Match, error)
}

type MatchUseCase struct {
	DB              *gorm.DB
	MatchRepository *repository.MatchRepository
}

func (uc *MatchUseCase) Create(ctx context.Context, input model.CreateMatchInput) (*entity.Match, error) {
	homeTeamID, err := uuid.Parse(input.HomeTeamID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid home team ID")
	}

	awayTeamID, err := uuid.Parse(input.AwayTeamID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid away team ID")
	}

	match := &entity.Match{
		Date:       input.Date,
		Venue:      input.Venue,
		HomeTeamID: homeTeamID,
		AwayTeamID: awayTeamID,
	}

	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := uc.MatchRepository.Create(ctx, tx, match); err != nil {
		return nil, status.Errorf(codes.Aborted, "can not create match")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Aborted, "can not commit create match")
	}

	return match, nil
}

func (uc *MatchUseCase) GetByID(ctx context.Context, id string) (*entity.Match, error) {
	matchID, err := uuid.Parse(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid match ID")
	}

	match, err := uc.MatchRepository.GetByID(ctx, uc.DB, matchID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "match not found")
		}
		return nil, err
	}
	return match, nil
}

func (uc *MatchUseCase) Update(ctx context.Context, id string, input model.UpdateMatchInput) (*entity.Match, error) {
	matchID, err := uuid.Parse(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid match ID")
	}

	homeTeamID, err := uuid.Parse(input.HomeTeamID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid home team ID")
	}

	awayTeamID, err := uuid.Parse(input.AwayTeamID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid away team ID")
	}

	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	existing, err := uc.MatchRepository.GetByID(ctx, tx, matchID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "match not found")
		}
		return nil, err
	}

	existing.Date = input.Date
	existing.Venue = input.Venue
	existing.HomeTeamID = homeTeamID
	existing.AwayTeamID = awayTeamID

	if err := uc.MatchRepository.Update(ctx, tx, existing); err != nil {
		return nil, status.Errorf(codes.Aborted, "can not update match")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Aborted, "can not commit update match")
	}

	return existing, nil
}

func (uc *MatchUseCase) Delete(ctx context.Context, id string) error {
	matchID, err := uuid.Parse(id)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid match ID")
	}

	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	_, err = uc.MatchRepository.GetByID(ctx, uc.DB, matchID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return status.Errorf(codes.NotFound, "match not found")
		}
		return err
	}

	err = uc.MatchRepository.Delete(ctx, tx, matchID)
	if err != nil {
		return status.Errorf(codes.Aborted, "can not delete match: %v", err)
	}

	if err = tx.Commit().Error; err != nil {
		return status.Errorf(codes.Aborted, "can not commit delete match")
	}

	return nil
}

func (uc *MatchUseCase) List(ctx context.Context, page, pageSize int) ([]*entity.Match, error) {
	offset := (page - 1) * pageSize
	return uc.MatchRepository.List(ctx, uc.DB, offset, pageSize)
}

var _ IMatchUseCase = &MatchUseCase{}

func NewMatchUseCase(DB *gorm.DB, matchRepository *repository.MatchRepository) *MatchUseCase {
	return &MatchUseCase{DB: DB, MatchRepository: matchRepository}
}
