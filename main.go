package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/rafaelhl/library-gorm-lib/books/bookfinder"
	"github.com/rafaelhl/library-gorm-lib/books/booklistfinder"
	"github.com/rafaelhl/library-gorm-lib/books/booksinserter"
	"github.com/rafaelhl/library-gorm-lib/books/bookupdater"
	"github.com/rafaelhl/library-gorm-lib/books/handler/bookfind"
	"github.com/rafaelhl/library-gorm-lib/books/handler/bookinsert"
	"github.com/rafaelhl/library-gorm-lib/books/handler/booklistfind"
	"github.com/rafaelhl/library-gorm-lib/books/handler/bookupdate"
	"github.com/rafaelhl/library-gorm-lib/books/repository"
)

var dsn = "root:root@(localhost:3306)/library?parseTime=true"

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, "pong")
	})

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	repository := repository.New(db)

	router.Method(http.MethodPost, "/books", bookinsert.NewHandler(booksinserter.New(repository, repository)))
	router.Method(http.MethodGet, "/books/{bookID}", bookfind.NewHandler(bookfinder.New(repository)))
	router.Method(http.MethodPut, "/books/{bookID}", bookupdate.NewHandler(bookupdater.New(repository)))
	router.Method(http.MethodGet, "/books", booklistfind.NewHandler(booklistfinder.New(repository)))

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
