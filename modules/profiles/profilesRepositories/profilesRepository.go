package profilesrepositories

import (
	"context"
	"fmt"
	"time"

	"github.com/NattpkJsw/real-world-api-go/modules/profiles"
	"github.com/jmoiron/sqlx"
)

type IProfilesRepository interface {
	FindOneUserProfileByUsername(username string, curUserId int) (*profiles.Profile, error)
	FollowUser(username string, curUserId int) (*profiles.Profile, error)
	UnfollowUser(username string, curUserId int) (*profiles.Profile, error)
}

type profilesRepository struct {
	db *sqlx.DB
}

func ProfilesRepository(db *sqlx.DB) IProfilesRepository {
	return &profilesRepository{
		db: db,
	}
}

func (r *profilesRepository) FindOneUserProfileByUsername(username string, curUserId int) (*profiles.Profile, error) {
	query := `
		SELECT
		"username",
		"bio",
		"image",
		(
			SELECT
				CASE WHEN EXISTS (
					SELECT 1
					FROM "user_follows" "uf"
					WHERE "uf"."follower_id" = $2 AND 
						"uf"."following_id" = "users"."id"
				) THEN TRUE ELSE FALSE END
		) AS "following"
		FROM "users"
		WHERE "username" = $1;`

	profile := new(profiles.Profile)
	if err := r.db.Get(profile, query, username, curUserId); err != nil {
		return nil, fmt.Errorf("user profile not found")
	}
	return profile, nil
}

func (r *profilesRepository) FollowUser(username string, curUserId int) (*profiles.Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
	INSERT INTO "user_follows" 
		("following_id", "follower_id")
	SELECT
		(SELECT "id" FROM "users" WHERE "username" = $1),$2
	WHERE NOT EXISTS (
		SELECT 1
		FROM "user_follows"
		WHERE "following_id" = (SELECT "id" FROM "users" WHERE "username" = $1)
		AND "follower_id" = $2)
	AND (SELECT "id" FROM "users" WHERE "username" = $1) <> $2;`

	if _, err := r.db.ExecContext(ctx, query, username, curUserId); err != nil {
		return nil, fmt.Errorf("fail to follow the user")
	}

	return r.FindOneUserProfileByUsername(username, curUserId)
}

func (r *profilesRepository) UnfollowUser(username string, curUserId int) (*profiles.Profile, error) {
	query := `
	DELETE FROM "user_follows"
	WHERE "follower_id" = $2 AND
	"following_id" = (SELECT 
					"id" 
					FROM "users" 
					WHERE "username" = $1);
	`
	if _, err := r.db.Exec(query, username, curUserId); err != nil {
		return nil, err
	}

	return r.FindOneUserProfileByUsername(username, curUserId)
}
