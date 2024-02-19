package articleshandlers

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/NattpkJsw/real-world-api-go/config"
	"github.com/NattpkJsw/real-world-api-go/modules/articles"
	articlesusecases "github.com/NattpkJsw/real-world-api-go/modules/articles/articlesUsecases"
	"github.com/NattpkJsw/real-world-api-go/modules/entities"
	"github.com/gofiber/fiber/v2"
)

type articlesHandlersErrCode string

const (
	getSingleArticleErr  articlesHandlersErrCode = "article-001"
	getArticlesErr       articlesHandlersErrCode = "article-002"
	getArticlesFeedErr   articlesHandlersErrCode = "article-003"
	createArticleErr     articlesHandlersErrCode = "article-004"
	updateArticleErr     articlesHandlersErrCode = "article-005"
	deleteArticleErr     articlesHandlersErrCode = "article-006"
	favoriteArticleErr   articlesHandlersErrCode = "article-007"
	unfavoriteArticleErr articlesHandlersErrCode = "article-008"
	getTagslistErr       articlesHandlersErrCode = "article-009"
)

type IArticleshandler interface {
	GetSingleArticle(c *fiber.Ctx) error
	GetArticlesList(c *fiber.Ctx) error
	GetArticlesFeed(c *fiber.Ctx) error
	CreateArticle(c *fiber.Ctx) error
	UpdateArticle(c *fiber.Ctx) error
	DeleteArticle(c *fiber.Ctx) error
	FavoriteArticle(c *fiber.Ctx) error
	UnfavoriteArticle(c *fiber.Ctx) error
	GetTagsList(c *fiber.Ctx) error
}

type articlesHandler struct {
	cfg             config.IConfig
	articlesUsecase articlesusecases.IArticlesUsecase
}

func ArticlesHandler(cfg config.IConfig, articlesUsecase articlesusecases.IArticlesUsecase) IArticleshandler {
	return &articlesHandler{
		cfg:             cfg,
		articlesUsecase: articlesUsecase,
	}
}

func (h *articlesHandler) GetSingleArticle(c *fiber.Ctx) error {
	pathVariable := strings.TrimSpace(c.Params("slug"))
	slug, err := url.PathUnescape(pathVariable)

	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(getSingleArticleErr),
			err.Error(),
		).Res()
	}
	userId := c.Locals("userId").(int)

	article, err := h.articlesUsecase.GetSingleArticle(slug, userId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(getSingleArticleErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, article).Res()
}

func (h *articlesHandler) GetArticlesList(c *fiber.Ctx) error {
	req := &articles.ArticleFilter{}
	userId := c.Locals("userId").(int)

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(getArticlesErr),
			err.Error(),
		).Res()
	}

	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Offset <= 0 {
		req.Offset = 0
	}

	articlesOut, err := h.articlesUsecase.GetArticlesList(req, userId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(getArticlesErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, articlesOut).Res()

}

func (h *articlesHandler) GetArticlesFeed(c *fiber.Ctx) error {
	req := &articles.ArticleFeedFilter{}
	userId := c.Locals("userId").(int)

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(getArticlesFeedErr),
			err.Error(),
		).Res()
	}

	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Offset <= 0 {
		req.Offset = 0
	}
	articlesOut, err := h.articlesUsecase.GetArticlesFeed(req, userId)

	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(getArticlesFeedErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, articlesOut).Res()

}

func (h *articlesHandler) CreateArticle(c *fiber.Ctx) error {
	userID := c.Locals("userId").(int)
	req := new(articles.JSONArticleCredential)

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(createArticleErr),
			err.Error(),
		).Res()
	}
	req.Article.Author = userID

	article, err := h.articlesUsecase.CreateArticle(req.Article)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(createArticleErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, article).Res()
}

func (h *articlesHandler) UpdateArticle(c *fiber.Ctx) error {
	pathVariable := strings.TrimSpace(c.Params("slug"))
	slug, err := url.QueryUnescape(pathVariable)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateArticleErr),
			err.Error(),
		).Res()
	}
	userID := c.Locals("userId").(int)
	req := new(articles.JSONArticleCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateArticleErr),
			err.Error(),
		).Res()
	}
	req.Article.Slug = slug

	article, err := h.articlesUsecase.UpdateArticle(req.Article, userID)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(updateArticleErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, article).Res()
}

func (h *articlesHandler) DeleteArticle(c *fiber.Ctx) error {
	userId := c.Locals("userId").(int)
	pathVariable := strings.TrimSpace(c.Params("slug"))
	slug, err := url.QueryUnescape(pathVariable)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(deleteArticleErr),
			err.Error(),
		).Res()
	}
	if err := h.articlesUsecase.DeleteArticle(slug, userId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteArticleErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusNoContent, nil).Res()
}

func (h *articlesHandler) FavoriteArticle(c *fiber.Ctx) error {

	fmt.Println("favoriteeee")
	userID := c.Locals("userId").(int)
	pathVariable := strings.TrimSpace(c.Params("slug"))
	slug, err := url.QueryUnescape(pathVariable)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(favoriteArticleErr),
			err.Error(),
		).Res()
	}

	article, err := h.articlesUsecase.FavoriteArticle(slug, userID)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(favoriteArticleErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, article).Res()
}

func (h *articlesHandler) UnfavoriteArticle(c *fiber.Ctx) error {
	fmt.Println("unfavoriteeee")
	userID := c.Locals("userId").(int)
	pathVariable := strings.TrimSpace(c.Params("slug"))
	slug, err := url.QueryUnescape(pathVariable)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(unfavoriteArticleErr),
			err.Error(),
		).Res()
	}

	article, err := h.articlesUsecase.UnfavoriteArticle(slug, userID)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(unfavoriteArticleErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, article).Res()

}

func (h *articlesHandler) GetTagsList(c *fiber.Ctx) error {
	tagsList, err := h.articlesUsecase.GetTagsList()
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(getArticlesErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, tagsList).Res()
}
