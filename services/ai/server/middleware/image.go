package middleware

import (
	"warehouseai/ai/dataservice"
	e "warehouseai/ai/errors"
	"warehouseai/ai/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func Image(logger *logrus.Logger, picture dataservice.PictureInterface) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		pic, err := c.FormFile("image")

		if err != nil {
			resp := e.NewErrorResponse(e.HttpBadRequest, "Image must be provided in image")
			return c.Status(resp.ErrorCode).JSON(resp)
		}

		url, svcErr := service.UploadImage(pic, picture, logger)

		if svcErr != nil {
			return c.Status(svcErr.ErrorCode).JSON(svcErr)
		}

		c.Locals("imageUrl", url)
		return c.Next()
	}
}
