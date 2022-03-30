package repository

import (
	"context"

	m "github.com/vaibhawj/todos-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ToDoRepository struct {
	collection *mongo.Collection
}

func NewTodoRepository(db *mongo.Database) ToDoRepository {
	return ToDoRepository{collection: db.Collection("todos")}
}

func (r ToDoRepository) GetTodos(c context.Context) ([]m.Todo, error) {
	cur, err := r.collection.Find(c, bson.D{}, options.Find())
	if err != nil {
		return nil, err
	}
	var todos []m.Todo
	cur.All(c, &todos)
	return todos, nil
}

func (r ToDoRepository) CreateTodo(c context.Context, todo m.Todo) error {
	_, err := r.collection.InsertOne(c, todo)
	if err != nil {
		return err
	}
	return nil
}

func (r ToDoRepository) FindById(c context.Context, id string) (m.Todo, error) {
	todo := m.Todo{}
	err := r.collection.FindOne(c, bson.M{"_id": id}).Decode(&todo)
	if err != nil {
		return m.Todo{}, err
	}
	return todo, nil
}

func (r ToDoRepository) UpdateTodo(c context.Context, todo m.Todo) error {
	res, err := r.collection.UpdateByID(c, todo.Id, bson.M{"$set": todo})
	if res.MatchedCount == 0 {
		return m.NotFoundError{}
	}
	return err
}

func (r ToDoRepository) DeleteTodo(c context.Context, id string) error {
	res, err := r.collection.DeleteOne(c, bson.M{"_id": id})
	if res.DeletedCount == 0 {
		return m.NotFoundError{}
	}
	return err
}
