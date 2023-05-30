package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"go-server/books"
)

type bookServer struct {
	store *books.BookStore
}

func NewTaskServer() *bookServer {
	store := books.New()
	return &bookServer{store: store}
}

func (ts *bookServer) bookHandler(w http.ResponseWriter, req *http.Request) {
	pathStrings := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
	fmt.Println(pathStrings, pathStrings[len(pathStrings)-1], pathStrings[len(pathStrings)-1] == "books")
	if pathStrings[len(pathStrings)-1] == "books" && req.Method == http.MethodGet {
		ts.GetAllBooksHandler(w, req)
		return
	}
	if pathStrings[len(pathStrings)-1] == "new" && req.Method == http.MethodPost {
		fmt.Println("NEW")
		ts.CreateBookHandler(w, req)
		return
	}
	if req.Method == http.MethodGet {
		fmt.Println("GET ONE")
		ts.GetBookHandler(w, req, pathStrings[len(pathStrings)-1])
		return
	}
	if req.Method == http.MethodPost {
		fmt.Println("UPDATE")
		ts.UpdateBookHandler(w, req, pathStrings[len(pathStrings)-1])
	}
	if req.Method == http.MethodDelete {
		fmt.Println("DELETE")
		ts.DeleteBookHandler(w, req, pathStrings[len(pathStrings)-1])
		return
	}
	// if req.URL.Path == "/books" {

	// } else {
	// 	http.Error(w, fmt.Sprintf("expect method GET, DELETE or POST at /task/, got %v", req.Method), http.StatusMethodNotAllowed)
	// 	return
	// }
}

func (ts *bookServer) GetAllBooksHandler(w http.ResponseWriter, req *http.Request) {
	allBooks := ts.store.GetAllBooks()
	js, err := json.Marshal(allBooks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ts *bookServer) CreateBookHandler(w http.ResponseWriter, req *http.Request) {
	type ResponseId struct {
		Id string `json:"id"`
	}

	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	var rt books.ResponseBody
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&rt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := ts.store.CreateBook(rt.Name, rt.Author, rt.Total, rt.CurrentPage)
	w.Header().Set("Content-Type", "application/text")
	w.Write([]byte(id.String()))
}

func (ts *bookServer) GetBookHandler(w http.ResponseWriter, req *http.Request, id string) {
	uuidId, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	task, err := ts.store.GetBook(uuidId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	js, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ts *bookServer) UpdateBookHandler(w http.ResponseWriter, req *http.Request, id string) {
	uuidId, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var rt books.ResponseBody
	fmt.Println(rt)
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&rt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ok := ts.store.UpdateBook(uuidId, rt)
	if ok != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func (ts *bookServer) DeleteBookHandler(w http.ResponseWriter, req *http.Request, id string) {
	uuidId, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ok := ts.store.DeleteBook(uuidId)
	if ok != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func main() {
	mux := http.NewServeMux()
	server := NewTaskServer()
	mux.HandleFunc("/books/", server.bookHandler)

	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}
