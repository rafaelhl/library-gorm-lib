//go:generate mockery -name=BookFinder

package bookfind

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/rafaelhl/library-gorm-lib/books"
)

type (
	BookFinder interface {
		FindBook(ctx context.Context, bookID uint) (books.Book, error)
	}

	handler struct {
		finder BookFinder
	}
)

func NewHandler(finder BookFinder) handler {
	return handler{
		finder: finder,
	}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "bookID")
	bookID, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	book, err := h.finder.FindBook(r.Context(), uint(bookID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, book)
}
