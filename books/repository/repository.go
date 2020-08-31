package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/rafaelhl/library-gorm-lib/books"
)

type BooksRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) BooksRepository {
	return BooksRepository{
		db: db,
	}
}

func (r BooksRepository) FindShelf(ctx context.Context, shelfID uint) (shelf books.Shelf, err error) {
	err = r.db.WithContext(ctx).Preload(clause.Associations).First(&shelf, shelfID).Error
	return
}

func (r BooksRepository) InsertBook(ctx context.Context, book books.Book) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		book.ShelfID = book.BookShelf.ID
		if err := tx.WithContext(ctx).Create(&book).Error; err != nil {
			return err
		}

		return nil
	})
}
