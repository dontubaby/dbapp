package dbpkg

import (
	"context"
	"fmt"

	//"database/sql"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	Db *pgxpool.Pool
}

type Task struct {
	Id          int    `db:"id"`
	Opened      int64  `db:"opened"`
	Closed      int64  `db:"closed"`
	Author_id   int    `db:"author_id"`
	Assigned_id int    `db:"assigned_id"`
	Title       string `db:"title"`
	Content     string `db:"content"`
}

// Прописал SQL запрос на создание таблиц, согласно схеме БД.
// Добавил в начале запроса DROP, но это для универсальности запроса, детали задания не
// уточнялись. В принципе можно и без него.
// Так-же запрос можно легко распилить на создание отдельных таблиц, но это тоже не уточнялось
// в задании.

func (s *Storage) NewTable(ctx context.Context) error {
	_, err := s.Db.Exec(ctx, `DROP TABLE IF EXISTS tasks_labels,tasks,labels,users;

CREATE TABLE users(
  id SERIAL PRIMARY KEY,
  name TEXT NOT NUll
);

CREATE TABLE labels(
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL
);

CREATE TABLE tasks(
  id SERIAL PRIMARY KEY,
  opened BIGINT NOT NULL DEFAULT extract(epoch from now()),
  closed BIGINT DEFAULT 0,
  author_id INTEGER REFERENCES users(id) DEFAULT 0,
  assigned_id INTEGER REFERENCES users(id) DEFAULT 0,
  title TEXT,
  content TEXT  
);

CREATE TABLE tasks_labels(
  task_id INTEGER REFERENCES tasks(id),
  label_id INTEGER REFERENCES labels(id)
);`)
	if err != nil {
		log.Fatalf("Error!Cant create new tables:  %v\n", err)
		return err
	}
	return nil
}

func NewDb(ctx context.Context, connString string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		log.Fatalf("Cant create new instance of DB: %v\n", err)
	}
	s := Storage{
		Db: db,
	}
	return &s, nil
}

func (s *Storage) NewTask(ctx context.Context, t Task) error {
	_, err := s.Db.Exec(ctx, `INSERT INTO tasks(opened,closed,author_id,assigned_id,title,content) 
		VALUES ($1,$2,$3,$4,$5,$6);`,
		t.Opened, t.Closed, t.Author_id, t.Assigned_id, t.Title, t.Content)
	if err != nil {
		log.Fatalf("Error!Cant write new Task: %v\n", err)
		return err
	}
	return nil
}

func (s *Storage) GetAllTasks(ctx context.Context) ([]Task, error) {
	rows, err := s.Db.Query(ctx, `SELECT * FROM tasks;`)
	if err != nil {
		return nil, fmt.Errorf("Unable to query tasks: %w", err)
	}
	defer rows.Close()

	tasks := []Task{}

	for rows.Next() {
		task := Task{}
		err = rows.Scan(&task.Id, &task.Opened, &task.Closed, &task.Author_id,
			&task.Assigned_id, &task.Title, &task.Content)
		if err != nil {
			return nil, fmt.Errorf("Unable scan row: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *Storage) GetTasksByAuthor(ctx context.Context, author string) ([]Task, error) {
	rows, err := s.Db.Query(ctx, `SELECT tasks.id, tasks.opened,tasks.closed,tasks.author_id,tasks.assigned_id,
	tasks.title,tasks.content FROM tasks 
	JOIN users ON tasks.author_id=users.id WHERE users.name=$1;`, author)
	if err != nil {
		return nil, fmt.Errorf("Unable to query tasks: %w", err)
	}
	defer rows.Close()

	tasks := []Task{}

	for rows.Next() {
		task := Task{}
		err = rows.Scan(&task.Id, &task.Opened, &task.Closed, &task.Author_id,
			&task.Assigned_id, &task.Title, &task.Content)
		if err != nil {
			return nil, fmt.Errorf("Unable scan row: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *Storage) GetTasksByLabel(ctx context.Context, label string) ([]Task, error) {
	rows, err := s.Db.Query(ctx, `SELECT tasks.id, tasks.opened,tasks.author_id,tasks.assigned_id,
tasks.title,tasks.content
FROM tasks 
JOIN tasks_labels ON tasks_labels.task_id=tasks.id
JOIN labels 	  ON tasks_labels.label_id=labels.id
                  WHERE labels.name=$1;`, label)
	if err != nil {
		return nil, fmt.Errorf("Unable to query tasks: %w", err)
	}
	defer rows.Close()

	tasks := []Task{}

	for rows.Next() {
		task := Task{}
		err = rows.Scan(&task.Id, &task.Opened, &task.Closed, &task.Author_id, &task.Assigned_id, &task.Title, &task.Content)
		if err != nil {
			return nil, fmt.Errorf("Unable scan row: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *Storage) UpdateTaskById(ctx context.Context, id int, newTitle string, newContent string) error {
	_, err := s.Db.Exec(ctx, `UPDATE tasks SET title=$1,content=$2 WHERE tasks.id=$3;`,
		newTitle, newContent, id)
	if err != nil {
		log.Fatalf("Error!Cant write new Task: %v\n", err)
		return err
	}
	return nil
}

func (s *Storage) DeleteTaskById(ctx context.Context, id int) error {
	_, err := s.Db.Exec(ctx, `DELETE FROM tasks WHERE id=$1;`, id)
	if err != nil {
		log.Fatalf("Error!Cant write new Task: %v\n", err)
		return err
	}
	return nil
}


