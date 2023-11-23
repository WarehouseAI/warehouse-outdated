package middleware

import (
	"warehouseai/auth/dataservice"
	"warehouseai/auth/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func Image(logger *logrus.Logger, picture dataservice.PictureInterface) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		pic, err := c.FormFile("image")

		if err != nil {
			return c.Next()
		}

		url, svcErr := service.UploadImage(pic, picture, logger)

		if svcErr != nil {
			return c.Status(svcErr.ErrorCode).JSON(svcErr)
		}

		c.Locals("imageUrl", url)
		return c.Next()
	}
}
