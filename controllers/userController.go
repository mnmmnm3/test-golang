package controllers

import (
	"go-api/initializers"
	"go-api/models"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context, studentID uint, password string) {
	// hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})

		return
	}

	newUserID := "STD" + strconv.FormatUint(uint64(studentID), 10)

	// create the user
	user := models.User{UserID: newUserID, Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user, try to use different username",
		})

		return
	}

	// respond
	c.JSON(http.StatusOK, gin.H{"success": "sign up successfull"})
}

func Login(c *gin.Context) {
	// get the email/pass of req body
	var body struct {
		UserID   string `json:"user_id"`
		Password string `json:"password"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	// Look up requested user
	var user models.User
	initializers.DB.First(&user, "user_id = ?", body.UserID)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid username or password",
		})

		return
	}

	// compare sent in pass with saved user pass hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})

		return
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// sign and get the complete encoded token as a string using the secret
	tokenstring, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})

		return
	}

	// send it back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenstring, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully logged in",
	})
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"authenticaed_user": user,
	})
}
