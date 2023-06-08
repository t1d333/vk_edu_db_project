
CREATE TABLE IF NOT EXISTS user (
    id serial PRIMARY KEY,
    nickname varchar(255) NOT NULL UNIQUE,
    fullname varchar(255),
    about varchar(255),
    email varchar(255) NOT NULL UNIQUE
);
