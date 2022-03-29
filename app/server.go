package app

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	m "github.com/vaibhawj/todos-service/model"
	r "github.com/vaibhawj/todos-service/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type server struct {
	db     *mongo.Database
	router *gin.Engine
	repo   r.Repository
}

func NewServer(db *mongo.Database, router *gin.Engine) server {
	server := server{db: db, router: router}
	server.repo = r.Repository{Collection: db.Collection("todos")}
	return server
}

func (s server) Start() {
	router := s.router
	router.GET("/todos", s.getTodos)
	router.POST("/todos", s.postTodo)
	router.GET("/todos/:id", s.getTodo)
	router.PUT("/todos/:id", s.updateTodo)
	router.DELETE("/todos/:id", s.deleteTodo)

	router.Run("localhost:8080")
}

func (s server) getTodos(c *gin.Context) {
	todos, err := s.repo.GetTodos(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, m.ErrorResponse{Error: err.Error()})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (s server) postTodo(c *gin.Context) {
	reqBody := m.TodoRequest{}
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, m.ErrorResponse{Error: err.Error()})
		return
	}
	newToDo := m.Todo{Id: uuid.NewString(), Text: reqBody.Text, Done: reqBody.Done}

	err = s.repo.CreateTodo(c.Request.Context(), newToDo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, m.ErrorResponse{Error: err.Error()})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusCreated, newToDo)
}

func (s server) getTodo(c *gin.Context) {
	id := c.Param("id")
	todo, err := s.repo.FindById(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			c.JSON(http.StatusNotFound, m.ErrorResponse{Error: "Todo with specified id not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, m.ErrorResponse{Error: err.Error()})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (s server) updateTodo(c *gin.Context) {
	id := c.Param("id")
	reqBody := m.TodoRequest{}
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, m.ErrorResponse{Error: err.Error()})
		return
	}

	newToDo := m.Todo{Id: id, Text: reqBody.Text, Done: reqBody.Done}
	err = s.repo.UpdateTodo(c.Request.Context(), newToDo)
	if err != nil {
		if err.Error() == "Todo with specified id not found" {
			c.JSON(http.StatusNotFound, m.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, m.ErrorResponse{Error: err.Error()})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, newToDo)
}

func (s server) deleteTodo(c *gin.Context) {
	id := c.Param("id")
	err := s.repo.DeleteTodo(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "Todo with specified id not found" {
			c.JSON(http.StatusNotFound, m.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, m.ErrorResponse{Error: err.Error()})
		fmt.Println(err)
		return
	}
	c.Status(http.StatusNoContent)
}
