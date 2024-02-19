package articlespatterns

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/NattpkJsw/real-world-api-go/modules/articles"
	"github.com/jmoiron/sqlx"
)

type IAddArticleBuilder interface {
	initTransaction() error
	addArticle() error
	addTag(req []*string, articleId int) error
	commit() error
	getArticleId() int
}

type addArticleBuilder struct {
	db  *sqlx.DB
	tx  *sqlx.Tx
	req *articles.ArticleCredential
}

type addArticleEngineer struct {
	builder IAddArticleBuilder
}

func AddArticleBuilder(db *sqlx.DB, req *articles.ArticleCredential) IAddArticleBuilder {
	return &addArticleBuilder{
		db:  db,
		req: req,
	}
}

func AddArticleEngineer(b IAddArticleBuilder) *addArticleEngineer {
	return &addArticleEngineer{
		builder: b,
	}
}

func (en *addArticleEngineer) AddArticle() (int, error) {
	if err := en.builder.initTransaction(); err != nil {
		return 0, err
	}
	if err := en.builder.addArticle(); err != nil {
		return 0, err
	}
	if err := en.builder.commit(); err != nil {
		return 0, err
	}

	return en.builder.getArticleId(), nil
}

func (b *addArticleBuilder) getArticleId() int {
	return b.req.Id
}

func (b *addArticleBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}
func (b *addArticleBuilder) addArticle() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "articles"(
		"title",
		"slug",
		"description",
		"body",
		"author_id"
	)
	VALUES ($1, $1, $2, $3, $4)
	ON CONFLICT ("title", "slug") DO NOTHING
	RETURNING "id";`
	if err := b.tx.QueryRowxContext(
		ctx,
		query,
		b.req.Title,
		b.req.Description,
		b.req.Body,
		b.req.Author,
	).Scan(&b.req.Id); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert article failed: %v", err)
	}
	if len(b.req.TagList) != 0 {
		err := b.addTag(b.req.TagList, b.req.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *addArticleBuilder) addTag(req []*string, articleId int) error {
	query := `INSERT INTO "tags"("name") 
				VALUES ($1) 
				ON CONFLICT ("name") 
				DO NOTHING 
				RETURNING "id";`
	var generatedTagIDs []int
	for _, name := range req {
		var tagId int
		err := b.tx.QueryRowx(query, *name).Scan(&tagId)
		if err != nil && err != sql.ErrNoRows {
			b.tx.Rollback()
			return fmt.Errorf("insert tags fail:%v", err)
		}

		if tagId != 0 {
			generatedTagIDs = append(generatedTagIDs, tagId)
		} else {
			subQuery := `
			SELECT "id" 
			FROM "tags"
			WHERE "name" = $1;`
			var existingTagID int
			if err := b.db.Get(&existingTagID, subQuery, *name); err != nil {
				return fmt.Errorf("get existing tagID fail:%v", err)
			}
			generatedTagIDs = append(generatedTagIDs, existingTagID)
		}

	}

	if err := b.addArticleTag(generatedTagIDs, articleId); err != nil {
		return fmt.Errorf("add article tag error:%v", err)
	}

	return nil
}

func (b *addArticleBuilder) addArticleTag(tagId []int, articeId int) error {
	if len(tagId) == 0 {
		return fmt.Errorf("no tags to add")
	}
	query := `INSERT INTO "article_tags" ("article_id", "tag_id") VALUES ($1, $2)`

	for _, tag := range tagId {
		_, err := b.tx.Exec(query, articeId, tag)
		if err != nil {
			b.tx.Rollback()
			return fmt.Errorf("add article tag failed:%v", err)
		}
	}
	return nil
}

func (b *addArticleBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return fmt.Errorf("commit error: %v", err)
	}
	return nil
}
