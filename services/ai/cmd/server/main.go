package server

import (
	"warehouseai/ai/adapter/grpc/client/auth"
	"warehouseai/ai/adapter/grpc/client/user"
	"warehouseai/ai/dataservice/aidata"
	"warehouseai/ai/dataservice/commanddata"
	"warehouseai/ai/dataservice/picturedata"
	"warehouseai/ai/dataservice/ratingdata"
	"warehouseai/ai/server/handlers/ai"
	"warehouseai/ai/server/handlers/commands"
	"warehouseai/ai/server/handlers/rating"
	"warehouseai/ai/server/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
)

func StartServer(port string, ratingDB *ratingdata.Database, aiDB *aidata.Database, commandDB *commanddata.Database, pictureStorage *picturedata.Storage, logger *logrus.Logger) error {
	aiHandler := newHttpAiHandler(aiDB, pictureStorage, logger)
	commandHandler := newHttpCommandHandler(commandDB, aiDB, logger)
	ratingHandler := newRatingHandler(ratingDB, aiDB, logger)
	app := fiber.New()
	app.Use(setupCORS())

	sessionStrictMw := middleware.SessionStrict(logger, aiHandler.AuthClient)
	sessionMw := middleware.Session(logger, aiHandler.AuthClient)
	pictureMW := middleware.Image(logger, pictureStorage)

	route := app.Group("/ai")
	route.Post("/create/generate", sessionStrictMw, pictureMW, aiHandler.CreateAiWithoutKeyHandler)
	route.Post("/create/exist", sessionStrictMw, pictureMW, aiHandler.CreateAiWithKeyHandler)
	route.Get("/get", sessionMw, aiHandler.GetAIHandler)
	route.Get("/get/many", aiHandler.GetAisHandler)
	route.Get("/search", aiHandler.SearchHandler)
	route.Post("/command/create", sessionStrictMw, commandHandler.CreateCommandHandler)
	route.Post("/command/execute", sessionStrictMw, commandHandler.ExecuteCommandHandler)
	route.Get("/rating/get", ratingHandler.GetAiRatingHandler)
	route.Post("/rating/set", sessionStrictMw, ratingHandler.SetRatingForAiHandler)

	return app.Listen(port)
}

func newHttpAiHandler(db *aidata.Database, pictureStorage *picturedata.Storage, logger *logrus.Logger) *ai.Handler {
	authClient := auth.NewAuthGrpcClient("auth:8041")
	userClient := user.NewUserGrpcClient("user:8001")

	return &ai.Handler{
		DB:             db,
		Logger:         logger,
		PictureStorage: pictureStorage,
		UserClient:     userClient,
		AuthClient:     authClient,
	}
}

func newHttpCommandHandler(commandDB *commanddata.Database, aiDB *aidata.Database, logger *logrus.Logger) *commands.Handler {
	authClient := auth.NewAuthGrpcClient("auth:8041")

	return &commands.Handler{
		CommandDB:  commandDB,
		AiDB:       aiDB,
		Logger:     logger,
		AuthClient: authClient,
	}
}

func newRatingHandler(ratingDB *ratingdata.Database, aiDB *aidata.Database, logger *logrus.Logger) *rating.Handler {
	return &rating.Handler{
		RatingRepository: ratingDB,
		AiRepository:     aiDB,
		Logger:           logger,
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
