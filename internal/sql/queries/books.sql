-- TODO: change this to get book by isbn
-- name: GetBookByID :one
SELECT * FROM books WHERE id = $1;

-- name: CreateBook :one
-- TODO: year_of_publishing better naming,
-- number_of_pages?
-- year_of_publishing to publish_date
-- read_status and read_date should be nullable
-- personal_rating should be nullable
-- personal_notes should be nullable
-- img should be nullable
-- number_of_pages should be nullable
-- publisher should be nullable
-- category should be nullable
INSERT INTO books
(user_id,isbn, title, author, category, publisher, year_of_publishing,
img, number_of_pages, personal_rating, personal_notes, read_status, read_date)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING id;


-- name: GetBooks :many
SELECT * FROM books ORDER BY id DESC LIMIT 10 OFFSET $1;
