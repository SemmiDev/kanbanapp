package entity

import "time"

type Category struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Type      string    `json:"type" gorm:"type:varchar(255);not null"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryRequest struct {
	Type string `json:"type" binding:"required"`
}

type CategoryData struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	Tasks []Task `json:"tasks"`
}

func DataToCategoryData(categories []Category, tasks []Task) []CategoryData {
	categoryData := make([]CategoryData, 0, len(categories))

	for _, category := range categories {
		tasksData := make([]Task, 0, len(tasks))

		for _, task := range tasks {
			if task.CategoryID == category.ID {
				tasksData = append(tasksData, task)
			}
		}

		categoryData = append(categoryData, CategoryData{
			ID:    category.ID,
			Type:  category.Type,
			Tasks: tasksData,
		})
	}

	return categoryData
}
