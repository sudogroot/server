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

package home

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	cError "gitlab.com/joyread/ultimate/error"
	"gitlab.com/joyread/ultimate/models"
)

type BookInfo struct {
	Title     string
	FilePath  string
	CoverPath string
}

// Home ...
func Home(c *gin.Context) {
	token, _ := c.Cookie("joyread-token")

	db, ok := c.MustGet("db").(*sql.DB)
	if !ok {
		fmt.Println("Middleware db error")
	}

	var userID int
	userID = models.GetUserIDFromToken(db, token)

	if userID != 0 {
		books := models.GetBooks(db, userID)

		c.HTML(http.StatusOK, "index.html", gin.H{
			"books": books,
		})
	} else {
		c.Redirect(http.StatusMovedPermanently, "/signup")
	}
}

// SendFile ...
func SendFile(c *gin.Context) {
	err := downloadFile()
	cError.CheckError(err)

	c.HTML(http.StatusOK, "pdf-wrapper.html", gin.H{
		"bookName": "NextcloudManual.pdf",
	})
}

// downloadFile ...
func downloadFile() error {
	req, err := http.NewRequest("GET", "http://139.59.68.38/remote.php/webdav/NextcloudManual.pdf", nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth("mysticmode", "XzBxa-F8e56-eGRwo-WFijE-xws7m")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("bad status: %s", resp.Status)
	}

	// Create the file
	out, err := os.Create("uploads/NextcloudManual.pdf")
	if err != nil {
		return err
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// ServeBook ...
func ServeBook(c *gin.Context) {
	userPresent, ok := c.MustGet("userPresent").(bool)
	if !ok {
		fmt.Println("Middleware user error")
	}

	if userPresent {
		bookName := c.Param("bookName")
		bookPath := fmt.Sprintf("./uploads/%s", bookName)
		http.ServeFile(c.Writer, c.Request, bookPath)
	} else {
		c.Redirect(http.StatusMovedPermanently, "/signup")
	}
}

// ServeCover ...
func ServeCover(c *gin.Context) {
	userPresent, ok := c.MustGet("userPresent").(bool)
	if !ok {
		fmt.Println("Middleware user error")
	}

	if userPresent {
		coverName := c.Param("coverName")
		coverPath := fmt.Sprintf("./uploads/img/%s", coverName)
		http.ServeFile(c.Writer, c.Request, coverPath)
	} else {
		c.Redirect(http.StatusMovedPermanently, "/signup")
	}
}
