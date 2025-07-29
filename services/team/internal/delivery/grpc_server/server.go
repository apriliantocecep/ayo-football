package grpc_server

import (
	"context"
	"errors"
	"github.com/apriliantocecep/ayo-football/services/team/internal/entity"
	"github.com/apriliantocecep/ayo-football/services/team/internal/model"
	"github.com/apriliantocecep/ayo-football/services/team/internal/usecase"
	"github.com/apriliantocecep/ayo-football/services/team/pkg/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type TeamServer struct {
	TeamUseCase *usecase.TeamUseCase
	pb.UnimplementedTeamServiceServer
}

func (s *TeamServer) CreateTeam(ctx context.Context, in *pb.CreateTeamRequest) (*pb.Team, error) {
	input := model.CreateTeamInput{
		Name:      in.GetName(),
		Logo:      in.GetLogo(),
		FoundedAt: int(in.GetFoundedAt()),
		Address:   in.GetAddress(),
		City:      in.GetCity(),
	}

	team, err := s.TeamUseCase.CreateTeam(ctx, input)
	if err != nil {
		return nil, err
	}

	return toProtoTeam(team), nil
}

func (s *TeamServer) GetTeam(ctx context.Context, in *pb.GetTeamRequest) (*pb.Team, error) {
	team, err := s.TeamUseCase.GetTeamByID(ctx, in.GetId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "team not found")
		}
		return nil, err
	}

	return toProtoTeam(team), nil
}

func (s *TeamServer) UpdateTeam(ctx context.Context, in *pb.UpdateTeamRequest) (*pb.Team, error) {
	input := model.UpdateTeamInput{
		Name:      in.GetName(),
		Logo:      in.GetLogo(),
		FoundedAt: int(in.GetFoundedAt()),
		Address:   in.GetAddress(),
		City:      in.GetCity(),
	}

	team, err := s.TeamUseCase.UpdateTeam(ctx, in.GetId(), input)
	if err != nil {
		return nil, err
	}

	return toProtoTeam(team), nil
}

func (s *TeamServer) DeleteTeam(ctx context.Context, in *pb.DeleteTeamRequest) (*pb.DeleteTeamResponse, error) {
	err := s.TeamUseCase.DeleteTeam(ctx, in.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.DeleteTeamResponse{Status: "success"}, nil
}

func (s *TeamServer) ListTeams(ctx context.Context, in *pb.ListTeamsRequest) (*pb.ListTeamsResponse, error) {
	page := int(in.GetPage())
	if page <= 0 {
		page = 1
	}
	pageSize := int(in.GetPageSize())
	if pageSize <= 0 {
		pageSize = 10
	}

	teams, err := s.TeamUseCase.ListTeams(ctx, page, pageSize)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list teams: %v", err)
	}

	var protoTeams []*pb.Team
	for _, t := range teams {
		protoTeams = append(protoTeams, toProtoTeam(&t))
	}

	return &pb.ListTeamsResponse{
		Teams: protoTeams,
		Total: int32(len(protoTeams)),
	}, nil
}

func NewTeamServer(teamUseCase *usecase.TeamUseCase) *TeamServer {
	return &TeamServer{TeamUseCase: teamUseCase}
}

func toProtoTeam(t *entity.Team) *pb.Team {
	return &pb.Team{
		Id:        t.ID.String(),
		Name:      t.Name,
		Logo:      t.Logo,
		FoundedAt: int32(t.FoundedAt),
		Address:   t.Address,
		City:      t.City,
		CreatedAt: t.CreatedAt.String(),
		UpdatedAt: t.UpdatedAt.String(),
	}
}
