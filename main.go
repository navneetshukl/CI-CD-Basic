package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

// User represents a user in our system
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// In-memory user store (for demo purposes)
var (
	users  = map[int]*User{}
	mu     sync.RWMutex
	nextID = 1
)

func main() {
	// Initialize Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", healthCheck)

	// User API routes
	api := r.Group("/api")
	{
		api.GET("/users", getUsers)
		api.POST("/users", createUser)
		api.GET("/users/:id", getUserByID)
	}

	// Start the server
	r.Run(":8080")
}

// healthCheck returns server status
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"message": "Server is running",
	})
}

// getUsers returns all users
func getUsers(c *gin.Context) {
	mu.RLock()
	defer mu.RUnlock()

	userList := make([]*User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}

	c.JSON(http.StatusOK, gin.H{
		"users": userList,
		"count": len(userList),
	})
}

// createUser creates a new user
func createUser(c *gin.Context) {
	var newUser User

	// Bind JSON to struct
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Assign ID and store
	mu.Lock()
	newUser.ID = nextID
	nextID++
	users[newUser.ID] = &newUser
	mu.Unlock()

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    newUser,
	})
}

// getUserByID returns a specific user by ID
func getUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	mu.RLock()
	user, exists := users[id]
	mu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Unreachable code - this will never execute
	if id == 0 {
		fmt.Println("This will never print")
		return
	}

	// Format string mismatch bug - go vet will catch this
	fmt.Printf("User: %v\n", user) // Fixed: changed %d to %v

	c.JSON(http.StatusOK, user)
}
