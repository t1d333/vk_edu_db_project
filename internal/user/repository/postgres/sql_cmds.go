package postgres

var (
	createUserCmd = `
		INSERT INTO users (nickname, fullname, about, email)
		VALUES ($1, $2, $3, $4)
		RETURNING id, nickname, fullname, about, email;`

	getUserCmd = `
            SELECT id, nickname, fullname, about, email
            FROM users
            WHERE nickname = $1;
    `

	updateUserCmd = `
        UPDATE users
        SET fullname = $2,
            about = $3,
            email = $4
        WHERE nickname = $1
        RETURNING id, nickname, fullname, about, email;
    `
)
