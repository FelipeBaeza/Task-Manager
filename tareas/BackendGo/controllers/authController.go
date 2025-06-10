package controllers

import (
	"net/http"
	"go-template/models"
	"go-template/middleware"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

)


// CreateUser handles user registration by creating a new user account
func CreateUser(c *gin.Context){
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if user.Name == "" || user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name or email or password is required"})
		return
	}

	users := getUserCollection(c)
	if users == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to user collection"})
		return
	}

	user.ID = primitive.NewObjectID()
	
	existingUser := models.User{}
	err := users.FindOne(c, bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}
	if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking existing user"})
		return
	}
	_, err = users.InsertOne(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": user})

}

// LoginUser handles user authentication by validating credentials
func LoginUser(c *gin.Context) {
	// Define the request structure DTO
	type LoginRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	var user LoginRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	
	// Validate email and password
	if user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	// Authenticate user
	collection := getUserCollection(c)
	if collection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to user collection"})
		return
	}

	var existingUser models.User
	err := collection.FindOne(c, bson.M{"email": user.Email}).Decode(&existingUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}
	
	// Generate JWT token
	token, err := middleware.GenerateToken(existingUser.ID.Hex())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
