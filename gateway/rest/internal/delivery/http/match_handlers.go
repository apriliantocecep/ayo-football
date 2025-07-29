package http

import (
	"context"
	"github.com/apriliantocecep/ayo-football/gateway/rest/internal/delivery/grpc_client"
	"github.com/apriliantocecep/ayo-football/gateway/rest/internal/model"
	"github.com/apriliantocecep/ayo-football/services/match/pkg/pb"
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
}

type MatchHandler struct {
	Validate           *validator.Validate
	MatchServiceClient *grpc_client.MatchServiceClient
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

	// TODO validate team id home and away

	grpcReq := &pb.CreateMatchRequest{
		Date:       req.Date,
		Venue:      req.Venue,
		HomeTeamId: req.HomeTeamID,
		AwayTeamId: req.AwayTeamID,
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

func NewMatchHandler(validate *validator.Validate, matchServiceClient *grpc_client.MatchServiceClient) *MatchHandler {
	return &MatchHandler{Validate: validate, MatchServiceClient: matchServiceClient}
}
