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

func TestBooksRepository_FindBookByID(t *testing.T) {
	dbFile := fmt.Sprintf(testDB, "find-book")
	db := buildTestDB(dbFile)
	defer os.Remove(dbFile)

	shelf := createShelf(db)

	repo := repository.New(db)
	book := shelf.Books[0]
	b, err := repo.FindBookByID(context.Background(), book.ID)

	assert.NoError(t, err)
	cleanBookTime(&b)
	cleanShelfTime(&b.BookShelf)
	cleanBookTime(&book)
	cleanShelfTime(&shelf)
	shelf.Books = nil
	book.BookShelf = shelf
	assert.Equal(t, book, b)
}

func TestBooksRepository_UpdateBook(t *testing.T) {
	dbFile := fmt.Sprintf(testDB, "update-book")
	db := buildTestDB(dbFile)
	defer os.Remove(dbFile)

	shelf := createShelf(db)

	repo := repository.New(db)
	book := shelf.Books[0]
	book.Edition = 2
	book.Author = "Anonymous"
	err := repo.UpdateBook(context.Background(), book)

	assert.NoError(t, err)

	var storedBook books.Book
	err = db.First(&storedBook, book.ID).Error
	assert.NoError(t, err)

	cleanBookTime(&book)
	cleanBookTime(&storedBook)
	assert.Equal(t, book, storedBook)
}

func TestBooksRepository_FindAllBooks(t *testing.T) {
	dbFile := fmt.Sprintf(testDB, "find-all-books")
	db := buildTestDB(dbFile)
	defer os.Remove(dbFile)

	shelf := createShelf(db)

	repo := repository.New(db)

	allBooks, err := repo.FindAllBooks(context.Background())

	assert.NoError(t, err)
	cleanShelfTime(&shelf)
	cleanBooksTime(allBooks)
	shelfBooks := make([]books.Book, len(shelf.Books))
	for i, b := range shelf.Books {
		b.BookShelf = shelf
		b.BookShelf.Books = nil
		shelfBooks[i] = b
	}
	assert.Equal(t, shelfBooks, allBooks)
}

func TestBooksRepository_DeleteBook(t *testing.T) {
	dbFile := fmt.Sprintf(testDB, "delete-book")
	db := buildTestDB(dbFile)
	defer os.Remove(dbFile)

	shelf := createShelf(db)

	repo := repository.New(db)
	ctx := context.Background()
	book, _ := repo.FindBookByID(ctx, shelf.Books[0].ID)

	err := repo.DeleteBook(ctx, book)

	assert.NoError(t, err)
	var deletedBook books.Book
	db.Unscoped().Find(&deletedBook, book.ID)
	assert.NotNil(t, deletedBook.DeletedAt)
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
