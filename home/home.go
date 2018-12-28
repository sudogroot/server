package home

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	cError "gitlab.com/joyread/server/error"
	"gitlab.com/joyread/server/models"
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
