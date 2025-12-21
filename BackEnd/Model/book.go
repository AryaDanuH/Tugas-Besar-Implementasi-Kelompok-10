package main

import (
	"database/sql"
	"fmt"
)

var db *sql.DB

type Book struct {
	BookID        int    `json:"book_id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	Publisher     string `json:"publisher"`
	YearPublished int    `json:"year_published"`
	ISBN          string `json:"isbn"`
	CategoryID    int    `json:"category_id"`
	CategoryName  string `json:"category_name"`
	CoverImage    string `json:"cover_image"`
	Location      string `json:"location"`
	Status        string `json:"status"`
	Description   string `json:"description"`
	UploadedBy    int    `json:"uploaded_by"`
	UploaderName  string `json:"uploader_name"`
	UploaderEmail string `json:"uploader_email"`
	UploaderPhone string `json:"uploader_phone"`
}

type Category struct {
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
}

type Location struct {
	LocationID   int    `json:"location_id"`
	LocationName string `json:"location_name"`
	Address      string `json:"address"`
}

type BookLocation struct {
	BookLocationID int `json:"book_location_id"`
	BookID         int `json:"book_id"`
	LocationID     int `json:"location_id"`
	Stock          int `json:"stock"`
}

// GetAllBooks mengambil semua buku
func GetAllBooks() ([]Book, error) {
	rows, err := db.Query(`
		SELECT b.book_id, b.title, b.author, b.publisher, b.year_published, b.isbn, b.category_id, c.category_name, 
		       COALESCE(b.cover_image, ''), COALESCE(b.location, ''), COALESCE(b.status, 'pending'), COALESCE(b.description, ''),
		       COALESCE(b.uploaded_by, 0), COALESCE(b.uploader_name, ''), COALESCE(b.uploader_email, ''), COALESCE(b.uploader_phone, '')
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
		err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.Publisher, &book.YearPublished, &book.ISBN, &book.CategoryID, &book.CategoryName,
			&book.CoverImage, &book.Location, &book.Status, &book.Description, &book.UploadedBy, &book.UploaderName, &book.UploaderEmail, &book.UploaderPhone)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

// GetBookByID mengambil buku berdasarkan ID
func GetBookByID(id int) (*Book, error) {
	book := &Book{}
	err := db.QueryRow(`
		SELECT b.book_id, b.title, b.author, b.publisher, b.year_published, b.isbn, b.category_id, c.category_name,
		       COALESCE(b.cover_image, ''), COALESCE(b.location, ''), COALESCE(b.status, 'pending'), COALESCE(b.description, ''),
		       COALESCE(b.uploaded_by, 0), COALESCE(b.uploader_name, ''), COALESCE(b.uploader_email, ''), COALESCE(b.uploader_phone, '')
		FROM book b
		LEFT JOIN category c ON b.category_id = c.category_id
		WHERE b.book_id = ?
	`, id).Scan(&book.BookID, &book.Title, &book.Author, &book.Publisher, &book.YearPublished, &book.ISBN, &book.CategoryID, &book.CategoryName,
		&book.CoverImage, &book.Location, &book.Status, &book.Description, &book.UploadedBy, &book.UploaderName, &book.UploaderEmail, &book.UploaderPhone)
	if err != nil {
		return nil, err
	}
	return book, nil
}

// GetBooksByCategory mengambil buku berdasarkan kategori
func GetBooksByCategory(categoryID int) ([]Book, error) {
	rows, err := db.Query(`
		SELECT b.book_id, b.title, b.author, b.publisher, b.year_published, b.isbn, b.category_id, c.category_name,
		       COALESCE(b.cover_image, ''), COALESCE(b.location, ''), COALESCE(b.status, 'pending'), COALESCE(b.description, ''),
		       COALESCE(b.uploaded_by, 0), COALESCE(b.uploader_name, ''), COALESCE(b.uploader_email, ''), COALESCE(b.uploader_phone, '')
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
		err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.Publisher, &book.YearPublished, &book.ISBN, &book.CategoryID, &book.CategoryName,
			&book.CoverImage, &book.Location, &book.Status, &book.Description, &book.UploadedBy, &book.UploaderName, &book.UploaderEmail, &book.UploaderPhone)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

// SearchBooks mencari buku berdasarkan judul atau author
func SearchBooks(query string) ([]Book, error) {
	rows, err := db.Query(`
		SELECT b.book_id, b.title, b.author, b.publisher, b.year_published, b.isbn, b.category_id, c.category_name,
		       COALESCE(b.cover_image, ''), COALESCE(b.location, ''), COALESCE(b.status, 'pending'), COALESCE(b.description, ''),
		       COALESCE(b.uploaded_by, 0), COALESCE(b.uploader_name, ''), COALESCE(b.uploader_email, ''), COALESCE(b.uploader_phone, '')
		FROM book b
		LEFT JOIN category c ON b.category_id = c.category_id
		WHERE b.title LIKE ? OR b.author LIKE ?
	`, "%"+query+"%", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.Publisher, &book.YearPublished, &book.ISBN, &book.CategoryID, &book.CategoryName,
			&book.CoverImage, &book.Location, &book.Status, &book.Description, &book.UploadedBy, &book.UploaderName, &book.UploaderEmail, &book.UploaderPhone)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

// GetPopularBooks mengambil buku populer berdasarkan jumlah peminjaman
func GetPopularBooks(limit int) ([]Book, error) {
	rows, err := db.Query(`
		SELECT b.book_id, b.title, b.author, b.publisher, b.year_published, b.isbn, b.category_id, c.category_name,
		       COALESCE(b.cover_image, ''), COALESCE(b.location, ''), COALESCE(b.status, 'pending'), COALESCE(b.description, ''),
		       COALESCE(b.uploaded_by, 0), COALESCE(b.uploader_name, ''), COALESCE(b.uploader_email, ''), COALESCE(b.uploader_phone, '')
		FROM book b
		LEFT JOIN category c ON b.category_id = c.category_id
		LEFT JOIN borrow br ON b.book_id = br.book_id
		GROUP BY b.book_id
		ORDER BY COUNT(br.borrow_id) DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.Publisher, &book.YearPublished, &book.ISBN, &book.CategoryID, &book.CategoryName,
			&book.CoverImage, &book.Location, &book.Status, &book.Description, &book.UploadedBy, &book.UploaderName, &book.UploaderEmail, &book.UploaderPhone)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

// AddBook menambahkan buku baru (for admin use only, use uploadBook handler for user uploads)
func AddBook(title, author, publisher, isbn string, yearPublished int, categoryID int) (int, error) {
	result, err := db.Exec(`
		INSERT INTO book (title, author, publisher, isbn, year_published, category_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`, title, author, publisher, isbn, yearPublished, categoryID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// EditBook mengubah data buku (consolidated to handle both with and without description)
func EditBook(bookID int, title, author, publisher, isbn, description string, yearPublished int, categoryID int) error {
	fmt.Printf("[v0] EditBook called - BookID: %d, Description length: %d, Description: %s\n", bookID, len(description), description)

	sqlQuery := `UPDATE book SET title = ?, author = ?, publisher = ?, isbn = ?, description = ?, year_published = ?, category_id = ? WHERE book_id = ?`
	fmt.Printf("[v0] SQL Query: %s\n", sqlQuery)
	fmt.Printf("[v0] Parameters: title=%s, author=%s, publisher=%s, isbn=%s, description=%s, year=%d, catID=%d, bookID=%d\n",
		title, author, publisher, isbn, description, yearPublished, categoryID, bookID)

	result, err := db.Exec(sqlQuery,
		title, author, publisher, isbn, description, yearPublished, categoryID, bookID,
	)

	if err != nil {
		fmt.Printf("[v0] ERROR executing UPDATE: %v\n", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("[v0] ERROR getting rows affected: %v\n", err)
		return err
	}

	fmt.Printf("[v0] Rows affected: %d\n", rowsAffected)

	var savedDescription string
	err = db.QueryRow("SELECT description FROM book WHERE book_id = ?", bookID).Scan(&savedDescription)
	if err != nil {
		fmt.Printf("[v0] ERROR reading back description: %v\n", err)
	} else {
		fmt.Printf("[v0] Description saved in DB: %s\n", savedDescription)
	}

	return nil
}

// DeleteBook menghapus buku
func DeleteBook(bookID int) error {
	// First, return all active borrowed books
	_, err := db.Exec("UPDATE borrow SET status = 'returned', return_date = NOW() WHERE book_id = ? AND status = 'active'", bookID)
	if err != nil {
		return err
	}

	// Then delete the book
	_, err = db.Exec("DELETE FROM book WHERE book_id = ?", bookID)
	return err
}

// ValidateBookData memvalidasi data buku
func ValidateBookData(title, author, isbn string) bool {
	if title == "" || author == "" {
		return false
	}
	return true
}

// NotifyOwner mengirim notifikasi ke pemilik buku
func NotifyOwner(bookID int, userID int, message string) error {
	_, err := db.Exec(`
		INSERT INTO notification (user_id, book_id, message, created_at)
		VALUES (?, ?, ?, NOW())
	`, userID, bookID, message)
	return err
}

// GetAllLocations mengambil semua lokasi
func GetAllLocations() ([]Location, error) {
	rows, err := db.Query("SELECT location_id, location_name, address FROM location")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []Location
	for rows.Next() {
		var loc Location
		err := rows.Scan(&loc.LocationID, &loc.LocationName, &loc.Address)
		if err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}
	return locations, nil
}

// GetLocationByID mengambil lokasi berdasarkan ID
func GetLocationByID(id int) (*Location, error) {
	loc := &Location{}
	err := db.QueryRow("SELECT location_id, location_name, address FROM location WHERE location_id = ?", id).
		Scan(&loc.LocationID, &loc.LocationName, &loc.Address)
	if err != nil {
		return nil, err
	}
	return loc, nil
}

// UploadBook menambahkan buku baru dari pengguna
func UploadBook(title, author, publisher, isbn, coverImage, location, description string, yearPublished int, categoryID int, uploadedBy int, uploaderName, uploaderEmail, uploaderPhone string) (int, error) {
	result, err := db.Exec(`
		INSERT INTO book (title, author, publisher, isbn, year_published, category_id, cover_image, location, description, uploaded_by, uploader_name, uploader_email, uploader_phone)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, title, author, publisher, isbn, yearPublished, categoryID, coverImage, location, description, uploadedBy, uploaderName, uploaderEmail, uploaderPhone)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// GetBooksByUploader mengambil buku yang diunggah oleh pengguna tertentu
func GetBooksByUploader(userID int) ([]Book, error) {
	rows, err := db.Query(`
		SELECT b.book_id, b.title, b.author, b.publisher, b.year_published, b.isbn, b.category_id, c.category_name,
		       COALESCE(b.cover_image, ''), COALESCE(b.location, ''), COALESCE(b.status, 'pending'), COALESCE(b.description, ''),
		       COALESCE(b.uploaded_by, 0), COALESCE(b.uploader_name, ''), COALESCE(b.uploader_email, ''), COALESCE(b.uploader_phone, '')
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
		err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.Publisher, &book.YearPublished, &book.ISBN, &book.CategoryID, &book.CategoryName,
			&book.CoverImage, &book.Location, &book.Status, &book.Description, &book.UploadedBy, &book.UploaderName, &book.UploaderEmail, &book.UploaderPhone)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

// GetPendingBooks mengambil semua buku dengan status pending
func GetPendingBooks() ([]Book, error) {
	rows, err := db.Query(`
		SELECT b.book_id, b.title, b.author, b.publisher, b.year_published, b.isbn, b.category_id, c.category_name,
		       COALESCE(b.cover_image, ''), COALESCE(b.location, ''), COALESCE(b.status, 'pending'), COALESCE(b.description, ''),
		       COALESCE(b.uploaded_by, 0), COALESCE(b.uploader_name, ''), COALESCE(b.uploader_email, ''), COALESCE(b.uploader_phone, '')
		FROM book b
		LEFT JOIN category c ON b.category_id = c.category_id
		WHERE b.status = 'pending'
		ORDER BY b.book_id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.Publisher, &book.YearPublished, &book.ISBN, &book.CategoryID, &book.CategoryName,
			&book.CoverImage, &book.Location, &book.Status, &book.Description, &book.UploadedBy, &book.UploaderName, &book.UploaderEmail, &book.UploaderPhone)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

// GetAcceptedBooks mengambil semua buku dengan status accepted
func GetAcceptedBooks() ([]Book, error) {
	rows, err := db.Query(`
		SELECT b.book_id, b.title, b.author, b.publisher, b.year_published, b.isbn, b.category_id, c.category_name,
		       COALESCE(b.cover_image, ''), COALESCE(b.location, ''), COALESCE(b.status, 'pending'), COALESCE(b.description, ''),
		       COALESCE(b.uploaded_by, 0), COALESCE(b.uploader_name, ''), COALESCE(b.uploader_email, ''), COALESCE(b.uploader_phone, '')
		FROM book b
		LEFT JOIN category c ON b.category_id = c.category_id
		WHERE b.status = 'accepted'
		ORDER BY b.book_id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.Publisher, &book.YearPublished, &book.ISBN, &book.CategoryID, &book.CategoryName,
			&book.CoverImage, &book.Location, &book.Status, &book.Description, &book.UploadedBy, &book.UploaderName, &book.UploaderEmail, &book.UploaderPhone)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}
