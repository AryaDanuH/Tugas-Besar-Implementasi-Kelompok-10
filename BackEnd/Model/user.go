package main

import (
	"database/sql"
	"time"
)

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

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

// GetUserByEmail mengambil user berdasarkan email
func GetUserByEmail(email string) (*User, error) {
	user := &User{}
	err := db.QueryRow("SELECT user_id, name, email, password, phone, address, role, COALESCE(profile_image, '') FROM user WHERE email = ?", email).
		Scan(&user.UserID, &user.Name, &user.Email, &user.Password, &user.Phone, &user.Address, &user.Role, &user.ProfileImage)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// GetUserByID mengambil user berdasarkan ID
func GetUserByID(id int) (*User, error) {
	user := &User{}
	err := db.QueryRow("SELECT user_id, name, email, password, phone, address, role, COALESCE(profile_image, '') FROM user WHERE user_id = ?", id).
		Scan(&user.UserID, &user.Name, &user.Email, &user.Password, &user.Phone, &user.Address, &user.Role, &user.ProfileImage)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// CreateUser membuat user baru
func CreateUser(name, email, password, phone, address string) (int, error) {
	result, err := db.Exec("INSERT INTO user (name, email, password, phone, address, role) VALUES (?, ?, ?, ?, ?, ?)",
		name, email, password, phone, address, "member")
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

// UpdateUser mengupdate data user
func UpdateUser(id int, name, phone, address string) error {
	_, err := db.Exec("UPDATE user SET name = ?, phone = ?, address = ? WHERE user_id = ?",
		name, phone, address, id)
	return err
}

func UpdateUserProfileImage(id int, profileImage string) error {
	_, err := db.Exec("UPDATE user SET profile_image = ? WHERE user_id = ?",
		profileImage, id)
	return err
}

// UpdateUserPassword mengupdate password user
func UpdateUserPassword(id int, newPassword string) error {
	_, err := db.Exec("UPDATE user SET password = ? WHERE user_id = ?",
		newPassword, id)
	return err
}

func ChangeUsername(id int, newName string) error {
	_, err := db.Exec("UPDATE user SET name = ? WHERE user_id = ?", newName, id)
	return err
}

func ForgotPassword(email string) (string, error) {
	user, err := GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", sql.ErrNoRows
	}

	// Generate reset token
	token := generateResetToken()
	err = SaveResetToken(user.UserID, token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func VerifyUserData(name, email, phone string) bool {
	if name == "" || email == "" || phone == "" {
		return false
	}
	return true
}

// SaveResetToken menyimpan token reset password
func SaveResetToken(userID int, token string) error {
	_, err := db.Exec("UPDATE user SET reset_token = ?, reset_token_expiry = DATE_ADD(NOW(), INTERVAL 1 HOUR) WHERE user_id = ?",
		token, userID)
	return err
}

// VerifyResetToken memverifikasi token reset password
func VerifyResetToken(userID int, token string) (bool, error) {
	var storedToken string
	var expiry time.Time
	err := db.QueryRow("SELECT reset_token, reset_token_expiry FROM user WHERE user_id = ?", userID).
		Scan(&storedToken, &expiry)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	if storedToken == token && time.Now().Before(expiry) {
		return true, nil
	}
	return false, nil
}

// DeleteResetToken menghapus token reset password
func DeleteResetToken(userID int) error {
	_, err := db.Exec("UPDATE user SET reset_token = NULL, reset_token_expiry = NULL WHERE user_id = ?", userID)
	return err
}

func generateResetToken() string {
	return time.Now().Format("20060102150405") + "-" + string(rune(time.Now().UnixNano()))
}
