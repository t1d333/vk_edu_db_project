package postgres

var getStatusCmd = `
       SELECT
       (SELECT count(*) FROM users)  AS users,
       (SELECT count(*) FROM forums) AS forums,
       (SELECT count(*) FROM threads) AS threads,
       (SELECT count(*) FROM posts)  AS posts
       `

var clearDb = `
    TRUNCATE TABLE
    users,
    forums,
    threads,
    posts,
    votes
    CASCADE;
`
