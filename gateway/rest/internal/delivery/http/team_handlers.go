package http

import (
	"context"
	"github.com/apriliantocecep/ayo-football/gateway/rest/internal/delivery/grpc_client"
	"github.com/apriliantocecep/ayo-football/gateway/rest/internal/model"
	"github.com/apriliantocecep/ayo-football/services/team/pkg/pb"
	"github.com/apriliantocecep/ayo-football/shared/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

type ITeamHandler interface {
	CreateTeam(c *fiber.Ctx) error
	GetTeam(c *fiber.Ctx) error
	UpdateTeam(c *fiber.Ctx) error
	DeleteTeam(c *fiber.Ctx) error
	ListTeams(c *fiber.Ctx) error
}

type TeamHandler struct {
	Validate          *validator.Validate
	TeamServiceClient *grpc_client.TeamServiceClient
}

func (h *TeamHandler) CreateTeam(c *fiber.Ctx) error {
	var req model.CreateTeamRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if err := h.Validate.Struct(&req); err != nil {
		//return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		return err
	}

	grpcReq := &pb.CreateTeamRequest{
		Name:      req.Name,
		Logo:      req.Logo,
		FoundedAt: req.FoundedAt,
		Address:   req.Address,
		City:      req.City,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	team, err := h.TeamServiceClient.Client.CreateTeam(ctx, grpcReq)
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	return c.Status(fiber.StatusCreated).JSON(team)
}

func (h *TeamHandler) GetTeam(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing id"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	team, err := h.TeamServiceClient.Client.GetTeam(ctx, &pb.GetTeamRequest{Id: id})
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	return c.JSON(team)
}

func (h *TeamHandler) UpdateTeam(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.UpdateTeamRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	if err := h.Validate.Struct(&req); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	updated, err := h.TeamServiceClient.Client.UpdateTeam(ctx, &pb.UpdateTeamRequest{
		Id:        id,
		Name:      req.Name,
		Logo:      req.Logo,
		FoundedAt: req.FoundedAt,
		Address:   req.Address,
		City:      req.City,
	})
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	return c.JSON(updated)
}

func (h *TeamHandler) DeleteTeam(c *fiber.Ctx) error {
	id := c.Params("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := h.TeamServiceClient.Client.DeleteTeam(ctx, &pb.DeleteTeamRequest{Id: id})
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	return c.JSON(fiber.Map{"status": res.GetStatus()})
}

func (h *TeamHandler) ListTeams(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, err := h.TeamServiceClient.Client.ListTeams(ctx, &pb.ListTeamsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
	})
	if err != nil {
		return utils.HandleGrpcError(err)
	}

	return c.JSON(resp)
}

var _ ITeamHandler = &TeamHandler{}

func NewTeamHandler(validate *validator.Validate, teamServiceClient *grpc_client.TeamServiceClient) *TeamHandler {
	return &TeamHandler{Validate: validate, TeamServiceClient: teamServiceClient}
}
