package grpc_server

import (
	"context"
	"errors"
	"github.com/apriliantocecep/ayo-football/services/match/internal/entity"
	"github.com/apriliantocecep/ayo-football/services/match/internal/model"
	"github.com/apriliantocecep/ayo-football/services/match/internal/usecase"
	"github.com/apriliantocecep/ayo-football/services/match/pkg/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"time"
)

type MatchServer struct {
	MatchUseCase *usecase.MatchUseCase
	pb.UnimplementedMatchServiceServer
}

func (s *MatchServer) CreateMatch(ctx context.Context, in *pb.CreateMatchRequest) (*pb.Match, error) {
	date, err := time.Parse("2006-01-02 15:04:05", in.GetDate())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid date")
	}
	input := model.CreateMatchInput{
		Date:       date,
		Venue:      in.GetVenue(),
		HomeTeamID: in.GetHomeTeamId(),
		AwayTeamID: in.GetAwayTeamId(),
	}
	match, err := s.MatchUseCase.Create(ctx, input)
	if err != nil {
		return nil, err
	}
	return toProto(match), nil
}

func (s *MatchServer) GetMatch(ctx context.Context, in *pb.GetMatchRequest) (*pb.Match, error) {
	match, err := s.MatchUseCase.GetByID(ctx, in.GetId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "match not found")
		}
		return nil, err
	}

	return toProto(match), nil
}

func (s *MatchServer) UpdateMatch(ctx context.Context, in *pb.UpdateMatchRequest) (*pb.Match, error) {
	date, err := time.Parse("2006-01-02 15:04:05", in.GetDate())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid date")
	}
	input := model.UpdateMatchInput{
		Date:       date,
		Venue:      in.GetVenue(),
		HomeTeamID: in.GetHomeTeamId(),
		AwayTeamID: in.GetAwayTeamId(),
	}
	match, err := s.MatchUseCase.Update(ctx, in.GetId(), input)
	if err != nil {
		return nil, err
	}
	return toProto(match), nil
}

func (s *MatchServer) DeleteMatch(ctx context.Context, in *pb.DeleteMatchRequest) (*emptypb.Empty, error) {
	err := s.MatchUseCase.Delete(ctx, in.GetId())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *MatchServer) ListMatch(ctx context.Context, in *pb.ListMatchRequest) (*pb.ListMatchResponse, error) {
	page := int(in.GetPage())
	if page <= 0 {
		page = 1
	}
	pageSize := int(in.GetPageSize())
	if pageSize <= 0 {
		pageSize = 10
	}

	players, err := s.MatchUseCase.List(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	var protoMatches []*pb.Match
	for _, t := range players {
		protoMatches = append(protoMatches, toProto(t))
	}

	return &pb.ListMatchResponse{Matches: protoMatches}, nil
}

func (s *MatchServer) CreateGoal(ctx context.Context, in *pb.CreateGoalRequest) (*pb.Goal, error) {
	//TODO implement me
	panic("implement me")
}

func (s *MatchServer) DeleteGoal(ctx context.Context, in *pb.DeleteGoalRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func NewMatchServer(matchUseCase *usecase.MatchUseCase) *MatchServer {
	return &MatchServer{MatchUseCase: matchUseCase}
}

func toProto(match *entity.Match) *pb.Match {
	return &pb.Match{
		Id:         match.ID.String(),
		Date:       match.Date.String(),
		Venue:      match.Venue,
		HomeTeamId: match.HomeTeamID.String(),
		AwayTeamId: match.AwayTeamID.String(),
		CreatedAt:  match.CreatedAt.String(),
		UpdatedAt:  match.UpdatedAt.String(),
	}
}
