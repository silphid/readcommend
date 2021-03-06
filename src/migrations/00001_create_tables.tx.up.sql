CREATE TABLE era
(
  id INTEGER NOT NULL PRIMARY KEY,
  title TEXT NOT NULL,
  min_year SMALLINT,
  max_year SMALLINT
);

CREATE TABLE size
(
  id INTEGER NOT NULL PRIMARY KEY,
  title TEXT NOT NULL,
  min_pages SMALLINT,
  max_pages SMALLINT
);

CREATE TABLE genre
(
  id INTEGER NOT NULL PRIMARY KEY,
  title TEXT NOT NULL
); 

CREATE TABLE author
(
  id INTEGER NOT NULL PRIMARY KEY,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL
);

CREATE TABLE book
(
  id INTEGER NOT NULL PRIMARY KEY,
  title TEXT NOT NULL,
  year_published SMALLINT NOT NULL,
  rating NUMERIC(3, 2) NOT NULL,
  pages SMALLINT NOT NULL,
  genre_id INTEGER REFERENCES genre(id),
  author_id INTEGER REFERENCES author(id)
);

CREATE INDEX book_published ON book USING btree (year_published);
CREATE INDEX book_rating ON book USING btree (rating);
CREATE INDEX book_pages ON book USING btree (pages);
CREATE INDEX book_genre_id ON book USING btree (genre_id);
CREATE INDEX book_author_id ON book USING btree (author_id);
