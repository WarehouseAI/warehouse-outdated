package server

import (
	"warehouseai/user/adapter/broker/mail"
	"warehouseai/user/adapter/grpc/client/auth"
	"warehouseai/user/dataservice/userdata"
	h "warehouseai/user/server/handlers"
	"warehouseai/user/server/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
)

func StartServer(port string, db *userdata.Database, mailProducer *mail.MailProducer, logger *logrus.Logger) error {
	handler := newHttpHandler(db, mailProducer, logger)
	app := fiber.New()
	app.Use(setupCORS())

	sessionMw := middleware.Session(logger, handler.AuthClient)
	userMw := middleware.User(logger, handler.DB)

	route := app.Group("/user")
	route.Patch("/update", sessionMw, handler.UpdatePersonalDataHandler)
	route.Patch("/update/email", sessionMw, handler.UpdateEmailHandler)
	route.Patch("/update/password", sessionMw, userMw, handler.UpdatePasswordHandler)
	// route.Patch("/favorites/add", h.sessionMiddleware, h.userMiddleware, h.AddFavoriteAIHandler)
	// route.Patch("/favorites/remove", h.sessionMiddleware, h.userMiddleware, h.RemoveFavoriteAIHandler)
	// route.Get("/favorites", h.sessionMiddleware, h.GetFavoriteAIHandler)

	return app.Listen(port)
}

func newHttpHandler(db *userdata.Database, mailProducer *mail.MailProducer, logger *logrus.Logger) *h.Handler {
	authClient := auth.NewAuthGrpcClient("auth:8041")

	return &h.Handler{
		DB:           db,
		Logger:       logger,
		MailProducer: mailProducer,
		AuthClient:   authClient,
	}
}

func setupCORS() func(*fiber.Ctx) error {
	return cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
		AllowOrigins:     "http://localhost:3000, https://warehouse-ai-frontend.vercel.app",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	})
}
