package postgres

var (
	getPostById = `
        SELECT  id, parent, author, message, isEdited, forum, thread, created
        FROM posts
        WHERE id = $1;
    `

	getPostAuthor = `
        SELECT id, nickname, fullname, about, email
        FROM users
        WHERE nickname = $1;
    `

	getPostThread = `
        SELECT  id, title, author, forum, message, slug, votes, created
        FROM threads
        WHERE id = $1;
    `
	getPostForum = `
        SELECT id, title, user_nickname, slug, posts, threads
        FROM forums
        WHERE slug = $1;
    `

	updatePost = `
        UPDATE posts
        SET isEdited = case when (trim($2) = '') OR (trim($2) = trim(message)) then false else true end,
            message = case when trim($2) = '' then message else $2 end
        WHERE id = $1
        RETURNING id, parent, author, message, isEdited, forum, thread, created; 
    `
)
