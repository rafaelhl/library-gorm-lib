package bookfinder_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/rafaelhl/library-gorm-lib/books"
	"github.com/rafaelhl/library-gorm-lib/books/bookfinder"
	"github.com/rafaelhl/library-gorm-lib/books/bookfinder/mocks"
)

var (
	bookID   uint = 1
	expected      = books.Book{
		Title:       "Livro de Teste",
		Description: "Esse livro é de teste",
		Author:      "Rafael Holanda",
		Edition:     1,
		BookShelf:   books.Shelf{},
	}
)

func TestFinder_FindBook(t *testing.T) {
	fetcher := new(mocks.BookFetcher)
	fetcher.On("FindBookByID", mock.Anything, bookID).Return(expected, nil)
	finder := bookfinder.New(fetcher)

	book, err := finder.FindBook(context.Background(), bookID)

	assert.NoError(t, err)
	assert.Equal(t, expected, book)
	fetcher.AssertExpectations(t)
}

func TestFinder_FindBookFail(t *testing.T) {
	fetcher := new(mocks.BookFetcher)
	fetcher.On("FindBookByID", mock.Anything, bookID).Return(books.Book{}, errors.New("unexpected error|"))
	finder := bookfinder.New(fetcher)

	_, err := finder.FindBook(context.Background(), bookID)

	assert.Error(t, err)
}
