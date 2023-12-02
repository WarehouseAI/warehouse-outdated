package handlers

import (
	"warehouseai/auth/adapter/broker/mail"
	"warehouseai/auth/adapter/grpc/client/user"
	"warehouseai/auth/dataservice/picturedata"
	"warehouseai/auth/dataservice/sessiondata"
	"warehouseai/auth/dataservice/tokendata"
	e "warehouseai/auth/errors"
	"warehouseai/auth/model"
	"warehouseai/auth/service"
	"warehouseai/auth/service/login"
	"warehouseai/auth/service/register"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	ResetTokenDB        *tokendata.Database[model.ResetToken]
	VerificationTokenDB *tokendata.Database[model.VerificationToken]
	SessionDB           *sessiondata.Database
	PictureStorage      *picturedata.Storage
	MailProducer        *mail.MailProducer
	Logger              *logrus.Logger
	UserClient          *user.UserGrpcClient
}

func (h *Handler) RegisterHandler(c *fiber.Ctx) error {
	var req register.RegisterRequest
	form, err := c.MultipartForm()

	if err != nil {
		response := e.NewErrorResponse(e.HttpBadRequest, err.Error())
		return c.Status(response.ErrorCode).JSON(response)
	}

	imageUrl := c.Locals("imageUrl")

	if imageUrl == nil {
		req.Image = ""
	} else {
		req.Image = imageUrl.(string)
	}

	req.Username = form.Value["username"][0]
	req.Firstname = form.Value["firstname"][0]
	req.Lastname = form.Value["lastname"][0]
	req.Password = form.Value["password"][0]
	req.Email = form.Value["email"][0]
	req.ViaGoogle = false

	userId, svcErr := register.Register(&req, h.UserClient, h.VerificationTokenDB, h.MailProducer, h.Logger)

	if svcErr != nil {
		return c.Status(svcErr.ErrorCode).JSON(svcErr)
	}

	return c.Status(fiber.StatusCreated).JSON(userId)
}

func (h *Handler) LoginHandler(c *fiber.Ctx) error {
	var request login.LoginRequest

	if err := c.BodyParser(&request); err != nil {
		response := e.NewErrorResponse(e.HttpBadRequest, "Invalid request body")
		return c.Status(response.ErrorCode).JSON(response)
	}

	response, session, err := login.Login(&request, h.UserClient, h.SessionDB, h.Logger)

	if err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "sessionId",
		Value:    session.ID,
		SameSite: fiber.CookieSameSiteNoneMode,
		Secure:   true,
	})

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *Handler) RegisterVerifyHandler(c *fiber.Ctx) error {
	token := c.Query("token")
	user := c.Query("user")

	request := register.RegisterVerifyRequest{
		UserId: user,
		Token:  token,
	}

	response, err := register.RegisterVerify(request, h.UserClient, h.VerificationTokenDB, h.Logger)

	if err != nil {
		return c.Status(err.ErrorCode).JSON(err.ErrorMessage)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *Handler) PasswordReset(c *fiber.Ctx) error {
	resetTokenId := c.Query("token_id")
	var request service.ResetConfirmRequest

	if err := c.BodyParser(&request); err != nil {
		response := e.NewErrorResponse(e.HttpBadRequest, "Invalid request body")
		return c.Status(response.ErrorCode).JSON(response)
	}

	response, err := service.PasswordReset(&request, resetTokenId, h.UserClient, h.ResetTokenDB, h.Logger)

	if err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *Handler) VerifyReset(c *fiber.Ctx) error {
	verificationCode := c.Query("verification")
	resetTokenId := c.Query("token_id")

	resetToken, err := service.VerifyResetCode(verificationCode, resetTokenId, h.ResetTokenDB, h.Logger)

	if err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(resetToken)
}

func (h *Handler) SendResetHandler(c *fiber.Ctx) error {
	var request service.ResetAttemptRequest

	if err := c.BodyParser(&request); err != nil {
		response := e.NewErrorResponse(e.HttpBadRequest, "Invalid request body")
		return c.Status(response.ErrorCode).JSON(response)
	}

	resetToken, err := service.SendResetEmail(request, h.ResetTokenDB, h.UserClient, h.MailProducer, h.Logger)

	if err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(resetToken)
}

func (h *Handler) LogoutHandler(c *fiber.Ctx) error {
	sessionId := c.Cookies("sessionId")

	if sessionId == "" {
		response := e.NewErrorResponse(e.HttpUnauthorized, "Empty session key")
		return c.Status(response.ErrorCode).JSON(response)
	}

	if err := service.Logout(sessionId, h.SessionDB, h.Logger); err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}
	c.ClearCookie("sessionId")

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) WhoAmIHandler(c *fiber.Ctx) error {
	sessionId := c.Cookies("sessionId")

	if sessionId == "" {
		response := e.NewErrorResponse(e.HttpUnauthorized, "Empty session key")
		return c.Status(response.ErrorCode).JSON(response)
	}

	_, newSession, err := service.Authenticate(sessionId, h.SessionDB, h.Logger)

	if err != nil {
		return c.Status(err.ErrorCode).JSON(err)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "sessionId",
		Value:    newSession.ID,
		SameSite: fiber.CookieSameSiteNoneMode,
		Secure:   true,
	})

	return c.SendStatus(fiber.StatusOK)
}
