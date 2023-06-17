CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE
  IF NOT EXISTS users (
    id bigserial,
    nickname citext NOT NULL UNIQUE COLLATE "ucs_basic" PRIMARY KEY,
    fullname text,
    about text,
    email citext NOT NULL UNIQUE
  );

CREATE TABLE
  IF NOT EXISTS forums (
    id bigserial ,
    title text NOT NULL,
    slug citext NOT NULL UNIQUE PRIMARY KEY,
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
    path BIGINT[] NOT NULL DEFAULT ARRAY[]::BIGINT[],
    created timestamp
    with
      time zone DEFAULT now (),
    CONSTRAINT thread_check CHECK(thread IS NOT NULL)
  );


CREATE TABLE IF NOT EXISTS votes (
    id bigserial,
    nickname citext NOT NULL REFERENCES users (nickname),
    thread bigint NOT NULL REFERENCES threads (id),
    voice int NOT NULL,
    PRIMARY KEY(nickname, thread)
);


CREATE TABLE IF NOT EXISTS forum_users (
    nickname citext NOT NULL COLLATE "ucs_basic" REFERENCES users (nickname),
    fullname text NOT NULL,
    about text NOT NULL,
    email citext NOT NULL,
    forum citext NOT NULL REFERENCES forums (slug),
    PRIMARY KEY (nickname, forum)
);



CREATE OR REPLACE FUNCTION update_post_path()
RETURNS TRIGGER AS $$
BEGIN
    NEW.path = case when NEW.parent = 0 then array_append(NEW.path, NEW.id) else array_append((SELECT path FROM posts WHERE id = NEW.parent), NEW.id) end;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_post_path_trigger 
BEFORE INSERT ON posts
FOR EACH ROW
EXECUTE FUNCTION update_post_path();


CREATE OR REPLACE FUNCTION update_forum_users()
RETURNS TRIGGER AS $$
DECLARE
    nickname_tmp citext;
    fullname_tmp text;
    about_tmp text;
    email_tmp citext;
BEGIN
    SELECT nickname, fullname, about, email
    FROM users 
    WHERE nickname = NEW.author
    INTO nickname_tmp, fullname_tmp, about_tmp, email_tmp;

    INSERT INTO forum_users (nickname, fullname, about, email, forum)
    VALUES (nickname_tmp, fullname_tmp, about_tmp, email_tmp, NEW.forum)
    ON CONFLICT DO NOTHING;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_forum_users_by_post_trigger
AFTER INSERT ON posts
FOR EACH ROW
EXECUTE FUNCTION update_forum_users();

CREATE TRIGGER update_forum_users_by_thread_trigger
AFTER INSERT ON threads
FOR EACH ROW
EXECUTE FUNCTION update_forum_users();


-- counter updaters
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

CREATE OR REPLACE FUNCTION increment_thread_votes()
  RETURNS TRIGGER AS
$$
BEGIN
  UPDATE threads SET votes = threads.votes + new.voice WHERE id = new.thread;
  RETURN NEW;
END;
$$
LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_thread_votes()
  RETURNS TRIGGER AS
$$
BEGIN
  UPDATE threads
  SET votes =  votes + NEW.voice - OLD.voice
  WHERE id = NEW.thread;
  RETURN NEW;
END;
$$
LANGUAGE plpgsql;

-- counter triggers
CREATE TRIGGER increment_forum_threads_trigger
AFTER INSERT ON threads
FOR EACH ROW
EXECUTE FUNCTION increment_forum_threads();

CREATE TRIGGER increment_forum_posts_trigger
AFTER INSERT ON posts
FOR EACH ROW
EXECUTE FUNCTION increment_forum_posts();

CREATE TRIGGER increment_thread_votes_trigger
AFTER INSERT ON votes
FOR EACH ROW
EXECUTE FUNCTION increment_thread_votes();

CREATE TRIGGER update_thread_votes_trigger
AFTER UPDATE ON votes
FOR EACH ROW
EXECUTE FUNCTION update_thread_votes();

-- Indexes 

-- Users
CREATE INDEX IF NOT EXISTS user_nickname_hash ON users using hash (nickname);
CREATE INDEX IF NOT EXISTS  user_nickname_email ON users (nickname, email);

-- Threads
CREATE INDEX IF NOT EXISTS thread_slug_hash ON threads USING hash (slug);
CREATE INDEX IF NOT EXISTS thread_forum_hash ON threads USING hash (forum);
CREATE INDEX IF NOT EXISTS thread_forum_search ON threads (forum, created);

-- Posts
CREATE INDEX IF NOT EXISTS user_posts ON posts (forum, author);
CREATE INDEX IF NOT EXISTS flat_sort ON posts (thread, id);
CREATE INDEX IF NOT EXISTS tree_sort ON posts (thread, path);
CREATE INDEX IF NOT EXISTS parent_tree_sort ON posts ((path[1]), path);

-- Forums
CREATE INDEX IF NOT EXISTS forum_slug_hash ON forums using hash (slug);

-- Votes
CREATE INDEX IF NOT EXISTS user_vote ON votes (nickname, thread);

-- Posts
CREATE INDEX IF NOT EXISTS post_id_hash ON posts using hash (id);
CREATE INDEX IF NOT EXISTS post_thread_hash ON posts using hash (thread);

VACUUM ANALYZE;
