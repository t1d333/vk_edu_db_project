package postgres

var (
	createPostBeginCmd = "INSERT INTO posts (parent, author, message, thread, forum, created) VALUES "

	checkPostAuthor = `
        SELECT id
        FROM users
        WHERE nickname = $1;
    `

	checkPostParent = `
        SELECT id
        FROM posts
        WHERE id = $1 AND thread = $2;
    `

	createThreadCmd = `
        INSERT INTO threads (title, author, forum, message, slug, created)
        VALUES($1, $2, $3, $4, $5, $6)
        RETURNING  id, title, author, (SELECT slug from forums WHERE slug = $3), message, slug, created;
    `

	getThreadBySlugCmd = `
        SELECT  id, title, author, forum, message, slug, votes, created
        FROM threads
        WHERE slug = $1;
    `

	getThreadByIdCmd = `
        SELECT  id, title, author, forum, message, slug, votes, created
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

	getPostsAscCmd = `
        SELECT id, parent, author, message, isEdited, forum, thread, created
        FROM posts
        WHERE thread = $1 AND id > $2
        ORDER BY created, id
        LIMIT $3;
    `

	getPostsDescWithSinceCmd = `
        SELECT id, parent, author, message, isEdited, forum, thread, created
        FROM posts
        WHERE thread = $1 AND id < $2
        ORDER BY created DESC, id DESC 
        LIMIT $3;
    `
	getPostsDescCmd = `
        SELECT id, parent, author, message, isEdited, forum, thread, created
        FROM posts
        WHERE thread = $1
        ORDER BY created DESC, id DESC 
        LIMIT $2;   
    `

	getPostsTreeAscCmd = `
        SELECT id, parent, author, message, isEdited, forum, thread, created
        FROM posts
        WHERE thread = $1
        ORDER BY path, id
        LIMIT $2;
    `

	getPostsTreeWithSinceAscCmd = `
        SELECT id, parent, author, message, isEdited, forum, thread, created
        FROM posts
        WHERE thread = $1 AND path > (SELECT path FROM posts WHERE id = $2) 
        ORDER BY path, id
        LIMIT $3;
    `

	getPostsTreeDescCmd = `
        SELECT id, parent, author, message, isEdited, forum, thread, created
        FROM posts
        WHERE thread = $1
        ORDER BY path DESC, id
        LIMIT $2;
    `

	getPostsTreeWithSinceDescCmd = `
        SELECT id, parent, author, message, isEdited, forum, thread, created
        FROM posts
        WHERE thread = $1 AND path < (SELECT path FROM posts WHERE id = $2) 
        ORDER BY path DESC, id
        LIMIT $3;
    `

	getPostsParentTreeAscCmd = `
        SELECT id, parent, author, message, isEdited, forum, thread, created
        FROM posts
		WHERE path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent = 0 ORDER BY id ASC LIMIT $2)
		ORDER BY path, id;
    `

	getPostsParentTreeWithSinceAscCmd = `
        SELECT id, parent, author, message, isEdited, forum, thread, created
        FROM posts
	    WHERE path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent = 0 AND path[1] >
	    (SELECT path[1] FROM posts WHERE id = $2) ORDER BY id ASC LIMIT $3) 
	    ORDER BY path, id;
    `

	getPostsParentTreeDescCmd = `
        SELECT id, parent, author, message, isEdited, forum, thread,created
        FROM posts WHERE path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent = 0 ORDER BY id DESC LIMIT $2)
		ORDER BY path[1] DESC, path, id ;
    `

	getPostsParentTreeWithSinceDescCmd = `
        SELECT id, parent, author, message, isEdited, forum, thread, created FROM posts
	    WHERE path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent = 0 AND path[1] <
	    (SELECT path[1] FROM posts WHERE id = $2) ORDER BY id DESC LIMIT $3)
	    ORDER BY path[1] DESC, path, id;
    `

	addVoteCmd = `
        INSERT INTO votes
        (nickname, voice, thread)
        VALUES ($1, $2, $3)
        RETURNING id
    `

	updateVoteCmd = `
        UPDATE votes
        SET voice = $1
        WHERE nickname = $2 AND thread = $3 AND voice != $1
        RETURNING id;
    `

	getVoteCmd = `
        SELECT voice
        FROM votes
        WHERE nickname = $1 AND thread = $2;
    `
)
