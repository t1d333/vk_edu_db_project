package postgres

var (
	createPostCmd = `
        INSERT INTO posts
        (parent, author, message, thread, forum, created)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, parent, author, message, isEdited, forum, thread; 
    `
	createPostByIdCmd = `
        INSERT INTO posts
        (parent, author, message, thread, forum, created)
        VALUES ($1, $2, $3, $4, (SELECT forum FROM threads WHERE id = $4))
        RETURNING id, parent, author, message, isEdited, forum, thread; 
    `

	createPostBySlugCmd = `
        WITH thread AS (
            SELECT id, forum
            FROM threads
            WHERE slug = $4
        )
        INSERT INTO posts
        (parent, author, message, thread, forum, created)
        VALUES ($1, $2, $3, (SELECT id FROM thread) , (SELECT forum FROM thread))
        RETURNING id, parent, author, message, isEdited, forum, thread; 
    `

	createThreadCmd = `
        INSERT INTO threads (title, author, forum, message, slug, created)
        VALUES($1, $2, $3, $4, $5, $6)
        RETURNING  id, title, author, (SELECT slug from forums WHERE slug = $3), message, slug, created;
    `

	getThreadBySlugCmd = `
        SELECT  id, title, author, forum, message, slug, created
        FROM threads
        WHERE slug = $1;
    `

	getThreadByIdCmd = `
        SELECT  id, title, author, forum, message, slug, created
        FROM threads
        WHERE id = $1;
    `

	updateThreadByIdCmd = `
        UPDATE threads
        SET message = case when trim($2) = '' then message else $2 end, 
            title = case when trim($3) = '' then title else $3 end
        WHERE id = $1
        RETURNING  id, title, author, forum, message, slug, votes, created;
    `

	updateThreadBySlugCmd = `
        UPDATE threads
        SET message = case when trim($2) = '' then message else $2 end, 
            title = case when trim($3) = '' then title else $3 end
        WHERE slug = $1
        RETURNING  id, title, author, forum, message, slug, votes, created;
    `
)
