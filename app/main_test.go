package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetTasks(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)

	router := setupRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var tasks []task
	err := json.Unmarshal(w.Body.Bytes(), &tasks)
	assert.NoError(t, err)
}

func TestPostTask(t *testing.T) {
	newTask := task{
		TaskName:    "Test Task",
		Description: "This is a test task",
		DueDate:     time.Now().Add(24 * time.Hour),
	}
	taskJson, _ := json.Marshal(newTask)

	req, _ := http.NewRequest(http.MethodPost, "/task", bytes.NewBuffer(taskJson))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var createdTask task
	err := json.Unmarshal(w.Body.Bytes(), &createdTask)
	assert.NoError(t, err)
	assert.Equal(t, newTask.TaskName, createdTask.TaskName)
	assert.Equal(t, newTask.Description, createdTask.Description)
}

func TestGetTaskById(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/task/1", nil)
	w := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var specifiedTask task
	err := json.Unmarshal(w.Body.Bytes(), &specifiedTask)
	assert.NoError(t, err)
	assert.Equal(t, "1", specifiedTask.Id)
}

func TestUpdateTaskById(t *testing.T) {
	updatedTask := task{
		Id:                 "1",
		TaskName:           "Updated Task",
		Description:        "This task has been updated",
		Status:             "Completed",
		DueDate:            time.Now().Add(48 * time.Hour),
		CompletionDateTime: time.Now(),
		ModifiedDateTIme:   time.Now(),
	}

	taskJSON, _ := json.Marshal(updatedTask)

	req, _ := http.NewRequest(http.MethodPut, "/task/1", bytes.NewBuffer(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var updatedResponse task
	err := json.Unmarshal(w.Body.Bytes(), &updatedResponse)
	assert.NoError(t, err)
	assert.Equal(t, updatedTask.TaskName, updatedResponse.TaskName)
	assert.Equal(t, updatedTask.Description, updatedResponse.Description)
	assert.Equal(t, updatedTask.Status, updatedResponse.Status)
}

func TestDeleteTaskById(t *testing.T) {
	req, _ := http.NewRequest(http.MethodDelete, "/task/1", nil)
	w := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Succeed to delete task")
}
