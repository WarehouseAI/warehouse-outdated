package http

import (
	u "warehouse/src/internal/utils"
	svc "warehouse/src/services/user/internal/service/user"

	"github.com/gofiber/fiber/v2"
)

type UserPublicAPI struct {
	svc svc.UserService
}

func NewUserPublicAPI(svc svc.UserService) *UserPublicAPI {
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
	app.Use(u.SetupCORS())
	route := app.Group("/user")

	route.Post("/create", api.CreateHandler)

	return app
}
