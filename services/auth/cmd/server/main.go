package server

import (
	"warehouseai/auth/adapter/broker"
	"warehouseai/auth/adapter/grpc/client/user"
	"warehouseai/auth/dataservice/picturedata"
	"warehouseai/auth/dataservice/sessiondata"
	"warehouseai/auth/dataservice/tokendata"
	m "warehouseai/auth/model"
	h "warehouseai/auth/server/handlers"
	"warehouseai/auth/server/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
)

func StartServer(
	port string,
	resetTokenDB *tokendata.Database[m.ResetToken],
	verificationTokenDB *tokendata.Database[m.VerificationToken],
	sessionDB *sessiondata.Database,
	pictureStorage *picturedata.Storage,
	mailProducer *broker.Broker,
	logger *logrus.Logger,
) error {

	handler := newHttpHandler(resetTokenDB, verificationTokenDB, sessionDB, pictureStorage, mailProducer, logger)
	app := fiber.New()
	app.Use(setupCORS())

	route := app.Group("/auth")

	pictureMw := middleware.Image(logger, pictureStorage)

	route.Post("/register", pictureMw, handler.RegisterHandler)
	route.Get("/register/confirm", handler.RegisterVerifyHandler)
	route.Post("/login", handler.LoginHandler)
	route.Post("/reset/request", handler.SendResetHandler)
	route.Get("/reset/verify", handler.VerifyReset)
	route.Post("/reset/confirm", handler.PasswordReset)
	route.Delete("/logout", handler.LogoutHandler)
	route.Get("/whoami", handler.WhoAmIHandler)

	return app.Listen(port)
}

func newHttpHandler(
	resetTokenDB *tokendata.Database[m.ResetToken],
	verificationTokenDB *tokendata.Database[m.VerificationToken],
	sessionDB *sessiondata.Database,
	pictureStorage *picturedata.Storage,
	mailProducer *broker.Broker,
	logger *logrus.Logger,
) *h.Handler {

	userClient := user.NewUserGrpcClient("user:8001")

	return &h.Handler{
		ResetTokenDB:        resetTokenDB,
		VerificationTokenDB: verificationTokenDB,
		SessionDB:           sessionDB,
		PictureStorage:      pictureStorage,
		Broker:              mailProducer,
		Logger:              logger,
		UserClient:          userClient,
	}
}

func setupCORS() func(*fiber.Ctx) error {
	return cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
		AllowOrigins:     "http://localhost:3000, https://warehouse-ai-frontend.vercel.app, https://warehousai.com",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	})
}
