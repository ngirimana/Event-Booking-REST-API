package routes

import (
	"net/http"

	"example.com/rest-api/models"
	"example.com/rest-api/utils"
	"github.com/gin-gonic/gin"
)

// signUp godoc
// @Summary Create a new user
// @Description Registers a new user by saving their details in the database. It hashes the user's password before storing it.
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.User true "User data"
// @Success 201 {object} map[string]string "User created successfully"
// @Failure 400 {object} map[string]string "Could not parse request data"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/signup [post]
func signUp(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	err = user.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User created successfully!"})
}

// login godoc
// @Summary User Login
// @Description Authenticate a user and return a JWT token upon successful login.
// @Tags users
// @Accept  json
// @Produce  json
// @Param  user  body  models.User  true  "User credentials"
// @Success 200 {object} map[string]interface{} "message: Login successful, token: <JWT Token>"
// @Failure 400 {object} map[string]string "message: Could not parse request data"
// @Failure 401 {object} map[string]string "message: Invalid credentials"
// @Failure 500 {object} map[string]string "message: Could not generate token"
// @Router /users/login [post]
func login(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}

	err = user.Authenticate()
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials."})
		return
	}
	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token."})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Login successful!", "token": token})
}
