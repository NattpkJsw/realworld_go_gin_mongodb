package usershandlers

import (
	"github.com/NattpkJsw/real-world-api-go/config"
	"github.com/NattpkJsw/real-world-api-go/modules/entities"
	"github.com/NattpkJsw/real-world-api-go/modules/users"
	usersUsecases "github.com/NattpkJsw/real-world-api-go/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type userHandlersErrCode string

const (
	signUpErr          userHandlersErrCode = "users-001"
	logInErr           userHandlersErrCode = "users-002"
	refreshPassportErr userHandlersErrCode = "users-003"
	logOutErr          userHandlersErrCode = "users-004"
	getUserErr         userHandlersErrCode = "users-005"
	UpdateUserErr      userHandlersErrCode = "users-006"
)

type IUsersHandler interface {
	SignUp(c *fiber.Ctx) error
	LogIn(c *fiber.Ctx) error
	LogOut(c *fiber.Ctx) error
	GetUser(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
}

type usersHandler struct {
	cfg          config.IConfig
	usersUsecase usersUsecases.IUsersUsecase
}

func UsersHandler(cfg config.IConfig, usersUsecase usersUsecases.IUsersUsecase) IUsersHandler {
	return &usersHandler{
		cfg:          cfg,
		usersUsecase: usersUsecase,
	}
}

func (h *usersHandler) SignUp(c *fiber.Ctx) error {
	// Request body parser
	req := &users.RegisterReq{}

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpErr),
			err.Error(),
		).Res()
	}

	// Email validation
	if !req.User.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpErr),
			"email pattern is invalid",
		).Res()
	}
	// Insert
	result, err := h.usersUsecase.InsertCustomer(req.User)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpErr),
				err.Error(),
			).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code, //500
				string(signUpErr),
				err.Error(),
			).Res()
		}
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}

func (h *usersHandler) LogIn(c *fiber.Ctx) error {
	req := new(users.UserSignin)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(logInErr),
			err.Error(),
		).Res()
	}

	passport, err := h.usersUsecase.GetPassport(&req.User)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(logInErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *usersHandler) LogOut(c *fiber.Ctx) error {
	req := new(users.OauthToken)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(logOutErr),
			err.Error(),
		).Res()
	}
	if err := h.usersUsecase.DeleteOauth(req.AccessToken); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(logOutErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}

func (h *usersHandler) GetUser(c *fiber.Ctx) error {
	token := c.Locals("accessToken").(string)
	// Get profile
	result, err := h.usersUsecase.GetUser(token)
	if err != nil {
		switch err.Error() {
		case "get user failed: sql: no rows in result set":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(getUserErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(getUserErr),
				err.Error(),
			).Res()
		}
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()
}

func (h *usersHandler) UpdateUser(c *fiber.Ctx) error {
	userId := c.Locals("userId").(int)
	token := c.Locals("accessToken").(string)
	req := new(users.UserCredentialInput)

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(UpdateUserErr),
			err.Error(),
		).Res()
	}
	req.User.AccessToken = token
	req.User.Id = userId
	ret, err := h.usersUsecase.UpdateUser(&req.User)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(UpdateUserErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, ret).Res()
}
