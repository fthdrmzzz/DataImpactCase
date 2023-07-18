package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"dataimpact-backend/database"
	"dataimpact-backend/logger"
	"dataimpact-backend/models"
	"dataimpact-backend/utils"
)

// UserHandler handles the user-related API requests
type UserHandler struct {
	DB  *database.Database
	log *logrus.Logger
}

// NewUserHandler creates a new UserHandler instance with the provided dependencies
func NewUserHandler(db *database.Database) *UserHandler {
	return &UserHandler{
		DB:  db,
		log: logger.GetLogger(),
	}
}

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	File string `json:"file"`
}

// CreateUserResponse represents the response payload for creating a user
type CreateUserResponse struct {
	UserID string `json:"userId"`
}

// LoginRequest represents the request payload for user login
type LoginRequest struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

// LoginResponse represents the response payload for user login
type LoginResponse struct {
	Token string `json:"token"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
	File string `json:"file"`
}

// ListUsersResponse represents the response payload for listing users
type ListUsersResponse struct {
	Users []models.User `json:"users"`
}

// GetUserByIDResponse represents the response payload for getting a user by ID
type GetUserByIDResponse struct {
	User models.User `json:"user"`
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c *gin.Context) {
	h.log.Info("Creating a new user...")

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Deserialize user data
	var user models.User
	if err := json.Unmarshal([]byte(req.File), &user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the user exists
	existingUser := models.User{}
	err := h.DB.Database.Collection("users").FindOne(c.Request.Context(), bson.M{"_id": user.ID}).Decode(&existingUser)
	if err == nil {
		// User already exists, return conflict response
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	} else if err != mongo.ErrNoDocuments {
		// If the error is not ErrNoDocuments
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Encrypt password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
		return
	}
	user.Password = string(hashedPassword)

	// Insert user into the database
	result, err := h.DB.Database.Collection("users").InsertOne(c.Request.Context(), user)
	if err != nil {
		h.log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate user file
	fileName := user.ID.Hex() + ".txt"
	filePath := filepath.Join("users", fileName)
	if err := ioutil.WriteFile(filePath, []byte(user.Data), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user file"})
		h.log.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"userId": result.InsertedID.(primitive.ObjectID).Hex()})

}

// Login authenticates a user and generates a token
func (h *UserHandler) Login(c *gin.Context) {
	h.log.Info("Logging user in...")

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create primitive ObjectID from ID string
	primitiveObjectID, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Error parsing the ID"})
		return
	}

	// Find user by ID
	user := models.User{}
	err = h.DB.Database.Collection("users").FindOne(c.Request.Context(), bson.M{"_id": primitiveObjectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No user with id '" + req.ID + "' found"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong password"})
		}

		return
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate token
	token := utils.GenerateTokenWithDefaultLength()

	c.JSON(http.StatusOK, gin.H{"token": token})

}

// DeleteUser deletes a user by ID
func (h *UserHandler) DeleteUser(c *gin.Context) {
	h.log.Info("Deleting user...")

	userID := c.Param("id")
	// Create primitive ObjectID from ID string
	primitiveObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Error parsing the ID"})
		return
	}
	// Find user by ID
	user := models.User{}
	err = h.DB.Database.Collection("users").FindOne(c.Request.Context(), bson.M{"_id": primitiveObjectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		}
		return
	}

	// Delete user from the database
	_, err = h.DB.Database.Collection("users").DeleteOne(c.Request.Context(), bson.M{"_id": primitiveObjectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Delete user file
	filePath := filepath.Join("users", userID+".txt")
	err = os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})

}

// ListUsers lists all users
func (h *UserHandler) ListUsers(c *gin.Context) {
	h.log.Info("Listing users...")

	// Find all users
	cursor, err := h.DB.Database.Collection("users").Find(c.Request.Context(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list users"})
		return
	}
	defer cursor.Close(c.Request.Context())

	var users []models.User
	for cursor.Next(c.Request.Context()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user"})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, gin.H{"users": users})

}

// GetUserByID gets a user by ID
func (h *UserHandler) GetUserByID(c *gin.Context) {
	h.log.Info("Getting user by ID...")

	userID := c.Param("id")
	primitiveObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Error parsing the ID"})
		return
	}
	// Find user by ID
	user := models.User{}
	err = h.DB.Database.Collection("users").FindOne(c.Request.Context(), bson.M{"_id": primitiveObjectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})

}

// UpdateUser updates a user by ID
func (h *UserHandler) UpdateUser(c *gin.Context) {
	h.log.Info("Updating user...")

	userID := c.Param("id")
	primitiveObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Error parsing the ID"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Deserialize user data
	var updatedUser models.User
	if err := json.Unmarshal([]byte(req.File), &updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else if updatedUser.ID != primitive.NilObjectID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User Id can't be updated"})
		return
	}

	// Check if the user exists
	user := models.User{}
	err = h.DB.Database.Collection("users").FindOne(c.Request.Context(), bson.M{"_id": primitiveObjectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		}
		return
	}

	// Update user in the database
	_, err = h.DB.Database.Collection("users").ReplaceOne(c.Request.Context(), bson.M{"_id": primitiveObjectID}, updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Update user file
	filePath := filepath.Join("users", userID+".txt")
	if _, err := os.Stat(filePath); err == nil {
		err := ioutil.WriteFile(filePath, []byte(updatedUser.Data), 0644)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user file"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})

}
