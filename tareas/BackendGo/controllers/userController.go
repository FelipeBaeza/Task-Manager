package controllers

import (
	"go-template/services"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// getUserCollection returns the MongoDB users collection
func getUserCollection(c *gin.Context) *mongo.Collection {
	return services.Client.Database("task_db").Collection("users")
}

// UserMe retrieves the authenticated user's profile information
func UserMe(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	// Convert string to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID format"})
		return
	}

	users := getUserCollection(c)
	var userData map[string]interface{}

	// Use ObjectID instead of string
	err = users.FindOne(c, bson.M{"_id": objectID}).Decode(&userData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to retrieve user data"})
		return
	}

	c.JSON(200, gin.H{"user": userData})
}

// GetAllUsers retrieves all users from the database
func GetAllUsers(c *gin.Context) {
	users := getUserCollection(c)
	cursor, err := users.Find(c, bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve users"})
		return
	}
	defer cursor.Close(c)

	var userList []map[string]interface{}
	if err = cursor.All(c, &userList); err != nil {
		c.JSON(500, gin.H{"error": "Failed to decode user data"})
		return
	}

	c.JSON(200, gin.H{"users": userList})
}
