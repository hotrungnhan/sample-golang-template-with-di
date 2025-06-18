BEGIN;

CREATE TABLE
    IF NOT EXISTS shorten_urls (
        id VARCHAR NOT NULL PRIMARY KEY,
        original_url TEXT NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT now ()
    );

END;