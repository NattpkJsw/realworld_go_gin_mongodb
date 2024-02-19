package usersrepositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/NattpkJsw/real-world-api-go/modules/users"
	userspatterns "github.com/NattpkJsw/real-world-api-go/modules/users/usersPatterns"
	"github.com/jmoiron/sqlx"
)

type IUsersRepository interface {
	InsertUser(req *users.UserRegisterReq) (*users.User, error)
	FindOneUserByEmail(email string) (*users.UserCredentialCheck, error)
	InsertOauth(req *users.UserToken) error
	FindOneOath(accessToken string) (*users.Oauth, error)
	UpdateOauth(req *users.UserToken) error
	GetProfile(userId int) (*users.User, error)
	DeleteOauth(accessToken string) error
	UpdateUser(user *users.UserCredentialCheck) (*users.User, error)
}

type usersRepository struct {
	db *sqlx.DB
}

func UsersRepository(db *sqlx.DB) IUsersRepository {
	return &usersRepository{
		db: db,
	}
}

func (r *usersRepository) InsertUser(req *users.UserRegisterReq) (*users.User, error) {
	result := userspatterns.InsertUser(r.db, req)

	var err error

	result, err = result.Customer()
	if err != nil {
		return nil, err
	}

	// Get result from inserting
	user, err := result.Result()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *usersRepository) FindOneUserByEmail(email string) (*users.UserCredentialCheck, error) {
	query := `
	SELECT
		"id",
		"email",
		"password",
		"username",
		"image",
		"bio"
	FROM "users"
	WHERE "email" = $1;`

	user := new(users.UserCredentialCheck)
	if err := r.db.Get(user, query, strings.ToLower(email)); err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (r *usersRepository) InsertOauth(req *users.UserToken) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
	INSERT INTO "oauth" (
		"user_id",
		"access_token"
	)
	VALUES ($1, $2)
		RETURNING "id";`

	if err := r.db.QueryRowContext(
		ctx,
		query,
		req.User_Id,
		req.AccessToken,
	).Scan(&req.Id); err != nil {
		return fmt.Errorf("insert oauth failed: %v", err)
	}
	return nil
}

func (r *usersRepository) DeleteOauth(accessToken string) error {
	query := `
	DELETE FROM "oauth" WHERE "access_token" = $1;`
	if _, err := r.db.ExecContext(context.Background(), query, accessToken); err != nil {
		return fmt.Errorf("oauth not found ")
	}
	return nil
}

func (r *usersRepository) GetProfile(userId int) (*users.User, error) {
	query := `
	SELECT
		"id",
		"email",
		"username",
		"image",
		"bio"
	FROM "users"
	WHERE "id" = $1;`

	profile := new(users.User)
	if err := r.db.Get(profile, query, userId); err != nil {
		return nil, fmt.Errorf("get user failed: %v", err)
	}
	return profile, nil
}

func (r *usersRepository) UpdateUser(user *users.UserCredentialCheck) (*users.User, error) {
	// Build the update query dynamically
	query, params, err := buildUpdateQuery(user)
	if err != nil {
		return nil, err
	}

	// Execute the query
	if _, err := r.db.NamedExec(query, params); err != nil {
		return nil, fmt.Errorf("update user failed: %v", err)
	}

	return r.GetProfile(user.Id)
}

func buildUpdateQuery(user *users.UserCredentialCheck) (string, map[string]interface{}, error) {
	// ... (same as your original code)
	query := `UPDATE "users" SET`
	params := make(map[string]any)
	params["id"] = user.Id

	if user.Username != "" {
		query += ` 
		"username" = :username,`
		params["username"] = user.Username
	}

	if user.Email != "" {
		query += ` 
		"email" = :email,`
		params["email"] = user.Email
	}

	if user.Bio != nil {
		query += ` 
		"bio" = :bio,`
		params["bio"] = user.Bio
	}

	if user.Image != nil {
		query += ` 
		"image" = :image,`
		params["image"] = user.Image
	}

	if user.Password != "" {
		if err := user.BcryptHashingUpdate(); err != nil {
			return "", nil, fmt.Errorf("hash password failed: %v", err)
		}
		query += ` "password" = :password,`
		params["password"] = user.Password
	}
	// Remove the trailing comma
	query = query[:len(query)-1]
	query += ` 
	WHERE "id" = :id;`

	return query, params, nil
}

func (r *usersRepository) UpdateOauth(req *users.UserToken) error {
	query := `
	UPDATE "oauth" SET
		"access_token" = :access_token
	WHERE "id" = :id;`

	if _, err := r.db.NamedExecContext(context.Background(), query, req); err != nil {
		return fmt.Errorf("update oauth failed: %v", err)
	}
	return nil
}

func (r *usersRepository) FindOneOath(accessToken string) (*users.Oauth, error) {
	query := `
	SELECT
		"id",
		"user_id"
	FROM "oauth"
	WHERE "access_token" = $1;`

	oauth := new(users.Oauth)
	if err := r.db.Get(oauth, query, accessToken); err != nil {
		return nil, fmt.Errorf("oauth not found, %v", err)
	}
	return oauth, nil
}
