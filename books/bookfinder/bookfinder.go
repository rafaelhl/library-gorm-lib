//go:generate mockery -name=BookFetcher

package bookfinder

import (
	"context"

	"github.com/rafaelhl/library-gorm-lib/books"
)

type (
	BookFetcher interface {
		FindBookByID(ctx context.Context, bookID uint) (books.Book, error)
	}

	Finder struct {
		fetcher BookFetcher
	}
)

func New(fetcher BookFetcher) Finder {
	return Finder{
		fetcher: fetcher,
	}
}

func (f Finder) FindBook(ctx context.Context, bookID uint) (books.Book, error) {
	return f.fetcher.FindBookByID(ctx, bookID)
}
