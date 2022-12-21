package repository

import (
	"context"

	"github.com/SemmiDev/kanbanapp/internal/entity"

	"gorm.io/gorm"
)

type TaskRepository interface {
	GetTasks(ctx context.Context, id int) ([]entity.Task, error)
	StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error)
	GetTaskByID(ctx context.Context, id int) (entity.Task, error)
	GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error)
	UpdateTask(ctx context.Context, task *entity.Task) error
	DeleteTask(ctx context.Context, id int) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db}
}

func (r *taskRepository) GetTasks(ctx context.Context, id int) ([]entity.Task, error) {
	tasks := []entity.Task{}
	err := r.db.WithContext(ctx).Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepository) StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error) {
	err = r.db.WithContext(ctx).Create(task).Error
	if err != nil {
		return 0, err
	}
	return task.ID, nil
}

func (r *taskRepository) GetTaskByID(ctx context.Context, id int) (entity.Task, error) {
	task := entity.Task{}
	err := r.db.WithContext(ctx).First(&task, id).Error
	return task, err
}

func (r *taskRepository) GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error) {
	tasks := []entity.Task{}
	err := r.db.WithContext(ctx).Where("category_id = ?", catId).Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) UpdateTask(ctx context.Context, task *entity.Task) error {
	if task.CategoryID != 0 {
		err := r.db.WithContext(ctx).Model(&task).Update("category_id", task.CategoryID).Error
		return err
	}

	// update task title and description without category id
	err := r.db.WithContext(ctx).Model(&task).Updates(
		entity.Task{
			Title:       task.Title,
			Description: task.Description,
		},
	).Error
	return err
}

func (r *taskRepository) DeleteTask(ctx context.Context, id int) error {
	err := r.db.WithContext(ctx).Delete(&entity.Task{}, id).Error
	return err
}
