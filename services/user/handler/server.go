package handler

import (
	"warehouseai/user/adapter/grpc/client"
	"warehouseai/user/dataservice"
	"warehouseai/user/handler/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	db           *dataservice.Database
	logger       *logrus.Logger
	authReceiver *client.GrpcClient
}

func NewHttpHandler(db *dataservice.Database, logger *logrus.Logger) *Handler {
	authReceiver := client.NewGrpcClient("auth-service:8041")

	return &Handler{
		db:           db,
		logger:       logger,
		authReceiver: authReceiver,
	}
}

func (h *Handler) Start(port string) error {
	app := fiber.New()
	app.Use(setupCORS())

	sessionMw := middleware.Session(h.logger, h.authReceiver)
	userMw := middleware.User(h.logger, h.db)

	route := app.Group("/user")
	route.Patch("/update", sessionMw, h.UpdatePersonalDataHandler)
	route.Patch("/update/email", sessionMw, h.UpdateEmailHandler)
	route.Patch("/update/password", sessionMw, userMw, h.UpdatePasswordHandler)
	// route.Patch("/favorites/add", h.sessionMiddleware, h.userMiddleware, h.AddFavoriteAIHandler)
	// route.Patch("/favorites/remove", h.sessionMiddleware, h.userMiddleware, h.RemoveFavoriteAIHandler)
	// route.Get("/favorites", h.sessionMiddleware, h.GetFavoriteAIHandler)
	route.Get("/verify/:code", sessionMw, userMw, h.UpdateVerificationHandler)

	return app.Listen(port)
}

func setupCORS() func(*fiber.Ctx) error {
	return cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
		AllowOrigins:     "http://localhost:3000, https://warehouse-ai-frontend.vercel.app",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	})
}
