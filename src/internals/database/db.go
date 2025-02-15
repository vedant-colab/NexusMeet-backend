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
		fmt.Println("error deactivating token", dberr)
		return fmt.Errorf("failed to deactivate existing sessions: %w", dberr)
	}
	return nil
}

func GetUsers() ([]map[string]string, error) {
	if DB == nil {
		return nil, fmt.Errorf("connection not initialized from database")
	}

	query := `SELECT u.name, u.email 
	           FROM users u 
	           INNER JOIN sessions s ON u.name = s.username 
	           WHERE s.isActive = true;`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []map[string]string

	for rows.Next() {
		var name, email string
		if err := rows.Scan(&name, &email); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		user := map[string]string{
			"name":  name,
			"email": email,
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func UpdateProfileInDB(username string, data struct {
	Bio            string `json:"bio"`
	Location       string `json:"location"`
	Website        string `json:"website"`
	ProfilePicture string `json:"profile_picture"`
}) error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	// Prepare the SQL query
	query := `
		UPDATE profile
		SET bio = $1, location = $2, website = $3, profile_picture = $4
		WHERE username = $5;
	`

	// Execute the query
	result, err := DB.Exec(query, data.Bio, data.Location, data.Website, data.ProfilePicture, username)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func GetUserByUsername(username string) (map[string]string, error) {
	if DB == nil {
		return nil, fmt.Errorf("connection not initialized from database")
	}

	query := `SELECT p.username, p.bio, p.location, p.website, p.profile_picture 
              FROM profile p 
              WHERE p.username = $1;`

	row := DB.QueryRow(query, string(username)) // Use QueryRow since we expect one row

	var (
		// bio, location, website, profile_picture string
		user = make(map[string]string) // Initialize the map
	)

	// Handle NULL values from database
	var nullBio, nullLocation, nullWebsite, nullProfilePic sql.NullString

	err := row.Scan(&username, &nullBio, &nullLocation, &nullWebsite, &nullProfilePic)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}

	// Convert NULL values to empty strings if necessary
	user["username"] = username
	user["bio"] = nullBio.String
	user["location"] = nullLocation.String
	user["website"] = nullWebsite.String
	user["profile_picture"] = nullProfilePic.String

	return user, nil
}
