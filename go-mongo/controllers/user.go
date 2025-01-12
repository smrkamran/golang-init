package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/smrkamran/go-mongo/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserController struct {
	mongoClient *mongo.Client
}

func NewUserController(c *mongo.Client) *UserController {
	return &UserController{mongoClient: c}
}

func (uc UserController) getCollection() *mongo.Collection {
	return uc.mongoClient.Database("mongo-golang").Collection("users")
}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	// Validate ObjectID
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusNotFound)
		return
	}

	// Fetch user from database
	u := models.User{}
	collection := uc.getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&u)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Respond with JSON
	uj, err := json.Marshal(u)
	if err != nil {
		http.Error(w, "Failed to marshal user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", uj)
}

// CreateUser
func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Decode request body
	u := models.User{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Assign a new ObjectID
	u.Id = bson.NewObjectID()

	// Insert into database
	collection := uc.getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, u)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Respond with created user
	uj, err := json.Marshal(u)
	if err != nil {
		http.Error(w, "Failed to marshal user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n", uj)
}

// DeleteUser
func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	// Validate ObjectID
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusNotFound)
		return
	}

	// Delete user from database
	collection := uc.getCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Deleted user %s\n", id)
}
