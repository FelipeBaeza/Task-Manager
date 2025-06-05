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

// ------------------- functions to interact with MongoDB -------------------
func getTasksCollection(c *gin.Context) *mongo.Collection {
	if services.Client == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not available"})
		return nil
	}
	return services.Client.Database("task_db").Collection("tasks")
}

// ------------------- Task Controller Functions -------------------

// CreateTask creates a new task
func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if task.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	collection := getTasksCollection(c)
	if collection == nil {
		return
	}

	// Generate values for the task
	task.ID = primitive.NewObjectID()
	task.Completed = false
	task.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err := collection.InsertOne(context.TODO(), task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// GetTasks retrieves all tasks
func GetTasks(c *gin.Context) {
	collection := getTasksCollection(c)
	if collection == nil {
		return
	}

	// Find all tasks
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}
	defer cursor.Close(context.TODO())

	// Convert cursor to slice of tasks
	var tasks []models.Task
	if err := cursor.All(context.TODO(), &tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// CompletedTask marks a task as completed
func CompletedTask(c *gin.Context) {
	taskID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	collection := getTasksCollection(c)
	if collection == nil {
		return
	}

	// Update task to completed status
	update := bson.M{"$set": bson.M{"completed": true}}
	result, err := collection.UpdateOne(context.TODO(), bson.M{"_id": taskID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	// Check if task was found and updated
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}
