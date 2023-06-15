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
    slug citext NOT NULL UNIQUE,
    user_nickname citext CONSTRAINT user_nickname NOT NULL REFERENCES users (nickname),
    threads int DEFAULT 0,
    posts int DEFAULT 0
  );

CREATE TABLE
  IF NOT EXISTS threads (
    id bigserial PRIMARY KEY,
    title text NOT NULL,
    author citext NOT NULL REFERENCES users (nickname),
    forum citext NOT NULL REFERENCES forums (slug),
    message text NOT NULL,
    votes int DEFAULT 0,
    slug citext,
    created timestamp
    with
      time zone DEFAULT now ()
  );

CREATE TABLE
  IF NOT EXISTS posts (
    id bigserial PRIMARY KEY,
    parent int,
    author citext NOT NULL REFERENCES users (nickname),
    message text NOT NULL,
    isEdited boolean DEFAULT false,
    forum citext REFERENCES forums (slug),
    thread bigint REFERENCES threads (id),
    created timestamp
    with
      time zone DEFAULT now (),
    CONSTRAINT thread_check CHECK(thread IS NOT NULL)
  );


CREATE OR REPLACE FUNCTION validate_parent_thread()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.parent <> 0 AND NOT EXISTS (
    SELECT 1 FROM posts WHERE id = NEW.parent AND thread = NEW.thread
  ) THEN
    RAISE EXCEPTION 'Invalid parent';
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER validate_parent_thread_trigger
BEFORE INSERT OR UPDATE ON posts
FOR EACH ROW
EXECUTE FUNCTION validate_parent_thread();


CREATE OR REPLACE FUNCTION increment_forum_threads()
  RETURNS TRIGGER AS
$$
BEGIN
  UPDATE forums
  SET threads = threads + 1
  WHERE slug = NEW.forum;
  RETURN NEW;
END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION increment_forum_posts()
  RETURNS TRIGGER AS
$$
BEGIN
  UPDATE forums
  SET posts = posts + 1
  WHERE slug = NEW.forum;
  RETURN NEW;
END;
$$
LANGUAGE plpgsql;


CREATE TRIGGER increment_forum_threads_trigger
AFTER INSERT ON threads
FOR EACH ROW
EXECUTE FUNCTION increment_forum_threads();


CREATE TRIGGER increment_forum_posts_trigger
AFTER INSERT ON posts
FOR EACH ROW
EXECUTE FUNCTION increment_forum_posts();



-- TODO: пофиксить этот триггер
-- CREATE FUNCTION update_thread_counter() RETURNS trigger as $thread_counter_updater$
--     UPDATE forums
--     SET threads = threads + 1
--     WHERE slug = NEW.slug
-- $thread_counter_updater$ LANGUAGE plpsql;
--
-- CREATE TRIGGER thread_counter_updater BEFORE INSERT ON thread
--     FOR EACH ROW EXECUTE PROCEDURE update_thread_counter();
--

-- TODO: добавить триггер для обновления счетчика posts
-- TODO: добавить триггер для обновления счетчика votes
