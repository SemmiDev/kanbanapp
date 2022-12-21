package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SemmiDev/kanbanapp/internal/entity"
	"github.com/SemmiDev/kanbanapp/internal/service"
)

type TaskAPI interface {
	GetTask(w http.ResponseWriter, r *http.Request)
	CreateNewTask(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
	UpdateTaskCategory(w http.ResponseWriter, r *http.Request)
}

type taskAPI struct {
	taskService service.TaskService
}

func NewTaskAPI(taskService service.TaskService) *taskAPI {
	return &taskAPI{taskService}
}

func (t *taskAPI) GetTask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("id")
	if userID == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	userIDInInt, err := strconv.Atoi(userID.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	taskID := r.URL.Query().Get("task_id")
	if taskID == "" {
		tasks, err := t.taskService.GetTasks(r.Context(), userIDInInt)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tasks)
		return
	}

	taskIDInInt, err := strconv.Atoi(taskID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task id"))
		return
	}

	task, err := t.taskService.GetTaskByID(r.Context(), taskIDInInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (t *taskAPI) CreateNewTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}

	if task.Title == "" || task.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}

	userID := r.Context().Value("id")
	if userID == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	userIDInInt, err := strconv.Atoi(userID.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	var newTask = entity.Task{
		Title:       task.Title,
		Description: task.Description,
		CategoryID:  task.CategoryID,
		UserID:      userIDInInt,
	}

	taskResult, err := t.taskService.StoreTask(r.Context(), &newTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	taskResponse := map[string]interface{}{
		"user_id": taskResult.UserID,
		"task_id": taskResult.ID,
		"message": "success create new task",
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(taskResponse)
}

func (t *taskAPI) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("id")
	if userID == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	userIDInInt, err := strconv.Atoi(userID.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	taskID := r.URL.Query().Get("task_id")
	if taskID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task id"))
		return
	}

	taskIDInInt, err := strconv.Atoi(taskID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task id"))
		return
	}

	err = t.taskService.DeleteTask(r.Context(), taskIDInInt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	taskResponse := map[string]interface{}{
		"user_id": userIDInInt,
		"task_id": taskIDInInt,
		"message": "success delete task",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(taskResponse)
}

func (t *taskAPI) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	userID := r.Context().Value("id")
	if userID == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	userIDInInt, err := strconv.Atoi(userID.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	taskID := r.URL.Query().Get("task_id")
	if taskID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task id"))
		return
	}

	taskIDInInt, err := strconv.Atoi(taskID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task id"))
		return
	}

	var newTask = entity.Task{
		ID:          taskIDInInt,
		Title:       task.Title,
		Description: task.Description,
		CategoryID:  task.CategoryID,
		UserID:      userIDInInt,
	}

	taskResult, err := t.taskService.UpdateTask(r.Context(), &newTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	taskResponse := map[string]interface{}{
		"user_id": taskResult.UserID,
		"task_id": taskResult.ID,
		"message": "success update task",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(taskResponse)
}

func (t *taskAPI) UpdateTaskCategory(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskCategoryRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	var updateTask = entity.Task{
		ID:         task.ID,
		CategoryID: task.CategoryID,
		UserID:     int(idLogin),
	}

	_, err = t.taskService.UpdateTask(r.Context(), &updateTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userId,
		"task_id": task.ID,
		"message": "success update task category",
	})
}
