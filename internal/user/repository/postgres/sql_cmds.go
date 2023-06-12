package postgres

var (
	// createUserCmd = `
	//
	//        WITH new_user AS (
	//            VALUES ('JohnDoe', 'johndoe@example.com')
	//        )
	//        INSERT INTO users (nickname, email)
	//        SELECT nickname, email
	//        FROM new_user
	//        UNION ALL
	//        SELECT nickname, email
	//        FROM users
	//        WHERE nickname = (SELECT nickname FROM new_user)
	//        OR email = (SELECT email FROM new_user)
	//        RETURNING *
	//    `
	createUserCmd = `
		INSERT INTO users (nickname, fullname, about, email)
		VALUES ($1, $2, $3, $4)
		RETURNING id, nickname, fullname, about, email;`

	getUserCmd = `
        SELECT id, nickname, fullname, about, email
        FROM users
        WHERE nickname = $1;
    `
	getUserByEmailCmd = `
        SELECT id, nickname, fullname, about, email
        FROM users
        WHERE email = $1;
    `

	updateUserCmd = `
        UPDATE users
        SET fullname = case when trim($2) = '' then fullname else $2 end,
            about = case when trim($3) = '' then about else $3 end,
            email = case when trim($4) = '' then email else $4 end
        WHERE nickname = $1
        RETURNING id, nickname, fullname, about, email;
    `
)
