package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/gookit/color.v1"
)

var collection *mongo.Collection
var ctx = context.TODO()

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("tasker").Collection("tasks")
}

func main() {
	app := &cli.App{
		Name:  "tasker",
		Usage: "A simple CLI Program to manage your tasks",
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Action: func(c *cli.Context) error {
					str := c.Args().First()
					if str == "" {
						return errors.New("cannot add an empty task")
					}

					task := &Task{
						ID:        primitive.NewObjectID(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						Text:      str,
						Completed: false,
					}
					return createTask(task)
				},
			},
			{
				Name:    "all",
				Aliases: []string{"l"},
				Usage:   "list all tasks",
				Action: func(c *cli.Context) error {
					tasks, err := getAll()
					if err != nil {
						if err == mongo.ErrNoDocuments {
							fmt.Print("Nothing to see here. \nRun `add 'task'` to add a task")
							return nil
						}
						return err
					}
					printTasks(tasks)
					return nil
				},
			},
			{
				Name:    "done",
				Aliases: []string{"d"},
				Usage:   "complete a task on the list",
				Action: func(c *cli.Context) error {
					text := c.Args().First()
					return completeTask(text)
				},
			},
			{
				Name:    "pending",
				Aliases: []string{"pd"},
				Usage:   "Pending tasks on the list",
				Action: func(c *cli.Context) error {
					tasks, err := getPending()
					if err != nil {
						if err == mongo.ErrNoDocuments {
							fmt.Print("Nothing to see here. \nRun `add 'task'` to add a task")
							return nil
						}
						return err
					}
					printTasks(tasks)
					return nil
				},
			},
			{
				Name:    "finish",
				Aliases: []string{"fn"},
				Usage:   "Finish tasks on the list",
				Action: func(c *cli.Context) error {
					tasks, err := getFinished()
					if err != nil {
						if err == mongo.ErrNoDocuments {
							fmt.Print("Nothing to see here. \nRun `add 'task'` to add a task")
							return nil
						}
						return err
					}
					printTasks(tasks)
					return nil
				},
			},
			{
				Name:  "rm",
				Usage: "Hard deletes a task on the list",
				Action: func(c *cli.Context) error {
					text := c.Args().First()
					err := deleteTask(text)
					if err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:  "del",
				Usage: "Soft deletes a task on the list",
				Action: func(c *cli.Context) error {
					text := c.Args().First()
					err := softDeleteTask(text)
					if err != nil {
						return err
					}
					return nil
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

type Task struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	DeletedAt time.Time          `bson:"deleted_at"`
	Text      string             `bson:"text"`
	Completed bool               `bson:"completed"`
	Deleted   bool               `bson:"deleted"`
}

func createTask(task *Task) error {
	_, err := collection.InsertOne(ctx, task)
	return err
}

func getAll() ([]*Task, error) {
	filter := bson.D{{}}
	return filterTasks(filter)
}

func filterTasks(filter interface{}) ([]*Task, error) {
	var tasks []*Task
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return tasks, err
	}
	for cur.Next(ctx) {
		var t Task
		err := cur.Decode(&t)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, &t)
	}
	if err := cur.Err(); err != nil {
		return tasks, err
	}

	cur.Close(ctx)
	if len(tasks) == 0 {
		return tasks, mongo.ErrNoDocuments
	}

	return tasks, nil
}

func printTasks(tasks []*Task) {
	for i, v := range tasks {
		if v.Completed {
			color.Green.Printf("%d: %s\n", i+1, v.Text)
		} else if v.Deleted {
			color.Red.Printf("%d: %s\n", i+1, v.Text)
		} else {
			color.Yellow.Printf("%d: %s\n", i+1, v.Text)
		}
	}
}

func completeTask(text string) error {
	filter := bson.D{primitive.E{Key: "text", Value: text}}

	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "completed", Value: true},
	}}}

	t := &Task{}
	return collection.FindOneAndUpdate(ctx, filter, update).Decode(t)
}

func getPending() ([]*Task, error) {
	filter := bson.D{
		primitive.E{Key: "completed", Value: false},
		primitive.E{Key: "deleted", Value: false},
	}

	return filterTasks(filter)
}

func getFinished() ([]*Task, error) {
	filter := bson.D{
		primitive.E{Key: "completed", Value: true},
	}

	return filterTasks(filter)
}
func deleteTask(text string) error {
	filter := bson.D{primitive.E{Key: "text", Value: text}}

	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("no tasks were deleted")
	}
	return nil
}

func softDeleteTask(text string) error {
	filter := bson.D{primitive.E{Key: "text", Value: text}}

	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "deleted", Value: true},
		primitive.E{Key: "deleted_at", Value: time.Now()},
	}}}

	t := &Task{}
	return collection.FindOneAndUpdate(ctx, filter, update).Decode(t)
}
