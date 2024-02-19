CREATE TABLE "users" (
  "id" int PRIMARY KEY,
  "email" varchar,
  "username" varchar UNIQUE,
  "password" varchar,
  "image" varchar,
  "bio" text,
  "createdat" timestamp,
  "updatedat" timestamp
);

CREATE TABLE "user_follows" (
  "follower_id" int,
  "following_id" int
);

CREATE TABLE "articles" (
  "id" int PRIMARY KEY,
  "author_id" int,
  "title" varchar,
  "description" varchar,
  "body" text,
  "createdat" timestamp,
  "updatedat" timestamp
);

CREATE TABLE "article_favorites" (
  "user_id" int,
  "article_id" int
);

CREATE TABLE "comments" (
  "id" int PRIMARY KEY,
  "body" text,
  "article_id" int,
  "author_id" int,
  "createdat" timestamp,
  "updatedat" timestamp
);

CREATE TABLE "tags" (
  "id" int PRIMARY KEY,
  "name" varchar
);

CREATE TABLE "article_tags" (
  "article_id" int,
  "tag_id" int
);

CREATE TABLE "oauth" (
  "id" varchar PRIMARY KEY,
  "user_id" varchar,
  "access_token" varchar,
  "createdat" timestamp,
  "updatedat" timestamp
);

ALTER TABLE "user_follows" ADD FOREIGN KEY ("follower_id") REFERENCES "users" ("id");

ALTER TABLE "user_follows" ADD FOREIGN KEY ("following_id") REFERENCES "users" ("id");

ALTER TABLE "articles" ADD FOREIGN KEY ("author_id") REFERENCES "users" ("id");

ALTER TABLE "article_favorites" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "article_favorites" ADD FOREIGN KEY ("article_id") REFERENCES "articles" ("id");

ALTER TABLE "comments" ADD FOREIGN KEY ("article_id") REFERENCES "articles" ("id");

ALTER TABLE "comments" ADD FOREIGN KEY ("author_id") REFERENCES "users" ("id");

ALTER TABLE "article_tags" ADD FOREIGN KEY ("article_id") REFERENCES "articles" ("id");

ALTER TABLE "article_tags" ADD FOREIGN KEY ("tag_id") REFERENCES "tags" ("id");

ALTER TABLE "oauth" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
