package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/igromanas/go-final/pkg/model"
)

func AddTask(task *model.Task) (string, error) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat);"
	res, err := DB.ExecContext(context.Background(), query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)

	if err != nil {
		return "", fmt.Errorf("error inserting row: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("error getting row id: %v", err)
	}

	return strconv.FormatInt(id, 10), nil
}

func Tasks(limit int) ([]*model.Task, error) {
	query := `SELECT * 
				FROM scheduler 
				ORDER BY date DESC 
				LIMIT :limit;`
	rows, err := DB.QueryContext(context.Background(), query, sql.Named("limit", limit))
	if err != nil {
		return []*model.Task{}, fmt.Errorf("tasks query execution error: %v", err)
	}
	defer rows.Close()

	tasks := []*model.Task{}
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return []*model.Task{}, fmt.Errorf("getting task row error: %v", err)
		}
		tasks = append(tasks, &t)
	}
	if err := rows.Err(); err != nil {
		return []*model.Task{}, fmt.Errorf("tasks rows processing error: %v", err)
	}

	return tasks, nil
}

func GetTaskByID(id string) (*model.Task, error) {
	query := `SELECT * 
				FROM scheduler 
				WHERE id = :id;`

	row := DB.QueryRowContext(context.Background(), query, sql.Named("id", id))
	task := &model.Task{}
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getting task row error: %v", err)
	}
	return task, nil
}

func UpdateTask(task *model.Task) error {
	query := `UPDATE scheduler 
				SET date = :date, 
					title = :title, 
					comment = :comment, 
					repeat = :repeat
				WHERE id = :id;`
	res, err := DB.ExecContext(context.Background(), query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID),
	)
	if err != nil {
		return fmt.Errorf("updating task error: %v", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("counting updated rows error: %v", err)
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}
	return nil
}

func UpdateDate(id string, date string) error {
	query := `UPDATE scheduler 
				SET date = :date
				WHERE id = :id;`
	res, err := DB.ExecContext(context.Background(), query, sql.Named("date", date), sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("updating task date error: %v", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("counting updated rows error: %v", err)
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating task date")
	}
	return nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = :id;`
	res, err := DB.ExecContext(context.Background(), query, sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("deleting task error")
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("counting deleted rows error: %v", err)
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for deleting task")
	}
	return nil
}
