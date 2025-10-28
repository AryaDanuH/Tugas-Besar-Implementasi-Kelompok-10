package main

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
		SELECT b.book_id, b.title, b.author, b.publisher, b.year_published, b.isbn, b.category_id, c.category_name
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

// GetBookByID mengambil buku berdasarkan ID
func GetBookByID(id int) (*Book, error) {
	book := &Book{}
	err := db.QueryRow(`
		SELECT b.book_id, b.title, b.author, b.publisher, b.year_published, b.isbn, b.category_id, c.category_name
		FROM book b
		LEFT JOIN category c ON b.category_id = c.category_id
		WHERE b.book_id = ?
	`, id).Scan(&book.BookID, &book.Title, &book.Author, &book.Publisher, &book.YearPublished, &book.ISBN, &book.CategoryID, &book.CategoryName)
	if err != nil {
		return nil, err
	}
	return book, nil
}

// GetBooksByCategory mengambil buku berdasarkan kategori
func GetBooksByCategory(categoryID int) ([]Book, error) {
	rows, err := db.Query(`
		SELECT b.book_id, b.title, b.author, b.publisher, b.year_published, b.isbn, b.category_id, c.category_name
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

// SearchBooks mencari buku berdasarkan judul atau author
func SearchBooks(query string) ([]Book, error) {
	rows, err := db.Query(`
		SELECT b.book_id, b.title, b.author, b.publisher, b.year_published, b.isbn, b.category_id, c.category_name
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
		err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.Publisher, &book.YearPublished, &book.ISBN, &book.CategoryID, &book.CategoryName)
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
		SELECT b.book_id, b.title, b.author, b.publisher, b.year_published, b.isbn, b.category_id, c.category_name
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
		err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.Publisher, &book.YearPublished, &book.ISBN, &book.CategoryID, &book.CategoryName)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

// AddBook menambahkan buku baru
func AddBook(title, author, publisher string, yearPublished int, isbn string, categoryID int) (int, error) {
	result, err := db.Exec(`
		INSERT INTO book (title, author, publisher, year_published, isbn, category_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`, title, author, publisher, yearPublished, isbn, categoryID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// EditBook mengubah data buku
func EditBook(bookID int, title, author, publisher string, yearPublished int, isbn string, categoryID int) error {
	_, err := db.Exec(`
		UPDATE book
		SET title = ?, author = ?, publisher = ?, year_published = ?, isbn = ?, category_id = ?
		WHERE book_id = ?
	`, title, author, publisher, yearPublished, isbn, categoryID, bookID)
	return err
}

// DeleteBook menghapus buku
func DeleteBook(bookID int) error {
	_, err := db.Exec("DELETE FROM book WHERE book_id = ?", bookID)
	return err
}

// ValidateBookData memvalidasi data buku
func ValidateBookData(title, author, isbn string) bool {
	if title == "" || author == "" || isbn == "" {
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
