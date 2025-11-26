package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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
	router.HandleFunc("/api/users/{id}/upload-profile", uploadProfileImage).Methods("POST")

	router.HandleFunc("/api/books", getBooks).Methods("GET")
	router.HandleFunc("/api/books", addBook).Methods("POST")
	router.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/api/books/{id}", editBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")
	router.HandleFunc("/api/books/search", searchBooks).Methods("GET")
	router.HandleFunc("/api/books/popular", getPopularBooks).Methods("GET")
	router.HandleFunc("/api/books/category/{categoryId}", getBooksByCategory).Methods("GET")

	router.HandleFunc("/api/locations", getLocations).Methods("GET")
	router.HandleFunc("/api/locations/{id}", getLocation).Methods("GET")

	router.HandleFunc("/api/borrows", createBorrow).Methods("POST")
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
	var req BorrowRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	borrowID, err := CreateBorrow(req.UserID, req.BookLocationID, req.DeliveryType)
	if err != nil {
		http.Error(w, "Failed to create borrow", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"message":   "Borrow created successfully",
		"borrow_id": borrowID,
	})
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
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Profile image uploaded successfully",
		"user_id": userID,
	})
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
		UPDATE book SET title = ?, author = ?, publisher = ?, year_published = ?, isbn = ?, category_id = ?
		WHERE book_id = ?
	`, book.Title, book.Author, book.Publisher, book.YearPublished, book.ISBN, book.CategoryID, bookID)
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
		SELECT b.book_id, b.title, b.author, b.publisher, b.year_published, b.isbn, c.category_id, c.category_name
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
		err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.Publisher, &book.YearPublished, &book.ISBN, &book.CategoryID, &book.CategoryName)
		if err != nil {
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getPopularBooks(w http.ResponseWriter, r *http.Request) {
	books, err := GetAllBooks()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Book status updated successfully",
		"book_id": bookID,
		"status":  req.Status,
	})
}

// ============ DATA STRUCTURES ============

type User struct {
	UserID   int    `json:"user_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Role     string `json:"role"`
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
	BorrowID       int        `json:"borrow_id"`
	UserID         int        `json:"user_id"`
	BookLocationID int        `json:"book_location_id"`
	BorrowDate     time.Time  `json:"borrow_date"`
	ReturnDate     *time.Time `json:"return_date"`
	Status         string     `json:"status"`
	DeliveryType   string     `json:"delivery_type"`
}

type BorrowRequest struct {
	UserID         int    `json:"user_id"`
	BookLocationID int    `json:"book_location_id"`
	DeliveryType   string `json:"delivery_type"`
}

type Review struct {
	ReviewID  int       `json:"review_id"`
	BookID    int       `json:"book_id"`
	UserID    int       `json:"user_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

// ============ DATABASE FUNCTIONS ============

func GetUserByEmail(email string) (*User, error) {
	var user User
	err := db.QueryRow("SELECT user_id, name, email, password, phone, address, role FROM user WHERE email = ?", email).
		Scan(&user.UserID, &user.Name, &user.Email, &user.Password, &user.Phone, &user.Address, &user.Role)
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
	err := db.QueryRow("SELECT user_id, name, email, password, phone, address, role FROM user WHERE user_id = ?", userID).
		Scan(&user.UserID, &user.Name, &user.Email, &user.Password, &user.Phone, &user.Address, &user.Role)
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
		SELECT b.book_id, b.title, b.author, b.publisher, b.year_published, b.isbn, c.category_id, c.category_name
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
		err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.Publisher, &book.YearPublished, &book.ISBN, &book.CategoryID, &book.CategoryName)
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
		SELECT b.book_id, b.title, b.author, b.publisher, b.year_published, b.isbn, c.category_id, c.category_name
		FROM book b
		LEFT JOIN category c ON b.category_id = c.category_id
		WHERE b.book_id = ?
	`, bookID).Scan(&book.BookID, &book.Title, &book.Author, &book.Publisher, &book.YearPublished, &book.ISBN, &book.CategoryID, &book.CategoryName)
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
		SELECT b.book_id, b.title, b.author, b.publisher, b.year_published, b.isbn, c.category_id, c.category_name
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
		err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.Publisher, &book.YearPublished, &book.ISBN, &book.CategoryID, &book.CategoryName)
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

func CreateBorrow(userID, bookLocationID int, deliveryType string) (int, error) {
	result, err := db.Exec(`
		INSERT INTO borrow (user_id, book_location_id, borrow_date, status, delivery_type)
		VALUES (?, ?, NOW(), 'borrowed', ?)
	`, userID, bookLocationID, deliveryType)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func GetBorrowByID(borrowID int) (*Borrow, error) {
	var borrow Borrow
	err := db.QueryRow(`
		SELECT borrow_id, user_id, book_location_id, borrow_date, return_date, status, delivery_type
		FROM borrow WHERE borrow_id = ?
	`, borrowID).Scan(&borrow.BorrowID, &borrow.UserID, &borrow.BookLocationID, &borrow.BorrowDate, &borrow.ReturnDate, &borrow.Status, &borrow.DeliveryType)
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
		SELECT borrow_id, user_id, book_location_id, borrow_date, return_date, status, delivery_type
		FROM borrow WHERE user_id = ? ORDER BY borrow_date DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var borrows []Borrow
	for rows.Next() {
		var borrow Borrow
		err := rows.Scan(&borrow.BorrowID, &borrow.UserID, &borrow.BookLocationID, &borrow.BorrowDate, &borrow.ReturnDate, &borrow.Status, &borrow.DeliveryType)
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

func GetBookReviews(bookID int) ([]Review, error) {
	rows, err := db.Query(`
		SELECT review_id, book_id, user_id, rating, comment, created_at
		FROM review WHERE book_id = ? ORDER BY created_at DESC
	`, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []Review
	for rows.Next() {
		var review Review
		err := rows.Scan(&review.ReviewID, &review.BookID, &review.UserID, &review.Rating, &review.Comment, &review.CreatedAt)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}
