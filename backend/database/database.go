package database

import (
	"database/sql"
	"log"
	"mimic/backend/misc"
	"mimic/backend/types"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func Init() {
	var err error
	dbPath := os.Getenv("DB_PATH")
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Could not connect to database:", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
		password_hash TEXT,
		session_token TEXT
    )`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS uploads (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER,
        file_name TEXT,
        file_path TEXT,
        short_code TEXT
    )`)
	if err != nil {
		log.Fatal(err)
	}

}

func Close() {
	db.Close()
}

func Insert(upload *types.Upload) error {
	stmt, err := db.Prepare(`
		INSERT INTO uploads (user_id, file_name, file_path)
		VALUES (?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(upload.UserID, upload.FileName, upload.FilePath)
	if err != nil {
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// Generate a unique short code from the auto-incremented ID
	upload.ShortCode = misc.EncodeBase62(lastID)

	// Update the record with the generated short code
	_, err = db.Exec(`UPDATE uploads SET short_code = ? WHERE id = ?`, upload.ShortCode, lastID)
	if err != nil {
		return err
	}

	return nil
}

func GetUploads(userID int) []types.Upload {
	rows, err := db.Query("SELECT * FROM uploads WHERE user_id = ? ORDER BY id DESC LIMIT 10 OFFSET 0", userID)
	if err != nil {
		log.Fatal("Error querying database:", err)
	}
	defer rows.Close()

	var uploads []types.Upload
	for rows.Next() {
		var upload types.Upload
		err = rows.Scan(&upload.ID, &upload.UserID, &upload.FileName, &upload.FilePath, &upload.ShortCode)
		if err != nil {
			log.Fatal("Error scanning row:", err)
		}
		uploads = append(uploads, upload)
	}

	return uploads
}

func GetUpload(shortCode string) types.Upload {
	row := db.QueryRow("SELECT * FROM uploads WHERE short_code = ?", shortCode)
	var upload types.Upload
	err := row.Scan(&upload.ID, &upload.UserID, &upload.FileName, &upload.FilePath, &upload.ShortCode)
	if err != nil {
		log.Fatal("Error scanning row:", err)
	}
	return upload
}

func GetUserID(name string) (int, error) {
	row := db.QueryRow("SELECT id FROM users WHERE name = ?", name)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func GetUserByToken(token string) (types.User, error) {
	row := db.QueryRow("SELECT id, name, password_hash FROM users WHERE session_token = ?", token)
	var user types.User
	err := row.Scan(&user.ID, &user.Name, &user.PasswordHash)
	if err != nil {
		return user, err
	}
	return user, nil
}

func UpdateSessionToken(userID int, token string) error {
	stmt, err := db.Prepare(`
		UPDATE users SET session_token = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(token, userID)
	if err != nil {
		return err
	}

	return nil
}

func CreateUser(name string, password string) (types.User, error) {
	var user types.User
	stmt, err := db.Prepare(`
		INSERT INTO users (name, password_hash)
		VALUES (?, ?)
	`)
	if err != nil {
		return user, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(name, password)
	if err != nil {
		return user, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return user, err
	}

	user.ID = int(lastID)

	return user, nil
}
