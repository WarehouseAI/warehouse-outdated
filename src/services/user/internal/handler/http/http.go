package http

import (
	"context"
	pg "warehouse/src/internal/database/postgresdb"
	u "warehouse/src/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserServiceProvider struct {
	userDatabase *pg.PostgresDatabase[pg.User]
	logger       *logrus.Logger
	ctx          context.Context
}

func NewUserPublicAPI(userDatabase *pg.PostgresDatabase[pg.User], logger *logrus.Logger) *UserServiceProvider {
	ctx := context.Background()

	return &UserServiceProvider{
		userDatabase: userDatabase,
		logger:       logger,
		ctx:          ctx,
	}
}

func (pvd *UserServiceProvider) CreateHandler(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusCreated)
}

// INIT
func (pvd *UserServiceProvider) Init() *fiber.App {
	app := fiber.New()
	app.Use(u.SetupCORS())
	route := app.Group("/user")

	route.Post("/create", pvd.CreateHandler)

	return app
}
