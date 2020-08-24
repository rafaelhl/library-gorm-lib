// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	books "github.com/rafaelhl/library-gorm-lib/books"

	context "context"

	mock "github.com/stretchr/testify/mock"
)

// UpdaterEngine is an autogenerated mock type for the UpdaterEngine type
type UpdaterEngine struct {
	mock.Mock
}

// UpdateBook provides a mock function with given fields: ctx, book
func (_m *UpdaterEngine) UpdateBook(ctx context.Context, book books.Book) error {
	ret := _m.Called(ctx, book)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, books.Book) error); ok {
		r0 = rf(ctx, book)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
