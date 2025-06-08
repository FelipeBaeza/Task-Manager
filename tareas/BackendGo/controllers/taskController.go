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

	// Validate required fields
	if task.AssignedTo.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "AssignedTo is required"})
		return
	}

	if task.CreatedBy.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CreatedBy is required"})
		return
	}

	collection := getTasksCollection(c)
	if collection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to tasks collection"})
		return
	}

	// Set default values if not provided
	if task.Status == "" {
		task.Status = "pendiente"
	}
	if task.Priority == "" {
		task.Priority = "media"
	}

	// Validate status and priority - CORREGIDO AQUÍ
	if task.Status != "pendiente" && task.Status != "en_progreso" && task.Status != "completada" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Valid values: pendiente, en_progreso, completada"})
		return
	}

	if task.Priority != "baja" && task.Priority != "media" && task.Priority != "alta" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid priority. Valid values: baja, media, alta"})
		return
	}

	// Generate values for the task
	task.ID = primitive.NewObjectID()
	now := primitive.NewDateTimeFromTime(time.Now())
	task.CreatedAt = now
	task.UpdatedAt = now

	_, err := collection.InsertOne(context.TODO(), task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// GetTasks retrieves all tasks with optional filtering
func GetTasks(c *gin.Context) {
	collection := getTasksCollection(c)
	if collection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to tasks collection"})
		return
	}

	// Build filter based on query parameters
	filter := bson.M{}

	// Filter by status if provided
	if status := c.Query("status"); status != "" {
		// Validate status values according to schema
		if status != "pendiente" && status != "en_progreso" && status != "completada" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Valid values: pendiente, en_progreso, completada"})
			return
		}
		filter["status"] = status
	}

	// Filter by priority if provided
	if priority := c.Query("priority"); priority != "" {
		// Validate priority values
		if priority != "baja" && priority != "media" && priority != "alta" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid priority. Valid values: baja, media, alta"})
			return
		}
		filter["priority"] = priority
	}

	// Filter by assignedTo if provided
	if assignedTo := c.Query("assignedTo"); assignedTo != "" {
		assignedToID, err := primitive.ObjectIDFromHex(assignedTo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignedTo ID"})
			return
		}
		filter["assignedTo"] = assignedToID
	}

	// Filter by createdBy if provided
	if createdBy := c.Query("createdBy"); createdBy != "" {
		createdByID, err := primitive.ObjectIDFromHex(createdBy)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid createdBy ID"})
			return
		}
		filter["createdBy"] = createdByID
	}

	// Find tasks with filter
	cursor, err := collection.Find(context.TODO(), filter)
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

	// Return empty array if no tasks found instead of null
	if tasks == nil {
		tasks = []models.Task{}
	}

	c.JSON(http.StatusOK, tasks)
}

func GetTaskByID(c *gin.Context) {
	var task models.Task
	taskID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}
	collection := getTasksCollection(c)
	if collection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to tasks collection"})
		return
	}
	err = collection.FindOne(context.TODO(), bson.M{"_id": taskID}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// DeleteTask deletes a task by ID
func DeleteTask(c *gin.Context) {
	taskID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	collection := getTasksCollection(c)
	if collection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to tasks collection"})
		return
	}

	result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": taskID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// UpdateTaskStatus updates a task status (pendiente, en_progreso, completada)
func UpdateTaskStatus(c *gin.Context) {
	taskID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Get status from request body
	var requestBody struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
		return
	}

	// Validate status values
	if requestBody.Status != "pendiente" && requestBody.Status != "en_progreso" && requestBody.Status != "completada" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Valid values: pendiente, en_progreso, completada"})
		return
	}

	collection := getTasksCollection(c)
	if collection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to tasks collection"})
		return
	}

	// Update task status and updatedAt
	update := bson.M{
		"$set": bson.M{
			"status":    requestBody.Status,
			"updatedAt": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

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

	c.JSON(http.StatusOK, gin.H{
		"message": "Task status updated successfully",
		"status":  requestBody.Status,
	})
}
