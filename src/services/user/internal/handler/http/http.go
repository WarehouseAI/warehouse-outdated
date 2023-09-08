package http

import (
	m "warehouse/src/services/user/pkg/model"

	"github.com/gofiber/fiber/v2"
)

type UserPublicAPI struct {
	svc m.UserService
}

func NewUserPublicAPI(svc m.UserService) *UserPublicAPI {
	return &UserPublicAPI{
		svc: svc,
	}
}

func (api *UserPublicAPI) CreateHandler(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusCreated)
}

// INIT
func (api *UserPublicAPI) Init() *fiber.App {
	app := fiber.New()
	route := app.Group("/user")

	route.Post("/create", api.CreateHandler)

	return app
}
