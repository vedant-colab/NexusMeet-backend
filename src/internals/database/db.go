package database

import (
	"database/sql"
	"fmt"
	"log"
	"src/internals/config"

	_ "github.com/lib/pq"
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
