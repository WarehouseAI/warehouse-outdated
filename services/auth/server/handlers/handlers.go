package handlers

import (
	"fmt"
	"path/filepath"
	"warehouseai/auth/adapter/grpc/client/user"
	"warehouseai/auth/dataservice/sessiondata"
	"warehouseai/auth/dataservice/tokendata"
	e "warehouseai/internal/errors"
	"warehouseai/internal/gen"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	ResetTokenDB *tokendata.Database
	SessionDB    *sessiondata.Database
	Logger       *logrus.Logger
	AuthClient   *user.UserGrpcClient
}

func (pvd *Handler) RegisterHandler(c *fiber.Ctx) error {
	form, err := c.MultipartForm()

	if err != nil {
		response := e.NewErrorResponse(e.HttpBadRequest, err.Error())
		return c.Status(response.ErrorCode).JSON(response)
	}

	username := form.Value["username"][0]
	rawPicture, err := c.FormFile("picture")

	if err != nil {
		response := e.NewErrorResponse(e.HttpInternalError, err.Error())
		return c.Status(response.ErrorCode).JSON(response)
	}

	picture, err := rawPicture.Open()

	if err != nil {
		response := e.NewErrorResponse(e.HttpInternalError, err.Error())
		return c.Status(response.ErrorCode).JSON(response)
	}

	defer picture.Close()

	fileName := fmt.Sprintf("%s_avatar%s", username, filepath.Ext(rawPicture.Filename))

	link, svcErr := register.UploadAvatar(picture, fileName, pvd.logger, pvd.s3)

	if svcErr != nil {
		return c.Status(svcErr.ErrorCode).JSON(svcErr)
	}

	userInfo := &gen.CreateUserMsg{
		Username:  username,
		Firstname: form.Value["firstname"][0],
		Lastname:  form.Value["lastname"][0],
		Password:  form.Value["password"][0],
		Email:     form.Value["email"][0],
		Picture:   link,
	}

	userId, svcErr := register.Register(userInfo, pvd.userGateway, pvd.logger, pvd.ctx)

	if svcErr != nil {
		return c.Status(svcErr.ErrorCode).JSON(svcErr)
	}

	return c.Status(fiber.StatusCreated).JSON(userId)
}

func (pvd *AuthServiceProvider) LoginHandler(c *fiber.Ctx) error {
	var creds login.Request

	if err := c.BodyParser(&creds); err != nil {
		statusCode := httputils.BadRequest
		return c.Status(statusCode).JSON(httputils.NewErrorResponse(statusCode, "Invalid request body"))
	}

	session, err := login.Login(&creds, pvd.userGateway, pvd.sessionDatabase, pvd.logger, pvd.ctx)

	if err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "sessionId",
		Value:    session.ID,
		SameSite: fiber.CookieSameSiteNoneMode,
		Secure:   true,
	})

	return c.SendStatus(fiber.StatusOK)
}

func (pvd *AuthServiceProvider) PasswordReset(c *fiber.Ctx) error {
	resetTokenId := c.Query("token_id")
	var request gen.ResetPasswordRequest

	if err := c.BodyParser(&request); err != nil {
		statusCode := httputils.BadRequest
		return c.Status(statusCode).JSON(httputils.NewErrorResponse(statusCode, "Invalid request body"))
	}

	response, err := recovery.PasswordReset(&request, resetTokenId, pvd.userGateway, pvd.tokenDatabase, pvd.logger, pvd.ctx)

	if err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (pvd *AuthServiceProvider) VerifyReset(c *fiber.Ctx) error {
	verificationCode := c.Query("verification")
	tokenId := c.Query("token_id")

	resetToken, err := recovery.VerifyResetCode(verificationCode, tokenId, pvd.tokenDatabase, pvd.logger, pvd.ctx)

	if err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(resetToken)
}

func (pvd *AuthServiceProvider) SendResetHandler(c *fiber.Ctx) error {
	var request recovery.ResetAttemptRequest

	if err := c.BodyParser(&request); err != nil {
		statusCode := httputils.BadRequest
		return c.Status(statusCode).JSON(httputils.NewErrorResponse(statusCode, "Invalid request body"))
	}

	resetToken, err := recovery.SendResetEmail(request, pvd.tokenDatabase, pvd.userGateway, pvd.logger, pvd.ctx)

	if err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(resetToken)
}

func (pvd *AuthServiceProvider) LogoutHandler(c *fiber.Ctx) error {
	sessionId := c.Cookies("sessionId")

	if sessionId == "" {
		return c.Status(httputils.Unauthorized).JSON(httputils.NewErrorResponse(httputils.Unauthorized, "Empty session key."))
	}

	if err := logout.Logout(sessionId, pvd.sessionDatabase, pvd.logger, pvd.ctx); err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}
	c.ClearCookie("sessionId")

	return c.SendStatus(fiber.StatusOK)
}

func (api *AuthServiceProvider) WhoAmIHandler(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}
