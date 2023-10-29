package middleware

import (
	"context"
	"time"
	db "warehouse/src/internal/database"
	r "warehouse/src/internal/database/redisdb"
	"warehouse/src/internal/utils/httputils"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Middleware func(c *fiber.Ctx) error

type SessionProvider interface {
	Get(ctx context.Context, sessionId string) (*r.Session, *db.DBError)
	Update(ctx context.Context, sessionId string) (*r.Session, *db.DBError)
}

func Session(sessionProvider SessionProvider, logger *logrus.Logger) Middleware {
	return func(c *fiber.Ctx) error {
		sessionId := c.Cookies("sessionId")

		if sessionId == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(httputils.NewErrorResponse(httputils.Abort, "Empty session key."))
		}

		session, dbErr := sessionProvider.Get(context.Background(), sessionId)

		if dbErr != nil {
			logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Session middleware")
			return c.Status(fiber.StatusInternalServerError).JSON(httputils.NewErrorResponse(httputils.ServerError, dbErr.Message))
		}

		if session == nil {
			c.ClearCookie("sessionId")
			return c.Status(fiber.StatusForbidden).JSON(httputils.NewErrorResponse(httputils.Abort, "Session has expired"))
		}

		newSession, dbErr := sessionProvider.Update(context.Background(), sessionId)

		if dbErr != nil {
			logger.WithFields(logrus.Fields{"time": time.Now(), "error": dbErr.Payload}).Info("Session middleware")
			return c.Status(fiber.StatusInternalServerError).JSON(httputils.NewErrorResponse(httputils.ServerError, dbErr.Message))
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
