-- TODO: decide what can be null
CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    isbn VARCHAR(13) NOT NULL DEFAULT '',
    author TEXT NOT NULL DEFAULT '',
    category TEXT NOT NULL DEFAULT '',
    publisher TEXT NOT NULL DEFAULT '',
    -- TODO: publishation_date, publish_date
    year_of_publishing INT NOT NULL DEFAULT 0,
    -- image should be a blob
    -- or a path to the image
    img TEXT NOT NULL DEFAULT '',
    -- only pages?
    number_of_pages INT NOT NULL DEFAULT 0,
    -- personal rating should be from 1 to 5 a floating number
    personal_rating FLOAT NOT NULL DEFAULT 0,
    personal_notes TEXT NOT NULL DEFAULT '',
    -- read status should be a boolean
    read_status BOOLEAN NOT NULL DEFAULT false,
    read_date DATE NOT NULL DEFAULT CURRENT_DATE,
    -- percentage of completion from 0 to 100
    -- reading_progress FLOAT NOT NULL,
    user_id SERIAL REFERENCES users (id) ON DELETE CASCADE NOT NULL
);
