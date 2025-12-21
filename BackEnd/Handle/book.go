package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// getBooks handler untuk ambil semua buku
func getBooks(w http.ResponseWriter, r *http.Request) {
	books, err := GetAllBooks()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// getBook handler untuk ambil detail buku
func getBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	book, err := GetBookByID(bookID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if book == nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// getBooksByCategory handler untuk ambil buku berdasarkan kategori
func getBooksByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryID, err := strconv.Atoi(vars["categoryId"])
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	books, err := GetBooksByCategory(categoryID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// searchBooks handler untuk mencari buku
func searchBooks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "Query parameter 'q' is required"})
		return
	}

	books, err := SearchBooks(query)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// getPopularBooks handler untuk ambil buku populer
func getPopularBooks(w http.ResponseWriter, r *http.Request) {
	limit := 10
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	books, err := GetPopularBooks(limit)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// addBook handler untuk menambah buku baru
func addBook(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title         string `json:"title"`
		Author        string `json:"author"`
		Publisher     string `json:"publisher"`
		YearPublished int    `json:"year_published"`
		ISBN          string `json:"isbn"`
		CategoryID    int    `json:"category_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if !ValidateBookData(req.Title, req.Author, req.ISBN) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid book data"})
		return
	}

	bookID, err := AddBook(req.Title, req.Author, req.Publisher, req.YearPublished, req.ISBN, req.CategoryID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to add book"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"book_id": bookID,
		"message": "Book added successfully",
	})
}

// editBook handler untuk mengubah data buku
func editBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[v0] ========== EDIT BOOK HANDLER CALLED ==========")

	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid book ID"})
		return
	}

	var req struct {
		Title         string `json:"title"`
		Author        string `json:"author"`
		Publisher     string `json:"publisher"`
		YearPublished int    `json:"year_published"`
		ISBN          string `json:"isbn"`
		CategoryID    int    `json:"category_id"`
		Description   string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	fmt.Printf("[v0] Handler received - BookID: %d, Description length: %d, Description: %s\n", bookID, len(req.Description), req.Description)

	if !ValidateBookData(req.Title, req.Author, req.ISBN) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid book data"})
		return
	}

	fmt.Printf("[v0] Calling EditBook with description: %s\n", req.Description)
	err = EditBook(bookID, req.Title, req.Author, req.Publisher, req.ISBN, req.Description, req.YearPublished, req.CategoryID)
	if err != nil {
		fmt.Printf("[v0] ERROR from EditBook: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to edit book"})
		return
	}

	fmt.Printf("[v0] Book updated successfully\n")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Book updated successfully",
	})
}

// deleteBook handler untuk menghapus buku
func deleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid book ID"})
		return
	}

	err = DeleteBook(bookID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete book"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Book deleted successfully",
	})
}

// getLocations handler untuk ambil semua lokasi
func getLocations(w http.ResponseWriter, r *http.Request) {
	locations, err := GetAllLocations()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}

// getLocation handler untuk ambil detail lokasi
func getLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	locationID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid location ID", http.StatusBadRequest)
		return
	}

	location, err := GetLocationByID(locationID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if location == nil {
		http.Error(w, "Location not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(location)
}
