package db

import (
	format "fmt"
	"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db, conn_err = sql.Open("mysql", "root:imonomy@/visadd_stats")

type BlogPost struct  {
	text string
	id int
	date time.Time
	title string
}

//func makeQuery (query string, args ...interface{}) (*driver.Rows, error){
//	go db.Query(query, args)
//}

func GetBlogPost(id int, channel chan BlogPost){
	var post BlogPost
	res, err := db.Prepare("SELECT id, text, date, title FROM blog WHERE id = ?")
	if err != nil {
		log.Panicf("Error: %s", err)
		panic("qeq")
	}
	defer res.Close()
	res.QueryRow(1).Scan(&post.id, &post.text, &post.date, &post.title)
	//format.Printf("Blog post is: %s", post.text)
	//format.Printf("Getting blog post by %d id...", id)
	format.Println("Received request")
	channel <- post
	close(channel)
}