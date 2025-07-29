package grpc_server

import (
	"context"
	"errors"
	"github.com/apriliantocecep/ayo-football/services/player/internal/entity"
	"github.com/apriliantocecep/ayo-football/services/player/internal/model"
	"github.com/apriliantocecep/ayo-football/services/player/internal/usecase"
	"github.com/apriliantocecep/ayo-football/services/player/pkg/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"strings"
)

type PlayerServer struct {
	PlayerUseCase *usecase.PlayerUseCase
	pb.UnimplementedPlayerServiceServer
}

func (s *PlayerServer) CreatePlayer(ctx context.Context, in *pb.CreatePlayerRequest) (*pb.Player, error) {
	input := model.CreatePlayerInput{
		TeamID:     in.GetTeamId(),
		Name:       in.GetName(),
		Height:     in.GetHeight(),
		Weight:     in.GetWeight(),
		Position:   in.GetPosition().String(),
		BackNumber: in.GetBackNumber(),
	}
	player, err := s.PlayerUseCase.CreatePlayer(ctx, input)
	if err != nil {
		return nil, err
	}
	return toProto(player), nil
}

func (s *PlayerServer) GetPlayer(ctx context.Context, in *pb.GetPlayerRequest) (*pb.Player, error) {
	player, err := s.PlayerUseCase.GetByID(ctx, in.GetId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "player not found")
		}
		return nil, err
	}

	return toProto(player), nil
}

func (s *PlayerServer) UpdatePlayer(ctx context.Context, in *pb.UpdatePlayerRequest) (*pb.Player, error) {
	input := model.UpdatePlayerInput{
		TeamID:     in.GetTeamId(),
		Name:       in.GetName(),
		Height:     in.GetHeight(),
		Weight:     in.GetWeight(),
		Position:   in.GetPosition().String(),
		BackNumber: in.GetBackNumber(),
	}
	player, err := s.PlayerUseCase.UpdatePlayer(ctx, in.GetId(), input)
	if err != nil {
		return nil, err
	}
	return toProto(player), nil
}

func (s *PlayerServer) DeletePlayer(ctx context.Context, in *pb.DeletePlayerRequest) (*emptypb.Empty, error) {
	err := s.PlayerUseCase.DeletePlayer(ctx, in.GetId())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *PlayerServer) ListPlayersByTeam(ctx context.Context, in *pb.ListPlayersByTeamRequest) (*pb.ListPlayersResponse, error) {
	page := int(in.GetPage())
	if page <= 0 {
		page = 1
	}
	pageSize := int(in.GetPageSize())
	if pageSize <= 0 {
		pageSize = 10
	}

	players, err := s.PlayerUseCase.ListByTeamID(ctx, in.GetTeamId(), page, pageSize)
	if err != nil {
		return nil, err
	}

	var protoPlayers []*pb.Player
	for _, t := range players {
		protoPlayers = append(protoPlayers, toProto(t))
	}

	return &pb.ListPlayersResponse{Players: protoPlayers}, nil
}

func NewPlayerServer(playerUseCase *usecase.PlayerUseCase) *PlayerServer {
	return &PlayerServer{PlayerUseCase: playerUseCase}
}

func toProto(player *entity.Player) *pb.Player {
	return &pb.Player{
		Id:         player.ID.String(),
		TeamId:     player.TeamID.String(),
		Name:       player.Name,
		Height:     player.Height,
		Weight:     player.Weight,
		Position:   mapStringToPosition(player.Position),
		BackNumber: player.BackNumber,
		CreatedAt:  player.CreatedAt.String(),
		UpdatedAt:  player.UpdatedAt.String(),
	}
}

func mapStringToPosition(pos string) pb.Position {
	switch strings.ToUpper(pos) {
	case "PENYERANG":
		return pb.Position_PENYERANG
	case "GELANDANG":
		return pb.Position_GELANDANG
	case "BERTAHAN":
		return pb.Position_BERTAHAN
	case "PENJAGA_GAWANG":
		return pb.Position_PENJAGA_GAWANG
	default:
		return pb.Position_UNKNOWN
	}
}
