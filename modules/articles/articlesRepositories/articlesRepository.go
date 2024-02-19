package articlesrepositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/NattpkJsw/real-world-api-go/modules/articles"
	articlespatterns "github.com/NattpkJsw/real-world-api-go/modules/articles/articlesPatterns"
	"github.com/jmoiron/sqlx"
)

type IArticlesRepository interface {
	GetSingleArticle(articleId int, userId int) (*articles.Article, error)
	GetArticlesList(req *articles.ArticleFilter, userId int) ([]*articles.Article, int, error)
	GetArticleIdBySlug(slug string) (int, error)
	CreateArticle(req *articles.ArticleCredential) (*articles.Article, error)
	UpdateArticle(req *articles.ArticleCredential, userID int) (*articles.Article, error)
	DeleteArticle(articleID, userID int) error
	FavoriteArticle(userID, articleID int) (*articles.Article, error)
	UnfavoriteArticle(userID, articleID int) (*articles.Article, error)
	GetTagsList() (*articles.TagList, error)
}

type articlesRepository struct {
	db *sqlx.DB
}

func ArticlesRepository(db *sqlx.DB) IArticlesRepository {
	return &articlesRepository{
		db: db,
	}
}

func (r *articlesRepository) GetSingleArticle(articleId int, userId int) (*articles.Article, error) {
	query := `
	SELECT
		to_jsonb("ar")
	FROM(
			SELECT
			"a"."slug",
			"a"."title",
			"a"."description",
			"a"."body",
			(
				SELECT coalesce(array_to_json(array_agg("t"."name")),'[]'::json)
				FROM "article_tags" "at"
				JOIN "tags" AS "t" ON "at"."tag_id" = "t"."id"
				WHERE "a"."id" = "at"."article_id"
			) AS "taglist",
			"a"."createdat",
			"a"."updatedat",
			(
				SELECT
				CASE WHEN EXISTS(
					SELECT 1
					FROM "article_favorites" "af"
					WHERE "af"."user_id" = $2 AND "af"."article_id" = "a"."id"
				) THEN TRUE ELSE FALSE END
			) AS "favorited",
			(
				SELECT COUNT(*)
				FROM "article_favorites" "af"
				WHERE "af"."article_id" = "a"."id"
			) AS "favoritesCount",
			(
				SELECT 
					json_build_object(
						'username', "u"."username",
						'bio', "u"."bio",
						'image', "u"."image",
						'following',
						CASE 
							WHEN EXISTS (
								SELECT 1
								FROM "user_follows" "uf"
								WHERE "a"."author_id" = "uf"."following_id"  AND "uf"."follower_id" = $2
							) THEN TRUE 
							ELSE FALSE 
						END
					)
				FROM "users" "u"
				WHERE "a"."author_id" = "u"."id"
			) AS "author"
			FROM "articles" "a"
			WHERE "a"."id" = $1
			LIMIT 1
	) AS "ar";`

	articleBytes := make([]byte, 0)
	article := new(articles.Article)

	if err := r.db.Get(&articleBytes, query, articleId, userId); err != nil {
		return nil, fmt.Errorf("get article failed: %v", err)
	}
	if err := json.Unmarshal(articleBytes, &article); err != nil {
		return nil, fmt.Errorf("unmarshal article failed: %v", err)
	}
	return article, nil
}

func (r *articlesRepository) GetArticleIdBySlug(slug string) (int, error) {
	query := `
	SELECT "a"."id"
	FROM "articles" "a"
	WHERE "a"."slug" = $1`

	var id int
	if err := r.db.Get(&id, query, slug); err != nil {
		return 0, fmt.Errorf("get articleID failed: %v", err)
	}
	return id, nil
}

func (r *articlesRepository) GetArticlesList(req *articles.ArticleFilter, userId int) ([]*articles.Article, int, error) {
	builder := articlespatterns.FindArticleBuilder(r.db, req)
	engineer := articlespatterns.FindProductEngineer(builder)

	result, err := engineer.FindArticle(userId).Result()
	count := len(result)

	return result, count, err
}

func (r *articlesRepository) CreateArticle(req *articles.ArticleCredential) (*articles.Article, error) {
	builder := articlespatterns.AddArticleBuilder(r.db, req)
	articleId, err := articlespatterns.AddArticleEngineer(builder).AddArticle()
	if err != nil {
		return nil, err
	}

	article, err := r.GetSingleArticle(articleId, req.Author)
	if err != nil {
		return nil, err
	}

	return article, nil
}

func (r *articlesRepository) UpdateArticle(req *articles.ArticleCredential, userID int) (*articles.Article, error) {
	query := `
	UPDATE "articles" SET`
	params := make(map[string]any)
	params["id"] = req.Id

	if req.Title != "" {
		query += " title = :title,"
		query += " slug = :slug,"
		params["title"] = req.Title
		params["slug"] = req.Title
	}

	if req.Body != "" {
		query += " body = :body,"
		params["body"] = req.Body
	}

	if req.Description != "" {
		query += " description = :description,"
		params["description"] = req.Description
	}

	query = query[:len(query)-1]
	query += " WHERE id = :id;"
	fmt.Println("query === ", query)
	if _, err := r.db.NamedExec(query, params); err != nil {
		return nil, fmt.Errorf("update article failed:%v", err)
	}

	return r.GetSingleArticle(req.Id, userID)
}

func (r *articlesRepository) DeleteArticle(articleID, userID int) error {
	query := `
	DELETE
	FROM "articles"
	WHERE "id" = $1 AND "author_id" = $2;`

	result, err := r.db.ExecContext(context.Background(), query, articleID, userID)
	if err != nil {
		return fmt.Errorf("delete article failed: %v", err)
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting number of affected rows failed: %v", err)
	}

	if rowAffected == 0 {
		return fmt.Errorf("the article doesn't exist")
	}

	return nil
}

func (r *articlesRepository) FavoriteArticle(userID, articleID int) (*articles.Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	INSERT INTO "article_favorites"(
		"user_id",
		"article_id"
	)
	SELECT
		$1,$2
	WHERE NOT EXISTS(
		SELECT 1
		FROM "article_favorites"
		WHERE "user_id" = $1 AND "article_id" = $2
	)
	`
	result, err := r.db.ExecContext(ctx, query, userID, articleID)
	if err != nil {
		return nil, fmt.Errorf("add favorite failed: %v", err)
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("getting number of affected rows failed: %v", err)
	}

	if rowAffected == 0 {
		return nil, fmt.Errorf("the favorite already exist")
	}

	return r.GetSingleArticle(articleID, userID)
}

func (r *articlesRepository) UnfavoriteArticle(userID, articleID int) (*articles.Article, error) {
	query := `
	DELETE
	FROM "article_favorites"
	WHERE "user_id" = $1 AND "article_id" = $2;`

	result, err := r.db.ExecContext(context.Background(), query, userID, articleID)
	if err != nil {
		return nil, fmt.Errorf("unfavorite article failed:%v", err)
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("getting number of affected rows failed: %v", err)
	}

	if rowAffected == 0 {
		return nil, fmt.Errorf("already unfavorite article")
	}

	return r.GetSingleArticle(articleID, userID)
}

func (r *articlesRepository) GetTagsList() (*articles.TagList, error) {
	query := `
		SELECT
		json_build_object('tags',array_to_json(array_agg("t"."name")))
		FROM "tags" "t";`

	var bytes string
	tagsResult := new(articles.TagList)
	if err := r.db.Get(&bytes, query); err != nil {
		return nil, fmt.Errorf("get tags list failed: %v", err)
	}

	if err := json.Unmarshal([]byte(bytes), &tagsResult); err != nil {
		return nil, fmt.Errorf("unmarshal tags list failed: %v", err)
	}

	return tagsResult, nil
}
