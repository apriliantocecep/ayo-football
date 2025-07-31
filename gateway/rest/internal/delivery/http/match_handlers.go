package http

import (
	"context"
	"github.com/apriliantocecep/ayo-football/gateway/rest/internal/delivery/grpc_client"
	"github.com/apriliantocecep/ayo-football/gateway/rest/internal/model"
	"github.com/apriliantocecep/ayo-football/services/match/pkg/pb"
	playerpb "github.com/apriliantocecep/ayo-football/services/player/pkg/pb"
	teampb "github.com/apriliantocecep/ayo-football/services/team/pkg/pb"
	"github.com/apriliantocecep/ayo-football/shared/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

type IMatchHandler interface {
	Create(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	List(c *fiber.Ctx) error
	CreateGoal(c *fiber.Ctx) error
	GetGoal(c *fiber.Ctx) error
	UpdateGoal(c *fiber.Ctx) error
	DeleteGoal(c *fiber.Ctx) error
}

type MatchHandler struct {
	Validate            *validator.Validate
	MatchServiceClient  *grpc_client.MatchServiceClient
	PlayerServiceClient *grpc_client.PlayerServiceClient
	TeamServiceClient   *grpc_client.TeamServiceClient
}

func (h *MatchHandler) CreateGoal(c *fiber.Ctx) error {
	id := c.Params("id") // match id
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing match id")
	}

	match, err := h.MatchServiceClient.Client.GetMatch(c.UserContext(), &pb.GetMatchRequest{Id: id})
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	var req model.CreateGoalRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	player, err := h.PlayerServiceClient.Client.GetPlayer(c.UserContext(), &playerpb.GetPlayerRequest{Id: req.PlayerID})
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	if err := h.Validate.Struct(&req); err != nil {
		//return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	grpcReq := &pb.CreateGoalRequest{
		MatchId:  match.Id,
		PlayerId: player.Id,
		ScoredAt: req.ScoredAt,
	}
	goal, err := h.MatchServiceClient.Client.CreateGoal(ctx, grpcReq)
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	return c.Status(fiber.StatusCreated).JSON(goal)
}

func (h *MatchHandler) GetGoal(c *fiber.Ctx) error {
	id := c.Params("goalId") // goal id
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing id")
	}

	matchId := c.Params("id") // match id
	if matchId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing match id")
	}

	_, err := h.MatchServiceClient.Client.GetMatch(c.UserContext(), &pb.GetMatchRequest{Id: matchId})
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	goal, err := h.MatchServiceClient.Client.GetGoal(ctx, &pb.GetGoalRequest{Id: id})
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	return c.JSON(goal)
}

func (h *MatchHandler) UpdateGoal(c *fiber.Ctx) error {
	id := c.Params("goalId") // goal id
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing id")
	}

	var req model.UpdateGoalRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	matchId := c.Params("id") // match id
	if matchId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing match id")
	}

	match, err := h.MatchServiceClient.Client.GetMatch(c.UserContext(), &pb.GetMatchRequest{Id: matchId})
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	if err := h.Validate.Struct(&req); err != nil {
		return err
	}

	player, err := h.PlayerServiceClient.Client.GetPlayer(c.UserContext(), &playerpb.GetPlayerRequest{Id: req.PlayerID})
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	updated, err := h.MatchServiceClient.Client.UpdateGoal(ctx, &pb.UpdateGoalRequest{
		Id:       id,
		MatchId:  match.Id,
		PlayerId: player.Id,
		ScoredAt: req.ScoredAt,
	})
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	return c.JSON(updated)
}

func (h *MatchHandler) DeleteGoal(c *fiber.Ctx) error {
	id := c.Params("goalId") // goal id
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing id")
	}

	matchId := c.Params("id") // match id
	if matchId == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing match id")
	}

	_, err := h.MatchServiceClient.Client.GetMatch(c.UserContext(), &pb.GetMatchRequest{Id: matchId})
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err = h.MatchServiceClient.Client.DeleteGoal(ctx, &pb.DeleteGoalRequest{Id: id})
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	return c.JSON(fiber.Map{"status": "deleted"})
}

func (h *MatchHandler) Create(c *fiber.Ctx) error {
	var req model.CreateMatchRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := h.Validate.Struct(&req); err != nil {
		//return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		return err
	}

	// validate a home team
	homeTeam, err := h.TeamServiceClient.Client.GetTeam(c.UserContext(), &teampb.GetTeamRequest{Id: req.HomeTeamID})
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	// validate an away team
	awayTeam, err := h.TeamServiceClient.Client.GetTeam(c.UserContext(), &teampb.GetTeamRequest{Id: req.AwayTeamID})
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	if homeTeam.Id == awayTeam.Id {
		return fiber.NewError(fiber.StatusInternalServerError, "home team and away team are the same")
	}

	grpcReq := &pb.CreateMatchRequest{
		Date:       req.Date,
		Venue:      req.Venue,
		HomeTeamId: homeTeam.Id,
		AwayTeamId: awayTeam.Id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	match, err := h.MatchServiceClient.Client.CreateMatch(ctx, grpcReq)
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	return c.Status(fiber.StatusCreated).JSON(match)
}

func (h *MatchHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing id")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	match, err := h.MatchServiceClient.Client.GetMatch(ctx, &pb.GetMatchRequest{Id: id})
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	return c.JSON(match)
}

func (h *MatchHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.UpdateMatchRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if err := h.Validate.Struct(&req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	updated, err := h.MatchServiceClient.Client.UpdateMatch(ctx, &pb.UpdateMatchRequest{
		Id:         id,
		Date:       req.Date,
		Venue:      req.Venue,
		HomeTeamId: req.HomeTeamID,
		AwayTeamId: req.AwayTeamID,
	})
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	return c.JSON(updated)
}

func (h *MatchHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := h.MatchServiceClient.Client.DeleteMatch(ctx, &pb.DeleteMatchRequest{Id: id})
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	return c.JSON(fiber.Map{"status": "deleted"})
}

func (h *MatchHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, err := h.MatchServiceClient.Client.ListMatch(ctx, &pb.ListMatchRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
	})
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	return c.JSON(resp)
}

var _ IMatchHandler = &MatchHandler{}

func NewMatchHandler(validate *validator.Validate, matchServiceClient *grpc_client.MatchServiceClient, playerServiceClient *grpc_client.PlayerServiceClient, teamServiceClient *grpc_client.TeamServiceClient) *MatchHandler {
	return &MatchHandler{
		Validate:            validate,
		MatchServiceClient:  matchServiceClient,
		PlayerServiceClient: playerServiceClient,
		TeamServiceClient:   teamServiceClient,
	}
}
