package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Post struct {
	Title string `json:"title"`
	Body string `json:"body"`
	Author User `json:"author"`
}

type User struct {
	UserName string `json:"userName"`
	Name string `json:"name"`
	Email string `json:"email"`
}

var posts []Post

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", getPosts).Methods("GET")
	router.HandleFunc("/posts/{id}", getAPost).Methods("GET")
	router.HandleFunc("/posts/{id}", updatePost).Methods("PUT")
	router.HandleFunc("/posts/{id}", patchAPost).Methods("PATCH")
	router.HandleFunc("/posts/{id}", deleteAPost).Methods("DELETE")
	router.HandleFunc("/add", addPost).Methods("POST")
	fmt.Println("Listening and serving")
	log.Fatal(http.ListenAndServe(":5000", router))
}

func addPost(w http.ResponseWriter, r *http.Request) {
	setHeader("json", w)
	//item := mux.Vars(r)["newItem"]
	var newItem Post
	json.NewDecoder(r.Body).Decode(&newItem)
	posts = append(posts, newItem)
	json.NewEncoder(w).Encode(posts)
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	setHeader("json", w)
	if len(posts) > 0 {
		json.NewEncoder(w).Encode(posts)
		fmt.Println(len(posts))
		fmt.Println(posts)
		return
	}
	setHeader("html", w)
	w.Write([]byte("<h1>Please add post</h1><br><p>then come back later</p>"))
}

func getAPost(w http.ResponseWriter, r *http.Request) {
	setHeader("json", w)
	idParams := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParams)
	if err != nil {
		//log.Fatal("Provided id for the post is invalid")
		w.Write([]byte("Provided id for the post is invalid"))
	}

	if id >= len(posts) || id < 0 {
		notFound(w)
	} else {
		json.NewEncoder(w).Encode(posts[id])
	}
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	setHeader("json", w)
	idParams := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParams)
	if err != nil {
		notFound(w)
	}

	checkID(id, w)

	var postToUpdate Post
	err = json.NewDecoder(r.Body).Decode(&postToUpdate)
	checkErr("An error occurred to update post", w, err)

	posts[id] = postToUpdate
	setHeader("json", w)
	json.NewEncoder(w).Encode(posts)

	fmt.Println(postToUpdate)
}

func patchAPost(w http.ResponseWriter, r *http.Request) {
	idParams := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParams)

	checkErr("Provided ID is invalid", w, err)
	checkID(id, w)

	post := &posts[id]
	json.NewDecoder(r.Body).Decode(&post)
	//posts[id] = post
	setHeader("json", w)
	json.NewEncoder(w).Encode(post)
}

func deleteAPost(w http.ResponseWriter, r *http.Request) {
	idParams := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParams)
	checkErr("Provided ID is invalid", w, err)
	checkID(id, w)
	fmt.Printf("%d\n", len(posts))
	posts = append(posts[:id], posts[id+1:]...)
	fmt.Printf("%d\n", len(posts))
}

func notFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(404)
	w.Write([]byte("<h1>Page not found<h1>"))
}

func checkErr(msg string,w http.ResponseWriter, err error) {
	if err != nil {
		w.Write([]byte(msg))
	}
}

func checkID(id int, w http.ResponseWriter) {
	if id >= len(posts) || id < 0 {
		w.Write([]byte("Provided id is out of range"))
	}
}

func setHeader(headerType string, w http.ResponseWriter) {
	switch headerType {
	case "json":
		w.Header().Set("Content-Type", "application/json")
	case "html":
		w.Header().Set("Content-Type", "text/html")
	}
}