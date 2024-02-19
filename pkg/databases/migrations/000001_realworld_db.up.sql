BEGIN;

SET TIME ZONE 'Asia/Bangkok';

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION set_updatedat_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updatedat = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';


CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "email" VARCHAR,
  "username" VARCHAR UNIQUE NOT NULL,
  "password" VARCHAR NOT NULL,
  "image" VARCHAR,
  "bio" TEXT,
  "createdat" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updatedat" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "oauth" (
  "id" uuid NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
  "user_id" INT NOT NULL,
  "access_token" VARCHAR NOT NULL,
  "createdat" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updatedat" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "user_follows" (
  "follower_id" INT NOT NULL,
  "following_id" INT NOT NULL
);

CREATE TABLE "articles" (
  "id" SERIAL PRIMARY KEY,
  "author_id" INT NOT NULL,
  "slug" VARCHAR UNIQUE NOT NULL ,
  "title" VARCHAR UNIQUE NOT NULL ,
  "description" VARCHAR,
  "body" TEXT,
  "createdat" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updatedat" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "article_favorites" (
  "user_id" INT NOT NULL,
  "article_id" INT NOT NULL
);

CREATE TABLE "comments" (
  "id" SERIAL PRIMARY KEY,
  "body" TEXT NOT NULL,
  "article_id" INT NOT NULL,
  "author_id" INT NOT NULL,
  "createdat" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updatedat" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "tags" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR NOT NULL UNIQUE
);

CREATE TABLE "article_tags" (
  "article_id" INT NOT NULL,
  "tag_id" INT NOT NULL
);

ALTER TABLE "user_follows" ADD FOREIGN KEY ("follower_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "user_follows" ADD FOREIGN KEY ("following_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "articles" ADD FOREIGN KEY ("author_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "articles" ADD CONSTRAINT unique_title_slug UNIQUE ("title", "slug");
ALTER TABLE "article_favorites" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "article_favorites" ADD FOREIGN KEY ("article_id") REFERENCES "articles" ("id") ON DELETE CASCADE;
ALTER TABLE "comments" ADD FOREIGN KEY ("article_id") REFERENCES "articles" ("id") ON DELETE CASCADE;
ALTER TABLE "comments" ADD FOREIGN KEY ("author_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "article_tags" ADD FOREIGN KEY ("article_id") REFERENCES "articles" ("id") ON DELETE CASCADE;
ALTER TABLE "article_tags" ADD FOREIGN KEY ("tag_id") REFERENCES "tags" ("id") ON DELETE CASCADE;
ALTER TABLE "oauth" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

CREATE TRIGGER set_updatedat_timestamp_users_table BEFORE UPDATE ON "users" FOR EACH ROW EXECUTE PROCEDURE set_updatedat_column();
CREATE TRIGGER set_updatedat_timestamp_oauth_table BEFORE UPDATE ON "oauth" FOR EACH ROW EXECUTE PROCEDURE set_updatedat_column();
CREATE TRIGGER set_updatedat_timestamp_articles_table BEFORE UPDATE ON "articles" FOR EACH ROW EXECUTE PROCEDURE set_updatedat_column();
CREATE TRIGGER set_updatedat_timestamp_comments_table BEFORE UPDATE ON "comments" FOR EACH ROW EXECUTE PROCEDURE set_updatedat_column();

COMMIT;