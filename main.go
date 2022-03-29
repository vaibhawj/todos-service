package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vaibhawj/todos-service/app"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@localhost:27017"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	server := app.NewServer(client.Database("todosdb"), gin.Default())

	server.Start()
}
