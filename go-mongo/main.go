package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/smrkamran/go-mongo/controllers"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	r := httprouter.New()
	uc := controllers.NewUserController(getMongoClient())
	r.GET("/user/:id", uc.GetUser)
	r.POST("/user", uc.CreateUser)
	r.DELETE("/user/:id", uc.DeleteUser)

	http.ListenAndServe("localhost:9000", r)
}

func getMongoClient() *mongo.Client {
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	return client
}
