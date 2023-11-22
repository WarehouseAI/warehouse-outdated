package server

import (
	"warehouseai/user/adapter/broker/mail"
	"warehouseai/user/adapter/grpc/client/ai"
	"warehouseai/user/adapter/grpc/client/auth"
	"warehouseai/user/dataservice/favoritesdata"
	"warehouseai/user/dataservice/userdata"
	h "warehouseai/user/server/handlers"
	"warehouseai/user/server/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
)

func StartServer(port string, userDb *userdata.Database, favoritesDb *favoritesdata.Database, mailProducer *mail.MailProducer, logger *logrus.Logger) error {
	handler := newHttpHandler(userDb, favoritesDb, mailProducer, logger)
	app := fiber.New()
	app.Use(setupCORS())

	sessionMw := middleware.Session(logger, handler.AuthClient)
	userMw := middleware.User(logger, handler.UserDB)

	route := app.Group("/user")
	route.Patch("/update", sessionMw, handler.UpdatePersonalDataHandler)
	route.Patch("/update/email", sessionMw, handler.UpdateEmailHandler)
	route.Patch("/update/password", sessionMw, userMw, handler.UpdatePasswordHandler)
	route.Patch("/favorites/add", sessionMw, handler.AddFavoriteHandler)
	route.Delete("/favorites/delete", sessionMw, handler.RemoveFavoriteHandler)
	route.Get("/favorites", sessionMw, handler.GetFavoritesHandler)

	return app.Listen(port)
}

func newHttpHandler(userDb *userdata.Database, favoritesDB *favoritesdata.Database, mailProducer *mail.MailProducer, logger *logrus.Logger) *h.Handler {
	authClient := auth.NewAuthGrpcClient("auth:8041")
	aiClient := ai.NewAiGrpcClient("ai:8021")

	return &h.Handler{
		UserDB:       userDb,
		FavoritesDB:  favoritesDB,
		Logger:       logger,
		MailProducer: mailProducer,
		AiClient:     aiClient,
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
