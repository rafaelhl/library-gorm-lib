//go:generate mockery -name=Finder
//go:generate mockery -name=Deleter

package bookdeleter

import (
	"context"

	"github.com/rafaelhl/library-gorm-lib/books"
)

type (
	Finder interface {
		FindBookByID(ctx context.Context, bookID uint) (books.Book, error)
	}

	Deleter interface {
		DeleteBook(ctx context.Context, book books.Book) error
	}

	BookDeleter struct {
		finder  Finder
		deleter Deleter
	}
)

func New(finder Finder, deleter Deleter) BookDeleter {
	return BookDeleter{
		finder:  finder,
		deleter: deleter,
	}
}

func (d BookDeleter) DeleteBook(ctx context.Context, id uint) error {
	book, err := d.finder.FindBookByID(ctx, id)
	if err != nil || book.ID == 0 {
		return err
	}

	return d.deleter.DeleteBook(ctx, book)
}
