package server

import (
	"warehouseai/ai/adapter/grpc/client/auth"
	"warehouseai/ai/dataservice/aidata"
	"warehouseai/ai/dataservice/commanddata"
	"warehouseai/ai/server/handlers/ai"
	"warehouseai/ai/server/handlers/commands"
	"warehouseai/ai/server/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
)

func StartServer(port string, aiDB *aidata.Database, commandDB *commanddata.Database, logger *logrus.Logger) error {
	aiHandler := newHttpAiHandler(aiDB, logger)
	commandHandler := newHttpCommandHandler(commandDB, aiDB, logger)
	app := fiber.New()
	app.Use(setupCORS())

	sessionMw := middleware.Session(logger, aiHandler.AuthClient)

	route := app.Group("/ai")
	route.Post("/create/generate", sessionMw, aiHandler.CreateAiWithoutKeyHandler)
	route.Post("/create/exist", sessionMw, aiHandler.CreateAiWithKeyHandler)
	route.Get("/get/:id", aiHandler.GetAIHandler)
	route.Get("/get/many", aiHandler.GetAisHandler)
	route.Post("/command/create", sessionMw, commandHandler.CreateCommandHandler)
	route.Post("/command/execute", sessionMw, commandHandler.ExecuteCommandHandler)

	return app.Listen(port)
}

func newHttpAiHandler(db *aidata.Database, logger *logrus.Logger) *ai.Handler {
	authClient := auth.NewAuthGrpcClient("auth:8041")

	return &ai.Handler{
		DB:         db,
		Logger:     logger,
		AuthClient: authClient,
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

func setupCORS() func(*fiber.Ctx) error {
	return cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
		AllowOrigins:     "http://localhost:3000, https://warehouse-ai-frontend.vercel.app",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	})
}
