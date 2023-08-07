package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mateenqazi/jwt-authenication/initializers"
	"github.com/mateenqazi/jwt-authenication/models"
	"golang.org/x/crypto/bcrypt"
)

type RequestBody struct {
	Email    string
	Password string
}

// SignupTags     godoc
// @Summary       Sign up Tags
// @Description   Save tags data in Db.
// @Produce       application/json
// @Tags          tags
// @Success       200 {object} map[string]interface{}
// @Router        /signup [post]
// @Param         email    formData  string   true  "User's email"
// @Param         password formData  string   true  "User's password"
func Signup(c *gin.Context) {
	// GET the email/pass from req body

	var body struct {
		Email    string `form:"email" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	// HASH the pass

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password  ",
		})
		return
	}

	// CREATE user

	user := models.User{Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user  ",
		})
		return
	}

	// response

	c.JSON(http.StatusOK, gin.H{})
}

// LoginTags      godoc
// @Summary       Login Tags
// @Description   login data.
// @Produce       application/json
// @Tags          tags
// @Success       200 {object} map[string]interface{}
// @Router        /login [post]
// @Param         email    formData  string   true  "User's email"
// @Param         password formData  string   true  "User's password"
func Login(c *gin.Context) {
	var body struct {
		Email    string `form:"email" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Look up requested user
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// compare sent in pass with saved user pass hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECERT_KEY")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	// response
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{})
}

// Validate godoc
// @Summary Validate user
// @Description Validate user's authentication.
// @Produce application/json
// @Tags authentication
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /validation [get]
func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}
