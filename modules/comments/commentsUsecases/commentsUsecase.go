package commentsusecases

import (
	"github.com/NattpkJsw/real-world-api-go/config"
	articlesrepositories "github.com/NattpkJsw/real-world-api-go/modules/articles/articlesRepositories"
	"github.com/NattpkJsw/real-world-api-go/modules/comments"
	commentsrepositories "github.com/NattpkJsw/real-world-api-go/modules/comments/commentsRepositories"
)

type ICommentUsecase interface {
	FindComments(slug string, userID int) (*comments.JSONComment, error)
	InsertComment(slug string, req *comments.CommentCredential) (*comments.JSONSingleComment, error)
	DeleteComment(commentID, userID int) error
}

type commentUsecase struct {
	cfg                config.IConfig
	commentRepository  commentsrepositories.ICommentsRepository
	articlesRepository articlesrepositories.IArticlesRepository
}

func CommentUsecase(cfg config.IConfig, commentRepository commentsrepositories.ICommentsRepository, articlesRepository articlesrepositories.IArticlesRepository) ICommentUsecase {
	return &commentUsecase{
		cfg:                cfg,
		commentRepository:  commentRepository,
		articlesRepository: articlesRepository,
	}
}

func (u *commentUsecase) FindComments(slug string, userID int) (*comments.JSONComment, error) {
	articleID, err := u.articlesRepository.GetArticleIdBySlug(slug)
	if err != nil {
		return nil, err
	}
	commentOut, err := u.commentRepository.FindComments(articleID, userID)
	if err != nil {
		return nil, err
	}

	jsonComments := &comments.JSONComment{
		Comments: commentOut,
	}

	return jsonComments, nil
}

func (u *commentUsecase) InsertComment(slug string, req *comments.CommentCredential) (*comments.JSONSingleComment, error) {
	articleID, err := u.articlesRepository.GetArticleIdBySlug(slug)
	if err != nil {
		return nil, err
	}
	req.ArticleID = articleID
	singleComment, err := u.commentRepository.InsertComment(req)
	if err != nil {
		return nil, err
	}

	jsonComment := &comments.JSONSingleComment{
		Comment: *singleComment,
	}

	return jsonComment, nil

}

func (u *commentUsecase) DeleteComment(commentID, userID int) error {
	return u.commentRepository.DeleteComment(commentID, userID)

}
