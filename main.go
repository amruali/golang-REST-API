package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type Article struct {
	ID                 int
	ArticleWriter      string
	PostedAt           time.Time
	UpdatedAt          time.Time
	ArticleWriterEmail string
}

const dbPath string = "your-database-path"

func ConnectDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		panic("failed to connect to database")
	}
	return db
}

func welcomeToWebsite(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my website")
}

func GetArticles(w http.ResponseWriter, r *http.Request) {
	db := ConnectDB()
	defer db.Close()
	var Articles []Article
	db.Find(&Articles)
	json.NewEncoder(w).Encode(Articles)
}

func AddNewArticle(w http.ResponseWriter, r *http.Request) {
	db := ConnectDB()
	defer db.Close()
	parameters := mux.Vars(r)
	name := parameters["name"]
	id, _ := strconv.Atoi(parameters["id"])
	email := parameters["email"]
	db.Create(Article{ArticleWriter: name, ID: id, ArticleWriterEmail: email})
	fmt.Fprintf(w, "New user created")
}
func UpdateArticle(w http.ResponseWriter, r *http.Request) {
	db := ConnectDB()
	defer db.Close()
	parameters := mux.Vars(r)
	id, _ := strconv.Atoi(parameters["id"])
	email := parameters["email"]
	var article Article
	db.Where("id = ?", id).Find(&article)
	article.ArticleWriterEmail = email
	db.Save(&article)
	fmt.Fprintf(w, "user information udated")

}

func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	db := ConnectDB()
	defer db.Close()
	parameters := mux.Vars(r)
	id, _ := strconv.Atoi(parameters["id"])
	var article Article
	db.Where("id = ?", id).Find(&article)
	db.Delete(&article)
	fmt.Fprintf(w, "user deleted")
}

func httpRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", welcomeToWebsite).Methods("GET")
	router.HandleFunc("/api/articles", GetArticles).Methods("GET")
	router.HandleFunc("/api/new_article/{id}/{name}/{email}", AddNewArticle).Methods("POST")
	router.HandleFunc("/api/update_article/{id}/{email}", UpdateArticle).Methods("PUT")
	router.HandleFunc("/api/delete_article/{id}", DeleteArticle).Methods("DELETE")

	http.ListenAndServe(":8080", router)
}

func main() {
	db := ConnectDB()
	db.AutoMigrate(&Article{})
	defer db.Close()
	httpRequests()
}
