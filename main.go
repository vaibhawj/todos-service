package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
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
	Text string `json:"text" binding:"required"`
	Done bool   `json:"done"`
}

type TodoRequest struct {
	Text string `json:"text" binding:"required"`
	Done bool   `json:"done"`
}

var todos = []Todo{}

func getTodos(c *gin.Context) {
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
	todos = append(todos, newToDo)
	c.JSON(http.StatusCreated, newToDo)
}

func getTodo(c *gin.Context) {
	id := c.Param("id")
	for _, todo := range todos {
		if todo.Id == id {
			c.JSON(http.StatusOK, todo)
			return
		}
	}
	c.JSON(http.StatusNotFound, ErrorResponse{Error: "No todo found with specified id"})
}
