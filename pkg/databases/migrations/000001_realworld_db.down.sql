BEGIN;

DROP TRIGGER IF EXISTS set_updatedat_timestamp_users_table  ON "users";
DROP TRIGGER IF EXISTS set_updatedat_timestamp_oauth_table  ON "oauth";
DROP TRIGGER IF EXISTS set_updatedat_timestamp_articles_table  ON "articles";
DROP TRIGGER IF EXISTS set_updatedat_timestamp_comments_table  ON "comments";

DROP FUNCTION IF EXISTS set_updatedat_column();

DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "oauth" CASCADE;
DROP TABLE IF EXISTS "user_follows" CASCADE;
DROP TABLE IF EXISTS "articles" CASCADE;
DROP TABLE IF EXISTS "article_favorites" CASCADE;
DROP TABLE IF EXISTS "comments" CASCADE;
DROP TABLE IF EXISTS "tags" CASCADE;
DROP TABLE IF EXISTS "article_tags" CASCADE;

COMMIT;