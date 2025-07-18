package grpc_server

import (
	"context"
	"github.com/apriliantocecep/posfin-blog/services/article/internal/model"
	"github.com/apriliantocecep/posfin-blog/services/article/internal/usecase"
	"github.com/apriliantocecep/posfin-blog/services/article/pkg/pb"
)

type ArticleServer struct {
	ArticleUseCase *usecase.ArticleUseCase
	pb.UnimplementedArticleServiceServer
}

func (a *ArticleServer) SubmitArticle(ctx context.Context, in *pb.SubmitArticleRequest) (*pb.SubmitArticleResponse, error) {
	req := model.ArticleRequest{
		Content: in.HtmlContent,
		UserId:  in.UserId,
		Title:   in.Title,
		Author:  in.Author,
	}
	res, err := a.ArticleUseCase.Insert(ctx, &req)
	if err != nil {
		return nil, err
	}

	return &pb.SubmitArticleResponse{
		ArticleId: res.ArticleId,
		Status:    res.Status,
	}, nil
}

func (a *ArticleServer) PublishArticle(ctx context.Context, in *pb.PublishArticleRequest) (*pb.PublishArticleResponse, error) {
	req := model.ModerationRequest{
		ArticleId: in.GetArticleId(),
		UserId:    in.GetUserId(),
	}
	err := a.ArticleUseCase.SendForModeration(ctx, &req)
	if err != nil {
		return nil, err
	}

	return &pb.PublishArticleResponse{Status: "sent for content moderation"}, nil
}

func NewArticleServer(articleUseCase *usecase.ArticleUseCase) *ArticleServer {
	return &ArticleServer{ArticleUseCase: articleUseCase}
}
