package main

import (
	"database/sql"
	"time"
)

type Borrow struct {
	BorrowID       int        `json:"borrow_id"`
	UserID         int        `json:"user_id"`
	BookLocationID int        `json:"book_location_id"`
	BorrowDate     time.Time  `json:"borrow_date"`
	ReturnDate     *time.Time `json:"return_date"`
	Status         string     `json:"status"`
	DeliveryType   string     `json:"delivery_type"`
	BookTitle      string     `json:"book_title"`
	LocationName   string     `json:"location_name"`
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
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateBorrow membuat peminjaman baru
func CreateBorrow(userID, bookLocationID int, deliveryType string) (int, error) {
	result, err := db.Exec(`
		INSERT INTO borrow (user_id, book_location_id, borrow_date, status, delivery_type)
		VALUES (?, ?, NOW(), 'pending', ?)
	`, userID, bookLocationID, deliveryType)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

// GetBorrowByID mengambil peminjaman berdasarkan ID
func GetBorrowByID(id int) (*Borrow, error) {
	borrow := &Borrow{}
	err := db.QueryRow(`
		SELECT br.borrow_id, br.user_id, br.book_location_id, br.borrow_date, br.return_date, br.status, br.delivery_type, b.title, l.location_name
		FROM borrow br
		JOIN book_location bl ON br.book_location_id = bl.book_location_id
		JOIN book b ON bl.book_id = b.book_id
		JOIN location l ON bl.location_id = l.location_id
		WHERE br.borrow_id = ?
	`, id).Scan(&borrow.BorrowID, &borrow.UserID, &borrow.BookLocationID, &borrow.BorrowDate, &borrow.ReturnDate, &borrow.Status, &borrow.DeliveryType, &borrow.BookTitle, &borrow.LocationName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return borrow, nil
}

// GetUserBorrows mengambil peminjaman user
func GetUserBorrows(userID int) ([]Borrow, error) {
	rows, err := db.Query(`
		SELECT br.borrow_id, br.user_id, br.book_location_id, br.borrow_date, br.return_date, br.status, br.delivery_type, b.title, l.location_name
		FROM borrow br
		JOIN book_location bl ON br.book_location_id = bl.book_location_id
		JOIN book b ON bl.book_id = b.book_id
		JOIN location l ON bl.location_id = l.location_id
		WHERE br.user_id = ?
		ORDER BY br.borrow_date DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var borrows []Borrow
	for rows.Next() {
		var borrow Borrow
		err := rows.Scan(&borrow.BorrowID, &borrow.UserID, &borrow.BookLocationID, &borrow.BorrowDate, &borrow.ReturnDate, &borrow.Status, &borrow.DeliveryType, &borrow.BookTitle, &borrow.LocationName)
		if err != nil {
			return nil, err
		}
		borrows = append(borrows, borrow)
	}
	return borrows, nil
}

// ApproveBorrow menyetujui peminjaman
func ApproveBorrow(borrowID int) error {
	_, err := db.Exec("UPDATE borrow SET status = 'approved' WHERE borrow_id = ?", borrowID)
	return err
}

// RejectBorrow menolak peminjaman
func RejectBorrow(borrowID int) error {
	_, err := db.Exec("UPDATE borrow SET status = 'rejected' WHERE borrow_id = ?", borrowID)
	return err
}

// ReturnBook mengembalikan buku
func ReturnBook(borrowID int) error {
	_, err := db.Exec("UPDATE borrow SET return_date = NOW(), status = 'returned' WHERE borrow_id = ?", borrowID)
	return err
}

// UpdateBookStatus mengupdate status buku
func UpdateBookStatus(bookID int, status string) error {
	_, err := db.Exec("UPDATE book SET status = ? WHERE book_id = ?", status, bookID)
	return err
}

// CreateReview membuat review buku
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

// GetBookReviews mengambil review buku
func GetBookReviews(bookID int) ([]Review, error) {
	rows, err := db.Query(`
		SELECT r.review_id, r.book_id, r.user_id, r.rating, r.comment, u.name, r.created_at
		FROM review r
		JOIN user u ON r.user_id = u.user_id
		WHERE r.book_id = ?
		ORDER BY r.created_at DESC
	`, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []Review
	for rows.Next() {
		var review Review
		err := rows.Scan(&review.ReviewID, &review.BookID, &review.UserID, &review.Rating, &review.Comment, &review.UserName, &review.CreatedAt)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}
