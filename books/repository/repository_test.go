package repository_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/rafaelhl/library-gorm-lib/books"
	"github.com/rafaelhl/library-gorm-lib/books/repository"
)

const testDB = "test-%s.db"

var (
	shelf = books.Shelf{
		Capacity: 2,
		Amount:   1,
		Books: []books.Book{
			{
				Title:       "Test",
				Description: "Test of database",
				Author:      "Test",
				Edition:     1,
			},
		},
	}
)

func TestBooksRepository_FindShelf(t *testing.T) {
	dbFile := fmt.Sprintf(testDB, "find-shelf")
	db := buildTestDB(dbFile)
	defer os.Remove(dbFile)

	shelf := createShelf(db)

	repo := repository.New(db)
	s, err := repo.FindShelf(context.Background(), 1)

	cleanShelfTime(&s)
	cleanShelfTime(&shelf)

	assert.NoError(t, err)
	assert.Equal(t, shelf, s)
}

func TestBooksRepository_InsertBook(t *testing.T) {
	dbFile := fmt.Sprintf(testDB, "insert-book")
	db := buildTestDB(dbFile)
	defer os.Remove(dbFile)

	shelf := createShelf(db)

	repo := repository.New(db)
	book := books.Book{
		Title:       "Insert Test",
		Description: "Test of insertion of a book",
		Edition:     1,
		Author:      "Tester of Insert",
		BookShelf:   shelf,
	}
	err := repo.InsertBook(context.Background(), book)

	assert.NoError(t, err)

	var (
		s books.Shelf
		b books.Book
	)

	db.First(&s)
	assert.Equal(t, shelf.ID, s.ID)
	assert.Equal(t, 2, s.Amount)

	var _books []books.Book
	db.Find(&_books)
	assert.Len(t, _books, 2)

	b = _books[1]
	cleanBookTime(&b)
	book.ID = 2
	book.ShelfID = 1
	book.BookShelf = books.Shelf{}
	assert.Equal(t, book, b)
}

func createShelf(db *gorm.DB) books.Shelf {
	shelfTest := shelf
	db.Create(&shelfTest)
	return shelfTest
}

func cleanShelfTime(s *books.Shelf) {
	s.CreatedAt = time.Time{}
	s.UpdatedAt = time.Time{}
	cleanBooksTime(s.Books)
}

func cleanBooksTime(books []books.Book) {
	for i, b := range books {
		cleanBookTime(&b)
		cleanShelfTime(&b.BookShelf)
		books[i] = b
	}
}

func cleanBookTime(b *books.Book) {
	b.CreatedAt = time.Time{}
	b.UpdatedAt = time.Time{}
}

func buildTestDB(dbFile string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect db")
	}
	_ = db.AutoMigrate(&books.Book{}, &books.Shelf{})
	return db
}
