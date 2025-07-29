package grpc_server

import (
	"context"
	"github.com/apriliantocecep/ayo-football/services/moderation/internal/usecase"
	"github.com/apriliantocecep/ayo-football/services/moderation/pkg/pb"
)

type ModerationServer struct {
	MetadataUseCase *usecase.MetadataUseCase
	pb.UnimplementedModerationServiceServer
}

func (m *ModerationServer) PublishArticle(ctx context.Context, in *pb.PublishArticleRequest) (*pb.PublishArticleResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewModerationServer(metadataUseCase *usecase.MetadataUseCase) *ModerationServer {
	return &ModerationServer{MetadataUseCase: metadataUseCase}
}
