package articlespatterns

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/NattpkJsw/real-world-api-go/modules/articles"
	"github.com/NattpkJsw/real-world-api-go/pkg/utils"
	"github.com/jmoiron/sqlx"
)

type IFindArticleBuilder interface {
	openJsonQuery(userId int)
	initQuery()
	countQuery()
	whereQuery()
	sort()
	closeJsonQuery()
	resetQuery()
	Result() ([]*articles.Article, error)
	PrintQUery()
}

type findArticleBuilder struct {
	db             *sqlx.DB
	req            *articles.ArticleFilter
	query          string
	lastStackIndex int
	values         []any
}

type findArticleEngineer struct {
	builder IFindArticleBuilder
}

func FindArticleBuilder(db *sqlx.DB, req *articles.ArticleFilter) IFindArticleBuilder {
	return &findArticleBuilder{
		db:  db,
		req: req,
	}
}

func FindProductEngineer(builder IFindArticleBuilder) *findArticleEngineer {
	return &findArticleEngineer{builder: builder}
}

func (en *findArticleEngineer) FindArticle(userId int) IFindArticleBuilder {
	en.builder.openJsonQuery(userId)
	en.builder.initQuery()
	en.builder.whereQuery()
	en.builder.sort()
	en.builder.closeJsonQuery()
	return en.builder
}

func (en *findArticleEngineer) CountArticle() IFindArticleBuilder {
	en.builder.countQuery()
	en.builder.whereQuery()
	return en.builder
}

func (b *findArticleBuilder) openJsonQuery(userId int) {
	b.values = make([]any, 0)
	b.values = append(b.values, userId)
	b.lastStackIndex = len(b.values)
	b.query += `
	SELECT
		array_to_json(array_agg("t"))
	FROM(`
}

func (b *findArticleBuilder) initQuery() {
	b.query += `
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
				WHERE "af"."user_id" = $1 AND "af"."article_id" = "a"."id"
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
							WHERE "a"."author_id" = "uf"."following_id"  AND "uf"."follower_id" = $1
						) THEN TRUE 
						ELSE FALSE 
					END
				)
			FROM "users" "u"
			WHERE "a"."author_id" = "u"."id"
		) AS "author"
		FROM "articles" "a"
		WHERE 1 = 1`
}
func (b *findArticleBuilder) countQuery() {
	b.query += `
	SELECT
			COUNT(*) AS "count"
	FROM "articles" "a"
	WHERE 1 = 1`
}
func (b *findArticleBuilder) whereQuery() {
	var queryWhere string
	queryWhereStack := make([]string, 0)

	if b.req.Tag != "" {
		b.values = append(b.values, b.req.Tag)
		queryWhereStack = append(queryWhereStack, ` 
		AND "a"."id" IN (
			SELECT "at"."article_id"
			FROM "article_tags" "at"
			JOIN "tags" AS "t" ON "t"."id" = "at"."tag_id"
			WHERE "t"."name" = ?)`)
	}

	if b.req.Author != "" {
		b.values = append(b.values, b.req.Author)
		queryWhereStack = append(queryWhereStack, ` 
		AND "a"."author_id" IN (
			SELECT "u"."id"
			FROM "users" "u"
			WHERE "u"."username" = ?
		)`)
	}

	if b.req.Favorited != "" {
		b.values = append(b.values, b.req.Favorited)
		queryWhereStack = append(queryWhereStack, ` 
		AND "a"."id" IN (
			SELECT "af"."article_id"
			FROM "article_favorites" "af"
			JOIN "users" AS "u" ON "u"."id" = "af"."user_id"
			WHERE "u"."username" = ?)`)
	}

	for i := range queryWhereStack {
		queryWhere += strings.Replace(queryWhereStack[i], "?", "$"+strconv.Itoa(i+2), 1)
	}

	if b.req.IsFeed {
		queryWhere += ` 
		AND EXISTS (
			SELECT 1
			FROM "user_follows" "uf"
			WHERE "a"."author_id" = "uf"."following_id" AND "uf"."follower_id" = $1
		)`
	}

	b.lastStackIndex = len(b.values)
	b.query += queryWhere

}

func (b *findArticleBuilder) sort() { // sort
	b.query += ` 
	ORDER BY "a"."createdat" DESC`
	//  set offset and limit
	b.values = append(b.values, b.req.Offset)
	b.values = append(b.values, b.req.Limit)
	b.query += fmt.Sprintf(` OFFSET $%d LIMIT $%d`, b.lastStackIndex+1, b.lastStackIndex+2)
}

func (b *findArticleBuilder) closeJsonQuery() {
	b.query += `
	) AS "t";`
}

func (b *findArticleBuilder) resetQuery() {
	b.query = ""
	b.values = make([]any, 0)
	b.lastStackIndex = 0
}

func (b *findArticleBuilder) Result() ([]*articles.Article, error) {
	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	bytes := make([]byte, 0)
	articelsResult := make([]*articles.Article, 0)
	if err := b.db.Get(&bytes, b.query, b.values...); err != nil {
		log.Printf("find articles failed: %v\n", err)
		return make([]*articles.Article, 0), err
	}

	if err := json.Unmarshal(bytes, &articelsResult); err != nil {
		log.Printf("unmarshal articles failed: %v\n", err)
		return make([]*articles.Article, 0), err
	}
	b.resetQuery()

	return articelsResult, nil
}

func (b *findArticleBuilder) PrintQUery() {
	utils.Debug(b.values)
	fmt.Println(b.query)
}
