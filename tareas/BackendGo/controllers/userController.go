package controllers

import (
	"go-template/services"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getUserCollection(c *gin.Context) *mongo.Collection {
	return services.Client.Database("task_db").Collection("users")
}

func UserMe(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	// Convertir string a ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID format"})
		return
	}

	users := getUserCollection(c)
	var userData map[string]interface{}

	// Usar ObjectID en lugar de string
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
