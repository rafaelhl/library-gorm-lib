package books

import (
	"time"

	"gorm.io/gorm"
)

type (
	Book struct {
		ID          uint           `gorm:"primarykey"  json:"id"`
		Title       string         `json:"title"`
		Description string         `json:"description"`
		Author      string         `json:"author"`
		Edition     int            `json:"edition"`
		ShelfID     uint           `json:"-"`
		BookShelf   Shelf          `gorm:"foreignkey:ShelfID" json:"shelf"`
		CreatedAt   time.Time      `json:"-"`
		UpdatedAt   time.Time      `json:"-"`
		DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	}

	Shelf struct {
		ID        uint           `gorm:"primarykey" json:"id"`
		Capacity  int            `json:"-"`
		Amount    int            `json:"-"`
		Books     []Book         `json:"-"`
		CreatedAt time.Time      `json:"-"`
		UpdatedAt time.Time      `json:"-"`
		DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	}
)

func (Book) TableName() string {
	return "book"
}

func (Shelf) TableName() string {
	return "shelf"
}
