package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	host     = "localhost"
	dbport   = 5432
	user     = "postgres"
	password = "chen"
	dbname   = "main"
)

type Env struct {
	db *sql.DB
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname = %s sslmode=disable", host, dbport, user, password, dbname)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}

	env := &Env{db: db}

	router := gin.New()
	router.Use(gin.Logger())

	router.Use(cors.Default())

	router.POST("/post/create", env.createPost)
	router.GET("/posts/getAll", env.getPosts)

	log.Fatalln(router.Run(fmt.Sprintf(":%v", port)))
}

type Post struct {
	Content string
}

func (e *Env) createPost(c *gin.Context) {
	newPost := Post{}
	if err := c.BindJSON(&newPost); err != nil {
		log.Printf("invalid JSON body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	q := `INSERT INTO posts(content) VALUES($1)`
	_, err := e.db.Exec(q, newPost.Content)
	if err != nil {
		log.Printf("error occurred while inserting new record: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("success")
	c.JSON(http.StatusOK, "success")
}

func (e *Env) getPosts(c *gin.Context) {
	rows, err := e.db.Query("SELECT id, content FROM posts")
	defer rows.Close()
	if err != nil {
		log.Fatalln(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var id int
	var content string
	var posts []string
	for rows.Next() {
		rows.Scan(&id, &content)
		log.Println(content)
		posts = append(posts, content)
	}
	c.JSON(http.StatusOK, posts)
}
