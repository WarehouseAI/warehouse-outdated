package middleware

import (
	"os"
	"time"
	pg "warehouse/src/internal/database/postgresdb"
	"warehouse/src/internal/utils/httputils"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func User(logger *logrus.Logger) Middleware {
	userDatabase, err := pg.NewPostgresDatabase[pg.User](os.Getenv("DATA_DB_HOST"), os.Getenv("DATA_DB_USER"), os.Getenv("DATA_DB_PASSWORD"), os.Getenv("DATA_DB_NAME"), os.Getenv("DATA_DB_PORT"))
	if err != nil {
		panic(err)
	}

	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId")
		user, dbErr := userDatabase.GetOneByPreload(map[string]interface{}{"id": userId}, "FavoriteAi")

		if dbErr != nil {
			statusCode := httputils.InternalError
			logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("User middleware")
			return c.Status(statusCode).JSON(httputils.NewErrorResponse(statusCode, dbErr.Message))
		}

		if user == nil {
			statusCode := httputils.NotFound
			return c.Status(statusCode).JSON(httputils.NewErrorResponse(statusCode, "User not found."))
		}

		c.Locals("user", user)

		return c.Next()
	}

}
