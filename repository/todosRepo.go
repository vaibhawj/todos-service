package repository

import (
	"context"

	m "github.com/vaibhawj/todos-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	Collection *mongo.Collection
}

func (r Repository) GetTodos(c context.Context) ([]m.Todo, error) {
	cur, err := r.Collection.Find(c, bson.D{}, options.Find())
	if err != nil {
		return nil, err
	}
	var todos []m.Todo
	cur.All(c, &todos)
	return todos, nil
}

func (r Repository) CreateTodo(c context.Context, todo m.Todo) error {
	_, err := r.Collection.InsertOne(c, todo)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) FindById(c context.Context, id string) (m.Todo, error) {
	todo := m.Todo{}
	err := r.Collection.FindOne(c, bson.M{"id": id}).Decode(&todo)
	if err != nil {
		return m.Todo{}, err
	}
	return todo, nil
}
