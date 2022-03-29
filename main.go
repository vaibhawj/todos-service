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

type server struct {
	db     *mongo.Database
	router *gin.Engine
}

func (s server) start() {
	router := s.router
	router.GET("/todos", s.getTodos)
	router.POST("/todos", s.postTodo)
	router.GET("/todos/:id", s.getTodo)

	router.Run("localhost:8080")
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@localhost:27017"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	server := server{db: client.Database("todosdb"), router: gin.Default()}

	server.start()
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

func (s server) getTodos(c *gin.Context) {
	cur, err := s.db.Collection("todos").Find(c.Request.Context(), bson.D{}, options.Find())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		fmt.Println(err)
		return
	}
	var todos []Todo
	cur.All(c.Request.Context(), &todos)
	c.JSON(http.StatusOK, todos)
}

func (s server) postTodo(c *gin.Context) {
	reqBody := TodoRequest{}
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	newToDo := Todo{Id: uuid.NewString(), Text: reqBody.Text, Done: reqBody.Done}

	_, err = s.db.Collection("todos").InsertOne(c.Request.Context(), newToDo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusCreated, newToDo)
}

func (s server) getTodo(c *gin.Context) {
	id := c.Param("id")
	todo := Todo{}
	err := s.db.Collection("todos").FindOne(c.Request.Context(), bson.M{"id": id}).Decode(&todo)
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
