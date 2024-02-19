package servers

import (
	articleshandlers "github.com/NattpkJsw/real-world-api-go/modules/articles/articlesHandlers"
	articlesrepositories "github.com/NattpkJsw/real-world-api-go/modules/articles/articlesRepositories"
	articlesusecases "github.com/NattpkJsw/real-world-api-go/modules/articles/articlesUsecases"
	"github.com/NattpkJsw/real-world-api-go/modules/middlewares"
)

type IArticleModule interface {
	Init()
	Repository() articlesrepositories.IArticlesRepository
	Usecase() articlesusecases.IArticlesUsecase
	Handler() articleshandlers.IArticleshandler
}

type articleModule struct {
	*moduleFactory
	repository articlesrepositories.IArticlesRepository
	usecase    articlesusecases.IArticlesUsecase
	handler    articleshandlers.IArticleshandler
}

func (m *moduleFactory) ArticlesModule() IArticleModule {
	articlesRepository := articlesrepositories.ArticlesRepository(m.server.db)
	articlesUsecase := articlesusecases.ArticlesUsecase(m.server.cfg, articlesRepository)
	articlesHandler := articleshandlers.ArticlesHandler(m.server.cfg, articlesUsecase)

	return &articleModule{
		moduleFactory: m,
		repository:    articlesRepository,
		usecase:       articlesUsecase,
		handler:       articlesHandler,
	}
}

func (a *articleModule) Init() {
	router := a.router.Group("/articles")

	router.Post("/", a.middle.JwtAuth(string(middlewares.WriteLevel)), a.handler.CreateArticle)
	router.Put("/:slug", a.middle.JwtAuth(string(middlewares.WriteLevel)), a.handler.UpdateArticle)
	router.Get("/:slug", a.middle.JwtAuth(string(middlewares.ReadLevel)), a.handler.GetSingleArticle)
	router.Get("/", a.middle.JwtAuth(string(middlewares.ReadLevel)), a.handler.GetArticlesList)
	router.Delete("/:slug/favorite", a.middle.JwtAuth(string(middlewares.WriteLevel)), a.handler.UnfavoriteArticle)
}

func (f *articleModule) Repository() articlesrepositories.IArticlesRepository { return f.repository }
func (f *articleModule) Usecase() articlesusecases.IArticlesUsecase           { return f.usecase }
func (f *articleModule) Handler() articleshandlers.IArticleshandler           { return f.handler }
