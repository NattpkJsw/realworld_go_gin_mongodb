package commentshandlers

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/NattpkJsw/real-world-api-go/config"
	"github.com/NattpkJsw/real-world-api-go/modules/comments"
	commentsusecases "github.com/NattpkJsw/real-world-api-go/modules/comments/commentsUsecases"
	"github.com/NattpkJsw/real-world-api-go/modules/entities"
	"github.com/gofiber/fiber/v2"
)

type commentsHandlersErrCode string

const (
	findCommentsErr   commentsHandlersErrCode = "comment-001"
	insertCommentsErr commentsHandlersErrCode = "comment-002"
	deleteCommentsErr commentsHandlersErrCode = "comment-003"
)

type ICommentsHandler interface {
	FindComments(c *fiber.Ctx) error
	InsertComment(c *fiber.Ctx) error
	DeleteComment(c *fiber.Ctx) error
}

type commentsHandler struct {
	cfg             config.IConfig
	commentsUsecase commentsusecases.ICommentUsecase
}

func CommentsHandler(cfg config.IConfig, commentsUsecase commentsusecases.ICommentUsecase) ICommentsHandler {
	return &commentsHandler{
		cfg:             cfg,
		commentsUsecase: commentsUsecase,
	}
}

func (h *commentsHandler) FindComments(c *fiber.Ctx) error {
	pathVariable := strings.TrimSpace(c.Params("slug"))
	slug, err := url.QueryUnescape(pathVariable)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findCommentsErr),
			err.Error(),
		).Res()
	}
	userID := c.Locals("userId").(int)

	commentsResult, err := h.commentsUsecase.FindComments(slug, userID)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findCommentsErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, commentsResult).Res()

}

func (h *commentsHandler) InsertComment(c *fiber.Ctx) error {
	userID := c.Locals("userId").(int)
	req := new(comments.JSONCommentCredential)
	pathVariable := strings.TrimSpace(c.Params("slug"))
	slug, err := url.QueryUnescape(pathVariable)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findCommentsErr),
			err.Error(),
		).Res()
	}
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertCommentsErr),
			err.Error(),
		).Res()
	}

	req.Comment.AuthorID = userID
	comment, err := h.commentsUsecase.InsertComment(slug, req.Comment)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(insertCommentsErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, comment).Res()

}

func (h *commentsHandler) DeleteComment(c *fiber.Ctx) error {
	userID := c.Locals("userId").(int)
	pathVariable := strings.Trim(c.Params("id"), " ")
	id, err := url.QueryUnescape(pathVariable)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findCommentsErr),
			err.Error(),
		).Res()
	}
	commentID, err := strconv.Atoi(id)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(deleteCommentsErr),
			err.Error(),
		).Res()
	}

	if err := h.commentsUsecase.DeleteComment(commentID, userID); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteCommentsErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusNoContent, nil).Res()
}
