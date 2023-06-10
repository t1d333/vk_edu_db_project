CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE
  IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    nickname citext NOT NULL UNIQUE,
    fullname text,
    about text,
    email citext NOT NULL UNIQUE
  );

CREATE TABLE
  IF NOT EXISTS forums (
    id bigserial PRIMARY KEY,
    title text NOT NULL,
    slug text NOT NULL UNIQUE,
    user_nickname citext NOT NULL REFERENCES users (nickname),
    threads int DEFAULT 0,
    posts int DEFAULT 0
  );

CREATE TABLE
  threads (
    id bigserial PRIMARY KEY,
    title text NOT NULL,
    author citext NOT NULL REFERENCES users (nickname),
    forum int REFERENCES forums (id),
    message text NOT NULL,
    votes int DEFAULT 0,
    slug text,
    created timestamp
    with
      time zone DEFAULT now (),
  );
