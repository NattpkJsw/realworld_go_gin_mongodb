package userspatterns

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/NattpkJsw/real-world-api-go/modules/users"
	"github.com/jmoiron/sqlx"
)

type IInsertUser interface {
	Customer() (IInsertUser, error)
	Result() (*users.User, error)
}

type userReq struct {
	id  string
	req *users.UserRegisterReq
	db  *sqlx.DB
}

type customer struct {
	*userReq
}

func InsertUser(db *sqlx.DB, req *users.UserRegisterReq) IInsertUser {
	return &customer{
		userReq: &userReq{
			req: req,
			db:  db,
		},
	}
}

func (f *userReq) Customer() (IInsertUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
	INSERT INTO "users" (
		"email",
		"username",
		"password"
	)
	VALUES
		($1, $2, $3)
	RETURNING "id";`

	if err := f.db.QueryRowContext(
		ctx,
		query,
		strings.ToLower(f.req.Email),
		strings.ToLower(f.req.Username),
		f.req.Password,
	).Scan(&f.id); err != nil {
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"users_username_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("username has been used")
		case "ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)":
			return nil, fmt.Errorf("email has been used")
		default:
			return nil, fmt.Errorf("insert user failed: %v", err)
		}
	}
	return f, nil
}

func (f *userReq) Result() (*users.User, error) {
	query := `
	SELECT
		to_jsonb("t")
	FROM(
		SELECT
			"u"."id",
			"u"."email",
			"u"."username",
			"u"."image",
			"u"."bio"
		FROM "users" "u"
		WHERE "u"."id" = $1
	) AS "t"`

	data := make([]byte, 0)
	if err := f.db.Get(&data, query, f.id); err != nil {
		return nil, fmt.Errorf("get user failed: %v", err)
	}

	user := new(users.User)
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("unmarshal user failed: %v", err)
	}

	return user, nil
}
