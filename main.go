package main

import (
	"context"
	"fmt"
	"os"

	//"database/sql"
	"Skillfactory/30-DB/pkg/storage"
	"log"
	//"os"
	//"github.com/jackc/pgx/v4"
	//"github.com/jackc/pgx/v4/pgxpool"
)

type User struct {
	id   int
	name string
}

func main() {

	task1 := dbpkg.Task{
		Opened:      0,
		Closed:      0,
		Author_id:   0,
		Assigned_id: 0,
		Title:       "Task1",
		Content:     "Task1 content",
	}

	ctx := context.Background()

	pwd := os.Getenv("DBPASS")
	log.Println(pwd)
	connString := "postgres://postgres:" + pwd + "@localhost:5432/gotest"

	pool, err := dbpkg.NewDb(ctx, connString)
	if err != nil {
		log.Fatalf("Cant connect to DB: %v", err)
	}
	defer pool.Db.Close()

	//Тест метода по добавлению задачи
	err = pool.NewTask(ctx, task1)
	if err != nil {
		log.Fatalf("Cant create new task: %v", err)
	}

	//Тест метода по получению списка задач
	tasks, err := pool.GetAllTasks(ctx)
	if err != nil {
		log.Fatalf("Cant read data from DB: %v", err)
	}
	fmt.Println(tasks)

	//Тест метода по получению задач по автору
	task, err := pool.GetTasksByAuthor(ctx, "Alex")
	if err != nil {
		log.Fatalf("Cant read data from DB: %v", err)
	}
	fmt.Println(task)

	//Тест метода по обновлению задачи
	err = pool.UpdateTaskById(ctx, 1, "NEW TASK", "NEW TASK CONTENT")

	//Тест метода по удалению задачи
	err = pool.DeleteTaskById(ctx, 10)
}
