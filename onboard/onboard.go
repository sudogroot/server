package onboard

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time" // vendor packages

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	cError "gitlab.com/joyread/ultimate/error"
	"gitlab.com/joyread/ultimate/models"
	"gitlab.com/joyread/ultimate/nextcloud"
	"golang.org/x/crypto/bcrypt" // custom packages
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateJWTToken(passwordHash string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{})
	tokenString, err := token.SignedString([]byte(passwordHash))
	return tokenString, err
}

func validateJWTToken(tokenString string, passwordHash string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(passwordHash), nil
	})

	return token.Valid, err
}

// ErrorResponse struct
type ErrorResponse struct {
	Error string `json:"error"`
}

// GetSignUp ...
func GetSignUp(c *gin.Context) {
	userPresent, ok := c.MustGet("userPresent").(bool)
	if !ok {
		fmt.Println("Middleware user error")
	}

	if userPresent {
		c.Redirect(http.StatusMovedPermanently, "/")
	} else {
		c.HTML(http.StatusOK, "signup.html", "")
	}
}

// SignUpRequest struct
type SignUpRequest struct {
	Username string `form:"username" binding:"required"`
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// PostSignUp ...
func PostSignUp(c *gin.Context) {
	var form SignUpRequest

	if err := c.Bind(&form); err == nil {
		// Generate password hash using bcrypt
		passwordHash, err := hashPassword(form.Password)
		cError.CheckError(err)

		// Generate JWT token using the hash password above
		token, err := generateJWTToken(passwordHash)
		cError.CheckError(err)

		db, ok := c.MustGet("db").(*sql.DB)
		if !ok {
			fmt.Println("Middleware db error")
		}

		signUpModel := models.SignUpModel{
			Username:     form.Username,
			Email:        form.Email,
			PasswordHash: passwordHash,
			Token:        token,
			IsAdmin:      1,
		}

		models.InsertAccount(db, signUpModel)

		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "joyread-token", Value: token, Expires: expiration}
		http.SetCookie(c.Writer, &cookie)

		c.Redirect(http.StatusMovedPermanently, "/")
	} else {
		errorResponse := &ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(http.StatusBadRequest, errorResponse)
	}
}

// // GetSMTP ...
// func GetSMTP(c *gin.Context) {
// 	db, ok := c.MustGet("db").(*sql.DB)
// 	if !ok {
// 		fmt.Println("Middleware db error")
// 	}

// 	if models.CheckSMTP(db) {
// 		c.Redirect(http.StatusMovedPermanently, "/storage")
// 	} else {
// 		c.HTML(http.StatusOK, "smtp.html", "")
// 	}
// }

// // SMTPRequest struct
// type SMTPRequest struct {
// 	SMTPHostname string `form:"hostname" binding:"required"`
// 	SMTPPort     string `form:"port" binding:"required"`
// 	SMTPUsername string `form:"username" binding:"required"`
// 	SMTPPassword string `form:"password" binding:"required"`
// }

// // PostSMTP ...
// func PostSMTP(c *gin.Context) {
// 	var form SMTPRequest

// 	if err := c.Bind(&form); err == nil {
// 		db, ok := c.MustGet("db").(*sql.DB)
// 		if !ok {
// 			fmt.Println("Middleware db error")
// 		}

// 		smtpPort, _ := strconv.Atoi(form.SMTPPort)

// 		smtpModel := models.SMTPModel{
// 			Hostname: form.SMTPHostname,
// 			Port:     smtpPort,
// 			Username: form.SMTPUsername,
// 			Password: form.SMTPPassword,
// 		}

// 		models.InsertSMTP(db, smtpModel)

// 		c.Redirect(http.StatusMovedPermanently, "/storage")
// 	} else {
// 		errorResponse := &ErrorResponse{
// 			Error: err.Error(),
// 		}
// 		c.JSON(http.StatusBadRequest, errorResponse)
// 	}
// }

// // TestEmailRequest struct
// type TestEmailRequest struct {
// 	SMTPHostname  string `json:"smtp_hostname" binding:"required"`
// 	SMTPPort      string `json:"smtp_port" binding:"required"`
// 	SMTPUsername  string `json:"smtp_username" binding:"required"`
// 	SMTPPassword  string `json:"smtp_password" binding:"required"`
// 	SMTPTestEmail string `json:"smtp_test_email" binding:"required"`
// }

// // TestEmailResponse struct
// type TestEmailResponse struct {
// 	IsEmailSent bool `json:"is_email_sent"`
// }

// // TestEmail ...
// func TestEmail(c *gin.Context) {
// 	var form TestEmailRequest

// 	if err := c.BindJSON(&form); err == nil {
// 		smtpPort, _ := strconv.Atoi(form.SMTPPort)
// 		emailSubject := "Joyread - Test email for your SMTP configuration"
// 		emailBody := "Congratulations mate!<br /><br />Your test email has been succesfully received, you could submit the Mail configuration form now.<br /><br />Cheers,<br/>Joyread"

// 		sendEmailRequest := email.SendEmailRequest{
// 			From:         form.SMTPUsername,
// 			To:           form.SMTPTestEmail,
// 			Subject:      emailSubject,
// 			Body:         emailBody,
// 			SMTPHostname: form.SMTPHostname,
// 			SMTPPort:     smtpPort,
// 			SMTPUsername: form.SMTPUsername,
// 			SMTPPassword: form.SMTPPassword,
// 		}

// 		isEmailSent := email.SendSyncEmail(sendEmailRequest)

// 		testEmailResponse := &TestEmailResponse{
// 			IsEmailSent: isEmailSent,
// 		}

// 		c.JSON(http.StatusOK, testEmailResponse)
// 	} else {
// 		errorResponse := &ErrorResponse{
// 			Error: err.Error(),
// 		}
// 		c.JSON(http.StatusBadRequest, errorResponse)
// 	}
// }

// GetStorage ...
func GetStorage(c *gin.Context) {
	c.HTML(http.StatusOK, "storage.html", "")
}

// NextcloudRequest struct
type NextcloudRequest struct {
	UserID                int    `form:"user_id" binding:"required"`
	NextcloudURL          string `form:"nextcloud_url" binding:"required"`
	NextcloudClientID     string `form:"nextcloud_client_id" binding:"required"`
	NextcloudClientSecret string `form:"nextcloud_client_secret" binding:"required"`
	NextcloudDirectory    string `form:"nextcloud_directory" binding:"required"`
	JoyreadURL            string `form:"joyread_url" binding:"required"`
}

// NextcloudResponse struct
type NextcloudResponse struct {
	Status  string `json:"status"`
	AuthURL string `json:"auth_url"`
}

// PostNextcloud ...
func PostNextcloud(c *gin.Context) {
	var form NextcloudRequest

	if err := c.Bind(&form); err == nil {
		db, ok := c.MustGet("db").(*sql.DB)
		if !ok {
			fmt.Println("Middleware db error")
		}

		// Redirect URI - https://myjoyread.com/nextcloud-auth/:user_id
		redirectURI := fmt.Sprintf("%s/nextcloud-auth/%d", form.JoyreadURL, form.UserID)

		nextcloudModel := models.NextcloudModel{
			UserID:       form.UserID,
			URL:          form.NextcloudURL,
			ClientID:     form.NextcloudClientID,
			ClientSecret: form.NextcloudClientSecret,
			Directory:    form.NextcloudDirectory,
			RedirectURI:  redirectURI,
		}
		models.InsertNextcloud(db, nextcloudModel)

		authURLRequest := nextcloud.AuthURLRequest{
			URL:         form.NextcloudURL,
			ClientID:    form.NextcloudClientID,
			RedirectURI: redirectURI,
		}
		authURL := nextcloud.GetAuthURL(authURLRequest)

		nextcloudResponse := &NextcloudResponse{
			Status:  "registered",
			AuthURL: authURL,
		}
		c.JSON(http.StatusMovedPermanently, nextcloudResponse)
	} else {
		errorResponse := &ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(http.StatusBadRequest, errorResponse)
	}
}

// NextcloudAuthCode ...
func NextcloudAuthCode(c *gin.Context) {
	// Get UserID from the URL
	userIDString := c.Param("user_id")
	var userID int
	if len(userIDString) > 0 {
		userID, _ = strconv.Atoi(userIDString)
	}

	// Get authorization code from the URL
	code := c.Query("code")

	db, ok := c.MustGet("db").(*sql.DB)
	if !ok {
		fmt.Println("Middleware db error")
	}

	selectNextcloudModel := models.SelectNextcloudModel{
		UserID: userID,
	}
	selectNextcloudResponse := models.SelectNextcloud(db, selectNextcloudModel)

	accessTokenRequest := nextcloud.AccessTokenRequest{
		URL:          selectNextcloudResponse.URL,
		ClientID:     selectNextcloudResponse.ClientID,
		ClientSecret: selectNextcloudResponse.ClientSecret,
		AuthCode:     code,
		RedirectURI:  selectNextcloudResponse.RedirectURI,
	}
	accessTokenResponse := nextcloud.GetAccessToken(accessTokenRequest)

	fmt.Println(accessTokenResponse)

	nextcloudTokenModel := models.NextcloudTokenModel{
		AccessToken:  accessTokenResponse.AccessToken,
		RefreshToken: accessTokenResponse.RefreshToken,
		UserID:       userID,
	}
	models.UpdateNextcloudToken(db, nextcloudTokenModel)

	c.Redirect(http.StatusMovedPermanently, "/")
}

// SignInRequest struct
type SignInRequest struct {
	UsernameOrEmail string `json:"usernameoremail" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

// SignInResponse struct
type SignInResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

// PostSignIn ...
func PostSignIn(c *gin.Context) {
	var form SignInRequest

	if err := c.BindJSON(&form); err == nil {
		db, ok := c.MustGet("db").(*sql.DB)
		if !ok {
			fmt.Println("Middleware db error")
		}

		selectPasswordHashAndJWTTokenModel := models.SelectPasswordHashAndJWTTokenModel{
			UsernameOrEmail: form.UsernameOrEmail,
		}
		selectPasswordHashAndJWTTokenResponse := models.SelectPasswordHashAndJWTToken(db, selectPasswordHashAndJWTTokenModel)

		if isPasswordValid := checkPasswordHash(form.Password, selectPasswordHashAndJWTTokenResponse.PasswordHash); isPasswordValid == true {
			isTokenValid, err := validateJWTToken(selectPasswordHashAndJWTTokenResponse.Token, selectPasswordHashAndJWTTokenResponse.PasswordHash)
			cError.CheckError(err)

			signInResponse := &SignInResponse{
				Status: "authorized",
				Token:  selectPasswordHashAndJWTTokenResponse.Token,
			}

			if isTokenValid == true {
				c.JSON(http.StatusMovedPermanently, signInResponse)
			} else {
				c.JSON(http.StatusMovedPermanently, gin.H{"status": "unauthorized"})
			}
		} else {
			c.JSON(http.StatusMovedPermanently, gin.H{"status": "unauthorized"})
		}
	} else {
		errorResponse := &ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(http.StatusBadRequest, errorResponse)
	}
}

// SignOut ...
func SignOut(c *gin.Context) {
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "joyread-token", Value: "", Expires: expiration}
	http.SetCookie(c.Writer, &cookie)

	c.Redirect(http.StatusMovedPermanently, "/signup")
}

// IsAdminPresent ...
func IsAdminPresent(c *gin.Context) {
	db, ok := c.MustGet("db").(*sql.DB)
	if !ok {
		fmt.Println("Middleware db error")
	}

	userID := models.SelectAdmin(db)
	isAdminPresent := false

	if userID > 0 {
		isAdminPresent = true
	}

	c.JSON(http.StatusOK, gin.H{"user_id": userID, "is_admin_present": isAdminPresent})
}

// IsSMTPPresent ...
// func IsSMTPPresent(c *gin.Context) {
// 	db, ok := c.MustGet("db").(*sql.DB)
// 	if !ok {
// 		fmt.Println("Middleware db error")
// 	}

// 	c.JSON(http.StatusOK, gin.H{"is_smtp_present": models.CheckSMTP(db)})
// }

// // IsStoragePresent ...
// func IsStoragePresent(c *gin.Context) {
// 	db, ok := c.MustGet("db").(*sql.DB)
// 	if !ok {
// 		fmt.Println("Middleware db error")
// 	}

// 	userID := models.SelectAdmin(db)
// 	isStoragePresent := false

// 	if storage := models.CheckStorage(db, userID); storage != "none" {
// 		isStoragePresent = true
// 	}

// 	c.JSON(http.StatusOK, gin.H{"user_id": userID, "is_storage_present": isStoragePresent})
// }

// // CheckOnboard ...
// func CheckOnboard(c *gin.Context) {
// 	db, ok := c.MustGet("db").(*sql.DB)
// 	if !ok {
// 		fmt.Println("Middleware db error")
// 	}
// 	userID := models.SelectAdmin(db)

// 	var currentProgress string

// 	if storage := models.CheckStorage(db, userID); storage != "none" {
// 		currentProgress = "onboarded"
// 	} else if models.CheckSMTP(db) {
// 		currentProgress = "smtp"
// 	} else if userID > 0 {
// 		currentProgress = "signup"
// 	}

// 	c.JSON(http.StatusOK, gin.H{"user_id": userID, "current_progress": currentProgress})
// }
