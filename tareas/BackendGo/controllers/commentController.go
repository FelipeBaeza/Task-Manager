package controllers

import (
	"context"
	"go-template/models"
	"go-template/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getCommentsCollection(c *gin.Context) *mongo.Collection {
	return services.Client.Database("task_db").Collection("comments")
}

// CreateComment creates a new comment for a task
func CreateComment(c *gin.Context) {
	taskID, err := primitive.ObjectIDFromHex(c.Param("id")) // Cambiado de "taskId" a "id"
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if comment.Text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment text is required"})
		return
	}

	// Get authenticated user ID
	userID := c.MustGet("userID").(string)
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Verify task exists and user has access
	tasksCollection := getTasksCollection(c)
	var task models.Task
	err = tasksCollection.FindOne(context.TODO(), bson.M{
		"_id":       taskID,
		"createdBy": userObjectID,
	}).Decode(&task)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found or no permission"})
		return
	}

	// Set comment values
	comment.ID = primitive.NewObjectID()
	comment.TaskID = taskID
	comment.AuthorID = userObjectID
	comment.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	collection := getCommentsCollection(c)
	_, err = collection.InsertOne(context.TODO(), comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// GetCommentsByTask gets all comments for a specific task
func GetCommentsByTask(c *gin.Context) {
	taskID, err := primitive.ObjectIDFromHex(c.Param("id")) // Cambiado de "taskId" a "id"
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Get authenticated user ID
	userID := c.MustGet("userID").(string)
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Verify task exists and user has access
	tasksCollection := getTasksCollection(c)
	var task models.Task
	err = tasksCollection.FindOne(context.TODO(), bson.M{
		"_id":       taskID,
		"createdBy": userObjectID,
	}).Decode(&task)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found or no permission"})
		return
	}

	collection := getCommentsCollection(c)
	cursor, err := collection.Find(context.TODO(), bson.M{"taskId": taskID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	defer cursor.Close(context.TODO())

	var comments []models.Comment
	if err := cursor.All(context.TODO(), &comments); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode comments"})
		return
	}

	if comments == nil {
		comments = []models.Comment{}
	}

	c.JSON(http.StatusOK, comments)
}

// DeleteComment deletes a comment
func DeleteComment(c *gin.Context) {
	commentID, err := primitive.ObjectIDFromHex(c.Param("commentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	// Get authenticated user ID
	userID := c.MustGet("userID").(string)
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	collection := getCommentsCollection(c)

	// Only allow deletion of own comments
	filter := bson.M{
		"_id":      commentID,
		"authorId": userObjectID,
	}

	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found or no permission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
