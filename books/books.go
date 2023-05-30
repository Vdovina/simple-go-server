package books

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type BookItem struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Author      string    `json:"author"`
	Total       int16     `json:"total"`
	CurrentPage int16     `json:"current_page"`
	UpdateDate  time.Time `json:"update_date"`
}

type ResponseBody struct {
	Name        string `json:"name"`
	Author      string `json:"author"`
	Total       int16  `json:"total"`
	CurrentPage int16  `json:"current_page"`
}

type BookStore struct {
	sync.Mutex
	books map[int]BookItem
}

func (ts *BookStore) foundById(id uuid.UUID) int {
	index := -1
	for i, item := range ts.books {
		if item.Id == id {
			index = i
			break
		}
	}
	return index
}

func New() *BookStore {
	ts := &BookStore{}
	ts.books = make(map[int]BookItem)
	return ts
}

func (ts *BookStore) CreateBook(name string, author string, total int16, currentPage int16) uuid.UUID {
	ts.Lock()
	defer ts.Unlock()

	id := uuid.New()
	book := BookItem{id, name, author, total, currentPage, time.Now()}
	ts.books[len(ts.books)] = book
	return book.Id
}

func (ts *BookStore) GetBook(id uuid.UUID) (BookItem, error) {
	ts.Lock()
	defer ts.Unlock()

	var book BookItem
	for _, item := range ts.books {
		if item.Id == id {
			book = item
			break
		}
	}

	if (book == BookItem{}) {
		return BookItem{}, fmt.Errorf("Book wasn't found")
	}
	return book, nil
}

func (ts *BookStore) GetBookByName(name string) (BookItem, error) {
	ts.Lock()
	defer ts.Unlock()

	var book BookItem
	for _, item := range ts.books {
		if item.Name == name {
			book = item
			break
		}
	}

	if (book == BookItem{}) {
		return BookItem{}, fmt.Errorf("Book wasn't found")
	}
	return book, nil
}

func (ts *BookStore) UpdateBook(id uuid.UUID, newData ResponseBody) error {
	ts.Lock()
	defer ts.Unlock()

	index := ts.foundById(id)
	if index == -1 {
		return fmt.Errorf("Book wasn't found")
	}

	// updated := ts.books[index]
	// updated.CurrentPage = (map[bool]int{true: a, false: a - 1})[]
	// ts.books[index] = updated
	return nil
}

func (ts *BookStore) DeleteBook(id uuid.UUID) error {
	ts.Lock()
	defer ts.Unlock()

	index := ts.foundById(id)
	if index == -1 {
		return fmt.Errorf("Book wasn't found")
	}

	delete(ts.books, index)
	return nil
}

func (ts *BookStore) GetAllBooks() []BookItem {
	// TODO load list from file
	ts.Lock()
	defer ts.Unlock()

	allBooks := make([]BookItem, 0, len(ts.books))
	for _, task := range ts.books {
		allBooks = append(allBooks, task)
	}
	return allBooks
}
