package profileshandlers

import (
	"strings"

	"github.com/NattpkJsw/real-world-api-go/config"
	"github.com/NattpkJsw/real-world-api-go/modules/entities"
	profilesusecases "github.com/NattpkJsw/real-world-api-go/modules/profiles/profilesUsecases"
	"github.com/gofiber/fiber/v2"
)

type profileHandlersErrCode string

const (
	getProfileErr   profileHandlersErrCode = "profiles-001"
	followUserErr   profileHandlersErrCode = "profiles-002"
	unfollowUserErr profileHandlersErrCode = "profiles-003"
)

type IProfileHandler interface {
	GetProfile(c *fiber.Ctx) error
	FollowUser(c *fiber.Ctx) error
	UnfollowUser(c *fiber.Ctx) error
}

type profileHandler struct {
	cfg            config.IConfig
	profileUsecase profilesusecases.IProfilesUsecase
}

func ProfileHandler(cfg config.IConfig, profileUsecase profilesusecases.IProfilesUsecase) IProfileHandler {
	return &profileHandler{
		cfg:            cfg,
		profileUsecase: profileUsecase,
	}
}

func (h *profileHandler) GetProfile(c *fiber.Ctx) error {
	username := strings.Trim(c.Params("username"), " ")
	result, err := h.profileUsecase.GetProfile(username, c.Locals("userId").(int))
	if err != nil {
		switch err.Error() {
		case "get user failed: sql: no rows in result set":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(getProfileErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(getProfileErr),
				err.Error(),
			).Res()
		}
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()
}

func (h *profileHandler) FollowUser(c *fiber.Ctx) error {
	username := strings.Trim(c.Params("username"), " ")
	result, err := h.profileUsecase.FollowUser(username, c.Locals("userId").(int))
	if err != nil {
		switch err.Error() {
		case "get user failed: sql: no rows in result set":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(followUserErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(followUserErr),
				err.Error(),
			).Res()
		}
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()
}

func (h *profileHandler) UnfollowUser(c *fiber.Ctx) error {
	username := strings.Trim(c.Params("username"), " ")
	result, err := h.profileUsecase.UnfollowUser(username, c.Locals("userId").(int))
	if err != nil {
		switch err.Error() {
		case "get user failed: sql: no rows in result set":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(unfollowUserErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(unfollowUserErr),
				err.Error(),
			).Res()
		}
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()
}
