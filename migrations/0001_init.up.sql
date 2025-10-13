CREATE EXTENSION IF NOT EXISTS citext;

DO $$ BEGIN
    IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'media_type') THEN
        CREATE TYPE media_type AS ENUM ('FILM', 'BOOK', 'GAME', 'SHOW', 'MUSIC');
    end if;
    IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'film_type') THEN
        CREATE TYPE film_type AS ENUM ('FEATURE FILM', 'SHORT FILM');
    end if;
    IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'book_type') THEN
        CREATE TYPE book_type AS ENUM ('NOVEL', 'SHORT STORY', 'COLLECTION');
    end if;
    IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'show_type') THEN
        CREATE TYPE show_type AS ENUM ('SEASON', 'MINI SERIES');
    end if;
    IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'music_type') THEN
        CREATE TYPE music_type AS ENUM ('ALBUM', 'EXTENDED PLAY');
    end if;
END $$;

CREATE TABLE IF NOT EXISTS users(
    id BIGSERIAL PRIMARY KEY,
    email CITEXT UNIQUE NOT NULL,
    display_name TEXT UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    password_hash TEXT NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    last_login TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS media(
    id BIGSERIAL PRIMARY KEY,
    media_type media_type NOT NULL,
    title TEXT NOT NULL,
    release_date DATE NOT NULL,
    description TEXT,
    external_references JSONB DEFAULT '{}'::JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by BIGINT REFERENCES users(id) ON DELETE RESTRICT,
    updated_at TIMESTAMPTZ,
    updated_by BIGINT REFERENCES users(id) ON DELETE RESTRICT,
);

CREATE TABLE IF NOT EXISTS films(
    media_id BIGINT PRIMARY KEY REFERENCES media(id) ON DELETE CASCADE,
    film_type film_type NOT NULL,
    runtime_minutes INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS books(
    media_id BIGINT PRIMARY KEY REFERENCES media(id) ON DELETE CASCADE,
    book_type book_type NOT NULL,
    pages INTEGER
);

CREATE TABLE IF NOT EXISTS games(
    media_id BIGINT PRIMARY KEY REFERENCES media(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS shows(
    media_id BIGINT PRIMARY KEY REFERENCES media(id) ON DELETE CASCADE,
    show_type show_type NOT NULL,
    season INTEGER NOT NULL,
    end_date DATE,
    episodes INTEGER,
    runtime_minutes INTEGER
);

CREATE TABLE IF NOT EXISTS musics(
    media_id BIGINT PRIMARY KEY REFERENCES media(id) ON DELETE CASCADE,
    music_type music_type NOT NULL,
    runtime_minutes INTEGER,
    tracks INTEGER
);

CREATE TABLE IF NOT EXISTS user_media(
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    media_id BIGINT REFERENCES media(id) ON DELETE CASCADE,
    rating INTEGER,
    PRIMARY KEY (user_id, media_id)
);

CREATE TABLE IF NOT EXISTS diary_entry(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    media_id BIGINT REFERENCES media(id) ON DELETE CASCADE,
    started_at TIMESTAMPTZ NOT NULL,
    ended_at TIMESTAMPTZ,
    rating INTEGER,
    description TEXT,
    is_private BOOLEAN NOT NULL DEFAULT FALSE
)