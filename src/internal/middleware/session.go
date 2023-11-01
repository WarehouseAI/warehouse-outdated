package middleware

import (
	"context"
	"os"
	"time"
	db "warehouse/src/internal/database"
	r "warehouse/src/internal/database/redisdb"
	"warehouse/src/internal/utils/httputils"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Middleware func(c *fiber.Ctx) error

func Session(logger *logrus.Logger) Middleware {
	sessionDatabase := r.NewRedisDatabase(os.Getenv("SESSION_DB_HOST"), os.Getenv("SESSION_DB_PORT"), os.Getenv("SESSION_DB_PASSWORD"))

	return func(c *fiber.Ctx) error {
		sessionId := c.Cookies("sessionId")

		if sessionId == "" {
			return c.Status(httputils.Unauthorized).JSON(httputils.NewErrorResponse(httputils.Unauthorized, "Empty session key."))
		}

		newSession, dbErr := sessionDatabase.Update(context.Background(), sessionId)

		if dbErr != nil {
			logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Session middleware")

			if dbErr.ErrorType == db.NotFound {
				errorResponse := httputils.NewErrorResponseFromDBError(dbErr.ErrorType, "Session has expired.")
				return c.Status(errorResponse.ErrorCode).JSON(errorResponse)
			}

			errorResponse := httputils.NewErrorResponseFromDBError(dbErr.ErrorType, dbErr.Message)
			return c.Status(errorResponse.ErrorCode).JSON(errorResponse)
		}

		c.Locals("userId", newSession.Payload.UserId)
		c.Cookie(&fiber.Cookie{
			Name:     "sessionId",
			Value:    newSession.ID,
			SameSite: fiber.CookieSameSiteNoneMode,
			Secure:   true,
		})

		return c.Next()
	}

}
