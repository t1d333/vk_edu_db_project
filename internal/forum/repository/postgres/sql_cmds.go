package postgres

var (
	getForumUserCmd = `
        SELECT nickname
        FROM users
        WHERE nickname = $1;
    `

	createCmd = `
        INSERT INTO forums (title, user_nickname, slug)
        VALUES($1, $2, $3)
        RETURNING id, title, user_nickname, slug, posts, threads;
    `

	getForumCmd = `
        SELECT id, title, user_nickname, slug, posts, threads
        FROM forums
        WHERE slug = $1;
    `
	createThreadCmd = `
        INSERT INTO threads (title, author, forum, message, slug, created)
        VALUES($1, $2, $3, $4, $5, $6)
        RETURNING  id, title, author, forum, message, slug, created;
    `

	getThreadsDescCmd = `
        SELECT id, title, author, forum, message, slug, votes, created
        FROM threads
        WHERE forum = $1
        ORDER BY created DESC
        LIMIT $2;
    `
	getThreadsDescWithFilterCmd = `
        SELECT id, title, author, forum, message, slug, votes, created
        FROM threads
        WHERE forum = $1 AND created <= $3
        ORDER BY created DESC
        LIMIT $2;
    `

	getThreadsAscCmd = `
        SELECT id, title, author, forum, message, slug, votes, created
        FROM threads
        WHERE forum = $1
        ORDER BY created
        LIMIT $2;
        `

	getThreadsAscWithFilterCmd = `
        SELECT id, title, author, forum, message, slug, votes, created
        FROM threads
        WHERE forum = $1 AND created >= $3
        ORDER BY created
        LIMIT $2;
    `

	getForumUsersAsc = `
		SELECT nickname, fullname, about, email
        FROM forum_users
        WHERE forum = $1
        ORDER BY nickname
        LIMIT $2;
	   `

	getForumUsersWithSinceAsc = `
		SELECT nickname, fullname, about, email
        FROM forum_users
        WHERE forum = $1 AND nickname > $2
        ORDER BY nickname
        LIMIT $3;
	   `

	getForumUsersDesc = `
		SELECT nickname, fullname, about, email
        FROM forum_users
        WHERE forum = $1
        ORDER BY nickname DESC
        LIMIT $2;
	   `

	getForumUsersWithSinceDesc = `
		SELECT nickname, fullname, about, email
        FROM forum_users
        WHERE forum = $1 AND nickname < $2
        ORDER BY nickname DESC
        LIMIT $3;
	   `
)
