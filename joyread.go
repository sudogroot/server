package joyread

import (
	// built-in packages
	"database/sql"
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv" // vendor packages

	"github.com/gin-gonic/gin" // custom packages
	"github.com/joyread/server/books"
	cError "github.com/joyread/server/error"
	"github.com/joyread/server/home"
	"github.com/joyread/server/middleware"
	"github.com/joyread/server/models"
	"github.com/joyread/server/onboard"
	"github.com/joyread/server/settings"
)

// StartServer handles the URL routes and starts the server
func StartServer() {
	// Gin initiate
	r := gin.Default()

	conf := settings.GetConf()

	// Serve static files
	r.Static("/assets", path.Join(conf.BaseValues.AssetPath, "assets"))

	// HTML rendering
	r.LoadHTMLGlob(path.Join(conf.BaseValues.AssetPath, "assets/templates/*"))

	// Create uploads/img directory
	os.MkdirAll(path.Join(conf.BaseValues.AssetPath, "uploads/img"), os.ModePerm)

	// Open postgres database
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", conf.BaseValues.DBValues.DBUsername, conf.BaseValues.DBValues.DBPassword, conf.BaseValues.DBValues.DBHostname, conf.BaseValues.DBValues.DBPort, conf.BaseValues.DBValues.DBName, conf.BaseValues.DBValues.DBSSLMode)
	db, err := sql.Open("postgres", connStr)
	cError.CheckError(err)
	defer db.Close()

	fmt.Println(runtime.GOOS)
	fmt.Println(runtime.GOARCH)

	if runtime.GOOS == "windows" {
		fmt.Println("Hello from Windows")
	}

	// models.CreateLegend(db)
	models.CreateAccount(db)
	models.CreateBooks(db)
	// models.CreateSMTP(db)
	// models.CreateNextcloud(db)

	r.Use(
		middleware.CORSMiddleware(),
		middleware.APIMiddleware(db),
		middleware.UserMiddleware(db),
	)

	// Gin handlers
	r.GET("/", home.Home)
	r.GET("/uploads/:bookName", home.ServeBook)
	r.GET("/cover/:coverName", home.ServeCover)
	r.GET("/signin", home.Home)
	r.GET("/send-file", home.SendFile)
	r.POST("/signin", onboard.PostSignIn)
	r.GET("/signup", onboard.GetSignUp)
	r.POST("/signup", onboard.PostSignUp)
	r.GET("/signout", onboard.SignOut)
	r.GET("/storage", onboard.GetStorage)
	r.POST("/nextcloud", onboard.PostNextcloud)
	r.GET("/nextcloud-auth/:user_id", onboard.NextcloudAuthCode)
	r.POST("/upload-books", books.UploadBooks)
	r.GET("/book/:bookName", books.GetBook)
	r.GET("/viewer/:bookName", books.Viewer)

	// Listen and serve
	port, err := strconv.Atoi(conf.BaseValues.ServerPort)
	if err != nil {
		fmt.Println("Invalid port specified")
		os.Exit(1)
	}
	r.Run(fmt.Sprintf(":%d", port))
}
