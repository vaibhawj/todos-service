package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var todosCollection *mongo.Collection

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@localhost:27017"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	todosCollection = client.Database("todosdb").Collection("todos")

	router := gin.Default()
	router.GET("/todos", getTodos)
	router.POST("/todos", postTodo)
	router.GET("/todos/:id", getTodo)

	router.Run("localhost:8080")
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Todo struct {
	Id   string `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type TodoRequest struct {
	Text string `json:"text" binding:"required"`
	Done bool   `json:"done"`
}

func getTodos(c *gin.Context) {
	cur, err := todosCollection.Find(c.Request.Context(), bson.D{}, options.Find())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		fmt.Println(err)
		return
	}
	var todos []Todo
	cur.All(c.Request.Context(), &todos)
	c.JSON(http.StatusOK, todos)
}

func postTodo(c *gin.Context) {
	reqBody := TodoRequest{}
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	newToDo := Todo{Id: uuid.NewString(), Text: reqBody.Text, Done: reqBody.Done}

	_, err = todosCollection.InsertOne(c.Request.Context(), newToDo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusCreated, newToDo)
}

func getTodo(c *gin.Context) {
	id := c.Param("id")
	todo := Todo{}
	err := todosCollection.FindOne(c.Request.Context(), bson.M{"id": id}).Decode(&todo)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Todo with specified id not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, todo)
}
