package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"httpServer/internal/app"
	"httpServer/internal/database"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Task struct {
	ID         int
	Name       string `json:"name"`
	Details    string `json:"detail"`
	Completion bool   `json:"completion"`
}

// Function at add the incoming task into PSQL
func AddTask(appState *app.AppState, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)
	task := Task{}
	err := decoder.Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	_, err = appState.DB.CreateTask(r.Context(), database.CreateTaskParams{
		Title:  task.Name,
		Detail: sql.NullString{String: task.Details, Valid: task.Details != ""},
	})
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Task has been added successfully"))

}

// Function to send list of tasks to user
func ListTasks(appState *app.AppState, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	tasks, err := appState.DB.ListTasks(r.Context())
	if err != nil {
		fmt.Printf("%s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

// Delete the task
func DeleteTask(appState *app.AppState, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	fmt.Println("Extracted Vars:", vars) // Debugging step

	taskID, exists := vars["task_id"]
	if !exists || taskID == "" {
		http.Error(w, "Please enter task id", http.StatusBadRequest)
		return
	}

	ID, err := strconv.Atoi(taskID)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	err = appState.DB.DeleteTask(r.Context(), int32(ID))
	if err != nil {
		fmt.Printf("%s", err)
	}
	w.Write([]byte("Task has been deleted successfully"))
	w.WriteHeader(http.StatusOK)
}

// Update the task Completion
func UpdateTask(appState *app.AppState, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	taskID := mux.Vars(r)["task_id"]
	if taskID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please enter task id"))
	}
	ID, err := strconv.Atoi(taskID)
	if err != nil {
		fmt.Printf("%s", err)
	}
	err = appState.DB.UpdateTask(r.Context(), int32(ID))
	if err != nil {
		fmt.Printf("%s", err)
	}
	w.Write([]byte("Task has been updated"))
	w.WriteHeader(http.StatusOK)
}
