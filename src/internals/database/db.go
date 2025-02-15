package database

import (
	"database/sql"
	"fmt"
	"log"
	"src/internals/config"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func ConnectDB() {
	var err error

	dsn := config.GetEnv("DATABASE_DSN", "")
	// dsn := os.Getenv("DATABASE_DSN")

	// Open a connection
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Verify the connection
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	fmt.Println("Successfully connected to the database")
}

func CreateUser(name string, email string, password string) (sql.Result, error) {
	if DB == nil {
		return nil, fmt.Errorf("connection not initialized from database")
	}

	rows := DB.QueryRow("SELECT * FROM users WHERE email=$1", email)
	rowErr := rows.Scan()
	if rowErr != sql.ErrNoRows {
		return nil, fmt.Errorf("user already exists")
	}

	query := "INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)"
	result, err := DB.Exec(query, name, email, string(password))
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	fmt.Println("result", result)
	return result, err
}

func SaveUserToken(token string, name string, password string) (sql.Result, error) {
	if DB == nil {
		return nil, fmt.Errorf("connection not initialized from database")
	}

	var hashedPassword string
	err := DB.QueryRow("SELECT password_hash FROM users WHERE name = $1", name).Scan(&hashedPassword)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user doesn't exist")
	} else if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	// Compare the provided password with the hashed password in the database
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	// Deactivate any existing active sessions for the user
	_, dberr := DB.Exec("UPDATE sessions SET isActive = false WHERE username = $1 AND isActive = true", name)
	if dberr != nil {
		return nil, fmt.Errorf("failed to deactivate existing sessions: %w", dberr)
	}

	// Insert the new session token
	result, err := DB.Exec("INSERT INTO sessions (token, username, isActive) VALUES ($1, $2, true)", token, name)
	if err != nil {
		return nil, fmt.Errorf("error saving user token: %w", err)
	}

	return result, nil
}

func DeactivateToken(token string, username string) error {
	if DB == nil {
		return fmt.Errorf("connection not initialized from database")
	}
	_, dberr := DB.Exec("UPDATE sessions SET isActive = false WHERE username = $1 AND isActive = true", username)
	if dberr != nil {
		return fmt.Errorf("failed to deactivate existing sessions: %w", dberr)
	}
	return nil
}
