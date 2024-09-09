package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

const (
	port = ":8080"
)

var DB *sql.DB

func main() {
	// Initialize the database
	var err error
	DB, err = initDB()
	if err != nil {
		panic("Failed to initialize the database: " + err.Error())
	}

	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	server.GET("/verse", getRandomVerse)
	// Uncomment the endpoints (For getting all data & posting data)
	// server.GET("/verses", getAllVerses)
	// server.POST("/verse", postVerses)

	server.Run(port)
}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "api.db")
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(10) //adjustable
	db.SetMaxOpenConns(5)  //adjustable

	err = createTable(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS verses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		reference TEXT NOT NULL,
		text TEXT NOT NULL
	)`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// Verse struct representing a verse in the database
type Verse struct {
	Id        int    `json:"id"`
	Reference string `json:"reference"`
	Text      string `json:"text"`
}

// Function to get a random verse from the database
func getRandomVerseInfo() (*Verse, error) {
	rows, err := DB.Query(`SELECT id, reference, text FROM verses ORDER BY RANDOM() LIMIT 1`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var verse Verse
	if rows.Next() {
		err := rows.Scan(&verse.Id, &verse.Reference, &verse.Text)
		if err != nil {
			return nil, err
		}
	}

	return &verse, nil
}

// Endpoint to get a random verse
func getRandomVerse(c *gin.Context) {
	verse, err := getRandomVerseInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, verse)
}

// Uncomment the code if you want more functionality (Getting all data & posting data)

// func postVerses(c *gin.Context) {
// 	var verses []Verse

// 	if err := c.ShouldBindJSON(&verses); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
// 		return
// 	}

// 	// Use a transaction for batch insertion
// 	tx, err := DB.Begin()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
// 		return
// 	}

// 	stmt, err := tx.Prepare(`INSERT INTO verses (reference, text) VALUES (?, ?)`)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
// 		return
// 	}
// 	defer stmt.Close()

// 	for _, verse := range verses {
// 		_, err := stmt.Exec(verse.Reference, verse.Text)
// 		if err != nil {
// 			tx.Rollback()
// 			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
// 			return
// 		}
// 	}

// 	if err := tx.Commit(); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{"message": "Verses successfully created!"})
// }
// func getAllVerses(c *gin.Context) {
// 	verses, err := allVersesInfo()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, verses)
// }

// func allVersesInfo() ([]Verse, error) {
// 	rows, err := DB.Query(`SELECT id, reference, text FROM verses`)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer rows.Close()

// 	var verses []Verse
// 	for rows.Next() {
// 		var t Verse
// 		if err := rows.Scan(&t.Id, &t.Reference, &t.Text); err != nil {
// 			return nil, err
// 		}
// 		verses = append(verses, t)
// 	}
// 	return verses, nil
// }
