package server

import (
	"warehouseai/stat/dataservice/statdata"
	h "warehouseai/stat/server/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
)

func StartServer(port string, statDb *statdata.Database, logger *logrus.Logger) error {
	handler := newHttpHandler(statDb, logger)
	app := fiber.New()
	app.Use(setupCORS())

	//Todo: сделать ручки для принятия проведённого времени на сайте конкретным порльзователем
	//Todo: сделать ручку для изменения флага дневной активности пользователя
	//Todo: сделать ручку увеличения количества просмотров каждой нейронки
	//Todo: сделать ручку для трекинга количестваа активной аудитории в данный момент времени
	//Todo: перенести логику трекинга количества использований каждой нейронки сюда
	route := app.Group("/stat")
	route.Get("user/num", handler.GetNumOfUsersHandler)
	route.Get("developer/num", handler.GetNumOfDevelopersHandler)
	route.Get("ai/uses", handler.GetNumOfAiUsesHandler)

	return app.Listen(port)
}

func newHttpHandler(db *statdata.Database, logger *logrus.Logger) *h.Handler {
	return &h.Handler{
		DB:     db,
		Logger: logger,
	}
}

// TODO: Подумать над тем, что может нужно что-то поменять
func setupCORS() func(*fiber.Ctx) error {
	return cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
		AllowOrigins:     "http://localhost:3000, https://warehouse-ai-frontend.vercel.app",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	})
}
