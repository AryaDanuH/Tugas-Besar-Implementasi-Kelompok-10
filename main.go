package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var db *sql.DB

func init() {
	godotenv.Load()

	var err error
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "root:@tcp(127.0.0.1:3306)/libmatch?parseTime=true"
	}

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Database connection error:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Database ping error:", err)
	}
	fmt.Println("Database connected successfully")
}

func main() {
	router := mux.NewRouter()

	router.Use(corsMiddleware)

	router.HandleFunc("/api/auth/register", registerUser).Methods("POST")
	router.HandleFunc("/api/auth/login", loginUser).Methods("POST")
	router.HandleFunc("/api/auth/change-password", changePassword).Methods("POST")
	router.HandleFunc("/api/auth/forgot-password", forgotPassword).Methods("POST")
	router.HandleFunc("/api/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/api/users/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/api/users/{id}/change-username", changeUsername).Methods("PUT")
	router.HandleFunc("/api/users/{id}/upload-profile-image", uploadProfileImage).Methods("POST")

	router.HandleFunc("/api/books", getBooks).Methods("GET")
	router.HandleFunc("/api/books", addBook).Methods("POST")
	router.HandleFunc("/api/books/upload", uploadBook).Methods("POST")
	router.HandleFunc("/api/users/{userId}/borrowed-books", getUserBorrowedBooks).Methods("GET")
	router.HandleFunc("/api/users/{userId}/books", getUserBooks).Methods("GET")
	router.HandleFunc("/api/books/pending", getPendingBooks).Methods("GET")
	router.HandleFunc("/api/books/accepted", getAcceptedBooks).Methods("GET")
	router.HandleFunc("/api/books/new-arrivals", getNewArrivals).Methods("GET")
	router.HandleFunc("/api/books/search", searchBooks).Methods("GET")
	router.HandleFunc("/api/books/popular", getMostViewedBooks).Methods("GET") // Changed from getPopularBooks
	router.HandleFunc("/api/books/top-borrowed", getTopBorrowedBooks).Methods("GET")
	router.HandleFunc("/api/books/category/{categoryId}", getBooksByCategory).Methods("GET")
	router.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/api/books/{id}", editBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")
	router.HandleFunc("/api/books/{bookId}/view", incrementBookView).Methods("POST") // Added new route

	router.HandleFunc("/api/categories", getCategories).Methods("GET")

	router.HandleFunc("/api/locations", getLocations).Methods("GET")
	router.HandleFunc("/api/locations/{id}", getLocation).Methods("GET")

	router.HandleFunc("/api/borrows", createBorrow).Methods("POST")
	router.HandleFunc("/api/borrows", getBorrows).Methods("GET")
	router.HandleFunc("/api/borrows/{id}", getBorrow).Methods("GET")
	router.HandleFunc("/api/borrows/user/{userId}", getUserBorrows).Methods("GET")
	router.HandleFunc("/api/borrows/{id}/approve", approveBorrow).Methods("PUT")
	router.HandleFunc("/api/borrows/{id}/reject", rejectBorrow).Methods("PUT")
	router.HandleFunc("/api/borrows/{id}/return", returnBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}/status", updateBookStatus).Methods("PUT")

	router.HandleFunc("/api/reviews", createReview).Methods("POST")
	router.HandleFunc("/api/reviews/book/{bookId}", getBookReviews).Methods("GET")

	router.PathPrefix("/FrontEnd/").Handler(http.StripPrefix("/FrontEnd/", http.FileServer(http.Dir("FrontEnd"))))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "FrontEnd/index.html")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ============ TYPE DEFINITIONS ============

type User struct {
	UserID       int    `json:"user_id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	Role         string `json:"role"`
	ProfileImage string `json:"profile_image"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Book struct {
	BookID        int    `json:"book_id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	Publisher     string `json:"publisher"`
	YearPublished int    `json:"year_published"`
	ISBN          string `json:"isbn"`
	CategoryID    int    `json:"category_id"`
	CategoryName  string `json:"category_name"`
	UploadedBy    int    `json:"uploaded_by"`
	UploaderName  string `json:"uploader_name"`
	UploaderEmail string `json:"uploader_email"`
	UploaderPhone string `json:"uploader_phone"`
	Description   string `json:"description"`
	CoverImage    string `json:"cover_image"`
	Location      string `json:"location"`
	Status        string `json:"status"`
	Views         int    `json:"views"` // Added Views field
}

type Category struct {
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
}

type Location struct {
	LocationID   int    `json:"location_id"`
	LocationName string `json:"location_name"`
	Address      string `json:"address"`
	OwnerID      int    `json:"owner_id"`
}

type BookLocation struct {
	BookLocationID int `json:"book_location_id"`
	BookID         int `json:"book_id"`
	LocationID     int `json:"location_id"`
	Stock          int `json:"stock"`
}

type Borrow struct {
	BorrowID        int            `json:"borrow_id"`
	UserID          int            `json:"user_id"`
	BookID          int            `json:"book_id"`
	BorrowDate      string         `json:"borrow_date"`
	ReturnDate      *string        `json:"return_date"`
	DueDate         string         `json:"due_date"`
	Status          string         `json:"status"`
	DeliveryType    string         `json:"delivery_type"`
	PickupLocation  sql.NullString `json:"pickup_location"`
	DeliveryAddress sql.NullString `json:"delivery_address"`
	TotalPrice      float64        `json:"total_price"`
}

type BorrowRequest struct {
	UserID       int    `json:"user_id"`
	BookID       int    `json:"book_id"`
	DeliveryType string `json:"delivery_type"`
}

type Review struct {
	ReviewID  int       `json:"review_id"`
	BookID    int       `json:"book_id"`
	UserID    int       `json:"user_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

type BorrowedBookDetail struct {
	TransactionID int        `json:"transaction_id"`
	BorrowDate    time.Time  `json:"borrow_date"`
	ReturnDate    *time.Time `json:"return_date"`
	Status        string     `json:"status"`
	DeliveryType  string     `json:"delivery_type"`
	BookID        int        `json:"book_id"`
	Title         string     `json:"title"`
	Author        string     `json:"author"`
	CoverImage    string     `json:"cover_image"`
	Description   string     `json:"description"`
	YearPublished int        `json:"year_published"`
	Category      string     `json:"category"`
}

// ============ USER HANDLERS ============

func registerUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	existingUser, err := GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if existingUser != nil {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	userID, err := CreateUser(req.Name, req.Email, req.Password, req.Phone, req.Address)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User registered successfully",
		"user_id": userID,
	})
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if user.Password != req.Password {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Login successful",
		"user":    user,
	})
}

func changePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Email       string `json:"email"`
		NewPassword string `json:"newPassword"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid request",
		})
		return
	}

	if req.Email == "" || req.NewPassword == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Email dan password harus diisi",
		})
		return
	}

	if len(req.NewPassword) < 6 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Password minimal 6 karakter",
		})
		return
	}

	user, err := GetUserByEmail(req.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Database error",
		})
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Email tidak ditemukan",
		})
		return
	}

	if user.Password == req.NewPassword {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Password baru tidak boleh sama dengan password lama",
		})
		return
	}

	err = UpdateUserPassword(user.UserID, req.NewPassword)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Gagal mengubah password",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Password berhasil diubah",
	})
}

func forgotPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		Email string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid request",
		})
		return
	}

	if req.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Email is required",
		})
		return
	}

	user, err := GetUserByEmail(req.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Database error",
		})
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Email not found",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Password reset link sent to email",
	})
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := GetUserByID(userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err = UpdateUser(userID, user.Name, user.Phone, user.Address)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User updated successfully",
	})
}

func changeUsername(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		NewName string `json:"newName"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err = UpdateUser(userID, req.NewName, "", "")
	if err != nil {
		http.Error(w, "Failed to update username", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Username updated successfully",
	})
}

func uploadProfileImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid user ID",
		})
		return
	}

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to parse form data",
		})
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "No file uploaded",
		})
		return
	}
	defer file.Close()

	uploadsDir := "./FrontEnd/uploads/profiles"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		os.MkdirAll(uploadsDir, 0755)
	}

	ext := strings.TrimPrefix(filepath.Ext(handler.Filename), ".")
	if ext == "" {
		ext = "jpg"
	}
	filename := fmt.Sprintf("profile_%d_%d.%s", userID, time.Now().Unix(), ext)
	filePath := filepath.Join(uploadsDir, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to save profile image",
		})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to write profile image",
		})
		return
	}

	profileImagePath := "/FrontEnd/uploads/profiles/" + filename

	_, err = db.Exec("UPDATE user SET profile_image = ? WHERE user_id = ?", profileImagePath, userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to update database",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":       true,
		"message":       "Profile image uploaded successfully",
		"profile_image": profileImagePath,
	})
}

// ============ BOOK HANDLERS ============

func getBooks(w http.ResponseWriter, r *http.Request) {
	books, err := GetAllBooks()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getUserBooks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	books, err := GetBooksByUploader(userID)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if books == nil {
		books = []Book{}
	}
	json.NewEncoder(w).Encode(books)
}

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

func addBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`
		INSERT INTO book (title, author, publisher, year_published, isbn, category_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`, book.Title, book.Author, book.Publisher, book.YearPublished, book.ISBN, book.CategoryID)
	if err != nil {
		http.Error(w, "Failed to add book", http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Book added successfully",
		"book_id": int(id),
	})
}

func uploadBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to parse form data",
		})
		return
	}

	title := r.FormValue("title")
	author := r.FormValue("author")
	yearPublishedStr := r.FormValue("year_published")
	categoryIDStr := r.FormValue("category_id")
	uploaderName := r.FormValue("uploader_name")
	uploaderEmail := r.FormValue("uploader_email")
	uploaderPhone := r.FormValue("uploader_phone")
	uploadedByStr := r.FormValue("uploaded_by")
	description := r.FormValue("description")
	publisher := r.FormValue("publisher")
	isbn := r.FormValue("isbn")
	location := r.FormValue("location")

	if title == "" || author == "" || categoryIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Title, author, and category are required",
		})
		return
	}

	categoryID := 0
	if categoryIDStr != "" {
		fmt.Sscanf(categoryIDStr, "%d", &categoryID)
	}
	yearPublished := 0
	if yearPublishedStr != "" {
		fmt.Sscanf(yearPublishedStr, "%d", &yearPublished)
	}
	uploadedBy := 0
	if uploadedByStr != "" {
		fmt.Sscanf(uploadedByStr, "%d", &uploadedBy)
	}

	var coverImagePath string
	file, handler, err := r.FormFile("cover_image")
	if err == nil {
		defer file.Close()

		uploadsDir := "./FrontEnd/uploads"
		if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
			os.MkdirAll(uploadsDir, 0755)
		}

		ext := strings.TrimPrefix(filepath.Ext(handler.Filename), ".")
		if ext == "" {
			ext = "jpg"
		}
		filename := fmt.Sprintf("book_cover_%d_%d.%s", time.Now().Unix(), rand.Intn(10000), ext)
		filepath := filepath.Join(uploadsDir, filename)

		dst, err := os.Create(filepath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Failed to save cover image",
			})
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Failed to write cover image",
			})
			return
		}

		coverImagePath = "/FrontEnd/uploads/" + filename
	}

	status := "pending"
	if uploadedBy == 11 {
		status = "accepted"
	}

	_, err = db.Exec(`
		INSERT INTO book (title, author, publisher, year_published, isbn, category_id, uploaded_by, uploader_name, uploader_email, uploader_phone, description, cover_image, location, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, title, author, publisher, yearPublished, isbn, categoryID, uploadedBy, uploaderName, uploaderEmail, uploaderPhone, description, coverImagePath, location, status)

	if err != nil {
		log.Printf("Database error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Failed to save book to database: " + err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"message":     "Book uploaded successfully and pending approval",
		"cover_image": coverImagePath,
	})
}

func editBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var book Book
	err = json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`
		UPDATE book SET title = ?, author = ?, publisher = ?, year_published = ?, isbn = ?, category_id = ?, description = ?
		WHERE book_id = ?
	`, book.Title, book.Author, book.Publisher, book.YearPublished, book.ISBN, book.CategoryID, book.Description, bookID)
	if err != nil {
		http.Error(w, "Failed to edit book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Book updated successfully",
	})
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM book WHERE book_id = ?", bookID)
	if err != nil {
		http.Error(w, "Failed to delete book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Book deleted successfully",
	})
}

func searchBooks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Search query required", http.StatusBadRequest)
		return
	}

	rows, err := db.Query(`
		SELECT b.book_id, b.title, b.author, COALESCE(b.publisher, ''), b.year_published, COALESCE(b.isbn, ''),
		       b.category_id, c.category_name, COALESCE(b.cover_image, '')
		FROM book b
		LEFT JOIN category c ON b.category_id = c.category_id
		WHERE b.title LIKE ? OR b.author LIKE ?
	`, "%"+query+"%", "%"+query+"%")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.Publisher, &book.YearPublished, &book.ISBN, &book.CategoryID, &book.CategoryName, &book.CoverImage)
		if err != nil {
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// getMostViewedBooks returns books sorted by view count
func getMostViewedBooks(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT b.book_id, b.title, b.author, b.publisher, b.year_published,
		       b.isbn, b.category_id, c.category_name, b.uploaded_by,
		       u.name AS uploader_name, u.email AS uploader_email,
		       u.phone AS uploader_phone, b.description, b.cover_image,
		       b.location, b.status, b.views
		FROM book b
		LEFT JOIN category c ON b.category_id = c.category_id
		LEFT JOIN user u ON b.uploaded_by = u.user_id
		WHERE b.status = 'accepted'
		ORDER BY b.views DESC
		LIMIT 10
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error querying most viewed books: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(
			&book.BookID, &book.Title, &book.Author, &book.Publisher,
			&book.YearPublished, &book.ISBN, &book.CategoryID, &book.CategoryName,
			&book.UploadedBy, &book.UploaderName, &book.UploaderEmail,
			&book.UploaderPhone, &book.Description, &book.CoverImage,
			&book.Location, &book.Status, &book.Views,
		)
		if err != nil {
			log.Printf("Error scanning book: %v", err)
			continue
		}
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// getPopularBooks returns books sorted by borrow count (most engaged/clicked books)
func getPopularBooks(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT b.book_id, b.title, b.author, b.publisher, b.year_published,
		       b.isbn, b.category_id, c.category_name, b.uploaded_by,
		       u.name AS uploader_name, u.email AS uploader_email,
		       u.phone AS uploader_phone, b.description, b.cover_image,
		       b.location, b.status, b.views
		FROM book b
		LEFT JOIN category c ON b.category_id = c.category_id
		LEFT JOIN user u ON b.uploaded_by = u.user_id
		WHERE b.status = 'accepted'
		ORDER BY b.views DESC
		LIMIT 10
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error querying popular books: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(
			&book.BookID, &book.Title, &book.Author, &book.Publisher,
			&book.YearPublished, &book.ISBN, &book.CategoryID, &book.CategoryName,
			&book.UploadedBy, &book.UploaderName, &book.UploaderEmail,
			&book.UploaderPhone, &book.Description, &book.CoverImage,
			&book.Location, &book.Status, &book.Views,
		)
		if err != nil {
			log.Printf("Error scanning book: %v", err)
			continue
		}
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// getTopBorrowedBooks returns most borrowed books
func getTopBorrowedBooks(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT b.book_id, b.title, b.author, b.publisher, b.year_published,
		       b.isbn, b.category_id, c.category_name, b.uploaded_by,
		       u.name AS uploader_name, u.email AS uploader_email,
		       u.phone AS uploader_phone, b.description, b.cover_image,
		       b.location, b.status, COUNT(br.borrow_id) as borrow_count
		FROM book b
		LEFT JOIN category c ON b.category_id = c.category_id
		LEFT JOIN user u ON b.uploaded_by = u.user_id
		LEFT JOIN borrow br ON b.book_id = br.book_id
		WHERE b.status = 'accepted'
		GROUP BY b.book_id
		ORDER BY borrow_count DESC
		LIMIT 10
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error querying top borrowed books: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		var borrowCount int
		err := rows.Scan(
			&book.BookID, &book.Title, &book.Author, &book.Publisher,
			&book.YearPublished, &book.ISBN, &book.CategoryID, &book.CategoryName,
			&book.UploadedBy, &book.UploaderName, &book.UploaderEmail,
			&book.UploaderPhone, &book.Description, &book.CoverImage,
			&book.Location, &book.Status, &borrowCount,
		)
		if err != nil {
			log.Printf("Error scanning book: %v", err)
			continue
		}
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getPendingBooks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT b.book_id, b.title, b.author, COALESCE(b.publisher, ''), b.year_published, COALESCE(b.isbn, ''),
		       b.category_id, c.category_name, COALESCE(b.uploaded_by, 0), COALESCE(b.uploader_name, ''),
		       COALESCE(b.uploader_email, ''), COALESCE(b.uploader_phone, ''), COALESCE(b.description, ''),
		       COALESCE(b.cover_image, ''), COALESCE(b.location, ''), COALESCE(b.status, 'pending')
		FROM book b
		LEFT JOIN category c ON b.category_id = c.category_id
		WHERE b.status = 'pending'
		ORDER BY b.book_id DESC
	`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Failed to fetch pending books",
		})
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.Publisher, &book.YearPublished, &book.ISBN,
			&book.CategoryID, &book.CategoryName, &book.UploadedBy, &book.UploaderName,
			&book.UploaderEmail, &book.UploaderPhone, &book.Description, &book.CoverImage,
			&book.Location, &book.Status)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Failed to scan book data",
			})
			return
		}
		books = append(books, book)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getAcceptedBooks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT b.book_id, b.title, b.author, COALESCE(b.publisher, ''), b.year_published, COALESCE(b.isbn, ''),
		       b.category_id, c.category_name, COALESCE(b.uploaded_by, 0), COALESCE(b.uploader_name, ''),
		       COALESCE(b.uploader_email, ''), COALESCE(b.uploader_phone, ''), COALESCE(b.description, ''),
		       COALESCE(b.cover_image, ''), COALESCE(b.location, ''), COALESCE(b.status, 'pending')
		FROM book b
		LEFT JOIN category c ON b.category_id = c.category_id
		WHERE b.status = 'accepted'
		ORDER BY b.book_id DESC
	`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Failed to fetch accepted books",
		})
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.Publisher, &book.YearPublished, &book.ISBN,
			&book.CategoryID, &book.CategoryName, &book.UploadedBy, &book.UploaderName,
			&book.UploaderEmail, &book.UploaderPhone, &book.Description, &book.CoverImage,
			&book.Location, &book.Status)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Failed to scan book data",
			})
			return
		}
		books = append(books, book)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getNewArrivals(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT b.book_id, b.title, b.author, COALESCE(b.publisher, ''), b.year_published, COALESCE(b.isbn, ''),
		       b.category_id, c.category_name, COALESCE(b.uploaded_by, 0), COALESCE(b.uploader_name, ''),
		       COALESCE(b.uploader_email, ''), COALESCE(b.uploader_phone, ''), COALESCE(b.description, ''),
		       COALESCE(b.cover_image, ''), COALESCE(b.location, ''), COALESCE(b.status, 'pending')
		FROM book b
		LEFT JOIN category c ON b.category_id = c.category_id
		WHERE b.status = 'accepted'
		ORDER BY b.book_id DESC
	`)
	if err != nil {
		log.Printf("Error fetching new arrivals: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.Publisher, &book.YearPublished, &book.ISBN,
			&book.CategoryID, &book.CategoryName, &book.UploadedBy, &book.UploaderName,
			&book.UploaderEmail, &book.UploaderPhone, &book.Description, &book.CoverImage,
			&book.Location, &book.Status)
		if err != nil {
			log.Printf("Error scanning book row: %v", err)
			http.Error(w, "Database scan error", http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	log.Printf("Successfully fetched %d new arrival books", len(books))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func updateBookStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Status != "pending" && req.Status != "accepted" && req.Status != "rejected" {
		http.Error(w, "Invalid status. Must be 'pending', 'accepted', or 'rejected'", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE book SET status = ? WHERE book_id = ?", req.Status, bookID)
	if err != nil {
		log.Printf("Error updating book status: %v", err)
		http.Error(w, "Failed to update book status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Book status updated successfully",
		"book_id": bookID,
		"status":  req.Status,
	})
}

func getUserBorrowedBooks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	rows, err := db.Query(`
		SELECT
			b.borrow_id,
			b.book_id,
			bk.title,
			bk.author,
			bk.cover_image,
			bk.description,
			bk.year_published,
			c.category_name,
			b.borrow_date,
			b.due_date,
			b.return_date,
			b.status,
			b.delivery_type,
			b.pickup_location,
			b.delivery_address,
			b.total_price
		FROM borrow b
		JOIN book bk ON b.book_id = bk.book_id
		LEFT JOIN category c ON bk.category_id = c.category_id
		WHERE b.user_id = ?
		ORDER BY b.borrow_date DESC
	`, userID)
	if err != nil {
		log.Printf("Error fetching borrowed books: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var borrowedBooks []map[string]interface{}
	for rows.Next() {
		var borrowID, bookID, yearPublished int
		var title, author, coverImage, description, categoryName, borrowDate, dueDate, status, deliveryType string
		var returnDate, pickupLocation, deliveryAddress sql.NullString
		var totalPrice float64

		err := rows.Scan(&borrowID, &bookID, &title, &author, &coverImage, &description, &yearPublished,
			&categoryName, &borrowDate, &dueDate, &returnDate, &status, &deliveryType, &pickupLocation, &deliveryAddress, &totalPrice)
		if err != nil {
			log.Printf("Error scanning borrowed book row: %v", err)
			http.Error(w, "Database scan error", http.StatusInternalServerError)
			return
		}

		borrowedBook := map[string]interface{}{
			"borrow_id":        borrowID,
			"book_id":          bookID,
			"title":            title,
			"author":           author,
			"cover_image":      coverImage,
			"description":      description,
			"year_published":   yearPublished,
			"category":         categoryName,
			"borrow_date":      borrowDate,
			"due_date":         dueDate,
			"return_date":      returnDate,
			"status":           status,
			"delivery_type":    deliveryType,
			"pickup_location":  pickupLocation,
			"delivery_address": deliveryAddress,
			"total_price":      totalPrice,
		}
		borrowedBooks = append(borrowedBooks, borrowedBook)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(borrowedBooks)
}

// ============ LOCATION HANDLERS ============

func getLocations(w http.ResponseWriter, r *http.Request) {
	locations, err := GetAllLocations()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}

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

// ============ BORROW HANDLERS ============

func createBorrow(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID       int    `json:"user_id"`
		BookID       int    `json:"book_id"`
		DeliveryType string `json:"delivery_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Determine borrowDate, dueDate, and totalPrice based on deliveryType
	borrowDate := time.Now()
	dueDate := borrowDate.AddDate(0, 0, 14) // 14 days from borrow date
	var totalPrice float64

	if req.DeliveryType == "delivery" {
		totalPrice = 113300.00 // Example price for delivery
	} else { // Pickup
		totalPrice = 88000.00 // Example price for pickup
	}

	borrowID, err := CreateBorrow(req.UserID, req.BookID, borrowDate, dueDate, req.DeliveryType, totalPrice)
	if err != nil {
		http.Error(w, "Failed to create borrow", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"borrow_id": borrowID,
		"message":   "Book borrowed successfully",
	})
}

func getBorrows(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT borrow_id, book_id, user_id, borrow_date, return_date, due_date, status
		FROM borrow
		ORDER BY borrow_date DESC
	`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Failed to fetch borrows",
		})
		return
	}
	defer rows.Close()

	var borrows []struct {
		BorrowID   int    `json:"borrow_id"`
		BookID     int    `json:"book_id"`
		UserID     int    `json:"user_id"`
		BorrowDate string `json:"borrow_date"`
		ReturnDate string `json:"return_date"`
		DueDate    string `json:"due_date"`
		Status     string `json:"status"`
	}

	for rows.Next() {
		var borrow struct {
			BorrowID   int    `json:"borrow_id"`
			BookID     int    `json:"book_id"`
			UserID     int    `json:"user_id"`
			BorrowDate string `json:"borrow_date"`
			ReturnDate string `json:"return_date"`
			DueDate    string `json:"due_date"`
			Status     string `json:"status"`
		}
		err := rows.Scan(&borrow.BorrowID, &borrow.BookID, &borrow.UserID,
			&borrow.BorrowDate, &borrow.ReturnDate, &borrow.DueDate, &borrow.Status)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Failed to scan borrow data",
			})
			return
		}
		borrows = append(borrows, borrow)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(borrows)
}

func getBorrow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	borrowID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid borrow ID", http.StatusBadRequest)
		return
	}

	borrow, err := GetBorrowByID(borrowID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if borrow == nil {
		http.Error(w, "Borrow not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(borrow)
}

func getUserBorrows(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	borrows, err := GetUserBorrows(userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(borrows)
}

func approveBorrow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	borrowID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid borrow ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE borrow SET status = 'approved' WHERE borrow_id = ?", borrowID)
	if err != nil {
		http.Error(w, "Failed to approve borrow", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Borrow approved successfully",
	})
}

func rejectBorrow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	borrowID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid borrow ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE borrow SET status = 'rejected' WHERE borrow_id = ?", borrowID)
	if err != nil {
		http.Error(w, "Failed to reject borrow", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Borrow rejected successfully",
	})
}

func returnBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	borrowID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid borrow ID", http.StatusBadRequest)
		return
	}

	err = ReturnBook(borrowID)
	if err != nil {
		http.Error(w, "Failed to return book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Book returned successfully",
	})
}

// ============ REVIEW HANDLERS ============

func createReview(w http.ResponseWriter, r *http.Request) {
	var review Review
	err := json.NewDecoder(r.Body).Decode(&review)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	reviewID, err := CreateReview(review.BookID, review.UserID, review.Rating, review.Comment)
	if err != nil {
		http.Error(w, "Failed to create review", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"message":   "Review created successfully",
		"review_id": reviewID,
	})
}

func getBookReviews(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["bookId"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	reviews, err := GetBookReviews(bookID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

// ============ CATEGORY HANDLERS ============

// getCategories fetches all categories from the database
func getCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT category_id, category_name FROM category ORDER BY category_name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var categories []map[string]interface{}
	for rows.Next() {
		var categoryID int
		var categoryName string

		if err := rows.Scan(&categoryID, &categoryName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		categories = append(categories, map[string]interface{}{
			"category_id":   categoryID,
			"category_name": categoryName,
		})
	}

	json.NewEncoder(w).Encode(categories)
}

// ============ DATABASE FUNCTIONS ============

func GetUserByEmail(email string) (*User, error) {
	var user User
	err := db.QueryRow("SELECT user_id, name, email, password, phone, address, role, COALESCE(profile_image, '') FROM user WHERE email = ?", email).
		Scan(&user.UserID, &user.Name, &user.Email, &user.Password, &user.Phone, &user.Address, &user.Role, &user.ProfileImage)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(userID int) (*User, error) {
	var user User
	err := db.QueryRow("SELECT user_id, name, email, password, phone, address, role, COALESCE(profile_image, '') FROM user WHERE user_id = ?", userID).
		Scan(&user.UserID, &user.Name, &user.Email, &user.Password, &user.Phone, &user.Address, &user.Role, &user.ProfileImage)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(name, email, password, phone, address string) (int, error) {
	result, err := db.Exec("INSERT INTO user (name, email, password, phone, address, role) VALUES (?, ?, ?, ?, ?, ?)",
		name, email, password, phone, address, "member")
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func UpdateUser(userID int, name, phone, address string) error {
	_, err := db.Exec("UPDATE user SET name = ?, phone = ?, address = ? WHERE user_id = ?",
		name, phone, address, userID)
	return err
}

func UpdateUserPassword(userID int, newPassword string) error {
	_, err := db.Exec("UPDATE user SET password = ? WHERE user_id = ?", newPassword, userID)
	return err
}

func GetAllBooks() ([]Book, error) {
	rows, err := db.Query(`
		SELECT b.book_id, b.title, b.author, COALESCE(b.publisher, ''), b.year_published, COALESCE(b.isbn, ''),
		       b.category_id, c.category_name, COALESCE(b.uploaded_by, 0), COALESCE(b.uploader_name, ''),
		       COALESCE(b.uploader_email, ''), COALESCE(b.uploader_phone, ''), COALESCE(b.description, ''),
		       COALESCE(b.cover_image, ''), COALESCE(b.location, ''), COALESCE(b.status, 'pending'), COALESCE(b.views, 0)
		FROM book b
		LEFT JOIN category c ON b.category_id = c.category_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.Publisher, &book.YearPublished, &book.ISBN,
			&book.CategoryID, &book.CategoryName, &book.UploadedBy, &book.UploaderName,
			&book.UploaderEmail, &book.UploaderPhone, &book.Description, &book.CoverImage,
			&book.Location, &book.Status, &book.Views)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func GetBookByID(bookID int) (*Book, error) {
	var book Book
	err := db.QueryRow(`
		SELECT b.book_id, b.title, b.author, COALESCE(b.publisher, ''), b.year_published, COALESCE(b.isbn, ''),
		       b.category_id, c.category_name, COALESCE(b.uploaded_by, 0), COALESCE(b.uploader_name, ''),
		       COALESCE(b.uploader_email, ''), COALESCE(b.uploader_phone, ''), COALESCE(b.description, ''),
		       COALESCE(b.cover_image, ''), COALESCE(b.location, ''), COALESCE(b.status, 'pending'), COALESCE(b.views, 0)
		FROM book b
		LEFT JOIN category c ON b.category_id = c.category_id
		WHERE b.book_id = ?
	`, bookID).Scan(&book.BookID, &book.Title, &book.Author, &book.Publisher, &book.YearPublished, &book.ISBN,
		&book.CategoryID, &book.CategoryName, &book.UploadedBy, &book.UploaderName,
		&book.UploaderEmail, &book.UploaderPhone, &book.Description, &book.CoverImage,
		&book.Location, &book.Status, &book.Views)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func GetBooksByCategory(categoryID int) ([]Book, error) {
	rows, err := db.Query(`
		SELECT b.book_id, b.title, b.author, COALESCE(b.publisher, ''), b.year_published, COALESCE(b.isbn, ''),
		       b.category_id, c.category_name, COALESCE(b.uploaded_by, 0), COALESCE(b.uploader_name, ''),
		       COALESCE(b.uploader_email, ''), COALESCE(b.uploader_phone, ''), COALESCE(b.description, ''),
		       COALESCE(b.cover_image, ''), COALESCE(b.location, ''), COALESCE(b.status, 'pending'), COALESCE(b.views, 0)
		FROM book b
		LEFT JOIN category c ON b.category_id = c.category_id
		WHERE b.category_id = ?
	`, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.Publisher, &book.YearPublished, &book.ISBN,
			&book.CategoryID, &book.CategoryName, &book.UploadedBy, &book.UploaderName,
			&book.UploaderEmail, &book.UploaderPhone, &book.Description, &book.CoverImage,
			&book.Location, &book.Status, &book.Views)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func GetBooksByUploader(userID int) ([]Book, error) {
	rows, err := db.Query(`
		SELECT b.book_id, b.title, b.author, COALESCE(b.publisher, ''), b.year_published, COALESCE(b.isbn, ''),
		       b.category_id, c.category_name, COALESCE(b.uploaded_by, 0), COALESCE(b.uploader_name, ''),
		       COALESCE(b.uploader_email, ''), COALESCE(b.uploader_phone, ''), COALESCE(b.description, ''),
		       COALESCE(b.cover_image, ''), COALESCE(b.location, ''), COALESCE(b.status, 'pending'), COALESCE(b.views, 0)
		FROM book b
		LEFT JOIN category c ON b.category_id = c.category_id
		WHERE b.uploaded_by = ?
		ORDER BY b.book_id DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.Publisher, &book.YearPublished, &book.ISBN,
			&book.CategoryID, &book.CategoryName, &book.UploadedBy, &book.UploaderName,
			&book.UploaderEmail, &book.UploaderPhone, &book.Description, &book.CoverImage,
			&book.Location, &book.Status, &book.Views)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func GetAllLocations() ([]Location, error) {
	rows, err := db.Query("SELECT location_id, location_name, address, owner_id FROM location")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []Location
	for rows.Next() {
		var loc Location
		err := rows.Scan(&loc.LocationID, &loc.LocationName, &loc.Address, &loc.OwnerID)
		if err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}
	return locations, nil
}

func GetLocationByID(locationID int) (*Location, error) {
	var loc Location
	err := db.QueryRow("SELECT location_id, location_name, address, owner_id FROM location WHERE location_id = ?", locationID).
		Scan(&loc.LocationID, &loc.LocationName, &loc.Address, &loc.OwnerID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &loc, nil
}

func CreateBorrow(userID, bookID int, borrowDate, dueDate time.Time, deliveryType string, totalPrice float64) (int, error) {
	result, err := db.Exec(`
		INSERT INTO borrow (user_id, book_id, borrow_date, due_date, status, delivery_type, total_price)
		VALUES (?, ?, ?, ?, 'active', ?, ?)
	`, userID, bookID, borrowDate, dueDate, deliveryType, totalPrice)

	if err != nil {
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastInsertID), nil
}

func GetBorrowByID(borrowID int) (*Borrow, error) {
	var borrow Borrow
	err := db.QueryRow(`
		SELECT borrow_id, user_id, book_id, borrow_date, due_date, return_date, status, delivery_type, COALESCE(pickup_location, ''), COALESCE(delivery_address, ''), total_price
		FROM borrow WHERE borrow_id = ?
	`, borrowID).Scan(&borrow.BorrowID, &borrow.UserID, &borrow.BookID, &borrow.BorrowDate, &borrow.DueDate, &borrow.ReturnDate, &borrow.Status, &borrow.DeliveryType, &borrow.PickupLocation, &borrow.DeliveryAddress, &borrow.TotalPrice)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &borrow, nil
}

func GetUserBorrows(userID int) ([]Borrow, error) {
	rows, err := db.Query(`
		SELECT borrow_id, user_id, book_id, borrow_date, due_date, return_date, status, delivery_type, COALESCE(pickup_location, ''), COALESCE(delivery_address, ''), total_price
		FROM borrow WHERE user_id = ? ORDER BY borrow_date DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var borrows []Borrow
	for rows.Next() {
		var borrow Borrow
		err := rows.Scan(&borrow.BorrowID, &borrow.UserID, &borrow.BookID, &borrow.BorrowDate, &borrow.DueDate, &borrow.ReturnDate, &borrow.Status, &borrow.DeliveryType, &borrow.PickupLocation, &borrow.DeliveryAddress, &borrow.TotalPrice)
		if err != nil {
			return nil, err
		}
		borrows = append(borrows, borrow)
	}
	return borrows, nil
}

func ReturnBook(borrowID int) error {
	_, err := db.Exec("UPDATE borrow SET return_date = NOW(), status = 'returned' WHERE borrow_id = ?", borrowID)
	return err
}

func CreateReview(bookID, userID, rating int, comment string) (int, error) {
	result, err := db.Exec(`
		INSERT INTO review (book_id, user_id, rating, comment, created_at)
		VALUES (?, ?, ?, ?, NOW())
	`, bookID, userID, rating, comment)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func GetBookReviews(bookID int) ([]map[string]interface{}, error) {
	rows, err := db.Query(`
		SELECT r.review_id, r.book_id, r.user_id, u.name, r.rating, r.comment, r.created_at
		FROM review r
		LEFT JOIN user u ON r.user_id = u.user_id
		WHERE r.book_id = ?
		ORDER BY r.created_at DESC
	`, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []map[string]interface{}
	for rows.Next() {
		var reviewID, bookID, userID, rating int
		var username, comment string
		var createdAt time.Time

		err := rows.Scan(&reviewID, &bookID, &userID, &username, &rating, &comment, &createdAt)
		if err != nil {
			return nil, err
		}

		review := map[string]interface{}{
			"review_id":  reviewID,
			"book_id":    bookID,
			"user_id":    userID,
			"username":   username,
			"rating":     rating,
			"comment":    comment,
			"created_at": createdAt,
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

// ============ BOOK VIEW HANDLERS ============

// incrementBookView increments the view count for a specific book.
func incrementBookView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["bookId"])
	if err != nil {
		log.Printf("[DEBUG] Invalid book ID: %v", err)
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	log.Printf("[DEBUG] Incrementing view for book ID: %d", bookID)

	result, err := db.Exec("UPDATE book SET views = views + 1 WHERE book_id = ?", bookID)
	if err != nil {
		log.Printf("Error incrementing book view: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("[DEBUG] Rows affected: %d", rowsAffected)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "View counted"})
}
