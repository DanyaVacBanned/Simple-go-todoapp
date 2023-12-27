package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
	"todoapp/internal/storage"
)

type Storage struct {
	dbname *sql.DB
}

type Task struct {
	TaskName        string
	TaskDescription string
	TaskDone        int
	TaskCreated     string
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	query, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS task(
	    id INTEGER PRIMARY KEY,
	    task_name VARCHAR(150) NOT NULL,
	    task_description TEXT NOT NULL,
	    task_done INTEGER DEFAULT 0,
	    task_created TEXT NOT NULL
	);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	_, err = query.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	return &Storage{dbname: db}, nil
}

func (s *Storage) CreateTask(
	taskName string, taskDescription string,
) (int64, error) {
	const op = "storage.sqlite.CreateTask"
	if taskName == "" || taskDescription == "" {
		return 0, fmt.Errorf("%s: required fields are not specified", op)
	}
	query, err := s.dbname.Prepare(`
	INSERT INTO task(task_name, task_description, task_created)
	VALUES (?, ?, ?);
`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	res, err := query.Exec(taskName, taskDescription, time.Now().String())
	id, _ := res.LastInsertId()
	return id, nil
}

func (s *Storage) DeleteTask(taskId int64) error {
	const op = "storage.sqlite.DeleteTask"
	query, err := s.dbname.Prepare(`
	DELETE FROM task WHERE id = ?;
`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = query.Exec(taskId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storage.ErrTaskNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) UpdateTask(taskId int64, taskName string, taskDescription string, taskDone int) (int64, error) {
	const op = "storage.sqlite.UpdateTask"
	if taskDone > 1 && taskDone < 0 {
		return 0, fmt.Errorf("taskDone must be 1 or 0, not %s", taskDone)
	}
	// TODO Придумать - как менять поля выборочно
	query, err := s.dbname.Prepare(`
	UPDATE task
	SET task_name = ?, task_description = ?, task_done = ?
	WHERE id = ?;
`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	result, err := query.Exec(taskName, taskDescription, taskDone, taskId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("task not found")
			return 0, storage.ErrTaskNotFound
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, _ := result.LastInsertId()
	return id, nil
}

func (s *Storage) GetTask(taskId string) ([]Task, error) {
	const op = "storage.sqlite.GetTask"
	query, err := s.dbname.Prepare(`
	SELECT task_name, task_description, task_done, task_created
	FROM task
	WHERE id = ?;
    `)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	var t Task
	err = query.QueryRow(taskId).Scan(&t.TaskName, &t.TaskDescription, &t.TaskDone, &t.TaskCreated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrTaskNotFound
		}
	}
	var result []Task
	result = append(result, t)
	return result, nil
}

func (s *Storage) GetTasks() ([]Task, error) {
	const op = "storage.sqlite.GetTasks"

	query, err := s.dbname.Prepare(`
	SELECT task_name, task_description, task_done, task_created
	FROM task
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	rows, err := query.Query()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrTaskNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)
	var got []Task
	for rows.Next() {
		var t Task
		err = rows.Scan(&t.TaskName, &t.TaskDescription, &t.TaskDone, &t.TaskCreated)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		got = append(got, t)

	}
	return got, nil
}
