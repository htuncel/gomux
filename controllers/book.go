package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"main/configs"
	"main/models"
	"main/utils"
)

// BookController entails all the methods concerning books a wise man
type BookController struct{}

// BookService interface of the book controller
type BookService interface {
	FindBooks(w http.ResponseWriter, r *http.Request)
	FindBook(w http.ResponseWriter, r *http.Request)
	CreateBook(w http.ResponseWriter, r *http.Request)
	UpdateBook(w http.ResponseWriter, r *http.Request)
	DeleteBook(w http.ResponseWriter, r *http.Request)
}

// NewBookController to get a book controller
func NewBookController() *BookController {
	return new(BookController)
}

// FindBooks godoc
// @Summary Get details of all books
// @Description Get details of all books
// @Tags books
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Book
// @Failure 400 {object} map[string]string
// @Router /books [get]
func (b *BookController) FindBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var books []models.Book
	configs.DB.Find(&books)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string][]models.Book{"data": books})
}

// FindBook godoc
// @Summary Get detail of book with given id
// @Description Get detail of book with given id
// @Tags books
// @Accept  json
// @Produce  json
// @Param id path string true "id of the book"
// @Success 200 {object} models.Book
// @Failure 400 {object} map[string]string
// @Router /books/{id} [get]
func (b *BookController) FindBook(w http.ResponseWriter, r *http.Request) {
	// Get model if exist
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	var book models.Book

	if err := configs.DB.Where("id = ?", id).First(&book).Error; err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Record not found!"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]models.Book{"data": book})
}

// CreateBook godoc
// @Summary Create a new book
// @Description Create a new book with the input paylod
// @Tags books
// @Accept  json
// @Produce  json
// @Param book body models.CreateBookInput true "Create book"
// @Success 200 {object} models.CreateBookInput
// @Failure 400 {object} map[string]string
// @Router /books [post]
func (b *BookController) CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Validate input
	var input models.CreateBookInput
	errDecode := json.NewDecoder(r.Body).Decode(&input)
	if errDecode != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": errDecode.Error()})
		return
	}

	errValidation := utils.Validate.Struct(input)
	if errValidation != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": errValidation.Error()})
		return
	}

	// Create book
	book := models.Book{Title: input.Title, Author: input.Author}
	configs.DB.Create(&book)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]models.Book{"data": book})
}

// UpdateBook godoc
// @Summary Update a book
// @Description Update a book with the input paylod
// @Tags books
// @Accept  json
// @Produce  json
// @Param id path string true "id of the book"
// @Param book body models.UpdateBookInput true "Update book"
// @Success 200 {object} models.UpdateBookInput
// @Failure 400 {object} map[string]string
// @Router /books/{id} [patch]
func (b *BookController) UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	// Validate input
	var input models.UpdateBookInput
	errDecode := json.NewDecoder(r.Body).Decode(&input)
	if errDecode != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": errDecode.Error()})
		return
	}

	errValidation := utils.Validate.Struct(input)
	if errValidation != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": errValidation.Error()})
		return
	}

	// Get model if exist
	var book models.Book
	if err := configs.DB.Where("id = ?", id).First(&book).Error; err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Record not found!"})
		return
	}
	configs.DB.Model(&book).Updates(input)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]models.Book{"data": book})
}

// DeleteBook godoc
// @Summary Delete a book
// @Description Delete a book with the input paylod
// @Tags books
// @Accept  json
// @Produce  json
// @Param id path string true "id of the book"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} map[string]string
// @Router /books/{id} [delete]
func (b *BookController) DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	// Get model if exist
	var book models.Book
	if err := configs.DB.Where("id = ?", id).First(&book).Error; err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Record not found!"})
		return
	}

	configs.DB.Delete(&book)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"data": true})
}
