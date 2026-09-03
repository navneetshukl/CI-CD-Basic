package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setUpRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/health", healthCheck)
	api := r.Group("/api")
	{
		api.GET("/users", getUsers)
		api.POST("/users", createUser)
		api.GET("/users/:id", getUserByID)
	}

	return r
}

func resetUsers() {
	mu.Lock()
	users = map[int]*User{}
	nextID = 1
	mu.Unlock()
}

func TestHealthCheck(t *testing.T) {
	router := setUpRouter()
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got %s", response["status"])
	}
}

func TestCreateUser(t *testing.T) {
	router := setUpRouter()
	resetUsers()

	user := User{Name: "John", Age: 30}
	jsonValue, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	createdUser := response["user"].(map[string]interface{})
	if createdUser["name"] != "John" {
		t.Errorf("Expected name 'John', got %s", createdUser["name"])
	}
	if int(createdUser["age"].(float64)) != 30 {
		t.Errorf("Expected age 30, got %d", int(createdUser["age"].(float64)))
	}
}

func TestGetUsers(t *testing.T) {
	router := setUpRouter()
	resetUsers()

	req, _ := http.NewRequest("GET", "/api/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if int(response["count"].(float64)) != 0 {
		t.Errorf("Expected count 0, got %d", int(response["count"].(float64)))
	}
}

func TestGetUserByID(t *testing.T) {
	router := setUpRouter()
	resetUsers()

	user := User{Name: "Alice", Age: 25}
	jsonValue, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	req2, _ := http.NewRequest("GET", "/api/users/1", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w2.Code)
	}
}

func TestGetUserByInvalidID(t *testing.T) {
	router := setUpRouter()
	resetUsers()

	req, _ := http.NewRequest("GET", "/api/users/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetUserNotFound(t *testing.T) {
	router := setUpRouter()
	resetUsers()

	req, _ := http.NewRequest("GET", "/api/users/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

// TestGetUserName will fail - intentionally checking for wrong name
func TestGetUserName(t *testing.T) {
	router := setUpRouter()
	resetUsers()

	user := User{Name: "Alice", Age: 25}
	jsonValue, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	createdUser := response["user"].(map[string]interface{})
	// INTENTIONALLY WRONG: Expecting "Bob" but getting "Alice"
	if createdUser["name"] != "Bob" {
		t.Errorf("Expected name 'Bob', got '%s'", createdUser["name"])
	}
}
