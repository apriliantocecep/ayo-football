package http

import (
	"context"
	"github.com/apriliantocecep/ayo-football/gateway/rest/internal/delivery/grpc_client"
	"github.com/apriliantocecep/ayo-football/gateway/rest/internal/model"
	"github.com/apriliantocecep/ayo-football/services/player/pkg/pb"
	teampb "github.com/apriliantocecep/ayo-football/services/team/pkg/pb"
	"github.com/apriliantocecep/ayo-football/shared/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
	"time"
)

type IPlayerHandler interface {
	CreatePlayer(c *fiber.Ctx) error
	GetPlayer(c *fiber.Ctx) error
	UpdatePlayer(c *fiber.Ctx) error
	DeletePlayer(c *fiber.Ctx) error
	ListPlayerByTeam(c *fiber.Ctx) error
}

type PlayerHandler struct {
	Validate            *validator.Validate
	PlayerServiceClient *grpc_client.PlayerServiceClient
	TeamServiceClient   *grpc_client.TeamServiceClient
}

func (h *PlayerHandler) CreatePlayer(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var req model.CreatePlayerRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := h.Validate.Struct(&req); err != nil {
		//return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		return err
	}

	// validate team
	team, err := h.TeamServiceClient.Client.GetTeam(ctx, &teampb.GetTeamRequest{Id: req.TeamID})
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	positionEnum := model.ValidPositions[strings.ToUpper(req.Position)]
	grpcReq := &pb.CreatePlayerRequest{
		TeamId:     team.Id,
		Name:       req.Name,
		Height:     req.Height,
		Weight:     req.Weight,
		Position:   positionEnum,
		BackNumber: req.BackNumber,
	}

	player, err := h.PlayerServiceClient.Client.CreatePlayer(ctx, grpcReq)
	if err != nil {
		return utils.HandleGrpcError(err)
	}
	return c.Status(fiber.StatusCreated).JSON(model.PlayerToResponse(player))
}

func (h *PlayerHandler) GetPlayer(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing id")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	player, err := h.PlayerServiceClient.Client.GetPlayer(ctx, &pb.GetPlayerRequest{Id: id})
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	return c.JSON(model.PlayerToResponse(player))
}

func (h *PlayerHandler) UpdatePlayer(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.UpdatePlayerRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if err := h.Validate.Struct(&req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// validate team
	team, err := h.TeamServiceClient.Client.GetTeam(ctx, &teampb.GetTeamRequest{Id: req.TeamID})
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	positionEnum := model.ValidPositions[strings.ToUpper(req.Position)]

	updated, err := h.PlayerServiceClient.Client.UpdatePlayer(ctx, &pb.UpdatePlayerRequest{
		Id:         id,
		Name:       req.Name,
		Height:     req.Height,
		Weight:     req.Weight,
		Position:   positionEnum,
		BackNumber: req.BackNumber,
		TeamId:     team.Id,
	})
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	return c.JSON(model.PlayerToResponse(updated))
}

func (h *PlayerHandler) DeletePlayer(c *fiber.Ctx) error {
	id := c.Params("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := h.PlayerServiceClient.Client.DeletePlayer(ctx, &pb.DeletePlayerRequest{Id: id})
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	return c.JSON(fiber.Map{"status": "deleted"})
}

func (h *PlayerHandler) ListPlayerByTeam(c *fiber.Ctx) error {
	teamID := c.Query("team_id")
	if teamID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing team_id")
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// validate team
	team, err := h.TeamServiceClient.Client.GetTeam(ctx, &teampb.GetTeamRequest{Id: teamID})
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	resp, err := h.PlayerServiceClient.Client.ListPlayersByTeam(ctx, &pb.ListPlayersByTeamRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		TeamId:   team.Id,
	})
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	var players []*model.PlayerResource
	for _, player := range resp.Players {
		players = append(players, model.PlayerToResponse(player))
	}

	return c.JSON(players)
}

var _ IPlayerHandler = &PlayerHandler{}

func NewPlayerHandler(validate *validator.Validate, playerServiceClient *grpc_client.PlayerServiceClient, teamServiceClient *grpc_client.TeamServiceClient) *PlayerHandler {
	return &PlayerHandler{Validate: validate, PlayerServiceClient: playerServiceClient, TeamServiceClient: teamServiceClient}
}
