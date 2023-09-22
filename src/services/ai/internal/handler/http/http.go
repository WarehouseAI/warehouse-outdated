package http

import (
	"context"
	"errors"
	dbm "warehouse/src/internal/db/models"
	"warehouse/src/internal/dto"
	mv "warehouse/src/internal/middleware"
	svc "warehouse/src/services/ai/internal/service"
	m "warehouse/src/services/ai/pkg/models"

	"github.com/gofiber/fiber/v2"
)

type APIInstance struct {
	svc svc.AIService
	sMw *mv.SessionMiddleware
	uMw *mv.UserMiddleware
}

func NewAiAPI(service svc.AIService, sessionMiddleware *mv.SessionMiddleware, userMiddleware *mv.UserMiddleware) *APIInstance {
	return &APIInstance{
		svc: service,
		sMw: sessionMiddleware,
		uMw: userMiddleware,
	}
}

func (api *APIInstance) CreateHandler(c *fiber.Ctx) error {
	user := c.Locals("user").(*dbm.User)
	var aiInfo m.CreateAIRequest

	if err := c.BodyParser(&aiInfo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: dto.BadRequestError.Error()})
	}

	apiInfo, err := api.svc.Create(context.Background(), &aiInfo, user)

	if err != nil && errors.Is(err, dto.InternalError) {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Message: err.Error()})
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(apiInfo)
}

func (api *APIInstance) AddCommandHandler(c *fiber.Ctx) error {
	var commandInfo m.AddCommandRequest

	if err := c.BodyParser(&commandInfo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Message: dto.BadRequestError.Error()})
	}

	if err := api.svc.AddCommand(context.Background(), &commandInfo); err != nil {
		statusCode := fiber.StatusInternalServerError
		return c.Status(statusCode).JSON(dto.ErrorResponse{Code: statusCode, Message: err.Error()})
	}

	return c.SendStatus(fiber.StatusCreated)
}

// INIT
func (api *APIInstance) Init() *fiber.App {
	app := fiber.New()
	route := app.Group("/ai")

	route.Post("/create", api.sMw.Session, api.uMw.User, api.CreateHandler)
	route.Post("/command/create", api.sMw.Session, api.AddCommandHandler)

	return app
}
