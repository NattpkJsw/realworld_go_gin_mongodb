BEGIN;

INSERT INTO "users"(
    "username",
    "email",
    "password",
    "image",
    "bio"
)
VALUES
('jake','jake@j.com','$2a$10$7M/tQynKrlcK6uT65AjqEetyS7D5C3CuHyHsnPRIYhEhomMNlzccq','https://i.stack.imgur.com/xHWG8.jpg','I work at statefarm'),
('joedoe','joedoe@j.com','$2a$10$7M/tQynKrlcK6uT65AjqEetyS7D5C3CuHyHsnPRIYhEhomMNlzccq','https://i.stack.imgur.com/xHWG8.jpg','I work at home'),
('doeja','dieja@j.com','$2a$10$7M/tQynKrlcK6uT65AjqEetyS7D5C3CuHyHsnPRIYhEhomMNlzccq','https://i.stack.imgur.com/xHWG8.jpg','fine');

INSERT INTO "articles"(
    "author_id",
    "slug",
    "title",
    "description",
    "body"
)
VALUES
(1,'how-to-train-your-dragon','How to train your dragon','Ever wonder how?','It takes a Jacobian'),
(1,'welcome-to-friday','Welcome to friday','Yawn','Good morning'),
(1,'the first','the first','the first one','the first one'),
(3,'the second','the second','the second one','the second one'),
(3,'the third','the third','the third one','the third one'),
(2,'the one','the one','the one one','the one one'),
(2,'the two','the two','the two one','the two one'),
(2,'the three','the three','the three one','the three one'),
(2,'the four','the four','the four one','the four one'),
(3,'the five','the five','the five one','the five one');

INSERT INTO "comments"(
    "body",
    "article_id",
    "author_id"
)
VALUES
('It takes a Jacobian',1,2),
('Wow WOw wow',2,1);

INSERT INTO "tags"(
    "name"
)
VALUES
('sun'),
('set'),
('greet'),
('morning'),
('good bye');

INSERT INTO "article_tags"(
    "article_id",
    "tag_id"
)
VALUES
(1,1),
(1,2),
(2,3),
(5,4),
(5,5),
(6,2);

INSERT INTO "article_favorites"(
    "user_id",
    "article_id"
)
VALUES
(1,2),
(1,1),
(2,2),
(3,5),
(3,6),
(3,1),
(1,5);

INSERT INTO "user_follows"(
    "follower_id",
    "following_id"
)
VALUES
(1,2),
(3,2),
(1,3),
(2,1);


COMMIT;