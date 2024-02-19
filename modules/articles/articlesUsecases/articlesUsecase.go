package articlesusecases

import (
	"github.com/NattpkJsw/real-world-api-go/config"
	"github.com/NattpkJsw/real-world-api-go/modules/articles"
	articlesrepositories "github.com/NattpkJsw/real-world-api-go/modules/articles/articlesRepositories"
)

type IArticlesUsecase interface {
	GetSingleArticle(slug string, userId int) (*articles.JSONArticle, error)
	GetArticlesList(req *articles.ArticleFilter, userId int) (*articles.ArticleList, error)
	GetArticlesFeed(req *articles.ArticleFeedFilter, userId int) (*articles.ArticleList, error)
	CreateArticle(req *articles.ArticleCredential) (*articles.JSONArticle, error)
	UpdateArticle(req *articles.ArticleCredential, userID int) (*articles.JSONArticle, error)
	DeleteArticle(slug string, userID int) error
	FavoriteArticle(slug string, userID int) (*articles.JSONArticle, error)
	UnfavoriteArticle(slug string, userID int) (*articles.JSONArticle, error)
	GetTagsList() (*articles.TagList, error)
}

type articlesUsecase struct {
	cfg                config.IConfig
	articlesRepository articlesrepositories.IArticlesRepository
}

func ArticlesUsecase(cfg config.IConfig, articlesRepository articlesrepositories.IArticlesRepository) IArticlesUsecase {
	return &articlesUsecase{
		cfg:                cfg,
		articlesRepository: articlesRepository,
	}
}

func (u *articlesUsecase) GetSingleArticle(slug string, userId int) (*articles.JSONArticle, error) {
	articleId, err := u.articlesRepository.GetArticleIdBySlug(slug)
	if err != nil {
		return nil, err
	}

	article, err := u.articlesRepository.GetSingleArticle(articleId, userId)
	if err != nil {
		return nil, err
	}
	jsonArticle := &articles.JSONArticle{
		Article: article,
	}
	return jsonArticle, nil
}

func (u *articlesUsecase) GetArticlesList(req *articles.ArticleFilter, userId int) (*articles.ArticleList, error) {
	articleList, count, err := u.articlesRepository.GetArticlesList(req, userId)

	articlesOut := &articles.ArticleList{
		Article:       articleList,
		ArticlesCount: count,
	}

	return articlesOut, err
}

func (u *articlesUsecase) GetArticlesFeed(req *articles.ArticleFeedFilter, userId int) (*articles.ArticleList, error) {
	input := &articles.ArticleFilter{
		Limit:  req.Limit,
		Offset: req.Offset,
		IsFeed: true,
	}
	articleList, count, err := u.articlesRepository.GetArticlesList(input, userId)

	articlesOut := &articles.ArticleList{
		Article:       articleList,
		ArticlesCount: count,
	}

	return articlesOut, err
}

func (u *articlesUsecase) CreateArticle(req *articles.ArticleCredential) (*articles.JSONArticle, error) {
	article, err := u.articlesRepository.CreateArticle(req)
	if err != nil {
		return nil, err
	}
	jsonArticle := &articles.JSONArticle{
		Article: article,
	}

	return jsonArticle, nil

}

func (u *articlesUsecase) UpdateArticle(req *articles.ArticleCredential, userID int) (*articles.JSONArticle, error) {
	articleID, err := u.articlesRepository.GetArticleIdBySlug(req.Slug)
	if err != nil {
		return nil, err
	}
	req.Id = articleID

	article, err := u.articlesRepository.UpdateArticle(req, userID)
	if err != nil {
		return nil, err
	}
	jsonArticle := &articles.JSONArticle{
		Article: article,
	}

	return jsonArticle, nil
}

func (u *articlesUsecase) DeleteArticle(slug string, userID int) error {
	artcleID, err := u.articlesRepository.GetArticleIdBySlug(slug)
	if err != nil {
		return err
	}
	return u.articlesRepository.DeleteArticle(artcleID, userID)
}

func (u *articlesUsecase) FavoriteArticle(slug string, userID int) (*articles.JSONArticle, error) {
	articleID, err := u.articlesRepository.GetArticleIdBySlug(slug)
	if err != nil {
		return nil, err
	}
	articleOut, err := u.articlesRepository.FavoriteArticle(userID, articleID)
	if err != nil {
		return nil, err
	}
	jsonArticle := &articles.JSONArticle{
		Article: articleOut,
	}
	return jsonArticle, nil
}

func (u *articlesUsecase) UnfavoriteArticle(slug string, userID int) (*articles.JSONArticle, error) {
	articleID, err := u.articlesRepository.GetArticleIdBySlug(slug)
	if err != nil {
		return nil, err
	}
	articleOut, err := u.articlesRepository.UnfavoriteArticle(userID, articleID)
	if err != nil {
		return nil, err
	}
	jsonArticle := &articles.JSONArticle{
		Article: articleOut,
	}
	return jsonArticle, nil
}

func (u *articlesUsecase) GetTagsList() (*articles.TagList, error) {
	return u.articlesRepository.GetTagsList()
}
