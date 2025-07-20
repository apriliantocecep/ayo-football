package grpc_server

import (
	"context"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/model"
	"github.com/apriliantocecep/posfin-blog/services/auth/internal/usecase"
	"github.com/apriliantocecep/posfin-blog/services/auth/pkg/pb"
)

type AuthServer struct {
	UserUseCase *usecase.UserUseCase
	pb.UnimplementedAuthServiceServer
}

func (a *AuthServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	req := model.LoginRequest{
		Identity: in.Identifier,
		Password: in.Password,
	}

	res, err := a.UserUseCase.Login(ctx, &req)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		AccessToken: res.AccessToken,
		ExpiresAt:   res.AccessTokenExpiresAt.String(),
	}, nil
}

func (a *AuthServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	req := model.RegisterRequest{
		Name:     in.Name,
		Email:    in.Email,
		Password: in.Password,
	}

	res, err := a.UserUseCase.Register(ctx, &req)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		UserId:   res.UserId,
		Username: res.Username,
	}, nil
}

func (a *AuthServer) RegisterWithQueue(ctx context.Context, in *pb.RegisterWithQueueRequest) (*pb.RegisterWithQueueResponse, error) {
	req := model.RegisterRequest{
		Name:     in.Name,
		Email:    in.Email,
		Password: in.Password,
	}

	res, err := a.UserUseCase.RegisterWithQueue(ctx, &req)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterWithQueueResponse{
		UserId:   res.UserId,
		Username: res.Username,
	}, nil
}

func (a *AuthServer) ValidateToken(ctx context.Context, in *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	userId, err := a.UserUseCase.ValidateToken(ctx, in.GetToken())
	if err != nil {
		return nil, err
	}
	return &pb.ValidateTokenResponse{UserId: userId}, nil
}

func (a *AuthServer) GetUserById(ctx context.Context, in *pb.GetUserByIdRequest) (*pb.User, error) {
	res, err := a.UserUseCase.GetUserById(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.User{
		Id:        res.ID,
		Name:      res.Name,
		Email:     res.Email,
		Username:  res.Username,
		CreatedAt: res.CreatedAt.String(),
	}, nil
}

func NewAuthServer(userUseCase *usecase.UserUseCase) *AuthServer {
	return &AuthServer{UserUseCase: userUseCase}
}
