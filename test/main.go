package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type postInfo struct {
	idx     int
	pseudo  string
	message string
}
type postMessage struct {
	id      int
	image   string
	pseudo  string
	titre   string
	message string
}

type oui struct {
	Poste   []postMessage
	Comment []postMessage
}

func databasePost(username string, newPost string) {

	database, _ :=
		sql.Open("sqlite3", "data.db")
	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS post (id INTEGER PRIMARY KEY, username TEXT, newPost TEXT)")
	statement.Exec()
	statement, _ =
		database.Prepare("INSERT INTO post (username, newPost) VALUES (?, ?)")
	fmt.Println("ici")
	statement.Exec(username, newPost)
	rows, _ :=
		database.Query("SELECT id, username, newPost FROM post")
	var id int
	var test []string
	for rows.Next() {
		rows.Scan(&id, &username, &newPost)
		test = append(test, strconv.Itoa(id)+": "+username+" "+newPost+"\n")
	}

}

func databaseComment(username string, newComment string) {

	database, _ :=
		sql.Open("sqlite3", "data.db")
	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS comment (id INTEGER PRIMARY KEY, username TEXT, newComment TEXT)")
	statement.Exec()
	statement, _ =
		database.Prepare("INSERT INTO comment (username, newComment) VALUES (?, ?)")
	fmt.Println("ici")
	statement.Exec(username, newComment)
	rows, _ :=
		database.Query("SELECT id, username, newComment FROM comment")
	var id int
	var test []string
	for rows.Next() {
		rows.Scan(&id, &username, &newComment)
		test = append(test, strconv.Itoa(id)+": "+username+" "+newComment+"\n")
	}

}

func PostHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connect")

	userName := r.FormValue("username")
	newPost := r.FormValue("newPost")

	if userName != "" && newPost != "" {
		databasePost(userName, newPost)
		fmt.Println("tu sors de  wsh")
	}

	tpl := template.Must(template.ParseFiles("assets/signIn.html"))

	data := oui{
		Poste: getPostInfo(),
	}

	err := tpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}

}

func CommentHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connect")

	userName := r.FormValue("username")
	newComment := r.FormValue("newComment")

	if userName != "" && newComment != "" {
		databaseComment(userName, newComment)
		fmt.Println("tu sors de  wsh")
	}

	tpl := template.Must(template.ParseFiles("assets/comment.html"))

	data := oui{
		Poste:   getPostInfo(),
		Comment: getCommentInfo(),
	}

	err := tpl.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}

}

func getPostInfo() []postMessage {

	database, _ :=
		sql.Open("sqlite3", "data.db")
	rows, _ :=
		database.Query("SELECT id, username, newPost FROM post")

	var _id int
	var test []postMessage
	var _pseudo string
	var _message string
	for rows.Next() {
		rows.Scan(&_id, &_pseudo, &_message)
		data := postMessage{
			id:      _id,
			pseudo:  _pseudo,
			message: _message,
		}
		test = append(test, data)
	}

	return test

}

func getCommentInfo() []postMessage {

	database, _ :=
		sql.Open("sqlite3", "data.db")
	rows, _ :=
		database.Query("SELECT id, username, newComment FROM comment")

	var _id int
	var test []postMessage
	var _pseudo string
	var _message string
	for rows.Next() {
		rows.Scan(&_id, &_pseudo, &_message)
		data := postMessage{
			id:      _id,
			pseudo:  _pseudo,
			message: _message,
		}
		test = append(test, data)
	}
	return test
}

func main() {
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.Handle("/", fs)
	http.HandleFunc("/account", PostHandle)
	http.HandleFunc("/comment", CommentHandle)
	http.ListenAndServe(":8080", nil)
}
