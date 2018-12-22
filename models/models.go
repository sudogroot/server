/*
	Copyright (C) 2018 Nirmal Almara

    This file is part of Joyread.

    Joyread is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    Joyread is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
	along with Joyread.  If not, see <https://www.gnu.org/licenses/>.
*/

package models

import (
	"database/sql" // custom packages

	cError "gitlab.com/joyread/server/error"
)

// CreateLegend ...
func CreateLegend(db *sql.DB) {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS legend (server_url VARCHAR(255) NOT NULL, storage VARCHAR(255) NOT NULL DEFAULT 'local')")
	cError.CheckError(err)

	stmt.Exec()
}

// CreateAccount ...
func CreateAccount(db *sql.DB) {
	_, err := db.Query("CREATE TABLE IF NOT EXISTS account (id BIGSERIAL PRIMARY KEY, username VARCHAR(255) UNIQUE NOT NULL, email VARCHAR(255) UNIQUE NOT NULL, password_hash VARCHAR(255) NOT NULL, jwt_token VARCHAR(255) NOT NULL, is_admin BOOLEAN NOT NULL DEFAULT FALSE)")
	cError.CheckError(err)
}

// SignUpModel struct
type SignUpModel struct {
	Username     string
	Email        string
	PasswordHash string
	Token        string
	IsAdmin      int
}

// InsertAccount ...
func InsertAccount(db *sql.DB, signUpModel SignUpModel) {
	var lastInsertId int

	err := db.QueryRow("INSERT INTO account (username, email, password_hash, jwt_token, is_admin) VALUES ($1, $2, $3, $4, $5) returning id", signUpModel.Username, signUpModel.Email, signUpModel.PasswordHash, signUpModel.Token, signUpModel.IsAdmin).Scan(&lastInsertId)
	cError.CheckError(err)
}

// GetUserIDFromToken ...
func GetUserIDFromToken(db *sql.DB, token string) int {
	rows, err := db.Query("SELECT id FROM account WHERE jwt_token = $1", token)
	cError.CheckError(err)

	var userID int

	if rows.Next() {
		err := rows.Scan(&userID)
		cError.CheckError(err)
	}
	rows.Close()

	return userID
}

// SelectAdmin ...
func SelectAdmin(db *sql.DB) int {
	// Check for Admin in the user table
	rows, err := db.Query("SELECT id FROM account WHERE is_admin = $1", true)
	cError.CheckError(err)

	var userID int

	if rows.Next() {
		err := rows.Scan(&userID)
		cError.CheckError(err)
	}
	rows.Close()

	return userID
}

// SelectPasswordHashAndJWTTokenModel struct
type SelectPasswordHashAndJWTTokenModel struct {
	UsernameOrEmail string
}

// SelectPasswordHashAndJWTTokenResponse struct
type SelectPasswordHashAndJWTTokenResponse struct {
	PasswordHash string
	Token        string
}

// SelectPasswordHashAndJWTToken ...
func SelectPasswordHashAndJWTToken(db *sql.DB, selectPasswordHashAndJWTTokenModel SelectPasswordHashAndJWTTokenModel) *SelectPasswordHashAndJWTTokenResponse {
	// Search for username in the 'account' table with the given string
	rows, err := db.Query("SELECT password_hash, jwt_token FROM account WHERE username = $1", selectPasswordHashAndJWTTokenModel.UsernameOrEmail)
	cError.CheckError(err)

	var selectPasswordHashAndJWTTokenResponse SelectPasswordHashAndJWTTokenResponse

	if rows.Next() {
		err := rows.Scan(&selectPasswordHashAndJWTTokenResponse.PasswordHash, &selectPasswordHashAndJWTTokenResponse.Token)
		cError.CheckError(err)
	} else {
		// if username doesn't exist, search for email in the 'account' table with the given string
		rows, err := db.Query("SELECT password_hash, jwt_token FROM account WHERE email = $1", selectPasswordHashAndJWTTokenModel.UsernameOrEmail)
		cError.CheckError(err)

		if rows.Next() {
			err := rows.Scan(&selectPasswordHashAndJWTTokenResponse.PasswordHash, &selectPasswordHashAndJWTTokenResponse.Token)
			cError.CheckError(err)
		}
		rows.Close()
	}
	rows.Close()

	return &selectPasswordHashAndJWTTokenResponse
}

// CreateBooks ...
func CreateBooks(db *sql.DB) {
	_, err := db.Query("CREATE TABLE IF NOT EXISTS books (id BIGSERIAL PRIMARY KEY, title VARCHAR(255) NOT NULL, author VARCHAR(255) NOT NULL, pages INTEGER NOT NULL DEFAULT 0, file_path VARCHAR(255) NOT NULL, cover_path VARCHAR(255) NOT NULL, account_id INTEGER REFERENCES account(id))")
	cError.CheckError(err)
}

type BooksModel struct {
	Title     string
	Author    string
	Pages     string
	FilePath  string
	CoverPath string
	AccountID int
}

// InsertBooks ...
func InsertBooks(db *sql.DB, booksModel BooksModel) {
	var lastInsertId int

	err := db.QueryRow("INSERT INTO books (title, author, pages, file_path, cover_path, account_id) VALUES ($1, $2, $3, $4, $5, $6) returning id", booksModel.Title, booksModel.Author, booksModel.Pages, booksModel.FilePath, booksModel.CoverPath, booksModel.AccountID).Scan(&lastInsertId)
	cError.CheckError(err)
}

type BookInfo struct {
	Title     string
	FilePath  string
	CoverPath string
}

// GetBooks ...
func GetBooks(db *sql.DB, userID int) (books []BookInfo) {
	rows, err := db.Query("SELECT title, file_path, cover_path FROM books WHERE account_id = $1", userID)
	cError.CheckError(err)

	var title, filePath, coverPath string

	for rows.Next() {
		err := rows.Scan(&title, &filePath, &coverPath)
		cError.CheckError(err)

		books = append(books, BookInfo{
			Title:     title,
			FilePath:  filePath,
			CoverPath: coverPath,
		},
		)
	}
	rows.Close()

	return
}

// // CreateSMTP ...
// func CreateSMTP(db *sql.DB) {
// 	_, err := db.Query("CREATE TABLE IF NOT EXISTS smtp (hostname VARCHAR(255) NOT NULL, port INTEGER NOT NULL, username VARCHAR(255) NOT NULL, password VARCHAR(255) NOT NULL)")
// 	cError.CheckError(err)
// }

// // SMTPModel struct
// type SMTPModel struct {
// 	Hostname string
// 	Port     int
// 	Username string
// 	Password string
// }

// // InsertSMTP ...
// func InsertSMTP(db *sql.DB, smtpModel SMTPModel) {
// 	_, err := db.Query("INSERT INTO smtp (hostname, port, username, password) VALUES ($1, $2, $3, $4)", smtpModel.Hostname, smtpModel.Port, smtpModel.Username, smtpModel.Password)
// 	cError.CheckError(err)
// }

// // CheckSMTP ...
// func CheckSMTP(db *sql.DB) bool {

// 	// Check for Admin in the user table
// 	rows, err := db.Query("SELECT hostname FROM smtp")
// 	cError.CheckError(err)

// 	var isSMTPPresent = false

// 	if rows.Next() {
// 		isSMTPPresent = true
// 	}
// 	rows.Close()

// 	return isSMTPPresent
// }

// CreateNextcloud ...
func CreateNextcloud(db *sql.DB) {
	_, err := db.Query("CREATE TABLE IF NOT EXISTS nextcloud (id BIGSERIAL, user_id INTEGER REFERENCES account(id), url VARCHAR(255) NOT NULL, client_id VARCHAR(1200) NOT NULL, client_secret VARCHAR(1200) NOT NULL, directory VARCHAR(255) NOT NULL, redirect_uri VARCHAR(255) NOT NULL, access_token VARCHAR(255), refresh_token VARCHAR(255), PRIMARY KEY (id, user_id))")
	cError.CheckError(err)
}

// NextcloudModel struct
type NextcloudModel struct {
	UserID       int
	URL          string
	ClientID     string
	ClientSecret string
	Directory    string
	RedirectURI  string
}

// InsertNextcloud ...
func InsertNextcloud(db *sql.DB, nextcloudModel NextcloudModel) {
	_, err := db.Query("INSERT INTO nextcloud (user_id, url, client_id, client_secret, directory, redirect_uri) VALUES ($1, $2, $3, $4, $5, $6)", nextcloudModel.UserID, nextcloudModel.URL, nextcloudModel.ClientID, nextcloudModel.ClientSecret, nextcloudModel.Directory, nextcloudModel.RedirectURI)
	cError.CheckError(err)

	_, err = db.Query("UPDATE account SET storage=$1 WHERE id=$2", "nextcloud", nextcloudModel.UserID)
	cError.CheckError(err)
}

// SelectNextcloudModel struct
type SelectNextcloudModel struct {
	UserID int
}

// SelectNextcloudResponse struct
type SelectNextcloudResponse struct {
	URL          string
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// SelectNextcloud ...
func SelectNextcloud(db *sql.DB, selectNextcloudModel SelectNextcloudModel) *SelectNextcloudResponse {
	rows, err := db.Query("SELECT url, client_id, client_secret, redirect_uri FROM nextcloud WHERE user_id=$1", selectNextcloudModel.UserID)
	cError.CheckError(err)

	var selectNextcloudResponse SelectNextcloudResponse

	if rows.Next() {
		err := rows.Scan(&selectNextcloudResponse.URL, &selectNextcloudResponse.ClientID, &selectNextcloudResponse.ClientSecret, &selectNextcloudResponse.RedirectURI)
		cError.CheckError(err)
	}
	rows.Close()

	return &selectNextcloudResponse
}

// NextcloudTokenModel struct
type NextcloudTokenModel struct {
	AccessToken  string
	RefreshToken string
	UserID       int
}

// UpdateNextcloudToken ...
func UpdateNextcloudToken(db *sql.DB, nextcloudTokenModel NextcloudTokenModel) {
	_, err := db.Query("UPDATE nextcloud SET access_token=$1, refresh_token=$2 WHERE user_id=$3", nextcloudTokenModel.AccessToken, nextcloudTokenModel.RefreshToken, nextcloudTokenModel.UserID)
	cError.CheckError(err)
}

// CheckStorage ...
func CheckStorage(db *sql.DB, userID int) string {
	rows, err := db.Query("SELECT storage FROM account WHERE id=$1", userID)
	cError.CheckError(err)

	var storage string

	if rows.Next() {
		err := rows.Scan(&storage)
		cError.CheckError(err)
	}
	rows.Close()

	return storage
}

// CheckNextcloudToken ...
func CheckNextcloudToken(db *sql.DB, userID int) string {
	rows, err := db.Query("SELECT access_token FROM nextcloud WHERE user_id=$1", userID)
	cError.CheckError(err)

	var accessToken string

	if rows.Next() {
		err := rows.Scan(&accessToken)
		cError.CheckError(err)
	}
	rows.Close()

	return accessToken
}
