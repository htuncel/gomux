package models

// Book struct defines schema
type Book struct {
	ID     uint   `json:"id" gorm:"primary_key"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

// CreateBookInput struct defines book create information schema
type CreateBookInput struct {
	Title  string `json:"title" validate:"required"`
	Author string `json:"author" validate:"required"`
}

// UpdateBookInput struct defines book update information schema
type UpdateBookInput struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}
